package root

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"net/rpc"
	"os"
	"strings"

	"github.com/ButterHost69/PKr-Cli/config"
	"github.com/ButterHost69/PKr-Cli/dialer"
	"github.com/ButterHost69/PKr-Cli/encrypt"
	"github.com/ButterHost69/PKr-Cli/filetracker"
	"github.com/ButterHost69/PKr-Cli/models"
	"github.com/ButterHost69/PKr-Cli/pb"
	"github.com/ButterHost69/PKr-Cli/utils"

	"github.com/ButterHost69/kcp-go"
)

const DATA_CHUNK = 1024 // 1 KB

func connectToAnotherUser(workspace_owner_username, server_ip, username, password string) (string, string, *net.UDPConn, *kcp.UDPSession, error) {
	local_port := rand.Intn(16384) + 16384
	fmt.Println("My Local Port:", local_port)

	// Get My Public IP
	myPublicIP, err := dialer.GetMyPublicIP(local_port)
	if err != nil {
		fmt.Println("Error while Getting my Public IP:", err)
		fmt.Println("Source: connectToAnotherUser()")
		return "", "", nil, nil, err
	}
	fmt.Println("My Public IP Addr:", myPublicIP)

	myPublicIPSplit := strings.Split(myPublicIP, ":")
	myPublicIPOnly := myPublicIPSplit[0]
	myPublicPortOnly := myPublicIPSplit[1]

	// New GRPC Client
	gRPC_cli_service_client, err := dialer.NewGRPCClients(server_ip)
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println("Description: Cannot Create New GRPC Client")
		fmt.Println("Source: connectToAnotherUser()")
		return "", "", nil, nil, err
	}

	// Prepare req
	req := &pb.RequestPunchFromReceiverRequest{
		WorkspaceOwnerUsername: workspace_owner_username,
		ListenerUsername:       username,
		ListenerPassword:       password,
		ListenerPublicIp:       myPublicIPOnly,
		ListenerPublicPort:     myPublicPortOnly,
	}

	// Request Timeout
	ctx, cancelFunc := context.WithTimeout(context.Background(), CONTEXT_TIMEOUT)
	defer cancelFunc()

	// Sending Request to Server
	res, err := gRPC_cli_service_client.RequestPunchFromReceiver(ctx, req)
	if err != nil {
		fmt.Println("Error while Requesting Punch from Receiver:", err)
		fmt.Println("Source: connectToAnotherUser()")
		return "", "", nil, nil, err

	}
	fmt.Println("Remote Addr:", res.WorkspaceOwnerPublicIp+":"+res.WorkspaceOwnerPublicPort)

	// Creating UDP Conn to Perform UDP NAT Hole Punching
	udp_conn, err := net.ListenUDP("udp", &net.UDPAddr{
		Port: local_port,
		IP:   net.IPv4zero, // or nil
	})
	if err != nil {
		fmt.Printf("Error while Listening to %d: %v\n", local_port, err)
		fmt.Println("Source: connectToAnotherUser()")
		return "", "", nil, nil, err
	}

	workspace_owner_public_ip := res.WorkspaceOwnerPublicIp + ":" + res.WorkspaceOwnerPublicPort
	client_handler_name, err := dialer.WorkspaceListenerUdpNatHolePunching(udp_conn, workspace_owner_public_ip)
	if err != nil {
		fmt.Println("Error while Punching to Remote Addr:", err)
		fmt.Println("Source: connectToAnotherUser()")
		return "", "", nil, nil, err

	}
	fmt.Println("UDP NAT Hole Punching Completed Successfully")

	// Creating KCP-Conn, KCP = Reliable UDP
	kcp_conn, err := kcp.DialWithConnAndOptions(workspace_owner_public_ip, nil, 0, 0, udp_conn)
	if err != nil {
		fmt.Println("Error while Dialing KCP Connection to Remote Addr:", err)
		fmt.Println("Source: connectToAnotherUser()")
		return "", "", nil, nil, err
	}

	// KCP Params for Congestion Control
	kcp_conn.SetWindowSize(128, 512)
	kcp_conn.SetNoDelay(1, 10, 1, 1)

	return client_handler_name, workspace_owner_public_ip, udp_conn, kcp_conn, nil
}

