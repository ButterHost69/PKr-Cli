package dialer

import (
	// "errors"
	"fmt"

	"github.com/ButterHost69/PKr-cli/myrpc"
	"github.com/ButterHost69/PKr-cli/utils"
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

func RequestPublicKey(rcv_username, server_ip, my_username, my_password string) ([]byte, error) {
	callHandler := myrpc.ServerCallHandler{}

	port := utils.GetRandomPort()
	receivers_ip, err := callHandler.CallPunchFromReciever(server_ip, rcv_username, my_username, my_password, port)
	if err != nil {
		return nil, err
	}


	handler := myrpc.ClientCallHandler{}
	public_key, err := handler.CallGetPublicKey(receivers_ip, fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	return public_key, nil
}

func RequestInitNewWorkSpaceConnection(server_ip, my_username, my_password, rcv_username, workspace_name, workspace_password string, public_key []byte) (int, error) {
	port := utils.GetRandomPort()

	callHandler := myrpc.ServerCallHandler{}
	receivers_ip, err := callHandler.CallPunchFromReciever(server_ip, rcv_username, my_username, my_password, port)

	handler := myrpc.ClientCallHandler{}
	lipaddr := fmt.Sprintf(":%d", port)
	response, err := handler.CallInitNewWorkSpaceConnection(workspace_name, my_username, server_ip, workspace_password, receivers_ip, lipaddr, public_key)
	return response, err
}