# Lab 3A: Paxos-based Key/Value Service

## Introduction

Welcome to Lab 3A! This lab builds upon your primary/backup key-value service from Lab 2 to create a truly fault-tolerant distributed system using the Paxos consensus algorithm. You'll learn how to implement one of the most fundamental algorithms in distributed systems.

### Learning Objectives

By the end of this lab, you will understand:
- How to implement the Paxos consensus algorithm
- How to use Paxos to create a fault-tolerant key-value service
- The challenges of distributed consensus and how Paxos solves them
- Memory management in long-running distributed systems
- Network partition handling and majority-based decision making

### The Problem with Lab 2

Your Lab 2 implementation has a critical weakness: it depends on a single master view server. If the view server crashes or becomes unreachable, your entire key-value service stops working, even if the primary and backup servers are perfectly healthy. This creates a single point of failure.

Additionally, Lab 2's approach to handling temporary server unavailability is inefficient - it either blocks operations or declares servers dead and requires expensive complete database transfers.

### The Paxos Solution

Paxos eliminates the need for a central coordinator by using a distributed consensus algorithm. Instead of a single view server making decisions, all servers participate in reaching agreement on every operation. This means:

- **No Single Point of Failure**: No single server can bring down the entire system
- **Network Partition Tolerance**: The system continues working as long as a majority of servers can communicate
- **Automatic Recovery**: Servers that were temporarily unavailable can catch up by learning about operations they missed

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
    │              Paxos Consensus Layer               │
    │  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌──────┐ │
    │  │ Paxos   │  │ Paxos   │  │ Paxos   │  │ ...  │ │
    │  │ Peer 0  │  │ Peer 1  │  │ Peer 2  │  │      │ │
    │  └─────────┘  └─────────┘  └─────────┘  └──────┘ │
    └──────────────────────────────────────────────────┘
           │                  │                  │
           │   Agreed Order   │   Agreed Order   │   Agreed Order
           │                  │                  │
    ┌──────▼──────────────────▼──────────────────▼─────┐
    │              Key-Value Replicas                  │
    │  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌──────┐ │
    │  │   KV    │  │   KV    │  │   KV    │  │ ...  │ │
    │  │ Server 0│  │ Server 1│  │ Server 2│  │      │ │
    │  └─────────┘  └─────────┘  └─────────┘  └──────┘ │
    └──────────────────────────────────────────────────┘
```

### How Paxos Works

Paxos operates on a sequence of "instances" (numbered 0, 1, 2, ...). Each instance represents one operation that needs to be agreed upon. The algorithm ensures that:

1. **Safety**: All servers that decide on an instance agree on the same value
2. **Liveness**: If a majority of servers are available, some operation will eventually be decided

Your key-value service will use Paxos to agree on the order of operations. Each Put() or Get() request becomes a Paxos instance, and all servers apply these operations in the same order to maintain consistency.

### System Components

- **Clients**: Send Put(), PutHash(), and Get() RPCs to any available server
- **KVPaxos Servers**: Each contains a key-value database replica and a Paxos peer
- **Paxos Peers**: Implement the consensus algorithm and communicate via RPC

### Key Concepts

**Paxos Instances**: Each operation (Put/Get) becomes a numbered instance in the Paxos log. Instance 0 might be "Put(key1, value1)", instance 1 might be "Get(key2)", etc.

**Majority Rule**: Paxos requires agreement from a majority of servers. With 5 servers, you need at least 3 to agree. This ensures that even if some servers fail or become unreachable, the system can still make progress.

**Operation Log**: Think of Paxos as maintaining a log of operations. All servers apply operations from this log in the same order, ensuring consistency.

### Limitations of This Implementation

This lab implementation has some limitations that would need to be addressed in a production system:
- **No Persistence**: Data is lost if all servers crash simultaneously
- **Fixed Server Set**: Cannot add or remove servers dynamically  
- **Performance**: Each operation requires multiple Paxos messages
- **Memory Usage**: Must keep track of all operations until they can be forgotten

These limitations are acceptable for learning purposes and can be addressed in more advanced implementations.

### References

For deeper understanding, consult:
- Paxos notes and readings
- **Chubby**: Google's lock service using Paxos
- **Paxos Made Live**: Google's experience implementing Paxos
- **Spanner**: Google's globally-distributed database
- **Zookeeper**: Apache's coordination service
- **Harp** and **Viewstamped Replication**: Alternative consensus algorithms

## Part A: Implementing the Paxos Library

### Getting Started

The skeleton code for this lab is in the `src/paxos` directory. Let's start by examining what we have:

```shell
$ cd src/paxos
$ go test
Single proposer: --- FAIL: TestBasic (5.02 seconds)
        test_test.go:48: too few decided; seq=0 ndecided=0 wanted=3
