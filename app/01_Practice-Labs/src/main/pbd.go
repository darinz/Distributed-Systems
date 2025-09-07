// Package main implements the primary/backup daemon (pbd).
// This is a standalone application that runs a primary/backup key-value server
// that participates in the distributed key-value system.
//
// The primary/backup server is responsible for:
//   - Maintaining a local key-value store
//   - Replicating data between primary and backup servers
//   - Handling client requests (Get, Put, PutHash operations)
//   - Coordinating with the view service for role management
//   - Ensuring consistency through primary/backup replication
//
// Usage: pbd <viewport> <myport>
//   viewport: Unix socket path of the view service
//   myport:   Unix socket path where this server will listen
//
// Example:
//   ./pbd /tmp/viewservice-socket /tmp/pbserver-1
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"distributed-systems/app/01_Practice-Labs/src/pbservice"
)

func main() {
	// Validate command line arguments
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <viewport> <myport>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  viewport: Unix socket path of the view service\n")
		fmt.Fprintf(os.Stderr, "  myport:   Unix socket path where this server will listen\n")
		fmt.Fprintf(os.Stderr, "Example: %s /tmp/viewservice-socket /tmp/pbserver-1\n", os.Args[0])
		os.Exit(1)
	}

	viewport := os.Args[1]
	myport := os.Args[2]

	fmt.Printf("Starting primary/backup daemon\n")
	fmt.Printf("  View service: %s\n", viewport)
	fmt.Printf("  Server port:  %s\n", myport)

	// Start the primary/backup server
	pbServer := pbservice.StartServer(viewport, myport)
	if pbServer == nil {
		fmt.Fprintf(os.Stderr, "Failed to start primary/backup server on %s\n", myport)
		os.Exit(1)
	}

	fmt.Printf("Primary/backup daemon started successfully on %s\n", myport)
	fmt.Printf("Press Ctrl+C to stop the server\n")

	// Set up graceful shutdown handling
	setupGracefulShutdown(pbServer)

	// Keep the server running
	keepAlive()
}

// setupGracefulShutdown configures signal handling for graceful shutdown.
// It listens for SIGINT (Ctrl+C) and SIGTERM signals and shuts down
// the primary/backup server cleanly when received.
func setupGracefulShutdown(pbServer *pbservice.PBServer) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		fmt.Printf("\nReceived signal %v, shutting down primary/backup server...\n", sig)
		pbServer.Kill()
		fmt.Println("Primary/backup daemon stopped")
		os.Exit(0)
	}()
}

// keepAlive keeps the main goroutine alive while the server runs.
// This function runs in an infinite loop, sleeping periodically to avoid
// consuming CPU resources unnecessarily.
func keepAlive() {
	for {
		time.Sleep(100 * time.Second)
	}
}
