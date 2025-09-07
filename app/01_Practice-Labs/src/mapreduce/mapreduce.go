// Package mapreduce provides a distributed MapReduce framework implementation.
//
// This package implements the MapReduce programming model for processing
// large datasets in parallel across multiple worker processes. It provides
// both sequential and distributed execution modes.
//
// The MapReduce framework consists of:
//   - Master: Coordinates the overall job execution and manages workers
//   - Workers: Execute map and reduce tasks assigned by the master
//   - File System: Shared storage for intermediate and final results
//
// Key Components:
//   - Split: Divides input files into smaller chunks for parallel processing
//   - Map: Processes input chunks and produces key-value pairs
//   - Reduce: Aggregates values for each key across all map outputs
//   - Merge: Combines reduce outputs into a single final result
//
// Example usage:
//   mr := MakeMapReduce(nMap, nReduce, inputFile, masterAddress)
//   <-mr.DoneChannel // Wait for completion
package mapreduce

import (
	"bufio"
	"container/list"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"log"
	"net"
	"net/rpc"
	"os"
	"sort"
	"strconv"
)

// import "os/exec" // Reserved for future process management features

// MapReduce Framework Overview:
//
// The MapReduce framework processes large datasets by dividing work into
// parallel map and reduce phases:
//
// 1. Split Phase:
//    Input file is divided into nMap chunks:
//    - mrtmp.<filename>-0, mrtmp.<filename>-1, ..., mrtmp.<filename>-<nMap-1>
//
// 2. Map Phase:
//    Each map task processes one input chunk and produces nReduce intermediate files:
//    - mrtmp.<filename>-<mapJob>-0, mrtmp.<filename>-<mapJob>-1, ...,
//      mrtmp.<filename>-<mapJob>-<nReduce-1>
//    Total intermediate files: nMap × nReduce
//
// 3. Reduce Phase:
//    Each reduce task processes all intermediate files for its partition:
//    - Collects mrtmp.<filename>-*-<reduceJob> from all map tasks
//    - Produces mrtmp.<filename>-res-<reduceJob>
//
// 4. Merge Phase:
//    Combines all reduce outputs into final result: mrtmp.<filename>
//
// The framework supports both sequential execution (single process) and
// distributed execution (multiple worker processes coordinated by a master).

// Debug controls the verbosity of debug output.
// Set to 1 or higher to enable debug logging.
const Debug = 0

// DPrintf provides conditional debug logging.
// Only prints debug messages when Debug > 0.
//
// Parameters:
//   - format: Printf-style format string
//   - a: Arguments for the format string
//
// Returns:
//   - n: Number of bytes written
//   - err: Any error that occurred during writing
func DPrintf(format string, a ...interface{}) (n int, err error) {
	if Debug > 0 {
		n, err = fmt.Printf(format, a...)
	}
	return
}

// KeyValue represents a key-value pair used throughout the MapReduce framework.
// Both Map and Reduce functions work with these pairs to process data.
//
// Fields:
//   - Key: The key component of the pair (typically a string identifier)
//   - Value: The value component of the pair (typically data associated with the key)
type KeyValue struct {
	Key   string // The key component
	Value string // The value component
}

// MapReduce represents the master process that coordinates a MapReduce job.
// It manages worker registration, job distribution, and result collection.
//
// The master maintains state about the job configuration, active workers,
// and communication channels for coordinating the distributed execution.
type MapReduce struct {
	// Job configuration
	nMap          int    // Number of map tasks to create
	nReduce       int    // Number of reduce tasks to create
	file          string // Input file name
	MasterAddress string // Network address for the master process

	// Communication channels
	registerChannel chan string // Channel for worker registration
	DoneChannel     chan bool   // Channel to signal job completion
	readyChannel    chan string // Channel for available workers
	doneChannel     chan bool   // Channel for task completion notifications

	// Runtime state
	alive bool           // Whether the master is still running
	l     net.Listener   // Network listener for RPC connections
	stats *list.List     // Statistics about job execution

	// Worker management
	Workers map[string]*WorkerInfo // Map of registered worker addresses to info
}

// InitMapReduce initializes a new MapReduce master with the specified configuration.
//
// This function creates and initializes all necessary data structures and
// communication channels for the MapReduce job, but does not start the
// registration server or begin job execution.
//
// Parameters:
//   - nmap: Number of map tasks to create
//   - nreduce: Number of reduce tasks to create
//   - file: Input file name to process
//   - master: Network address for the master process
//
// Returns:
//   - *MapReduce: Initialized MapReduce master instance
func InitMapReduce(nmap int, nreduce int, file string, master string) *MapReduce {
	mr := &MapReduce{
		nMap:            nmap,
		nReduce:         nreduce,
		file:            file,
		MasterAddress:   master,
		alive:           true,
		registerChannel: make(chan string),
		DoneChannel:     make(chan bool),
		Workers:         make(map[string]*WorkerInfo),
		readyChannel:    make(chan string),
		doneChannel:     make(chan bool),
	}
	return mr
}

