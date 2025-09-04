# "Paxos Made Moderately Complex" Made Moderately Simple: Supplementary Notes

## State Machine Replication

### Core Concept
- **Goal**: Agree on order of operations
- **Think of operations** as a log
- **Challenge**: Ensure all replicas execute operations in same order

## Basic Paxos Review

### Two-Phase Protocol
- **Phase 1**: Send prepare messages, pick value to accept
- **Phase 2**: Send accept messages
- **Problem**: Two round-trips for each operation

## Can We Do Better?

### Optimized Approach
- **Phase 1**: "Leader election" - deciding whose value we will use
- **Phase 2**: "Commit" - leader makes sure it's still leader, commits value
- **Key insight**: Split these phases to enable one round-trip operations

### Benefits
- **Single round-trip**: For normal operations
- **Better performance**: Reduced latency
- **Maintained safety**: Still guarantees consensus

## PMMC Architecture

### Roles in PMMC

#### Replicas (like learners)
- **Keep log** of operations
- **Maintain state machine**
- **Store configurations**
- **Execute operations** in agreed order

#### Leaders (like proposers)
- **Get elected** through consensus
- **Drive consensus protocol**
- **Propose operations** for log slots
- **Handle client requests**

#### Acceptors (simplified!)
- **"Vote" on leaders**
- **Simpler than basic Paxos**
- **Maintain ballot information**

## Ballot System

### Ballot Numbers
- **Format**: (leader, seqnum) pairs
- **Isomorphic**: To previous system we discussed
- **Example**: (Alice, 1), (Bob, 2), (Alice, 3)

### Ballot Properties
- **Ballot numbers increase**: Over time
- **Only accept values**: From current ballot
- **Never remove ballots**: Maintain history
- **Consistency guarantee**: If value v chosen by majority on ballot b, then any value accepted by any acceptor in same slot on ballot b' > b has same value

## Leader Election

### When to Run for Office
- **At beginning of time**: Initial leader election
- **When current leader seems failed**: Based on timeouts and pings
- **Algorithm**: Paper describes ping-based timeout mechanism

### Election Strategy
- **If preempted**: Don't immediately try for election again
- **Avoid thrashing**: Prevent rapid leader changes
- **Stability**: Ensure leader can make progress

## Leader Responsibilities

### Value Proposals
- **Only propose one value** per ballot and slot
- **Consistency guarantee**: If value v chosen by majority on ballot b, then any value proposed by any leader in same slot on ballot b' > b has same value
- **Maintain safety**: Ensure no conflicting values

### Operation Handling
- **Process client requests**
- **Propose operations** for log slots
- **Ensure progress**: Drive consensus forward

## Reconfiguration

### The Challenge
- **All replicas must agree** on who the leader and acceptors are
- **Configuration changes**: Must be coordinated
- **Safety**: Ensure consistent view of system

### Solution: Use the Log
- **Commit special reconfiguration command**
- **New configuration applies** after WINDOW slots
- **Delayed activation**: Ensures all replicas see change

### Handling Edge Cases
- **What if no client requests?**: Commit no-ops until WINDOW is cleared
- **Ensure progress**: System can always reconfigure when needed
- **Maintain liveness**: Prevent deadlock during reconfiguration

## Practical Considerations

### State Simplifications
- **Can track much less information**: Especially on replicas
- **Reduced memory overhead**: Compared to basic Paxos
- **Simpler implementation**: Fewer edge cases to handle

### Garbage Collection
- **Problem**: Unbounded memory growth is bad
- **Solution**: Track finished slots across all instances
- **Garbage collect**: When everyone has learned result
- **Memory management**: Prevent unbounded growth

### Read-Only Commands
- **Can't just read from replica**: Why? (Consistency issues)
- **Don't need their own slot**: Can be handled specially
- **Optimization**: Avoid unnecessary consensus rounds

## Key Questions

### Implementation Details
- **What should be in stable storage?**: Critical state for crash recovery
- **What are the costs?**: Performance implications of Paxos
- **Is it practical enough?**: Real-world applicability

### Trade-offs
- **Performance vs Correctness**: Paxos guarantees vs latency
- **Complexity vs Simplicity**: More complex but more efficient
- **Memory vs Speed**: Storage requirements vs performance

## Key Takeaways

### PMMC Improvements
- **Single round-trip**: For normal operations
- **Simplified acceptors**: Easier to implement
- **Better performance**: Reduced latency
- **Maintained safety**: All Paxos guarantees preserved

### Design Principles
- **Split phases**: Separate leader election from value commitment
- **Use logs**: For reconfiguration and state management
- **Optimize common case**: One round-trip for normal operations
- **Handle edge cases**: Reconfiguration, garbage collection

### Practical Benefits
- **Real-world applicability**: More practical than basic Paxos
- **Performance optimization**: Reduced latency
- **Simplified implementation**: Fewer complex edge cases
- **Memory efficiency**: Better garbage collection

### Applications
- **Distributed databases**: Consistent replication
- **Configuration management**: Coordinated system changes
- **State machine replication**: Ordered operation execution
- **High-availability systems**: Fault-tolerant consensus

### Limitations
- **Still complex**: Despite "moderately complex" name
- **Leader dependency**: Requires stable leader for performance
- **Reconfiguration overhead**: Delayed activation of changes
- **Memory management**: Need for garbage collection