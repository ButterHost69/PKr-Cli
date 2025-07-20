package main

import (
	"bufio"
	"fmt"
	"os"
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
			var server_ip, username, password string

			fmt.Print("> Enter Username: ")
			fmt.Scan(&username)

			fmt.Print("> Enter Password: ")
			fmt.Scan(&password)

			fmt.Print("> Enter Server IP: ")
			fmt.Scan(&server_ip)

			root.Install(server_ip, username, password)
		}

	case "init":
		{
			var workspace_password, push_desc string
			reader := bufio.NewReader(os.Stdin)

			fmt.Print("> Enter Workspace Password: ")
			workspace_password, _ = reader.ReadString('\n')
			workspace_password = strings.TrimSpace(workspace_password)

			fmt.Print("> Enter Push Description: ")
			push_desc, _ = reader.ReadString('\n')
			push_desc = strings.TrimSpace(push_desc)

			root.InitWorkspace(workspace_password, push_desc)
		}

	case "clone":
		{
			var workspace_owner_username string
			var workspace_name string
			var workspace_password string

			fmt.Println("WARNING: All Previous files'll be DELETED & REPLACED by files Received from Workspace Owner")
			fmt.Print("> Enter the Workspace Owner Username: ")
			fmt.Scan(&workspace_owner_username)

			fmt.Print("> Enter Workspace Name: ")
			fmt.Scan(&workspace_name)

			fmt.Print("> Enter Workspace Password: ")
			fmt.Scan(&workspace_password)

			fmt.Println("Cloning ...")
			root.Clone(workspace_owner_username, workspace_name, workspace_password)
		}

	case "list":
		{
			fmt.Println("Fetching All Workspaces from Server ...")
			root.ListAllWorkspaces()
		}

	case "push":
		{
			var workspace_name, push_desc string
			reader := bufio.NewReader(os.Stdin)

			fmt.Print("> Enter Workspace Name: ")
			workspace_name, _ = reader.ReadString('\n')
			workspace_name = strings.TrimSpace(workspace_name)

			fmt.Print("> Enter Push Description: ")
			push_desc, _ = reader.ReadString('\n')
			push_desc = strings.TrimSpace(push_desc)

			fmt.Printf("Pushing Workpace: %s ...\n", workspace_name)
			root.Push(workspace_name, push_desc)
		}

	default:
		printArguments()
	}
}
