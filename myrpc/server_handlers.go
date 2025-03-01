package myrpc

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"
)

func (h *ServerCallHandler) CallRegisterWorkspace(server_ip, username, password, workspace_name string) error {

	var req RegisterWorkspaceRequest
	var res RegisterUserResponse

	req.Username = username
	req.Password = password
	req.WorkspaceName = workspace_name

	if err := call(SERVER_HANDLER_NAME+".RegisterWorkspace", req, &res, server_ip, ""); err != nil {
		return fmt.Errorf("Error in Calling RPC...\nError: %v", err)
	}

	if res.Response != 200 {
		return fmt.Errorf("Calling Ping Method was not Successful.\nReturn Code - %d", res.Response)
	}

	return nil
}

func (h *ServerCallHandler) CallRegisterUser(server_ip, username, password string) (string, error) {
	var req RegisterUserRequest
	var res RegisterUserResponse

	req.Username = username
	req.Password = password

	if err := call(SERVER_HANDLER_NAME+".RegisterUser", req, &res, server_ip, ""); err != nil {
		return "", fmt.Errorf("Error in Calling RPC...\nError: %v", err)
	}

	if res.Response != 200 {
		return "", fmt.Errorf("Calling CallRegisterUser Method was not Successful.\nReturn Code - %d", res.Response)
	}

	return res.UniqueUsername, nil
}

func (h *ServerCallHandler) CallRequestPunchFromReciever(server_ip, recieverUsername, username, password string, myPublicIP string, conn *net.UDPConn) (string, error) {
	var req RequestPunchFromRecieverRequest
	var res RequestPunchFromRecieverResponse

	req.Username = username
	req.Password = password
	req.RecieversUsername = recieverUsername

	ip_split := strings.Split(myPublicIP, ":")
	req.SendersIP = ip_split[0]
	req.SendersPort = ip_split[1]

	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	if err := callWithContextAndConn(ctx, SERVER_HANDLER_NAME+".RequestPunchFromReciever", req, &res, server_ip, conn); err != nil {
		return "", fmt.Errorf("Error while calling %s.RequestPunchFromReciever RPC...\nSource: CallRequestPunchFromReciever\nError: %v", SERVER_HANDLER_NAME, err)
	}

	if res.Response != 200 {
		return "", fmt.Errorf("Calling CallPunchFromReciever Method was not Successful.\nReturn Code - %d", res.Response)
	}

	ip_addr := fmt.Sprintf("%s:%d", res.RecieversPublicIP, res.RecieversPublicPort)
	return ip_addr, nil
}
