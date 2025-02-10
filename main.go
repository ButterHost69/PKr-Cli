package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ButterHost69/PKr-cli/logger"
	"github.com/ButterHost69/PKr-cli/root"
)

// TODO: [ ] While Zipping File or Unzipping files, empty directory is ignored (Moit)
// TODO: [ ] Shift everything to flag based, support terminal inputs and CLI - Server Commands Remaining
// TODO: [X] Refactor the code.
//
// TODO: [ ] Why the fuck are there two model files ?? Make it in to 1
// TODO: [ ] Why are there print statements in files other than main.
// TODO: [ ] Write Tests, bro why am I doing this manual... Use docker maybe to simulate the whole thing ???
// TODO: [ ] Add verbose option that allows for prints in "inside" functions (anything aside from main)
// TODO: [ ] Use a Better Logging Method(preferbly homemade), also logs are made partially, not everything is logged.

var (
	TUI                    bool
	CLI                    bool
	BACKGROUND_SERVER_PORT int
	LOG_IN_TERMINAL        bool
	LOG_LEVEL              int
)

var (
	workspace_logger   *logger.WorkspaceLogger
	userconfing_logger *logger.UserLogger
)

const (
	ROOT_DIR = "..\\tmp"
	LOG_FILE = ROOT_DIR + "\\logs.txt"
)

func Init() {
	// flag.IntVar(&BACKGROUND_SERVER_PORT, "ip", 9000, "Other Users BACKGROUND Port")
	value := os.Getenv("PKR-IP")
	if value == "" {
		value = ":9000"
	}
	BACKGROUND_SERVER_PORT, _ = strconv.Atoi(value)

	flag.BoolVar(&TUI, "tui", false, "Use Application in TUI Mode")
	flag.BoolVar(&CLI, "cli", false, "Use Application in CLI Mode")
	flag.BoolVar(&LOG_IN_TERMINAL, "lt", false, "Log Events in Terminal.")
	flag.IntVar(&LOG_LEVEL, "ll", 4, "Set Log Levels.") // 4 -> No Logs

	flag.Parse()

	workspace_logger = logger.InitWorkspaceLogger()
	userconfing_logger = logger.InitUserLogger(LOG_FILE)

	workspace_logger.SetLogLevel(logger.IntToLog(LOG_LEVEL))
	userconfing_logger.SetLogLevel(logger.IntToLog(LOG_LEVEL))

	workspace_logger.SetPrintToTerminal(LOG_IN_TERMINAL)
	userconfing_logger.SetPrintToTerminal(LOG_IN_TERMINAL)
}

func PrintMode() {
	fmt.Println("Must Define Mode to use PKr in")
	fmt.Println("	1] `PKr -tui` -> For Terminal User Interface. Takes Input through stdout, requires less flags")
	fmt.Println("	2] `PKr -cli` -> For Command Line Interface. Requires Input as flags")
}

// [ ]: Add PUSH Command Info
func PrintArguments() {
	fmt.Printf("Required Minimum 3 Args\n\n")
	fmt.Println("Valid Parameters:")
	fmt.Println("	1] install -> Create User and Install PKr")
	fmt.Println("	2] init -> Initialize a Workspace, allows other Users to connect")
	fmt.Println("	3] clone -> Clone an existing Workspace of a different User")
	fmt.Println("	4] list -> List all Send and Get Workspaces")
	fmt.Println("	5] server -> Connect With a Server to Manage Multiple Dynamic Connections")
}

func PrintServerOptions() {
	fmt.Println("server requires additional arguments")
	fmt.Println("Valid Arguments:")
	fmt.Println("	1] setup -> Initialize Connection with New PKr Server")
	fmt.Println("	2] register_workspace -> Initialize An Existing Workspace with a connected PKr Server")
}

func main() {
	Init()

	// TUI -> Takes Input from stdin and print output/errors on stdout
	// CLI -> Input Passed as flag with command, output/errors on stdout
	if !TUI && !CLI {
		PrintMode()
		return
	}

	if len(os.Args) < 3 {
		PrintArguments()
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

		// [ ] Test this code thoroughly - Use VMWare Maybe ??
		case "push":
			{
				dir, err := os.Getwd()
				if err != nil {
					fmt.Println("Could not Get Directory Name: ")
					fmt.Println(err)
					return
				}

				workspace_namel := strings.Split(dir, "\\")
				workspace_name := workspace_namel[len(workspace_namel)-1]
				fmt.Println("Pushing Workpace: ", workspace_name)

				success, err := root.Push(workspace_name, workspace_logger)
				if err != nil {
					fmt.Printf("Error Occured in Pushing Workspace: %s\n", workspace_name)
					fmt.Println(err)
					return
				}

				fmt.Printf("\nNotified %d Users !!\n", success)
			}

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

				err := root.Clone(BACKGROUND_SERVER_PORT, workspace_ip, workspace_name, workspace_password)
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
					PrintServerOptions()
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
						PrintServerOptions()
						return
					}
				}
			}
		default:
			{
				PrintArguments()
				return
			}
		}
	}

	if CLI {
		// TODO: [X] Required %s Parameter than display usage when flags are not provided, by checking if val == ""
		switch cmd {
		case "install":
			{
				installCmd := flag.NewFlagSet("install", flag.ExitOnError)
				username := installCmd.String("u", "", "Username to Install PKr")

				installCmd.Parse(os.Args[3:])
				if *username == "" {
					fmt.Println("Error: Username is required for install")
					fmt.Println(`Usage: PKr -cli install -u="username"`)
					return
				}

				fmt.Println("Creating User: ", *username)
				err := root.Install(*username)
				if err != nil {
					fmt.Println("Could not Install PKr.")
					fmt.Println(err)
					return
				}

				fmt.Printf(" ~ Created User : %s\n", *username)
				return
			}

		case "uninstall":

		case "get":

		case "push":

		case "clone":
			{
				cloneCmd := flag.NewFlagSet("clone", flag.ExitOnError)

				workspace_ip := cloneCmd.String("ip", "", "(*) Clone Workspace IP")
				workspace_name := cloneCmd.String("wn", "", "(*) Workspace Name")
				workspace_password := cloneCmd.String("wp", "", "(*) Workspace Password ")

				cloneCmd.Parse(os.Args[3:])
				if *workspace_ip == "" && *workspace_name == "" && *workspace_password == "" {
					fmt.Println("Error: Workspace IP and Name and Password required")
					cloneCmd.Usage()
					return
				}

				err := root.Clone(BACKGROUND_SERVER_PORT, *workspace_ip, *workspace_name, *workspace_password)
				if err != nil {
					fmt.Printf("Error Occured in Cloning Workspace: %s at IP: %s\n", *workspace_name, *workspace_ip)
					fmt.Println(err)
					return
				}

				fmt.Printf("Successfully Cloned Workspace: %s\n", *workspace_name)
				return

			}

		// Maybe its Done
		case "init":
			{
				initCmd := flag.NewFlagSet("init", flag.ExitOnError)
				workspace_password := initCmd.String("wp", "", "(*) Workspace Password ")

				initCmd.Parse(os.Args[3:])
				if *workspace_password == "" {
					fmt.Println("Error: Workspace Password required")
					initCmd.Usage()
					return
				}

				if err := root.Init(*workspace_password); err != nil {
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
		// TODO: [ ] CLI remaining For server part. Flags are not taken
		case "server":
			{
				if len(os.Args) < 4 {
					PrintServerOptions()
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
						PrintServerOptions()
					}
				}
			}
		default:
			{
				PrintArguments()
			}
		}
	}
}
