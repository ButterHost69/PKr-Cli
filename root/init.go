package root

import (
	"fmt"
	"os"
	"strings"

	"github.com/ButterHost69/PKr-Base/config"
	"github.com/ButterHost69/PKr-cli/dialer"
)

// FIXME: Not Converted
func Init(server_alias, workspace_password string) error {
	// [X] Register Folder to Send Workspace / Export Folder
	// [X] Create a .PKr Folder
	// [X] Create a Keys Folder [Will Store Other Users Public Keys]
	// [X] Create a config file -> Store Shit like User info
	// [X] Log this entry ... [in the .PKr folder of each workspace]

	server_ip, server_username, server_password, err := config.GetServerDetails(server_alias)
	if err != nil {
		return fmt.Errorf("could not fetch server details...\nError: %v", err)
	}

	callHandler := dialer.CallHandler{
		Lipaddr: "",
	}


	workspace_path, err := os.Getwd()
	workspace_path_split := strings.Split(workspace_path, "\\")
	workspaceName := workspace_path_split[len(workspace_path_split)-1]

	if err = callHandler.CallRegisterWorkspace(server_ip, server_username, server_password, workspaceName); err != nil {
		return fmt.Errorf("could not call register new Workspace...\nError: %v", err)
	}

	// Check if .PKr folder already exists; if so then do nothing ...
	// FIXME: [ ] This Doesnt Work Please Check Why Later
	if _, err := os.Stat(".PKr"); os.IsExist(err) {
		return fmt.Errorf(".PKr Already Exists...\nIt seems PKr is already Initialized in this Directory....\nError: %v", err)
	}

	// Create .Pkr Folder ; return if error occured
	if err := os.Mkdir(".PKr", os.ModePerm); err != nil {
		return fmt.Errorf("error Occured In Creating Folder .PKr\nError: %v", err)
	}

	// Create Keys Folder
	if err := os.Mkdir(".PKr/Keys/", os.ModePerm); err != nil {
		return fmt.Errorf("error Occured In Creating Folder Keys\nError: %v", err)
	}

	if err != nil {
		return fmt.Errorf("unable to Identify The Current Working Directory Name.\nError: %v", err)
	}

	// Register the workspace in the main userConfig file
	if err := config.RegisterNewSendWorkspace(server_alias, workspaceName, workspace_path, workspace_password); err != nil {
		return fmt.Errorf("could Not Register The Workspace To the userConfig File.\nError: %v", err)
	}

	// Create the workspace config file
	if err := config.CreatePKRConfigIfNotExits(workspaceName, workspace_path); err != nil {
		return fmt.Errorf("could Not Create .PKr/PKRConfig.json.\nError: %v", err)
	}

	// log := "Workspace '" + workspaceName + "' Created"

	// // Add Entry to the Main File ??? I dont know the Main file path of rn /tmp dir
	// if err := config.AddUsersLogEntry(workspaceName, log); err != nil {
	// 	return fmt.Errorf("could Not add Entry to the Users Logs File.\nError:%v", err)
	// }

	// // Add Entry to Workspace logs
	// if err := config.AddLogEntry(workspaceName, log); err != nil {
	// 	return fmt.Errorf("could Not add Entry to the Users Logs File.\nError:%v", err)
	// }

	return nil
}
