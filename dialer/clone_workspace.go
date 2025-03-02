package dialer

import (
	// "errors"
	"fmt"
	"net"

	"github.com/ButterHost69/PKr-cli/myrpc"
)

// Sender Side - CLI
// TODO 1. Request Reciever (of the request) to Punch - Server
// TODO 2. Send Get Public Key Request - Client
// TODO 3. Encrpyt and Send Password - Client
// TODO 4. Send Password and MY_Public_Key and Request for new Init Connection - Client
// TODO 5. Register as Get Workspace - Locally
// TODO 6. Get File
// TODO 7. Decrypt and store

// Reciever Side = Base
// TODO 1. Recieve Request from server to punch - Server
// TODO 2. Puch and respond with allocated port to server - Server
// TODO 3. Auth Connection
// TODO 4. Register Connection to Send Workspace and Store ClientsKey
// TODO 5. Encrypt and Send File

// MAYBE - Treat Get Data as Seperate Service ...

func RequestPublicKey(receivers_ip string, udpConn *net.UDPConn, rpcClientHandler myrpc.ClientCallHandler) ([]byte, error) {
	public_key, err := rpcClientHandler.CallGetPublicKey(receivers_ip, udpConn)
	if err != nil {
		fmt.Println("Error while Calling Get Public Key\nSource: RequestPublicKey\nError:", err)
		return nil, err
	}
	return public_key, nil
}

func RequestInitNewWorkSpaceConnection(server_ip, my_username, my_password, rcv_username, workspace_name, workspace_password string, public_key []byte, udpConn *net.UDPConn, workspace_owner_ip string, rpcClientHandler myrpc.ClientCallHandler) (int, error) {
	response, err := rpcClientHandler.CallInitNewWorkSpaceConnection(workspace_name, my_username, server_ip, workspace_password, workspace_owner_ip, public_key, udpConn)
	if err != nil {
		fmt.Println("Error while Requesting Init New Workspace Connection\nSource: RequestInitNewWorkSpaceConnection\nError:", err)
		return -1, err
	}
	return response, err
}
