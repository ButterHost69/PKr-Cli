package root

import (
	"context"
	"fmt"
	"time"

	"github.com/ButterHost69/PKr-Base/config"
	"github.com/ButterHost69/PKr-Cli/dialer"
	"github.com/ButterHost69/PKr-Cli/pb"
)

const CONTEXT_TIMEOUT = 60 * time.Second

func Install(server_alias, server_ip, username, password string) {
	config.CreateUserIfNotExists()

	if server_alias == "" || server_ip == "" || username == "" || password == "" {
		fmt.Println("Username or Password or Server IP MUST NOT be Empty")
		return
	}

	// New GRPC Client
	gRPC_cli_service_client, err := dialer.NewGRPCClients(server_ip)
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println("Description: Cannot Create New GRPC Client")
		fmt.Println("Source: Install()")
		return
	}

	// Prepare req
	req := &pb.RegisterRequest{
		Username: username,
		Password: password,
	}

	// Request Timeout
	ctx, cancelFunc := context.WithTimeout(context.Background(), CONTEXT_TIMEOUT)
	defer cancelFunc()

	// Sending Request ...
	_, err = gRPC_cli_service_client.Register(ctx, req)
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println("Description: Cannot Register User")
		fmt.Println("Source: Install()")
		return
	}

	// Adding New Server to Config
	err = config.AddNewServerToConfig(server_alias, server_ip, username, password)
	if err != nil {
		fmt.Println("Error Occured in Adding Server to serverConfig.json:", err)
		fmt.Println("Source: Install()")
		return
	}

	fmt.Println("Entry added to userConfig.json file")
}
