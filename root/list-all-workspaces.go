package root

import (
	"context"
	"fmt"

	"github.com/ButterHost69/PKr-Base/config"
	"github.com/ButterHost69/PKr-Base/dialer"
	"github.com/ButterHost69/PKr-Base/pb"
)

func ListAllWorkspaces(server_alias string) {
	// Get Details from Config
	server_ip, username, password, err := config.GetServerDetails(server_alias)
	if err != nil {
		fmt.Println("Error while getting Server Details from Config:", err)
		fmt.Println("Source: ListAllWorkspaces()")
		return
	}

	// New GRPC Client
	gRPC_cli_service_client, err := dialer.NewGRPCClients(server_ip)
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println("Description: Cannot Create New GRPC Client")
		fmt.Println("Source: ListAllWorkspaces()")
		return
	}

	// Prepare req
	req := &pb.GetAllWorkspacesRequest{
		Username: username,
		Password: password,
	}

	// Request Timeout
	ctx, cancelFunc := context.WithTimeout(context.Background(), CONTEXT_TIMEOUT)
	defer cancelFunc()

	res, err := gRPC_cli_service_client.GetAllWorkspaces(ctx, req)
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println("Description: Cannot Notify New Push to Listeners")
		fmt.Println("Source: ListAllWorkspaces()")
	}

	for _, workspace := range res.Workspaces {
		fmt.Printf("Workspace Owner: %s, Workspace Name: %s\n", workspace.WorkspaceOwner, workspace.WorkspaceName)
	}
	fmt.Println("Workspace Fetched Successfully")
}