// MakeMapReduce creates and starts a new MapReduce master.
//
// This is the main entry point for creating a distributed MapReduce job.
// It initializes the master, starts the registration server to accept
// worker connections, and begins job execution in a separate goroutine.
//
// Parameters:
//   - nmap: Number of map tasks to create
//   - nreduce: Number of reduce tasks to create
//   - file: Input file name to process
//   - master: Network address for the master process
//
// Returns:
//   - *MapReduce: Started MapReduce master instance
//
// The master will run asynchronously. Use mr.DoneChannel to wait for completion.
func MakeMapReduce(nmap int, nreduce int, file string, master string) *MapReduce {
	mr := InitMapReduce(nmap, nreduce, file, master)
	mr.StartRegistrationServer()
	go mr.Run()
	return mr
}

// Register handles worker registration requests via RPC.
//
// When a worker starts up, it calls this method to register itself
// with the master. The master adds the worker to its available
// worker pool and can then assign tasks to it.
//
// Parameters:
//   - args: Registration arguments containing worker address
//   - res: Reply structure to return registration status
//
// Returns:
//   - error: Always nil (registration is always accepted)
func (mr *MapReduce) Register(args *RegisterArgs, res *RegisterReply) error {
	DPrintf("Register: worker %s\n", args.Worker)
	mr.registerChannel <- args.Worker
	res.OK = true
	return nil
}

// Shutdown gracefully shuts down the MapReduce master.
//
// This method stops accepting new worker registrations and
// closes the network listener, causing the registration
// server goroutine to exit.
//
// Parameters:
//   - args: Shutdown arguments (unused)
//   - res: Shutdown reply (unused)
//
// Returns:
//   - error: Always nil
func (mr *MapReduce) Shutdown(args *ShutdownArgs, res *ShutdownReply) error {
	DPrintf("Shutdown: registration server\n")
	mr.alive = false
	mr.l.Close() // causes the Accept to fail
	return nil
}

// StartRegistrationServer starts the RPC server for worker registration.
//
// This method sets up a Unix domain socket listener to accept
// worker registration requests. It runs the RPC server in a
// separate goroutine to handle concurrent worker connections.
//
// The server uses Unix domain sockets for local communication
// between master and worker processes on the same machine.
func (mr *MapReduce) StartRegistrationServer() {
	rpcs := rpc.NewServer()
	rpcs.Register(mr)
	
	// Remove any existing socket file (Unix domain sockets only)
	os.Remove(mr.MasterAddress)
	
	l, err := net.Listen("unix", mr.MasterAddress)
	if err != nil {
		log.Fatalf("RegistrationServer %s error: %v", mr.MasterAddress, err)
	}
	mr.l = l

	// Start accepting connections in a separate goroutine
	go func() {
		for mr.alive {
			conn, err := mr.l.Accept()
			if err == nil {
				// Handle each connection in its own goroutine
				go func() {
					rpcs.ServeConn(conn)
					conn.Close()
				}()
			} else {
				DPrintf("RegistrationServer: accept error %v", err)
				break
			}
		}
		DPrintf("RegistrationServer: done\n")
	}()
}

// MapName generates the filename for a specific map task input.
//
// This function creates standardized filenames for map task inputs
// by combining the base filename with the map job number.
//
// Parameters:
//   - fileName: Base name of the input file
//   - MapJob: Zero-based index of the map task
//
// Returns:
//   - string: Generated filename in format "mrtmp.<fileName>-<MapJob>"
//
// Example:
//   MapName("input.txt", 2) returns "mrtmp.input.txt-2"
func MapName(fileName string, MapJob int) string {
	return "mrtmp." + fileName + "-" + strconv.Itoa(MapJob)
}

