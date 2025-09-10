# Raft: A Consensus Algorithm for Replicated Logs

## Introduction: The Consensus Problem

### Why Consensus Matters

**Machines can work together as a coherent group if they are able to achieve consensus**

Consensus is the fundamental problem in distributed systems: how do multiple machines agree on a single value or sequence of values, even in the presence of failures? This seemingly simple problem is actually quite complex and has been the subject of decades of research.

### The Landscape of Consensus Algorithms

**Many different consensus algorithms exist**:

- **Paxos**: The classic consensus algorithm, widely studied but notoriously difficult to understand
- **Viewstamped Replication**: An alternative approach that predates Paxos
- **Byzantine fault tolerance**: Algorithms that handle malicious failures
- **Raft**: The focus of this paper, designed for understandability
- **And many others**: Each with different trade-offs and characteristics

### The Paxos Problem

**Before Raft, Paxos was synonymous with consensus**

#### Why Paxos is Problematic

**Paxos is notoriously difficult to understand**:

**The algorithm is complex**:
- **Multiple roles**: Proposers, acceptors, and learners
- **Complex interactions**: Intricate message passing between roles
- **Subtle correctness conditions**: Easy to get wrong in implementation

**The original paper is not written in a clear and concise manner**:
- **Academic style**: Focuses on mathematical proofs rather than practical implementation
- **Missing details**: Leaves many implementation details unspecified
- **Poor examples**: Lacks clear, concrete examples of how the algorithm works

#### The Practical Impact

**The result**: While Paxos is theoretically sound, it's extremely difficult to implement correctly in real systems. Many systems that claim to use Paxos actually use simplified or incorrect versions.

### Raft's Mission: Understandability

**Main goal of Raft is to improve understandability**

#### Specific Objectives

**Make Raft easier to learn**:
- **Clear structure**: Algorithm should be easy to understand from the ground up
- **Good examples**: Concrete examples that illustrate how the algorithm works
- **Intuitive reasoning**: Design decisions should make intuitive sense

**Make it easier to develop intuition about Raft so that it can be used and extended when building real systems**:
- **Practical focus**: Algorithm should be designed with real-world implementation in mind
- **Extensibility**: Should be easy to extend for specific use cases
- **Debugging**: Should be easy to reason about when things go wrong

### Raft's Key Design Techniques

**Raft has two key techniques to improve understandability**:

#### (1) Decomposition

**Separate leader election, log replication, and safety**:
- **Leader election**: How to choose a leader when the current leader fails
- **Log replication**: How the leader replicates log entries to followers
- **Safety**: How to ensure the algorithm maintains correctness properties

**Benefits of decomposition**:
- **Easier to understand**: Each component can be understood independently
- **Easier to implement**: Each component can be implemented and tested separately
- **Easier to debug**: Problems can be isolated to specific components

#### (2) State Space Reduction

**Reduce nondeterminism**:
- **Fewer choices**: Algorithm makes fewer arbitrary decisions
- **Clearer logic**: Each decision point has clear reasoning
- **Predictable behavior**: Algorithm behavior is more predictable

**Reduce number of ways servers can be inconsistent with each other**:
- **Simplified state**: Servers have fewer possible states
- **Clearer transitions**: State transitions are more obvious
- **Easier recovery**: Inconsistent states are easier to detect and fix

### Novel Features in Raft

**Raft introduces several novel features that improve both understandability and performance**:

#### Strong Leader

**Logs flow only from leader to other servers**:
- **Unidirectional flow**: All log entries come from the leader
- **Simplified protocol**: Followers only receive, never send log entries
- **Clear responsibility**: Leader has complete control over the log

#### Leader Election

**Randomized timers avoid conflicts**:
- **Split vote prevention**: Random timeouts prevent indefinite split votes
- **Fast convergence**: Elections typically complete quickly
- **Simple logic**: Election process is straightforward to understand

#### Membership Changes

**Servers can be added to or removed from the cluster without harming performance with joint consensus**:
- **Online reconfiguration**: Cluster can be reconfigured without downtime
- **Safety guarantees**: Joint consensus ensures no split-brain scenarios
- **Practical necessity**: Real systems need to handle membership changes
## Replicated State Machines

### The Foundation of Distributed Systems

