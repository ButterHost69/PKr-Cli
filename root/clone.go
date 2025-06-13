package root

import (
	"context"
	"encoding/base64"
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
	"github.com/ButterHost69/PKr-Cli/pb"

	"github.com/ButterHost69/kcp-go"
)

func connectToAnotherUser(workspace_owner_username, server_ip, username, password string) (string, *kcp.UDPSession, error) {
	local_port := rand.Intn(16384) + 16384
	fmt.Println("My Local Port:", local_port)

	// Get My Public IP
	myPublicIP, err := dialer.GetMyPublicIP(local_port)
	if err != nil {
		fmt.Println("Error while Getting my Public IP:", err)
		fmt.Println("Source: connectToAnotherUser()")
		return "", nil, err
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
		return "", nil, err
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
		return "", nil, err

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
		return "", nil, err
	}

	client_handler_name, err := dialer.WorkspaceListenerUdpNatHolePunching(udp_conn, res.WorkspaceOwnerPublicIp+":"+res.WorkspaceOwnerPublicPort)
	if err != nil {
		fmt.Println("Error while Punching to Remote Addr:", err)
		fmt.Println("Source: connectToAnotherUser()")
		return "", nil, err

	}
	fmt.Println("UDP NAT Hole Punching Completed Successfully")

	// Creating KCP-Conn, KCP = Reliable UDP
	kcp_conn, err := kcp.DialWithConnAndOptions(res.WorkspaceOwnerPublicIp+":"+res.WorkspaceOwnerPublicPort, nil, 0, 0, udp_conn)
	if err != nil {
		fmt.Println("Error while Dialing KCP Connection to Remote Addr:", err)
		fmt.Println("Source: connectToAnotherUser()")
		return "", nil, err
	}

	// KCP Params for Congestion Control
	kcp_conn.SetWindowSize(128, 512)
	kcp_conn.SetNoDelay(1, 20, 0, 1)
	kcp_conn.SetACKNoDelay(false)

	return client_handler_name, kcp_conn, nil
}

func storeDataIntoWorkspace(res *dialer.GetDataResponse) error {
	data_bytes := res.Data
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

	workspacePath := "."

	zip_file_path := workspacePath + "\\.PKr\\" + res.NewHash + ".zip"
	if err = filetracker.SaveDataToFile(data, zip_file_path); err != nil {
		fmt.Println("Error while Saving Data into '.PKr/abc.zip':", err)
		fmt.Println("Source: storeDataIntoWorkspace()")
		return err
	}

	if err = filetracker.CleanFilesFromWorkspace(workspacePath); err != nil {
		fmt.Println("Error while Cleaning Workspace :", err)
		fmt.Println("Source: storeDataIntoWorkspace()")
		return err
	}

	// Unzip Content
	if err = filetracker.UnzipData(zip_file_path, workspacePath+"\\"); err != nil {
		fmt.Println("Error while Unzipping Data into Workspace:", err)
		fmt.Println("Source: storeDataIntoWorkspace()")
		return err
	}
	return nil
}

func Clone(workspace_owner_username, workspace_name, workspace_password, server_alias string) {
	// Get Details from Config
	server_ip, username, password, err := config.GetServerDetails(server_alias)
	if err != nil {
		fmt.Println("Error while getting Server Details from Config:", err)
		fmt.Println("Source: Clone()")
		return
	}

	client_handler_name, kcp_conn, err := connectToAnotherUser(workspace_owner_username, server_ip, username, password)
	if err != nil {
		fmt.Println("Error while Connecting to Another User:", err)
		fmt.Println("Source: Clone()")
		return
	}
	defer kcp_conn.Close()

	// Creating RPC Client
	rpc_client := rpc.NewClient(kcp_conn)
	defer rpc_client.Close()

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

	fmt.Println("Calling GetData ...")
	// Calling GetData
	res, err := rpcClientHandler.CallGetData(username, server_ip, workspace_name, encrypted_password, "", client_handler_name, rpc_client)
	if err != nil {
		fmt.Println("Error while Calling GetData:", err)
		fmt.Println("Source: Clone()")
		return
	}

	fmt.Println("Get Data Responded, now storing files into workspace")
	// Store Data into workspace
	err = storeDataIntoWorkspace(res)
	if err != nil {
		fmt.Println("Error while Storing Requested Data into Workspace:", err)
		fmt.Println("Source: Clone()")
		return
	}

	// Register New User to Workspace
	// Add New Workspace Connection
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
		fmt.Println("Error while Requesting Punch from Receiver:", err)
		fmt.Println("Source: Clone()")
		return
	}

	// Update tmp/userConfig.json
	err = config.RegisterNewGetWorkspace(server_alias, workspace_name, currDir, workspace_password, res.NewHash)
	if err != nil {
		fmt.Println("Error while Registering New GetWorkspace:", err)
		fmt.Println("Source: Clone()")
		return
	}
	fmt.Println("Clone Done")
}
