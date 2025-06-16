package root

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ButterHost69/PKr-Base/config"
	"github.com/ButterHost69/PKr-Base/encrypt"
	"github.com/ButterHost69/PKr-Base/filetracker"
	"github.com/ButterHost69/PKr-Base/dialer"
	"github.com/ButterHost69/PKr-Base/pb"
)

func Push(workspace_name, server_alias string) {
	// Getting Workspace Absolute Path
	workspace_path, err := config.GetSendWorkspaceFilePath(workspace_name)
	if err != nil {
		log.Println("Error while getting Absolute Workspace Path:", err)
		log.Println("Source: Push()")
		return
	}

	fmt.Println("Checking if Changes are Present in the Workspace")
	fmt.Println("Creating File Tree Structure for the Workspace")
	new_tree, err := config.GetNewTree(workspace_path)
	if err != nil {
		fmt.Println("Could Not Create New Tree of the Current Workspace")
		fmt.Println(err)
		return
	}

	fmt.Println("Fetching Old File Tree")
	old_tree, err := config.ReadFromTreeFile(workspace_path)
	if err != nil {
		fmt.Println("Could Not Read Old Tree of the file_tree.json")
		fmt.Println(err)
		return
	}

	fmt.Println("Comparing Trees ...")
	updates := config.CompareTrees(old_tree, new_tree, "NotRequired")
	if len(updates.Changes) == 0 {
		fmt.Println("No New Changes Detected in 'PUSH'")
		return
	}

	fmt.Println("Updates Detected ...")
	fmt.Println("Clearing CURRENT")

	// Reading Last Hash from Config
	conf, err := config.ReadFromPKRConfigFile(filepath.Join(workspace_path, ".PKr", "workspaceConfig.json"))
	if err != nil {
		log.Println("Error while Reading from PKr Config File:", err)
		log.Println("Source: Push()")
		return
	}

	old_zipped_filepath := filepath.Join(workspace_path, ".PKr", "Files", "Current", conf.LastHash + ".enc")
	err = os.Remove(old_zipped_filepath)
	if err != nil {
		fmt.Println("Error deleting old zip file:", err)
		fmt.Println("Source: Push()")
		return
	}

	OLDIV := filepath.Join(workspace_path, ".PKr", "Files", "Current", "AES_IV")
	err = os.Remove(OLDIV)
	if err != nil {
		fmt.Println("Error deleting old IV:", err)
		fmt.Println("Source: Push()")
		return
	}

	OLDKEY := filepath.Join(workspace_path, ".PKr", "Files", "Current", "AES_KEY")
	err = os.Remove(OLDKEY)
	if err != nil {
		fmt.Println("Error deleting old Key file:", err)
		fmt.Println("Source: Push()")
		return
	}

	fmt.Println("Creating Entire Zip File for Current ...")
	
	destination_path := filepath.Join(workspace_path, ".PKr", "Files", "Current") + string(filepath.Separator)
	hash_zipfile, err := filetracker.ZipData(workspace_path, destination_path)
	if err != nil {
		fmt.Println("Error Creating Zip file :", err)
		fmt.Println("Source: Push()")
		return
	}
	new_hash := strings.Split(hash_zipfile, ".")[0]
	fmt.Println("Zip File Created")
	fmt.Println("Hash: ", new_hash)

	fmt.Println("Encrypting CURRENT zip file")
	fmt.Println("Generating and Storing AES Keys")

	key, err := encrypt.AESGenerakeKey(16)
	if err != nil {
		fmt.Println("Failed to Generate AES Keys:", err)
		fmt.Println("Source: Push()")
		return
	}

	// Storing Key
	err = os.WriteFile(filepath.Join(destination_path,"AES_KEY"), key, 0644)
	if err != nil {
		fmt.Println("Failed to Write AES Key to File:", err)
		fmt.Println("Source: Push()")
		return
	}

	iv, err := encrypt.AESGenerateIV()
	if err != nil {
		fmt.Println("Failed to Generate IV Keys:", err)
		fmt.Println("Source: Push()")
		return
	}

	// Storing IV
	err = os.WriteFile(filepath.Join(destination_path, "AES_IV"), iv, 0644)
	if err != nil {
		fmt.Println("Failed to Write AES IV to File:", err)
		fmt.Println("Source: Push()")
		return
	}

	// Encrypting Zip File
	fmt.Println("Encrypting Zip and Storing for Workspace ...")
	zipped_filepath := filepath.Join(destination_path, new_hash + ".zip")
	destination_filepath := strings.Replace(zipped_filepath, ".zip", ".enc", 1)
	if err := encrypt.AESEncrypt(zipped_filepath, destination_filepath, key, iv); err != nil {
		fmt.Println("Failed to Encrypt Data using AES:", err)
		fmt.Println("Source: Push()")
		return
	}

	// Removing Zip File
	err = os.Remove(zipped_filepath)
	if err != nil {
		fmt.Println("Error deleting zip file:", err)
		fmt.Println("Source: Push()")
		return
	}
	fmt.Println("Removed Zip File - ", zipped_filepath)

	fmt.Println("Registering New Push to Workspace Config")

	fmt.Println("Updating Last Hash")
	err = config.UpdateLastHash(workspace_name, new_hash)
	if err != nil {
		fmt.Println("Error while Updating Last Hash to Config:", err)
		fmt.Println("Source: Push()")
		return
	}

	fmt.Println("Updating File Tree")
	err = config.WriteToFileTree(workspace_path, new_tree)
	if err != nil {
		fmt.Println("Error Write Tree to file:", err)
		fmt.Println("Source: InitWorkspace()")
		return
	}

	fmt.Println("Updating New Changes to the Workspace Config")
	updates.Hash = new_hash
	err = config.AppendWorkspaceUpdates(updates, filepath.Join(workspace_path, ".PKr", "workspaceConfig.json"))
	if err != nil {
		fmt.Println("Error Write Tree to file:", err)
		fmt.Println("Source: InitWorkspace()")
		return
	}

	// TODO:
	// Cache the Changes into Zip
	fmt.Println("Merging Updates ...")
	mupdates, err := config.MergeUpdates(filepath.Join(workspace_path, ".PKr", "workspaceConfig.json"), conf.LastHash, new_hash)
	if err != nil {
		fmt.Println("Unable to Merge Updates:", err)
		fmt.Println("Source: InitWorkspace()")
		return
	}

	fmt.Println("Zipping Changed Files")
	fmt.Println("Generating Changes Hash Name ...")
	files_hash_list := []string{}
	for _, changes := range mupdates.Changes {
		if changes.Type == "Updated" {
			files_hash_list = append(files_hash_list, updates.Hash)
		}
	}

	changes_hash_name := encrypt.GeneratHashFromFileNames(files_hash_list)
	fmt.Println("Changes Hash:", changes_hash_name)

	err = filetracker.ZipUpdates(mupdates, workspace_path, changes_hash_name)
	if err != nil {
		log.Println("Error while Creating Zip for Changes:", err)
		log.Println("Source: Push()")
		return
	}

	// Encrypt Zip and Store Keys
	fmt.Println("Encrypting Changes Zip File...")
	changes_path := filepath.Join(workspace_path, ".PKr", "Files", "Changes", changes_hash_name)
	// Generating Key
	fmt.Println("Generating Keys for Changes File ...")
	changes_key, err := encrypt.AESGenerakeKey(16)
	if err != nil {
		fmt.Println("Failed to Generate AES Keys:", err)
		fmt.Println("Source: Push()")
		return
	}

	// Storing Key
	err = os.WriteFile(filepath.Join(changes_path,"AES_KEY"), changes_key, 0644)
	if err != nil {
		fmt.Println("Failed to Write AES Key to File:", err)
		fmt.Println("Source: Push()")
		return
	}

	changes_iv, err := encrypt.AESGenerateIV()
	if err != nil {
		fmt.Println("Failed to Generate IV Keys:", err)
		fmt.Println("Source: Push()")
		return
	}

	// Storing IV
	err = os.WriteFile(filepath.Join(changes_path, "AES_IV"), changes_iv, 0644)
	if err != nil {
		fmt.Println("Failed to Write AES IV to File:", err)
		fmt.Println("Source: Push()")
		return
	}

	// Encrypting Zip File
	fmt.Println("Encrypting Zip and Storing for Workspace ...")
	changes_zipped_filepath := filepath.Join(changes_path, changes_hash_name+".zip")
	changes_destination_filepath := strings.Replace(changes_zipped_filepath, ".zip", ".enc", 1)
	if err := encrypt.AESEncrypt(changes_zipped_filepath, changes_destination_filepath, changes_key, changes_iv); err != nil {
		fmt.Println("Failed to Encrypt Data using AES:", err)
		fmt.Println("Source: Push()")
		return
	}

	// Removing Zip File
	err = os.Remove(changes_zipped_filepath)
	if err != nil {
		fmt.Println("Error deleting zip file:", err)
		fmt.Println("Source: Push()")
		return
	}
	fmt.Println("Removed Changes Zip File - ", changes_zipped_filepath)



	fmt.Println("Notifying to Listeners")
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

	err = config.UpdateLastHash(workspace_name, new_hash)
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println("Description: Cannot Add New Push to Config")
		fmt.Println("Source: Push()")
		return
	}
	fmt.Println("New Push Registered Successfully")
}
