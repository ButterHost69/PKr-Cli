package myrpc

import "fmt"

func (h *ClientCallHandler) CallGetPublicKey(ripaddr, lipaddr string) ([]byte, error) {
	var req PublicKeyRequest
	var res PublicKeyResponse

	if err := call(SERVER_HANDLER_NAME+".GetPublicKey", req, &res, ripaddr, lipaddr); err != nil {
		return []byte{}, fmt.Errorf("Error in Calling RPC...\nError: %v", err)
	}

	return res.PublicKey, nil
}

func (h *ClientCallHandler) CallInitNewWorkSpaceConnection(workspace_name, my_username, server_ip, workspace_password, ripaddr, lipaddr string, my_public_key []byte) (int, error){
	var req InitWorkspaceConnectionRequest
	var res InitWorkspaceConnectionResponse

	req.WorkspaceName = workspace_name
	req.MyUsername    = my_username
	req.MyPublicKey   = my_public_key

	req.ServerIP          = server_ip
	req.WorkspacePassword = workspace_password

	if err := call(SERVER_HANDLER_NAME+".GetPublicKey", req, &res, ripaddr, lipaddr); err != nil {
		return 400, fmt.Errorf("Error in Calling RPC...\nError: %v", err)
	}

	return int(res.Response), nil
}