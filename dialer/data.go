package dialer

import (
	"fmt"
	"net"

	"github.com/ButterHost69/PKr-cli/encrypt"
	"github.com/ButterHost69/PKr-cli/filetracker"
	"github.com/ButterHost69/PKr-cli/myrpc"
)

func RequestGetData(receivers_ip, myusername, mypassword, workspace_name, workspace_password, server_ip string, udpConn *net.UDPConn, rpcClientHandler myrpc.ClientCallHandler, clientHandlerName string) (int, error) {
	// Get Data, Key, IV
	// Decrypt Key, IV
	// Decrypt Data
	// Store Zip
	// Clear the GetFolder except .PKr
	// Unzip data on the GetFolder
	// Store Last Hash into Config Files

	// FIXME Encrypt Workspace Password before Sending - Store when Calling GetPublicKey for the first time for the user
	// TODO: Instead of Last Hash as "", fetch it from config(ALSO CHECK IF LAST_HASH IS EVEN NEEDED)
	res, err := rpcClientHandler.CallGetData(myusername, server_ip, workspace_name, workspace_password, "", receivers_ip, udpConn, clientHandlerName)
	if err != nil {
		return res.Response, err
	}
	fmt.Println("Response of Get Data from Workspace Owner")

	if res.Response != 200 {
		return res.Response, nil
	}

	data_bytes := res.Data
	key_bytes := res.KeyBytes
	iv_bytes := res.IVBytes

	decrypted_key, err := encrypt.DecryptData(string(key_bytes))
	if err != nil {
		return 400, err
	}

	decrypted_iv, err := encrypt.DecryptData(string(iv_bytes))
	if err != nil {
		return 400, err
	}

	data, err := encrypt.AESDecrypt(data_bytes, decrypted_key, decrypted_iv)
	if err != nil {
		return 400, err
	}

	workspacePath := "."

	zip_file_path := workspacePath + "\\.PKr\\" + res.NewHash + ".zip"
	if err = filetracker.SaveDataToFile(data, zip_file_path); err != nil {
		return 400, err
	}

	if err = filetracker.CleanFilesFromWorkspace(workspacePath); err != nil {
		return 400, err
	}

	// Unzip Content
	// unzip_file_path := getworkspace.WorkspacePath
	if err = filetracker.UnzipData(zip_file_path, workspacePath+"\\"); err != nil {
		return 400, err
	}

	return 200, nil
}
