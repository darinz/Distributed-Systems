# Lab 1: MapReduce

## Introduction

Welcome to Lab 1! In this lab, you'll build a distributed MapReduce system using Go. This lab is designed to teach you:

1. **Go Programming**: Learn Go syntax, goroutines, channels, and RPC
2. **Distributed Systems**: Understand how to coordinate multiple processes
3. **Fault Tolerance**: Handle worker failures gracefully
4. **MapReduce Pattern**: Implement the classic distributed computing paradigm

### What You'll Build

You'll create a MapReduce framework with three main components:
- **Part I**: A word counting program using the MapReduce pattern
- **Part II**: A master process that distributes jobs to worker processes
- **Part III**: Fault-tolerant job distribution that handles worker failures

### Learning Objectives

By the end of this lab, you should understand:
- How to use Go's concurrency primitives (goroutines and channels)
- How to implement RPC communication between processes
- How to design fault-tolerant distributed systems
- How the MapReduce programming model works in practice

The implementation follows the design principles from the original [MapReduce paper](http://research.google.com/archive/mapreduce-osdi04.pdf), adapted for educational purposes.

## Getting Started

### Prerequisites
- Go 1.23.x installed on your system
- Basic understanding of Go syntax (if you're new to Go, check out the [Go Tour](https://tour.golang.org/))

### Setup Instructions

1. **Navigate to the project directory:**
   ```bash
   $ cd src/main
   ```

2. **Try running the initial code:**
   ```bash
   $ go run wc.go master kjv12.txt sequential
   ```

3. **You'll see compilation errors:**
   ```
   # command-line-arguments
   ./wc.go:11: missing return at end of function
   ./wc.go:15: missing return at end of function
   ```

   **Don't worry!** These errors are expected because the `Map` and `Reduce` functions are incomplete - that's what you'll implement in Part I.

### Understanding the Input Data

The lab uses `kjv12.txt`, which contains the King James Version of the Bible. This is a large text file (about 4MB) that's perfect for testing MapReduce because:
- It contains many repeated words
- It's large enough to benefit from parallel processing
- The word frequency results are easy to verify

### Project Structure

```
src/
├── main/
│   ├── wc.go          # Your word count implementation (Part I)
│   └── kjv12.txt      # Input text file
└── mapreduce/
    ├── mapreduce.go   # Core MapReduce framework
    ├── master.go      # Master implementation (Parts II & III)
    ├── worker.go      # Worker implementation (provided)
    ├── common.go      # RPC message types (provided)
    └── test_test.go   # Test suite
```

## Part I: Word Count Implementation

### Objective
Implement the `Map` and `Reduce` functions in `wc.go` to count word frequencies in the input text.

### Understanding MapReduce

MapReduce works in two phases:
1. **Map Phase**: Each map task processes a chunk of input and produces key-value pairs
2. **Reduce Phase**: Each reduce task aggregates values for the same key

For word counting:
- **Map**: Takes a text chunk → produces (word, "1") pairs
- **Reduce**: Takes (word, ["1", "1", "1"]) → produces (word, "3")

### Step-by-Step Implementation

#### Step 1: Study the Framework
Before coding, understand how the framework works:

1. **Read the MapReduce paper** (Section 2): [MapReduce Paper](http://research.google.com/archive/mapreduce-osdi04.pdf)
2. **Study the code** in `mapreduce.go`, especially:
   - `RunSingle()` function
   - `DoMap()` and `DoReduce()` functions
   - How the framework calls your Map and Reduce functions

#### Step 2: Implement the Map Function

The `Map` function receives a string (text chunk) and should return a list of `KeyValue` pairs.

**Key Requirements:**
- Split text into words (sequences of letters only)
- Return each word as a key with value "1"
- Use the provided `KeyValue` struct

**Implementation Hints:**
```go
// Use this to split text into words
func(r rune) bool { return !unicode.IsLetter(r) }

// Example usage:
words := strings.FieldsFunc(value, separator)
```

**Optimization Hint:** You can do local counting in the Map function to reduce network traffic:
```go
// Count words locally, then emit (word, count) pairs
wordCounts := make(map[string]int)
for _, word := range words {
    wordCounts[word]++
}
```

#### Step 3: Implement the Reduce Function

The `Reduce` function receives a key (word) and a list of values (counts) and should return the total count as a string.

**Key Requirements:**
- Sum all the count values for the given word
- Return the total as a string
- Handle the case where values might be empty

**Implementation Hints:**
```go
// Convert string to int
count, err := strconv.Atoi(valueStr)

// Convert int back to string
return strconv.Itoa(totalCount)
```

### Testing Your Implementation

#### Method 1: Run the Program
```bash
$ go run wc.go master kjv12.txt sequential
```

You should see output like:
```
Split kjv12.txt
DoMap: read split mrtmp.kjv12.txt-0 966954
DoMap: read split mrtmp.kjv12.txt-1 966953
...
Merge: read mrtmp.kjv12.txt-res-0
Merge: read mrtmp.kjv12.txt-res-1
Merge: read mrtmp.kjv12.txt-res-2
```

#### Method 2: Verify Results
```bash
$ sort -n -k2 mrtmp.kjv12.txt | tail -10
```

Expected output:
```
unto: 8940
he: 9666
shall: 9760
in: 12334
that: 12577
And: 12846
to: 13384
of: 34434
and: 38850
the: 62075
```

#### Method 3: Use the Test Script
```bash
$ ./test-wc.sh
```

### Success Criteria
- ✅ Program compiles without errors
- ✅ Produces correct word frequency counts
- ✅ Most frequent words match expected results
- ✅ Code is concise (should be ~10 lines total)

### Cleanup
Remove temporary files when done:
```bash
$ rm mrtmp.*
```

### Common Issues and Solutions

- **Issue**: "missing return at end of function"
- **Solution**: Make sure both `Map` and `Reduce` functions have return statements

- **Issue**: Wrong word counts
- **Solution**: Check that you're splitting on non-letter characters and handling empty strings

- **Issue**: Program hangs
- **Solution**: Ensure `Reduce` function returns a string, not an integer

## Part II: Distributed MapReduce Master

### Objective
Implement a master process that coordinates multiple worker processes to execute MapReduce jobs in parallel.

### Understanding the Architecture

The distributed MapReduce system consists of:
- **Master**: Coordinates job distribution and tracks progress
- **Workers**: Execute map and reduce tasks assigned by the master
- **RPC Communication**: Master and workers communicate via Unix domain sockets

### Key Concepts

#### Worker Registration
- Workers register with the master when they start up
- The master maintains a list of available workers
- Workers are assigned jobs as they become available

#### Job Distribution
- Map jobs are distributed first, then reduce jobs
- Each job runs in its own goroutine for parallel execution
- The master waits for all jobs to complete before proceeding

#### Channel-Based Coordination
- `registerChannel`: Receives worker registration requests
- `readyChannel`: Tracks available workers for job assignment
- `doneChannel`: Signals when jobs complete

### Implementation Guide

#### Step 1: Study the Provided Code

Before implementing, understand the existing code:

1. **`common.go`**: RPC message types (`DoJobArgs`, `DoJobReply`, etc.)
2. **`worker.go`**: Worker implementation (already complete)
3. **`mapreduce.go`**: Core framework with `Register` RPC handler
4. **`test_test.go`**: Test cases to validate your implementation

#### Step 2: Implement Worker Registration Handler

Create a goroutine that continuously accepts worker registrations:

```go
func (mr *MapReduce) registerWorkers() {
    for {
        address := <-mr.registerChannel
        // Add worker to available pool
        mr.Workers[address] = &WorkerInfo{address: address}
        // Make worker available for jobs
        mr.readyChannel <- address
    }
}
```

#### Step 3: Implement Job Delegation

Create a function that assigns jobs to available workers:

```go
func (mr *MapReduce) delegateJob(jtype JobType, jno int) {
    for {
        // Wait for available worker
        worker := <-mr.readyChannel
        
        // Try to assign job
        if mr.assignJob(worker, jtype, jno) {
            // Job completed successfully
            mr.doneChannel <- true
            mr.readyChannel <- worker  // Return worker to pool
            return
        }
        // Job failed, try with another worker
    }
}
```

#### Step 4: Implement Job Assignment

Create a function that sends RPC calls to workers:

```go
func (mr *MapReduce) assignJob(worker string, jtype JobType, jno int) bool {
    var args DoJobArgs
    switch jtype {
    case Map:
        args = DoJobArgs{
            File:          mr.file,
            Operation:     Map,
            JobNumber:     jno,
            NumOtherPhase: mr.nReduce,
        }
    case Reduce:
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
```

#### Step 5: Implement the Main Coordination Function

The `RunMaster()` function orchestrates the entire process:

```go
func (mr *MapReduce) RunMaster() *list.List {
    // Start worker registration handler
    go mr.registerWorkers()
    
    // Execute all map jobs in parallel
    for i := 0; i < mr.nMap; i++ {
        go mr.delegateJob(Map, i)
    }
    
    // Wait for all map jobs to complete
    for i := 0; i < mr.nMap; i++ {
        <-mr.doneChannel
    }
    
    // Execute all reduce jobs in parallel
    for i := 0; i < mr.nReduce; i++ {
        go mr.delegateJob(Reduce, i)
    }
    
    // Wait for all reduce jobs to complete
    for i := 0; i < mr.nReduce; i++ {
        <-mr.doneChannel
    }
    
    // Shutdown workers and return statistics
    return mr.KillWorkers()
}
```

### Testing Your Implementation

#### Run the Tests
```bash
$ cd src/mapreduce
$ go test
```

#### Expected Output
You should see:
```
Test: Basic mapreduce ...
  ... Basic Passed
PASS
```

#### Debugging Tips
- Add `log.Printf()` statements to track execution flow
- Use `go test > out 2>&1` to capture all output
- Look for "PASS" at the end of test output

### Success Criteria
- ✅ All tests pass (`go test` shows PASS)
- ✅ Workers are properly registered and managed
- ✅ Jobs are distributed in parallel
- ✅ Master waits for all jobs to complete
- ✅ Workers are properly shut down

### Common Issues and Solutions

**Issue**: Tests hang or timeout
**Solution**: Ensure `RunMaster()` waits for all jobs to complete before returning

**Issue**: "Some worker didn't do any work"
**Solution**: Make sure workers are properly returned to the ready channel after job completion

**Issue**: RPC errors
**Solution**: Check that worker addresses are correct and workers are running

### Important Notes

- **Ignore RPC reflection errors**: The error messages about "wrong number of ins" are safe to ignore
- **Unix domain sockets**: Communication is local to the machine
- **No failure handling yet**: Part II assumes workers don't fail (that's Part III)
- **Design for Part III**: Structure your code to easily add failure handling later

## Part III: Fault-Tolerant MapReduce

### Objective
Enhance the master to handle worker failures gracefully by reassigning failed jobs to other workers.

### Understanding Fault Tolerance

#### Why Worker Failures Are Manageable
- **Stateless Workers**: Workers don't maintain persistent state between jobs
- **Idempotent Jobs**: Running the same job multiple times produces the same result
- **RPC Timeouts**: Failed RPC calls indicate worker problems

#### Failure Scenarios
1. **Worker Crashes**: Worker process terminates unexpectedly
2. **Network Issues**: Worker becomes unreachable but may still be running
3. **Worker Overload**: Worker becomes unresponsive due to resource constraints

#### Key Insight: Job Reassignment
When an RPC call to a worker fails, the master should:
1. Assume the worker is no longer available
2. Reassign the job to another available worker
3. Continue until the job completes successfully

### Implementation Strategy

#### The Good News: Minimal Changes Required!

If you implemented Part II correctly, you likely already have fault tolerance! Here's why:

```go
func (mr *MapReduce) delegateJob(jtype JobType, jno int) {
    for {
        worker := <-mr.readyChannel
        
        // This is the key: if assignJob fails, we loop and try again
        if mr.assignJob(worker, jtype, jno) {
            mr.doneChannel <- true
            mr.readyChannel <- worker
            return
        }
        // If assignJob returns false (RPC failed), we try with another worker
    }
}
```

#### How It Works
1. **RPC Failure Detection**: The `call()` function returns `false` when RPC fails
2. **Automatic Retry**: The `delegateJob` loop continues until a worker succeeds
3. **Worker Pool Management**: Failed workers are simply not returned to the ready channel

### Testing Fault Tolerance

#### Test Cases
The test suite includes two fault tolerance tests:

1. **`TestOneFailure`**: One worker fails after 10 jobs, another continues
2. **`TestManyFailures`**: Multiple workers fail continuously, new workers are added

#### Running the Tests
```bash
$ cd src/mapreduce
$ go test
```

Expected output:
```
Test: Basic mapreduce ...
  ... Basic Passed
Test: One Failure mapreduce ...
  ... One Failure Passed
Test: Many Failures mapreduce ...
  ... Many Failures Passed
PASS
```

### Understanding the Test Behavior

#### TestOneFailure
- Starts 2 workers
- Worker 0 fails after 10 jobs
- Worker 1 continues and handles remaining jobs
- Verifies that all jobs complete despite worker failure

#### TestManyFailures
- Continuously starts new workers (2 every second)
- Each worker fails after 10 jobs
- Tests that the system makes progress despite continuous failures
- Demonstrates that new workers can join and contribute

### Common Implementation Issues

#### Issue 1: Not Handling RPC Failures
**Problem**: Code assumes RPC calls always succeed
**Solution**: Check the return value of `call()` and retry on failure

#### Issue 2: Deadlock on Worker Failures
**Problem**: Failed workers are not returned to the ready channel
**Solution**: Only return workers to the ready channel after successful job completion

#### Issue 3: Infinite Loops
**Problem**: If all workers fail, the system might hang
**Solution**: The test framework ensures new workers are continuously added

### Success Criteria
- ✅ All three tests pass (`TestBasic`, `TestOneFailure`, `TestManyFailures`)
- ✅ System handles single worker failures gracefully
- ✅ System makes progress despite continuous worker failures
- ✅ No deadlocks or infinite loops

### Advanced Considerations

#### Idempotency
- Jobs can be safely executed multiple times
- No need to worry about duplicate job execution
- Output files are overwritten, not appended

#### Master Failure
- This lab assumes the master never fails
- Master fault tolerance is much more complex (covered in later labs)
- Would require state replication and consensus protocols

#### Performance Implications
- Failed jobs cause delays due to retries
- Multiple workers may execute the same job
- Overall system throughput may decrease under high failure rates

### Debugging Fault Tolerance

#### Add Logging
```go
func (mr *MapReduce) assignJob(worker string, jtype JobType, jno int) bool {
    // ... existing code ...
    
    success := call(worker, "Worker.DoJob", args, &reply)
    if !success {
        log.Printf("Job assignment failed for worker %s, job %d", worker, jno)
    }
    return success
}
```

#### Monitor Worker Behavior
- Watch for RPC timeout messages
- Track which workers are being used
- Verify that failed workers are not reused

### Final Notes

- **Code Size**: The complete solution (Parts II + III) should be around 60 lines
- **Design Philosophy**: Simple retry logic is often sufficient for fault tolerance
- **Real-World Applications**: This pattern is used in production distributed systems
- **Next Steps**: Later labs will cover more complex fault tolerance scenarios

## Implementation Design Summary

### Part I: Word Count Implementation

**Map Function Strategy:**
- Split input text into words using `strings.FieldsFunc` with a custom separator
- Perform local word counting to reduce network traffic
- Emit (word, count) pairs instead of (word, "1") pairs for efficiency

**Reduce Function Strategy:**
- Sum all count values for each word
- Return the total count as a string
- Handle empty value lists gracefully

**Optimization Benefits:**
- Local counting reduces data transfer between map and reduce phases
- Less noticeable in parallel execution but improves sequential performance

### Part II & III: Distributed Master Implementation

**Architecture Overview:**
The master uses a channel-based architecture with three main components:

1. **Worker Registration Handler** (Goroutine)
   - Continuously listens on `registerChannel`
   - Adds new workers to the worker pool
   - Places workers on `readyChannel` for job assignment

2. **Job Delegation** (One goroutine per job)
   - Each map/reduce job gets its own goroutine
   - Waits for available worker from `readyChannel`
   - Attempts job assignment via RPC
   - On success: signals completion and returns worker to pool
   - On failure: retries with next available worker

3. **Master Coordination** (Main thread)
   - Starts worker registration handler
   - Launches all map job goroutines
   - Waits for all map jobs to complete
   - Launches all reduce job goroutines
   - Waits for all reduce jobs to complete
   - Shuts down workers and returns statistics

**Fault Tolerance Design:**
- **Automatic Retry**: Failed RPC calls trigger job reassignment
- **Worker Pool Management**: Failed workers are not returned to the ready pool
- **Idempotent Jobs**: Duplicate job execution is safe and expected
- **No State Recovery**: Workers are stateless, so failures don't require state recovery

**Key Design Principles:**
- **Parallel Execution**: All jobs run concurrently using goroutines
- **Channel Synchronization**: Channels coordinate worker availability and job completion
- **Simple Retry Logic**: RPC failures automatically trigger job reassignment
- **Resource Management**: Workers are properly returned to the pool after successful jobs

**Scalability Considerations:**
- New workers can register at any time during execution
- The system can handle dynamic worker addition and removal
- Job distribution is load-balanced across available workers
- The design supports both single-machine and distributed deployments

This design provides a robust, fault-tolerant MapReduce implementation that can handle worker failures gracefully while maintaining high performance through parallel execution.
