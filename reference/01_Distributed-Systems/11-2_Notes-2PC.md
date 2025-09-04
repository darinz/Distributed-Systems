# Two-Phase Commit: Supplementary Notes

## Implications of Two Generals

### The Problem
- **Cannot get agreement** in distributed system to perform action at same time
- **What if we want** to update data stored in multiple locations?
- **Linearizable fashion**: Consistent with sequential order
- **Solution**: Perform group of operations at logical instant in time, not physical instant

### Key Insight
- **Logical time**: More important than physical time
- **Atomic updates**: All or nothing across multiple locations
- **Consistency**: Maintained through coordination

## Setting and Requirements

### Use Cases
- **Multikey update**: To sharded key-value store
- **Bank transfer**: Between different accounts
- **Distributed transactions**: Across multiple systems

### Requirements
- **Atomicity**: All or none - either all operations succeed or all fail
- **Linearizability**: Consistent with sequential order
- **No stale reads**: Readers see consistent state
- **No write buffering**: Changes are immediately visible

### Assumption
- **For now, ignore availability**: Focus on correctness first

## One Phase Commit?

### Simple Approach
- **Central coordinator decides**, tells everyone else
- **Problem**: What if some participants can't do the request?
  - Bank account has zero balance
  - Bank account doesn't exist
  - Resource unavailable

### Challenges
- **How to get atomicity/linearizability?**
- **Need to apply changes** at same logical point in time
- **Need all other changes** to appear before/after
- **Acquire read/write lock** on each location
- **For linearizability**: Need read/write lock on all locations at same time

## Two-Phase Commit Protocol

### Basic Structure
- **Central coordinator asks**: Can you commit this transaction?
- **Participants commit to commit**:
  - Acquire any necessary locks
  - In meantime, no other operations allowed on that key
  - Delay other concurrent 2PC operations
- **Central coordinator decides**, tells everyone else:
  - Release locks
  - Apply changes or abort

### Calendar Example
- **Doug has three advisors**: Tom, Zach, Mike
- **Want to schedule meeting**: With all of them
- **Try Tuesday at 11**: People usually free then
- **Calendars on different nodes**: Distributed system
- **Other students scheduling**: Concurrent operations
- **Nodes can fail**: Messages can be dropped

## Atomic Commit Protocol (ACP)

### Properties
- **Every node arrives** at same decision
- **Once node decides**: It never changes
- **Transaction committed**: Only if all nodes vote Yes
- **Normal operation**: If all processes vote Yes, transaction committed
- **Failure recovery**: If all failures eventually repaired, transaction eventually committed or aborted

### Roles
- **Participants** (Mike, Tom, Zach): Nodes that must update data relevant to transaction
- **Coordinator** (Doug): Node responsible for executing protocol (might also be participant)

### Messages
- **Prepare**: "Can you commit this transaction?"
- **Commit**: "Commit this transaction"
- **Abort**: "Abort this transaction"

## Failure Scenarios

### Simple Case
- **In absence of failures**: 2PC is pretty simple!
- **Interesting failures**:
  - Participant failures
  - Coordinator failures
  - Message drops

### Do We Need the Coordinator?
- **Can participants decide** amongst themselves?
- **Yes, if participants can know** for certain that coordinator has failed
- **Problem**: What if coordinator is just slow?
  - Participants decide to commit!
  - Coordinator times out, declares abort!
  - **Inconsistency**: Different decisions

## Blocking Protocol

### Definition
- **Blocking protocol**: Cannot make progress if some participants are unavailable
- **Unavailable**: Either down or partitioned
- **Fault-tolerance**: But not availability
- **Fundamental limitation**: Cannot be avoided

### Implications
- **2PC is blocking**: Cannot proceed if any participant fails
- **Trade-off**: Correctness vs availability
- **Real-world impact**: System can get stuck

## Making 2PC Non-Blocking

### Paxos Solution
- **Paxos is non-blocking**: Can make progress despite failures
- **Use Paxos to update** individual keys
- **Can we use Paxos** to update multiple keys?
  - **Same shard**: Easy
  - **Different shards**: More complex

### Two-Phase Commit on Paxos

#### Process
1. **Client requests** multi-key operation at coordinator
2. **Coordinator logs request**:
   - Paxos: Available despite node failures
3. **Coordinator sends prepare**
4. **Replicas decide** to commit/abort, log result:
   - Paxos: Available despite node failures
5. **Coordinator collects replies**, log result:
   - Paxos: Available despite node failures
6. **Coordinator sends** commit/abort
7. **Replicas record result**:
   - Paxos: Available despite node failures

#### Benefits
- **Non-blocking**: Can make progress despite failures
- **Consistent**: All nodes reach same decision
- **Available**: System doesn't get stuck
- **Fault-tolerant**: Handles node failures

## Key Takeaways

### 2PC Properties
- **Atomic**: All or nothing execution
- **Consistent**: All nodes reach same decision
- **Blocking**: Cannot proceed if any participant fails
- **Simple**: Easy to understand and implement

### Design Principles
- **Two phases**: Prepare then commit/abort
- **Central coordinator**: Makes final decision
- **Participant voting**: Can vote yes or no
- **Lock acquisition**: During prepare phase

### Trade-offs
- **Correctness vs Availability**: 2PC prioritizes correctness
- **Simplicity vs Performance**: Simple but can block
- **Centralized vs Distributed**: Coordinator is single point of failure
- **Blocking vs Non-blocking**: 2PC is inherently blocking

### Applications
- **Distributed databases**: ACID transactions
- **Banking systems**: Money transfers
- **Distributed file systems**: Atomic updates
- **Microservices**: Distributed transactions

### Limitations
- **Blocking**: Cannot proceed with failures
- **Single coordinator**: Potential bottleneck
- **Performance**: Two round-trips required
- **Availability**: System can get stuck

### Improvements
- **Paxos-based 2PC**: Non-blocking variant
- **Three-phase commit**: Reduces blocking
- **Optimistic approaches**: Reduce coordination overhead
- **Saga pattern**: Alternative to 2PC for microservices

### When to Use
- **Strong consistency**: When ACID properties required
- **Simple systems**: When complexity is acceptable
- **Low failure rates**: When blocking is acceptable
- **Traditional databases**: When 2PC is built-in

### When Not to Use
- **High availability**: When system cannot block
- **High performance**: When latency is critical
- **Microservices**: When distributed transactions are complex
- **Eventual consistency**: When strong consistency not needed