package root

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/ButterHost69/PKr-Base/config"
	"github.com/ButterHost69/PKr-Base/dialer"
	"github.com/ButterHost69/PKr-Base/encrypt"
	"github.com/ButterHost69/PKr-Base/filetracker"
	"github.com/ButterHost69/PKr-Base/pb"
)

// Encrypts Zip File, Stores Enc Zip File & Deletes Old Zip File
func EncryptZipFileAndStore(zipped_filepath, zip_enc_path string, key, iv []byte) error {
	zipped_filepath_obj, err := os.Open(zipped_filepath)
	if err != nil {
		fmt.Println("Failed to Open Zipped File:", err)
		fmt.Println("Source: encryptZipFileAndStore()")
		return err
	}
	defer zipped_filepath_obj.Close()

	zip_enc_file_obj, err := os.Create(zip_enc_path)
	if err != nil {
		fmt.Println("Failed to Create & Open Enc Zipped File:", err)
		fmt.Println("Source: encryptZipFileAndStore()")
		return err
	}
	defer zip_enc_file_obj.Close()

	buffer := make([]byte, DATA_CHUNK)
	reader := bufio.NewReader(zipped_filepath_obj)
	writer := bufio.NewWriter(zip_enc_file_obj)

	// Reading from Zip File, Encrypting it & Writing it to Enc Zip File
	offset := 0
	for {
		n, err := reader.Read(buffer)
		if err != nil {
			if err == io.EOF {
				fmt.Println("File Encryption Completed ...")
				break
			}
			fmt.Println("Error while Reading Zip File:", err)
			fmt.Println("Source: encryptZipFileAndStore()")
			return err
		}
		encrypted, err := encrypt.EncryptDecryptChunk(buffer[:n], key, iv)
		if err != nil {
			fmt.Println("Failed to Encrypt Chunk:", err)
			fmt.Println("Source: encryptZipFileAndStore()")
			return err
		}

		_, err = writer.Write(encrypted)
		if err != nil {
			fmt.Println("Failed to Write Chunk to File:", err)
			fmt.Println("Source: encryptZipFileAndStore()")
			return err
		}

		// Flush buffer to disk after 'FLUSH_AFTER_EVERY_X_CHUNK'
		if offset%FLUSH_AFTER_EVERY_X_MB == 0 {
			err = writer.Flush()
			if err != nil {
				fmt.Println("Error flushing 'writer' after X KB/MB buffer:", err)
				fmt.Println("Soure: encryptZipFileAndStore()")
				return err
			}
		}
		offset += n
	}

	// Flush buffer to disk at end
	err = writer.Flush()
	if err != nil {
		fmt.Println("Error flushing 'writer' buffer:", err)
		fmt.Println("Soure: encryptZipFileAndStore()")
		return err
	}
	zipped_filepath_obj.Close() // Close Obj now, so we can delete zip file
	zip_enc_file_obj.Close()

	// Removing Zip File
	err = os.Remove(zipped_filepath)
	if err != nil {
		fmt.Println("Error deleting zip file:", err)
		fmt.Println("Source: encryptZipFileAndStore()")
		return err
	}
	fmt.Println("Removed Zip File - ", zipped_filepath)
	return nil
}

