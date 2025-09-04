// package main

// import (
// 	"fmt"
// 	"log"
// 	"net/http"
// 	faulttolerance "user_entry/helper/Fault_tolerance"
// 	"user_entry/router"
// )

// func main() {
// 	fmt.Println("User interface server")
// 	go func() {
// 		b := faulttolerance.StartHealthMonitor()
// 		if !b {
// 			//Shut down the server
// 		}
// 	}()
// 	r := router.Router()
// 	log.Fatal(http.ListenAndServe(":3000", r))
// 	return
// }

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	faulttolerance "user_entry/helper/Fault_tolerance"
	"user_entry/router"
)

func main() {
	fmt.Println("User interface server")

	// create server instance
	srv := &http.Server{
		Addr:    ":3000",
		Handler: router.Router(),
	}

	// run health monitor in background
	go func() {
		b := faulttolerance.StartHealthMonitor()
		if !b {
			fmt.Println("Health monitor failed, shutting down server...")
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := srv.Shutdown(ctx); err != nil {
				log.Printf("Server Shutdown Failed:%+v", err)
			}
		}
	}()

	// also handle Ctrl+C to stop server gracefully
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		fmt.Println("\nReceived interrupt, shutting down...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("Server Shutdown Failed:%+v", err)
		}
	}()

	// start server (blocking)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe(): %s", err)
	}
}
