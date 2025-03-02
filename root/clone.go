package root

import (
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/ButterHost69/PKr-Base/config"
	"github.com/ButterHost69/PKr-cli/dialer"
	"github.com/ButterHost69/PKr-cli/encrypt"
	"github.com/ButterHost69/PKr-cli/myrpc"
	"github.com/ButterHost69/PKr-cli/utils"

	baseDialer "github.com/ButterHost69/PKr-Base/dialer"
)

func Clone(workspace_owner_username, workspace_name, workspace_password, server_alias string) error {
	// [X] Get Public Key From the Host Original Source PC
	// [X]  Encrypt Password
	// [X]  Read Our Key
	// [X]  Send Password and request for InitConnection -> return port
	// [X]  Register The import Folder
	// [X]  Connect to DataServer from the Port
	// [X]  Decrypt the file
	// [X]  Unzip the File

	// [X] Get and Encrypt Key

	server, err := config.GetServerDetailsUsingServerAlias(server_alias)
	if err != nil {
		return fmt.Errorf("Error Occured while Fetching Server Details Using Server Alias\nSource: Clone\nError:%v", err)
	}

	localPort := utils.GetRandomPort()
	log.Println("My Local Port:", localPort)

	myPublicIP, err := baseDialer.GetMyPublicIP(localPort)
	if err != nil {
		return fmt.Errorf("Error Occured while Getting My Public IP\nSource: Clone\nError:%v", err)
	}
	log.Println("My Public IP Addr:", myPublicIP)

	privateIPStr := ":" + strconv.Itoa(localPort)
	privateIP, err := net.ResolveUDPAddr("udp", privateIPStr)
	if err != nil {
		return fmt.Errorf("Error Occured while Resolving Private UDP Addr\nSource: Clone\nError:%v", err)
	}

	udpConn, err := net.ListenUDP("udp", privateIP)
	if err != nil {
		return fmt.Errorf("Error Occured while Listening to UDP\nSource: Clone\nError:%v", err)
	}

	log.Println("MOIT Calling Request Punch From Receiver ...")

	serverClient := myrpc.ServerCallHandler{}
	workspace_owner_ip, err := serverClient.CallRequestPunchFromReciever(server.ServerIP, workspace_owner_username, server.Username, server.Password, myPublicIP)
	if err != nil {
		return fmt.Errorf("Error Occured while Calling Request Punch From Reciever\nSource: Clone\nError:%v", err)
	}
	log.Println("Receivers IP:", workspace_owner_ip)

	err = dialer.UdpNatPunching(udpConn, workspace_owner_ip)
	if err != nil {
		return fmt.Errorf("Error Occured while Performing NAT Hole Punching\nSource: Clone\nError:%v")
	}
	fmt.Println("Punched Successfully ...")
	rpcClientHandler := myrpc.ClientCallHandler{}

	log.Println("Requesting Public Key ...")
	public_key, err := dialer.RequestPublicKey(workspace_owner_ip, udpConn, rpcClientHandler)
	if err != nil {
		return fmt.Errorf("error Occured in Retrieving Public Key.\nerror:%v", err)
	}
	fmt.Println("Retrieved Public Key From the Source PC\nNow Encrypting Data ...")

	encrypted_password, err := encrypt.EncryptData(workspace_password, string(public_key))
	if err != nil {
		return fmt.Errorf("error Occured in Encrypting Password.\nerror: %v", err)
	}

	fmt.Println("Data Encrypted ...\nReading my Public Key ...")

	my_public_key, err := os.ReadFile("./tmp/mykeys/publickey.pem")
	if err != nil {
		return fmt.Errorf("error Occured in Reading Our Public Key\nPlease Ensure Key is Present at ./tmp/mykeys/publickey.pem\nerror:%v", err)
	}

	fmt.Println("My Public Key is Retrieved ...")
	base64_public_key := []byte(base64.StdEncoding.EncodeToString(my_public_key))

	fmt.Println("Request Init New Work Space Connection START")
	response, err := dialer.RequestInitNewWorkSpaceConnection(server.ServerIP, server.Username, server.ServerIP, workspace_owner_username, workspace_name, encrypted_password, base64_public_key, udpConn, workspace_owner_ip, rpcClientHandler)
	if err != nil {
		return fmt.Errorf("error Occured in Dialing Init New Workspace Connection.\nerror:%v", err)
	}
	fmt.Println("Request Init New Work Space Connection END ith Response Code:", response)
	if response != 200 {
		return fmt.Errorf("could not init new workspace")
	}

	currDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error in Retrieving Current Working Directory.\nerror:%v", err)
	}
	err = os.MkdirAll(currDir+"\\.PKr\\", 0777)
	if err != nil {
		return fmt.Errorf("error Occured in Creating .PKr Folder.\nerror: %v", err)
	}

	fmt.Println("Initialized Workspace With the Source PC END\nSending Request of Get Data ...")

	// TODO Request to GetData() Separately... Pass null string as last hash
	res, err := dialer.RequestGetData(workspace_owner_ip, server.Username, server.Password, workspace_name, workspace_password, "", server.ServerIP, udpConn, rpcClientHandler)
	if err != nil {
		return err
	}

	if res != 200 {
		return fmt.Errorf("could not get data, response code - %d", res)
	}

	fmt.Println("Request Get Data Done ...")
	return nil
}
