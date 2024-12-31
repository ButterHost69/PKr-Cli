package models

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	// "github.com/ButterHost69/PKr-cli/models"
)

type PKRConfig struct {
	WorkspaceName  string       `json:"workspace_name"`
	AllConnections []Connection `json:"all_connections"`
	LastHash       string       `json:"last_hash"`
}

type Connection struct {
	Username      string `json:"username"`
	CurrentIP     string `json:"current_ip"`
	CurrentPort   string `json:"current_port"`
	PublicKeyPath string `json:"public_key_path"`
}

const (
	WORKSPACE_PKR_DIR          = ".PKr"
	LOGS_PKR_FILE_PATH         = WORKSPACE_PKR_DIR + "\\logs.txt"
	WORKSPACE_CONFIG_FILE_PATH = WORKSPACE_PKR_DIR + "\\workspaceConfig.json"
)

func CreatePKRConfigIfNotExits(workspace_name string, workspace_file_path string) error {
	pkr_config_file_path := workspace_file_path + "\\" + WORKSPACE_CONFIG_FILE_PATH
	if _, err := os.Stat(pkr_config_file_path); os.IsExist(err) {
		fmt.Println("~ workspaceConfig.json already Exists")
		return err
	}

	pkrconf := PKRConfig{
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

// TODO: [X] Fix the log part
func readFromPKRConfigFile(workspace_config_path string) (PKRConfig, error) {
	file, err := os.Open(workspace_config_path)
	if err != nil {
		AddUsersLogEntry(workspace_config_path, "error in opening PKR config file.... pls check if .PKr/workspaceConfig.json available ")
		return PKRConfig{}, err
	}
	defer file.Close()

	var pkrConfig PKRConfig
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&pkrConfig)
	if err != nil {
		AddUsersLogEntry(workspace_config_path, "error in decoding json data")
		return PKRConfig{}, err
	}

	// fmt.Println(pkrConfig)
	return pkrConfig, nil
}

func writeToPKRConfigFile(workspace_config_path string, newPKRConfing PKRConfig) error {
	jsonData, err := json.MarshalIndent(newPKRConfing, "", "	")
	// fmt.Println(jsonData)
	if err != nil {
		fmt.Println("error occured in Marshalling the data to JSON")
		fmt.Println(err)
		return err
	}

	// fmt.Println(string(jsonData))
	err = os.WriteFile(workspace_config_path, jsonData, 0777)
	if err != nil {
		fmt.Println("error occured in storing data in userconfig file")
		fmt.Println(err)
		return err
	}

	return nil
}

// Logs Entry of all the events occurred related to the workspace
// Also Creates the Log File by default
func AddLogEntry(workspace_name string, log_entry string) error {
	workspace_path, err := GetWorkspaceFilePath(workspace_name)
	if err != nil {
		return err
	}

	// Adds the ".Pkr/logs.txt"
	workspace_path += "\\" + LOGS_PKR_FILE_PATH

	// Opens or Creates the Log File
	file, err := os.OpenFile(workspace_path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	defer file.Close()
	log.SetOutput(file)
	log.Printf(log_entry+"\n", log.LstdFlags)

	return nil
}

func AddNewPushToConfig(workspace_name, zipfile_path string) error {
	workspace_path, err := GetWorkspaceFilePath(workspace_name)
	if err != nil {
		return err
	}

	workspace_path = workspace_path + "\\" + WORKSPACE_CONFIG_FILE_PATH
	// fmt.Println("[LOG DELETE LATER]Workspace Path: ", workspace_path)

	workspace_json, err := readFromPKRConfigFile(workspace_path)
	if err != nil {
		return fmt.Errorf("could not add entry to config file.\nError: %v", err)
	}

	workspace_json.LastHash = zipfile_path

	if err := writeToPKRConfigFile(workspace_path, workspace_json); err != nil {
		return fmt.Errorf("error in writing the update hash to file: %s.\nError: %v", workspace_path, err)
	}
	return nil
}

func GetWorkspaceConnectionsIP(workspace_path string) ([]string, error) {
	ipList := []string{}
	pkrconfig, err := readFromPKRConfigFile(workspace_path + "\\" + WORKSPACE_CONFIG_FILE_PATH)
	if err != nil {
		return ipList, err
	}

	for _, connection := range pkrconfig.AllConnections {
		ipList = append(ipList, connection.CurrentIP)
	}

	return ipList, nil
}