**"State machines on a collection of servers compute identical copies of the same state and can continue operating even if some of the servers are down"**

Replicated state machines are the fundamental building blocks of many distributed systems. The idea is simple but powerful: if multiple servers can maintain identical copies of the same state, the system can continue operating even when individual servers fail.

### Real-World Examples

#### Systems That Use Replicated State Machines

**Examples of systems that use replicated state machines**:
- **GFS (Google File System)**: Uses replication for metadata consistency
- **HDFS (Hadoop Distributed File System)**: Replicates file system metadata
- **RAMCloud**: Uses replication for in-memory data storage

#### Examples of Replicated State Machines

**Examples of replicated state machines**:
- **Chubby**: Google's lock service, provides distributed locking and configuration
- **ZooKeeper**: Apache's coordination service, provides distributed synchronization primitives

### The Role of Consensus

**Consensus algorithm keeps the replicated log consistent across machines**

**Multiple servers can have the appearance of being a single state machine**

The key insight is that if all servers apply the same sequence of commands to their state machines, they will all end up in the same state. The consensus algorithm ensures that all servers agree on the same sequence of commands.

### Properties of Consensus Algorithms

**Consensus algorithms generally have these properties**:

#### Safety

**Never return an incorrect result**:
- **Correctness guarantee**: Algorithm must never produce wrong results
- **Robustness**: Must do this even when there are network delays, partitions, packet loss, packet duplication, and packet reordering
- **Consistency**: All servers must agree on the same sequence of commands

#### Availability

**Consensus algorithm must work as long as a majority of the servers are functioning and can communicate with each other and with clients**:
- **Fault tolerance**: System continues operating despite individual server failures
- **Majority requirement**: Only needs a majority of servers to be operational
- **Network connectivity**: Servers must be able to communicate with each other and clients

**Raft assumes that when a server has an issue, the server just stops**:
- **Fail-stop model**: Server doesn't start sending arbitrary values out on the network or act maliciously
- **Byzantine vs. crash failures**: These cases are covered by Byzantine fault tolerant protocols
- **Simplified model**: Makes the algorithm easier to understand and implement

#### Timing Independence

**Do not depend on timing**:
- **Network delays**: Network delays, faulty clocks, etc. do not impact the correctness of the protocol
- **Performance impact**: At worst, these issues just slow the algorithm down
- **Robustness**: Algorithm remains correct even with timing variations

#### Performance

**Common case is fast**:
- **Majority requirement**: Only a majority of the servers need to respond for a command to be committed
- **Fault tolerance**: A small number of slow servers has no impact
- **Efficiency**: System can make progress even with some servers being slow or unresponsive
## What's Wrong with Paxos?

### Paxos Basics

**Paxos defines a protocol that can reach consensus on a single decision: single-decree Paxos**

**Multiple instances of this protocol can reach consensus on a series of decisions: multi-Paxos**

Paxos is theoretically sound and has been proven correct. The basic idea is elegant: it can reach consensus on a single value, and multiple instances can be combined to reach consensus on a sequence of values.

### Paxos Works, But...

**Paxos works properly and is efficient**

Paxos is not fundamentally broken. It's a correct algorithm that can achieve consensus efficiently. However, it has significant practical problems that make it difficult to use in real systems.

### The Downsides of Paxos

#### 1. Extreme Difficulty of Understanding

**Paxos is very difficult to understand**:
- **Complex roles**: Proposers, acceptors, and learners with intricate interactions
- **Subtle correctness conditions**: Easy to get wrong in implementation
- **Poor documentation**: Original paper is not written in a clear and concise manner
- **Implementation challenges**: Many systems that claim to use Paxos actually use simplified or incorrect versions

#### 2. Lack of Standard Multi-Paxos

**Lack of widely agreed-upon algorithm for multi-Paxos, so it is difficult to build systems with Paxos**:
- **No standard**: Different implementations use different approaches
- **Incompatibility**: Systems using different multi-Paxos variants can't interoperate
- **Implementation risk**: Each implementation must figure out the details independently

**Plus multi-Paxos basically just combines multiple individual decisions from single-decree Paxos, which is an inefficient and confusing way to build a log**:
- **Inefficient**: Each log entry requires a separate consensus decision
- **Confusing**: The relationship between individual decisions and the overall log is not clear
- **Complex**: Managing multiple consensus instances adds significant complexity