// Split divides the input file into nMap approximately equal-sized chunks.
//
// This method reads the input file and splits it into multiple files,
// one for each map task. The splitting is done on line boundaries to
// ensure that complete lines are preserved in each chunk.
//
// The method creates files named "mrtmp.<fileName>-0", "mrtmp.<fileName>-1",
// etc., where each file contains a portion of the original input.
//
// Parameters:
//   - fileName: Name of the input file to split
//
// The split files are created in the current working directory.
func (mr *MapReduce) Split(fileName string) {
	fmt.Printf("Split %s\n", fileName)
	
	infile, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Split: failed to open input file %s: %v", fileName, err)
	}
	defer infile.Close()
	
	fi, err := infile.Stat()
	if err != nil {
		log.Fatalf("Split: failed to stat input file %s: %v", fileName, err)
	}
	
	size := fi.Size()
	nchunk := size/int64(mr.nMap) + 1 // Ensure we don't have empty chunks

	outfile, err := os.Create(MapName(fileName, 0))
	if err != nil {
		log.Fatalf("Split: failed to create output file: %v", err)
	}
	writer := bufio.NewWriter(outfile)
	currentMap := 1
	bytesWritten := 0

	scanner := bufio.NewScanner(infile)
	for scanner.Scan() {
		// Check if we need to start a new output file
		if int64(bytesWritten) > nchunk*int64(currentMap) {
			writer.Flush()
			outfile.Close()
			
			outfile, err = os.Create(MapName(fileName, currentMap))
			if err != nil {
				log.Fatalf("Split: failed to create output file %d: %v", currentMap, err)
			}
			writer = bufio.NewWriter(outfile)
			currentMap++
		}
		
		line := scanner.Text() + "\n"
		writer.WriteString(line)
		bytesWritten += len(line)
	}
	
	if err := scanner.Err(); err != nil {
		log.Fatalf("Split: error reading input file: %v", err)
	}
	
	writer.Flush()
	outfile.Close()
}

// ReduceName generates the filename for intermediate data from a map task.
//
// This function creates standardized filenames for intermediate files
// that contain key-value pairs destined for a specific reduce task.
//
// Parameters:
//   - fileName: Base name of the input file
//   - MapJob: Zero-based index of the map task that produced this data
//   - ReduceJob: Zero-based index of the reduce task that will process this data
//
// Returns:
//   - string: Generated filename in format "mrtmp.<fileName>-<MapJob>-<ReduceJob>"
//
// Example:
//   ReduceName("input.txt", 2, 1) returns "mrtmp.input.txt-2-1"
func ReduceName(fileName string, MapJob int, ReduceJob int) string {
	return MapName(fileName, MapJob) + "-" + strconv.Itoa(ReduceJob)
}

// hash computes a 32-bit hash value for a string using the FNV-1a algorithm.
//
// This function is used to distribute key-value pairs across reduce tasks.
// Keys with the same hash value (modulo nReduce) will be processed by
// the same reduce task, ensuring that all values for a given key are
// processed together.
//
// Parameters:
//   - s: The string to hash
//
// Returns:
//   - uint32: 32-bit hash value of the input string
func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

// DoMap executes a map task for a specific input chunk.
//
// This function reads the input chunk for the specified map job,
// applies the user-provided Map function to process the data,
// and distributes the resulting key-value pairs across nReduce
// intermediate files based on hash partitioning.
//
// Parameters:
//   - JobNumber: Zero-based index of the map task to execute
//   - fileName: Base name of the input file
//   - nreduce: Number of reduce tasks (determines partitioning)
//   - Map: User-provided map function that processes input and returns key-value pairs
//
// The function creates nReduce intermediate files, each containing
// key-value pairs destined for a specific reduce task.
func DoMap(JobNumber int, fileName string, nreduce int, Map func(string) *list.List) {
	name := MapName(fileName, JobNumber)
	
	file, err := os.Open(name)
	if err != nil {
		log.Fatalf("DoMap: failed to open input file %s: %v", name, err)
	}
	
	fi, err := file.Stat()
	if err != nil {
		log.Fatalf("DoMap: failed to stat input file %s: %v", name, err)
	}
	
	size := fi.Size()
	fmt.Printf("DoMap: read split %s %d\n", name, size)
	
	// Read the entire input chunk into memory
	b := make([]byte, size)
	_, err = file.Read(b)
	if err != nil {
		log.Fatalf("DoMap: failed to read input file %s: %v", name, err)
	}
	file.Close()
	
	// Apply the user-provided Map function
	res := Map(string(b))
	
	// Distribute key-value pairs across reduce tasks using hash partitioning
	// This ensures that all values for the same key go to the same reduce task
	for r := 0; r < nreduce; r++ {
		file, err = os.Create(ReduceName(fileName, JobNumber, r))
		if err != nil {
			log.Fatalf("DoMap: failed to create intermediate file %s: %v", 
				ReduceName(fileName, JobNumber, r), err)
		}
		
		enc := json.NewEncoder(file)
		for e := res.Front(); e != nil; e = e.Next() {
			kv := e.Value.(KeyValue)
			// Use hash partitioning to determine which reduce task gets this key
			if hash(kv.Key)%uint32(nreduce) == uint32(r) {
				err := enc.Encode(&kv)
				if err != nil {
					log.Fatalf("DoMap: failed to encode key-value pair: %v", err)
				}
			}
		}
		file.Close()
	}
}

