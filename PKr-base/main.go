// Listening Server
// Listen For other Connections
// Responsible to Create the Server that will Send Data

package main

import (
	"ButterHost69/PKr-base/models"
	"ButterHost69/PKr-base/pb"
	"ButterHost69/PKr-base/services"
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
	IP_ADDR	string
)

func Init(){
	IP_ADDR = os.Getenv("PKR-IP")
	if IP_ADDR == "" {
		IP_ADDR = ":9000"
	}
}

func main(){
	Init()

	lis, err := net.Listen("tcp", IP_ADDR)
	if err != nil {
		services.AddUserLogEntry(err)
	}

	grpcServer := grpc.NewServer()
	backgroundService := services.BackgroundServer{}
	
	pb.RegisterBackgroundServiceServer(grpcServer, &backgroundService)
	fmt.Println("Server Started")

	// TODO: [ ] Test this code, neither human test nor code test done....
	// All The functions written with it are not tested
	go func ()  {
		services.AddUserLogEntry("Update me Service Started")	
		for {
			serverList, err := models.GetAllServers()
			if err != nil {
				services.AddUserLogEntry(err)	
			}

			for _, server := range serverList{
				services.SendUpdateMeRequest(server.ServerIP, server.Username, server.Password)
			}
			time.Sleep(5 * time.Minute)
		}
	}()
	if err := grpcServer.Serve(lis); err != nil {
		services.AddUserLogEntry(err)
		services.AddUserLogEntry("Closing Listening Server")
		os.Exit(1)
	}

}