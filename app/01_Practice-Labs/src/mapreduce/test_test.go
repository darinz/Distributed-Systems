// Package mapreduce provides comprehensive tests for the MapReduce framework.
//
// This test suite validates the correctness and fault tolerance of the
// MapReduce implementation through various scenarios including basic
// functionality, single worker failures, and multiple worker failures.
//
// The tests use a simple word-splitting MapReduce job to verify that:
//   - Basic MapReduce execution works correctly
//   - Worker failures are handled gracefully
//   - Jobs are redistributed when workers fail
//   - Final results are consistent regardless of failures
//
// Test Environment:
//   - Uses Unix domain sockets for local communication
//   - Creates temporary input files with sequential numbers
//   - Validates output correctness and worker participation
package mapreduce

import (
	"bufio"
	"container/list"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"
)

// Test configuration constants
const (
	nNumber = 100000 // Number of input numbers to generate for testing
	nMap    = 100    // Number of map tasks to create
	nReduce = 50     // Number of reduce tasks to create
)

// MapFunc is a test map function that splits input text into words.
//
// This function serves as a simple test case for the MapReduce framework.
// It takes a string input, splits it into words using whitespace as
// delimiters, and returns each word as a key-value pair with an empty value.
//
// Parameters:
//   - value: Input string to process
//
// Returns:
//   - *list.List: List of KeyValue pairs containing individual words
//
// Example:
//   Input: "hello world test"
//   Output: [{"hello", ""}, {"world", ""}, {"test", ""}]
func MapFunc(value string) *list.List {
	DPrintf("Map %v\n", value)
	
	res := list.New()
	words := strings.Fields(value)
	
	for _, w := range words {
		kv := KeyValue{Key: w, Value: ""}
		res.PushBack(kv)
	}
	
	return res
}

// ReduceFunc is a test reduce function that processes word counts.
//
// This function serves as a simple test case for the reduce phase.
// It receives a key (word) and a list of values (empty strings in this test),
// logs the processing for debugging, and returns an empty string.
//
// In a real word count application, this would count the number of
// occurrences of each word, but for testing purposes, it simply
// validates that the reduce phase is called correctly.
//
// Parameters:
//   - key: The word being processed
//   - values: List of values associated with the key
//
// Returns:
//   - string: Empty string (for testing purposes)
func ReduceFunc(key string, values *list.List) string {
	for e := values.Front(); e != nil; e = e.Next() {
		DPrintf("Reduce %s %v\n", key, e.Value)
	}
	return ""
}

// check validates that the MapReduce output matches the expected input.
//
// This function compares the input file against the MapReduce output file
// to ensure that all input numbers appear in the output in the correct
// sorted order. It verifies the correctness of the MapReduce execution
// by checking that no data is lost or corrupted during processing.
//
// Parameters:
//   - t: Testing context for reporting failures
//   - file: Base name of the input file (output file is "mrtmp." + file)
//
// The function will fail the test if:
//   - Input or output files cannot be opened
//   - Output contains incorrect or missing numbers
//   - Numbers are not in sorted order
//   - Output has wrong number of lines
func check(t *testing.T, file string) {
	// Open input file
	input, err := os.Open(file)
	if err != nil {
		log.Fatalf("check: failed to open input file %s: %v", file, err)
	}
	defer input.Close()
	
	// Open output file
	output, err := os.Open("mrtmp." + file)
	if err != nil {
		log.Fatalf("check: failed to open output file mrtmp.%s: %v", file, err)
	}
	defer output.Close()

	// Read and sort input lines
	var lines []string
	inputScanner := bufio.NewScanner(input)
	for inputScanner.Scan() {
		lines = append(lines, inputScanner.Text())
	}
	sort.Strings(lines)

	// Validate output against sorted input
	outputScanner := bufio.NewScanner(output)
	lineIndex := 0
	for outputScanner.Scan() {
		var inputValue int
		var outputValue int
		outputText := outputScanner.Text()
		
		// Parse input line
		n, err := fmt.Sscanf(lines[lineIndex], "%d", &inputValue)
		if n != 1 || err != nil {
			t.Fatalf("line %d: failed to parse input value '%s': %v", 
				lineIndex, lines[lineIndex], err)
		}
		
		// Parse output line
		n, err = fmt.Sscanf(outputText, "%d", &outputValue)
		if n != 1 || err != nil {
			t.Fatalf("line %d: failed to parse output value '%s': %v", 
				lineIndex, outputText, err)
		}
		
		// Compare values
		if inputValue != outputValue {
			t.Fatalf("line %d: input value %d != output value %d", 
				lineIndex, inputValue, outputValue)
		}
		
		lineIndex++
	}
	
	// Check that we processed the expected number of lines
	if lineIndex != nNumber {
		t.Fatalf("Expected %d lines in output, got %d", nNumber, lineIndex)
	}
}

