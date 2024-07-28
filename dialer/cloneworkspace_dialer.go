package dialer

import (
	"archive/zip"
	"context"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/ButterHost69/PKr-cli/encrypt"
	"github.com/ButterHost69/PKr-cli/models"
	"github.com/ButterHost69/PKr-cli/pb"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func GetPublicKey(workspace_ip string) ([]byte, error) {
	conn, err := grpc.NewClient(workspace_ip)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pb.NewBackgroundServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute * 5)	
	defer cancel()

	response, err := client.GetPublicKey(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}

	return response.Key, nil
}

func InitNewWorkSpaceConnection(workspace_ip, workspace_name, workspace_password, port string, public_key []byte)(int, error){
	conn, err := grpc.NewClient(workspace_ip)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	client := pb.NewBackgroundServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute * 5)
	defer cancel()

	response, err := client.InitNewWorkSpaceConnection(ctx, &pb.InitRequest{
		WorkspaceName: workspace_name,
		Username: "",
		Password: workspace_password,
		PublicKey: public_key,
		Port: port,
	})
	
	if err != nil {
		return 0, err
	}

	if response.Response == 4000 {
		return 0, errors.New("error: incorrect data/connection error")
	}

	return int(response.Port), nil
}

func getMyIP()(string, error){
	conn, err := net.Dial("udp","8.8.8.8:80")
	if err != nil {
		return "nil", err
	}
	// defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}

func UnzipData(src, dest string)(error) {
	zipped_data, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	zipped_data.Close()

	for _, file := range zipped_data.File {
		if file.FileInfo().IsDir() {
			continue
		}

		filepath := filepath.Join(dest, file.Name)
		if err := os.MkdirAll(filepath, 0777); err != nil {
			return nil
		}
		
		f, err := os.Create(filepath)
		if err != nil {
			return nil
		}

		content, err := file.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(f, content)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetData(workspace_name, workspace_ip, port string)(error){
	// Get Data, Key, IV
	// Decrypt Key, IV
	// Decrypt Data
	// Store Zip
	// Clear the GetFolder except .PKr
	// Unzip data on the GetFolder

	new_addr := workspace_ip + port
	my_ip, err := getMyIP()
	if err != nil {
		return err
	}

	conn, err := grpc.NewClient(new_addr)	
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewDataServiceClient(conn)
	
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

	stream, err := client.GetData(ctx, &pb.DataRequest{
		WorkspaceName: workspace_name,
		ConnectionIp: my_ip,
	})
	
	if err != nil {
		return err
	}

	var data_bytes 	[]byte
	var key_bytes	[]byte
	var iv_bytes	[]byte
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
		}
	}

	decrypted_key, err := encrypt.DecryptData(string(key_bytes))
	if err != nil {
		return err
	}

	decrypted_iv , err:= encrypt.DecryptData(string(iv_bytes))
	if err != nil {
		return err
	}

	data, err := encrypt.AESDecrypt(data_bytes, decrypted_key, decrypted_iv)
	if err != nil {
		return err
	}

	getworkspace, err := models.GetGetWorkspaceFolder(workspace_name)
	if err != nil {
		return err
	}

	zipFileName := strings.Split(time.Now().String(), " ")[0] + ".zip"
	zip_file_path := getworkspace.WorkspacePath + "\\.PKr\\" + zipFileName
	if err = os.WriteFile(zip_file_path, data, os.ModeAppend); err != nil {
		return err
	}

	// Delete all files in the GetWorkspace Dir except for .PKr
	files, err := ioutil.ReadDir(getworkspace.WorkspaceName)
	if err != nil {
		return err
	}

	for _, file:= range files{
		if file.Name() != ".PKr" {
			if err = os.RemoveAll(path.Join([]string{getworkspace.WorkspaceName, file.Name()}...)); err != nil {
				return err
			}
		}
	}

	// Unzip Content
	// unzip_file_path := getworkspace.WorkspacePath
	if err = UnzipData(zip_file_path, getworkspace.WorkspaceName); err != nil {
		return err
	}

	return nil
}
