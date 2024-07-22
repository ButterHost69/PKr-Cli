// Listening Server
// Listen For other Connections
// Responsible to Create the Server that will Send Data

package main

import (
	"ButterHost69/PKr-base/pb"
	"ButterHost69/PKr-base/services"
	"fmt"
	"net"
	"os"

	"google.golang.org/grpc"
)

// const (
// 	ROOT_DIR     = "..\\tmp"
// 	MY_KEYS_PATH = ROOT_DIR + "\\mykeys"
// 	CONFIG_FILE  = ROOT_DIR + "\\userConfig.json"
// 	LOG_FILE = ROOT_DIR + "\\logs.txt"
// 	SERVER_LOG_FILE = ROOT_DIR + "\\serverlogs.txt"
// )



func main(){
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		services.AddUserLogEntry(err)
	}

	grpcServer := grpc.NewServer()
	backgroundService := services.BackgroundServer{}
	
	pb.RegisterBackgroundServiceServer(grpcServer, &backgroundService)
	fmt.Println("Server Started")
	if err := grpcServer.Serve(lis); err != nil {
		services.AddUserLogEntry(err)
		services.AddUserLogEntry("Closing Listening Server")
		os.Exit(1)
	}

}