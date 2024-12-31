// Listening Server
// Listen For other Connections
// Responsible to Create the Server that will Send Data

package main

import (
	"ButterHost69/PKr-base/dialer"
	"ButterHost69/PKr-base/models"
	"ButterHost69/PKr-base/pb"
	"ButterHost69/PKr-base/services"
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
)

// const (
// 	ROOT_DIR     = "..\\tmp"
// 	MY_KEYS_PATH = ROOT_DIR + "\\mykeys"
// 	CONFIG_FILE  = ROOT_DIR + "\\userConfig.json"
// 	LOG_FILE = ROOT_DIR + "\\logs.txt"
// 	SERVER_LOG_FILE = ROOT_DIR + "\\serverlogs.txt"
// )

var (
	IP_ADDR string
)


// README: Pkr-base.exe -ip :9001
func Init() {
	flag.StringVar(&IP_ADDR, "ip", "", "Use Application in TUI Mode")
	flag.Parse()

	fmt.Println("IP_ADDR:", IP_ADDR)

	if IP_ADDR != "" {
		return
	}

	IP_ADDR = os.Getenv("PKR-IP")
	if IP_ADDR == "" {
		IP_ADDR = ":9000"
	}

	models.UpdateBasePort(IP_ADDR)
}

// TODO: [ ] Write "Push" Command notification server
func main() {
	Init()

	lis, err := net.Listen("tcp", IP_ADDR)
	if err != nil {
		fmt.Println("Error: ", err) // [ ]: For Debugging Only
	}

	grpcServer := grpc.NewServer()
	backgroundService := services.BackgroundServer{}

	pb.RegisterBackgroundServiceServer(grpcServer, &backgroundService)
	fmt.Println("Server Started")
	fmt.Println("Server running on: " + IP_ADDR)

	// TODO: [ ] Test this code, neither human test nor code test done....
	// All The functions written with it are not tested
	go func() {
		services.AddUserLogEntry("Update me Service Started")

		for {
			// Read Each Time... So can automatically detect changes without manual anything....
			serverList, err := models.GetAllServers()
			if err != nil {
				services.AddUserLogEntry(err)
			}

			// Quit For Loop if no Server list
			if len(serverList) == 0 {
				break
			}

			for _, server := range serverList {
				services.SendUpdateMeRequest(server.ServerIP, server.Username, server.Password)
			}
			time.Sleep(5 * time.Minute)
		}
	}()

	// [ ] Look for a better way to call this function instead of using go-routines
	dialer.ScanForUpdatesOnStart()

	if err := grpcServer.Serve(lis); err != nil {
		services.AddUserLogEntry(err)
		services.AddUserLogEntry("Closing Listening Server")
		os.Exit(1)
	}
}