#### 3. Performance Issues

**Servers talk to each other to form a decision without a leader… this harms performance**:
- **No leader**: All servers participate in each decision
- **Message complexity**: More messages required for each decision
- **Coordination overhead**: Servers must coordinate with each other for every decision
- **Slower consensus**: More time required to reach consensus on each decision

### The Practical Impact

**The result**: While Paxos is theoretically elegant, it's extremely difficult to implement correctly and efficiently in real systems. This has led to many systems using simplified or incorrect versions of Paxos, which can lead to bugs and inconsistencies.
## Designing for Understandability

### Raft's Design Goals

**Raft was designed with four key goals in mind**:

#### 1. Complete and Practical Foundation

**Must be complete and practical foundation for building a system**:
- **Complete**: Algorithm must handle all cases and edge conditions
- **Practical**: Must be implementable in real systems
- **Foundation**: Must provide a solid base for building distributed systems

#### 2. Safety and Availability

**Safe and available under normal conditions**:
- **Safety**: Algorithm must never produce incorrect results
- **Availability**: System must continue operating despite failures
- **Normal conditions**: Must work reliably in typical operating conditions

#### 3. Performance

**Efficient in the common case**:
- **Common case optimization**: Must be fast for typical operations
- **Reasonable performance**: Must not be significantly slower than alternatives
- **Scalability**: Performance should not degrade significantly with system size

#### 4. Understandability

**Good understandability**:
- **Easy to learn**: Algorithm should be easy to understand from the ground up
- **Intuitive**: Design decisions should make intuitive sense
- **Debuggable**: Should be easy to reason about when things go wrong
- **Implementable**: Should be easy to implement correctly

### The Trade-off

**The key insight**: Raft prioritizes understandability without sacrificing the other goals. This makes it easier to build correct, reliable systems while maintaining good performance and safety properties.
## The Raft Consensus Algorithm

### The Strong Leader Approach

**Raft elects a leader and the leader has complete control over the log**

#### How It Works

**Leader accepts log entries from clients and sends them to the followers**:
- **Client interaction**: Clients send requests to the leader
- **Log replication**: Leader replicates log entries to all followers
- **Unidirectional flow**: All log entries flow from leader to followers

**Once a majority of followers have received the log entries, the leader tells them to apply the entries (and the entries are now considered committed)**:
- **Majority requirement**: Only needs majority of followers to acknowledge
- **Commitment**: Once majority acknowledges, entry is committed
- **Durability**: Committed entries are guaranteed to be durable

#### Why This Reduces Complexity

**Log entries only flow from the leader to the followers!**

**This reduces complexity. Viewstamped Replication, for example, allows entries to flow in both directions, and the protocol as a whole is more complicated**:
- **Simplified protocol**: Only one direction of log flow
- **Clear responsibility**: Leader has complete control
- **Easier reasoning**: Simpler to understand and implement

### Raft Basics

#### Server States

**A server can be in one of three states**:
- **Leader**: Receives log entries from clients and sends them to followers
- **Follower**: Receives log entries from leader and applies them to state machine
- **Candidate**: Server trying to become the leader during an election

#### Terms

**Time is split into individual terms of arbitrary length**:
- **Each term starts with an election**: New term begins when election is triggered
- **Two outcomes**:
  1. **A leader is elected**: Election succeeds, leader takes control
  2. **There is a split vote**: Term ends, new election is conducted

**Each term has an unique identifying integer, and the integer increases with each successive term**:
- **Term numbers**: Each server stores the current term
- **Used to determine when a leader and/or followers are out-of-date**: Higher term numbers indicate more recent state

#### The Fail-Stop Assumption

**Again, note that the fail-stop failure assumption is important here!**

**If servers could lie about their term numbers, the protocol would fall apart**:
- **Trust requirement**: Servers must be honest about their term numbers
- **Byzantine failures**: Would require more complex protocols
- **Simplified model**: Makes the algorithm much easier to understand and implement

### RPC Types

**Two types of RPCs (well, 4 types in total, if you also count the responses)**:
- **RequestVote**: Used during leader elections
- **AppendEntries**: Append entries to the logs of followers and also used as a heartbeat mechanism

