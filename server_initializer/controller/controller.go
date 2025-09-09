package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"server_initializer/model"
	"strconv"
)

func Data_receive(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	filename := r.Header.Get("X-Filename")
	Port := r.Header.Get("Port-number")
	user := r.Header.Get("Username")
	chunk := r.Header.Get("Chunk_num")
	fmt.Println(Port)
	if err != nil {
		http.Error(w, "Error reading chunk", http.StatusInternalServerError)
		return
	}
	p := "./chunks/" + Port + "/" + user + "/" + chunk
	os.MkdirAll(p, os.ModePerm)
	os.WriteFile(fmt.Sprintf("%s/%s", p, filename), data, 0644)
	fmt.Println("Server")
	return
}

func Data_retreive(w http.ResponseWriter, r *http.Request) {
	fmt.Println("QWERTYUIOP")

	var a model.X
	err := json.NewDecoder(r.Body).Decode(&a)
	if err != nil {
		http.Error(w, "Invalid", http.StatusBadRequest)
		return
	}

	file_name := a.File_name
	username := a.Username
	port := a.Port
	details := a.Details

	var u [][]byte

	for _, val := range details {
		p := "./chunks/" + strconv.Itoa(port) + "/" + username + "/" + strconv.Itoa(val) + "/" + file_name
		file, err := os.Open(p)
		if err != nil {
			fmt.Println("File not found:", err)
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}

		data, err := io.ReadAll(file)
		file.Close()
		if err != nil {
			fmt.Println("Error reading file:", err)
			http.Error(w, "Error reading file", http.StatusInternalServerError)
			return
		}

		u = append(u, data)
	}

	// return as JSON array of base64 strings
	fmt.Println("Returning chunks:", len(u))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(u); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	return
}
