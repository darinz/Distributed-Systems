// Package mapreduce implements the worker process for distributed MapReduce execution.
//
// This file contains the worker implementation that executes map and reduce tasks
// assigned by the master. Workers register with the master, accept job assignments,
// execute the requested operations, and report completion back to the master.
//
// The worker process:
//   - Registers with the master upon startup
//   - Accepts RPC calls for job execution and shutdown
//   - Executes map and reduce tasks using user-provided functions
//   - Tracks job completion statistics
//   - Gracefully shuts down when requested
//
// Workers communicate with the master using Unix domain sockets for local
// communication and RPC for job coordination.
package mapreduce

import (
	"container/list"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
)

// Worker represents a MapReduce worker process that executes tasks assigned by the master.
//
// The Worker struct contains all the state and configuration needed for a worker
// to participate in the MapReduce framework. It maintains connections, tracks
// job statistics, and holds references to the user-provided map and reduce functions.
//
// Fields:
//   - name: Network address of this worker process
//   - Reduce: User-provided reduce function for processing key-value pairs
//   - Map: User-provided map function for processing input chunks
//   - nRPC: Number of RPC calls remaining before shutdown (-1 for infinite)
//   - nJobs: Number of jobs completed by this worker
//   - l: Network listener for accepting RPC connections
type Worker struct {
	name   string                                    // Worker's network address
	Reduce func(string, *list.List) string          // User-provided reduce function
	Map    func(string) *list.List                  // User-provided map function
	nRPC   int                                       // RPC calls remaining (-1 = infinite)
	nJobs  int                                       // Number of jobs completed
	l      net.Listener                             // Network listener for RPC connections
}

// DoJob executes a map or reduce task assigned by the master.
//
// This method is called by the master via RPC to assign a specific job
// to this worker. The worker executes the appropriate operation (map or reduce)
// using the user-provided functions and reports success back to the master.
//
// Parameters:
//   - arg: Job assignment arguments containing job details
//   - res: Reply structure to return job completion status
//
// Returns:
//   - error: Always nil (job execution errors are handled internally)
//
// The method logs job details for debugging and executes either DoMap
// or DoReduce based on the operation type specified in the arguments.
func (wk *Worker) DoJob(arg *DoJobArgs, res *DoJobReply) error {
	fmt.Printf("Dojob %s job %d file %s operation %v N %d\n",
		wk.name, arg.JobNumber, arg.File, arg.Operation,
		arg.NumOtherPhase)
	
	// Execute the appropriate operation based on job type
	switch arg.Operation {
	case Map:
		// Execute map task using user-provided map function
		DoMap(arg.JobNumber, arg.File, arg.NumOtherPhase, wk.Map)
	case Reduce:
		// Execute reduce task using user-provided reduce function
		DoReduce(arg.JobNumber, arg.File, arg.NumOtherPhase, wk.Reduce)
	}
	
	// Report successful completion
	res.OK = true
	return nil
}

// Shutdown gracefully terminates the worker and reports job statistics.
//
// This method is called by the master via RPC to request the worker to
// shut down. The worker reports the number of jobs it has completed
// and prepares for termination.
//
// Parameters:
//   - args: Shutdown arguments (unused)
//   - res: Shutdown reply containing job statistics
//
// Returns:
//   - error: Always nil
//
// The method adjusts the RPC counter to allow the shutdown RPC to complete
// and decrements the job counter to exclude the shutdown RPC itself.
func (wk *Worker) Shutdown(args *ShutdownArgs, res *ShutdownReply) error {
	DPrintf("Shutdown %s\n", wk.name)
	
	// Report the number of jobs completed by this worker
	res.Njobs = wk.nJobs
	res.OK = true
	
	// Allow the shutdown RPC to complete by setting nRPC to 1
	// This is safe because the same thread reads nRPC
	wk.nRPC = 1
	
	// Don't count the shutdown RPC as a completed job
	wk.nJobs--
	
	return nil
}

// Register notifies the master that this worker is available for job assignment.
//
// This function sends a registration RPC to the master, informing it that
// this worker is ready to accept and execute map or reduce tasks. The
// registration is a prerequisite for receiving job assignments.
//
// Parameters:
//   - master: Network address of the master process
//   - me: Network address of this worker process
//
// The function will log an error if the registration RPC fails, but will
// not terminate the worker process, allowing for retry mechanisms.
func Register(master string, me string) {
	args := &RegisterArgs{
		Worker: me,
	}
	var reply RegisterReply
	
	ok := call(master, "MapReduce.Register", args, &reply)
	if !ok {
		fmt.Printf("Register: RPC %s register error\n", master)
	}
}

// RunWorker starts a worker process and begins accepting job assignments.
//
// This function initializes a worker process, registers it with the master,
// and starts the main event loop to accept and execute RPC calls for
// job assignments and shutdown requests.
//
// Parameters:
//   - MasterAddress: Network address of the master process
//   - me: Network address of this worker process
//   - MapFunc: User-provided map function for processing input chunks
//   - ReduceFunc: User-provided reduce function for aggregating values
//   - nRPC: Number of RPC calls to accept before shutdown (-1 for infinite)
//
// The worker will:
//   1. Initialize its internal state and RPC server
//   2. Set up a Unix domain socket listener
//   3. Register with the master
//   4. Accept RPC calls until nRPC limit is reached or shutdown is requested
//   5. Clean up resources and exit
//
// Note: The RPC handling loop below should not be modified as it contains
// critical synchronization logic for job counting and shutdown handling.
func RunWorker(MasterAddress string, me string,
	MapFunc func(string) *list.List,
	ReduceFunc func(string, *list.List) string, nRPC int) {
	
	DPrintf("RunWorker %s\n", me)
	
	// Initialize worker with provided configuration
	wk := &Worker{
		name:   me,
		Map:    MapFunc,
		Reduce: ReduceFunc,
		nRPC:   nRPC,
		nJobs:  0,
	}
	
	// Set up RPC server for handling job assignments
	rpcs := rpc.NewServer()
	rpcs.Register(wk)
	
	// Remove any existing socket file (Unix domain sockets only)
	os.Remove(me)
	
	// Create Unix domain socket listener
	l, err := net.Listen("unix", me)
	if err != nil {
		log.Fatalf("RunWorker: worker %s error: %v", me, err)
	}
	wk.l = l
	
	// Register with the master to indicate availability
	Register(MasterAddress, me)

	// Main RPC handling loop - DO NOT MODIFY
	// This loop accepts RPC connections and processes job assignments
	// until the worker is shut down or the RPC limit is reached
	for wk.nRPC != 0 {
		conn, err := wk.l.Accept()
		if err == nil {
			// Decrement RPC counter and handle connection in separate goroutine
			wk.nRPC -= 1
			go rpcs.ServeConn(conn)
			wk.nJobs += 1
		} else {
			// Connection error - exit the loop
			break
		}
	}
	
	// Clean up resources
	wk.l.Close()
	DPrintf("RunWorker %s exit\n", me)
}
