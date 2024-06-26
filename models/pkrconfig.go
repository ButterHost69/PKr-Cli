package models

import (
	"encoding/json"
	"fmt"
	"os"
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
)

func CreatePKRConfigIfNotExits(workspace_name string) (error){
	pkr_config_file_path := WORKSPACE_PKR_DIR + "/workspaceConfig.json"
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

	err = os.WriteFile(pkr_config_file_path, jsonBytes, 0777)
	if err != nil {
		fmt.Println("~ Unable to Write PKrConfig to File")
		return err
	}

	return nil
}