// MergeName generates the filename for a reduce task output.
//
// This function creates standardized filenames for the output files
// produced by reduce tasks, which contain the final aggregated results.
//
// Parameters:
//   - fileName: Base name of the original input file
//   - ReduceJob: Zero-based index of the reduce task
//
// Returns:
//   - string: Generated filename in format "mrtmp.<fileName>-res-<ReduceJob>"
//
// Example:
//   MergeName("input.txt", 2) returns "mrtmp.input.txt-res-2"
func MergeName(fileName string, ReduceJob int) string {
	return "mrtmp." + fileName + "-res-" + strconv.Itoa(ReduceJob)
}

// DoReduce executes a reduce task for a specific partition.
//
// This function reads all intermediate files for the specified reduce job
// from all map tasks, groups key-value pairs by key, sorts the keys,
// and applies the user-provided Reduce function to aggregate values
// for each key.
//
// Parameters:
//   - job: Zero-based index of the reduce task to execute
//   - fileName: Base name of the input file
//   - nmap: Number of map tasks (determines how many intermediate files to read)
//   - Reduce: User-provided reduce function that aggregates values for a key
//
// The function creates a single output file containing the final
// aggregated results for all keys assigned to this reduce task.
func DoReduce(job int, fileName string, nmap int, Reduce func(string, *list.List) string) {
	// Map to group values by key across all map task outputs
	kvs := make(map[string]*list.List)
	
	// Read intermediate files from all map tasks for this reduce partition
	for i := 0; i < nmap; i++ {
		name := ReduceName(fileName, i, job)
		fmt.Printf("DoReduce: read %s\n", name)
		
		file, err := os.Open(name)
		if err != nil {
			log.Fatalf("DoReduce: failed to open intermediate file %s: %v", name, err)
		}
		
		dec := json.NewDecoder(file)
		for {
			var kv KeyValue
			err = dec.Decode(&kv)
			if err != nil {
				// EOF or other error - end of file
				break
			}
			
			// Group values by key
			if kvs[kv.Key] == nil {
				kvs[kv.Key] = list.New()
			}
			kvs[kv.Key].PushBack(kv.Value)
		}
		file.Close()
	}
	
	// Sort keys for deterministic output order
	var keys []string
	for k := range kvs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	
	// Create output file and write results
	p := MergeName(fileName, job)
	file, err := os.Create(p)
	if err != nil {
		log.Fatalf("DoReduce: failed to create output file %s: %v", p, err)
	}
	
	enc := json.NewEncoder(file)
	for _, k := range keys {
		// Apply user-provided Reduce function to aggregate values
		res := Reduce(k, kvs[k])
		err := enc.Encode(KeyValue{Key: k, Value: res})
		if err != nil {
			log.Fatalf("DoReduce: failed to encode result for key %s: %v", k, err)
		}
	}
	file.Close()
}

// Merge combines the results from all reduce tasks into a single output file.
//
// This method reads the output files from all reduce tasks, sorts the
// keys alphabetically, and writes the final results to a single file
// in a human-readable format.
//
// The final output file is named "mrtmp.<inputFileName>" and contains
// one line per key-value pair in the format "key: value".
func (mr *MapReduce) Merge() {
	DPrintf("Merge phase")
	
	// Collect all key-value pairs from reduce outputs
	kvs := make(map[string]string)
	for i := 0; i < mr.nReduce; i++ {
		p := MergeName(mr.file, i)
		fmt.Printf("Merge: read %s\n", p)
		
		file, err := os.Open(p)
		if err != nil {
			log.Fatalf("Merge: failed to open reduce output %s: %v", p, err)
		}
		
		dec := json.NewDecoder(file)
		for {
			var kv KeyValue
			err = dec.Decode(&kv)
			if err != nil {
				// EOF or other error - end of file
				break
			}
			kvs[kv.Key] = kv.Value
		}
		file.Close()
	}
	
	// Sort keys for deterministic output order
	var keys []string
	for k := range kvs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Write final results to output file
	outputFile := "mrtmp." + mr.file
	file, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("Merge: failed to create final output file %s: %v", outputFile, err)
	}
	
	w := bufio.NewWriter(file)
	for _, k := range keys {
		fmt.Fprintf(w, "%s: %s\n", k, kvs[k])
	}
	w.Flush()
	file.Close()
}

