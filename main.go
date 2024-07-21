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
	case "install": {

		// -> Setup Username [X]
		// -> Generate Public and Private Keys [X]
		// -> Register gRPC Server as a service 
		var username string
		fmt.Print("Enter a Username : ")
		fmt.Scan(&username)
		fmt.Printf("Okay %s, Setting Up Your System...\n", username)
		fmt.Println("This Might Take Some Time...")
		models.CreateUserIfNotExists(username)
	}

	case "uninstall":
		
	case "get":

	case "push":
		
	// Maybe its Done	
	case "init" : {
		// Register Folder to Send Workspace / Export Folder [X]
		// Create a .PKr Folder [X]
		// Create a Keys Folder [Will Store Other Users Public Keys] [X]
		// Create a config file -> Store Shit like User info [X]
		// Log this entry ... [in the .PKr folder of each workspace] [X]
		var workspace_password string
		fmt.Print("Please Enter A Password: ")
		fmt.Scan(&workspace_password)

		// Check if .PKr folder already exists; if so then do nothing ...
		// This Doesnt Work Please Check Why Later
		if _, err := os.Stat(".PKr"); os.IsExist(err) {
			fmt.Println(".PKr Already Exists...")
			fmt.Println("It seems PKr is already Initialized...")
			return 
		}
		
		// Create .Pkr Folder ; return if error occured
		if err := os.Mkdir(".PKr", os.ModePerm) ; err != nil {
			fmt.Printf("Error Occured In Creating Folder .PKr\nError: %v\n", err)
			return
		}
		
		// Create Keys Folder
		if err := os.Mkdir(".PKr/Keys/", os.ModePerm) ; err != nil {
			fmt.Printf("Error Occured In Creating Folder Keys\nError: %v\n", err)
			return
		}

		workspace_path, err := os.Getwd()
		workspace_path_split := strings.Split(workspace_path, "\\")
		workspaceName := workspace_path_split[len(workspace_path_split)	- 1]
		if err != nil {
			fmt.Println("Unable to Identify The Current Working Directory Name")
			fmt.Printf("Error: %v\n",err)
			return	
		}

		// Register the workspace in the main userConfig file
		if err :=  models.RegisterNewSendWorkspace(workspaceName, workspace_path, workspace_password); err != nil {
			fmt.Println("Could Not Register The Workspace To the userConfig File")
			fmt.Printf("Error: %v\n",err)
			return
		}

		// Create the workspace config file
		if err := models.CreatePKRConfigIfNotExits(workspaceName, workspace_path); err != nil {
			fmt.Println("Could Not Create .PKr/PKRConfig.json")
			fmt.Printf("Error: %v\n",err)
			return
		}

		log := "Workspace '" + workspaceName + "' Created"

		// Add Entry to the Main File ??? I dont know the Main file path of rn /tmp dir
		if err := models.AddUsersLogEntry(workspaceName, log) ; err != nil {
			fmt.Println("Could Not add Entry to the Users Logs File")
			fmt.Printf("Error: %v\n",err)
			return
		}

		// Add Entry to Workspace logs
		if err := models.AddLogEntry(workspaceName, log) ; err != nil {
			fmt.Println("Could Not add Entry to the Workspace Logs File")
			fmt.Printf("Error: %v\n",err)
			return
		}

		fmt.Println("Workspace Created Successfully !!")
		return

	}	
  }
}