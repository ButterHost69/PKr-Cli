package services

import (
	"ButterHost69/PKr-base/pb"
	"context"
	
	"google.golang.org/grpc/peer"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type BackgroundServer struct {
	pb.UnimplementedBackgroundServiceServer
}



func (s *BackgroundServer) GetPublicKey(ctx context.Context, message *emptypb.Empty)(*pb.PublicKey, error) {
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