**AppendEntries serves dual purposes**:
- **Log replication**: Sends new log entries to followers
- **Heartbeat**: Keeps followers alive and detects leader failures
### Raft's Three Components

**Raft decomposes the protocol into three (mostly) independent parts**:

#### 1. Leader Election
How to choose a leader when the current leader fails

#### 2. Log Replication
How the leader replicates log entries to followers

#### 3. Safety
How to ensure the algorithm maintains correctness properties

**See "State Machine Safety" in Figure 3**

**"if a server has applied a log entry at a given index to its state machine, no other server will ever apply a different log entry for the same index"**

This is the fundamental safety property that Raft must maintain. It ensures that all servers apply the same sequence of commands to their state machines, which is essential for maintaining consistency.
## Leader Election

### Triggering an Election

**If a follower does not receive heartbeats for a duration of election timeout, it triggers an election**:
- **Heartbeat mechanism**: Leader sends heartbeats to followers
- **Timeout detection**: If no heartbeat received, leader may have failed
- **Election trigger**: Follower starts election process

### The Election Process

**Increments the term and becomes a candidate**:
- **Term increment**: New term number indicates new election
- **State change**: Follower becomes candidate
- **Election start**: Candidate begins election process

**Votes for itself and uses RequestVote RPCs to ask other servers to vote for it**:
- **Self-vote**: Candidate votes for itself
- **Request votes**: Sends RequestVote RPCs to all other servers
- **Vote collection**: Collects votes from other servers

**Other servers vote based on a FCFS policy**:
- **First-come-first-served**: First candidate to ask gets the vote
- **Term-based**: Only vote for candidates with higher or equal term numbers
- **Log-based**: Only vote for candidates with up-to-date logs

### Election Outcomes

**Three outcomes**:

#### (1) Wins Election
- **Majority votes**: Receives votes from majority of servers
- **Becomes leader**: Candidate becomes the new leader
- **Takes control**: Starts sending heartbeats and log entries

#### (2) Another Server Wins Election
**Case (2) is detected if a candidate receives an AppendEntries RPC from a "leader" whose term number is at least as great as the candidate's**:
- **Leader detection**: Receives heartbeat from current leader
- **Term comparison**: Leader's term is at least as high as candidate's
- **State change**: The candidate becomes a follower

#### (3) No Winner After Some Amount of Time
**Case (3) is triggered after a timeout is exceeded**:
- **Split vote**: No candidate receives majority of votes
- **Timeout**: Election timeout expires without winner
- **Another vote is necessary**: New election must be conducted

**Generally happens due to a split vote**:
- **Multiple candidates**: Several servers become candidates simultaneously
- **Vote splitting**: Votes are split among multiple candidates
- **No majority**: No candidate receives majority of votes

**Must have randomized election timeouts or split votes could continue indefinitely (or at least for a very long time)**:
- **Randomization**: Prevents repeated split votes
- **Convergence**: Ensures elections eventually succeed
- **Performance**: Reduces time to elect a leader

### Election Timeout Configuration

**Election timeouts chosen randomly between 150ms-300ms**:
- **Random range**: Each server chooses random timeout in this range
- **Split vote prevention**: Different timeouts prevent simultaneous elections
- **Reasonable duration**: Long enough to avoid unnecessary elections, short enough for quick recovery

**Authors found that even though the election timeout introduces randomness to the protocol, this approach was the easiest to understand of all the approaches they considered**:
- **Simplicity**: Random timeouts are easy to understand and implement
- **Effectiveness**: Successfully prevents split votes
- **Trade-off**: Acceptable complexity for significant benefit
## Log Replication

### The Replication Process

**Clients send log entries to the leader, which replicates them among the followers**:
- **Client interaction**: Clients send commands to the leader
- **Leader responsibility**: Leader replicates entries to all followers
- **Unidirectional flow**: All replication flows from leader to followers

### Log Entry Structure

**Log entry contains**:
1. **State machine command**: The actual command to be executed
2. **Current leader's term number**: Term when entry was created
3. **Index of position in log**: Position in the log sequence

### Commitment and Safety

