# Microsecond Consensus: Achieving High Performance in Distributed Systems

## Introduction: The Microsecond Challenge

### The Scale of Modern Applications

**Many modern applications operate on the microsecond-scale**

To understand the significance of this statement, let's review the time scales we're working with:

- **1 second** = 1,000 milliseconds
- **1 second** = 1,000,000 microseconds  
- **1 second** = 1,000,000,000 nanoseconds

When we talk about microsecond-scale operations, we're dealing with operations that complete in millionths of a second. This is an incredibly demanding performance requirement that pushes the boundaries of what's possible with current hardware and software.

### Real-World Applications

#### Finance Applications (Especially Trading Systems)

**High-frequency trading systems** operate in this microsecond realm where every nanosecond of latency can mean the difference between profit and loss. In financial markets, being first to execute a trade can provide significant competitive advantages, making microsecond-level performance not just desirable but essential for survival.

#### Embedded Computing (Control Systems)

**Real-time control systems** in industrial automation, automotive systems, and robotics require microsecond-level response times to maintain safety and precision. A delay of even a few microseconds in a control loop could cause system instability or safety violations.

#### Microservices (Key-Value Stores)

**Modern microservice architectures** often rely on high-performance key-value stores as their foundation. For example, a GET request in RocksDB for a small key-value pair takes something like 4 µs. When your entire application's performance depends on these microsecond-level operations, any additional overhead becomes significant.

### The Infrastructure Challenge

**Discuss concerns around network stacks, kernel scheduling, memory management (hugepages), etc.**

Achieving microsecond-level performance requires careful attention to every layer of the system:

- **Network stacks**: Traditional TCP/IP stacks add significant latency
- **Kernel scheduling**: Context switches and system calls introduce delays
- **Memory management**: Page faults and memory allocation can cause unpredictable delays
- **Hardware optimization**: CPU cache behavior, NUMA effects, and memory bandwidth all matter

### The Replication Challenge

**Want to replicate these services so they have high availability**

**e.g., in a microservice application, some of the services are critical to the functioning of the entire system (such as a key-value store)**

The challenge is that while we can achieve microsecond-level performance for single-node operations, adding replication for high availability typically introduces significant overhead. This creates a fundamental tension between performance and reliability.
## State Machine Replication (SMR): The Performance Problem

### What is State Machine Replication?

**"Replicas execute requests in the same total order determined by a consensus protocol"**

State Machine Replication is a fundamental technique for building fault-tolerant distributed systems. The idea is elegant: if multiple replicas execute the same sequence of requests in the same order, they will all end up in the same state. The consensus protocol ensures that all replicas agree on the same total order of requests.

### The Performance Problem

#### Traditional SMR Systems Are Too Slow

**Traditional SMR systems are way too slow (add hundreds of µs of overhead)**

**This happens even on a "fast" network, so the network itself is not responsible for the overhead**

The problem isn't the network latency - even on high-speed networks with sub-microsecond latency, traditional SMR systems add hundreds of microseconds of overhead. This overhead comes from:

- **Protocol complexity**: Multiple round-trips for consensus
- **Software overhead**: Kernel crossings, context switches, memory copies
- **Synchronization**: Locks, barriers, and coordination between threads
- **Logging and persistence**: Writing to disk or persistent memory

#### Even Optimized Systems Add Significant Overhead

**Even optimized SMR systems (Hermes, DARE, Hovercraft) add a few µs of overhead**

**If a RocksDB GET request takes 4 µs, you are cutting the total possible throughput roughly in half if you then add a few µs of consensus overhead**

This is the crux of the problem. When your base operation takes 4 microseconds, adding even 2-3 microseconds of consensus overhead means you're roughly halving your throughput. In high-performance systems where every microsecond matters, this is unacceptable.

#### Recovery Time is Even Worse

**When a failure occurs, it takes tens of milliseconds to recover**

**These systems rely on timeouts to detect failure, and the timeout needs to be sufficiently large to not get triggered by natural delays in the network**

The recovery problem is even more severe. Traditional consensus systems use timeouts to detect failures, but these timeouts must be large enough to avoid false positives from natural network delays. This means:

- **Slow failure detection**: Takes tens of milliseconds to detect a failure
- **Conservative timeouts**: Must account for worst-case network conditions
- **Long recovery**: Even after detection, recovery takes additional time

### The Fundamental Challenge

**Thus, it's not clear how to improve this overhead with existing techniques… new techniques must be devised**

The traditional approaches to consensus have fundamental limitations that make them unsuitable for microsecond-scale applications. We need entirely new techniques that can achieve consensus with minimal overhead while maintaining safety and liveness properties.
## Mu: A New Approach to Microsecond Consensus

### Introducing Mu

**New SMR system**

Mu represents a breakthrough in consensus algorithm design, specifically targeting microsecond-scale performance requirements. It's not just an incremental improvement over existing systems - it's a fundamentally new approach that rethinks how consensus can be achieved at the hardware level.

### Performance Achievements

**1.3 µs to replicate a small request**

**99th percentile overhead of 1.6 µs**

