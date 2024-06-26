package main

import (
	"fmt"
	"os"
	"strings"
	"github.com/ButterHost69/PKr-cli/models"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Required Minimum 2 Args")
		return
	}
	cmd := strings.ToLower(os.Args[1])
	switch cmd {
	case "install":
		// -> Setup Username [X]
		// -> Generate Public and Private Keys [X]
		// -> Register gRPC Server as a service 
		var username string
		fmt.Print("Enter a Username : ")
		fmt.Scan(&username)
		fmt.Printf("Okay %s, Setting Up Your System...\n", username)
		fmt.Println("This Might Take Some Time...")
		models.CreateUserIfNotExists(username)

	case "uninstall":
		
	case "get":

	case "push":
		
	case "init" :
		// Register Folder to Send Workspace / Export Folder [X]
		// Create a .PKr Folder [X]
		// Create a Log Folder ??? or a Log File ???
		// Create a Keys Folder [Will Store Other Users Public Keys] 
		// Create a config file -> Store Shit like User info ... ??
		var workspace_password string
		fmt.Print("Please Enter A Password: ")
		fmt.Scan(&workspace_password)

		if _, err := os.Stat(".PKr"); os.IsExist(err) {
			fmt.Println(".PKr Already Exists...")
			fmt.Println("It seems PKr is already Initialized...")
			return 
		}

		if err := os.Mkdir(".PKr", os.ModePerm) ; err != nil {
			fmt.Printf("Error Occured In Creating Folder .PKr\nError: %v\n", err)

			workspace_path, err := os.Getwd()
			workspace_path_split := strings.Split(workspace_path, "/")
			workspaceName := workspace_path_split[len(workspace_path_split)	- 1]
			if err != nil {
				fmt.Println("Unable to Identify The Current Working Directory Name")
				fmt.Printf("Error: %v\n",err)
				return	
			}

			if err :=  models.RegisterNewSendWorkspace(workspaceName, workspace_path, workspace_password); err != nil {
				fmt.Println("Could Not Register The Workspace To the userConfig File")
				fmt.Printf("Error: %v\n",err)
				return
			}

			if err := models.CreatePKRConfigIfNotExits(workspaceName); err != nil {
				fmt.Println("Could Not Create .PKr/PKRConfig.json")
				fmt.Printf("Error: %v\n",err)
				return
			}
		}
	}

		


}