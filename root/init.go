package root

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ButterHost69/PKr-Base/config"
	"github.com/ButterHost69/PKr-Cli/dialer"
	"github.com/ButterHost69/PKr-Base/encrypt"
	"github.com/ButterHost69/PKr-Base/filetracker"
	"github.com/ButterHost69/PKr-Cli/pb"
)

func InitWorkspace(server_alias, workspace_password string) {
	// Get Details from Config
	server_ip, username, password, err := config.GetServerDetails(server_alias)
	if err != nil {
		log.Println("Error while getting Server Details from Config:", err)
		log.Println("Source: InitWorkspace()")
		return
	}

	// Check if .PKr folder already exists; if so then do nothing ...
	// FIXME: [ ] This Doesnt Work Please Check Why Later
	if _, err := os.Stat(".PKr"); os.IsExist(err) {
		fmt.Println("Error:", err)
		fmt.Println("Description: '.PKr' Already Exists...\nIt seems PKr is already Initialized in this Directory")
		fmt.Println("Source: InitWorkspace()")
		return
	}

	// Create .Pkr Folder ; return if error occured
	if err := os.Mkdir(".PKr", os.ModePerm); err != nil {
		fmt.Println("Error:", err)
		fmt.Println("Description: Cannot Create '.PKr' Directory")
		fmt.Println("Source: InitWorkspace()")
		return
	}

	// Create Keys Folder
	if err := os.Mkdir(".PKr/Keys/", os.ModePerm); err != nil {
		fmt.Println("Error:", err)
		fmt.Println("Description: Cannot Create Keys Folder")
		fmt.Println("Source: InitWorkspace()")
		return
	}

	// Create Files Folder
	if err := os.Mkdir(".PKr/Files/", os.ModePerm); err != nil {
		fmt.Println("Error:", err)
		fmt.Println("Description: Cannot Create Files Folder")
		fmt.Println("Source: InitWorkspace()")
		return
	}

	// Create Current Folder
	if err := os.Mkdir(".PKr/Files/Current/", os.ModePerm); err != nil {
		fmt.Println("Error:", err)
		fmt.Println("Description: Cannot Create Files/Current Folder")
		fmt.Println("Source: InitWorkspace()")
		return
	}

	// Create Changes Folder
	if err := os.Mkdir(".PKr/Files/Changes/", os.ModePerm); err != nil {
		fmt.Println("Error:", err)
		fmt.Println("Description: Cannot Create Files/Changes Folder")
		fmt.Println("Source: InitWorkspace()")
		return
	}

	// Get Curr Directory as Workspace Path
	workspace_path, err := os.Getwd()
	if err != nil {
		log.Println("Error Cannot Call Getwd():", err)
		log.Println("Source: InitWorkspace()")
		return
	}
	workspace_path_split := strings.Split(workspace_path, "\\")
	workspace_name := workspace_path_split[len(workspace_path_split)-1]

	zip_destination_path := workspace_path + "\\.PKr\\Files\\Current\\"
	fmt.Println("Destination For Current Snapshot: ", zip_destination_path)

	// Getting Hash of Zip File
	hash_zipfile, err := filetracker.ZipData(workspace_path, zip_destination_path)
	if err != nil {
		log.Println("Error while Getting Hash of Zipped Data:", err)
		log.Println("Source InitWorkspace()")
		return
	}
	hash_zipfile = strings.Split(hash_zipfile, ".")[0]

	// Create New gRPC Client
	gRPC_cli_service_client, err := dialer.NewGRPCClients(server_ip)
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println("Description: Cannot Create New GRPC Client")
		fmt.Println("Source: InitWorkspace()")
		return
	}

	// Prepare gRPC Request
	req := &pb.RegisterWorkspaceRequest{
		Username:      username,
		Password:      password,
		WorkspaceName: workspace_name,
		LastHash:      hash_zipfile,
	}

	// Request Timeout
	ctx, cancelFunc := context.WithTimeout(context.Background(), CONTEXT_TIMEOUT)
	defer cancelFunc()

	// Sending Request ...
	_, err = gRPC_cli_service_client.RegisterWorkspace(ctx, req)
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println("Description: Cannot Register User")
		fmt.Println("Source: InitWorkspace()")
		return
	}

	// Register the workspace in the main userConfig file
	if err := config.RegisterNewSendWorkspace(server_alias, workspace_name, workspace_path, workspace_password); err != nil {
		fmt.Println("Error:", err)
		fmt.Println("Description: Cannot Register Workspace to userConfig File")
		fmt.Println("Source: InitWorkspace()")
		return
	}

	// Create the workspace config file
	if err := config.CreatePKRConfigIfNotExits(workspace_name, workspace_path); err != nil {
		fmt.Println("Error:", err)
		fmt.Println("Description: Cannot Create .Pkr/PKRConfig.json")
		fmt.Println("Source: InitWorkspace()")
		return
	}

	fmt.Println("Adding New Push to Config ...")
	fmt.Println("Current Main Hash: ", hash_zipfile)
	
	err = config.AddNewPushToConfig(workspace_name, hash_zipfile)
	if err != nil {
		fmt.Println("Error while Adding New Init to Config:", err)
		fmt.Println("Source: InitWorkspace()")
		return
	}

	fmt.Println("Encrypting Zip File...")

	// Generating Key
	fmt.Println("Generating Keys ...")
	key, err := encrypt.AESGenerakeKey(16)
	if err != nil {
		fmt.Println("Failed to Generate AES Keys:", err)
		fmt.Println("Source: InitWorkspace()")
		return 
	}

	// Storing Key
	err = os.WriteFile(zip_destination_path + "AES_KEY", key, 0644)
	if err != nil {
		fmt.Println("Failed to Write AES Key to File:", err)
		fmt.Println("Source: InitWorkspace()")
		return 
	}

	iv, err := encrypt.AESGenerateIV()
	if err != nil {
		fmt.Println("Failed to Generate IV Keys:", err)
		fmt.Println("Source: InitWorkspace()")
		return 
	}

	// Storing IV
	err = os.WriteFile(zip_destination_path + "AES_IV", key, 0644)
	if err != nil {
		fmt.Println("Failed to Write AES IV to File:", err)
		fmt.Println("Source: InitWorkspace()")
		return 
	}

	// Encrypting Zip File
	fmt.Println("Encrypting Zip and Storing for Workspace ...")
	zipped_filepath := zip_destination_path + hash_zipfile + ".zip"
	destination_filepath := strings.Replace(zipped_filepath, ".zip", ".enc", 1)
	if err := encrypt.AESEncrypt(zipped_filepath, destination_filepath, key, iv); err != nil {
		fmt.Println("Failed to Encrypt Data using AES:", err)
		fmt.Println("Source: InitWorkspace()")
		return 
	}

	// Removing Zip File
	err = os.Remove(zipped_filepath ) 
	if err != nil {
		fmt.Println("Error deleting zip file:", err)
		fmt.Println("Source: InitWorkspace()")
		return
	}
	fmt.Println("Removed Zip File - ", zipped_filepath)
	
	fmt.Println("New Workspace Registered Successfully")
}
