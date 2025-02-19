package root

import (
	"encoding/base64"
	"fmt"
	"os"
	
	"github.com/ButterHost69/PKr-Base/config"
	"github.com/ButterHost69/PKr-cli/dialer"
	"github.com/ButterHost69/PKr-cli/encrypt"
)

// TODO: [ ] Temporary, Find a Better Way ---
// const (
// 	BACKGROUND_SERVER_PORT = 9000
// )
// ---

func Clone(workspace_owner_username, workspace_name, workspace_password, server_alias string) error {
	// [X] Get Public Key From the Host Original Source PC
	// [X]  Encrypt Password
	// [X]  Read Our Key
	// [X]  Send Password and request for InitConnection -> return port
	// [X]  Register The import Folder
	// [X]  Connect to DataServer from the Port
	// [X]  Decrypt the file
	// [X]  Unzip the File

	// [X] Get and Encrypt Key

	server, err := config.GetServerDetailsUsingServerAlias(server_alias)
	if err != nil {
		return fmt.Errorf("error Occured in Retrieving Public Key.\nerror:%v", err)
	}

	public_key, err := dialer.RequestPublicKey(workspace_owner_username, server.ServerIP, server.Username, server.Password)
	if err != nil {
		return fmt.Errorf("error Occured in Retrieving Public Key.\nerror:%v", err)
	}

	fmt.Println("Retrieved Public Key From the Source PC")

	encrypted_password, err := encrypt.EncryptData(workspace_password, string(public_key))
	if err != nil {
		return fmt.Errorf("error Occured in Encrypting Password.\nerror: %v", err)
	}

	my_public_key, err := os.ReadFile("./tmp/mykeys/publickey.pem")
	if err != nil {
		return fmt.Errorf("error Occured in Reading Our Public Key\nPlease Ensure Key is Present at ./tmp/mykeys/publickey.pem\nerror:%v", err)
	}

	base64_public_key := []byte(base64.StdEncoding.EncodeToString(my_public_key))

	// []byte(base64_public_key)
	response, err := dialer.RequestInitNewWorkSpaceConnection(server.ServerIP, server.Username, server.ServerIP, workspace_owner_username, workspace_name, encrypted_password, base64_public_key)
	if err != nil {
		return fmt.Errorf("error Occured in Dialing Init New Workspace Connection.\nerror:%v", err)
	}
	if response != 200 {
		return fmt.Errorf("could not init new workspace")
	}

	currDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error in Retrieving Current Working Directory.\nerror:%v", err)
	}
	err = os.MkdirAll(currDir+"\\.PKr\\", 0777)
	if err != nil {
		return fmt.Errorf("error Occured in Creating .PKr Folder.\nerror: %v", err)
	}

	fmt.Println("Initialized Workspace With the Source PC")

	// TODO Request to GetData() Separately... Pass null string as last hash
	// only_ip := strings.Split(workspace_ip, ":")[0] + ":"
	// fmt.Printf("Data Port: %d\n", port)

	// // [ ]: // For now the workspace's path is currDir, change this later
	// if err = dialer.GetData(workspace_name, only_ip, strconv.Itoa(port), currDir); err != nil {
	// 	return fmt.Errorf("error: Could not Retrieve Data From the Source PC.\nerror: %v", err)
	// }

	return nil
}