These numbers are remarkable. Mu achieves consensus overhead of just 1.3 microseconds in the common case, with 99% of operations completing within 1.6 microseconds. This represents a dramatic improvement over traditional systems that add hundreds of microseconds of overhead.

### The RDMA Advantage

**Mu requires only one-sided RDMA communication in the common case**

**Other systems require multiple roundtrips or two-sided communication**

**Not possible to further improve overhead in Mu without improving hardware**

The key insight is that Mu leverages RDMA (Remote Direct Memory Access) to achieve one-sided communication. This means:

- **No CPU involvement**: The leader can write directly to follower memory without involving the follower's CPU
- **Single round-trip**: No need for multiple message exchanges
- **Hardware optimization**: The network interface card (NIC) handles the communication directly

This is a fundamental shift from traditional consensus protocols that require multiple round-trips and CPU involvement on both sides.

### Leader Uniqueness Through RDMA Permissions

**Mu uses RDMA write permissions to ensure that only one leader at a time can write to a replica (multiple leaders cannot race)**

This is a clever use of RDMA's permission system. By granting write permissions to only one leader at a time, Mu ensures that:

- **No race conditions**: Multiple leaders cannot simultaneously write to the same replica
- **Hardware enforcement**: The NIC hardware enforces the permission, not software
- **Atomic operations**: Write permissions are granted and revoked atomically

### Novel Mechanisms in Mu

**New mechanisms in Mu**:

#### 1. RDMA Write Permissions
**Use of RDMA write permissions (one leader at a time can write to a replica)**
- **Hardware-level enforcement**: The NIC ensures only one leader can write
- **Atomic permission changes**: Permissions are changed atomically
- **No software coordination**: No need for software locks or coordination

#### 2. Leader Change Mechanism
**Mechanism to change leader**
- **Permission revocation**: Old leader's write permissions are revoked
- **Permission grant**: New leader is granted write permissions
- **State reconstruction**: New leader reconstructs any partial work from the old leader

#### 3. Garbage Collection
**Garbage collection mechanism for logs**
- **Log recycling**: Old log entries are recycled to prevent memory exhaustion
- **Safe reclamation**: Ensures committed entries are never reclaimed
- **Efficient management**: Minimizes memory overhead

### Fast Failure Detection and Recovery

**Failover overhead in Mu is 873 µs (99% latency of 945 µs)**

**Order of magnitude improvement over prior systems**

This is another remarkable achievement. While traditional systems take tens of milliseconds to recover from failures, Mu achieves recovery in under a millisecond.

#### The Pull-Score Mechanism

**Replicas read a heartbeat counter from the leader and calculate a badness score**

**When the score goes above a threshold, the replica considers the leader to be failed**

**Can use aggressive timeouts because reads are delayed by the network rather than heartbeats**

This is a key insight. Instead of using traditional heartbeats, Mu uses a pull-based approach:

- **Replicas pull**: Followers actively read a counter from the leader
- **Network delay**: The read operation itself is delayed by the network
- **Aggressive timeouts**: Can use much shorter timeouts because the read operation provides the timing information

**So the replica just needs to wait for its own read to complete and then it can update the score**

**In the common case, this will be fast**

The replica doesn't need to wait for a timeout - it just needs to wait for its own read operation to complete. If the leader is alive, the read will complete quickly. If the leader has failed, the read will timeout or fail.

#### Why This is Better Than Heartbeats

**If a heartbeat were used, the replica would need to wait for the timeout because for all it knows, the leader is alive and has sent the heartbeat, but the heartbeat is delayed by the network**

**In reality, the leader may have failed and the heartbeat was never sent, but the replica cannot assume this until the timeout is exceeded**

With traditional heartbeats, there's ambiguity about whether a missing heartbeat means the leader failed or the heartbeat was just delayed. With the pull-based approach, the replica knows immediately when its read operation completes or fails.

### Hardware Limitations

**Most of the failover overhead is from changing RDMA write permissions**

**The NIC hardware needs to be improved to make this faster**

Even with Mu's optimizations, the hardware itself becomes the bottleneck. Changing RDMA permissions requires hardware-level operations that take time. This suggests that future improvements will require hardware changes, not just software optimizations.
### Limitations of Mu

#### RDMA Requirement

**Requires RDMA**

**May work on local network, but will not work across the Internet**

This is a significant limitation. RDMA requires specialized network hardware and protocols that are typically only available in local area networks (LANs). This means:

- **Limited deployment**: Mu can only be used in controlled environments with RDMA-capable networks
- **No Internet deployment**: Cannot be used across the Internet or in cloud environments without RDMA
- **Infrastructure dependency**: Requires specific network hardware and configuration

#### Persistence Challenge

**In-memory system–does not write data to persistent storage**

**People are working RDMA extensions that will support writes to persistent memory**

This is another significant limitation. Mu currently only works with in-memory data, which means:

- **Data loss on failure**: If a machine fails, all data in memory is lost
- **No durability**: Cannot provide durability guarantees that many applications require
- **Limited use cases**: Only suitable for applications that can tolerate data loss

#### Future Possibilities