// checkWorker validates that all workers participated in job execution.
//
// This function checks the statistics returned by workers during shutdown
// to ensure that each worker processed at least one job. This validates
// that the load balancing and job distribution mechanisms are working
// correctly.
//
// Parameters:
//   - t: Testing context for reporting failures
//   - l: List of integers representing jobs completed by each worker
//
// The function will fail the test if any worker reports 0 completed jobs,
// indicating that the worker was registered but never assigned work.
func checkWorker(t *testing.T, l *list.List) {
	workerCount := 0
	for e := l.Front(); e != nil; e = e.Next() {
		jobsCompleted := e.Value.(int)
		if jobsCompleted == 0 {
			t.Fatalf("Worker %d didn't complete any jobs", workerCount)
		}
		workerCount++
	}
	
	if workerCount == 0 {
		t.Fatalf("No workers reported job completion statistics")
	}
}

// makeInput creates a test input file with sequential numbers.
//
// This function generates a temporary input file containing nNumber
// sequential integers, one per line. This provides a predictable
// dataset for testing MapReduce correctness.
//
// Returns:
//   - string: Name of the created input file
//
// The generated file contains numbers from 0 to nNumber-1, which
// allows for easy validation of MapReduce output correctness.
func makeInput() string {
	name := "824-mrinput.txt"
	
	file, err := os.Create(name)
	if err != nil {
		log.Fatalf("makeInput: failed to create input file %s: %v", name, err)
	}
	
	w := bufio.NewWriter(file)
	for i := 0; i < nNumber; i++ {
		fmt.Fprintf(w, "%d\n", i)
	}
	
	w.Flush()
	file.Close()
	
	return name
}

// port generates a unique Unix domain socket path for testing.
//
// This function creates a unique socket path in /var/tmp to avoid
// conflicts between concurrent test runs. The path includes the
// user ID and process ID to ensure uniqueness.
//
// Parameters:
//   - suffix: Additional suffix to make the path unique (e.g., "master", "worker0")
//
// Returns:
//   - string: Unique Unix domain socket path
//
// Note: Uses /var/tmp instead of current directory because AFS
// (Andrew File System) doesn't support Unix domain sockets.
func port(suffix string) string {
	s := "/var/tmp/824-"
	s += strconv.Itoa(os.Getuid()) + "/"
	
	// Create directory if it doesn't exist
	os.Mkdir(s, 0777)
	
	s += "mr"
	s += strconv.Itoa(os.Getpid()) + "-"
	s += suffix
	
	return s
}

// setup initializes a MapReduce instance for testing.
//
// This function creates a test environment by generating an input file
// and setting up a MapReduce master with the specified configuration.
// It provides a clean setup for each test case.
//
// Returns:
//   - *MapReduce: Configured MapReduce instance ready for testing
func setup() *MapReduce {
	file := makeInput()
	master := port("master")
	mr := MakeMapReduce(nMap, nReduce, file, master)
	return mr
}

// cleanup removes all temporary files created during testing.
//
// This function ensures that test artifacts are properly cleaned up
// after each test, preventing disk space issues and avoiding
// interference between test runs.
//
// Parameters:
//   - mr: MapReduce instance to clean up
func cleanup(mr *MapReduce) {
	mr.CleanupFiles()
	RemoveFile(mr.file)
}

