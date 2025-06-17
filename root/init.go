package root

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ButterHost69/PKr-Base/config"
	"github.com/ButterHost69/PKr-Base/dialer"
	"github.com/ButterHost69/PKr-Base/encrypt"
	"github.com/ButterHost69/PKr-Base/filetracker"
	"github.com/ButterHost69/PKr-Base/pb"
)

func InitWorkspace(server_alias, workspace_password string) {
	// Get Details from Config
	server_ip, username, password, err := config.GetServerDetails(server_alias)
	if err != nil {
		fmt.Println("Error while getting Server Details from Config:", err)
		fmt.Println("Source: InitWorkspace()")
		return
	}

	// Check if .PKr folder already exists; if so then do nothing ...
	_, err = os.Stat(".PKr")
	if err == nil {
		fmt.Println("'.PKr' file already exists")
		fmt.Println("Workspace is already Initialized")
		return
	} else if os.IsNotExist(err) {
		fmt.Println("'.PKr' file doesn't exists")
	} else {
		fmt.Println("Error while checking Existence of Destination file:", err)
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
	if err := os.Mkdir(filepath.Join(".PKr", "Keys"), os.ModePerm); err != nil {
		fmt.Println("Error:", err)
		fmt.Println("Description: Cannot Create Keys Folder")
		fmt.Println("Source: InitWorkspace()")
		return
	}

	// Create Files Folder
	if err := os.Mkdir(filepath.Join(".PKr", "Files"), os.ModePerm); err != nil {
		fmt.Println("Error:", err)
		fmt.Println("Description: Cannot Create Files Folder")
		fmt.Println("Source: InitWorkspace()")
		return
	}

	// Create Current Folder
	if err := os.Mkdir(filepath.Join(".PKr", "Files", "Current"), os.ModePerm); err != nil {
		fmt.Println("Error:", err)
		fmt.Println("Description: Cannot Create Files/Current Folder")
		fmt.Println("Source: InitWorkspace()")
		return
	}

	// Create Changes Folder
	if err := os.Mkdir(filepath.Join(".PKr", "Files", "Changes"), os.ModePerm); err != nil {
		fmt.Println("Error:", err)
		fmt.Println("Description: Cannot Create Files/Changes Folder")
		fmt.Println("Source: InitWorkspace()")
		return
	}

	// Get Curr Directory as Workspace Path
	workspace_path, err := os.Getwd()
	if err != nil {
		fmt.Println("Error Cannot Call Getwd():", err)
		fmt.Println("Source: InitWorkspace()")
		return
	}

	workspace_path_split := strings.Split(workspace_path, string(filepath.Separator))
	workspace_name := workspace_path_split[len(workspace_path_split)-1]

	zip_destination_path := filepath.Join(workspace_path, ".PKr", "Files", "Current") + string(filepath.Separator)
	fmt.Println("Destination For Current Snapshot: ", zip_destination_path)

	// Getting Hash of Zip File
	hash_zipfile, err := filetracker.ZipData(workspace_path, zip_destination_path)
	if err != nil {
		fmt.Println("Error while Getting Hash of Zipped Data:", err)
		fmt.Println("Source InitWorkspace()")
		return
	}
	hash_zipfile = strings.Split(hash_zipfile, ".")[0]

	// Create the workspace config file
	if err := config.CreatePKRConfigIfNotExits(workspace_name, workspace_path); err != nil {
		fmt.Println("Error:", err)
		fmt.Println("Description: Cannot Create .Pkr/PKRConfig.json")
		fmt.Println("Source: InitWorkspace()")
		return
	}

	fmt.Println("Current Main Hash: ", hash_zipfile)
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
	err = os.WriteFile(zip_destination_path+"AES_KEY", key, 0644)
	if err != nil {
		fmt.Println("Failed to Write AES Key to File:", err)
		fmt.Println("Source: InitWorkspace()")
		return
	}

	// Generating IV
	iv, err := encrypt.AESGenerateIV()
	if err != nil {
		fmt.Println("Failed to Generate IV Keys:", err)
		fmt.Println("Source: InitWorkspace()")
		return
	}

	// Storing IV
	err = os.WriteFile(zip_destination_path+"AES_IV", iv, 0644)
	if err != nil {
		fmt.Println("Failed to Write AES IV to File:", err)
		fmt.Println("Source: InitWorkspace()")
		return
	}

	// Encrypting Zip File of Entire in Chunks
	fmt.Println("Encrypting Zip and Storing for Workspace ...")
	zipped_filepath := zip_destination_path + hash_zipfile + ".zip"
	zip_enc_path := strings.Replace(zipped_filepath, ".zip", ".enc", 1)

	err = encryptZipFileAndStore(zipped_filepath, zip_enc_path, key, iv)
	if err != nil {
		fmt.Println("Error while Encrypting Zip File of Entire Workspace, Storing it & Deleting Zip File:", err)
		fmt.Println("Source: InitWorkspace()")
		return
	}

	// Create Tree
	tree, err := config.GetNewTree(workspace_path)
	if err != nil {
		fmt.Println("Error Could not Create Tree:", err)
		fmt.Println("Source: InitWorkspace()")
		return
	}

	// Store Tree
	err = config.WriteToFileTree(workspace_path, tree)
	if err != nil {
		fmt.Println("Error Write Tree to file:", err)
		fmt.Println("Source: InitWorkspace()")
		return
	}

	// Write Updates in PKr Config
	fmt.Println("Comparing Changes And Updating to workspace Config")
	changes := config.CompareTrees(config.FileTree{}, tree, hash_zipfile)
	err = config.AppendWorkspaceUpdates(changes, filepath.Join(workspace_path, ".PKr", "workspaceConfig.json"))
	if err != nil {
		fmt.Println("Error while Adding Changes to PKr Config:", err)
		fmt.Println("Source: InitWorkspace()")
		return
	}

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

	err = config.UpdateLastHash(workspace_name, hash_zipfile)
	if err != nil {
		fmt.Println("Error while Updating Last to Config:", err)
		fmt.Println("Source: InitWorkspace()")
		return
	}

	fmt.Println("Workspace Initialized Successfully")
}
