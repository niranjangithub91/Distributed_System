package main

import (
	"fmt"
	"log"
	"net/http"
	"user_entry/router"
)

func main() {
	fmt.Println("User interface server")
	r := router.Router()
	log.Fatal(http.ListenAndServe(":3000", r))
	return
}
