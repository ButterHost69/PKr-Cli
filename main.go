package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ButterHost69/PKr-cli/dialer"
	"github.com/ButterHost69/PKr-cli/encrypt"
	"github.com/ButterHost69/PKr-cli/models"
)

const (
	BACKGROUND_SERVER_PORT = 9000
)

// TODO: [ ] Shift everything to flag based, no terminal inputs, take all inputs as flags
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Required Minimum 2 Args")
		return
	}
	cmd := strings.ToLower(os.Args[1])
	switch cmd {
	case "install":
		{
			// -> Setup Username [X]
			// -> Generate Public and Private Keys [X]
			// -> Register gRPC Server as a service
			var username string
			fmt.Print("Enter a Username : ")
			fmt.Scan(&username)
			fmt.Printf("Okay %s, Setting Up Your System...\n", username)
			fmt.Println("This Might Take Some Time...")
			models.CreateUserIfNotExists(username)
			models.CreateServerConfigFiles()
		}

	case "uninstall":

	case "get":

	case "push":

	case "clone":
		{
			// Get Public Key From the Host Original Source PC [X]
			// Encrypt Password [X]
			// Read Our Key [X]
			// Send Password and request for InitConnection -> return port [X]
			// Register The import Folder [X]
			// Connect to DataServer from the Port  [X]
			// Decrypt the file [X]
			// Unzip the File [X]

			var workspace_ip string
			var workspace_name string
			var workspace_password string

			fmt.Print("> Enter the Workspace IP [addr:port]: ")
			fmt.Scan(&workspace_ip)

			fmt.Print("> Enter the Workspace Name: ")
			fmt.Scan(&workspace_name)

			fmt.Print("> Enter the Workspace Password: ")
			fmt.Scan(&workspace_password)

			// Get and Encrypt Key
			public_key, err := dialer.GetPublicKey(workspace_ip)
			if err != nil {
				fmt.Println("Error Occured in Retrieving Public Key")
				fmt.Println(err)
				return
			}

			fmt.Println("Retrieved Public Key From the Source PC")
			encrypted_password, err := encrypt.EncryptData(workspace_password, string(public_key))
			if err != nil {
				fmt.Println("Error Occured in Encrypting Password")
				fmt.Println(err)
				return
			}

			my_public_key, err := os.ReadFile("./tmp/mykeys/publickey.pem")
			if err != nil {
				fmt.Println("Error Occured in Reading Our Public Key")
				fmt.Println("Please Ensure Key is Present at ./tmp/mykeys/publickey.pem")
				fmt.Println(err)
			}

			base64_public_key := base64.StdEncoding.EncodeToString(my_public_key)

			port, err := dialer.InitNewWorkSpaceConnection(workspace_ip, workspace_name, encrypted_password, strconv.Itoa(BACKGROUND_SERVER_PORT), []byte(base64_public_key))
			if err != nil {
				fmt.Println("Error Occured in Dialing Init New Workspace Connection")
				fmt.Println(err)
			}
			currDir, err := os.Getwd()
			if err != nil {
				fmt.Println("Error in Retrieving Current Working Directory")
				fmt.Println(err)
			}
			err = os.MkdirAll(currDir+"\\.PKr\\", 0777)
			if err != nil {
				fmt.Println("Error Occured in Creating .PKr Folder")
				fmt.Println(err)
			}
			if err = models.AddGetWorkspaceFolderToUserConfig(workspace_name, currDir, workspace_ip); err != nil {
				fmt.Println("Error in adding GetConnection to the Main User Config Folder")
				fmt.Println(err)
			}

			fmt.Println("Initialized Workspace With the Source PC")
			only_ip := strings.Split(workspace_ip, ":")[0] + ":"
			fmt.Printf("Data Port: %d\n", port)
			if err = dialer.GetData(workspace_name, only_ip, strconv.Itoa(port)); err != nil {
				fmt.Println("Error: Could not Retrieve Data From the Source PC")
				fmt.Println(err)
			}
		}

	// Maybe its Done
	case "init":
		{
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
			if err := os.Mkdir(".PKr", os.ModePerm); err != nil {
				fmt.Printf("Error Occured In Creating Folder .PKr\nError: %v\n", err)
				return
			}

			// Create Keys Folder
			if err := os.Mkdir(".PKr/Keys/", os.ModePerm); err != nil {
				fmt.Printf("Error Occured In Creating Folder Keys\nError: %v\n", err)
				return
			}

			workspace_path, err := os.Getwd()
			workspace_path_split := strings.Split(workspace_path, "\\")
			workspaceName := workspace_path_split[len(workspace_path_split)-1]
			if err != nil {
				fmt.Println("Unable to Identify The Current Working Directory Name")
				fmt.Printf("Error: %v\n", err)
				return
			}

			// Register the workspace in the main userConfig file
			if err := models.RegisterNewSendWorkspace(workspaceName, workspace_path, workspace_password); err != nil {
				fmt.Println("Could Not Register The Workspace To the userConfig File")
				fmt.Printf("Error: %v\n", err)
				return
			}

			// Create the workspace config file
			if err := models.CreatePKRConfigIfNotExits(workspaceName, workspace_path); err != nil {
				fmt.Println("Could Not Create .PKr/PKRConfig.json")
				fmt.Printf("Error: %v\n", err)
				return
			}

			log := "Workspace '" + workspaceName + "' Created"

			// Add Entry to the Main File ??? I dont know the Main file path of rn /tmp dir
			if err := models.AddUsersLogEntry(workspaceName, log); err != nil {
				fmt.Println("Could Not add Entry to the Users Logs File")
				fmt.Printf("Error: %v\n", err)
				return
			}

			// Add Entry to Workspace logs
			if err := models.AddLogEntry(workspaceName, log); err != nil {
				fmt.Println("Could Not add Entry to the Workspace Logs File")
				fmt.Printf("Error: %v\n", err)
				return
			}

			fmt.Println("Workspace Created Successfully !!")
			return
		}

	// For server
	// Mainly for IP and shit
	// Try to make the code as swappable as possible
	//
	// Was working on server setup ~ check models, it is still partial
	case "server":
		{
			if len(os.Args) < 3 {
				fmt.Println("server requires additional arguments")
				return
			}
			opts := strings.ToLower(os.Args[2])
			switch opts {
			case "setup":
				{
					// [ ] Create a main Server Json file in tmp(root) dir
					// [ ] Allow user to connect to multiple server
					// [ ] Store Server IP, and your username and password (user can have multiple username and password)
					// [ ] Send Create User request to server, save to Server Json file and display the username to user at terminal

					var server_ip string
					fmt.Print("Please Enter Server IP: ")
					fmt.Scan(&server_ip)

					var server_username string
					fmt.Print("Please Enter A Username for Server Connection: ")
					fmt.Scan(&server_username)

					var server_password string
					fmt.Print("Please Enter A Password for Server Connection: ")
					fmt.Scan(&server_password)

					if server_ip == "" || server_username == "" || server_password == "" {
						fmt.Println("ip or username or password cannot be empty")
						return
					}

					if err := models.AddNewServerToConfig(server_username, server_password, server_ip); err != nil {
						fmt.Println("Error Occured in Adding Server to serverConfig.json")
						fmt.Println(err)
					}
				}
			}
		}
	}
}
