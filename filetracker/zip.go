package filetracker

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ButterHost69/PKr-cli/encrypt"
	"github.com/ButterHost69/PKr-cli/models"
)

// I dont Know if this works. Check it later
// I copied From : https://gosamples.dev/zip-file/
// CHEECHK THIISS LAAATTERRR
// Running This Function Twice Makes a Zip File Whose Size keeps increasing until the Entire Disk
// is filled
// Dont USE THISSSSS
func zippToInfiniteSize(workspace_path string) (string, error) {
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

func addFilesToZip(writer *zip.Writer, dirpath string, relativepath string) error {
	files, err := ioutil.ReadDir(dirpath)
	if err != nil {
		// log.Println(err)
		return err
	}

	for _, file := range files {
		// Comment This Later ... Only For Debugging
		// models.AddUsersLogEntry(fmt.Sprintf("File: %s", file.Name()))
		// ..........
		if file.Name() == ".PKr" || file.Name() == "PKr-base.exe" || file.Name() == "PKr-cli.exe" || file.Name() == "tmp" {
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
	fullZipPath := workspace_path + "\\.PKr\\" + zipFileName

	zip_file, err := os.Create(fullZipPath)
	if err != nil {
		// models.AddLogEntry(workspace_name, err)
		return "", err
	}

	fmt.Println("OLD Zip File Path: " + fullZipPath)

	writer := zip.NewWriter(zip_file)

	addFilesToZip(writer, workspace_path+"\\", "")

	if err = writer.Close(); err != nil {
		return "", err
	}

	hashFileName, err := encrypt.GenerateHashWithFileIO(zip_file)
	// hashFileName, err := encrypt.GenerateHashWithFilePath(fullZipPath)
	if err != nil {
		return "", err
	}

	zip_file.Close()
	
	hashFileName = hashFileName + ".zip"
	fullHashFilePath := workspace_path + "\\.PKr\\" + hashFileName

	workspace_split := strings.Split(workspace_path, "\\")
	workspace_name := workspace_split[len(workspace_split) - 1]
	
	if err = os.Rename(fullZipPath, fullHashFilePath); err != nil {
		logdata := fmt.Sprintf("could rename zip file to new hash name: %s | zipped file path: %s.\nError: %v", fullHashFilePath, fullZipPath, err)
		models.AddLogEntry(workspace_name, logdata)
		return "", err
	}

	fmt.Println("New Zip File Path: " + fullHashFilePath)
	return hashFileName, nil
}