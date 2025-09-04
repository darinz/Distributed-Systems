# Multi-Writer, Multi-Reader Atomic Registers: Supplementary Notes

## Drawbacks of Paxos

### Performance Issues
- **Leader bottleneck**: Single point of contention
- **O(n) messages**: On every request
- **High latency**: Due to leader coordination

### Liveness Problems
- **FLP impossibility**: Liveness not guaranteed
- **Bad availability**: During failure scenarios
- **Leader election delay**: Takes time to elect new leader when current fails

## Alternatives to Paxos

### Options
- **Allow randomness**: See Ben-Or algorithm
- **Weaken safety guarantees**: Accept weaker consistency (risky)
- **Constrain the problem**: Focus on specific use cases

### Register Approach
- **Simpler problem**: Read/write operations only
- **Wait-free guarantees**: Better liveness properties
- **No consensus required**: For basic operations

## Register Semantics

### Definition
- **Hold single value**: Want multiple values? Use multiple registers
- **Operations**: Read and write only
- **No appends**: Or other read-modify-write operations
- **Semantics**: Safe, regular, and atomic/linearizable
- **Goal**: Linearizability

### Why No Appends?
- **Simple consensus implementation**:
  1. All processes append their input value
  2. All processes read the value
  3. They all decide the first value that was appended
- **Impossibility**: If you can wait-free implement appendable register, you can solve consensus
- **FLP result**: Consensus is impossible in asynchronous systems

## Implementation Model

### Client/Server Architecture
- **Servers**: Replicas storing the value
- **Clients**: Send reads and writes
- **Goal**: Linearizability of reads and writes
- **Fault tolerance**: Up to f server crash failures
- **Client failures**: Can also fail by crashing

## Non-Blocking Algorithms

### Lock-Free vs Wait-Free
- **Lock-free**: Guarantees system-wide progress
- **Wait-free**: Guarantees per-client progress
- **Wait-free property**: No matter what other processes do, correct client's operations complete in finite steps

### Benefits
- **No deadlock**: Operations always complete
- **No starvation**: Individual clients make progress
- **Better liveness**: Compared to blocking algorithms

## Server Requirements

### How Many Servers?
- **Progress requirement**: Can wait for at most n - f responses
- **Write requirement**: Must send to > f replicas (otherwise lost forever)
- **Minimum servers**: 2f + 1
- **Quorum overlap**: Read quorum + write quorum > n
- **Simple solution**: Use majorities

### Quorum Properties
- **Majorities intersect**: Ensures consistency
- **Fault tolerance**: Can handle f failures
- **Progress**: Always have majority available

## Single Reader, Single Writer (SRSW)

### Basic Algorithm
- **Writer**: Sends value to majority
- **Reader**: Reads value from majority
- **Consistency**: Since majorities intersect, read gets writer's value

### Timestamped Version
- **Writer**: Sends timestamped value to majority
- **Reader**: Reads from majority, takes highest timestamp
- **Consistency**: Majorities intersect, so reader gets writer's value

### Server Algorithm
- **On write**: Update local timestamp and value if write's timestamp is greater; send ack
- **On read**: Respond with local timestamp and value

### Client Algorithms
- **Writer**: Increment local timestamp, send to all, wait for majority acks
- **Reader**: Read from majority, take value with highest timestamp

## Multiple Readers, Single Writer (MRSW)

### The Problem
- **Multiple reads**: By different processes overlapping same write
- **Consistency issues**: Readers might get different values
- **Write-back solution**: Readers write back the value they read

### Write-Back Algorithm
1. **Reader reads** value from majority, takes highest timestamp
2. **Reader performs write-back**: Writes value to majority
3. **Return value**: After write-back is complete
4. **Later readers**: Guaranteed to read value at least as new

### Why Write-Back?
- **Ensures consistency**: All readers see same value
- **Prevents stale reads**: Later reads get fresh values
- **Sequential consistency**: May not need write-back for weaker guarantees

## Multiple Readers, Multiple Writers (MRMW)

### Challenges
- **Same timestamps**: Writers might use identical timestamps
- **Timestamp ordering**: Write starting after previous write might use smaller timestamp
- **Tie-breaking**: Need consistent ordering

### Solutions
- **Writer IDs**: Break ties using writer ID (same as PMMC)
- **Timestamp ordering**: Ensure timestamps increase over time

### Ensuring Timestamp Ordering
1. **Writer queries majority**: Updates timestamp to be larger than largest found
2. **Writer writes value**: To majority as usual
3. **Guarantee**: Written value has timestamp larger than previously written values
4. **Readers**: Will read latest value (writer IDs break ties)

## ABD Algorithm (Attiya, Bar-Noy, Dolev 1995)

### Key Insight
- **Read and write methods**: Are essentially the same
- **Only difference**: Read writes and returns value that was read, write writes value to be written
- **Processes**: Can be both readers and writers

### Algorithm Structure
- **Two-phase protocol**: Query then write
- **Majority quorums**: For both phases
- **Timestamp ordering**: Ensures consistency
- **Wait-free**: Guarantees progress

## ABD vs Paxos

### ABD Advantages
- **Wait-freedom**: Guarantees progress even with multiple writers
- **No leader bottleneck**: Removes single point of contention
- **Same latency**: As leader-based Paxos
- **Better liveness**: No FLP impossibility issues

### Paxos Advantages
- **Arbitrary state machines**: Can support complex operations
- **Consensus**: Solves general agreement problem
- **Mature**: Well-understood and widely implemented

### Trade-offs
- **ABD**: Simpler, wait-free, but limited to read/write
- **Paxos**: More general, but has liveness issues and leader bottleneck

## Applications

### What Can We Do With Registers?
- **Read/write key-value store**: Implement distributed storage
- **Emulate shared memory**: Build distributed data structures
- **Coordination primitives**: Without full consensus

### Key Insight
> **Consensus isn't always the right problem! Don't solve it if you don't have to!**

## Key Takeaways

### Register Benefits
- **Wait-free**: Better liveness than Paxos
- **No leader**: Removes bottleneck
- **Simple**: Read/write operations only
- **Efficient**: Same latency as Paxos

### Design Principles
- **Use majorities**: For quorum overlap
- **Timestamp ordering**: Ensure consistency
- **Write-back**: For multiple readers
- **Two-phase protocol**: Query then write

### Trade-offs
- **Simplicity vs Generality**: Registers are simpler but less general
- **Liveness vs Complexity**: Wait-free but limited operations
- **Performance vs Features**: Better performance but fewer features

### When to Use
- **Read/write workloads**: Perfect fit for registers
- **High availability**: When liveness is critical
- **Simple coordination**: When consensus is overkill
- **Performance critical**: When latency matters

### Limitations
- **Limited operations**: Only read/write, no appends
- **No consensus**: Can't solve general agreement
- **Memory overhead**: Need to store timestamps
- **Complexity**: Still non-trivial to implement correctly