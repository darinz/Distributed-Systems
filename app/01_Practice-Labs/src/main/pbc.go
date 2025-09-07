// Package main implements the primary/backup client (pbc).
// This is a command-line client application for interacting with the
// primary/backup key-value service.
//
// The client supports two operations:
//   - Get: Retrieve the value for a given key
//   - Put: Store a key-value pair
//
// Usage:
//   pbc <viewport> <key>                    # Get value for key
//   pbc <viewport> <key> <value>            # Put key-value pair
//
// Arguments:
//   viewport: Unix socket path of the view service
//   key:      The key to get or set
//   value:    The value to set (for Put operation only)
//
// Examples:
//   ./pbc /tmp/viewservice-socket key1                    # Get value for key1
//   ./pbc /tmp/viewservice-socket key1 value1             # Set key1=value1
//
// Building and Running the System:
//   go build viewd.go    # Build view service daemon
//   go build pbd.go      # Build primary/backup daemon
//   go build pbc.go      # Build client
//
//   ./viewd /tmp/viewservice-socket &                     # Start view service
//   ./pbd /tmp/viewservice-socket /tmp/pbserver-1 &       # Start server 1
//   ./pbd /tmp/viewservice-socket /tmp/pbserver-2 &       # Start server 2
//   ./pbc /tmp/viewservice-socket key1 value1             # Put operation
//   ./pbc /tmp/viewservice-socket key1                    # Get operation
//
// Fault Tolerance Testing:
//   Start multiple pbd programs in separate terminals and kill/restart
//   them to test the system's fault tolerance capabilities.
package main

import (
	"fmt"
	"os"

	"distributed-systems/app/01_Practice-Labs/src/pbservice"
)

func main() {
	// Validate command line arguments and execute appropriate operation
	switch len(os.Args) {
	case 3:
		// Get operation: pbc viewport key
		executeGetOperation()
	case 4:
		// Put operation: pbc viewport key value
		executePutOperation()
	default:
		// Invalid number of arguments
		usage()
	}
}

// usage displays the usage information and exits the program.
func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s <viewport> <key>\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "       %s <viewport> <key> <value>\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "Arguments:\n")
	fmt.Fprintf(os.Stderr, "  viewport: Unix socket path of the view service\n")
	fmt.Fprintf(os.Stderr, "  key:      The key to get or set\n")
	fmt.Fprintf(os.Stderr, "  value:    The value to set (for Put operation only)\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "Examples:\n")
	fmt.Fprintf(os.Stderr, "  %s /tmp/viewservice-socket key1\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s /tmp/viewservice-socket key1 value1\n", os.Args[0])
	os.Exit(1)
}

// executeGetOperation performs a Get operation to retrieve a value for a key.
func executeGetOperation() {
	viewport := os.Args[1]
	key := os.Args[2]

	// Create client and connect to the primary/backup service
	ck := pbservice.MakeClerk(viewport, "")
	if ck == nil {
		fmt.Fprintf(os.Stderr, "Failed to create client for viewport %s\n", viewport)
		os.Exit(1)
	}

	// Perform Get operation
	value := ck.Get(key)
	
	// Output the result
	fmt.Printf("%s\n", value)
}

// executePutOperation performs a Put operation to store a key-value pair.
func executePutOperation() {
	viewport := os.Args[1]
	key := os.Args[2]
	value := os.Args[3]

	// Create client and connect to the primary/backup service
	ck := pbservice.MakeClerk(viewport, "")
	if ck == nil {
		fmt.Fprintf(os.Stderr, "Failed to create client for viewport %s\n", viewport)
		os.Exit(1)
	}

	// Perform Put operation
	ck.Put(key, value)

	// Confirm successful operation
	fmt.Printf("Successfully stored %s=%s\n", key, value)
}
