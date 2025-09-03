package main

import (
	"fmt"
	"log"
	"net/http"
	"server_initializer/model"
	"server_initializer/router"
)

func main() {
	var a model.D
	for i := 3001; i < 3004; i++ {
		r := router.Router()
		host := fmt.Sprintf(":%d", i)

		go func(h string, rr http.Handler) {
			fmt.Printf("Starting server on %s\n", h)
			if err := http.ListenAndServe(h, rr); err != nil {
				log.Fatalf("Failed to start server on %s: %v", h, err)
				a[i] = false
			} else {
				a[i] = true
			}
		}(host, r)
	}
	select {}
}
