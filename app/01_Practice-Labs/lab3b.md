# Lab 3B: Paxos-based Key/Value Server

## Introduction

Welcome to Lab 3B! This lab builds upon your Paxos implementation from Lab 3A to create a truly fault-tolerant distributed key-value storage system. You'll learn how to use Paxos consensus to build a replicated key-value service that maintains consistency even in the presence of network failures and server crashes.

### Learning Objectives

By the end of this lab, you will understand:
- How to use Paxos consensus to build a distributed key-value store
- How to implement sequential consistency across multiple replicas
- How to handle duplicate client requests and ensure at-most-once semantics
- How to manage server catch-up and log synchronization
- The challenges of building fault-tolerant distributed systems

### The Problem with Centralized Systems

Your Lab 2 implementation had a critical weakness: it depended on a single view server. If the view server crashed, your entire key-value service would stop working. Lab 3A solved this by implementing Paxos consensus, but now we need to use that consensus to build an actual application.

### The KVPaxos Solution

KVPaxos eliminates the single point of failure by using Paxos consensus for every operation. Instead of a central coordinator, all servers participate in reaching agreement on the order of operations. This means:

- **No Single Point of Failure**: No single server can bring down the entire system
- **Sequential Consistency**: All replicas see operations in the same order
- **Automatic Recovery**: Servers that were temporarily unavailable can catch up
- **At-Most-Once Semantics**: Operations execute exactly once, even with retries

### System Architecture

```
    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
    │   Client 1  │    │   Client 2  │    │   Client 3  │
    └──────┬──────┘    └──────┬──────┘    └──────┬──────┘
           │                  │                  │
           │    Put/Get       │    Put/Get       │    Put/Get
           │    Requests      │    Requests      │    Requests
           │                  │                  │
    ┌──────▼──────────────────▼──────────────────▼─────┐
    │              KVPaxos Servers                     │
    │  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌──────┐ │
    │  │KVPaxos  │  │KVPaxos  │  │KVPaxos  │  │ ...  │ │
    │  │Server 0 │  │Server 1 │  │Server 2 │  │      │ │
    │  └────┬────┘  └────┬────┘  └────┬────┘  └──────┘ │
    │       │            │            │                │
    │       │            │            │                │
    │  ┌────▼────┐  ┌────▼────┐  ┌────▼────┐  ┌──────┐ │
    │  │ Paxos   │  │ Paxos   │  │ Paxos   │  │ ...  │ │
    │  │ Peer 0  │  │ Peer 1  │  │ Peer 2  │  │      │ │
    │  └─────────┘  └─────────┘  └─────────┘  └──────┘ │
    └──────────────────────────────────────────────────┘
```

### Key Concepts

**Sequential Consistency**: All replicas must see operations in the same order. If client A does `Put("x", "1")` and then client B does `Put("x", "2")`, all servers must see these operations in that exact order.

**At-Most-Once Semantics**: Each operation must execute exactly once, even if the client retries the request due to network failures.

**Operation Log**: Think of Paxos as maintaining a log of operations. All servers apply operations from this log in the same order, ensuring consistency.

**Catch-Up Mechanism**: Servers that fall behind must be able to catch up by learning about operations they missed.

## Part B: Implementing the KVPaxos Server

### Getting Started

The skeleton code for this lab is in the `src/kvpaxos` directory. Let's start by examining what we have:

```shell
$ cd src/kvpaxos
$ go test
TestBasic: --- FAIL: TestBasic (5.02 seconds)
        test_test.go:XX: Get("a") -> "", expected "aa"
...
$
```

**Expected Initial Results**: The tests will fail initially because the KVPaxos implementation is incomplete. This is normal! You'll see failures because the server logic isn't implemented yet.

### Your Goal

Implement the KVPaxos key-value service in `server.go`. When complete, you should pass all tests:

