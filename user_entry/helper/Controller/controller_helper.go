package controller_helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	db "user_entry/helper/Database"
	"user_entry/model"
)

func Validate_signup(t model.Data) bool {
	if password, ok := t["password"].(string); ok {
		if len(password) < 5 {
			return false
		}
	}
	return true
}

func SendChunk(url string, chunk []byte, filename string, port int, name string) error {
	body := bytes.NewReader(chunk)
	fmt.Println(port)
	// Create POST request
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/octet-stream") // raw bytes
	req.Header.Set("X-Filename", filename)
	req.Header.Set("Port-number", strconv.Itoa(port))
	req.Header.Set("Username", name)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	a := resp.StatusCode
	if a != http.StatusAccepted {

	}
	respData, _ := io.ReadAll(resp.Body)
	fmt.Printf("Response from %s: %s\n", url, string(respData))

	//DB update
	var q model.Collection
	q.Name = name
	q.File_name = filename
	db.Collection_update(q)
	return nil
}

func Send_Chunk_request(url string, file_name string, username string, port int) (error, []byte) {
	var output []byte
	data := map[string]string{
		"username":  username,
		"file_name": file_name,
		"port":      strconv.Itoa(port),
	}

	// Convert Go map into JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err, output
	}

	// Send POST request with JSON body
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error sending request:", err)
		return err, output
	}
	defer resp.Body.Close()

	q, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return err, output
	}
	return nil, q
}