**Leader distributes log entries to followers, but this does not mean that it is safe for the log entries to be applied to the replicated state machine yet**:
- **Distribution vs. commitment**: Replication doesn't imply safety
- **Safety requirement**: Entries must be committed before application
- **Consistency**: All servers must agree on committed entries

**Only safe once a log entry has been committed**:

**"Raft guarantees that committed entries are durable and will eventually be executed by all of the available state machines"**:
- **Durability**: Committed entries are guaranteed to survive failures
- **Execution**: All available state machines will eventually execute committed entries
- **Consistency**: All state machines will execute the same sequence of commands

**Entry is committed once leader that created entry has replicated it across majority of servers**:
- **Majority requirement**: Only needs majority of servers to acknowledge
- **Leader requirement**: Only entries created by current leader can be committed
- **Safety guarantee**: Majority ensures entry will survive leader failures

**If entry was not created by leader, then the entry is still not considered committed even if it is replicated on a majority of the cluster**:
- **Leader requirement**: Only current leader's entries can be committed
- **Safety**: Prevents committed entries from being overwritten
- **See Figure 8 for what can go wrong**: Shows why this restriction is necessary

### Commitment Propagation

**Once an entry is committed, it is assumed that all preceding entries are committed, too**:
- **Prefix property**: Committed entries form a prefix of the log
- **Implicit commitment**: Earlier entries are automatically committed
- **Consistency**: All servers agree on committed prefix

**"The leader keeps track of the highest index it knows to be committed, and it includes that index in future AppendEntries RPCs (including heartbeats) so that the other servers eventually find out."**:
- **Commitment tracking**: Leader maintains highest committed index
- **Propagation**: Includes committed index in all AppendEntries RPCs
- **Eventual consistency**: All servers eventually learn about committed entries

**"Once a follower learns that a log entry is committed, it applies the entry to its local state machine (in log order)."**:
- **Application**: Followers apply committed entries to state machine
- **Order preservation**: Entries are applied in log order
- **Consistency**: All state machines execute same sequence

### Log Matching Property

**"If two entries in different logs have the same index and term, then they store the same command."**:
- **Uniqueness**: Same index and term implies same command
- **Consistency**: All servers agree on entries with same index and term

**"If two entries in different logs have the same index and term, then the logs are identical in all preceding entries."**:
- **Prefix consistency**: Logs are identical up to matching entries
- **Induction**: Consistency propagates backwards through log
- **Safety**: Ensures all servers have consistent log prefixes

### Log Consistency Maintenance

**Leader includes index and term of previous entry in log when sending AppendEntries RPC**:
- **Consistency check**: Followers can verify log consistency
- **Previous entry**: Includes information about preceding entry
- **Verification**: Followers can check if logs match

**Followers use this to check log consistency**:
- **Reject the RPC if the index and term do not match the current last entry in their log**:
  - **Mismatch detection**: Identifies inconsistent logs
  - **Rejection**: Prevents inconsistent entries from being applied
  - **Recovery**: Triggers log repair process

**Logs are generally consistent, but they can become inconsistent when a leader crashes**:
- **See Figure 7**: Shows how logs can become inconsistent
- **Leader failure**: Crashes can leave logs in inconsistent state
- **Recovery**: New leader must repair inconsistent logs

### Log Repair

**Raft handles this by forcing followers to have same log as leader**:
- **Overwrite**: Conflicting entries in followers' logs are overwritten
- **Consistency**: All followers end up with same log as leader
- **Safety**: Ensures all servers have consistent logs

**Leader maintains nextIndex for each follower (index of next entry that will be sent to follower)**:
- **Tracking**: Leader tracks what to send to each follower
- **Individual state**: Each follower can be at different position
- **Efficiency**: Only sends entries that follower needs

**If AppendEntries RPC fails, leader decrements nextIndex for that follower and tries again**:
- **Retry logic**: Continues until AppendEntries succeeds
- **Backtracking**: Moves backwards through log to find consistent point
- **Eventually, nextIndex will be low enough that this will succeed**: Guaranteed to find consistent point

**This seems kind of slow… can you think of a simple way to make it faster?**:
- **Optimization opportunity**: Linear search through log is inefficient
- **Binary search**: Could use binary search to find consistent point faster
- **Trade-off**: Simplicity vs. performance
## Safety