```shell
$ cd src/kvpaxos
$ go test
Test: Basic put/puthash/get ...
  ... Passed
Test: Concurrent clients ...
  ... Passed
Test: server frees Paxos log memory...
  ... Passed
Test: Partitioned servers and client ...
  ... Passed
Test: Unreliable ...
  ... Passed
Test: Tolerates holes in paxos sequence ...
  ... Passed
Test: Many clients, changing partitions ...
  ... Passed
PASS
ok      kvpaxos   45.123s
$
```

### Required Interface

Your KVPaxos implementation must support these RPC methods:

```go
func (kv *KVPaxos) Get(args *GetArgs, reply *GetReply) error
func (kv *KVPaxos) Put(args *PutArgs, reply *PutReply) error
```

### Key Requirements

**Sequential Consistency**: All completed operations must appear to have affected all replicas in the same order.

**At-Most-Once Semantics**: Each operation must execute exactly once, even with client retries.

**Majority Rule**: A server should only complete operations if it's part of a majority that can communicate.

**Catch-Up**: Servers that fall behind must be able to catch up by learning about missed operations.

## Implementation Strategy

### Step-by-Step Approach

Here's a recommended plan for implementing KVPaxos:

#### Step 1: Define the Operation Structure
Fill in the `Op` struct in `server.go` with the information needed for each operation:
- Client identifier and operation ID for duplicate detection
- Operation type (Get, Put, or PutHash)
- Key and value for the operation

#### Step 2: Implement the Put Handler
The `Put()` handler should:
- Create an `Op` struct for the operation
- Use Paxos to get agreement on a sequence number
- Execute the operation and return the result
- Handle duplicate detection

#### Step 3: Implement the Get Handler
The `Get()` handler should:
- Create an `Op` struct for the operation
- Use Paxos to get agreement on a sequence number
- Ensure the server is caught up with all previous operations
- Return the current value for the key

#### Step 4: Implement Duplicate Detection
Add mechanisms to ensure operations execute at most once:
- Track completed operation IDs
- Check operation logs for previously executed operations
- Return cached results for duplicate requests

#### Step 5: Implement Catch-Up Mechanism
Add logic for servers to catch up when they fall behind:
- Force agreement on missing sequence numbers
- Execute all operations up to the current sequence
- Maintain operation logs for duplicate detection

## Detailed Implementation Hints

### Operation Structure Design
```go
type Op struct {
    Client uint64 // Client identifier for duplicate detection
    OpId   uint64 // Unique operation ID for duplicate detection
    Put    int    // Operation type: 0 = Get, 1 = Put, 2 = PutHash
    Key    string // The key for the operation
    Value  string // The value for Put operations (empty for Gets)
}
```

### Sequence Number Management
- **Hint**: Your server should try to assign the next available Paxos instance (sequence number) to each incoming client RPC
- **Challenge**: Other replicas may also be trying to use that instance for different operations
- **Solution**: Be prepared to try different sequence numbers until you get agreement

### Duplicate Detection Strategy
- **Client Identification**: Use the client ID and operation ID to uniquely identify operations
- **Operation Logging**: Keep a log of completed operations per client
- **Cache Results**: Return cached results for operations that have already been executed

### Paxos Integration
- **No Direct Communication**: Servers should only interact through the Paxos log
- **Majority Requirement**: Only complete operations if you're part of a majority
- **Done() Calls**: Don't forget to call `px.Done(seq)` when you're finished with a sequence

### Waiting for Paxos Agreement
```go
to := 10 * time.Millisecond
for {
    decided, _ := kv.px.Status(seq)
    if decided {
        // Process the decided operation
        return
    }
    time.Sleep(to)
    if to < 10 * time.Second {
        to *= 2
    }
}
```

### Catch-Up Mechanism
- **Force Agreement**: If a server falls behind, use `px.Start(seq, Nop)` to force agreement on missing sequences
- **Execute in Order**: Process all operations from your last completed sequence to the current one
- **No-Op Operations**: Use no-op operations to fill gaps in the sequence

### Memory Management
- **Operation Logs**: Keep operation logs to prevent memory growth
- **Done() Tracking**: Track Done() values from all peers to determine what can be forgotten
- **Garbage Collection**: Periodically clean up old operation logs