**I am not sure what the performance will look like for that, but maybe you could envision some sort of theoretical setup where the NIC and the persistent memory are powered by a battery, which will allow the writes to be persisted even if the machine fails?**

**Thus, the NIC can respond immediately to the leader even if the write has not yet been persisted… this is important for performance**

This is an interesting theoretical possibility. The idea would be:

- **Battery-backed NIC**: The network interface card and persistent memory are powered by a battery
- **Immediate response**: The NIC can respond to the leader immediately, even before the write is fully persisted
- **Crash safety**: If the machine fails, the battery ensures the write is completed to persistent memory
- **Performance**: Maintains the performance benefits of immediate response while providing durability

This would be a significant hardware innovation that could make Mu suitable for applications requiring durability guarantees.
## Background: State Machine Replication and RDMA

### State Machine Replication (SMR)

**Replicates a service across replicas**

**"Provides strong consistency in the form of linearizability"**

State Machine Replication is a fundamental technique for building fault-tolerant distributed systems. The core idea is to maintain multiple copies (replicas) of a service and ensure they all execute the same sequence of operations in the same order.

#### How SMR Works

**Consensus protocol ensures that the logs (up through the last committed entry) are identical across the majority of replicas**

**Each log entry holds a request**

**The replica applies the requests in order to the state machine**

The process works as follows:

1. **Request logging**: Each request is logged as an entry in a replicated log
2. **Consensus**: A consensus protocol ensures all replicas agree on the same log entries
3. **State application**: Each replica applies the logged requests to its local state machine
4. **Consistency**: Since all replicas apply the same requests in the same order, they maintain identical state

#### Deterministic Execution

**We are assuming that each request is handled deterministically, so that all state machines that handle the same requests in the same order have the same state**

This is a crucial assumption. For SMR to work correctly, the state machine must be deterministic:

- **Same input, same output**: Given the same sequence of requests, all replicas must produce the same state
- **No randomness**: The state machine cannot use random numbers or other non-deterministic operations
- **Deterministic ordering**: The order of request execution must be deterministic

### Mu's SMR Model

**Mu assumes there is a single leader**

**Clients send requests to the leader**

**The leader replicates the requests across the majority of the cluster**

**The leader responds to the clients**

Mu uses a leader-based approach to SMR:

1. **Single leader**: One replica acts as the leader at any given time
2. **Client interaction**: Clients send requests only to the leader
3. **Replication**: The leader replicates requests to followers
4. **Response**: The leader responds to clients after successful replication

This is similar to Raft's approach but optimized for microsecond-scale performance.

#### Failure Model

**Crash-failure model**

**When a server fails, it just crashes. It does not emit additional, arbitrary values.**

Mu assumes a crash-failure model, which means:

- **Clean failures**: Servers fail by stopping, not by sending incorrect messages
- **No Byzantine behavior**: Failed servers don't send arbitrary or malicious messages
- **Simplified protocol**: This assumption allows for simpler and more efficient protocols

This is a reasonable assumption for many practical systems where hardware failures are more common than malicious behavior.
### RDMA (Remote Direct Memory Access)

#### Understanding DMA First

**DMA (not RDMA) already exists in computers**

**The DMA engine is a hardware accelerator that transfers data between memory (e.g., DRAM) and an I/O device**

**The CPU could do this itself, but it would need to coordinate everything and move chunks of data (e.g., 64 bit values, or cache lines, etc.) one at a time**

**This requires significant involvement from the CPU, which now cannot spend its cycles running other workloads**

**Thus, given that data transfer is a common operation, the DMA engine was implemented to handle this**

DMA (Direct Memory Access) is a fundamental hardware feature that allows I/O devices to transfer data directly to and from memory without involving the CPU. This is important because:

- **CPU efficiency**: The CPU doesn't need to handle every byte of data transfer
- **Parallelism**: Data transfer can happen while the CPU does other work
- **Performance**: Hardware-accelerated transfers are much faster than CPU-based transfers

**DMA only works within a single machine!**

**The I/O device and the memory must be in the same machine**

This is a key limitation of DMA - it only works within a single machine.

#### RDMA: Extending DMA Across Machines

**RDMA allows one machine to access the memory of another machine, without involvement from the remote machine's CPU. This is the point.**

RDMA extends the DMA concept across the network, allowing one machine to directly read from or write to the memory of another machine. This is revolutionary because:

- **No CPU involvement**: The remote machine's CPU doesn't need to be involved in the data transfer
- **Hardware acceleration**: The network interface card (NIC) handles the transfer directly
- **Low latency**: Eliminates software overhead and context switches
- **High throughput**: Can achieve near-hardware limits for data transfer

#### RDMA Operations

**RDMA operations**:
- **Send/Receive**: Traditional message passing
- **Read/Write**: Direct memory access operations
- **Atomics**: Atomic memory operations

**Used for two-sided communication**

**Used for one-sided communication**

RDMA supports both traditional two-sided communication (where both sides are involved) and one-sided communication (where only one side initiates the operation). Mu primarily uses one-sided operations for maximum performance.

#### RDMA Transports