// RemoveFile removes a file and logs any errors.
//
// This helper function is used during cleanup to remove temporary
// files created during MapReduce execution.
//
// Parameters:
//   - n: Name of the file to remove
//
// The function will terminate the program if file removal fails.
func RemoveFile(n string) {
	err := os.Remove(n)
	if err != nil {
		log.Fatalf("CleanupFiles: failed to remove file %s: %v", n, err)
	}
}

// CleanupFiles removes all temporary files created during MapReduce execution.
//
// This method removes:
//   - Input split files (mrtmp.<file>-0, mrtmp.<file>-1, ...)
//   - Intermediate files (mrtmp.<file>-<map>-<reduce>)
//   - Reduce output files (mrtmp.<file>-res-0, mrtmp.<file>-res-1, ...)
//   - Final output file (mrtmp.<file>)
//
// This cleanup is typically called after the MapReduce job completes
// to free up disk space.
func (mr *MapReduce) CleanupFiles() {
	// Remove input split files
	for i := 0; i < mr.nMap; i++ {
		RemoveFile(MapName(mr.file, i))
		// Remove intermediate files from this map task
		for j := 0; j < mr.nReduce; j++ {
			RemoveFile(ReduceName(mr.file, i, j))
		}
	}
	
	// Remove reduce output files
	for i := 0; i < mr.nReduce; i++ {
		RemoveFile(MergeName(mr.file, i))
	}
	
	// Remove final output file
	RemoveFile("mrtmp." + mr.file)
}

// RunSingle executes a MapReduce job sequentially on a single process.
//
// This function provides a simple way to run MapReduce jobs without
// the complexity of distributed execution. It's useful for testing,
// debugging, and processing small datasets.
//
// The execution follows the standard MapReduce phases:
//   1. Split the input file into chunks
//   2. Execute all map tasks sequentially
//   3. Execute all reduce tasks sequentially
//   4. Merge the results into a single output file
//
// Parameters:
//   - nMap: Number of map tasks to create
//   - nReduce: Number of reduce tasks to create
//   - file: Input file name to process
//   - Map: User-provided map function
//   - Reduce: User-provided reduce function
func RunSingle(nMap int, nReduce int, file string,
	Map func(string) *list.List,
	Reduce func(string, *list.List) string) {
	
	mr := InitMapReduce(nMap, nReduce, file, "")
	
	// Phase 1: Split input file into chunks
	mr.Split(mr.file)
	
	// Phase 2: Execute all map tasks sequentially
	for i := 0; i < nMap; i++ {
		DoMap(i, mr.file, mr.nReduce, Map)
	}
	
	// Phase 3: Execute all reduce tasks sequentially
	for i := 0; i < mr.nReduce; i++ {
		DoReduce(i, mr.file, mr.nMap, Reduce)
	}
	
	// Phase 4: Merge results into final output
	mr.Merge()
}

// CleanupRegistration shuts down the master's registration server.
//
// This method sends a shutdown request to the master process to
// gracefully stop accepting new worker registrations and close
// the network listener.
//
// If the RPC call fails, it logs an error but does not terminate
// the program, as the master may have already shut down.
func (mr *MapReduce) CleanupRegistration() {
	args := &ShutdownArgs{}
	var reply ShutdownReply
	ok := call(mr.MasterAddress, "MapReduce.Shutdown", args, &reply)
	if !ok {
		fmt.Printf("Cleanup: RPC %s error\n", mr.MasterAddress)
	}
	DPrintf("CleanupRegistration: done\n")
}

// Run executes the complete MapReduce job in distributed mode.
//
// This method orchestrates the entire MapReduce workflow:
//   1. Split the input file into chunks
//   2. Run the master to coordinate distributed execution
//   3. Merge results from all reduce tasks
//   4. Clean up the registration server
//   5. Signal completion via DoneChannel
//
// The method assumes a shared file system is available for
// communication between master and worker processes.
func (mr *MapReduce) Run() {
	fmt.Printf("Run mapreduce job %s %s\n", mr.MasterAddress, mr.file)

	// Phase 1: Split input file into chunks for map tasks
	mr.Split(mr.file)
	
	// Phase 2: Run master to coordinate distributed execution
	mr.stats = mr.RunMaster()
	
	// Phase 3: Merge results from all reduce tasks
	mr.Merge()
	
	// Phase 4: Clean up registration server
	mr.CleanupRegistration()

	fmt.Printf("%s: MapReduce done\n", mr.MasterAddress)

	// Signal completion to waiting goroutines
	mr.DoneChannel <- true
}
