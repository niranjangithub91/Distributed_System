package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"server_initializer/model"
)

func Data_receive(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	filename := r.Header.Get("X-Filename")
	Port := r.Header.Get("Port-number")
	user := r.Header.Get("Username")
	fmt.Println(Port)
	if err != nil {
		http.Error(w, "Error reading chunk", http.StatusInternalServerError)
		return
	}
	p := "./chunk/" + Port + "/" + user
	os.MkdirAll(p, os.ModePerm)
	os.WriteFile(fmt.Sprintf("%s/%s", p, filename), data, 0644)
	fmt.Println("Server")
	return
}

func Data_retreive(w http.ResponseWriter, r *http.Request) {
	var a model.Data
	err := json.NewDecoder(r.Body).Decode(&a)
	if err != nil {
		http.Error(w, "Invalid", http.StatusBadGateway)
		return
	}
	file_name := a["file_name"].(string)
	username := a["username"].(string)
	port := a["port"].(string)
	p := fmt.Sprintf("./chunk/%s/%s/%s", port, username, file_name)
	file, err := os.Open(p)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	defer file.Close()
}
