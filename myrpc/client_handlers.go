package myrpc

import (
	"context"
	"fmt"
	"net"
	"time"
)

// TODO - Not Important ; Take in req and res(pointer) structure as Parameters

func (h *ClientCallHandler) CallGetPublicKey(ripaddr string, conn *net.UDPConn, clientHandlerName string) ([]byte, error) {
	var req PublicKeyRequest
	var res PublicKeyResponse

	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	if err := callWithContextAndConn(ctx, CLIENT_BASE_HANDLER_NAME+clientHandlerName+".GetPublicKey", req, &res, ripaddr, conn); err != nil {
		return nil, fmt.Errorf("Error while Calling %s.GetPublicKey RPC\nSource: CallGetPublicKey\nError: %v", CLIENT_BASE_HANDLER_NAME, err)
	}
	return res.PublicKey, nil
}

func (h *ClientCallHandler) CallInitNewWorkSpaceConnection(workspace_name, my_username, server_ip, workspace_password, ripaddr string, my_public_key []byte, udpConn *net.UDPConn, clientHandlerName string) (int, error) {
	var req InitWorkspaceConnectionRequest
	var res InitWorkspaceConnectionResponse

	req.WorkspaceName = workspace_name
	req.MyUsername = my_username
	req.MyPublicKey = my_public_key

	req.ServerIP = server_ip
	req.WorkspacePassword = workspace_password

	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	if err := callWithContextAndConn(ctx, CLIENT_BASE_HANDLER_NAME+clientHandlerName+".InitNewWorkSpaceConnection", req, &res, ripaddr, udpConn); err != nil {
		return 400, fmt.Errorf("Error while Calling InitNewWorkSpaceConnection RPC...\nError: %v", err)
	}

	return int(res.Response), nil
}

func (h *ClientCallHandler) CallGetData(myusername, server_ip, workspace_name, workspace_password, last_hash, ripaddr string, udpConn *net.UDPConn, clientHandlerName string) (*GetDataResponse, error) {
	var req GetDataRequest
	var res GetDataResponse

	req.Username = myusername
	req.WorkspaceName = workspace_name
	req.WorkspacePassword = workspace_password
	req.LastHash = last_hash
	req.ServerIP = server_ip

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()

	if err := callWithContextAndConn(ctx, CLIENT_BASE_HANDLER_NAME+clientHandlerName+".GetData", req, &res, ripaddr, udpConn); err != nil {
		res.Response = 400
		return &res, fmt.Errorf("Error in Calling RPC...\nError: %v", err)
	}

	return &res, nil
}
