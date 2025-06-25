package root

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ButterHost69/PKr-Base/config"
	"github.com/ButterHost69/PKr-Base/dialer"
	"github.com/ButterHost69/PKr-Base/pb"
)

const CONTEXT_TIMEOUT = 60 * time.Second

func Install(server_ip, username, password string) {
	user_config_file_path := filepath.Join(os.Getenv("LOCALAPPDATA"), "PKr", "Config", "user-config.json")
	_, err := os.Stat(user_config_file_path)
	if err == nil {
		fmt.Println("It Seems PKr is Already Installed...")
		return
	} else if os.IsNotExist(err) {
		fmt.Println("Installing PKr ...")
	} else {
		fmt.Println("Error while checking Existence of user-config file:", err)
		fmt.Println("Source: Install()")
		return
	}

	if server_ip == "" || username == "" || password == "" {
		fmt.Println("Username or Password or Server IP MUST NOT be Empty")
		return
	}

	fmt.Println("Registering User, Sending Request to Server ...")
	// New GRPC Client
	gRPC_cli_service_client, err := dialer.GetNewGRPCClient(server_ip)
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

	// Add Credentials to Config
	err = config.CreateUserConfigIfNotExists(username, password, server_ip)
	if err != nil {
		fmt.Println("Error while Adding Credentials to user-config:", err)
		fmt.Println("Source: Install()")
		return
	}
	fmt.Println("PKr Installed Successfully")
}
