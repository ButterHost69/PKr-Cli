package services

import (
	"ButterHost69/PKr-base/encrypt"
	"ButterHost69/PKr-base/models"
	"ButterHost69/PKr-base/pb"
	"context"
	"time"

	"google.golang.org/grpc/peer"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type BackgroundServer struct {
	pb.UnimplementedBackgroundServiceServer
}



func (s *BackgroundServer) GetPublicKey(ctx context.Context, request *emptypb.Empty)(*pb.PublicKey, error) {
	keyData, err := ReadPublicKey(); 
	p, _ := peer.FromContext(ctx)
  	ip := p.Addr.String()
	if err != nil {
		logentry := "Could Not Provide Public Key To IP: " + ip
		AddUserLogEntry(logentry)
		AddUserLogEntry(err)

		return &pb.PublicKey{
			Key: "",
		}, err
	}
	logentry := "Successfully Provided Public Key To IP: " + ip
	AddUserLogEntry(logentry)

	return &pb.PublicKey{
		Key: keyData,
	}, nil

}

func (s *BackgroundServer) InitNewWorkSpaceConnection (ctx context.Context, request *pb.InitRequest)(*pb.InitResponse, error){
	// 1. Decrypt password [X]
	// 2. Authenticate Request [X]
	// 3. Add the New Connection to the .PKr Config File [X]
	// 4. Store the Public Key [X]
	// 5. Send the Response with port []
	// 6. Open a Data Transfer Port and shit [Will be a separate Function not here] []
	
	p, _ := peer.FromContext(ctx)
  	ip := p.Addr.String()

	encrypted_password := request.Password
	password, err := encrypt.DecryptData(encrypted_password)
	if err != nil {
		AddUserLogEntry("Failed to Init Workspace Connection for User IP: " + ip)
		AddUserLogEntry(err)	
		return &pb.InitResponse{
			Response: 4000,
			Port: 0000,
		}, nil
	}

	// Authenticates Workspace Name and Password and Get the Workspace File Path
	file_path := models.AuthenticateWorkspaceInfo(request.WorkspaceName, password)
	if file_path == "" {
		AddUserLogEntry("Failed to Init Workspace Connection for User IP: " + ip)
		return &pb.InitResponse{
			Response: 4000,
			Port: 0000,
		}, nil
	}

	var connection models.Connection
	connection.CurrentIP = ip
	connection.CurrentPort = request.Port

	// Save Public Key
	keysPath, err := models.StorePublicKeys(file_path + "\\keys\\", request.PublicKey)
	if err != nil {
		AddUserLogEntry("Failed to Init Workspace Connection for User IP: " + ip)
		AddUserLogEntry(err)
		return &pb.InitResponse{
			Response: 4000,
			Port: 0000,
		}, err
	}	
	
	// Store the New Connection in the .PKr Config file
	connection.PublicKeyPath = keysPath
	if err := models.AddConnectionToPKRConfigFile(file_path + "\\workspaceConfig.json", connection); err != nil {
		AddUserLogEntry("Failed to Init Workspace Connection for User IP: " + ip)
		AddUserLogEntry(err)
		return &pb.InitResponse{
			Response: 4000,
			Port: 0000,
		}, err
	}

	// Start New Data grpc Server
	port, err := StartDataServer(1024 * time.Second, file_path)
	// 
}