func Push(workspace_name, server_alias string) {
	// Getting Workspace's Absolute Path
	workspace_path, err := config.GetSendWorkspaceFilePath(workspace_name)
	if err != nil {
		fmt.Println("Error while getting Absolute Workspace Path:", err)
		fmt.Println("Source: Push()")
		return
	}

	// Creating New Tree
	fmt.Println("Checking if Changes are Present in the Workspace")
	fmt.Println("Creating File Tree Structure for the Workspace")
	new_tree, err := config.GetNewTree(workspace_path)
	if err != nil {
		fmt.Println("Could Not Create New Tree of the Current Workspace:", err)
		fmt.Println("Source: Push()")
		return
	}

	// Getting Old Tree
	fmt.Println("Fetching Old File Tree")
	old_tree, err := config.ReadFromTreeFile(workspace_path)
	if err != nil {
		fmt.Println("Could Not Read Old Tree of the file_tree.json:", err)
		fmt.Println("Source: Push()")
		return
	}

	// Comparing Old & New Trees
	fmt.Println("Comparing Trees ...")
	updates := config.CompareTrees(old_tree, new_tree, "NotRequired")
	if len(updates.Changes) == 0 {
		fmt.Println("No New Changes Detected in 'PUSH'")
		return
	}

	fmt.Println("Updates Detected ...")
	// Reading Last Hash from Config
	workspace_conf, err := config.ReadFromPKRConfigFile(filepath.Join(workspace_path, ".PKr", "workspaceConfig.json"))
	if err != nil {
		fmt.Println("Error while Reading from PKr Config File:", err)
		fmt.Println("Source: Push()")
		return
	}

	// Creating Zip of Entire Workspace
	zip_destination_path := filepath.Join(workspace_path, ".PKr", "Files", "Current") + string(filepath.Separator)
	fmt.Println("Creating Zip of Entire Workspace ...")
	hash_entire_workspace, err := filetracker.ZipData(workspace_path, zip_destination_path)
	if err != nil {
		fmt.Println("Error while Creating Zip File:", err)
		fmt.Println("Source: Push()")
		return
	}
	hash_entire_workspace = strings.Split(hash_entire_workspace, ".")[0]
	zipped_filepath := zip_destination_path + hash_entire_workspace + ".zip"
	fmt.Println("Zip File Created")

	// Generating AES Key
	fmt.Println("Generating Keys ...")
	key, err := encrypt.AESGenerakeKey(16)
	if err != nil {
		fmt.Println("Failed to Generate AES Keys:", err)
		fmt.Println("Source: Push()")
		return
	}

	// Storing AES Key
	err = os.WriteFile(zip_destination_path+"AES_KEY", key, 0644)
	if err != nil {
		fmt.Println("Failed to Write AES Key to File:", err)
		fmt.Println("Source: Push()")
		return
	}

	// Generating AES IV
	iv, err := encrypt.AESGenerateIV()
	if err != nil {
		fmt.Println("Failed to Generate IV Keys:", err)
		fmt.Println("Source: Push()")
		return
	}

	// Storing AES IV
	err = os.WriteFile(zip_destination_path+"AES_IV", iv, 0644)
	if err != nil {
		fmt.Println("Failed to Write AES IV to File:", err)
		fmt.Println("Source: Push()")
		return
	}

	// Encrypting Entire Workspace's Zip File in Chunks
	fmt.Println("Encrypting Zip and Storing it ...")
	zip_enc_path := strings.Replace(zipped_filepath, ".zip", ".enc", 1)
	err = EncryptZipFileAndStore(zipped_filepath, zip_enc_path, key, iv)
	if err != nil {
		fmt.Println("Error while Encrypting Zip File of Entire Workspace, Storing it & Deleting Zip File:", err)
		fmt.Println("Source: Push()")
		return
	}

	// Removing Previous Commit's Entire Encrypted Zip File
	old_zipped_filepath := zip_destination_path + workspace_conf.LastHash + ".enc"
	err = os.Remove(old_zipped_filepath)
	if err != nil {
		fmt.Println("Error while Deleting Old Enc Zip File:", err)
		fmt.Println("Source: Push()")
		return
	}
	fmt.Println("Removed Prev Commit's Enc Zip File - ", old_zipped_filepath)

	// Writing New Tree to Config
	fmt.Println("Updating File Tree")
	err = config.WriteToFileTree(workspace_path, new_tree)
	if err != nil {
		fmt.Println("Error Write Tree to file:", err)
		fmt.Println("Source: Push()")
		return
	}

	// Adding Changes to PKr Config
	fmt.Println("Updating New Changes to the Workspace Config")
	updates.Hash = hash_entire_workspace
	err = config.AppendWorkspaceUpdates(updates, workspace_path)
	if err != nil {
		fmt.Println("Error Write Tree to file:", err)
		fmt.Println("Source: Push()")
		return
	}

	fmt.Println("Generating Changes Hash Name ...")
	files_hash_list := []string{}
	for _, changes := range updates.Changes {
		if changes.Type == "Updated" {
			fmt.Println(changes.FilePath)
			files_hash_list = append(files_hash_list, changes.FilePath)
			files_hash_list = append(files_hash_list, changes.FileHash)
		}
	}
	fmt.Println("updates.Hash:", updates.Hash)
	fmt.Println("Files Hash List:", files_hash_list)
	fmt.Println("Updates:", updates)

	changes_hash_name := encrypt.GeneratHashFromFileNames(files_hash_list)
	changes_path := filepath.Join(workspace_path, ".PKr", "Files", "Changes", changes_hash_name)
	fmt.Println("Changes Hash:", changes_hash_name)
	fmt.Println(files_hash_list)

	is_updates_cache_present, err := filetracker.AreUpdatesCached(workspace_path, changes_hash_name)
	if err != nil {
		fmt.Println("Error while Checking Whether Updates're Already Cached or Not")
		fmt.Println("Source: Push()")
		return
	}

	// Skip Encryption & Zip Creation if the same changes already exists
	if !is_updates_cache_present {
		fmt.Println("Not Skipping Encrypting & Zipping ...")
		// Create Zip File of Changes
		err = filetracker.ZipUpdates(updates, workspace_path, changes_hash_name)
		if err != nil {
			fmt.Println("Error while Creating Zip for Changes:", err)
			fmt.Println("Source: Push()")
			return
		}
		// Create AES Key for Changes Zip
		changes_key, err := encrypt.AESGenerakeKey(16)
		if err != nil {
			fmt.Println("Failed to Generate AES Keys:", err)
			fmt.Println("Source: Push()")
			return
		}

		// Storing Key
		err = os.WriteFile(filepath.Join(changes_path, "AES_KEY"), changes_key, 0644)
		if err != nil {
			fmt.Println("Failed to Write AES Key to File:", err)
			fmt.Println("Source: Push()")
			return
		}

		// Creating AES IV for Changes Zip
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

		// Encrypting Changes Zip
		changes_zipped_filepath := filepath.Join(changes_path, changes_hash_name+".zip")
		changes_enc_zip_filepath := strings.Replace(changes_zipped_filepath, ".zip", ".enc", 1)

		err = EncryptZipFileAndStore(changes_zipped_filepath, changes_enc_zip_filepath, changes_key, changes_iv)
		if err != nil {
			fmt.Println("Error while Encrypting 'Changes' Zip File, Storing it & Deleting Zip File:", err)
			fmt.Println("Source: Push()")
			return
		}
	}

	// Get Details from Config
	server_ip, username, password, err := config.GetServerDetails(server_alias)
	if err != nil {
		fmt.Println("Error while getting Server Details from Config:", err)
		fmt.Println("Source: Push()")
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
		NewWorkspaceHash:       hash_entire_workspace,
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

	err = config.UpdateLastHash(workspace_name, hash_entire_workspace)
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println("Description: Cannot Add New Push to Config")
		fmt.Println("Source: Push()")
		return
	}
	fmt.Println("Registered New Push Successfully")
}