Forgetting: --- FAIL: TestForget (5.03 seconds) 
        test_test.go:48: too few decided; seq=0 ndecided=0 wanted=6
...
$
```

**Expected Initial Results**: The tests will fail initially because the Paxos implementation is incomplete. This is normal! You'll see failures like "too few decided" because the Paxos algorithm isn't implemented yet.

**Ignore These Errors**: You'll also see many "has wrong number of ins" and "type Paxos has no exported methods" errors. These are Go compilation warnings that you can safely ignore - they don't affect the functionality.

### Your Goal

Implement the Paxos consensus algorithm in `paxos.go`. When complete, you should pass all tests in the paxos directory:

```shell
$ cd src/paxos
$ go test
Test: Single proposer ...
  ... Passed
Test: Many proposers, same value ...
  ... Passed
Test: Many proposers, different values ...
  ... Passed
Test: Out-of-order instances ...
  ... Passed
Test: Deaf proposer ...
  ... Passed
Test: Forgetting ...
  ... Passed
Test: Lots of forgetting ...
  ... Passed
Test: Paxos frees forgotten instance memory ...
  ... Passed
Test: RPC counts aren't too high ...
  ... Passed
Test: Many instances ...
  ... Passed
Test: Minority proposal ignored ...
  ... Passed
Test: Many instances, unreliable RPC ...
  ... Passed
Test: No decision if partitioned ...
  ... Passed
Test: Decision in majority partition ...
  ... Passed
Test: All agree after full heal ...
  ... Passed
Test: One peer switches partitions ...
  ... Passed
Test: One peer switches partitions, unreliable ...
  ... Passed
Test: Many requests, changing partitions ...
  ... Passed
