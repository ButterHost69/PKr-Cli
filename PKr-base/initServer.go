package PKr_base

import (
	"net"

	"github.com/ButterHost69/PKr-cli/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// Responsible for listening to for other connections trying to establish a workspace connection
// Will Handle -> Establishing Connections with new Users
// This Means -> Verify New User that Are Connecting to the Workspace
// 			  -> Verifying Passwords
// 			  -> Receiving Keys and Shit...

//  Requester -> [Sends Username, Workspace_name, Password, Public Key] -> Reciever
//  Reciever ->  (IF ALL CREDENTIALS Correct) -> {
// 													Store Username, IP, ?? [In the Parent Folder],
//		   											 and Same under [Connections in userConfig.json ???],
//      	                                         Log this Transactions in [Parent,  and in tmp],
// 													Store Public Key in parent Folder
// 												 }

// Files 	-> tmp/userConfig.json
// 		 	-> workspace/.PKr/ -> Keys/
// 							   -> connections.json
// 							   -> LOG


type Server struct {

}


func InitConnectionServer() (error){
	lis, err := net.Listen("tcp", ":4040")
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer() 
	if err := grpcServer.Serve(lis) ; err != nil {
		return err 
	}

	return nil
}


func (s *Server) InitNewWorkSpaceConnection (ctx context.Context, request *pb.Request) (*pb.Response, error){
	
}