**RDMA transports**

**Mu uses Reliable Connection (RC)**

**In-order, reliable delivery**

RDMA provides different transport types with different guarantees. Mu uses Reliable Connection (RC) which provides:

- **Reliability**: Messages are guaranteed to be delivered
- **Ordering**: Messages are delivered in order
- **Error detection**: Network errors are detected and handled

#### RDMA Queues and Memory Management

**RDMA queues**

**Queue Pairs (QPs) are RDMA connection endpoints**

**Each QP has its own Completion Queue (CQ)**

**Work Requests (WRs) are posted to a QP**

**RDMA hardware performs the work and then write a Work Completion (WC) to the CQ**

The RDMA programming model uses a queue-based approach:

1. **Queue Pairs**: Each connection has a send queue and receive queue
2. **Work Requests**: Applications post work requests to queues
3. **Hardware execution**: RDMA hardware executes the requests
4. **Completion notification**: Hardware writes completion notifications to completion queues

**Applications register local virtual memory regions (MRs) with RDMA**

**QPs and MRs have different access modes**

**Examples:**
- **Read-only**
- **Read-write**
- **etc.**

**The access mode is set when the QP is initialized and when the MR is registered**

**Can be changed later (important for leader write permissions)**

**Can register same memory multiple times to get multiple MRs. Each MR can have a different access mode.**

**Thus, can give different permissions to each remote machine**

This is a key feature that Mu exploits. By registering the same memory region multiple times with different access modes, Mu can give different permissions to different remote machines. This is how Mu ensures that only the current leader can write to a replica's memory.
## Mu: The System Design

### High-Level Architecture

**Show Figure 1**

**A replica can be a leader or a follower**

**Each replica grants RDMA write permissions to its current leader**

**Leader uses RDMA write to write new request to log of each replica**

**When replica detects this, it applies the request to the application**

Mu's architecture is elegantly simple:

1. **Role assignment**: Each replica is either a leader or a follower
2. **Permission management**: Each replica grants RDMA write permissions to the current leader
3. **Direct replication**: The leader uses RDMA writes to directly write requests to follower logs
4. **Application execution**: Followers detect new log entries and apply them to the application

### Write Detection Mechanisms

**Paper does not state how the replica detects a write**

**Can use either interrupts or polling**

The paper doesn't specify how followers detect new log entries, but there are several possible approaches:

#### Interrupt-Based Detection

**Cool idea: Use RDMA to trigger an MSI-X on the replica**

**Explain what an MSI-X is and how it works**

MSI-X (Message Signaled Interrupts - Extended) is a hardware mechanism that allows devices to send interrupts to the CPU. The idea would be:

- **Hardware interrupt**: RDMA hardware triggers an interrupt when a write completes
- **Immediate notification**: The follower is immediately notified of new log entries
- **Low latency**: Interrupts provide the fastest possible notification
- **CPU overhead**: Interrupts require CPU involvement and context switching

#### Polling-Based Detection

**Cool idea: I read some other work that uses a compiler technique to manually inject manual checks into an application. Each check looks to see if a flag has been set to indicate new work.**

**After the first check, subsequent checks are fast because the cache line is already in the L1 cache.**

**Once the RDMA write occurs, there will be a cache miss, but you'd need to take the cache miss anyway. The important point is that nearly all of the checks hit the L1 cache (they are fast) up until the RDMA write occurs.**

This is a clever optimization technique:

- **Compiler injection**: The compiler automatically inserts checks for new work
- **Cache efficiency**: Most checks hit the L1 cache and are very fast
- **Cache miss on write**: When an RDMA write occurs, it causes a cache miss, but this is necessary anyway
- **Performance**: The overhead of polling is minimal because most checks are cache hits

#### Memory Coherence Requirements

**This requires the RDMA engine to be coherent with host memory**

**Will be easy to get this with CXL**

**Explain what CXL is**

CXL (Compute Express Link) is a new interconnect standard that provides:

- **Memory coherence**: Ensures that RDMA operations are coherent with host memory
- **Low latency**: Provides high-bandwidth, low-latency connections
- **Cache coherence**: Maintains cache coherence between different components
- **Future compatibility**: Will make RDMA coherence easier to achieve

### Trade-offs Between Detection Methods

**Go over tradeoffs**

The choice between interrupt-based and polling-based detection involves several trade-offs:

#### Interrupt-Based Advantages
- **Immediate notification**: Fastest possible detection of new work
- **CPU efficiency**: CPU can do other work until interrupted
- **Precise timing**: Exact timing of when new work arrives

#### Interrupt-Based Disadvantages
- **Context switching**: Interrupts require context switches
- **Cache pollution**: Interrupts can cause cache misses
- **Complexity**: Interrupt handling is more complex

#### Polling-Based Advantages
- **Simple implementation**: Easy to implement and debug
- **Cache efficiency**: Most checks are cache hits
- **Predictable performance**: No unexpected context switches

#### Polling-Based Disadvantages
- **CPU overhead**: CPU must continuously check for new work
- **Latency**: Slight delay between work arrival and detection
- **Power consumption**: Continuous polling consumes more power
### The Leader Failure Challenge