PASS
ok      paxos   59.523s
$
```

**Note**: The tests take about a minute to complete due to the extensive fault tolerance testing.

### Required Interface

Your Paxos implementation must support this interface:

```go
px = paxos.Make(peers []string, me int, rpcs *rpc.Server, saveToDisk bool, dir string, restart bool)
px.Start(seq int, v interface{}) // start agreement on new instance
px.Status(seq int) (decided bool, v interface{}) // get info about an instance
px.Done(seq int) // ok to forget all instances <= seq
px.Max() int // highest instance seq known, or -1
px.Min() int // instances before this have been forgotten
```

### Interface Explanation

- **`Make(peers, me, rpcs, saveToDisk, dir, restart)`**: Creates a new Paxos peer
  - `peers`: List of all peer addresses (including this one)
  - `me`: Index of this peer in the peers array
  - `rpcs`: RPC server to register with (nil to create new one)
  - `saveToDisk`: Whether to persist state to disk
  - `dir`: Directory for persistent storage
  - `restart`: Whether this is a restart from crash

- **`Start(seq, v)`**: Begins agreement on instance `seq` with value `v`
  - Returns immediately without waiting for agreement
  - Use `Status()` to check if agreement is reached

- **`Status(seq)`**: Checks if instance `seq` has been decided
  - Returns `(true, value)` if decided, `(false, nil)` if not
  - Only checks local state, doesn't communicate with other peers

- **`Done(seq)`**: Indicates this peer is done with instances ≤ `seq`
  - Allows Paxos to forget old instances and free memory
  - See `Min()` explanation for details

- **`Max()`**: Returns highest sequence number this peer has seen
  - Returns -1 if no instances have been seen

- **`Min()`**: Returns smallest sequence number that hasn't been forgotten
  - All instances < `Min()` have been garbage collected

### Key Requirements

**Concurrent Instances**: Your implementation must handle multiple instances concurrently. If peers call `Start()` with different sequence numbers simultaneously, your implementation should run the Paxos protocol for all of them in parallel. Don't wait for instance i to complete before starting instance i+1.

**Memory Management**: Long-running Paxos servers must forget old instances to prevent memory leaks. The forgetting mechanism works as follows:

1. When a peer application is done with instances ≤ x, it calls `Done(x)`
2. Each peer tracks the highest `Done()` value from every other peer
3. A peer can forget instances ≤ min(all Done values from all peers)
4. `Min()` returns this minimum + 1

**Piggybacking**: It's acceptable to piggyback Done values in Paxos protocol messages rather than sending separate messages.

**Error Handling**: 
- If `Start()` is called with seq < `Min()`, ignore the call
- If `Status()` is called with seq < `Min()`, return `(false, nil)`
## The Paxos Algorithm

### Three-Phase Protocol

Paxos uses a three-phase protocol to reach consensus:

1. **Prepare Phase**: Proposer asks acceptors to promise not to accept proposals with lower numbers
2. **Accept Phase**: Proposer asks acceptors to accept a specific value
3. **Decide Phase**: Proposer notifies all peers that a value has been decided

### Pseudo-code

Here's the Paxos pseudo-code for a single instance:

```go
// PROPOSER SIDE
proposer(v):
  while not decided:
    choose n, unique and higher than any n seen so far
    send prepare(n) to all servers including self
    if prepare_ok(n_a, v_a) from majority:
      v' = v_a with highest n_a; choose own v otherwise
      send accept(n, v') to all
      if accept_ok(n) from majority:
        send decided(v') to all

// ACCEPTOR STATE
acceptor's state:
  n_p (highest prepare seen)
  n_a, v_a (highest accept seen)

// ACCEPTOR HANDLERS
acceptor's prepare(n) handler:
  if n > n_p
    n_p = n
    reply prepare_ok(n_a, v_a)
  else
    reply prepare_reject

acceptor's accept(n, v) handler:
  if n >= n_p
    n_p = n
    n_a = n
    v_a = v
    reply accept_ok(n)
  else
    reply accept_reject
```

### Understanding the Algorithm

**Proposal Numbers**: Each proposal has a unique number `n`. Proposers must choose numbers higher than any they've seen before.

**Majority Rule**: Each phase requires responses from a majority of acceptors. This ensures that even if some acceptors fail, the system can still make progress.

**Value Selection**: In the prepare phase, if any acceptor has already accepted a value, the proposer must use that value (the one with the highest proposal number). This ensures consistency.

**Safety**: The algorithm guarantees that if any value is decided, all peers will eventually agree on that same value.

## Implementation Strategy

### Step-by-Step Approach

Here's a recommended plan for implementing Paxos:

#### Step 1: Define Data Structures
Add elements to the `Paxos` struct in `paxos.go` to hold the state you'll need:
- Define a struct to hold information about each agreement instance
- Include fields for tracking proposal numbers, accepted values, and decision status
- Add maps to track instances and Done values from other peers

#### Step 2: Define RPC Types
Create RPC argument/reply types for Paxos protocol messages:
- `PrepareArgs`/`PrepareReply` for the prepare phase
- `AcceptArgs`/`AcceptReply` for the accept phase  
- `DecideArgs`/`DecideReply` for the decide phase
- Remember: field names in RPC structures must start with capital letters
- Include sequence numbers to identify which instance each message refers to

#### Step 3: Implement the Proposer
Write a proposer function that drives the Paxos protocol for an instance:
- Start a proposer function in its own goroutine for each instance (in `Start()`)
- Implement the three-phase protocol: prepare, accept, decide
- Handle retries when proposals are rejected
- Choose unique proposal numbers to avoid conflicts

#### Step 4: Implement Acceptors
Write RPC handlers that implement the acceptor side:
- `Prepare()` handler: respond to prepare requests
- `Accept()` handler: respond to accept requests
- `Decided()` handler: handle decide notifications
- All handlers should be thread-safe (use mutexes)

#### Step 5: Test Basic Functionality
At this point you should be able to pass the first few tests:
- Single proposer
- Multiple proposers with same value
- Multiple proposers with different values

#### Step 6: Implement Forgetting
Add the memory management mechanism:
- Track Done values from all peers
- Implement garbage collection to free old instances
- Update `Min()` to return the correct value

## Detailed Implementation Hints

### Concurrency and Ordering
- **Out-of-order execution**: Multiple Paxos instances may execute simultaneously and may be decided out of order (e.g., instance 10 might be decided before instance 5)
- **Concurrent instances**: Each instance should run independently in its own goroutine
- **Thread safety**: Use mutexes to protect shared state, especially in acceptor handlers

### Network Reliability
- **Local calls**: For unreliable network tests, call the local acceptor through a function call rather than RPC to avoid network failures
- **RPC vs LPC**: Implement a `send()` function that chooses between RPC (for remote peers) and local procedure calls (for self)

### Proposal Number Management
- **Unique numbers**: Use `px.me` (peer index) to help ensure proposal numbers are unique across peers
- **Higher numbers**: When retrying, choose proposal numbers higher than any seen so far
- **Liveness**: Paxos doesn't guarantee liveness - proposers might keep proposing higher numbers. Choose retry strategies carefully to improve chances of success

### Instance Management
- **Multiple proposals**: Multiple peers may call `Start()` on the same instance with different values
- **Already decided**: Applications may call `Start()` for instances that are already decided
- **Data structures**: Use maps (not slices) to store instance information for easy deletion during garbage collection

### Memory Management
- **Forgetting strategy**: Plan your forgetting mechanism before implementing - you need to track Done values from all peers
- **Garbage collection**: Implement a background goroutine to periodically clean up forgotten instances
- **Min() calculation**: Find the minimum Done value across all peers to determine what can be forgotten

### Performance Optimization
- **Message efficiency**: Use the minimum number of messages for agreement in non-failure cases
- **Piggybacking**: Include Done values in existing protocol messages rather than sending separate messages

### Testing and Debugging
- **Graceful shutdown**: Check `px.dead` in long-running loops and goroutines
- **Error handling**: Handle RPC failures gracefully
- **Code size**: Expect around 300 lines of code for a complete implementation

### Common Pitfalls to Avoid
- **Deadlocks**: Be careful with mutex ordering to avoid deadlocks
- **Race conditions**: Ensure all shared state access is properly synchronized
- **Infinite loops**: Always check termination conditions in proposer loops
- **Memory leaks**: Properly implement the forgetting mechanism to prevent memory growth


## Reference Implementation Design

### Architecture Overview

This implementation uses a clean separation between proposer and acceptor roles:

**Send Interface**: A unified `send()` function that either makes an RPC call (for remote peers) or a local procedure call (for self). This abstraction allows the proposer to treat all peers uniformly.

**Data Structures**: 
- Use a map (instance number → instance) to track Paxos instances for easy garbage collection
- Use a map (peer index → highest Done sequence) to track forgetting information

### Proposer Implementation

**Threading Model**: `Start()` launches a new goroutine running `px.propose()` and returns immediately. Each instance runs independently.

**Communication Strategy**: Send messages to peers sequentially (one at a time) rather than concurrently. This is slower but simpler and less prone to concurrency bugs.

**Done Piggybacking**: Include Done values in Decide() messages to efficiently distribute forgetting information.

### Acceptor Implementation

**RPC Handlers**: Three RPC methods handle the acceptor side:
- `Prepare()`: Responds to prepare requests
- `Accept()`: Responds to accept requests  
- `Decided()`: Handles decide notifications

**Isolation**: Even when a proposer calls its own acceptor, use the same `send()` interface to maintain isolation between proposer and acceptor roles.

**Thread Safety**: All acceptor methods are protected by a global mutex to ensure thread safety.

### State Management

**Read-Only Operations**: `Status()`, `Min()`, and `Max()` only read local state without modification.

**Done Updates**: `Done()` updates the local entry in the distributed Done table and returns immediately.

### Key Design Principles

1. **Separation of Concerns**: Keep proposer and acceptor logic separate
2. **Uniform Interface**: Use the same send interface for all peer communication
3. **Memory Efficiency**: Use maps for easy garbage collection of old instances
4. **Thread Safety**: Protect all shared state with appropriate synchronization
5. **Graceful Degradation**: Handle network failures and peer unavailability