### The Leader Election Restriction

**Must have a restriction on which servers can be elected leader**:
- **Safety requirement**: Not all servers can be elected leader
- **Consistency**: Must ensure new leader has all committed entries
- **Prevention**: Prevents committed entries from being overwritten

**Otherwise, an out-of-date follower could be elected leader and overwrite committed entries with its own**:
- **Danger**: Out-of-date leader could overwrite committed entries
- **Consistency violation**: Could cause state machines to execute different sequences of commands across different followers
- **Real-world consequences**: Committed entries could have already had real-world consequences

### The Election Restriction

**Election restriction: Raft "guarantees that all the committed entries from previous terms are present on each new leader from the moment of its election, without the need to transfer those entries to the leader"**:
- **Automatic guarantee**: New leader automatically has all committed entries
- **No transfer needed**: No need to copy entries to new leader
- **Safety**: Ensures committed entries are never lost

**Leader is only elected if it has all committed entries**:
- **Requirement**: Candidate must have all committed entries to be elected
- **Vote condition**: Servers only vote for candidates with all committed entries
- **Safety**: Prevents out-of-date leaders from being elected

**Leader must have vote from majority to win election**:
- **Majority requirement**: Needs majority of servers to vote for it
- **Consensus**: Majority ensures leader has all committed entries
- **Safety**: Majority guarantee ensures consistency

### The Up-to-Date Requirement

**So at least one of the servers in the majority must have all committed log entries**:
- **Majority property**: At least one server in majority has all committed entries
- **Vote requirement**: This server will only vote for candidates with all committed entries
- **Safety**: Ensures new leader has all committed entries

**This server will only vote for a candidate whose log is at least as "up-to-date" as its own**:
- **Vote condition**: Only votes for up-to-date candidates
- **Consistency**: Ensures new leader has all committed entries
- **Safety**: Prevents out-of-date leaders from being elected

**Up-to-date has two requirements**:

#### (1) Term-Based Comparison
**If both logs have last entries with different terms, the log with the later term is more up-to-date**:
- **Term priority**: Later terms indicate more recent state
- **Consistency**: Later terms are more likely to have committed entries
- **Safety**: Ensures new leader has most recent state

#### (2) Length-Based Comparison
**If the logs have last entries with the same term, the longer log is more up-to-date**:
- **Length priority**: Longer logs have more entries
- **Consistency**: More entries indicate more complete state
- **Safety**: Ensures new leader has all entries from same term

### Safety Properties

**Safety argument proof**:
- **Mathematical proof**: Raft's safety properties are mathematically proven
- **Correctness**: Algorithm is guaranteed to maintain consistency
- **Reliability**: System will never produce incorrect results

### Handling Crashes

**Follower and candidate crashes**:
- **Idempotent RPCs**: It is ok to keep resending RequestVote and AppendEntries RPCs to crashed followers and candidates because Raft RPCs are idempotent
- **No harm**: Resending RPCs doesn't cause problems
- **Recovery**: Crashed servers can recover and rejoin

### Safety and Timing

**Timing does not affect the correctness of the protocol!**:
- **Correctness**: Algorithm remains correct regardless of timing
- **Robustness**: Works even with timing variations
- **Reliability**: Timing issues don't cause incorrect results

**But it can prevent the protocol from making reasonable progress**:
- **Performance impact**: Timing issues can slow down the algorithm
- **Progress**: System may not make progress if timing is too bad
- **Practical considerations**: Real systems need reasonable timing

**"Leader election is the aspect of Raft where timing is most critical"**:
- **Election sensitivity**: Elections are most affected by timing
- **Critical path**: Election timing affects overall system performance
- **Optimization**: Election timing is most important to optimize

### Timing Requirements

**Raft can elect and maintain a leader if the system meets the timing requirement**:
- **broadcastTime << electionTimeout << MTBF**
- **broadcastTime**: Average time it takes to send RPCs in parallel to all servers and receive a response
- **electionTimeout**: Election timeout duration
- **MTBF**: Mean Time Between Failures

