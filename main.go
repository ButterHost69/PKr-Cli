package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ButterHost69/PKr-Cli/root"
)

func printArguments() {
	fmt.Println("Valid Parameters:")
	fmt.Println("	1] install -> Create User and Install PKr")
	fmt.Println("	2] init -> Initialize a Workspace, allows other Users to connect")
	fmt.Println("	3] clone -> Clone an existing Workspace of a different User")
	fmt.Println("	4] list -> List all Send and Get Workspaces")
	fmt.Println("	5] push -> Push new Changes to Listeners")
}

func main() {
	if len(os.Args) < 2 {
		printArguments()
		return
	}

	cmd := strings.ToLower(os.Args[1])
	switch cmd {
	case "install":
		{
			var server_alias, server_ip, username, password string

			fmt.Print("> Enter Server Alias: ")
			fmt.Scan(&server_alias)

			fmt.Print("> Enter Server IP: ")
			fmt.Scan(&server_ip)

			fmt.Print("> Enter Username: ")
			fmt.Scan(&username)

			fmt.Print("> Enter Password: ")
			fmt.Scan(&password)

			fmt.Println("Installing ...")
			root.Install(server_alias, server_ip, username, password)
		}

	case "init":
		{
			var server_alias, workspace_password string

			fmt.Print("> Enter Server Alias: ")
			fmt.Scan(&server_alias)

			fmt.Print("> Enter Workspace Password: ")
			fmt.Scan(&workspace_password)

			fmt.Println("Initializing New Workspace ...")
			root.InitWorkspace(server_alias, workspace_password)
		}

	case "clone":
		{
			var workspace_owner_username string
			var workspace_name string
			var workspace_password string
			var server_alias string

			fmt.Print("> Enter the Workspace Owner Username: ")
			fmt.Scan(&workspace_owner_username)

			fmt.Print("> Enter Server Alias: ")
			fmt.Scan(&server_alias)

			fmt.Print("> Enter Workspace Name: ")
			fmt.Scan(&workspace_name)

			fmt.Print("> Enter Workspace Password: ")
			fmt.Scan(&workspace_password)

			fmt.Println("Cloning ...")
			root.Clone(workspace_owner_username, workspace_name, workspace_password, server_alias)
		}

	case "list":
		{
			var server_alias string

			fmt.Print("> Enter Server Alias: ")
			fmt.Scan(&server_alias)
			fmt.Println("Fetching All Workspaces ...")

			root.ListAllWorkspaces(server_alias)
		}

	case "push":
		{
			var server_alias string

			fmt.Print("> Enter Server Alias: ")
			fmt.Scan(&server_alias)

			current_working_directory, err := os.Getwd()
			if err != nil {
				fmt.Println("Could not Get Current Working Directory :", err)
				fmt.Println("Source: main()")
				return
			}

			current_working_directory_split := strings.Split(current_working_directory, string(filepath.Separator))
			workspace_name := strings.TrimSpace(current_working_directory_split[len(current_working_directory_split)-1])

			fmt.Printf("Pushing Workpace: %s ...\n", workspace_name)
			root.Push(workspace_name, server_alias)
		}

	default:
		printArguments()
	}
}