### Error Handling
- **Gob Registration**: Make sure to register your `Op` struct with `gob.Register(Op{})`
- **RPC Failures**: Handle RPC failures gracefully
- **Network Partitions**: Ensure the system works correctly under network partitions

## Common Pitfalls to Avoid

### Concurrency Issues
- **Race Conditions**: Use mutexes to protect shared state
- **Deadlocks**: Be careful with mutex ordering
- **Infinite Loops**: Always check termination conditions

### Paxos Misuse
- **Direct Communication**: Don't have servers communicate directly
- **Missing Done()**: Always call `px.Done()` when finished with a sequence
- **Wrong Sequence Numbers**: Make sure you're using the correct sequence numbers

### Duplicate Detection
- **Operation ID Collisions**: Ensure operation IDs are truly unique
- **Cache Inconsistency**: Keep operation logs consistent across servers
- **Memory Leaks**: Don't let operation logs grow indefinitely

### Performance Issues
- **Busy Waiting**: Use exponential backoff when waiting for Paxos agreement
- **Unnecessary Operations**: Don't execute operations that have already been completed
- **Memory Growth**: Implement proper garbage collection

## Testing and Debugging

### Test Categories
- **Basic Operations**: Put, Get, and PutHash functionality
- **Concurrent Access**: Multiple clients accessing the service simultaneously
- **Fault Tolerance**: Server failures and network partitions
- **Memory Management**: Proper cleanup of old operations
- **Unreliable Networks**: Message loss and network failures

### Debugging Tips
- **Enable Debug Output**: Set `Debug = 1` in `server.go` to see detailed logs
- **Check Gob Errors**: Look for "gob: type not registered" errors in logs
- **Verify Sequence Numbers**: Make sure sequence numbers are being used correctly
- **Monitor Memory Usage**: Watch for memory leaks in operation logs

### Expected Code Size
- **Target**: Around 200 lines of code for the server implementation
- **Focus**: Clean, well-documented code over complex optimizations
- **Structure**: Separate concerns between Paxos integration and key-value logic

## Reference Implementation Design

### Architecture Overview
The implementation uses a clean separation between Paxos consensus and key-value operations:

**Operation Flow**:
1. Client sends RPC request to any server
2. Server creates `Op` struct for the operation
3. Server uses Paxos to get agreement on a sequence number
4. Server executes all operations up to that sequence
5. Server returns result to client

**State Management**:
- **Key-Value Store**: In-memory map of key-value pairs
- **Operation Log**: Per-client log of completed operations
- **Sequence Tracking**: Current sequence number and last completed sequence

### Key Design Principles

1. **Consensus First**: Always use Paxos to agree on operation order
2. **Catch-Up Mechanism**: Servers must be able to recover from gaps
3. **Duplicate Detection**: Operations must execute at most once
4. **Memory Efficiency**: Implement proper cleanup of old operations
5. **Fault Tolerance**: System must work despite server failures

### Implementation Checklist

- [ ] Define `Op` struct with all necessary fields
- [ ] Implement `Put()` handler with Paxos integration
- [ ] Implement `Get()` handler with catch-up mechanism
- [ ] Add duplicate detection using operation logs
- [ ] Implement sequence number management
- [ ] Add memory management and garbage collection
- [ ] Test with various failure scenarios
- [ ] Verify all tests pass

## Getting Started

1. **Examine the Skeleton Code**: Look at the existing `Op` struct and RPC handlers
2. **Understand the Paxos Interface**: Review how to use the Paxos library from Lab 3A
3. **Start with Basic Operations**: Implement simple Put and Get operations first
4. **Add Duplicate Detection**: Ensure operations execute at most once
5. **Implement Catch-Up**: Add logic for servers to recover from gaps
6. **Test Thoroughly**: Run all tests and debug any failures

Remember: The key insight is that every operation must go through Paxos consensus to ensure all replicas stay synchronized. The challenge is managing the complexity of sequence numbers, duplicate detection, and server catch-up while maintaining good performance. 
