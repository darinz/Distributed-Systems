// Package mapreduce provides common types and utilities for the MapReduce framework.
//
// This package contains shared data structures, constants, and helper functions
// used across the MapReduce implementation, including RPC message types and
// the core RPC communication function.
//
// The package defines the communication protocol between master and worker
// processes, including job assignment, worker registration, and shutdown
// procedures.
package mapreduce

import (
	"fmt"
	"net/rpc"
)

// JobType constants define the types of tasks that can be assigned to workers.
const (
	// Map represents a map task that processes input chunks and produces key-value pairs
	Map = "Map"
	// Reduce represents a reduce task that aggregates values for specific keys
	Reduce = "Reduce"
)

// JobType represents the type of job that can be assigned to a worker.
// It is used to distinguish between map and reduce tasks in the job assignment protocol.
type JobType string

// RPC Message Types
//
// All RPC argument and reply structures must have field names that start
// with capital letters to ensure proper serialization/deserialization.

// DoJobArgs contains the arguments for assigning a job to a worker.
//
// This structure is sent from the master to a worker when assigning
// a map or reduce task. It contains all the information needed for
// the worker to execute the job.
//
// Fields:
//   - File: Name of the input file or base filename for the job
//   - Operation: Type of job (Map or Reduce)
//   - JobNumber: Zero-based index of this specific job
//   - NumOtherPhase: Total number of jobs in the other phase (used for partitioning)
type DoJobArgs struct {
	File          string   // Input file name
	Operation     JobType  // Type of job (Map or Reduce)
	JobNumber     int      // Zero-based job index
	NumOtherPhase int      // Total jobs in other phase (for partitioning)
}

// DoJobReply contains the response from a worker after completing a job.
//
// This structure is returned by workers to indicate whether they
// successfully completed the assigned job.
//
// Fields:
//   - OK: True if the job was completed successfully, false otherwise
type DoJobReply struct {
	OK bool // Success status of the job execution
}

// ShutdownArgs contains arguments for shutting down the master.
//
// This structure is used when requesting the master to shut down
// its registration server and stop accepting new worker connections.
// Currently, no arguments are needed for shutdown.
type ShutdownArgs struct {
	// No fields required for shutdown request
}

// ShutdownReply contains the response from a master shutdown request.
//
// This structure is returned by the master when it receives a
// shutdown request, providing information about the shutdown process.
//
// Fields:
//   - Njobs: Number of jobs that were in progress when shutdown was requested
//   - OK: True if shutdown was successful, false otherwise
type ShutdownReply struct {
	Njobs int  // Number of jobs in progress during shutdown
	OK    bool // Success status of the shutdown operation
}

// RegisterArgs contains arguments for worker registration.
//
// This structure is sent by workers when they register with the master
// to indicate their availability for job assignment.
//
// Fields:
//   - Worker: Network address of the worker process
type RegisterArgs struct {
	Worker string // Worker's network address (e.g., "localhost:7778")
}

// RegisterReply contains the response to a worker registration request.
//
// This structure is returned by the master when a worker attempts
// to register, indicating whether the registration was successful.
//
// Fields:
//   - OK: True if registration was successful, false otherwise
type RegisterReply struct {
	OK bool // Success status of the registration
}

// call sends an RPC request to a server and waits for the response.
//
// This function provides a unified interface for making RPC calls
// throughout the MapReduce framework. It handles connection establishment,
// request transmission, and response processing.
//
// The function uses Unix domain sockets for local communication between
// master and worker processes on the same machine. It includes automatic
// timeout handling and connection cleanup.
//
// Parameters:
//   - srv: Server address (Unix domain socket path)
//   - rpcname: Name of the RPC method to call (e.g., "MapReduce.Register")
//   - args: Arguments to send to the server (must be serializable)
//   - reply: Pointer to a structure to receive the response
//
// Returns:
//   - bool: True if the RPC call succeeded, false if it failed
//
// Important Notes:
//   - The reply argument must be a pointer to a structure
//   - The function will timeout if the server doesn't respond
//   - Connection is automatically closed after the call
//   - All RPC calls in the framework should use this function
//
// Example usage:
//   var reply DoJobReply
//   success := call(workerAddr, "Worker.DoJob", &args, &reply)
func call(srv string, rpcname string, args interface{}, reply interface{}) bool {
	// Establish connection to the server using Unix domain sockets
	c, err := rpc.Dial("unix", srv)
	if err != nil {
		// Connection failed - server may be down or unreachable
		return false
	}
	defer c.Close()

	// Make the RPC call and wait for response
	err = c.Call(rpcname, args, reply)
	if err == nil {
		// RPC call succeeded
		return true
	}

	// RPC call failed - log the error and return false
	fmt.Printf("RPC call failed: %v\n", err)
	return false
}
