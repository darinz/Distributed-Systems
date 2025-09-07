// Package main implements a word count program using the MapReduce framework.
// This program demonstrates the MapReduce pattern by counting word frequencies
// in a text file, splitting the work across multiple map and reduce tasks.
//
// Usage:
//   Sequential: go run wc.go master <input-file> sequential
//   Distributed: go run wc.go master <input-file> <master-address>
//   Worker:     go run wc.go worker <master-address> <worker-address>
package main

import (
	"container/list"
	"fmt"
	"log"
	"mapreduce"
	"os"
	"strconv"
	"strings"
	"unicode"
)

// Map processes a text chunk and returns word frequency counts.
//
// This function implements the Map phase of the MapReduce algorithm:
// 1. Splits the input text into words using unicode.IsLetter as the delimiter
// 2. Performs local word counting to optimize data transfer to Reduce phase
// 3. Returns a list of KeyValue pairs where keys are words and values are counts
//
// Parameters:
//   - value: The input text chunk to process
//
// Returns:
//   - *list.List: A list of mapreduce.KeyValue pairs containing word counts
func Map(value string) *list.List {
	// Define separator function to split on non-letter characters
	// This ensures we only count actual words, ignoring punctuation and numbers
	separator := func(r rune) bool {
		return !unicode.IsLetter(r)
	}
	
	// Split the input text into individual words
	words := strings.FieldsFunc(value, separator)

	// Perform local word counting to reduce network traffic
	// This optimization minimizes the amount of data sent to the Reduce phase
	wordCounts := make(map[string]int)
	for _, word := range words {
		// Skip empty strings that might result from splitting
		if word == "" {
			continue
		}
		wordCounts[word]++
	}

	// Convert the map to a list of KeyValue pairs for the MapReduce framework
	result := list.New()
	for word, count := range wordCounts {
		result.PushBack(mapreduce.KeyValue{
			Key:   word,
			Value: strconv.Itoa(count),
		})
	}

	return result
}

// Reduce aggregates word counts for a specific word across all map outputs.
//
// This function implements the Reduce phase of the MapReduce algorithm:
// 1. Takes a word (key) and a list of count values from different map tasks
// 2. Sums all the counts to get the total frequency for that word
// 3. Returns the final count as a string
//
// Parameters:
//   - key: The word whose counts are being aggregated
//   - values: A list of count strings from different map tasks
//
// Returns:
//   - string: The total count for the word as a string
//
// Panics:
//   - If any value in the list cannot be converted to an integer
func Reduce(key string, values *list.List) string {
	// Initialize sum to accumulate all counts for this word
	sum := 0
	
	// Iterate through all count values for this word
	for elem := values.Front(); elem != nil; elem = elem.Next() {
		// Type assert the value to string and convert to integer
		countStr, ok := elem.Value.(string)
		if !ok {
			log.Fatalf("Reduce: expected string value, got %T", elem.Value)
		}
		
		count, err := strconv.Atoi(countStr)
		if err != nil {
			log.Fatalf("Reduce: failed to convert count '%s' to integer: %v", countStr, err)
		}
		
		sum += count
	}

	// Return the total count as a string for the MapReduce framework
	return strconv.Itoa(sum)
}

// main is the entry point for the word count MapReduce program.
//
// The program can be run in three different modes:
//
// 1. Sequential Mode:
//    go run wc.go master <input-file> sequential
//    Runs the MapReduce job sequentially on a single process
//
// 2. Master Mode:
//    go run wc.go master <input-file> <master-address>
//    Starts the MapReduce master that coordinates distributed workers
//
// 3. Worker Mode:
//    go run wc.go worker <master-address> <worker-address>
//    Starts a worker process that executes map and reduce tasks
//
// Configuration:
//   - Number of map tasks: 5
//   - Number of reduce tasks: 3
//   - Worker timeout: 100 seconds
func main() {
	// Validate command line arguments
	if len(os.Args) != 4 {
		fmt.Fprintf(os.Stderr, "Usage: %s <mode> <input-file> <address|sequential>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Modes:\n")
		fmt.Fprintf(os.Stderr, "  master <file> sequential     - Run sequentially\n")
		fmt.Fprintf(os.Stderr, "  master <file> <address>      - Run as master\n")
		fmt.Fprintf(os.Stderr, "  worker <master> <worker>     - Run as worker\n")
		os.Exit(1)
	}

	mode := os.Args[1]
	inputFile := os.Args[2]
	address := os.Args[3]

	switch mode {
	case "master":
		if address == "sequential" {
			// Run MapReduce sequentially on a single process
			// This is useful for testing and small datasets
			mapreduce.RunSingle(5, 3, inputFile, Map, Reduce)
		} else {
			// Start the MapReduce master process
			// The master coordinates workers and manages job distribution
			mr := mapreduce.MakeMapReduce(5, 3, inputFile, address)
			
			// Block until the MapReduce job completes
			// The DoneChannel will be closed when all tasks are finished
			<-mr.DoneChannel
		}
		
	case "worker":
		// Start a worker process that executes map and reduce tasks
		// Workers connect to the master and request work
		masterAddr := os.Args[2]
		workerAddr := os.Args[3]
		mapreduce.RunWorker(masterAddr, workerAddr, Map, Reduce, 100)
		
	default:
		fmt.Fprintf(os.Stderr, "Invalid mode: %s. Must be 'master' or 'worker'\n", mode)
		os.Exit(1)
	}
}
