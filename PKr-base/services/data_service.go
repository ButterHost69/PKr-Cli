package services

import (
	"ButterHost69/PKr-base/encrypt"
	"ButterHost69/PKr-base/models"
	"ButterHost69/PKr-base/pb"
	"archive/zip"
	"fmt"
	"io"

	// "io/ioutil"
	"path/filepath"
	"strings"
	"sync"

	// "context"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"

	// "github.com/brianvoe/gofakeit/v7/source"
	"google.golang.org/grpc"
)

type DataServer struct {
	pb.UnimplementedDataServiceServer
	wg             *sync.WaitGroup
	workspace_path string
}

// I dont Know if this works. Check it later
// I copied From : https://gosamples.dev/zip-file/
// CHEECHK THIISS LAAATTERRR
func ZipData(workspace_path string) (string, error) {
	// 2024-01-20.zip
	zipFileName := strings.Split(time.Now().String(), " ")[0] + ".zip"

	file, err := os.Create(workspace_path + "\\.PKr\\" + zipFileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	writer := zip.NewWriter(file)

	// This Might Break in Linux...
	return zipFileName, filepath.Walk(workspace_path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			header, err := zip.FileInfoHeader(info)
			if err != nil {
				return err
			}

			header.Method = zip.Deflate

			header.Name, err = filepath.Rel(filepath.Dir(workspace_path), path)
			if err != nil {
				return err
			}

			if info.IsDir() {
				header.Name += "/"
			}

			headerWriter, err := writer.CreateHeader(header)
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			f, err := os.Open(path)
			if err != nil {
				return err
			}

			defer f.Close()

			_, err = io.Copy(headerWriter, f)
			return err
		})
}

func (server *DataServer) GetData(request *pb.DataRequest, stream pb.DataService_GetDataServer) error {

	// 1. Zip Data [x]
	//		|-> Check if already zipped ???? Later Not Now
	// 2. Encrypt Data [x] -> Check Ip or ID than find that public key
	// 		Encrypt using AES ; Encrypt AES using RSA ; SEND AES KEY ; SEND DATA
	// 3. Encrypt AES key and iv [X]
	// 4. Send / Stream Data [X]

	// workspace_path + "\\.PKr\\" + zipFileName
	zipped_filepath, err := ZipData(server.workspace_path)
	if err != nil {
		logdata := fmt.Sprintf("Could Not Zip The File\nError: %v", err)
		models.AddLogEntry(request.WorkspaceName, logdata)
		return err
	}

	key, err := encrypt.AESGenerakeKey(16)
	if err != nil {
		logdata := fmt.Sprintf("Could Not Generate AES Key\nError: %v", err)
		models.AddLogEntry(request.WorkspaceName, logdata)
		return err
	}

	iv, err := encrypt.AESGenerateIV()
	if err != nil {
		logdata := fmt.Sprintf("Could Not Generate AES IV\nError: %v", err)
		models.AddLogEntry(request.WorkspaceName, logdata)
		return err
	}

	destination_filepath := strings.Replace(zipped_filepath, ".zip", ".enc", 1)
	if err := encrypt.AESEncrypt(zipped_filepath, destination_filepath, key, iv); err != nil {
		logdata := fmt.Sprintf("Could Not Encrypt File\nError: %v", err)
		models.AddLogEntry(request.WorkspaceName, logdata)
		return err
	}

	publicKeyPath, err := models.GetConnectionsPublicKeyUsingIP(server.workspace_path, request.ConnectionIp)
	if err != nil {
		logdata := fmt.Sprintf("Could Not Find Users Public Key\nError: %v", err)
		models.AddLogEntry(request.WorkspaceName, logdata)
		return err
	}
	publicKey, err := os.ReadFile(publicKeyPath)
	if err != nil {
		logdata := fmt.Sprintf("Could Not Read Users Public Key\nError: %v", err)
		models.AddLogEntry(request.WorkspaceName, logdata)
		return err
	}

	encrypt_key, err := encrypt.EncryptData(string(key), string(publicKey))
	if err != nil {
		logdata := fmt.Sprintf("Could Not Encrypt Key\nError: %v", err)
		models.AddLogEntry(request.WorkspaceName, logdata)
		return err
	}

	encrypt_iv, err := encrypt.EncryptData(string(iv), string(publicKey))
	if err != nil {
		logdata := fmt.Sprintf("Could Not Encrypt IV\nError: %v", err)
		models.AddLogEntry(request.WorkspaceName, logdata)
		return err
	}

	encrypt_file, err := os.Open(destination_filepath)
	if err != nil {
		logdata := fmt.Sprintf("Could Not Open The Encrypted File\nError: %v", err)
		models.AddLogEntry(request.WorkspaceName, logdata)
		return err
	}
	defer encrypt_file.Close()

	buff := make([]byte, 256)
	chunknumber := 1

	// Send the File
	for {
		num, err := encrypt_file.Read(buff)
		if err == io.EOF {
			// If file Send Over then Send the Key and IV
			if err := stream.Send(&pb.Data{Filetype: 1, Chunk: []byte(encrypt_key)}); err != nil {
				logdata := fmt.Sprintf("Could Not Send Encrypted Key\nError: %v", err)
				models.AddLogEntry(request.WorkspaceName, logdata)
				return err
			}

			if err := stream.Send(&pb.Data{Filetype: 1, Chunk: []byte(encrypt_iv)}); err != nil {
				logdata := fmt.Sprintf("Could Not Send Encrypted IV\nError: %v", err)
				models.AddLogEntry(request.WorkspaceName, logdata)
				return err
			}

			break
		}

		if err != nil {
			logdata := fmt.Sprintf("Could Not Read Encrypted File\nError: %v", err)
			models.AddLogEntry(request.WorkspaceName, logdata)
			return err
		}

		chunk := buff[:num]
		if err := stream.Send(&pb.Data{Filetype: 0, Chunk: chunk}); err != nil {
			logdata := fmt.Sprintf("Could Not Send Encrypted File\nError: %v", err)
			models.AddLogEntry(request.WorkspaceName, logdata)
			return err
		}

		chunknumber += 1
	}
	logdata := fmt.Sprintf("File Sent SuccessFully to User IP: %v\nTotol Size: %v", request.ConnectionIp, chunknumber*256)
	models.AddLogEntry(request.WorkspaceName, logdata)

	server.wg.Done()
	return nil
}

func isPortInUse(port int) bool {
	conn, err := net.DialTimeout("tcp", ":"+strconv.Itoa(port), time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func StartDataServer(time_till_wait time.Duration, workspace_path string) (int, error) {
	// Pass the port using channels
	// Look into it later, i mean soon, after the DataServer soon[]

	rand.Seed(time.Now().UnixNano())

	var port int
	for {
		port = rand.Intn(65535-1024) + 1024
		if !isPortInUse(port) {
			break
		}
	}

	str_port := strconv.Itoa(port)
	lis, err := net.Listen("tcp", ":"+str_port)
	if err != nil {
		AddUserLogEntry(err)
	}
	var wg sync.WaitGroup
	wg.Add(1)

	grpcServer := grpc.NewServer()
	dataServer := DataServer{
		wg:             &wg,
		workspace_path: workspace_path,
	}

	pb.RegisterDataServiceServer(grpcServer, &dataServer)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			AddUserLogEntry(err)
			AddUserLogEntry("Closing Data Server")
			os.Exit(1)
		}
	}()

	wg.Wait()
	AddUserLogEntry("Data Transfer Done... Closing Data Server")
	grpcServer.GracefulStop()
}