// TODO: Instead of writing whole data_bytes into file at once,
// Write received encrypted data in chunks, after the transfer is completed, read from encrpyted file
// & decrypt it
// We can use Cipher Block Methods to decrypt & encrpyt with AES
func storeDataIntoWorkspace(res *models.GetMetaDataResponse, data_bytes []byte) error {
	key_bytes := res.KeyBytes
	iv_bytes := res.IVBytes

	decrypted_key, err := encrypt.DecryptData(string(key_bytes))
	if err != nil {
		fmt.Println("Error while Decrypting Key:", err)
		fmt.Println("Source: storeDataIntoWorkspace()")
		return err
	}

	decrypted_iv, err := encrypt.DecryptData(string(iv_bytes))
	if err != nil {
		fmt.Println("Error while Decrypting 'IV':", err)
		fmt.Println("Source: storeDataIntoWorkspace()")
		return err
	}

	data, err := encrypt.AESDecrypt(data_bytes, decrypted_key, decrypted_iv)
	if err != nil {
		fmt.Println("Error while Decrypting Data:", err)
		fmt.Println("Source: storeDataIntoWorkspace()")
		return err
	}

	// README: For PKr-Base
	// fmt.Println("Workspace Name:", workspace_name)
	// workspace_path, err := config.GetGetWorkspaceFilePath(workspace_name)
	// if err != nil {
	// 	fmt.Println("Error while Fetching Workspace Path from Config:", err)
	// 	fmt.Println("Source: storeDataIntoWorkspace()")
	// 	return err
	// }
	// fmt.Println("Workspace Path: ", workspace_path)

	workspace_path := "."

	zip_file_path := workspace_path + "\\.PKr\\" + res.NewHash + ".zip"
	if err = filetracker.SaveDataToFile(data, zip_file_path); err != nil {
		fmt.Println("Error while Saving Data into '.PKr/abc.zip':", err)
		fmt.Println("Source: storeDataIntoWorkspace()")
		return err
	}

	if err = filetracker.CleanFilesFromWorkspace(workspace_path); err != nil {
		fmt.Println("Error while Cleaning Workspace :", err)
		fmt.Println("Source: storeDataIntoWorkspace()")
		return err
	}

	// Unzip Content
	if err = filetracker.UnzipData(zip_file_path, workspace_path+"\\"); err != nil {
		fmt.Println("Error while Unzipping Data into Workspace:", err)
		fmt.Println("Source: storeDataIntoWorkspace()")
		return err
	}
	return nil
}

