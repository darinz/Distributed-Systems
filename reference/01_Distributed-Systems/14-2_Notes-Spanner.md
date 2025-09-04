# Spanner: Supplementary Notes

## BigTable in Retrospect

### Success and Limitations
- **Definitely useful, scalable system**: Still in use at Google, motivated lots of NoSQL DBs
- **Biggest mistake in design** (per Jeff Dean, Google): Not supporting distributed transactions!
- **Became really important**: With incremental updates
- **Users wanted them**: Implemented themselves, often incorrectly!
- **At least 3 papers later**: Fixed this limitation

### Key Insight
- **Distributed transactions**: Critical for many applications
- **User demand**: Applications needed cross-shard consistency
- **Implementation complexity**: Users often got it wrong
- **Research opportunity**: Multiple papers addressed this gap

## Two-Phase Commit Recap

### Basic Process
- **Keys partitioned**: Over different hosts; one coordinator per transaction
- **Acquire locks**: On all data read/written; release after commit
- **To commit**: Coordinator first sends prepare message to all shards; they respond prepare_ok or abort
  - **If prepare_ok**: They must be able to commit transaction
- **If all prepare_ok**: Coordinator sends commit to all; they write commit record and release locks

### 2PC Problems
- **Availability**: What do we do if either some shard or the coordinator fails?
  - **Generally**: 2PC is a blocking protocol, can't make progress until it comes back up
- **Performance**: Can we really afford to take locks and hold them for the entire commit process?

## Spanner Overview

### Purpose
- **Backend for F1 database**: Which runs the ad system
- **Basic model**: 2PC over Paxos
- **Uses physical clocks**: For performance
- **Distributed transactions**: With global consistency

### Key Innovation
- **2PC over Paxos**: Combines transaction coordination with fault tolerance
- **Physical clocks**: Enable lock-free read-only transactions
- **Global consistency**: Across distributed shards

## Example: Social Network

### Scenario
- **Simple schema**: User posts and friends lists
- **But sharded**: Across thousands of machines
- **Each replicated**: Across multiple continents
- **Example**: Generate page of friends' recent posts

### Problem
- **What if I remove friend X, post mean comment?**
- **Maybe he sees**: Old version of friends list, new version of my posts?
- **Inconsistency**: Different views of the same data

### Solution with Locking
- **Acquire read locks**: On friends list, and on each friend's posts
- **Prevents modification**: Concurrent changes
- **But potentially really slow**: Lock contention and latency

## Spanner Architecture

### Basic Structure
- **Each shard stored**: In a Paxos group
  - **Replicated across data centers**
  - **Has a (relatively long-lived) leader**
- **Transactions span Paxos groups**: Using 2PC
  - **Use 2PC for transactions**
  - **Leader of each Paxos group tracks locks**
  - **One group leader becomes the 2PC coordinator, others participate**

### Key Components
- **Paxos groups**: Fault-tolerant shard replication
- **2PC coordination**: Cross-shard transaction management
- **Lock management**: Per-shard lock tracking
- **Leader election**: Within each Paxos group

## Basic 2PC/Paxos Approach

### Transaction Execution
1. **During execution**: Read and write objects
   - **Contact appropriate Paxos group leader**: Acquire locks
2. **Client decides to commit**: Notifies the coordinator
   - **Coordinator contacts all shards**: Sends PREPARE message
   - **They Paxos-replicate**: Prepare log entry (including locks)
   - **Vote either ok or abort**
3. **If all shards vote OK**: Coordinator sends commit message
   - **Each shard Paxos-replicates**: Commit entry
   - **Leader releases locks**

### Key Insight
- **Same as basic 2PC**: From before
- **Just replaced writes**: To a log on disk with writes to a Paxos replicated log!
- **It is linearizable**: (= strict serializable = externally consistent)
- **So what's left?**: Lock-free read-only transactions

## Lock-Free Read-Only Transactions

### Key Idea
- **Assign meaningful timestamp**: To transaction
  - **Such that timestamps**: Are enough to order transactions meaningfully
- **Keep history of versions**: Around on each node
- **Then, reasonable to say**: Read-only transaction X reads at timestamp 10

### Benefits
- **No locks required**: For read-only transactions
- **Consistent snapshots**: At specific timestamps
- **High performance**: No lock contention
- **Global consistency**: Across all shards

## Spanner Topics

### Core Areas
1. **Distributed transactions, in detail**:
   - **On one Paxos group**
   - **Between Paxos groups**
2. **Fast read-only transactions**: With TrueTime
3. **Discussion**: Performance and limitations

### Design Goals
- **Fast reads**: Lock-free read-only transactions
- **Distributed transactions**: Cross-shard consistency
- **Performance**: Minimize latency and lock contention

## How Can We Get Fast Reads?

### Problem
- **R/W transactions are complicated**: And slow
- **Can we do fast, lock-free reads?**: Real time to the rescue

### Solution
- **TrueTime**: Physical clock synchronization
- **Timestamp-based ordering**: Global transaction ordering
- **Lock-free reads**: Using consistent snapshots

## TrueTime

### API
- **Exposes real time**: With uncertainty
- **{earliest: e, latest: l}**: == TT.now()
- **"Real time" is between**: Earliest and latest
- **Time is an illusion!**: But a useful one

### Guarantees
- **If I call TT.now()**: On two nodes simultaneously, intervals guaranteed to overlap!
- **If intervals don't overlap**: The later one happened later!
- **Global ordering**: Based on time intervals

### Implementation
- **GPS, atomic clocks**: High-precision time sources
- **All local clocks synced**: With masters, and expose uncertainty to local apps
- **Assumptions made**: About local clock drift

