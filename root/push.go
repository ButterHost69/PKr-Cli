package root

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/ButterHost69/PKr-Base/config"
	"github.com/ButterHost69/PKr-Cli/dialer"
	"github.com/ButterHost69/PKr-Base/filetracker"
	"github.com/ButterHost69/PKr-Cli/pb"
)

func Push(workspace_name, server_alias string) {
	// Getting Workspace Absolute Path
	workspace_path, err := config.GetSendWorkspaceFilePath(workspace_name)
	if err != nil {
		log.Println("Error while getting Absolute Workspace Path:", err)
		log.Println("Source: Push()")
		return
	}

	fmt.Println("Creating Zip File ...")

	// TODO: Check if Destination is Proper -- If Zips Work ; Delete Later
	destination_path := workspace_path + "\\.PKr\\Files\\Current\\"
	hash_zipfile, err := filetracker.ZipData(workspace_path, destination_path)
	if err != nil {
		return
	}
	hash_zipfile = strings.Split(hash_zipfile, ".")[0]
	fmt.Println("Zip File Created")

	// Reading Last Hash from Config
	conf, err := config.ReadFromPKRConfigFile(workspace_path + "\\.PKr\\workspaceConfig.json")
	if err != nil {
		log.Println("Error while Reading from PKr Config File:", err)
		log.Println("Source: Push()")
		return
	}

	fmt.Println("Comparing Last Hash to Hash of New Pushed Files ...")
	if conf.LastHash == hash_zipfile {
		fmt.Println("No New Changes Detected in 'PUSH'")
		return
	}
	fmt.Println("Changes Detected, Notifying this to Listeners")

	// Get Details from Config
	server_ip, username, password, err := config.GetServerDetails(server_alias)
	if err != nil {
		log.Println("Error while getting Server Details from Config:", err)
		log.Println("Source: Push()")
		return
	}

	// New GRPC Client
	gRPC_cli_service_client, err := dialer.NewGRPCClients(server_ip)
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println("Description: Cannot Create New GRPC Client")
		fmt.Println("Source: Push()")
		return
	}

	// Prepare req
	req := &pb.NotifyNewPushToListenersRequest{
		WorkspaceOwnerUsername: username,
		WorkspaceOwnerPassword: password,
		WorkspaceName:          workspace_name,
		NewWorkspaceHash:       hash_zipfile,
	}

	// Request Timeout
	ctx, cancelFunc := context.WithTimeout(context.Background(), CONTEXT_TIMEOUT)
	defer cancelFunc()

	_, err = gRPC_cli_service_client.NotifyNewPushToListeners(ctx, req)
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println("Description: Cannot Notify New Push to Listeners")
		fmt.Println("Source: Push()")
		return
	}

	err = config.AddNewPushToConfig(workspace_name, hash_zipfile)
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println("Description: Cannot Add New Push to Config")
		fmt.Println("Source: Push()")
		return
	}
	fmt.Println("New Push Registered Successfully")
}
