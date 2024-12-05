package root

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

// TODO: [ ] Temporary, Find a Better Way ---
const (
	BACKGROUND_SERVER_PORT = 9000
)
// ---


func Clone(workspace_ip, workspace_name, workspace_password string) error {
	// [X] Get Public Key From the Host Original Source PC 
	// [X]  Encrypt Password 
	// [X]  Read Our Key 
	// [X]  Send Password and request for InitConnection -> return port 
	// [X]  Register The import Folder 
	// [X]  Connect to DataServer from the Port  
	// [X]  Decrypt the file 
	// [X]  Unzip the File 

	// [ ] Get and Encrypt Key
	public_key, err := dialer.GetPublicKey(workspace_ip)
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

	base64_public_key := base64.StdEncoding.EncodeToString(my_public_key)

	port, err := dialer.InitNewWorkSpaceConnection(workspace_ip, workspace_name, encrypted_password, strconv.Itoa(BACKGROUND_SERVER_PORT), []byte(base64_public_key))
	if err != nil {
		return fmt.Errorf("error Occured in Dialing Init New Workspace Connection.\nerror:%v", err)
	}
	currDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error in Retrieving Current Working Directory.\nerror:%v", err)
	}
	err = os.MkdirAll(currDir+"\\.PKr\\", 0777)
	if err != nil {
		return fmt.Errorf("error Occured in Creating .PKr Folder.\nerror: %v", err)
	}
	if err = models.AddGetWorkspaceFolderToUserConfig(workspace_name, currDir, workspace_ip); err != nil {
		return fmt.Errorf("error in adding GetConnection to the Main User Config Folder.\nerror:%v", err)
	}

	fmt.Println("Initialized Workspace With the Source PC")

	only_ip := strings.Split(workspace_ip, ":")[0] + ":"
	fmt.Printf("Data Port: %d\n", port)
	if err = dialer.GetData(workspace_name, only_ip, strconv.Itoa(port)); err != nil {
		return fmt.Errorf("error: Could not Retrieve Data From the Source PC.\nerror: %v", err)
	}

	return nil
}