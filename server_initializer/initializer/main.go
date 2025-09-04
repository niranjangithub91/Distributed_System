package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
)

func main() {
	// Step 1: Build the server binary
	serverBinary := "./myserver"
	buildCmd := exec.Command("go", "build", "-o", serverBinary, "../myserver.go")
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	fmt.Println("Building server binary...")
	if err := buildCmd.Run(); err != nil {
		log.Fatalf("Failed to build server: %v", err)
	}

	// Step 2: Start multiple server processes
	var procs []*exec.Cmd
	for i := 3001; i < 3004; i++ {
		cmd := exec.Command(serverBinary, "-port", fmt.Sprintf("%d", i))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Start(); err != nil {
			log.Fatalf("Failed to start server on port %d: %v", i, err)
		}
		procs = append(procs, cmd)
		fmt.Printf("Launched server on port %d with PID %d\n", i, cmd.Process.Pid)
	}

	// Step 3: Handle Ctrl+C (SIGINT) or kill (SIGTERM)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan // block until user presses Ctrl+C

	// Step 4: Cleanup (stop servers + delete binary)
	fmt.Println("\nShutting down servers...")
	for _, p := range procs {
		if p.Process != nil {
			_ = p.Process.Kill()
			fmt.Printf("Killed PID %d\n", p.Process.Pid)
		}
	}
	os.Remove(filepath.Clean(serverBinary))
	fmt.Println("Cleaned up: servers stopped and binary deleted.")
}