**Main challenge is handling leader failures**

**Can't have two leaders each updating replicas**

**This could happen if there are network delays**

The fundamental challenge in any consensus system is handling leader failures. In Mu, this is particularly critical because:

- **Multiple leaders**: If there are network delays, multiple replicas might think they are the leader
- **Race conditions**: Two leaders could simultaneously write to the same replica
- **Consistency violation**: This would violate the consistency guarantees of the system

### Mu's Solution: Pull-Based Failure Detection

**Solution:**

**Leader increments counter in local memory**

**Each replica does an RDMA read on this counter every so often**

**If a replica sees the same counter value across multiple RDMA reads, it assumes the leader has failed**

Mu uses a pull-based approach for failure detection:

1. **Counter mechanism**: The leader continuously increments a counter in its local memory
2. **Periodic reads**: Each replica periodically reads this counter using RDMA
3. **Failure detection**: If the counter value doesn't change across multiple reads, the replica assumes the leader has failed

#### Advantages of Pull-Based Detection

**When does an RDMA read time out?**

**Why not have the leader push counter updates to the readers? Readers can also pull the counter value themselves with RDMA reads. Could this be faster, especially if the leader piggybacks counter updates onto other RDMA writes to the replica logs?**

The pull-based approach has several advantages:

- **No additional messages**: No need for separate heartbeat messages
- **Piggybacking**: Counter updates can be piggybacked onto other RDMA writes
- **Efficiency**: Followers can read the counter when they need to, not when the leader sends it
- **Flexibility**: Each replica can adjust its read frequency based on its needs

#### Leader Change Process

**New leader revokes write permissions granted to old leader**

**Also reconstructs partial work performed by previous leader**

**Thus, two leaders cannot write to the same replica**

When a new leader is elected, it must:

1. **Permission revocation**: Revoke write permissions from the old leader
2. **Permission grant**: Grant write permissions to itself
3. **State reconstruction**: Reconstruct any partial work performed by the previous leader
4. **Consistency**: Ensure that only one leader can write to each replica

This process ensures that the system maintains consistency even during leader changes.
### Mu's Two-Plane Architecture

**Two major parts**

Mu is designed with a two-plane architecture that separates concerns and optimizes for different performance requirements:

#### Replication Plane

**Replication plane**

**Handles replication**

**Can be in either:**
- **Leader mode**
- **Follower mode**

**Three components:**

##### 1. Replicator (Leader Mode Only)
**Replicates requests by writing them to the followers' logs**

The replicator is responsible for:
- **Request replication**: Writing client requests to follower logs using RDMA
- **Consensus coordination**: Ensuring requests are replicated to a majority of followers
- **Response handling**: Responding to clients after successful replication

##### 2. Replayer (Follower Mode Only)
**Replays entries from local log**

The replayer is responsible for:
- **Log processing**: Reading entries from the local log
- **State machine execution**: Applying log entries to the local state machine
- **Consistency maintenance**: Ensuring the local state matches other replicas

##### 3. Logging
**Stores client requests to be applied to the state machine**

**Also keep a copy of remote logs (used by a new leader to complete partial work performed by older leaders)**

The logging component:
- **Local log management**: Stores client requests in the local log
- **Remote log tracking**: Maintains copies of remote logs for recovery
- **Recovery support**: Helps new leaders complete partial work from previous leaders

#### Background Plane

**Background plane**

**Monitors health of leader**

**Puts replication plane into either leader mode or follower mode**

**Handles permission changes**

**Two components:**

##### 1. Leader Election
**Detects failure of a leader**

**Selects another replica to become the leader**

The leader election component:
- **Failure detection**: Monitors the health of the current leader
- **Election coordination**: Selects a new leader when the current leader fails
- **Mode switching**: Transitions the replication plane between leader and follower modes

##### 2. Permission Management
**Handles RDMA write permissions**

**Maintains permissions array**

**A remote replica requests write permission on this machine by writing a 1 to a slot in the permissions array**

The permission management component:
- **Permission tracking**: Maintains an array of write permissions for each replica
- **Permission requests**: Handles requests for write permissions from remote replicas
- **Permission revocation**: Revokes permissions when leaders change

#### Thread and Queue Organization

**Each plane has its own threads and queues**

This separation provides several benefits:

- **Performance isolation**: Replication and background operations don't interfere with each other
- **Specialized optimization**: Each plane can be optimized for its specific workload
- **Fault isolation**: Problems in one plane don't necessarily affect the other
- **Scalability**: Each plane can be scaled independently
### RDMA Communication Architecture

**RDMA communication**

**On a given machine, there are two QPs for each remote replica**
- **One for the replication plane**
- **One for the background plane**

**All QPs for replication plane share same CQ**

**All QPs for management plane share same CQ**

**Each replica has two MRs–one for each plane**
- **MR for replication plane has consensus log**
- **MR for background plane has metadata for leader election and permission management (e.g., permissions array)**

**Replica can change QP permissions over time by changing QP access flags**

**All replicas can always read and write to the background plane MR on a replica**

