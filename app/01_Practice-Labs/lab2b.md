# Lab 2B: Primary/Backup Key/Value Service

## Introduction

Welcome to Lab 2B! This lab builds upon the View Service from Lab 2A to implement a fault-tolerant primary/backup key-value service. You'll learn how to build a distributed system that maintains consistency and availability even when servers fail.

### Learning Objectives

By the end of this lab, you will understand:
- How to implement primary/backup replication for fault tolerance
- At-most-once semantics and duplicate request handling
- State synchronization between primary and backup servers
- Client-side retry logic and primary discovery
- Network partition handling and split-brain prevention

### Architecture Overview

```
    ┌─────────────┐    ┌─────────────┐
    │   Client 1  │    │   Client 2  │
    └──────┬──────┘    └──────┬──────┘
           │                  │
           │    Get/Put       │
           │    Requests      │
           │                  │
    ┌──────▼──────────────────▼──────┐
    │         Primary Server         │ ◄── Handles all client requests
    │     (Key-Value Store)          │
    └─────────────┬──────────────────┘
                  │
                  │ Forward Operations
                  │
    ┌─────────────▼──────────────────┐
    │        Backup Server           │ ◄── Replicates primary state
    │     (Key-Value Store)          │
    └─────────────┬──────────────────┘
                  │
                  │ Ping for View Updates
                  │
    ┌─────────────▼──────────────────┐
    │      View Service              │ ◄── Manages primary/backup assignments
    └────────────────────────────────┘
```

## Part B: The Primary/Backup Key/Value Service

### System Requirements

Your key/value service must meet these critical requirements:

**Fault Tolerance**: The service should continue operating correctly as long as there has never been a time when no server was alive. It must handle:
- Server crashes and restarts
- Network partitions (servers that can't communicate with each other)
- Temporary network failures without server crashes

**Correctness**: 
- `Clerk.Get(k)` must return the latest value set by a successful `Clerk.Put(k,v)` or `Clerk.PutHash(k,v)`
- Return empty string if the key has never been stored
- All operations must provide **at-most-once semantics** (no duplicate execution)

**Availability**: If your service is operating with just one server, it should be able to incorporate a recovered or idle server as backup, enabling it to tolerate another server failure.

### Critical Design Constraints

**Single Active Primary**: Only one primary server can be active at any given time. This prevents split-brain scenarios where multiple servers think they're primary and serve different data to different clients.

**RPC-Only Communication**: All communication between clients and servers must use RPC with the `call()` function in `client.go`.

**View Service Dependency**: Assume the view service never halts or crashes.

**Non-Primary Behavior**: Servers that aren't the active primary should either not respond to clients or respond with an error (set `GetReply.Err` or `PutReply.Err` to something other than `OK`).

### Client Behavior Requirements

**Persistent Operations**: `Clerk.Get()`, `Clerk.Put()`, and `Clerk.PutHash()` must only return when they have completed the operation:
- Put operations should keep trying until the key/value database is updated
- Get operations should keep trying until they retrieve the current value (or confirm the key doesn't exist)

**Duplicate Handling**: Your server must filter out duplicate RPCs from client retries to ensure at-most-once semantics. You can assume each clerk has only one outstanding Put or Get operation.

**Commit Point**: Think carefully about what constitutes the commit point for a Put operation.

### Server Behavior Requirements

**Periodic View Updates**: Servers should not contact the view service for every Put/Get request (this would put the view service on the critical path). Instead, servers should ping the view service periodically in `pbservice/server.go`'s `tick()` function.

**View Transition Strategy**: Your one-primary-at-a-time strategy should rely on the view service only promoting the backup from view i to be primary in view i+1. This ensures:
- If the old primary tries to handle a client request, it forwards to its backup
- If the backup hasn't heard about view i+1, it's not acting as primary yet (no harm)
- If the backup has heard about view i+1 and is acting as primary, it rejects the old primary's forwarded requests

**State Synchronization**: Ensure the backup sees every update to the key/value database through:
- Primary initializing the backup with the complete key/value database
- Forwarding subsequent client Put operations to the backup

### Project Setup

The skeleton code for the key/value servers is in `src/pbservice`. It uses your view service from Lab 2A.

### Quick Start

```shell
# Navigate to the project directory
$ cd /path/to/your/lab/src/pbservice

# Run the tests to see the current state
$ go test -v
```

**Expected Initial Test Results**: The tests will initially fail because the implementation is incomplete. This is expected! You'll see failures like:

```shell
Single primary, no backup: --- FAIL: TestBasicFail (2.00 seconds)
        test_test.go:50: first primary never formed view
--- FAIL: TestFailPut (5.55 seconds)
        test_test.go:165: wrong primary or backup
Concurrent Put()s to the same key: --- FAIL: TestConcurrentSame (8.51 seconds)
    ...
Partition an old primary: --- FAIL: TestPartition (3.52 seconds)
        test_test.go:354: wrong primary or backup
    ...
```

Don't worry - these failures will guide your implementation!

## Implementation Roadmap

Here's a step-by-step plan to implement the primary/backup service:

### Step 1: Server View Discovery
**Goal**: Make servers aware of their role (primary/backup/neither)

**Implementation**:
- Modify `pbservice/server.go` to ping the view service in the `tick()` function
- Once a server knows the current view, it knows if it's the primary, backup, or neither
- Store the current view in the server's state

**Key Code Location**: `tick()` function in `server.go`

### Step 2: Basic Put/Get Handlers
**Goal**: Implement core key-value operations

**Implementation**:
- Implement Put and Get handlers in `pbservice/server.go`
- Store keys and values in a `map[string]string`
- Handle PutHash operations:
  - Use the `hash()` function from `common.go`
  - Hash the concatenation of previous value + new value
  - Convert result using `strconv.Itoa(int(h))`
  - If key doesn't exist, treat previous value as empty string `""`
- Handle PutHash() and Put() similarly (both use `PutArgs`)

**Key Code Location**: `Put()` and `Get()` methods in `server.go`

### Step 3: Primary-Backup Forwarding
**Goal**: Ensure backup sees all updates

**Implementation**:
- Modify Put handler so the primary forwards updates to the backup
- When a server becomes backup in a new view, the primary should send its complete key/value database
- Implement forwarding RPCs for state synchronization

**Key Code Location**: `ForwardPut()`, `ForwardGet()`, `ForwardState()` methods

### Step 4: Client Retry Logic
**Goal**: Make clients resilient to failures

**Implementation**:
- Modify `client.go` so clients keep retrying until they get an answer
- Include enough information in `PutArgs` and `GetArgs` for duplicate filtering
- Implement duplicate detection in the key/value service
- Add client logic to switch from failed primary to new primary

**Key Code Location**: `Get()`, `Put()`, `PutHash()` methods in `client.go`

### Step 5: Primary Discovery and Switching
**Goal**: Handle primary failures gracefully

**Implementation**:
- If current primary doesn't respond or doesn't think it's primary, consult the view service
- Sleep for `viewservice.PingInterval` between retries to avoid CPU overload
- Implement client-side primary discovery logic

**Key Code Location**: Client retry loops and view service pinging

## Success Criteria

You're done when you can pass all the pbservice tests:

```shell
$ cd /path/to/your/lab/src/pbservice
$ go test -v
Test: Single primary, no backup ...
  ... Passed
Test: Add a backup ...
  ... Passed
Test: Primary failure ...
  ... Passed
Test: Kill last server, new one should not be active ...
  ... Passed
Test: at-most-once Put; unreliable ...
  ... Passed
Test: Put() immediately after backup failure ...
  ... Passed
Test: Put() immediately after primary failure ...
  ... Passed
Test: Concurrent Put()s to the same key ...
  ... Passed
Test: Concurrent Put()s to the same key; unreliable ...
  ... Passed
Test: Repeated failures/restarts ...
  ... Put/Gets done ... 
  ... Passed
Test: Repeated failures/restarts; unreliable ...
  ... Put/Gets done ... 
  ... Passed
Test: Old primary does not serve Gets ...
  ... Passed
Test: Partitioned old primary does not complete Gets ...
  ... Passed
PASS
ok      pbservice       113.352s
```

**Note**: You'll see some "method Kill has wrong number of ins" complaints and lots of "rpc: client protocol error" and "rpc: writing response" complaints - these are expected and can be ignored.

## Detailed Implementation Hints

### RPC Design Hints

**Forwarding RPCs**: You'll need to create new RPCs to forward client requests from primary to backup:
- The backup should reject direct client requests but accept forwarded requests
- Consider the difference between client-facing RPCs and internal forwarding RPCs

**State Transfer RPCs**: You'll need RPCs to handle complete key/value database transfer:
- Send the whole database in one RPC (include a `map[string]string` in RPC arguments)
- Think about when to trigger state transfer (view changes, backup initialization)

### Duplicate Detection Hints

**State Replication**: The state to filter duplicates must be replicated along with the key/value state:
- Think carefully about how to coordinate this handoff
- Consider what happens during view transitions
- Ensure duplicate detection works across primary/backup failover

**Client Identification**: You'll need to generate unique identifiers for clients:
- Look at `rand.Intn` for generating random numbers
- Consider using a combination of client ID and operation sequence number

### Network Reliability Hints

**Unreliable Network Tests**: The tester arranges for RPC replies to be lost in "unreliable" tests:
- RPCs are executed by the receiver, but the sender sees no reply
- This simulates real-world network conditions
- Your retry logic must handle these scenarios

### Server Lifecycle Hints

**Graceful Shutdown**: Tests kill servers by setting the `dead` flag:
- Make sure your server terminates correctly when this flag is set
- Ensure proper cleanup of resources and goroutines
- This is critical for test completion

### Performance and Timing Hints

**Test Duration**: It will take more than 100 seconds to run all test cases:
- This is due to the large number of timeouts in this lab
- Be patient - the tests are designed to stress-test your implementation
- Use this time to observe your system's behavior

**Code Size**: The solution should require around 300 lines of code:
- Focus on correctness over optimization
- Don't over-engineer - the core logic is straightforward
- Use the existing skeleton code as a foundation

### Debugging Hints

**View Service Integration**: Even if your view service passed all tests in Part A, it may still have bugs:
- These bugs might only surface in Part B's more complex scenarios
- Test your view service thoroughly before implementing Part B
- Consider edge cases in view transitions

**Incremental Testing**: Test each step of your implementation:
- Start with basic Put/Get operations
- Add primary/backup forwarding
- Implement client retry logic
- Test failure scenarios

### Common Pitfalls

**Split-Brain Prevention**: Ensure only one primary is active at any time:
- Use view numbers to determine primary legitimacy
- Reject operations from servers that aren't the current primary
- Handle view transitions carefully

**State Consistency**: Maintain consistency between primary and backup:
- Forward all operations to backup before committing
- Handle backup failures gracefully
- Ensure state transfer is complete before allowing new operations

**Client Resilience**: Make clients robust to failures:
- Implement exponential backoff for retries
- Handle network partitions gracefully
- Don't get stuck in infinite retry loops

## Reference Implementation Design

Here's a high-level overview of how the system works:

### Client Behavior
1. **Initial Request**: When a client first attempts a Get/Put request, it pings the view service for the current view
2. **Request Routing**: Requests are sent to the Primary server for that view along with unique client and request IDs
3. **Failure Handling**: If the client has trouble communicating with the Primary server:
   - Sleep for a tick interval
   - Ping the view service for the latest view
   - Send the request to the Primary of the latest view

### Server Behavior
1. **Role Enforcement**: Servers that are not the Primary of their current view ignore Put/Get requests
2. **Primary Processing**: When a Primary receives a request:
   - Check the log for the request ID
   - If the ID has been seen before and matches the current Client ID, return the previously returned value
   - If it's a new request, perform the operation, forward to backup, log the request, and return to client
3. **Backup Processing**: Servers that are not the Backup of their current view ignore forwarded requests
4. **State Synchronization**: When a Primary receives a new view from the view service, it forwards its complete state to any backup server

## Getting Started

### Prerequisites
- Complete Lab 2A (View Service) successfully
- Understand RPC communication patterns
- Familiarity with Go concurrency (goroutines, channels, mutexes)

### Development Tips
1. **Start Small**: Begin with basic Put/Get operations without backup forwarding
2. **Test Incrementally**: Add one feature at a time and test thoroughly
3. **Use Debug Output**: Enable debug printing to trace execution flow
4. **Study the Tests**: The test cases provide excellent examples of expected behavior

### Key Files to Focus On
- `pbservice/server.go`: Core server implementation
- `pbservice/client.go`: Client retry and discovery logic
- `pbservice/common.go`: RPC definitions and data structures
- `pbservice/test_test.go`: Comprehensive test suite

### Next Steps
1. Read through the existing skeleton code
2. Understand the RPC definitions in `common.go`
3. Implement the basic server functionality
4. Add client retry logic
5. Test with the provided test suite

Good luck with your implementation! Remember that building distributed systems is challenging, but the concepts you'll learn here are fundamental to modern computing.