// TestBasic tests basic MapReduce functionality with no failures.
//
// This test verifies that the MapReduce framework can successfully
// process a dataset using multiple workers without any failures.
// It validates that:
//   - Workers can register with the master
//   - Jobs are distributed correctly
//   - All tasks complete successfully
//   - Output is correct and complete
//   - All workers participate in job execution
func TestBasic(t *testing.T) {
	fmt.Printf("Test: Basic mapreduce ...\n")
	
	// Set up test environment
	mr := setup()
	
	// Start 2 workers that will run indefinitely
	for i := 0; i < 2; i++ {
		go RunWorker(mr.MasterAddress, port("worker"+strconv.Itoa(i)),
			MapFunc, ReduceFunc, -1)
	}
	
	// Wait for MapReduce job to complete
	<-mr.DoneChannel
	
	// Validate results
	check(t, mr.file)
	checkWorker(t, mr.stats)
	
	// Clean up test artifacts
	cleanup(mr)
	fmt.Printf("  ... Basic Passed\n")
}

// TestOneFailure tests MapReduce behavior when one worker fails.
//
// This test verifies that the framework can handle worker failures
// gracefully by redistributing failed jobs to other workers.
// It starts one worker that fails after 10 jobs and one that
// runs indefinitely, ensuring that all jobs eventually complete.
//
// The test validates that:
//   - Failed jobs are reassigned to other workers
//   - The job completes successfully despite worker failure
//   - Output is correct and complete
//   - Remaining workers handle the additional load
func TestOneFailure(t *testing.T) {
	fmt.Printf("Test: One Failure mapreduce ...\n")
	
	// Set up test environment
	mr := setup()
	
	// Start one worker that fails after 10 jobs
	go RunWorker(mr.MasterAddress, port("worker"+strconv.Itoa(0)),
		MapFunc, ReduceFunc, 10)
	
	// Start one worker that runs indefinitely
	go RunWorker(mr.MasterAddress, port("worker"+strconv.Itoa(1)),
		MapFunc, ReduceFunc, -1)
	
	// Wait for MapReduce job to complete
	<-mr.DoneChannel
	
	// Validate results
	check(t, mr.file)
	checkWorker(t, mr.stats)
	
	// Clean up test artifacts
	cleanup(mr)
	fmt.Printf("  ... One Failure Passed\n")
}

// TestManyFailures tests MapReduce behavior under continuous worker failures.
//
// This test simulates a challenging scenario where workers continuously
// fail and new workers are started to replace them. It validates that
// the framework can maintain progress even under high failure rates.
//
// Test behavior:
//   - Starts 2 new workers every second
//   - Each worker fails after completing 10 jobs
//   - Continues until the MapReduce job completes
//   - Validates that progress is made despite continuous failures
//
// This test verifies that:
//   - The framework can handle high failure rates
//   - Jobs are continuously reassigned as workers fail
//   - New workers can join and contribute to progress
//   - The job eventually completes successfully
func TestManyFailures(t *testing.T) {
	fmt.Printf("Test: Many Failures mapreduce ...\n")
	
	// Set up test environment
	mr := setup()
	
	workerIndex := 0
	done := false
	
	for !done {
		select {
		case done = <-mr.DoneChannel:
			// Job completed successfully
			check(t, mr.file)
			cleanup(mr)
			break
			
		default:
			// Start 2 new workers that will fail after 10 jobs each
			workerAddr1 := port("worker" + strconv.Itoa(workerIndex))
			go RunWorker(mr.MasterAddress, workerAddr1, MapFunc, ReduceFunc, 10)
			workerIndex++
			
			workerAddr2 := port("worker" + strconv.Itoa(workerIndex))
			go RunWorker(mr.MasterAddress, workerAddr2, MapFunc, ReduceFunc, 10)
			workerIndex++
			
			// Wait 1 second before starting next batch of workers
			time.Sleep(1 * time.Second)
		}
	}

	fmt.Printf("  ... Many Failures Passed\n")
}
