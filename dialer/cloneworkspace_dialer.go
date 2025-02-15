package dialer

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/ButterHost69/PKr-cli/config"
	"github.com/ButterHost69/PKr-cli/encrypt"
	"github.com/ButterHost69/PKr-cli/filetracker"
	"github.com/ButterHost69/PKr-cli/pb"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func GetPublicKey(workspace_ip string) ([]byte, error) {
	conn, err := grpc.NewClient(workspace_ip, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pb.NewBackgroundServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()

	response, err := client.GetPublicKey(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}

	return response.Key, nil
}

func InitNewWorkSpaceConnection(workspace_ip, workspace_name, workspace_password, port string, public_key []byte) (int, error) {
	conn, err := grpc.NewClient(workspace_ip, grpc.WithInsecure())
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	client := pb.NewBackgroundServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()

	response, err := client.InitNewWorkSpaceConnection(ctx, &pb.InitRequest{
		WorkspaceName: workspace_name,
		Username:      "",
		Password:      workspace_password,
		PublicKey:     public_key,
		Port:          port,
	})

	if err != nil {
		return 0, err
	}

	if response.Response == 4000 {
		return 0, errors.New("error: incorrect data/connection error")
	}

	return int(response.Port), nil
}

func getMyIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "nil", err
	}
	// defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}

func GetData(workspace_name, workspace_only_ip, port string, workspacePath string) error {
	// Get Data, Key, IV
	// Decrypt Key, IV
	// Decrypt Data
	// Store Zip
	// Clear the GetFolder except .PKr
	// Unzip data on the GetFolder
	// Store Last Hash into Config Files

	workspace_ip_with_port := workspace_only_ip + port
	my_ip, err := getMyIP()
	if err != nil {
		return err
	}

	// [ ] : Send Port, from UserConfig File
	my_ip = my_ip + ":9000"
	// Comment This Later
	fmt.Println("IP: " + my_ip)
	// ...
	conn, err := grpc.NewClient(workspace_ip_with_port, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewDataServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

	stream, err := client.GetData(ctx, &pb.DataRequest{
		WorkspaceName: workspace_name,
		ConnectionIp:  my_ip,
	})

	if err != nil {
		return err
	}

	var data_bytes []byte
	var key_bytes []byte
	var iv_bytes []byte

	var last_hash string
	for {
		data, err := stream.Recv()
		if err != nil {
			return err
		}

		if data.Filetype == 0 {
			data_bytes = append(data_bytes, data.Chunk...)
		} else if data.Filetype == 1 {
			key_bytes = append(key_bytes, data.Chunk...)
		} else if data.Filetype == 2 {
			iv_bytes = append(iv_bytes, data.Chunk...)
			break
		} else if data.Filetype == 3 {
			last_hash = string(data.Chunk)
		}
	}

	decrypted_key, err := encrypt.DecryptData(string(key_bytes))
	if err != nil {
		return err
	}

	decrypted_iv, err := encrypt.DecryptData(string(iv_bytes))
	if err != nil {
		return err
	}

	data, err := encrypt.AESDecrypt(data_bytes, decrypted_key, decrypted_iv)
	if err != nil {
		return err
	}

	zip_file_path := workspacePath + "\\.PKr\\" + last_hash + ".zip"
	if err = filetracker.SaveDataToFile(data, zip_file_path); err != nil {
		return err
	}

	if err = filetracker.CleanFilesFromWorkspace(workspacePath); err != nil {
		return err
	}

	// Unzip Content
	// unzip_file_path := getworkspace.WorkspacePath
	if err = filetracker.UnzipData(zip_file_path, workspacePath+"\\"); err != nil {
		return err
	}

	if err = config.AddGetWorkspaceFolderToUserConfig(workspace_name, workspacePath, workspace_ip_with_port, last_hash); err != nil {
		return fmt.Errorf("error in adding GetConnection to the Main User Config Folder.\nerror:%v", err)
	}

	return nil
}
