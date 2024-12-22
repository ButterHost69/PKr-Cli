package root

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/ButterHost69/PKr-cli/encrypt"
	"github.com/ButterHost69/PKr-cli/models"
)

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

	addFilesToZip(writer, workspace_path+"\\", "")

	if err = writer.Close(); err != nil {
		return "", err
	}

	return workspace_path + "\\.PKr\\" + zipFileName, nil
}


// [ ] Generate Hash of Zip Files
// [ ] Add to their .PKr file
// [ ] Notify all Connections 

// TODO:
// 	[ ] Test the Entire code nothing is tested

// LATER UPDATE:
// 		[ ] Do the Zip and Create Hash in Memory before saving it
// 		[ ]
func Push(workspace_name string) error {
	workspace_path, err := models.GetWorkspaceFilePath(workspace_name)
	if err != nil {
		return fmt.Errorf("could find workspace.\nError: %v", err)
	}
	zipfile, err := ZipData(workspace_path)
	if err != nil {
		return fmt.Errorf("could not zip data.\nError: %v", err)
	}

	fmt.Println("[Log Delete Later] Zipfile Path: ", zipfile)

	generate_hash, err := encrypt.GenerateHash(zipfile)
	if err != nil {
		return fmt.Errorf("could hash file data: %s.\nError: %v", zipfile, err)
	}

	//  [ ] Rename Zip file to hash name
	err = models.AddNewPushToConfig(workspace_name, generate_hash)
	if err != nil {
		return fmt.Errorf("could not zip data.\nError: %v", err)
	}

	
	return nil
	// generate_sha1 :=
}