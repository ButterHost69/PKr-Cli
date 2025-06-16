package root

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ButterHost69/PKr-Base/config"
	"github.com/ButterHost69/PKr-Base/dialer"
	"github.com/ButterHost69/PKr-Base/encrypt"
	"github.com/ButterHost69/PKr-Base/filetracker"
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

	fmt.Println("Creating Zip File ...")

	zip_destination_path := workspace_path + "\\.PKr\\Files\\Current\\"
	hash_zipfile, err := filetracker.ZipData(workspace_path, zip_destination_path)
	if err != nil {
		return
	}
	hash_zipfile = strings.Split(hash_zipfile, ".")[0]
	zipped_filepath := zip_destination_path + hash_zipfile + ".zip"
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
		// Removing Zip File
		err = os.Remove(zipped_filepath)
		if err != nil {
			fmt.Println("Error deleting zip file:", err)
			fmt.Println("Source: Push()")
			return
		}
		fmt.Println("Deleting Zip File which is created during Push,Because there were no changes")
		fmt.Println("Removed Zip File:", zipped_filepath)
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

	err = config.UpdateLastHash(workspace_name, hash_zipfile)
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println("Description: Cannot Add New Push to Config")
		fmt.Println("Source: Push()")
		return
	}
	fmt.Println("New Push Registered Successfully")

	// Generating Key
	fmt.Println("Generating Keys ...")
	key, err := encrypt.AESGenerakeKey(16)
	if err != nil {
		fmt.Println("Failed to Generate AES Keys:", err)
		fmt.Println("Source: Push()")
		return
	}

	// Storing Key
	err = os.WriteFile(zip_destination_path+"AES_KEY", key, 0644)
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
	err = os.WriteFile(zip_destination_path+"AES_IV", iv, 0644)
	if err != nil {
		fmt.Println("Failed to Write AES IV to File:", err)
		fmt.Println("Source: Push()")
		return
	}

	// Encrypting Zip File
	fmt.Println("Encrypting Zip and Storing for Workspace ...")
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

	// Removing Previous Commit's Encrypted Zip File
	old_zipped_filepath := zip_destination_path + conf.LastHash + ".enc"
	err = os.Remove(old_zipped_filepath)
	if err != nil {
		fmt.Println("Error deleting old enc zip file:", err)
		fmt.Println("Source: Push()")
		return
	}
	fmt.Println("Removed Prev Commit's Enc Zip File - ", zipped_filepath)
}
