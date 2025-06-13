package dialer

import (
	"fmt"
	"net/rpc"
)

func (h *ClientCallHandler) CallGetPublicKey(clientHandlerName string, rpc_client *rpc.Client) ([]byte, error) {
	var req PublicKeyRequest
	var res PublicKeyResponse

	rpc_name := CLIENT_BASE_HANDLER_NAME + clientHandlerName + ".GetPublicKey"
	if err := CallKCP_RPC_WithContext(req, &res, rpc_name, rpc_client); err != nil {
		fmt.Println("Error while Calling GetPublicKey:", err)
		fmt.Println("Source: CallGetPublicKey()")
		return nil, err
	}
	return res.PublicKey, nil
}

func (h *ClientCallHandler) CallInitNewWorkSpaceConnection(workspace_name, my_username, server_ip, workspace_password string, my_public_key []byte, clientHandlerName string, rpc_client *rpc.Client) error {
	var req InitWorkspaceConnectionRequest
	var res InitWorkspaceConnectionResponse

	req.WorkspaceName = workspace_name
	req.MyUsername = my_username
	req.MyPublicKey = my_public_key

	req.ServerIP = server_ip
	req.WorkspacePassword = workspace_password

	rpc_name := CLIENT_BASE_HANDLER_NAME + clientHandlerName + ".InitNewWorkSpaceConnection"
	if err := CallKCP_RPC_WithContext(req, &res, rpc_name, rpc_client); err != nil {
		fmt.Println("Error while Calling Init New Workspace Connection:", err)
		fmt.Println("Source: CallInitNewWorkSpaceConnection()")
		return err
	}
	return nil
}

func (h *ClientCallHandler) CallGetData(my_username, server_ip, workspace_name, workspace_password, last_hash, clientHandlerName string, rpc_client *rpc.Client) (*GetDataResponse, error) {
	var req GetDataRequest
	var res GetDataResponse

	req.Username = my_username
	req.WorkspaceName = workspace_name
	req.WorkspacePassword = workspace_password
	req.LastHash = last_hash
	req.ServerIP = server_ip

	rpc_name := CLIENT_BASE_HANDLER_NAME + clientHandlerName + ".GetData"
	if err := CallKCP_RPC_WithContext(req, &res, rpc_name, rpc_client); err != nil {
		fmt.Println("Error while Calling Get Data:", err)
		fmt.Println("Source: CallGetData()")
		return nil, err
	}
	return &res, nil
}