## TrueTime Usage

### Transaction Timestamps
- **Assign timestamp**: To each transaction
  - **At each Paxos group**: Timestamp increases monotonically
  - **Globally**: If T1 returns before T2 starts, timestamp(T1) < timestamp(T2)
- **Timestamp for RW transaction**: Chosen by coordinator leader
- **Timestamps for R/W transactions**: Max of:
  - **Local time**: When client request reached coordinator
  - **Prepare timestamps**: At every participant
  - **Timestamp of any previous local transaction**

### Global Ordering
- **Monotonic timestamps**: Within each Paxos group
- **Global ordering**: Across all groups
- **Consistent snapshots**: At any timestamp

## Commit Wait

### Requirement
- **Need to ensure**: That all future transactions will get a higher timestamp
- **Therefore, need to wait until**: TT.now() > transaction timestamp
- **And only then release locks**

### Performance Impact
- **What does this mean for performance?**
- **Larger TrueTime uncertainty bound**: → longer commit wait
- **Longer commit wait**: → locks held longer → can't process conflicting transactions → lower throughput
- **i.e., if time is less certain**: Spanner is slower!

### Trade-off
- **Time precision**: vs. performance
- **Uncertainty bound**: vs. commit wait time
- **Lock duration**: vs. throughput

## What Does This Buy Us?

### Read-Only Transactions
- **Can now do a read-only transaction**: At a particular timestamp, have it be meaningful
- **Example**: Pick a timestamp T in the past, read version w/ timestamp T from all shards
  - **Since T is in the past**: They will never accept a transaction with timestamp < T
  - **Don't need locks**: While we do this!
- **What if we want the current time?**: Use latest available timestamp

### Benefits
- **Lock-free reads**: No lock contention
- **Consistent snapshots**: Global consistency
- **High performance**: Fast read-only transactions
- **Global ordering**: Meaningful timestamps

## TrueTime Implementation

### Hardware
- **GPS, atomic clocks**: High-precision time sources
- **All local clocks synced**: With masters, and expose uncertainty to local apps
- **Assumptions made**: About local clock drift

### Software
- **Uncertainty bounds**: Exposed to applications
- **Clock synchronization**: Across all nodes
- **Drift assumptions**: Based on hardware characteristics

## What If TrueTime Fails?

### Google's Argument
- **Picked using engineering considerations**: Less likely than a total CPU failure
- **But what if it went wrong anyway?**
  - **Can cause very long commit wait periods**
  - **Can break ordering guarantees**: No longer externally consistent
  - **But system will always be serializable**: Gathering many timestamps and taking the max is a Lamport clock

### Fallback Behavior
- **Serializable**: Always maintained
- **External consistency**: May be lost
- **Performance**: May degrade significantly
- **Lamport clock**: Fallback ordering mechanism

## Conclusions

### What's Cool About Spanner?
- **Distributed transactions**: With decent performance
  - **What makes that possible?**: 2PC over Paxos
- **Read-only transactions**: With great performance
  - **What makes that possible?**: TrueTime and lock-free reads

### Key Insight
- **Clocks are a form of communication!**: Time enables global ordering
- **Physical clocks**: Enable lock-free read-only transactions
- **2PC over Paxos**: Enables fault-tolerant distributed transactions
- **Global consistency**: Across distributed shards

## Key Takeaways

### Spanner Design Principles
- **2PC over Paxos**: Combines transaction coordination with fault tolerance
- **Physical clocks**: Enable global ordering and lock-free reads
- **Distributed transactions**: Cross-shard consistency
- **Lock-free read-only transactions**: High performance for reads
- **Global consistency**: External consistency guarantees

### Architecture Benefits
- **Fault tolerance**: Paxos replication
- **Distributed transactions**: Cross-shard consistency
- **High performance**: Lock-free read-only transactions
- **Global ordering**: Meaningful timestamps
- **External consistency**: Strong consistency guarantees

### TrueTime Innovation
- **Physical clock synchronization**: GPS and atomic clocks
- **Uncertainty bounds**: Exposed to applications
- **Global ordering**: Based on time intervals
- **Lock-free reads**: Using consistent snapshots
- **Performance trade-off**: Time precision vs. commit wait

### Transaction Model
- **Read-write transactions**: Use 2PC over Paxos
- **Read-only transactions**: Lock-free using TrueTime
- **Global timestamps**: Meaningful ordering
- **Consistent snapshots**: At any timestamp
- **External consistency**: Strong guarantees

### Performance Characteristics
- **Read-write transactions**: Slower due to 2PC and locking
- **Read-only transactions**: Fast and lock-free
- **Commit wait**: Depends on TrueTime uncertainty
- **Lock duration**: Affects throughput
- **Time precision**: Critical for performance

### Limitations
- **TrueTime dependency**: System performance depends on clock precision
- **Commit wait**: Can be long with high uncertainty
- **Complexity**: More complex than simple 2PC
- **Hardware requirements**: GPS and atomic clocks needed
- **Failure modes**: TrueTime failures can degrade performance

### Modern Relevance
- **Cloud databases**: Google Cloud Spanner, Amazon Aurora
- **Distributed transactions**: Foundation for modern systems
- **Global consistency**: Important for many applications
- **Time-based ordering**: Used in many distributed systems
- **Lock-free algorithms**: Influence on modern database design

### Lessons Learned
- **Physical clocks**: Can enable global ordering
- **2PC over Paxos**: Combines transaction coordination with fault tolerance
- **Lock-free reads**: Can provide high performance
- **Time precision**: Critical for performance
- **Global consistency**: Possible with careful design
- **Trade-offs**: Performance vs. consistency vs. complexity