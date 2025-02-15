package dialer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/ButterHost69/PKr-cli/config"
)

// Todo: [ ] Add record in Server Json File
func RegisterServer(server_ip, username, password string) (string, error) {
	url := "http://" + server_ip + "/register/user"
	method := "POST"

	// username := ctx.PostForm("username")
	// password := ctx.PostForm("password")

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	err := writer.WriteField("username", username)
	if err != nil {
		return "", fmt.Errorf("error writing field: %v", err)
	}

	err = writer.WriteField("password", password)
	if err != nil {
		return "", fmt.Errorf("error writing field: %v", err)
	}

	err = writer.Close()
	if err != nil {
		return "", fmt.Errorf("error in closing writer: %v", err)
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return "", fmt.Errorf("error failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error failed to make send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error failed to read from the response: %v", err)
	}

	var repsonse config.RegisterResp
	err = json.Unmarshal(body, &repsonse)
	if err != nil {
		return "", fmt.Errorf("error failed to umarshall repsonse: %v", err)
	}
	if resp.Status != "200 OK" && repsonse.Response == "success" {
		return "", fmt.Errorf("error Expected Status: 200 OK  ||  Body: 'response':'success,\nreceived: Status: %s, Body: %s", resp.Status, string(body))
	}

	return repsonse.Username, nil
}

func RegisterWorkspace(server_ip, username, password, workspace_name string) error {
	url := "http://" + server_ip + "/register/workspace"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	err := writer.WriteField("username", username)
	if err != nil {
		return fmt.Errorf("error writing field: %v", err)
	}

	err = writer.WriteField("password", password)
	if err != nil {
		return fmt.Errorf("error writing field: %v", err)
	}

	err = writer.WriteField("workspace_name", workspace_name)
	if err != nil {
		return fmt.Errorf("error writing field: %v", err)
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("error in closing writer: %v", err)
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return fmt.Errorf("error failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error failed to make send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error failed to read from the response: %v", err)
	}

	// t.Logf("Response status: %v", resp.Status)
	// t.Logf("Response body: %s", body)

	var repsonse config.GenericResp
	err = json.Unmarshal(body, &repsonse)
	if err != nil {
		return fmt.Errorf("Error failed to umarshall repsonse: %v", err)
	}
	if resp.Status != "200 OK" && repsonse.Response == "success" {
		return fmt.Errorf("Error Expected Status: 200 OK  ||  Body: 'response':'success,\nreceived: Status: %s, Body: %s", resp.Status, string(body))
	}

	return nil
}
