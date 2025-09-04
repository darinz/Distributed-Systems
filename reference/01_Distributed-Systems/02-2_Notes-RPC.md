# Remote Procedure Calls: Supplementary Notes

## Knowledge in Distributed Systems

### Core Concepts
- **Communication** = transforming system's state of knowledge
- **Distributed knowledge**: Knowledge "distributed" among group members
- **Common knowledge**: Fact that is "publicly known"
- **Key insight**: In practical systems, common knowledge cannot be attained

### The Muddy Foreheads Problem

#### Setup
- **n children**, **k get mud** on their foreheads
- Children sit in circle
- **Teacher announces**: "Someone has mud on their forehead"
  - Someone = 1 or more
  - No one can see their own forehead
  - **k is NOT common knowledge**

#### Process
- Teacher repeatedly asks: "Raise your hand if you know you have mud on your forehead"

#### Result
- **First k-1 times**: All children reply "No"
- **k-th time**: All dirty children reply "Yes"

#### Reasoning Examples
- **k = 1**: Child with mud says "Yes" immediately
- **k = 2**: 
  - Each sees exactly one other person with mud
  - In round 1, X notices Y didn't say "Yes"
  - Y must have seen someone with mud → X must have mud

#### The Paradox
> **If k > 1, the teacher didn't say anything anyone didn't already know!**

**Key insight**: The announcement creates common knowledge through the iterative process, even though the information was already distributed.

## Why Distributed Systems Are Hard

### Core Challenges

#### 1. **Asynchrony**
- Different nodes run at different speeds
- Messages can be unpredictably, arbitrarily delayed

#### 2. **Failures (Partial and Ambiguous)**
- Parts of the system can crash
- Can't tell crash from slowness

#### 3. **Concurrency and Consistency**
- Replicated state, cached on multiple nodes
- How to keep many copies of data consistent?

#### 4. **Performance**
- Have to efficiently coordinate many machines
- Performance is variable and unpredictable
- **Tail latency**: Only as fast as slowest machine

#### 5. **Testing and Verification**
- Almost impossible to test all failure cases
- Proofs (emerging field) are really hard

#### 6. **Security**
- Need to assume adversarial nodes

## MapReduce Computational Model

### Core Functions
```python
# For each key k with value v, compute new key-value pairs
map(k, v) -> list(k', v')

# For each key k' and list of values v', compute new values
reduce(k', list(v')) -> list(v'')
```

### Architecture
- **Scheduler**: Accepts MapReduce jobs, finds master and workers
- **MapReduce Master**: Farms tasks to workers, restarts failed jobs, syncs completion
- **Workers**: Execute Map and Reduce tasks
- **Storage**: Stores initial data, intermediate files, end results

**User writes**: map and reduce functions
**Framework handles**: Parallelism, distribution, fault tolerance

## Remote Procedure Call (RPC)

### Concept
**RPC**: A request from client to execute a function on the server
- **To client**: Looks like a procedure call
- **To server**: Looks like an implementation of a procedure call

### RPC Flow

#### Client Side
```python
result = DoMap(worker, i)  # Example call
```
1. **Parameters marshalled** into message (arbitrary types)
2. **Message sent** to server (can be multiple packets)
3. **Wait for reply**

#### Server Side
1. **Message parsed**
2. **Operation invoked**: DoMap(i)
3. **Result marshalled** into message (can be multiple packets)
4. **Message sent** to client

### RPC vs. Procedure Call

#### Equivalents Needed
- **Procedure name** → Function identifier
- **Calling convention** → Parameter marshalling
- **Return value** → Response marshalling
- **Return address** → Client identification

#### Key Differences

##### Binding
- Client needs connection to server
- Server must implement required function
- **Version mismatch**: What if server runs different code version?

##### Performance
- **Procedure call**: ~10 cycles = ~3ns
- **RPC in data center**: 10 microseconds = ~1K slower
- **RPC in wide area**: Millions of times slower

##### Failures
- Messages get dropped?
- Client crashes?
- Server crashes?
- Server crashes after performing op but before replying?
- Server appears to crash but is slow?
- Network partitions?

## RPC Semantics

### Semantics = Meaning
- **Reply == ok** → ???
- **Reply != ok** → ???

### Three Semantics Types

#### 1. **At Least Once** (NFS, DNS)
- **True**: Executed at least once
- **False**: Maybe executed, maybe multiple times

#### 2. **At Most Once**
- **True**: Executed once
- **False**: Maybe executed, but never more than once

#### 3. **Exactly Once**
- **True**: Executed once
- **False**: Never returns false

## At Least Once Implementation

### Process
1. **RPC library waits** for response
2. **If none arrives** → re-send request
3. **Do this a few times**
4. **Still no response** → return error to application

### Example: Non-replicated Key/Value Server
- Client sends `Put k v`
- Server gets request, but network drops reply
- Client sends `Put k v` again
- **Question**: Should server respond "yes" or "no"?
- **Problem**: What if operation is "append"?

### Does TCP Fix This?
**TCP provides**:
- Reliable bi-directional byte stream
- Retransmission of lost packets
- Duplicate detection

**But what if**:
- TCP times out and client reconnects?
- Browser connects to Amazon
- RPC to purchase book
- WiFi times out during RPC
- Browser reconnects

### When At-Least-Once Works
- **No side effects**: Read-only operations
- **Idempotent operations**: Operations that can be repeated safely
- **Examples**: MapReduce, NFS (`readFileBlock`, `writeFileBlock`)

## At Most Once Implementation

### Process
1. **Client includes unique ID (UID)** with each request
2. **Use same UID** for re-send
3. **Server detects duplicate requests**
4. **Return previous reply** instead of re-running handler

### UID Generation Issues
- **How to ensure UID is unique?**
  - Big random number?
  - Combine unique client ID (IP address) with sequence number?
- **What if client crashes and restarts?** Can it reuse the same UID?
- **Solution**: Every node gets new ID on start

### When Can Server Discard Old RPCs?

#### Option 1: Never
- **Problem**: Memory grows indefinitely

#### Option 2: Client Acknowledgments
- Unique client IDs
- Per-client RPC sequence numbers
- Client includes "seen all replies ≤ X" with every RPC

#### Option 3: One Outstanding RPC
- Only allow client one outstanding RPC at a time
- Arrival of seq+1 allows server to discard all ≤ seq

### Server Crash Problem
**Issue**: If at-most-once list of recent RPC results is stored in memory, server will forget and accept duplicate requests when it reboots

**Questions**:
- Does server need to write recent RPC results to disk?
- If replicated, does replica also need to store recent RPC results?

## Key Takeaways

### Knowledge and Common Knowledge
- **Distributed knowledge** vs. **common knowledge**
- **Muddy foreheads** illustrates how common knowledge emerges
- **Practical systems** cannot achieve true common knowledge

### RPC Challenges
- **Performance**: Orders of magnitude slower than local calls
- **Failure handling**: Multiple failure modes to consider
- **Semantics**: Choose appropriate guarantee (at-least-once, at-most-once, exactly-once)
- **Implementation complexity**: UID generation, duplicate detection, crash recovery

### Design Principles
- **Choose semantics** based on operation characteristics
- **Handle failures** explicitly in design
- **Consider performance** implications of reliability mechanisms
- **Plan for crashes** and restarts in all components