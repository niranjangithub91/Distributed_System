package controller_helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

func SendChunk(url string, chunk []byte, filename string, port int, name string, chunk_num int) error {
	fmt.Println("Hi")
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
	req.Header.Set("Chunk_num", strconv.Itoa(chunk_num))

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
	q.Chunk = chunk_num
	q.Port = port
	db.Collection_update(q)
	return nil
}

func Send_Chunk_request(url string, file_name string, username string, port int, z map[int][]int) (error, map[int][]byte) {
	fmt.Println("thfgyrnkldac")
	r := z[port]
	if len(r) == 0 {
		return nil, make(map[int][]byte) // return empty initialized map
	}

	data := map[string]interface{}{
		"username":  username,
		"file_name": file_name,
		"port":      port,
		"details":   r,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err), nil
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error sending request: %w", err), nil
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("RAW RESPONSE:", string(body))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned %d: %s", resp.StatusCode, string(body)), nil
	}

	var chunks [][]byte
	if err := json.Unmarshal(body, &chunks); err != nil {
		return fmt.Errorf("error decoding JSON: %w", err), nil
	}

	u := make(map[int][]byte)
	counters := 0
	for _, val := range r {
		if counters >= len(chunks) {
			break
		}
		u[val] = chunks[counters]
		fmt.Println("Chunk", counters, "=>", string(chunks[counters]))
		counters++
	}

	fmt.Println("qwsdbcjdcbwjhcwjcvewjcvedcvejcve")

	return nil, u
}
