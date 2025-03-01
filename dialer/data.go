package dialer

func RequestGetData(rcv_username, myusername, mypassword, server_ip, workspace_name, workspace_password, last_hash string) (int, error) {
	// Get Data, Key, IV
	// Decrypt Key, IV
	// Decrypt Data
	// Store Zip
	// Clear the GetFolder except .PKr
	// Unzip data on the GetFolder
	// Store Last Hash into Config Files

	// FIXME Encrypt Workspace Password before Sending - Store when Calling GetPublicKey for the first time for the user

	// callHandler := myrpc.ServerCallHandler{}

	// port := utils.GetRandomPort()
	// receivers_ip, err := callHandler.CallRequestPunchFromReciever(server_ip, rcv_username, myusername, mypassword, port)
	// if err != nil {
	// 	return 400, err
	// }

	// handler := myrpc.ClientCallHandler{}
	// res, err := handler.CallGetData(myusername, server_ip, workspace_name, workspace_password, last_hash, receivers_ip, fmt.Sprintf(":%d", port))
	// if err != nil {
	// 	return res.Response, err
	// }

	// if res.Response != 200 {
	// 	return res.Response, nil
	// }

	// data_bytes := res.Data
	// key_bytes := res.KeyBytes
	// iv_bytes := res.IVBytes

	// decrypted_key, err := encrypt.DecryptData(string(key_bytes))
	// if err != nil {
	// 	return 400, err
	// }

	// decrypted_iv, err := encrypt.DecryptData(string(iv_bytes))
	// if err != nil {
	// 	return 400, err
	// }

	// data, err := encrypt.AESDecrypt(data_bytes, decrypted_key, decrypted_iv)
	// if err != nil {
	// 	return 400, err
	// }

	// workspacePath, err := config.GetWorkspaceFilePath(workspace_name)
	// if err != nil {
	// 	return 400, err
	// }

	// zip_file_path := workspacePath + "\\.PKr\\" + last_hash + ".zip"
	// if err = filetracker.SaveDataToFile(data, zip_file_path); err != nil {
	// 	return 400, err
	// }

	// if err = filetracker.CleanFilesFromWorkspace(workspacePath); err != nil {
	// 	return 400, err
	// }

	// // Unzip Content
	// // unzip_file_path := getworkspace.WorkspacePath
	// if err = filetracker.UnzipData(zip_file_path, workspacePath+"\\"); err != nil {
	// 	return 400, err
	// }

	return 200, nil
}