func fetchData(workspace_owner_public_ip, workspace_name, workspace_hash string, udp_conn *net.UDPConn, len_data_bytes int) ([]byte, error) {
	// Now Transfer Data using KCP ONLY, No RPC in chunks
	fmt.Println("Connecting Again to Workspace Owner")
	kcp_conn, err := kcp.DialWithConnAndOptions(workspace_owner_public_ip, nil, 0, 0, udp_conn)
	if err != nil {
		fmt.Println("Error while Dialing Workspace Owner to Get Data:", err)
		fmt.Println("Source: fetchData()")
		return nil, err
	}
	fmt.Println("Connected Successfully to Workspace Owner")

	// KCP Params for Congestion Control
	kcp_conn.SetWindowSize(128, 512)
	kcp_conn.SetNoDelay(1, 10, 1, 1)

	// Sending the Type of Session
	kpc_buff := [3]byte{'K', 'C', 'P'}
	_, err = kcp_conn.Write(kpc_buff[:])
	if err != nil {
		fmt.Println("Error while Writing the type of Session(KCP-RPC or KCP-Plain):", err)
		fmt.Println("Source: fetchData()")
		return nil, err
	}

	fmt.Println("Sending Workspace Name & Hash to Workspace Owner")
	// Sending Workspace Name & Hash
	_, err = kcp_conn.Write([]byte(workspace_name))
	if err != nil {
		fmt.Println("Error while Sending Workspace Name to Workspace Owner:", err)
		fmt.Println("Source: fetchData()")
		return nil, err
	}

	_, err = kcp_conn.Write([]byte(workspace_hash))
	if err != nil {
		fmt.Println("Error while Sending Workspace Name to Workspace Owner:", err)
		fmt.Println("Source: fetchData()")
		return nil, err
	}
	fmt.Println("Workspace Name & Hash Sent to Workspace Owner")

	CHUNK_SIZE := min(DATA_CHUNK, len_data_bytes)

	fmt.Println("Len Data Bytes:", len_data_bytes)
	fmt.Println("Len Buffer:", len_data_bytes+CHUNK_SIZE)
	data_bytes := make([]byte, len_data_bytes+CHUNK_SIZE)
	offset := 0

	fmt.Println("Now Reading Data from Workspace Owner ...")
	for offset < len_data_bytes {

		n, err := kcp_conn.Read(data_bytes[offset : offset+CHUNK_SIZE])
		// Check for Errors on Workspace Owner's Side
		if n < 30 {
			msg := string(data_bytes[offset : offset+n])
			if msg == "Incorrect Workspace Name/Hash" || msg == "Internal Server Error" {
				fmt.Println("\nError while Reading from Workspace on his/her side:", msg)
				fmt.Println("Source: fetchData()")
				return nil, errors.New(msg)
			}
		}

		if err != nil {
			fmt.Println("\nError while Reading from Workspace Owner:", err)
			fmt.Println("Source: fetchData()")
			return nil, err
		}
		offset += n
		utils.PrintProgressBar(offset, len_data_bytes, 100)
	}
	fmt.Println("\nData Transfer Completed ...")

	_, err = kcp_conn.Write([]byte("Data Received"))
	if err != nil {
		fmt.Println("Error while Sending Data Received Message:", err)
		fmt.Println("Source: fetchData()")
		// Not Returning Error because, we got data, we don't care if workspace owner now is offline
	}
	return data_bytes[:offset], nil
}

