package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"server_initializer/router"
)

func main() {
	port := flag.Int("port", -1, "port to run the server on")
	flag.Parse()

	// if port not provided, exit
	if *port == -1 {
		fmt.Println("Error: You must provide a port using -port")
		os.Exit(1)
	}

	r := router.Router()
	host := fmt.Sprintf(":%d", *port)

	fmt.Printf("Starting independent server on %s\n", host)
	if err := http.ListenAndServe(host, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
