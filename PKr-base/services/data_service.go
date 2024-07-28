package services

import (
	"ButterHost69/PKr-base/encrypt"
	"ButterHost69/PKr-base/models"
	"ButterHost69/PKr-base/pb"
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	// "log"

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
// Running This Function Twice Makes a Zip File Whose Size keeps increasing until the Entire Disk
// is filled
// Dont USE THISSSSS
func ZipppData(workspace_path string) (string, error) {
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

			relPath, err := filepath.Rel(filepath.Dir(workspace_path), path)
			if err != nil {
				return err
			}
			header.Name = relPath

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

func addFilesToZip(writer *zip.Writer, dirpath string, relativepath string)(error){
	files, err := ioutil.ReadDir(dirpath)
	if err != nil {
		// log.Println(err)
		return err
	}

	for _, file := range(files) {
		// Comment This Later ... Only For Debugging 
		// models.AddUsersLogEntry(fmt.Sprintf("File: %s", file.Name()))
		// ..........
		if file.Name() == ".PKr" || file.Name() == "PKr-base.exe" || file.Name() == "PKr-cli.exe" || file.Name() == "tmp"{
			continue
		} else if !file.IsDir() {
			content, err := os.ReadFile(dirpath + file.Name())
			
			if err != nil {
				// log.Println(err)
				return err		
			}

			file, err := writer.Create(relativepath + file.Name())
			if err != nil {
				// log.Println(err)
				return err
			}
			file.Write(content)
		} else if file.IsDir() {
			newDirPath := dirpath + file.Name() + "\\"
			newRelativePath := relativepath + file.Name() + "\\"
			
			addFilesToZip(writer, newDirPath, newRelativePath)
		}
	}

	return nil
}

func ZipData(workspace_path string) (string, error) {
	zipFileName := strings.Split(time.Now().String(), " ")[0] + ".zip"

	zip_file, err := os.Create(workspace_path + "\\.PKr\\" + zipFileName)
	if err != nil {
		// models.AddLogEntry(workspace_name, err)
		return "", err
	}

	defer zip_file.Close()


	writer := zip.NewWriter(zip_file) 

	// cwd, err := os.Getwd() 
	// if err != nil {
	// 	log.Println(err)
	// 	return		
	// }

	addFilesToZip(writer, workspace_path + "\\" , "" )

	if err = writer.Close(); err != nil {
		return "", err
	}

	return zipFileName, nil
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

	
	models.AddLogEntry(request.WorkspaceName, "Workspace Zipped")

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
	
	models.AddLogEntry(request.WorkspaceName, "AES Keys Generated")

	zipped_filepath = server.workspace_path + "\\.PKr\\" + zipped_filepath
	destination_filepath := strings.Replace(zipped_filepath, ".zip", ".enc", 1)
	if err := encrypt.AESEncrypt(zipped_filepath, destination_filepath, key, iv); err != nil {
		logdata := fmt.Sprintf("Could Not Encrypt File\nError: %v\nFilePath: %v", err,zipped_filepath)
		models.AddLogEntry(request.WorkspaceName, logdata)
		return err
	}

	
	models.AddLogEntry(request.WorkspaceName, "Zip AES is Encrypted")

	// ## Uncomment Later ##
	publicKeyPath, err := models.GetConnectionsPublicKeyUsingIP(server.workspace_path , request.ConnectionIp)
	if err != nil {
		logdata := fmt.Sprintf("Could Not Find Users Public Key\nError: %v", err)
		models.AddLogEntry(request.WorkspaceName, logdata)
		return err
	}

	// ## Uncomment Later ##
	publicKey, err := os.ReadFile(publicKeyPath)
	if err != nil {
		logdata := fmt.Sprintf("Could Not Read Users Public Key\nError: %v", err)
		models.AddLogEntry(request.WorkspaceName, logdata)
		return err
	}

	// ## Uncomment Later ##
	encrypt_key, err := encrypt.EncryptData(string(key), string(publicKey))
	if err != nil {
		logdata := fmt.Sprintf("Could Not Encrypt Key\nError: %v", err)
		models.AddLogEntry(request.WorkspaceName, logdata)
		return err
	}

	// ## Uncomment Later ##
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


	models.AddLogEntry(request.WorkspaceName, "AES Keys Encrypted")

	buff := make([]byte, 2048)
	chunknumber := 1

	// logdata = fmt.Sprintf()
	models.AddLogEntry(request.WorkspaceName, "Sending Chunks...")

	// Send the File
	for {
		num, err := encrypt_file.Read(buff)
		// ## Uncomment Later ##
		if err == io.EOF {
			// If file Send Over then Send the Key and IV
			models.AddLogEntry(request.WorkspaceName, "Sending Encrypted Keys")
			if err := stream.Send(&pb.Data{Filetype: 1, Chunk: []byte(encrypt_key)}); err != nil {
				logdata := fmt.Sprintf("Could Not Send Encrypted Key\nError: %v", err)
				models.AddLogEntry(request.WorkspaceName, logdata)
				return err
			}

			models.AddLogEntry(request.WorkspaceName, "Sending Encrypted IV")
			if err := stream.Send(&pb.Data{Filetype: 2, Chunk: []byte(encrypt_iv)}); err != nil {
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
	logdata := fmt.Sprintf("File Sent SuccessFully to User IP: %v\nTotal Size: %v", request.ConnectionIp, chunknumber*256)
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

func StartDataServer(time_till_wait time.Duration, workspace_name string,workspace_path string, portchan chan int, errorchan chan error){
	// Pass the port using channels
	// Look into it later, i mean soon, after the DataServer soon[]

	// Generate a Random not Taken Port Number
	rand.Seed(time.Now().UnixNano())
	var port int
	for {
		port = rand.Intn(65535-1024) + 1024
		if !isPortInUse(port) {
			break
		}
	}

	// Register and Start gRPC Server
	str_port := strconv.Itoa(port)
	
	logdata := fmt.Sprintf("Port Number: %v Has Been Selected For Data Transfer", str_port)
	models.AddLogEntry(workspace_name, logdata)
	
	lis, err := net.Listen("tcp", ":"+str_port)
	if err != nil {
		models.AddLogEntry(workspace_name ,err)
		errorchan <- err
	}

	logdata = fmt.Sprintf("Started Listening to the Port: %v", str_port)
	models.AddLogEntry(workspace_name, logdata)

	var wg sync.WaitGroup
	wg.Add(1)

	grpcServer := grpc.NewServer()
	dataServer := DataServer{
		wg:             &wg,
		workspace_path: workspace_path,
	}

	pb.RegisterDataServiceServer(grpcServer, &dataServer)

	go func() {
		// Add Serve With Time Out Later
		if err := grpcServer.Serve(lis); err != nil {
			models.AddLogEntry(workspace_name,err)
			models.AddLogEntry(workspace_name, "Closing Data Server")
			errorchan <- err
			os.Exit(1)
		}
	}()

	logdata = fmt.Sprintf("gRPC Server Data Started on Port: %v", str_port)
	models.AddLogEntry(workspace_name, logdata)

	portchan <- port
	
	wg.Wait()
	models.AddLogEntry(workspace_name, "Data Transfer Done... Closing Data Server")
	grpcServer.GracefulStop()
}
