package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ButterHost69/PKr-cli/root"
)

const (
	BACKGROUND_SERVER_PORT = 9000
)

// TODO: [ ] Shift everything to flag based, support terminal inputs and CLI
// TODO: [ ] Refactor the code.
//
//	TODO: [ ] Why the fuck are there two model files ?? Make it in to 1
//	TODO: [ ] Why are there print statements in files other than main.
//	TODO: [ ] Write Tests, bro why am I doing this manual... Use docker maybe to simulate the whole thing ???
//  TODO: [ ] Add verbose option that allows for prints in "inside" functions (anything aside from main)
// 	TODO: [ ] Use a Better Logging Method(preferbly homemade), also logs are made partially, not everything is logged.

var (
	TUI bool
	CLI bool
)

func Init() {
	flag.BoolVar(&TUI, "tui", false, "Use Application in TUI Mode")
	flag.BoolVar(&CLI, "cli", false, "Use Application in CLI Mode")
	flag.Parse()
}

func main() {
	Init()

	// TUI -> Takes Input from stdin and print output/errors on stdout
	// CLI -> Input Passed as flag with command, output/errors on stdout
	if !TUI && !CLI {
		fmt.Println("Must Define Mode to use PKr in")
		fmt.Println("	1] `PKr -tui` -> For Terminal User Interface. Takes Input through stdout, requires less flags")
		fmt.Println("	2] `PKr -cli` -> For Command Line Interface. Requires Input as flags")
		return
	}

	if len(os.Args) < 3 {
		fmt.Printf("Required Minimum 3 Args\n\n")
		fmt.Println("Valid Parameters:")
		fmt.Println("	1] install -> Create User and Install PKr")
		fmt.Println("	2] init -> Initialize a Workspace, allows other Users to connect")
		fmt.Println("	3] clone -> Clone an existing Workspace of a different User")
		fmt.Println("	4] list -> List all Send and Get Workspaces")
		fmt.Println("	5] server -> Connect With a Server to Manage Multiple Dynamic Connections")
		return
	}

	cmd := strings.ToLower(os.Args[2])
	if TUI {
		// TODO: [ ] Check if a User is created or not before exec any command. Notify user with suggestion to install PKr first.
		switch cmd {
		case "install":
			{
				var username string
				fmt.Print("Enter a Username : ")
				fmt.Scan(&username)

				fmt.Printf("Okay %s, Setting Up Your System...\n", username)
				fmt.Println("This Might Take Some Time...")

				err := root.Install(username)
				if err != nil {
					fmt.Println("Could not Install PKr.")
					fmt.Println(err)
					return
				}

				fmt.Printf(" ~ Created User : %s\n", username)
				return
			}

		case "uninstall":

		case "get":

		case "push":

		case "clone":
			{
				var workspace_ip string
				var workspace_name string
				var workspace_password string

				fmt.Print("> Enter the Workspace IP [addr:port]: ")
				fmt.Scan(&workspace_ip)

				fmt.Print("> Enter the Workspace Name: ")
				fmt.Scan(&workspace_name)

				fmt.Print("> Enter the Workspace Password: ")
				fmt.Scan(&workspace_password)

				err := root.Clone(workspace_ip, workspace_name, workspace_password)
				if err != nil {
					fmt.Printf("Error Occured in Cloning Workspace: %s at IP: %s\n", workspace_name, workspace_ip)
					fmt.Println(err)
					return
				}

				fmt.Printf("Successfully Cloned Workspace: %s\n", workspace_name)
				return
			}

		case "init":
			{
				var workspace_password string
				fmt.Print("Please Enter A Password: ")
				fmt.Scan(&workspace_password)
				if err := root.Init(workspace_password); err != nil {
					fmt.Println("Error Occured in Initialize a New Workspace")
					fmt.Printf("error: %v\n", err)
					return
				}

				fmt.Println("Workspace Created Successfully !!")
				return
			}

			// list Created workspaces
		case "list":
			{
				if err := root.List(); err != nil {
					fmt.Println("Could Not List Workspace Info")
					fmt.Printf("Error: %v", err)
					return
				}
			}

		// For server
		// Mainly for IP and shit
		// Try to make the code as swappable as possible
		//
		// Was working on server setup ~ check models, it is still partial
		case "server":
			{
				if len(os.Args) < 4 {
					fmt.Println("server requires additional arguments")
					fmt.Println("Valid Arguments:")
					fmt.Println("	1] setup -> Initialize Connection with New PKr Server")
					fmt.Println("	2] register_workspace -> Initialize An Existing Workspace with a connected PKr Server")
					return
				}
				opts := strings.ToLower(os.Args[3])
				switch opts {
				case "setup":
					{
						root.Server_Setup()
						return
					}

					// Can Name better in future
				case "register_workspace":
					{
						root.Server_RegisterWorkspace()
						return
					}
				default:
					{
						fmt.Printf("Incorrect Argument %s provided...\n\n", opts)
						fmt.Println("Valid Arguments:")
						fmt.Println("	1] setup -> Initialize Connection with New PKr Server")
						fmt.Println("	2] register_workspace -> Initialize An Existing Workspace with a connected PKr Server")
					}
				}
			}
		default:
			{
				fmt.Printf("Incorrect Parameter %s provided...\n\n", cmd)
				fmt.Println("Valid Parameters:")
				fmt.Println("	1] install -> Create User and Install PKr")
				fmt.Println("	2] init -> Initialize a Workspace, allows other Users to connect")
				fmt.Println("	3] clone -> Clone an existing Workspace of a different User")
				fmt.Println("	4] list -> List all Send and Get Workspaces")
				fmt.Println("	5] server -> Connect With a Server to Manage Multiple Dynamic Connections")
			}
		}
	}

	if CLI {
		// TODO: [ ] Required %s Parameter than display usage when flags are not provided, by checking if val == ""
		switch cmd {
		case "install":
			{
				var username string
				flag.StringVar(&username, "-u", "", "Username to Install PKr")
				flag.Parse()
				
				err := root.Install(username)
				if err != nil {
					fmt.Println("Could not Install PKr.")
					fmt.Println(err)
					return
				}

				fmt.Printf(" ~ Created User : %s\n", username)
				return
			}

		case "uninstall":

		case "get":

		case "push":

		case "clone":
			{
				var workspace_ip string
				var workspace_name string
				var workspace_password string

				flag.StringVar(&workspace_ip, "ip", "", "(*) Clone Workspace IP")
				flag.StringVar(&workspace_name, "wn", "", "(*) Workspace Name")
				flag.StringVar(&workspace_password, "wp", "", "(*) Workspace Password ")

				err := root.Clone(workspace_ip, workspace_name, workspace_password)
				if err != nil {
					fmt.Printf("Error Occured in Cloning Workspace: %s at IP: %s\n", workspace_name, workspace_ip)
					fmt.Println(err)
					return
				}

				fmt.Printf("Successfully Cloned Workspace: %s\n", workspace_name)
				return

			}

		// Maybe its Done
		case "init":
			{
				var workspace_password string
				flag.StringVar(&workspace_password, "wp", "", "(*) Workspace Password ")
				
				if err := root.Init(workspace_password); err != nil {
					fmt.Println("Error Occured in Initialize a New Workspace")
					fmt.Printf("error: %v\n", err)
					return
				}

				fmt.Println("Workspace Created Successfully !!")
				return
			}

			// list Created workspaces
		case "list":
			{
				if err := root.List(); err != nil {
					fmt.Println("Could Not List Workspace Info")
					fmt.Printf("Error: %v", err)
					return
				}
			}

		// For server
		// Mainly for IP and shit
		// Try to make the code as swappable as possible
		//
		// Was working on server setup ~ check models, it is still partial
		case "server":
			{
				if len(os.Args) < 4 {
					fmt.Println("server requires additional arguments")
					fmt.Println("Valid Arguments:")
					fmt.Println("	1] setup -> Initialize Connection with New PKr Server")
					fmt.Println("	2] register_workspace -> Initialize An Existing Workspace with a connected PKr Server")
					return
				}
				opts := strings.ToLower(os.Args[3])
				switch opts {
				case "setup":
					{
						root.Server_Setup()
						return
					}

					// Can Name better in future
				case "register_workspace":
					{
						root.Server_RegisterWorkspace()
						return
					}
				default:
					{
						fmt.Printf("Incorrect Argument %s provided...\n\n", opts)
						fmt.Println("Valid Arguments:")
						fmt.Println("	1] setup -> Initialize Connection with New PKr Server")
						fmt.Println("	2] register_workspace -> Initialize An Existing Workspace with a connected PKr Server")
					}
				}
			}
		default:
			{
				fmt.Printf("Incorrect Parameter %s provided...\n\n", cmd)
				fmt.Println("Valid Parameters:")
				fmt.Println("	1] install -> Create User and Install PKr")
				fmt.Println("	2] init -> Initialize a Workspace, allows other Users to connect")
				fmt.Println("	3] clone -> Clone an existing Workspace of a different User")
				fmt.Println("	4] list -> List all Send and Get Workspaces")
				fmt.Println("	5] server -> Connect With a Server to Manage Multiple Dynamic Connections")
			}
		}
	}
}
