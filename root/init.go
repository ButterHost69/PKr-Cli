package root

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
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
		log.Println("Error while getting Server Details from Config:", err)
		log.Println("Source: InitWorkspace()")
		return
	}

	// Check if .PKr folder already exists; if so then do nothing ...
	_, err = os.Stat(".PKr")
	if err == nil {
		log.Println("'.PKr' file already exists")
		log.Println("Workspace is already Initialized")
		return
	} else if os.IsNotExist(err) {
		log.Println("'.PKr' file doesn't exists")
	} else {
		log.Println("Error while checking Existence of Destination file:", err)
		log.Println("Source: InitWorkspace()")
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

	err = config.UpdateLastHash(workspace_name, hash_zipfile)
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
	err = os.WriteFile(zip_destination_path+"AES_KEY", key, 0644)
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
	err = os.WriteFile(zip_destination_path+"AES_IV", iv, 0644)
	if err != nil {
		fmt.Println("Failed to Write AES IV to File:", err)
		fmt.Println("Source: InitWorkspace()")
		return
	}

	// Encrypting Zip File in Chunks
	fmt.Println("Encrypting Zip and Storing for Workspace ...")
	zipped_filepath := zip_destination_path + hash_zipfile + ".zip"
	zip_enc_path := strings.Replace(zipped_filepath, ".zip", ".enc", 1)

	zipped_filepath_obj, err := os.Open(zipped_filepath)
	if err != nil {
		fmt.Println("Failed to Open Zipped File:", err)
		fmt.Println("Source: InitWorkspace()")
		return
	}

	zip_enc_file_obj, err := os.Create(zip_enc_path)
	if err != nil {
		fmt.Println("Failed to Create & Open Enc Zipped File:", err)
		fmt.Println("Source: InitWorkspace()")
		return
	}
	defer zip_enc_file_obj.Close()

	buffer := make([]byte, DATA_CHUNK)
	reader := bufio.NewReader(zipped_filepath_obj)
	writer := bufio.NewWriter(zip_enc_file_obj)

	offset := 0
	for {
		n, err := reader.Read(buffer)
		if err != nil {
			if err == io.EOF {
				fmt.Println("File Encryption Completed ...")
				break
			}
			log.Println("Error while Reading Zip File:", err)
			log.Println("Source: InitWorkspace()")
			return
		}
		encrypted, err := encrypt.EncryptDecryptChunk(buffer[:n], key, iv)
		if err != nil {
			fmt.Println("Failed to Encrypt Chunk:", err)
			fmt.Println("Source: InitWorkspace()")
			return
		}

		_, err = writer.Write(encrypted)
		if err != nil {
			fmt.Println("Failed to Write Chunk to File:", err)
			fmt.Println("Source: InitWorkspace()")
			return
		}

		// Flush buffer to disk after 'FLUSH_AFTER_EVERY_X_CHUNK'
		if offset%FLUSH_AFTER_EVERY_X_MB == 0 {
			err = writer.Flush()
			if err != nil {
				fmt.Println("Error flushing 'writer' after X KB/MB buffer:", err)
				fmt.Println("Soure: InitWorkspace()")
				return
			}
		}
		offset += n
	}

	// Flush buffer to disk at end
	err = writer.Flush()
	if err != nil {
		fmt.Println("Error flushing 'writer' buffer:", err)
		fmt.Println("Soure: InitWorkspace()")
		return
	}
	zipped_filepath_obj.Close() // Close Obj now, so we can delete zip file

	// Removing Zip File
	err = os.Remove(zipped_filepath)
	if err != nil {
		fmt.Println("Error deleting zip file:", err)
		fmt.Println("Source: InitWorkspace()")
		return
	}
	fmt.Println("Removed Zip File - ", zipped_filepath)

	fmt.Println("New Workspace Registered Successfully")
}
