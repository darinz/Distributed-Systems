// Package mapreduce implements the master process for distributed MapReduce execution.
//
// This file contains the core logic for the MapReduce master, which coordinates
// the execution of map and reduce tasks across multiple worker processes.
// The master is responsible for:
//   - Accepting worker registrations
//   - Distributing map and reduce jobs to available workers
//   - Tracking job completion and handling worker failures
//   - Coordinating the overall MapReduce workflow
//
// The master uses a channel-based architecture to manage worker availability
// and job distribution, ensuring efficient parallel execution of tasks.
package mapreduce

import (
	"container/list"
	"fmt"
)

// WorkerInfo contains information about a registered worker process.
//
// This structure tracks the state and metadata for each worker that
// has registered with the master. It is used for job assignment,
// failure tracking, and cleanup operations.
//
// Fields:
//   - address: Network address of the worker process (e.g., "localhost:7778")
//
// Additional fields can be added here to track worker state, such as:
//   - Last heartbeat time
//   - Number of completed jobs
//   - Worker capabilities or load
//   - Failure count
type WorkerInfo struct {
	address string // Network address of the worker process
	// Additional worker metadata can be added here
}

// KillWorkers gracefully shuts down all registered workers and collects statistics.
//
// This method sends a shutdown RPC to each registered worker, requesting
// them to terminate gracefully. It collects the number of jobs each worker
// has performed and returns this information as a list.
//
// The method is typically called at the end of a MapReduce job to clean up
// worker processes and gather execution statistics.
//
// Returns:
//   - *list.List: List of integers representing the number of jobs completed by each worker
//
// Note: Workers that fail to respond to the shutdown request are logged
// but do not cause the method to fail.
func (mr *MapReduce) KillWorkers() *list.List {
	l := list.New()
	
	for _, w := range mr.Workers {
		DPrintf("DoWork: shutdown %s\n", w.address)
		
		args := &ShutdownArgs{}
		var reply ShutdownReply
		ok := call(w.address, "Worker.Shutdown", args, &reply)
		
		if !ok {
			fmt.Printf("DoWork: RPC %s shutdown error\n", w.address)
		} else {
			l.PushBack(reply.Njobs)
		}
	}
	
	return l
}

// RunMaster orchestrates the execution of all map and reduce tasks.
//
// This is the main coordination function that manages the entire MapReduce
// workflow. It starts worker registration handling, distributes map tasks,
// waits for map completion, distributes reduce tasks, waits for reduce
// completion, and finally shuts down all workers.
//
// The method uses goroutines to enable parallel execution of tasks while
// using channels to synchronize completion and manage worker availability.
//
// Execution Flow:
//   1. Start worker registration handler
//   2. Launch all map tasks in parallel
//   3. Wait for all map tasks to complete
//   4. Launch all reduce tasks in parallel
//   5. Wait for all reduce tasks to complete
//   6. Shutdown all workers and collect statistics
//
// Returns:
//   - *list.List: Statistics about jobs completed by each worker
func (mr *MapReduce) RunMaster() *list.List {
	// Start the worker registration handler in a separate goroutine
	// This allows workers to register at any time during execution
	go mr.registerWorkers()

	// Phase 1: Execute all map tasks in parallel
	for i := 0; i < mr.nMap; i++ {
		go mr.delegateJob(Map, i)
	}
	
	// Wait for all map tasks to complete
	for i := 0; i < mr.nMap; i++ {
		<-mr.doneChannel // Block until a map task completes
	}

	// Phase 2: Execute all reduce tasks in parallel
	for i := 0; i < mr.nReduce; i++ {
		go mr.delegateJob(Reduce, i)
	}
	
	// Wait for all reduce tasks to complete
	for i := 0; i < mr.nReduce; i++ {
		<-mr.doneChannel // Block until a reduce task completes
	}

	// Phase 3: Cleanup - shutdown all workers and collect statistics
	return mr.KillWorkers()
}

// registerWorkers continuously accepts worker registrations and makes them available for jobs.
//
// This method runs in a separate goroutine and handles the registration
// of new workers. When a worker registers, it is added to the worker
// registry and placed on the ready channel, making it available for
// job assignment.
//
// The method runs indefinitely, allowing workers to register at any time
// during the MapReduce execution. This supports dynamic worker scaling
// and handles workers that may start up after the master begins execution.
//
// Registration Process:
//   1. Wait for a worker to register via registerChannel
//   2. Create WorkerInfo entry for the worker
//   3. Add worker to the ready channel for job assignment
func (mr *MapReduce) registerWorkers() {
	for {
		// Wait for a worker to register
		address := <-mr.registerChannel
		
		// Create worker info entry for shutdown and tracking
		mr.Workers[address] = &WorkerInfo{address: address}
		
		// Make the worker available for job assignment
		mr.readyChannel <- address
	}
}

// delegateJob assigns a specific job to an available worker and handles retries.
//
// This method is responsible for finding an available worker and assigning
// a specific map or reduce job to it. If the assignment fails (e.g., worker
// is unresponsive), it will retry with the next available worker.
//
// The method implements fault tolerance by automatically retrying failed
// job assignments with different workers, ensuring that all jobs eventually
// complete successfully.
//
// Parameters:
//   - jtype: Type of job (Map or Reduce)
//   - jno: Zero-based job number
//
// Execution Flow:
//   1. Wait for an available worker from readyChannel
//   2. Attempt to assign the job to the worker
//   3. If successful: signal completion and return worker to ready pool
//   4. If failed: retry with next available worker
func (mr *MapReduce) delegateJob(jtype JobType, jno int) {
	// Retry loop: continue until job is successfully completed
	for {
		// Wait for an available worker
		worker := <-mr.readyChannel
		
		// Attempt to assign the job to the worker
		if ok := mr.assignJob(worker, jtype, jno); ok {
			// Job completed successfully
			// Signal completion to the main coordination loop
			mr.doneChannel <- true
			
			// Return the worker to the available pool
			mr.readyChannel <- worker
			
			// Job is complete, exit the retry loop
			return
		}
		
		// Job assignment failed - worker may be unresponsive
		// Continue loop to try with the next available worker
	}
}

// assignJob sends a job assignment RPC to a specific worker.
//
// This method constructs the appropriate job arguments based on the job type
// and sends an RPC request to the specified worker. The arguments include
// all necessary information for the worker to execute the job.
//
// Parameters:
//   - worker: Network address of the worker to assign the job to
//   - jtype: Type of job (Map or Reduce)
//   - jno: Zero-based job number
//
// Returns:
//   - bool: True if the RPC call succeeded and the worker accepted the job
//
// Job Arguments:
//   - Map jobs: File name, job number, and number of reduce tasks (for partitioning)
//   - Reduce jobs: File name, job number, and number of map tasks (for input collection)
func (mr *MapReduce) assignJob(worker string, jtype JobType, jno int) bool {
	var args DoJobArgs
	
	// Construct job arguments based on job type
	switch jtype {
	case Map:
		// Map job: needs to know how many reduce tasks exist for partitioning
		args = DoJobArgs{
			File:          mr.file,
			Operation:     Map,
			JobNumber:     jno,
			NumOtherPhase: mr.nReduce,
		}
	case Reduce:
		// Reduce job: needs to know how many map tasks exist for input collection
		args = DoJobArgs{
			File:          mr.file,
			Operation:     Reduce,
			JobNumber:     jno,
			NumOtherPhase: mr.nMap,
		}
	}
	
	var reply DoJobReply
	return call(worker, "Worker.DoJob", args, &reply)
}
