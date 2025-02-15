package root

import (
	"fmt"

	"github.com/ButterHost69/PKr-cli/config"
	"github.com/ButterHost69/PKr-cli/dialer"
)

func Server_RegisterWorkspace() {
	// TODO: [X] Show the workspaces available
	// TODO: [X] Denote the workspace as linked to server
	// TODO: [X] Send request and register the workspace

	// TODO: [ ] Test this code ...
	// TODO: [ ] to the added workspace add a field ??if_server could do, idk why though. Maybe...
	var workspace_name string
	var server_ip string

	fmt.Print("Enter Send Workpace Name to Register: ")
	fmt.Scan(&workspace_name)
	fmt.Print("Enter Server IP: ")
	fmt.Scan(&server_ip)
	ifDone := config.AddServerToWorkpace(workspace_name, server_ip)
	if !ifDone {
		fmt.Println("error")
		return
	}
	username, password, err := config.GetServerUsernamePassword(server_ip)
	if err != nil {
		fmt.Printf("error Could not get Server Usernamea and Password...\nError: %v\n", err)
		return
	}
	err = dialer.RegisterWorkspace(server_ip, username, password, workspace_name)
	if err != nil {
		fmt.Printf("error Could not get Server Usernamea and Password...\nError: %v\n", err)
		return
	}
	fmt.Println("Done.")
}