**Necessary to set permission bits**

Mu's RDMA communication is carefully designed to support the two-plane architecture:

#### Queue Pair Organization

**Two QPs per remote replica**: Each machine maintains two queue pairs for each remote replica:
- **Replication QP**: Handles high-frequency replication operations
- **Background QP**: Handles low-frequency management operations

**Shared completion queues**: All QPs for the same plane share a completion queue:
- **Replication CQ**: Shared by all replication QPs
- **Background CQ**: Shared by all background QPs

This design provides:
- **Efficiency**: Shared CQs reduce memory overhead
- **Simplified management**: Easier to manage completion notifications
- **Performance**: Optimized for the specific workload of each plane

#### Memory Region Organization

**Two MRs per replica**: Each replica has two memory regions:
- **Replication MR**: Contains the consensus log and application data
- **Background MR**: Contains metadata for leader election and permission management

**Access control**: Different access modes for different operations:
- **Replication MR**: Only the current leader can write
- **Background MR**: All replicas can read and write (for permission management)

#### Dynamic Permission Management

**QP permission changes**: Replicas can change queue pair permissions over time:
- **Permission updates**: Change access flags to grant or revoke permissions
- **Dynamic adaptation**: Permissions can be updated as the system state changes
- **Hardware enforcement**: RDMA hardware enforces the permissions

**Background plane access**: All replicas can always access the background plane:
- **Permission management**: Necessary for setting and updating permission bits
- **Leader election**: Required for coordination during leader changes
- **System maintenance**: Enables system-wide coordination operations
## Replication Plane: Paxos-Based Consensus

### Mu's Paxos Foundation

**Mu builds on top of Paxos**

**Basic Paxos algorithm**

Mu uses Paxos as its underlying consensus protocol, but with significant optimizations enabled by RDMA. Understanding the basic Paxos algorithm is crucial for understanding how Mu works.

### Basic Paxos Algorithm

**Leader sends a proposal number to all replicas**

**Prepare phase**

**Three cases:**

#### Case (a): Higher Proposal Number
**If replica has already seen a higher proposal number, it sends a nack so that the leader can abort and try a different (higher) proposal number**

**I guess this occurs when either a new leader is out-of-date or there are two leaders racing?**

This case occurs when:
- **Out-of-date leader**: A new leader is trying to propose a value but is behind other leaders
- **Leader race**: Multiple leaders are competing for leadership
- **Network delays**: Messages arrive out of order due to network delays

The replica rejects the proposal because it has already seen a higher proposal number, indicating that another leader is more recent.

#### Case (b): Previous Value Accepted
**If the replica has not seen a higher proposal number but has accepted some value in the past, it returns that value along with the proposal number for that value. This returned proposal number is lower than the proposal number that the leader sent to the replica**

**If this occurs, the leader must use the largest proposal number (and corresponding value) sent to it by a replica. I guess this is because this value was already considered committed by a previous leader, so the new leader must continue to replicate it**

This case occurs when:
- **Previous leader**: A previous leader had already accepted a value
- **Leader failure**: The previous leader failed before completing the consensus
- **Recovery**: The new leader must continue with the previously accepted value

The new leader must use the previously accepted value to maintain consistency.

#### Case (c): No Previous Value
**If the replica has not seen a higher proposal number and has not accepted any value in the past, it sends an ack to the leader**

**If this occurs, then the leader can replicate its own value since the replicas are already up to date**

This case occurs when:
- **Clean state**: No previous leader has accepted any value
- **Fresh start**: The system is starting fresh or has been fully committed
- **Optimal case**: The leader can propose its own value

### Accept Phase

**Accept phase**

**Leader sends proposal number and value to all replicas**

**Replica sends ack if it has not seen any prepare message with higher proposal number**

After the prepare phase, the leader enters the accept phase:

1. **Proposal**: Leader sends the proposal number and value to all replicas
2. **Acceptance**: Replicas accept the proposal if they haven't seen a higher proposal number
3. **Commitment**: If a majority accepts, the value is committed

### Mu's Paxos Optimizations

**Mu is simpler than vanilla Paxos because it ensures that only one leader at a time can write to a replica**

**Leader maintains a set of confirmed followers**

**These are follower replicas that have granted write permission to the leader**

**No other leader can have write permission. This is important!**

**Thus, a race between 2+ leaders is avoided**

Mu's key insight is that by using RDMA write permissions, it can eliminate many of the complexities of traditional Paxos:

- **No leader races**: Only one leader can write to each replica at a time
- **Simplified protocol**: No need to handle complex race conditions
- **Hardware enforcement**: RDMA hardware ensures only one leader can write
- **Performance**: Eliminates the overhead of handling multiple leaders
### Mu's Log Structure

**Log has three components:**

#### (1) Minimum Proposal Number
**Follower replica publishes minProposal to its local memory**

**This is the minimum proposal number the follower replica can accept**

Each follower maintains a minimum proposal number that indicates:
- **Acceptance threshold**: The minimum proposal number it will accept
- **Consistency**: Ensures followers don't accept outdated proposals
- **Recovery**: Helps new leaders understand what proposals are acceptable