**broadcastTime should be order of magnitude less than electionTimeout so that elections are not needlessly triggered due to slow heartbeats**:
- **Heartbeat timing**: Heartbeats must be faster than election timeout
- **Election prevention**: Prevents unnecessary elections due to slow heartbeats
- **Performance**: Ensures system can maintain stable leadership

**electionTimeout should be several orders of magnitude less than MTBF to ensure that a leader can be elected and make steady progress before failing (otherwise the system could not make reasonable progress)**:
- **Failure rate**: Elections must be much faster than failures
- **Progress**: System must make progress between failures
- **Stability**: Ensures system can maintain stable operation
## Cluster Membership Changes

### The Need for Dynamic Membership

**Cluster membership will likely change over time**:
- **Scaling**: Need to add servers as system grows
- **Maintenance**: Need to remove servers for maintenance
- **Failures**: Need to replace failed servers
- **Dynamic requirements**: Real systems need dynamic membership

**Want to adjust it without taking the entire system down**:
- **Availability**: System must remain available during membership changes
- **Online reconfiguration**: Changes must happen without downtime
- **Seamless operation**: Clients should not notice membership changes

**Also don't want to require manual intervention for this since it will be a source of error (e.g., Instagram went down a couple years ago due to this)**:
- **Automation**: Manual intervention is error-prone
- **Reliability**: Automated membership changes are more reliable
- **Real-world example**: Instagram outage due to manual configuration error

### The Split-Brain Problem

**When changing a configuration, need to ensure that there cannot be two leaders elected at once**:
- **Safety requirement**: Must prevent split-brain scenarios
- **Consistency**: Only one leader can exist at a time
- **Prevention**: Must ensure no two leaders are elected simultaneously

**Can't update all servers at once, so servers could be split into old vs. new configurations and separate majorities could elect different leaders**:
- **Gradual update**: Servers are updated one at a time
- **Configuration split**: Some servers have old config, others have new config
- **Split-brain risk**: Each configuration could elect its own leader
- **See Figure 10**: Shows the split-brain problem

### Joint Consensus Solution

**Raft uses joint consensus to handle this**:
- **Two-phase approach**: Raft first switches to the joint consensus configuration, then switches to the new configuration
- **Joint consensus combines old and new configurations**: Combines both configurations temporarily
- **Safety**: Prevents split-brain scenarios during transitions

#### How Joint Consensus Works

**Log entries duplicated to all servers in both configurations**:
- **Replication**: All log entries are sent to all servers in both configurations
- **Consistency**: All servers have same log entries
- **Safety**: Ensures consistency across both configurations

**Any server from either configuration can be a leader**:
- **Flexibility**: Leader can come from either old or new configuration
- **Availability**: System remains available during transition
- **Efficiency**: No need to wait for specific servers

**Agreement for elections and log entry commitment requires separate majorities from both the old and new configurations**:
- **Dual majority**: Need majority from both configurations
- **Safety**: Prevents split-brain scenarios
- **Consistency**: Ensures both configurations agree

### The Configuration Change Process

**When you want to change configuration from Cold to Cnew, leader creates Cold,new, adds it to its own log, and then replicates it across all followers**:
- **Joint configuration**: Creates configuration that includes both old and new
- **Log entry**: Adds joint configuration to log
- **Replication**: Replicates joint configuration to all followers

**Once a configuration is present in a log, it is used even if it hasn't been committed**:
- **Immediate use**: Configuration is used as soon as it's in the log
- **Safety**: Ensures consistent behavior across servers
- **Efficiency**: No need to wait for commitment

**If the leader crashes before Cold,new is committed, a new leader is elected either from Cold or Cold,new**:
- **Recovery**: System can recover from leader crashes during transition
- **Flexibility**: New leader can come from either configuration
- **Safety**: System remains consistent

#### Important Constraints

**Importantly, Cnew cannot make decisions on its own yet!**:
- **Safety**: New configuration cannot operate independently
- **Prevention**: Prevents split-brain scenarios
- **Requirement**: Must work with old configuration

**It's ok for decisions to be made by Cold or Cold,new because the majorities from both configurations have the same servers**:
- **Overlap**: Both configurations share servers
- **Consistency**: Shared servers ensure consistency
- **Safety**: No split-brain possible

### The Second Phase

