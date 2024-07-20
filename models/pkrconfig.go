package models

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	// "github.com/ButterHost69/PKr-cli/models"
)

type PKRConfig struct {
	WorkspaceName 	string		`json:"workspace_name"`
	AllConnections	Connection	`json:"all_connections"`
}

type Connection struct {
	Username      string `json:"username"`
	CurrentIP     string `json:"current_ip"`
	CurrentPort   string `json:"current_port"`
	PublicKeyPath string `json:"public_key_path"`
}

const (
	WORKSPACE_PKR_DIR = ".PKr"
	LOGS_PKR_FILE_PATH = WORKSPACE_PKR_DIR + "\\logs.txt"
	WORKSPACE_CONFIG_FILE_PATH = WORKSPACE_PKR_DIR + "\\workspaceConfig.json"
)

func CreatePKRConfigIfNotExits(workspace_name string, workspace_file_path string) (error){
	pkr_config_file_path := workspace_file_path + "\\" + WORKSPACE_CONFIG_FILE_PATH
	if _, err := os.Stat(pkr_config_file_path); os.IsExist(err) {
		fmt.Println("~ workspaceConfig.jso already Exists")
		return err
	}

	pkrconf := PKRConfig {
		WorkspaceName: workspace_name,
	}

	jsonBytes, err := json.Marshal(pkrconf)
	if err != nil {
		fmt.Println("~ Unable to Parse PKrConfig to JSON")
		return err
	}

	// Creating Workspace Config File
	err = os.WriteFile(pkr_config_file_path, jsonBytes, 0777)
	if err != nil {
		fmt.Println("~ Unable to Write PKrConfig to File")
		return err
	}

	return nil
}


// Logs Entry of all the events occurred related to the workspace
// Also Creates the Log File by default
func AddLogEntry(workspace_name string, log_entry string) (error){
	workspace_path, err := GetWorkspaceFilePath(workspace_name)
	if err != nil {
		return err
	}

	// Adds the ".Pkr/logs.txt"
	workspace_path += "\\" + LOGS_PKR_FILE_PATH

	// Opens or Creates the Log File 
	file, err := os.OpenFile(workspace_path,  os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	
	defer file.Close()
	log.SetOutput(file)
	log.Printf(log_entry + "\n", log.LstdFlags)

		
	return nil
}



