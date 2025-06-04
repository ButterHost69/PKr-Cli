package root

import (
	"fmt"
	"strings"

	"github.com/ButterHost69/PKr-Base/config"
	"github.com/ButterHost69/PKr-cli/dialer"
	"github.com/ButterHost69/PKr-cli/filetracker"
	"github.com/ButterHost69/PKr-cli/logger"
)

// [X] Generate Hash of Zip Files
// [X] Add to their .PKr file
// [X] Notify all Connections

// TODO:
// 	[ ] Test the Entire code nothing is tested

// LATER UPDATE:
//
//	[X] Do the Zip and Create Hash in Memory before saving it ]

// TODO MOveover to RPC and
func Push(workspace_name string, workspace_logger *logger.WorkspaceLogger) (int, error) {
	workspace_path, err := config.GetSendWorkspaceFilePath(workspace_name)
	if err != nil {
		return -1, fmt.Errorf("could find workspace.\nError: %v", err)
	}

	fmt.Println("Zip File Created ...")
	hash_zipfile, err := filetracker.ZipData(workspace_path)
	if err != nil {
		return -1, fmt.Errorf("could not zip data.\nError: %v", err)
	}
	hash_zipfile = strings.Split(hash_zipfile, ".")[0]
	fmt.Println("Zip File Created ...\nAdding New Push to Config ...")

	err = config.AddNewPushToConfig(workspace_name, hash_zipfile)
	if err != nil {
		return -1, fmt.Errorf("could add entry to PKR config file.\nError: %v", err)
	}
	fmt.Println("New Push Added into Config ...\nGetting Workspace Connections Using Workspace Path ...")

	// [ ] Notify all Connections
	conf, err := config.ReadFromPKRConfigFile(workspace_path + "\\.PKr\\workspaceConfig.json")
	if err != nil {
		return -1, fmt.Errorf("could not read from .Pkr\\workspaceConfig.json.\nError: %v", err)
	}
	fmt.Println("Comparing Last Hash & Hash of Current Files")
	fmt.Println(conf.LastHash)
	fmt.Println(hash_zipfile)
	if conf.LastHash == hash_zipfile {
		return 0, fmt.Errorf("no new changes detected in 'PUSH'")
	}

	fmt.Println("GetWorkspaceConnectionsUsingPath Done")
	success_count := dialer.PushToConnections(workspace_name, conf.AllConnections, workspace_logger)

	return success_count, nil
	// generate_sha1 :=
}