**Once Cold,new is committed, the leader replicates Cnew and waits for it to commit**:
- **Second phase**: Switch to new configuration only
- **Replication**: Send new configuration to all servers
- **Commitment**: Wait for new configuration to be committed

**While Cold,new is committed, neither Cold nor Cnew can make decisions without the other**:
- **Joint operation**: Both configurations must agree
- **Safety**: Prevents split-brain scenarios
- **Consistency**: Ensures consistent behavior

**Once Cnew is committed, the servers only in Cold can be taken down**:
- **Completion**: Transition is complete
- **Cleanup**: Old servers can be removed
- **Efficiency**: System now operates with new configuration only

### Three Issues Addressed

#### (1) New Server Catching Up

**New servers can take a while to get all of the log entries, and we do not want commits to be impacted in the meantime**:
- **Catch-up time**: New servers need time to receive all log entries
- **Performance**: Don't want to slow down commits during catch-up
- **Solution**: New servers initially join cluster as non-voting members. Once they have all the log entries, the cluster configuration change starts

#### (2) Leader Not in New Configuration

**Cluster leader may not be part of new configuration**:
- **Leader exclusion**: Leader might not be in new configuration
- **Transition**: Leader will step down once Cnew has been committed because now the new cluster can operate independently
- **Management**: Thus, there may be some time where the leader is managing a cluster that it is not a part of

**Replicates log entries and commits when majority has entry, but the leader does not count itself as part of the majority**:
- **Self-exclusion**: Leader doesn't count itself in majority
- **Safety**: Ensures new configuration can make decisions
- **Consistency**: Maintains consistency during transition

#### (3) Removed Servers Causing Disruptions

**Removed servers no longer receive heartbeats, so they can trigger elections that disrupt the new cluster**:
- **Heartbeat loss**: Removed servers don't receive heartbeats
- **Election triggers**: They may start elections
- **Disruption**: These elections could disrupt the new cluster

**Servers disregard RequestVote RPCs if received within the election timeout of hearing from a leader**:
- **Election prevention**: Servers ignore vote requests if they recently heard from a leader
- **Stability**: Prevents unnecessary elections
- **Efficiency**: Maintains stable leadership

**Does not affect normal elections (i.e., common case) because each server waits for the election timeout before starting an election**:
- **Normal operation**: Regular elections still work normally
- **Timeout**: Servers wait for timeout before starting elections
- **Efficiency**: No impact on normal operation
## Evaluation

### Raft's Success

Raft has been remarkably successful in achieving its primary goal of understandability while maintaining the safety and performance properties required for practical distributed systems. The algorithm has been widely adopted in production systems and has become the de facto standard for consensus in many distributed systems.

### Key Achievements

#### 1. Understandability
- **Clear structure**: The three-component decomposition makes the algorithm easy to understand
- **Intuitive design**: Design decisions make intuitive sense
- **Good documentation**: The paper provides clear explanations and examples
- **Wide adoption**: Many systems have successfully implemented Raft

#### 2. Safety and Correctness
- **Mathematically proven**: Safety properties are formally proven
- **Practical implementation**: Algorithm has been successfully implemented in many systems
- **Real-world validation**: Used in production systems with good results

#### 3. Performance
- **Efficient common case**: Fast operation when leader is stable
- **Reasonable recovery**: Quick recovery from leader failures
- **Scalable**: Performance scales well with system size

### Impact on Distributed Systems

Raft has had a significant impact on the distributed systems community:

- **Replaced Paxos**: Many systems now use Raft instead of Paxos
- **Easier implementation**: Reduced the barrier to implementing consensus
- **Better systems**: Enabled more reliable distributed systems
- **Educational value**: Made consensus algorithms more accessible

### Future Directions

While Raft has been successful, there are still areas for improvement:

- **Performance optimizations**: Further optimizations for specific workloads
- **Extended features**: Additional features for specific use cases
- **Integration**: Better integration with other distributed systems components
- **Monitoring**: Better tools for monitoring and debugging Raft systems

### Conclusion

Raft represents a significant achievement in making consensus algorithms more understandable and practical. By prioritizing understandability without sacrificing safety or performance, Raft has made it possible for more developers to build reliable distributed systems. The algorithm's success demonstrates that good design can make complex algorithms accessible while maintaining their essential properties.