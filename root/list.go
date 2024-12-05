package root

import (
	"fmt"

	"github.com/ButterHost69/PKr-cli/models"
)

func List() error {
	userConfigFile, err := models.ReadFromUserConfigFile()
	if err != nil {
		return fmt.Errorf("could not list Workspaces...\nError: %v", err)
	}

	fmt.Println(" -- Send Workspaces -- ")
	for idx, workspace := range userConfigFile.Sendworkspaces {
		fmt.Printf("%d] %s: \n", idx, workspace.WorkspaceName)
		fmt.Printf("	- %s\n\n", workspace.WorkspacePath)

	}

	fmt.Println(" -- Get Workspaces -- ")
	for idx, workspace := range userConfigFile.GetWorkspaces {
		fmt.Printf("%d] %s: \n", idx, workspace.WorkspaceName)
		fmt.Printf("	- %s\n\n", workspace.WorkspacePath)
	}

	return nil
}