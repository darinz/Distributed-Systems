// Package main implements the view service daemon (viewd).
// This is a standalone application that runs the view service server
// for the primary/backup key-value system.
//
// The view service is responsible for:
//   - Tracking which servers are alive or dead
//   - Managing view transitions (primary/backup assignments)
//   - Ensuring at most one primary is active at a time
//   - Coordinating server promotions and demotions
//
// Usage: viewd <port>
//   port: Unix socket path where the view service will listen
//
// Example:
//   ./viewd /tmp/viewservice-socket
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"distributed-systems/app/01_Practice-Labs/src/viewservice"
)

func main() {
	// Validate command line arguments
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <port>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  port: Unix socket path where the view service will listen\n")
		fmt.Fprintf(os.Stderr, "Example: %s /tmp/viewservice-socket\n", os.Args[0])
		os.Exit(1)
	}

	port := os.Args[1]
	fmt.Printf("Starting view service daemon on %s\n", port)

	// Start the view service server
	vs := viewservice.StartServer(port)
	if vs == nil {
		fmt.Fprintf(os.Stderr, "Failed to start view service on %s\n", port)
		os.Exit(1)
	}

	fmt.Printf("View service daemon started successfully on %s\n", port)
	fmt.Printf("Press Ctrl+C to stop the service\n")

	// Set up graceful shutdown handling
	setupGracefulShutdown(vs)

	// Keep the service running
	keepAlive()
}

// setupGracefulShutdown configures signal handling for graceful shutdown.
// It listens for SIGINT (Ctrl+C) and SIGTERM signals and shuts down
// the view service cleanly when received.
func setupGracefulShutdown(vs *viewservice.ViewServer) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		fmt.Printf("\nReceived signal %v, shutting down view service...\n", sig)
		vs.Kill()
		fmt.Println("View service daemon stopped")
		os.Exit(0)
	}()
}

// keepAlive keeps the main goroutine alive while the view service runs.
// This function runs in an infinite loop, sleeping periodically to avoid
// consuming CPU resources unnecessarily.
func keepAlive() {
	for {
		time.Sleep(100 * time.Second)
	}
}