func Clone(workspace_owner_username, workspace_name, workspace_password, server_alias string) {
	// Get Details from Config
	server_ip, username, password, err := config.GetServerDetails(server_alias)
	if err != nil {
		fmt.Println("Error while getting Server Details from Config:", err)
		fmt.Println("Source: Clone()")
		return
	}

	// Connecting to Workspace Owner
	client_handler_name, workspace_owner_public_ip, udp_conn, kcp_conn, err := connectToAnotherUser(workspace_owner_username, server_ip, username, password)
	if err != nil {
		fmt.Println("Error while Connecting to Another User:", err)
		fmt.Println("Source: Clone()")
		return
	}

	// Sending the Type of Session
	rpc_buff := [3]byte{'R', 'P', 'C'}
	_, err = kcp_conn.Write(rpc_buff[:])
	if err != nil {
		fmt.Println("Error while Writing the type of Session(KCP-RPC or KCP-Plain):", err)
		fmt.Println("Source: cloneWorkspace()")
		return
	}

	// Creating RPC Client
	rpc_client := rpc.NewClient(kcp_conn)
	rpcClientHandler := dialer.ClientCallHandler{}

	fmt.Println("Calling Get Public Key")
	// Get Public Key of Workspace Owner
	public_key, err := rpcClientHandler.CallGetPublicKey(client_handler_name, rpc_client)
	if err != nil {
		fmt.Println("Error while Calling GetPublicKey:", err)
		fmt.Println("Source: Clone()")
		return
	}

	// Encrypting Workspace Password with Public Key
	encrypted_password, err := encrypt.EncryptData(workspace_password, string(public_key))
	if err != nil {
		fmt.Println("Error while Encrypting Workspace Password via Public Key:", err)
		fmt.Println("Source: Clone()")
		return
	}

	// Reading my Public Key
	my_public_key, err := os.ReadFile("./tmp/mykeys/publickey.pem")
	if err != nil {
		fmt.Println("Error while Reading Public Key:", err)
		fmt.Println("Source: Clone()")
		return
	}
	base64_public_key := []byte(base64.StdEncoding.EncodeToString(my_public_key))

	fmt.Println("Calling InitWorkspaceConnection")
	// Requesting InitWorkspaceConnection
	err = rpcClientHandler.CallInitNewWorkSpaceConnection(workspace_name, username, server_ip, encrypted_password, base64_public_key, client_handler_name, rpc_client)
	if err != nil {
		fmt.Println("Error while Calling Init New Workspace Connection:", err)
		fmt.Println("Source: Clone()")
		return
	}

	// Create .PKr folder to store zipped data
	currDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error while Getting Current Directory:", err)
		fmt.Println("Source: Clone()")
		return
	}
	err = os.MkdirAll(currDir+"\\.PKr\\", 0777)
	if err != nil {
		fmt.Println("Error while using MkdirAll for '.PKr' folder:", err)
		fmt.Println("Source: Clone()")
		return
	}

	fmt.Println("Calling GetMetaData ...")
	// Calling GetMetaData
	res, err := rpcClientHandler.CallGetMetaData(username, server_ip, workspace_name, encrypted_password, "", client_handler_name, rpc_client)
	if err != nil {
		fmt.Println("Error while Calling GetMetaData:", err)
		fmt.Println("Source: Clone()")
		return
	}

	fmt.Println("Get Meta Data Responded")
	rpc_client.Close()
	defer kcp_conn.Close()

	data_bytes, err := fetchData(workspace_owner_public_ip, workspace_name, res.NewHash, udp_conn, res.LenData)
	if err != nil {
		fmt.Println("Error while Fetching Data:", err)
		fmt.Println("Source: Clone()")
		return
	}

	fmt.Println("Now Storing Data into Workspace ...")
	// Store Data into workspace
	err = storeDataIntoWorkspace(res, data_bytes)
	if err != nil {
		fmt.Println("Error while Storing Requested Data into Workspace:", err)
		fmt.Println("Source: Clone()")
		return
	}
	fmt.Println("Data Stored into Workspace")

	// Update tmp/userConfig.json
	err = config.RegisterNewGetWorkspace(server_alias, workspace_name, workspace_owner_username, currDir, workspace_password, res.NewHash)
	if err != nil {
		fmt.Println("Error while Registering New GetWorkspace:", err)
		fmt.Println("Source: Clone()")
		return
	}

	// Register New User to Workspace
	// New GRPC Client
	gRPC_cli_service_client, err := dialer.NewGRPCClients(server_ip)
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println("Description: Cannot Create New GRPC Client")
		fmt.Println("Source: Clone()")
		return
	}

	// Prepare req
	register_user_to_workspace_res_req := &pb.RegisterUserToWorkspaceRequest{
		ListenerUsername:       username,
		ListenerPassword:       password,
		WorkspaceName:          workspace_name,
		WorkspaceOwnerUsername: workspace_owner_username,
	}

	// Request Timeout
	ctx, cancelFunc := context.WithTimeout(context.Background(), CONTEXT_TIMEOUT)
	defer cancelFunc()

	// Sending Request to Server
	_, err = gRPC_cli_service_client.RegisterUserToWorkspace(ctx, register_user_to_workspace_res_req)
	if err != nil {
		fmt.Println("Error while Registering User To Workspace:", err)
		fmt.Println("Source: Clone()")
		return
	}

	fmt.Println("Clone Done")
}