#### (2) First Undecided Offset (FUO)
**First undecided offset (FUO)**

**Lowest log index that is undecided on the replica**

The FUO indicates:
- **Progress tracking**: How far the replica has progressed in the log
- **Recovery**: Helps new leaders understand what work needs to be completed
- **Consistency**: Ensures all replicas are working on the same log positions

#### (3) Request Slots
**Sequence of slots for client requests**

**Each slot is a tuple consisting of the proposal number and the value (i.e., the client request itself)**

Each log slot contains:
- **Proposal number**: The Paxos proposal number for this slot
- **Value**: The actual client request to be executed
- **Metadata**: Additional information needed for consensus

### Mu's Consensus Algorithm

**Go over the algorithm in Listing 2**

**Once a leader enters the accept phase and successfully writes a log entry to the majority of replicas, the value is considered committed**

Mu's consensus algorithm works as follows:

1. **Prepare phase**: Leader sends proposal numbers to followers
2. **Accept phase**: Leader writes log entries to followers using RDMA
3. **Commitment**: Once majority of followers have the entry, it's committed
4. **Execution**: Followers apply committed entries to their state machines

The key insight is that Mu uses RDMA writes to implement the accept phase, which eliminates the need for traditional message passing and reduces latency significantly.
### Mu's Extensions and Optimizations

**Extensions**

#### Stragglers

**A replica may not have all of the committed values**

**If elected leader, the replica brings itself up-to-date with the log on the replica that has the highest FUO**

**Also brings its other confirmed followers up to date if needed**

Stragglers are replicas that fall behind due to network delays, temporary failures, or other issues. Mu handles stragglers by:

- **Leader election**: If a straggler is elected leader, it must catch up first
- **Log synchronization**: The new leader synchronizes with the replica that has the highest FUO
- **Follower updates**: The leader then brings other followers up to date as needed

#### Commit Piggybacking

**Commits**

**Want to avoid having the leader and follows communicate to indicate a value is committed**

**A follower assumes a value is committed when the next slot in the log is replicated**

**Known as commit piggybacking**

Traditional consensus protocols require explicit communication to indicate when a value is committed. Mu optimizes this by:

- **Implicit commitment**: Followers assume a value is committed when the next slot is replicated
- **No extra messages**: Eliminates the need for explicit commit messages
- **Performance**: Reduces the number of round-trips required

#### Prepare Phase Optimization

**Omitting the prepare phase**

**Once leader sees empty slots at some index at all of its followers, the leader knows all of the slots after that index are empty**

**Thus, the leader does not need to re-enter the prepare phase**

Mu can optimize the prepare phase by:

- **Empty slot detection**: If all followers have empty slots at a certain index, the leader knows all subsequent slots are empty
- **Skip prepare**: The leader can skip the prepare phase for these slots
- **Performance**: Reduces the overhead of the prepare phase

#### Dynamic Follower Management

**Growing the confirm followers set**

**The leader can add replicas to its confirmed followers set and bring them up to date as needed**

Mu supports dynamic addition of followers:

- **Follower addition**: New followers can be added to the system
- **Synchronization**: The leader brings new followers up to date
- **Permission management**: New followers are granted appropriate permissions

#### Replayer Optimization

**Replayer**

**A follower must not apply a log entry until the log entry is fully written**

**Mu writes a non-zero value to a canary byte and relies on RDMA to ensure the log entry write is ordered before the write to the canary byte**

**This is not guaranteed by RDMA but does happen in practice**

**Could also write checksum of log entry to canary section instead**

**The follower will not apply the log entry until the checksum it computes over the log entry matches the checksum in the canary section**

Mu uses a canary mechanism to ensure log entries are fully written before being applied:

- **Canary byte**: A special byte that indicates when a log entry is complete
- **Write ordering**: RDMA ensures the log entry is written before the canary byte
- **Checksum alternative**: Could use checksums instead of canary bytes for better reliability
- **Safety**: Ensures followers don't apply partial log entries
## Background Plane: Leader Election and Management

### Leader Election

**Leader election**

**"replica i decides that j is leader if j is the replica with the lowest id, among those that i considers to be alive"**

Mu uses a deterministic leader election algorithm:

- **Lowest ID wins**: The replica with the lowest ID among alive replicas becomes the leader
- **Deterministic**: All replicas will agree on who the leader should be
- **Simple**: Easy to understand and implement

#### Pull-Score Mechanism

**Follower uses a pull-score mechanism to detect if other replicas are alive**

**Replicas update counter in local memory**

**Replica does RDMA read on this counter**

**If counter has been updated, the replica's score is incremented**

**If the counter has not been updated, the replica's score is decremented**

**Once score goes below failure threshold, replica is considered failed**

**Once score goes above recovery threshold, the replica is considered alive**

**This seems to contradict text earlier in the paper that says a high score means failure**

Mu uses a pull-based scoring mechanism for failure detection:

