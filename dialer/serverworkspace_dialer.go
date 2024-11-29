package dialer

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

// Todo: [ ] Add record in Server Json File
func RegisterServer(username, password, ip_addr string) error {
	url := "http://localhost:9069/whats/new"
	method := "GET"

	// username := ctx.PostForm("username")
	// password := ctx.PostForm("password")

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

	fmt.Printf("Response status: %v", resp.Status)
	fmt.Printf("Response body: %s", body)

	return nil
}