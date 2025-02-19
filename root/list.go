package root

import (
	"fmt"

	"github.com/ButterHost69/PKr-Base/config"
)

func List() error {
	getWorkspaces, err := config.GetAllGetWorkspaces()
	if err != nil {
		return fmt.Errorf("could not list Workspaces...\nError: %v", err)
	}

	fmt.Println(" -- Get Workspaces -- ")
	for idx, workspace := range getWorkspaces {
		fmt.Printf("%d] %s: \n", idx, workspace.WorkspaceName)
		fmt.Printf("	- %s\n\n", workspace.WorkspacePath)

	}

	sendWorkspaces, err := config.GetAllSendWorkspaces()
	if err != nil {
		return fmt.Errorf("could not list Workspaces...\nError: %v", err)
	}
	fmt.Println(" -- Send Workspaces -- ")
	for idx, workspace := range sendWorkspaces {
		fmt.Printf("%d] %s: \n", idx, workspace.WorkspaceName)
		fmt.Printf("	- %s\n\n", workspace.WorkspacePath)
	}

	return nil
}
