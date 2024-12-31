package root

import (
	"fmt"
	"strings"

	"github.com/ButterHost69/PKr-cli/filetracker"
	"github.com/ButterHost69/PKr-cli/dialer"
	"github.com/ButterHost69/PKr-cli/models"
)


// [X] Generate Hash of Zip Files
// [X] Add to their .PKr file
// [X] Notify all Connections

// TODO:
// 	[ ] Test the Entire code nothing is tested

// LATER UPDATE:
//
//	[X] Do the Zip and Create Hash in Memory before saving it ]
func Push(workspace_name string) (int, error) {
	workspace_path, err := models.GetWorkspaceFilePath(workspace_name)
	if err != nil {
		return -1, fmt.Errorf("could find workspace.\nError: %v", err)
	}

	fmt.Println("Zip File Created ...")
	hash_zipfile, err := filetracker.ZipData(workspace_path)
	if err != nil {
		return -1, fmt.Errorf("could not zip data.\nError: %v", err)
	}

	//  [X] Rename Zip file to hash name
	err = models.AddNewPushToConfig(workspace_name, strings.Split(hash_zipfile, ".")[0])
	if err != nil {
		return -1, fmt.Errorf("could add entry to PKR config file.\nError: %v", err)
	}

	// [ ] Notify all Connections
	connections, err := models.GetWorkspaceConnectionsIP(workspace_path)
	if err != nil {
		return -1, fmt.Errorf("could not get workspace connections IP.\nError: %v", err)
	}

	success_count := dialer.PushToConnections(workspace_name, connections)

	return success_count, nil
	// generate_sha1 :=
}
