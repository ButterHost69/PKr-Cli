package myrpc

import (
	"fmt"
	"strings"

	baseDialer "github.com/ButterHost69/PKr-Base/dialer"
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

	// port := utils.GetRandomPort()
	// my_pub_ip, err := baseDialer.GetMyPublicIP(port)
	// if err != nil {
	// 	return "", fmt.Errorf("cannot get public ip\nSource: CallRegisterUser")
	// }

	// ip_split := strings.Split(my_pub_ip, ":")
	// req.PublicIP = ip_split[0]
	// req.PublicPort = ip_split[1]

	if err := call(SERVER_HANDLER_NAME+".RegisterUser", req, &res, server_ip, ""); err != nil {
		return "", fmt.Errorf("Error in Calling RPC...\nError: %v", err)
	}

	if res.Response != 200 {
		return "", fmt.Errorf("Calling CallRegisterUser Method was not Successful.\nReturn Code - %d", res.Response)
	}

	return res.UniqueUsername, nil
}

func (h *ServerCallHandler) CallRequestPunchFromReciever(server_ip, recieverUsername, username, password string, port int) (string, error) {
	var req RequestPunchFromRecieverRequest
	var res RequestPunchFromRecieverResponse

	ipaddr, err := baseDialer.GetMyPublicIP(port)
	if err != nil {
		return "", err
	}

	ip_split := strings.Split(ipaddr, ":")

	req.Username = username
	req.Password = password
	req.RecieversUsername = recieverUsername
	req.SendersIP = ip_split[0]
	req.SendersPort = ip_split[1]

	if err := call(SERVER_HANDLER_NAME+".RequestPunchFromReciever", req, &res, server_ip, ""); err != nil {
		return "", fmt.Errorf("Error in Calling RPC...\nError: %v", err)
	}

	if res.Response != 200 {
		return "", fmt.Errorf("Calling CallPunchFromReciever Method was not Successful.\nReturn Code - %d", res.Response)
	}

	return res.RecieversPublicIP, err
}