1. **Counter updates**: Each replica continuously updates a counter in its local memory
2. **Periodic reads**: Other replicas periodically read this counter using RDMA
3. **Score calculation**: Based on whether the counter has been updated, replicas adjust their scores
4. **Threshold-based**: Replicas are considered failed or alive based on score thresholds

**Note**: There appears to be a contradiction in the paper about whether high or low scores indicate failure. This suggests the scoring mechanism may be more complex than described.

### Fate Sharing

**Fate sharing**

**If replication thread fails or gets stuck while the leader election thread runs fine, then no new entries can be committed and no new leader can be elected**

**Thus, the system is deadlocked**

**To fix this, the leader election thread checks for activity on the replication thread periodically (e.g., every 10,000 iterations)**

**If there is not activity on the replication thread, the leader election thread stops incrementing the counter**

**Thus, a new leader will be elected**

Fate sharing is a critical issue in Mu's two-plane architecture:

#### The Problem
- **Thread isolation**: The replication and leader election threads are separate
- **Deadlock risk**: If the replication thread fails but the leader election thread continues, the system can deadlock
- **No progress**: No new entries can be committed and no new leader can be elected

#### The Solution
- **Activity monitoring**: The leader election thread periodically checks for activity on the replication thread
- **Counter stopping**: If no activity is detected, the leader election thread stops incrementing its counter
- **Leader election**: This causes other replicas to consider the current leader failed and elect a new one

This mechanism ensures that the system can recover from failures in either plane.
### Permission Management

**Permission management**

**Can update permissions in three ways:**

#### Method 1: Re-register MR
**Re-register MR**

**Overhead scales linearly with size of MR (probably due to page table entry updates)**

Re-registering the memory region:
- **Complete re-registration**: Creates a new memory region with different permissions
- **Linear overhead**: Time scales with the size of the memory region
- **Page table updates**: Requires updating page table entries for the entire region
- **High cost**: Most expensive method for large memory regions

#### Method 2: Move QP to Non-Operational State
**Move QP to non-operational state**

**Overhead is independent of size of MR**

Moving the queue pair to non-operational state:
- **QP state change**: Changes the state of the queue pair
- **Constant overhead**: Time is independent of memory region size
- **Hardware operation**: Requires hardware-level state changes
- **Moderate cost**: Faster than re-registration but slower than flag changes

#### Method 3: Change QP Flags
**Grant/revoke access by changing flags on QP**

**Overhead is independent of size of MR**

Changing queue pair flags:
- **Flag modification**: Changes access flags on the queue pair
- **Constant overhead**: Time is independent of memory region size
- **Software operation**: Primarily a software operation
- **Low cost**: Fastest method for permission changes

#### Performance Comparison

**Changing flags is 10x faster than moving QP to non-operational state**

**But changing flags can fail if RDMA operations in flight**

**Thus, use fast-slow approach**

**Attempt to change flags first**

**If fails, move QP to non-operational state**

Mu uses a fast-slow approach for permission management:

1. **Fast path**: First attempt to change flags (10x faster)
2. **Fallback**: If flag changes fail (due to in-flight operations), fall back to moving QP to non-operational state
3. **Reliability**: Ensures permission changes always succeed
4. **Performance**: Optimizes for the common case where flag changes succeed

This approach provides the best of both worlds: fast performance in the common case and reliability when needed.
### Log Recycling

**Log recycling**

**Log head**

**minHead**

Mu implements log recycling to prevent memory exhaustion:

- **Log head**: The current position in the log where new entries are written
- **minHead**: The minimum position that can be safely recycled
- **Memory management**: Old log entries are recycled to prevent memory exhaustion
- **Safety**: Ensures committed entries are never recycled

### Adding and Removing Replicas

**Adding and removing replicas**

Mu supports dynamic addition and removal of replicas:

- **Replica addition**: New replicas can be added to the system
- **Replica removal**: Existing replicas can be removed from the system
- **Permission management**: Permissions are updated when replicas are added or removed
- **Log synchronization**: New replicas are synchronized with the current log state

## Evaluation

### Performance Results

Mu achieves remarkable performance improvements:

- **1.3 µs consensus overhead**: Dramatic improvement over traditional systems
- **99th percentile of 1.6 µs**: Consistent performance across all operations
- **873 µs failover time**: Order of magnitude improvement over prior systems
- **99% latency of 945 µs**: Fast recovery from failures

### Key Achievements

1. **Microsecond-scale consensus**: First system to achieve consensus in microseconds
2. **Hardware-optimized design**: Leverages RDMA for maximum performance
3. **Simplified protocol**: Eliminates many complexities of traditional consensus
4. **Fast failure recovery**: Sub-millisecond recovery from failures

### Limitations and Future Work

1. **RDMA requirement**: Limited to environments with RDMA support
2. **In-memory only**: No persistence guarantees
3. **Hardware dependency**: Performance limited by hardware capabilities
4. **Future improvements**: Battery-backed NICs for persistence

### Impact

Mu represents a significant breakthrough in consensus algorithm design, demonstrating that microsecond-scale consensus is achievable with the right hardware and software optimizations. It opens new possibilities for high-performance distributed systems and provides a foundation for future research in this area.