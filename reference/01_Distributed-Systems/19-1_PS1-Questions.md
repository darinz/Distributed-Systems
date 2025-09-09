# Problem Set 1: Fundamental Distributed Systems Concepts

## Learning Objectives
This problem set focuses on fundamental distributed systems concepts:
- Understanding RPC semantics and their implementation requirements
- Analyzing logical clocks and their relationship to real-time
- Comparing caching strategies and their communication patterns
- Evaluating transaction protocols and their consistency guarantees
- Exploring distributed storage systems and their consistency models
- Understanding consensus protocols and their correctness properties
- Analyzing Byzantine fault tolerance and attack vectors
- Examining time synchronization and its impact on system correctness

---

## Question 1: True/False Analysis with Detailed Reasoning

**Context**: Understanding the fundamental properties and trade-offs of distributed systems is crucial for designing and implementing reliable systems. This question tests your understanding of key concepts through true/false statements that require careful analysis.

Answer true or false, explain your reasoning, and (where appropriate) give an example (if true) or counterexample (if false).

### 1a. Exactly-once RPC requires both sender and receiver to write to stable storage

**Detailed Hints:**

### Understanding Exactly-once RPC
- **Think about**: What does "exactly-once" mean in the context of RPC?
- **Consider**: What are the challenges in implementing exactly-once semantics?
- **Key insight**: Exactly-once requires handling duplicate requests and ensuring idempotency
- **Analysis approach**: Consider what information needs to be persisted to achieve exactly-once semantics

### Step-by-step Analysis

#### Step 1: Understand Exactly-once Semantics
**Detailed Hints:**
- **Think about**: What guarantees does exactly-once RPC provide?
- **Consider**: How does the system handle duplicate requests?
- **Key insight**: Exactly-once means the operation is executed exactly once, even if the request is sent multiple times

#### Step 2: Analyze Implementation Requirements
**Detailed Hints:**
- **Think about**: What information does the sender need to track?
- **Consider**: What information does the receiver need to track?
- **Key insight**: Both sides need to maintain state to handle duplicates and ensure idempotency

### Detailed Answer

**True.**

**Explanation:**
Exactly-once RPC requires both sender and receiver to write to stable storage for the following reasons:

1. **Sender requirements**: The sender must track which requests have been sent and which have been acknowledged to handle network failures and retries
2. **Receiver requirements**: The receiver must track which requests have been processed to handle duplicate requests and ensure idempotency
3. **Crash recovery**: Both sides need persistent state to recover from crashes and maintain exactly-once semantics
4. **Duplicate detection**: The receiver needs to identify and ignore duplicate requests

**Example**: A banking system where a transfer operation must be executed exactly once, even if the network retries the request multiple times.

### 1b. Two events can have the same Lamport clock time value but occur minutes apart in real-time

**Detailed Hints:**

### Understanding Lamport Clocks
- **Think about**: How do Lamport clocks work?
- **Consider**: What relationship do Lamport clocks have to real-time?
- **Key insight**: Lamport clocks capture causal relationships, not real-time ordering
- **Analysis approach**: Consider scenarios where causally unrelated events can have the same timestamp

### Step-by-step Analysis

#### Step 1: Understand Lamport Clock Properties
**Detailed Hints:**
- **Think about**: What do Lamport clocks measure?
- **Consider**: How are Lamport timestamps assigned?
- **Key insight**: Lamport clocks only increment when events occur or messages are received

#### Step 2: Analyze Concurrent Events
**Detailed Hints:**
- **Think about**: What happens with concurrent events?
- **Consider**: Can concurrent events have the same timestamp?
- **Key insight**: Concurrent events that don't communicate can have the same timestamp

### Detailed Answer

**True.**

**Explanation:**
Lamport clocks capture causal relationships, not real-time ordering. Two events can have the same Lamport timestamp if they are concurrent (neither causally precedes the other).

**Example**: Consider two processes that don't communicate:
- Process A: Event X occurs at real-time 10:00 AM, gets Lamport timestamp 5
- Process B: Event Y occurs at real-time 10:05 AM, gets Lamport timestamp 5
- Both events have the same Lamport timestamp (5) but occur 5 minutes apart in real-time

### 1c. Write-back caching never does more communication than write-through caching

**Detailed Hints:**

### Understanding Caching Strategies
- **Think about**: How do write-back and write-through caching work?
- **Consider**: When does each strategy perform communication?
- **Key insight**: Write-back delays writes until eviction, write-through writes immediately
- **Analysis approach**: Consider scenarios where write-back might require more communication

### Step-by-step Analysis

#### Step 1: Understand Write-through Caching
**Detailed Hints:**
- **Think about**: When does write-through caching communicate?
- **Consider**: What happens on every write operation?
- **Key insight**: Write-through writes to both cache and storage on every write

#### Step 2: Understand Write-back Caching
**Detailed Hints:**
- **Think about**: When does write-back caching communicate?
- **Consider**: What triggers communication in write-back?
- **Key insight**: Write-back only communicates when cache blocks are evicted

#### Step 3: Analyze Communication Scenarios
**Detailed Hints:**
- **Think about**: Can write-back ever require more communication?
- **Consider**: What happens with frequent evictions?
- **Key insight**: Write-back can require more communication if blocks are evicted frequently

### Detailed Answer

**False.**

**Explanation:**
Write-back caching can do more communication than write-through caching in certain scenarios.

**Counterexample**: Consider a scenario with frequent cache evictions:
- Write-through: Each write requires 1 communication (write to storage)
- Write-back: If a cache block is evicted after every write, each write requires 2 communications (write to storage for evicted block + write to cache for new block)

**Additional scenarios where write-back does more communication:**
- Cache thrashing with frequent evictions
- Large cache blocks that are frequently modified and evicted
- Systems with limited cache capacity relative to working set

### 1d. Two-phase commit with locks released before commit reaches stable storage is serializable

**Detailed Hints:**

### Understanding Two-Phase Commit
- **Think about**: How does two-phase commit work?
- **Consider**: When are locks released in the protocol?
- **Key insight**: Locks are typically released after the commit decision is made
- **Analysis approach**: Consider what happens if locks are released before the commit is persistent

### Step-by-step Analysis

#### Step 1: Understand Two-Phase Commit Protocol
**Detailed Hints:**
- **Think about**: What are the two phases of 2PC?
- **Consider**: When is the commit decision made?
- **Key insight**: The commit decision is made in the second phase

#### Step 2: Analyze Lock Release Timing
**Detailed Hints:**
- **Think about**: What happens if locks are released before commit is persistent?
- **Consider**: Can other transactions access the data?
- **Key insight**: Early lock release can allow other transactions to see uncommitted data

### Detailed Answer

**False.**

**Explanation:**
Two-phase commit with locks released before commit reaches stable storage is not serializable.

**Counterexample**: Consider the following scenario:
1. Transaction T1 prepares and gets locks on data items
2. T1 receives "commit" from coordinator
3. T1 releases locks before writing commit to stable storage
4. T1 crashes before writing commit to stable storage
5. Transaction T2 starts and reads the data items (locks are released)
6. T1 recovers and aborts (commit was not in stable storage)
7. T2 has read data that was never actually committed

This violates serializability because T2 sees the effects of T1's operations, but T1 was never committed.

### 1e. Updates in Dynamo are not serializable

**Detailed Hints:**

### Understanding Dynamo's Consistency Model
- **Think about**: What consistency model does Dynamo provide?
- **Consider**: How does Dynamo handle concurrent updates?
- **Key insight**: Dynamo provides eventual consistency, not strong consistency
- **Analysis approach**: Consider what serializability means and whether Dynamo provides it

### Step-by-step Analysis

#### Step 1: Understand Dynamo's Design
**Detailed Hints:**
- **Think about**: What are Dynamo's design goals?
- **Consider**: How does Dynamo handle conflicts?
- **Key insight**: Dynamo prioritizes availability over consistency

#### Step 2: Analyze Serializability
**Detailed Hints:**
- **Think about**: What does serializability require?
- **Consider**: Can Dynamo guarantee serializable execution?
- **Key insight**: Serializability requires a total order of transactions

### Detailed Answer

**True.**

**Explanation:**
Updates in Dynamo are not serializable because Dynamo provides eventual consistency, not strong consistency.

**Example**: Consider two concurrent updates to the same key:
- Client A updates key "user123" with value "Alice"
- Client B updates key "user123" with value "Bob"
- Both updates succeed at different replicas
- The system eventually converges to one value, but there's no guarantee about which one
- This violates serializability because there's no total order of the updates

**Why Dynamo doesn't provide serializability:**
- Dynamo prioritizes availability over consistency
- It uses vector clocks to detect conflicts but doesn't resolve them deterministically
- The system provides eventual consistency, not strong consistency

### 1f. GFS may contain duplicate records even if application writes each record once

**Detailed Hints:**

### Understanding GFS Record Append
- **Think about**: How does GFS handle record appends?
- **Consider**: What can cause duplicates in GFS?
- **Key insight**: GFS's record append operation can create duplicates due to retries
- **Analysis approach**: Consider scenarios where the application writes once but GFS creates duplicates

### Step-by-step Analysis

#### Step 1: Understand GFS Record Append
**Detailed Hints:**
- **Think about**: How does GFS's record append work?
- **Consider**: What happens if the operation fails?
- **Key insight**: GFS doesn't guarantee exactly-once semantics for record appends

#### Step 2: Analyze Duplicate Scenarios
**Detailed Hints:**
- **Think about**: What can cause duplicates in GFS?
- **Consider**: How does GFS handle failures and retries?
- **Key insight**: Network failures and retries can cause duplicates

### Detailed Answer

**True.**

**Explanation:**
GFS may contain duplicate records even if the application writes each record once due to the nature of GFS's record append operation.

**Example**: Consider the following scenario:
1. Application sends a record append request to GFS
2. The request times out due to network issues
3. Application retries the same record append request
4. Both requests eventually reach GFS and are processed
5. The same record appears twice in the file

**Why this happens:**
- GFS's record append doesn't provide exactly-once semantics
- Network failures can cause timeouts and retries
- GFS prioritizes performance over strict consistency
- The system doesn't prevent duplicate appends from the same client

## Question 2: Distributed Systems Features Analysis

**Context**: Understanding how fundamental distributed systems concepts are implemented in real-world systems is crucial for appreciating the practical applications of theoretical concepts. This question explores how different systems implement key distributed systems features.

Consider the following systems: git, Facebook's use of memcache, GFS, BigTable, Spanner, Dynamo, and Bitcoin. For each of the following features, find one example of that feature in one of the systems, and provide a detailed analysis of its role in the system. Use each system as an answer at most once.

### 2g. RPC (Remote Procedure Call)

**Detailed Hints:**

### Understanding RPC
- **Think about**: Which of these systems uses RPC for communication?
- **Consider**: How does RPC enable distributed system communication?
- **Key insight**: RPC allows processes to call functions on remote machines as if they were local
- **Analysis approach**: Look for systems that need to communicate between distributed components

### Step-by-step Analysis

#### Step 1: Identify RPC Usage
**Detailed Hints:**
- **Think about**: Which systems have client-server architectures?
- **Consider**: Which systems need to make remote calls?
- **Key insight**: Systems with distributed components typically use RPC

#### Step 2: Analyze RPC Role
**Detailed Hints:**
- **Think about**: How does RPC enable the system's functionality?
- **Consider**: What would happen without RPC in this system?
- **Key insight**: RPC is essential for distributed system operation

### Detailed Answer

**System: BigTable**

**Role of RPC in BigTable:**
BigTable uses RPC extensively for communication between its distributed components. Clients use RPC to communicate with tablet servers to read and write data. The BigTable master uses RPC to communicate with tablet servers for metadata operations, tablet assignments, and load balancing. Tablet servers use RPC to communicate with each other for tablet splits and other coordination tasks.

**Why RPC is essential:**
- Enables clients to access data stored on remote tablet servers
- Allows the master to manage and coordinate tablet servers
- Provides a clean abstraction for distributed operations
- Enables fault tolerance through remote procedure calls

### 2h. Caching

**Detailed Hints:**

### Understanding Caching
- **Think about**: Which systems use caching to improve performance?
- **Consider**: How does caching reduce latency and improve throughput?
- **Key insight**: Caching stores frequently accessed data in faster storage
- **Analysis approach**: Look for systems that need to reduce access times

### Step-by-step Analysis

#### Step 1: Identify Caching Usage
**Detailed Hints:**
- **Think about**: Which systems have performance requirements?
- **Consider**: Which systems access data frequently?
- **Key insight**: Systems with high read loads often use caching

#### Step 2: Analyze Caching Role
**Detailed Hints:**
- **Think about**: How does caching improve system performance?
- **Consider**: What would happen without caching in this system?
- **Key insight**: Caching is crucial for meeting performance requirements

### Detailed Answer

**System: Facebook's use of memcache**

**Role of Caching in Facebook's memcache:**
Facebook uses memcache as a distributed caching layer to store frequently accessed data (user profiles, friend lists, posts) in memory. When a front-end server needs data, it first checks the memcache server. If the data is not cached, the front-end retrieves it from the storage server and stores it in memcache for future requests.

**Why caching is essential:**
- Dramatically reduces latency for frequently accessed data
- Reduces load on storage servers
- Enables Facebook to serve millions of users with low latency
- Improves overall system throughput and scalability

### 2i. Eventual Consistency

**Detailed Hints:**

### Understanding Eventual Consistency
- **Think about**: Which systems prioritize availability over strong consistency?
- **Consider**: How does eventual consistency enable high availability?
- **Key insight**: Eventual consistency allows systems to remain available during partitions
- **Analysis approach**: Look for systems that need high availability

### Step-by-step Analysis

#### Step 1: Identify Eventual Consistency Usage
**Detailed Hints:**
- **Think about**: Which systems need to remain available during failures?
- **Consider**: Which systems can tolerate temporary inconsistencies?
- **Key insight**: Systems prioritizing availability often use eventual consistency

#### Step 2: Analyze Eventual Consistency Role
**Detailed Hints:**
- **Think about**: How does eventual consistency enable availability?
- **Consider**: What are the trade-offs of eventual consistency?
- **Key insight**: Eventual consistency is a key design choice for availability

### Detailed Answer

**System: Dynamo**

**Role of Eventual Consistency in Dynamo:**
Dynamo uses eventual consistency to ensure high availability and partition tolerance. When multiple replicas receive different updates, Dynamo doesn't immediately resolve conflicts. Instead, it stores multiple versions and uses vector clocks to detect conflicts. The system eventually converges to a consistent state, but there's no guarantee about which version will be chosen.

**Why eventual consistency is essential:**
- Enables the system to remain available during network partitions
- Allows writes to succeed even when some replicas are unavailable
- Provides high availability, which is crucial for Amazon's e-commerce platform
- Trades strong consistency for availability and partition tolerance

### 2j. Serializability

**Detailed Hints:**

### Understanding Serializability
- **Think about**: Which systems need strong consistency guarantees?
- **Consider**: How does serializability ensure correctness?
- **Key insight**: Serializability ensures that concurrent transactions appear to execute in some serial order
- **Analysis approach**: Look for systems that need strong consistency

### Step-by-step Analysis

#### Step 1: Identify Serializability Usage
**Detailed Hints:**
- **Think about**: Which systems need strong consistency?
- **Consider**: Which systems handle financial or critical data?
- **Key insight**: Systems requiring strong consistency often provide serializability

#### Step 2: Analyze Serializability Role
**Detailed Hints:**
- **Think about**: How does serializability ensure correctness?
- **Consider**: What are the benefits and costs of serializability?
- **Key insight**: Serializability is crucial for maintaining data integrity

### Detailed Answer

**System: Spanner**

**Role of Serializability in Spanner:**
Spanner provides serializable transactions across globally distributed data. It uses TrueTime and two-phase commit to ensure that all transactions appear to execute in some serial order, even across multiple data centers. This guarantees that concurrent transactions don't interfere with each other and maintains data integrity.

**Why serializability is essential:**
- Ensures data consistency across global operations
- Prevents race conditions and data corruption
- Enables complex transactions across multiple data centers
- Provides strong consistency guarantees for critical applications

### 2k. Logging

**Detailed Hints:**

### Understanding Logging
- **Think about**: Which systems need to maintain logs for recovery or auditing?
- **Consider**: How does logging enable system recovery and consistency?
- **Key insight**: Logging records system state changes for recovery and debugging
- **Analysis approach**: Look for systems that need recovery mechanisms

### Step-by-step Analysis

#### Step 1: Identify Logging Usage
**Detailed Hints:**
- **Think about**: Which systems need to recover from failures?
- **Consider**: Which systems need to maintain audit trails?
- **Key insight**: Systems requiring reliability often use logging

#### Step 2: Analyze Logging Role
**Detailed Hints:**
- **Think about**: How does logging enable system recovery?
- **Consider**: What are the benefits of logging?
- **Key insight**: Logging is crucial for system reliability and debugging

### Detailed Answer

**System: GFS**

**Role of Logging in GFS:**
GFS uses logging extensively for recovery and consistency. The master maintains logs of all metadata operations, including file creation, deletion, and chunk assignments. Chunk servers maintain logs of all write operations. These logs enable the system to recover from failures and maintain consistency across the distributed file system.

**Why logging is essential:**
- Enables recovery from master and chunk server failures
- Maintains consistency of metadata and data
- Provides audit trails for debugging and analysis
- Ensures durability of operations in the distributed file system

### 2l. State Machine Replication

**Detailed Hints:**

### Understanding State Machine Replication
- **Think about**: Which systems need to maintain consistent state across replicas?
- **Consider**: How does state machine replication ensure consistency?
- **Key insight**: State machine replication ensures all replicas execute the same sequence of operations
- **Analysis approach**: Look for systems that need strong consistency across replicas

### Step-by-step Analysis

#### Step 1: Identify State Machine Replication Usage
**Detailed Hints:**
- **Think about**: Which systems need consistent replicas?
- **Consider**: Which systems use consensus protocols?
- **Key insight**: Systems requiring strong consistency often use state machine replication

#### Step 2: Analyze State Machine Replication Role
**Detailed Hints:**
- **Think about**: How does state machine replication ensure consistency?
- **Consider**: What are the benefits of state machine replication?
- **Key insight**: State machine replication is crucial for maintaining consistency

### Detailed Answer

**System: Bitcoin**

**Role of State Machine Replication in Bitcoin:**
Bitcoin uses state machine replication to maintain a consistent blockchain across all nodes. All nodes execute the same sequence of transactions in the same order, ensuring that the blockchain state is consistent across the network. This is achieved through the proof-of-work consensus mechanism, which ensures that all honest nodes agree on the same sequence of blocks.

**Why state machine replication is essential:**
- Ensures all nodes maintain the same blockchain state
- Prevents double-spending and other inconsistencies
- Enables decentralized consensus without a central authority
- Provides the foundation for Bitcoin's security and consistency

### 2m. Hint

**Detailed Hints:**

### Understanding Hints
- **Think about**: Which systems use hints to improve performance or recovery?
- **Consider**: How do hints help systems make better decisions?
- **Key insight**: Hints provide additional information to help systems optimize their behavior
- **Analysis approach**: Look for systems that use additional information for optimization

### Step-by-step Analysis

#### Step 1: Identify Hint Usage
**Detailed Hints:**
- **Think about**: Which systems use additional information for optimization?
- **Consider**: Which systems need to make decisions about data placement or access?
- **Key insight**: Systems with optimization needs often use hints

#### Step 2: Analyze Hint Role
**Detailed Hints:**
- **Think about**: How do hints improve system performance?
- **Consider**: What are the benefits of using hints?
- **Key insight**: Hints are crucial for system optimization

### Detailed Answer

**System: git**

**Role of Hints in git:**
Git uses hints in several ways to improve performance. For example, git uses hints about file similarity to optimize delta compression, storing only the differences between similar files. Git also uses hints about commit history to optimize merge operations and conflict resolution. Additionally, git uses hints about file access patterns to optimize pack file organization.

**Why hints are essential:**
- Improves compression efficiency by identifying similar files
- Optimizes merge operations by understanding commit relationships
- Enhances performance of common operations like diff and merge
- Enables git to scale to large repositories with millions of files

**Key Learning Points:**
- Different systems implement the same concepts in different ways
- Understanding real-world implementations helps appreciate theoretical concepts
- System design choices reflect different trade-offs and requirements
- Each feature serves a specific purpose in enabling system functionality

## Question 3: Facebook Three-Tier System Semantic Analysis

**Context**: Understanding the consistency semantics of distributed caching systems is crucial for designing reliable applications. This question explores how different ordering of cache invalidation and storage updates affects system behavior and consistency guarantees.

Facebook uses a three-tier system for implementing its website. An array of front-end servers interacts with web clients (each client is hashed into exactly one front-end server); these front-end servers gather the information needed to render the client web page from an array of cache servers and a separate array of storage servers. Hashing is used to locate which cache and storage server might have a particular object (e.g., a friend list, or set of postings). The number of front-end servers, cache servers, and storage servers is not identical (the numbers are chosen to balance the workload), so in general, all front-ends talk to all cache servers and all storage servers.

The cache servers (called memcache servers) are managed as a "lookaside" cache. When rendering an object on a page, the front-end first sends a message to the relevant memcache server; if the data is not available, the front-end (not the cache) then retrieves the data from the relevant storage server. The front-end then stores the fetched data into the memcache server. On update, the front-end invalidates the cached copy (if any) and updates the storage server.

### 3a. Cache Invalidation Before Storage Update

**Question**: What semantics would occur if the front-end first invalidates the cache, and then updates the storage server?

**Detailed Hints:**

### Understanding the Scenario
- **Think about**: What happens when the cache is invalidated before the storage is updated?
- **Consider**: What can go wrong during the time between cache invalidation and storage update?
- **Key insight**: There's a window where the cache is empty but storage hasn't been updated yet
- **Analysis approach**: Consider what happens if another request occurs during this window

### Step-by-step Analysis

#### Step 1: Understand the Operation Sequence
**Detailed Hints:**
- **Think about**: What is the exact sequence of operations?
- **Consider**: What state is the system in after each step?
- **Key insight**: The system goes through different states during the update process

#### Step 2: Analyze Potential Issues
**Detailed Hints:**
- **Think about**: What can happen if another request occurs during the update?
- **Consider**: What data will the other request see?
- **Key insight**: Other requests might see stale data or cause inconsistencies

### Detailed Answer

**Semantics: Potential for stale data reads**

**What happens:**
1. Front-end invalidates cache (cache becomes empty)
2. Front-end updates storage server (storage has new data)
3. If another request occurs between steps 1 and 2, it will:
   - Find cache empty
   - Read from storage (which still has old data)
   - Cache the old data
   - Return old data to the client

**Problems:**
- **Stale data reads**: Clients can read stale data if requests occur during the update window
- **Cache pollution**: The cache gets populated with stale data
- **Inconsistency**: Different clients might see different values depending on timing

**Example scenario:**
- Initial state: Cache has "Alice", Storage has "Alice"
- Update starts: Cache invalidated (empty), Storage still has "Alice"
- Another request: Reads "Alice" from storage, caches "Alice"
- Update completes: Storage updated to "Bob", but cache still has "Alice"
- Result: Cache has stale data "Alice" while storage has "Bob"

### 3b. Storage Update Before Cache Invalidation

**Question**: What semantics would occur if the front-end updates the storage server and then invalidates the cache?

**Detailed Hints:**

### Understanding the Scenario
- **Think about**: What happens when storage is updated before cache is invalidated?
- **Consider**: What can go wrong during the time between storage update and cache invalidation?
- **Key insight**: There's a window where storage has new data but cache still has old data
- **Analysis approach**: Consider what happens if another request occurs during this window

### Step-by-step Analysis

#### Step 1: Understand the Operation Sequence
**Detailed Hints:**
- **Think about**: What is the exact sequence of operations?
- **Consider**: What state is the system in after each step?
- **Key insight**: The system goes through different states during the update process

#### Step 2: Analyze Potential Issues
**Detailed Hints:**
- **Think about**: What can happen if another request occurs during the update?
- **Consider**: What data will the other request see?
- **Key insight**: Other requests might see stale data from cache

### Detailed Answer

**Semantics: Potential for stale data reads**

**What happens:**
1. Front-end updates storage server (storage has new data)
2. Front-end invalidates cache (cache becomes empty)
3. If another request occurs between steps 1 and 2, it will:
   - Find cache with old data
   - Return old data to the client
   - Not update the cache

**Problems:**
- **Stale data reads**: Clients can read stale data from cache during the update window
- **Inconsistency**: Storage has new data but cache serves old data
- **Temporary inconsistency**: The system is inconsistent until cache invalidation completes

**Example scenario:**
- Initial state: Cache has "Alice", Storage has "Alice"
- Update starts: Storage updated to "Bob", Cache still has "Alice"
- Another request: Reads "Alice" from cache, returns "Alice"
- Update completes: Cache invalidated (empty)
- Result: Client received stale data "Alice" while storage had "Bob"

### 3c. Double Cache Invalidation

**Question**: What semantics would occur if the front-end invalidates the cache, updates the storage server, and then re-invalidates the cache?

**Detailed Hints:**

### Understanding the Scenario
- **Think about**: What happens with double cache invalidation?
- **Consider**: Is the second invalidation necessary or harmful?
- **Key insight**: The second invalidation might be redundant or might cause issues
- **Analysis approach**: Consider what happens if another request occurs during the update process

### Step-by-step Analysis

#### Step 1: Understand the Operation Sequence
**Detailed Hints:**
- **Think about**: What is the exact sequence of operations?
- **Consider**: What state is the system in after each step?
- **Key insight**: The system goes through different states during the update process

#### Step 2: Analyze the Second Invalidation
**Detailed Hints:**
- **Think about**: What is the purpose of the second invalidation?
- **Consider**: Can the second invalidation cause problems?
- **Key insight**: The second invalidation might be unnecessary or might cause issues

### Detailed Answer

**Semantics: Redundant but safe**

**What happens:**
1. Front-end invalidates cache (cache becomes empty)
2. Front-end updates storage server (storage has new data)
3. Front-end re-invalidates cache (cache remains empty)

**Analysis:**
- **First invalidation**: Ensures cache is empty before update
- **Storage update**: Updates the authoritative data source
- **Second invalidation**: Redundant since cache is already empty

**Benefits:**
- **Safety**: Ensures cache is definitely empty after update
- **Consistency**: Guarantees that subsequent reads will get fresh data from storage
- **Defensive programming**: Protects against race conditions

**Potential issues:**
- **Redundancy**: The second invalidation is unnecessary
- **Performance**: Extra network round-trip for no benefit
- **Complexity**: More complex than necessary

### 3d. Write-Token Algorithm

**Question**: What semantics would occur with the write-token algorithm?

An employee at Facebook suggests adding a write-token to the memcache server. When a front-end wants to change a value, it sends a message to memcache to atomically invalidate the entry and set the write-token; subsequent accesses to the server stall. The front-end releases the write-token when the data is updated at the server, allowing stalled accesses to proceed. What semantics would occur in this algorithm?

**Detailed Hints:**

### Understanding the Write-Token Algorithm
- **Think about**: How does the write-token mechanism work?
- **Consider**: What happens when a write-token is acquired?
- **Key insight**: The write-token prevents concurrent access during updates
- **Analysis approach**: Consider how this affects consistency and performance

### Step-by-step Analysis

#### Step 1: Understand the Write-Token Mechanism
**Detailed Hints:**
- **Think about**: What does acquiring a write-token do?
- **Consider**: What happens to other requests during the write-token?
- **Key insight**: The write-token acts as a lock on the cache entry

#### Step 2: Analyze the Semantics
**Detailed Hints:**
- **Think about**: What consistency guarantees does this provide?
- **Consider**: What are the performance implications?
- **Key insight**: The write-token ensures atomic updates but may impact performance

### Detailed Answer

**Semantics: Strong consistency with potential performance impact**

**What happens:**
1. Front-end acquires write-token and invalidates cache entry
2. Subsequent read requests to the same key stall (wait)
3. Front-end updates storage server
4. Front-end releases write-token
5. Stalled read requests proceed and read fresh data from storage

**Benefits:**
- **Strong consistency**: No stale data reads possible
- **Atomic updates**: Cache invalidation and storage update are coordinated
- **Race condition prevention**: No window for inconsistent reads
- **Guaranteed freshness**: All reads after update get fresh data

**Potential issues:**
- **Performance impact**: Read requests stall during updates
- **Latency increase**: Users experience higher latency during updates
- **Scalability concerns**: High update frequency could cause many stalled requests
- **Deadlock potential**: If write-token holder crashes, requests could stall indefinitely

**Example scenario:**
- Initial state: Cache has "Alice", Storage has "Alice"
- Update starts: Write-token acquired, cache invalidated, reads stall
- Storage updated: Storage now has "Bob"
- Write-token released: Stalled reads proceed, read "Bob" from storage, cache "Bob"
- Result: All clients see consistent data "Bob"

**Key Learning Points:**
- Different update orderings provide different consistency guarantees
- Cache invalidation timing affects system consistency
- Write-token algorithms provide strong consistency but may impact performance
- Understanding consistency semantics is crucial for system design

## Question 4: Paxos Maximum Values Analysis

**Context**: Understanding the theoretical limits of consensus protocols is crucial for analyzing their behavior and designing systems that use them. This question explores the maximum number of unique values that can be proposed in a single Paxos instance.

**Question**: What is the maximum number of unique values that can be proposed to a group of k Paxos acceptors (for a single instance of the protocol)?

**Detailed Hints:**

### Understanding Paxos Value Proposals
- **Think about**: How does Paxos handle multiple value proposals?
- **Consider**: What constraints does Paxos place on value proposals?
- **Key insight**: Paxos ensures that only one value can be chosen, but multiple values can be proposed
- **Analysis approach**: Consider the worst-case scenario for value proposals

### Step-by-step Analysis

#### Step 1: Understand Paxos Safety Guarantees
**Detailed Hints:**
- **Think about**: What does Paxos guarantee about chosen values?
- **Consider**: How does Paxos prevent multiple values from being chosen?
- **Key insight**: Paxos ensures that at most one value is chosen

#### Step 2: Analyze Value Proposal Scenarios
**Detailed Hints:**
- **Think about**: What is the worst-case scenario for value proposals?
- **Consider**: How many different values can be proposed before one is chosen?
- **Key insight**: The maximum number of unique values depends on the failure scenario

### Detailed Answer

**The maximum number of unique values that can be proposed is k (the number of acceptors).**

**Reasoning:**

1. **Worst-case scenario**: Each acceptor can propose a different value
2. **Paxos safety**: Only one value can ultimately be chosen
3. **Proposal process**: Each acceptor can propose a unique value during the prepare phase
4. **Acceptance process**: Acceptors can accept different values before consensus is reached

**Detailed analysis:**
- **Initial state**: All acceptors start with no accepted proposals
- **Proposal phase**: Each of the k acceptors can propose a different value
- **Acceptance phase**: Acceptors can accept different values initially
- **Consensus**: Eventually, one value will be chosen by a majority

**Example with k=3 acceptors:**
- Acceptor A proposes value "x"
- Acceptor B proposes value "y"  
- Acceptor C proposes value "z"
- Total: 3 unique values proposed
- Result: Only one value (e.g., "x") will be chosen by majority

**Key insight**: While k unique values can be proposed, Paxos guarantees that only one value will be chosen, ensuring safety.

## Question 5: Paxos Sequence Scenario Analysis

**Context**: Understanding how Paxos prevents certain problematic execution scenarios is crucial for appreciating its safety guarantees. This question explores why certain sequences of events cannot occur in Paxos and what actually happens instead.

**Question**: In Paxos, suppose that the acceptors are A, B, and C. A and B are also proposers, and there is a distinguished learner L. According to the Paxos paper, a value is chosen when a majority of acceptors accept it, and only a single value is chosen. How does Paxos ensure that the following sequence of events cannot happen? What actually happens, and which value is ultimately chosen?

**Given sequence:**
- a) A proposes sequence number 1, and gets responses from A, B, and C.
- b) A sends accept(1, "foo") messages to A and C and gets responses from both. Because a majority accepted, A tells L that "foo" has been chosen. However, A crashes before sending an accept to B.
- c) B proposes sequence number 2, and gets responses from B and C.
- d) B sends accept(2, "bar") messages to B and C and gets responses from both, so B tells L that "bar" has been chosen.

**Detailed Hints:**

### Understanding the Problematic Sequence
- **Think about**: What does this sequence claim happened?
- **Consider**: What are the key events that seem problematic?
- **Key insight**: The sequence suggests that both "foo" and "bar" could be chosen
- **Analysis approach**: Consider how Paxos prevents this scenario

### Step-by-step Analysis

#### Step 1: Analyze the First Phase (A's Proposal)
**Detailed Hints:**
- **Think about**: What does A learn from the prepare responses?
- **Consider**: What information do acceptors include in their prepare responses?
- **Key insight**: Prepare responses include the highest-numbered proposal each acceptor has accepted

#### Step 2: Analyze A's Accept Phase
**Detailed Hints:**
- **Think about**: What value should A propose in the accept phase?
- **Consider**: What does the Paxos algorithm require when no previous proposal has been accepted?
- **Key insight**: A can propose its own value "foo" since no previous proposal was accepted

#### Step 3: Analyze B's Prepare Phase
**Detailed Hints:**
- **Think about**: What does B learn from B and C's prepare responses?
- **Consider**: What information do B and C have about previously accepted proposals?
- **Key insight**: C accepted (1, "foo"), so it must report this to B

#### Step 4: Analyze B's Accept Phase
**Detailed Hints:**
- **Think about**: What value must B propose in the accept phase?
- **Consider**: What does the Paxos algorithm require when a previous proposal has been accepted?
- **Key insight**: B must propose "foo" (not "bar") to maintain safety

### Detailed Answer

**How Paxos prevents this sequence:**

Paxos prevents this sequence through the prepare phase, which informs proposers about previously accepted proposals.

**What actually happens:**

1. **A's prepare phase**: A sends prepare(1) to A, B, C
   - All respond with "no previous proposal accepted"
   - A can propose its own value "foo"

2. **A's accept phase**: A sends accept(1, "foo") to A, B, C
   - A and C accept (1, "foo")
   - A crashes before sending to B
   - A tells L that "foo" is chosen

3. **B's prepare phase**: B sends prepare(2) to B, C
   - B responds with "no previous proposal accepted"
   - **C responds with "I accepted (1, 'foo')"**

4. **B's accept phase**: B learns that C accepted (1, "foo")
   - B must propose "foo" (not "bar") to maintain safety
   - B sends accept(2, "foo") to B, C
   - Both accept (2, "foo")
   - B tells L that "foo" is chosen

**Which value is ultimately chosen:**
Value "foo" is ultimately chosen.

**Why the problematic sequence cannot happen:**
- The prepare phase ensures that proposers learn about previously accepted proposals
- Proposers must propose values from previously accepted proposals to maintain safety
- This prevents conflicting values from being chosen

## Question 6: Raft Byzantine Node Attack Analysis

**Context**: Understanding how Byzantine faults can compromise consensus protocols is crucial for designing secure distributed systems. This question explores how a Byzantine node can violate Raft's correctness constraints.

**Question**: For the Raft algorithm described in the reading list, outline how a Byzantine node would be able to cause each of the correctness constraints to be violated.

**Detailed Hints:**

### Understanding Raft Correctness Constraints
- **Think about**: What are Raft's main correctness constraints?
- **Consider**: How can a Byzantine node violate these constraints?
- **Key insight**: Byzantine nodes can behave arbitrarily and maliciously
- **Analysis approach**: Consider each correctness constraint and how it can be violated

### Step-by-step Analysis

#### Step 1: Identify Raft's Correctness Constraints
**Detailed Hints:**
- **Think about**: What guarantees does Raft provide?
- **Consider**: What are the main safety and liveness properties?
- **Key insight**: Raft provides election safety, leader append-only, log matching, leader completeness, and state machine safety

#### Step 2: Analyze Byzantine Attack Vectors
**Detailed Hints:**
- **Think about**: How can a Byzantine node attack each constraint?
- **Consider**: What malicious behaviors are possible?
- **Key insight**: Byzantine nodes can send conflicting messages, lie about their state, etc.

### Detailed Answer

**Raft's correctness constraints and how a Byzantine node can violate them:**

### 1. Election Safety
**Constraint**: At most one leader can be elected in a given term.

**Byzantine attack**: A Byzantine node can send conflicting vote requests to different nodes, claiming to be a candidate in the same term. This can cause multiple nodes to believe they are the leader, violating election safety.

### 2. Leader Append-Only
**Constraint**: A leader never overwrites or deletes entries in its log.

**Byzantine attack**: A Byzantine leader can maliciously modify or delete log entries, violating the append-only property. It can send different log entries to different followers.

### 3. Log Matching
**Constraint**: If two logs contain an entry with the same index and term, then the logs are identical in all preceding entries.

**Byzantine attack**: A Byzantine node can send conflicting log entries to different followers, creating inconsistent logs that violate the log matching property.

### 4. Leader Completeness
**Constraint**: If a log entry is committed in a given term, then it will be present in the logs of all leaders for all higher-numbered terms.

**Byzantine attack**: A Byzantine leader can refuse to replicate committed entries to new followers, violating leader completeness.

### 5. State Machine Safety
**Constraint**: If a server has applied a log entry at a given index to its state machine, no other server will ever apply a different log entry for the same index.

**Byzantine attack**: A Byzantine node can apply different state machine commands for the same log index, violating state machine safety.

**Key insight**: Raft assumes non-Byzantine failures and is vulnerable to Byzantine attacks. Real-world systems using Raft need additional mechanisms (like digital signatures) to handle Byzantine faults.

## Question 7: Spanner TrueTime Error Bounds Analysis

**Context**: Understanding how time synchronization affects distributed system correctness is crucial for designing systems that rely on time-based ordering. This question explores the impact of TrueTime error bounds on Spanner's performance and correctness.

**Question**: In Spanner, explain what would happen to the system performance/correctness if the error bound with true time is either zero or infinite.

**Detailed Hints:**

### Understanding TrueTime Error Bounds
- **Think about**: What does TrueTime error bound represent?
- **Consider**: How does Spanner use TrueTime for ordering?
- **Key insight**: TrueTime error bounds affect Spanner's ability to order transactions
- **Analysis approach**: Consider the implications of perfect vs. no time synchronization

### Step-by-step Analysis

#### Step 1: Understand TrueTime in Spanner
**Detailed Hints:**
- **Think about**: How does Spanner use TrueTime?
- **Consider**: What role does time play in Spanner's consistency model?
- **Key insight**: Spanner uses TrueTime to provide external consistency

#### Step 2: Analyze Error Bound Scenarios
**Detailed Hints:**
- **Think about**: What happens with zero error bound?
- **Consider**: What happens with infinite error bound?
- **Key insight**: Error bounds affect Spanner's ability to determine transaction ordering

### Detailed Answer

**Impact of TrueTime error bounds on Spanner:**

### Zero Error Bound (Perfect Time Synchronization)

**Performance impact:**
- **Optimal performance**: No waiting time required for external consistency
- **Immediate commits**: Transactions can commit immediately without waiting
- **Maximum throughput**: System achieves maximum possible throughput

**Correctness impact:**
- **Perfect external consistency**: Transactions are ordered exactly according to real time
- **No consistency violations**: All consistency guarantees are maintained perfectly
- **Ideal behavior**: System behaves as if all operations happen instantaneously

**Why this is ideal but unrealistic:**
- Perfect time synchronization is impossible in practice
- Network delays and clock drift make zero error bound unachievable
- Real systems must handle some degree of time uncertainty

### Infinite Error Bound (No Time Synchronization)

**Performance impact:**
- **Severe performance degradation**: System must wait indefinitely for external consistency
- **Blocking behavior**: Transactions cannot commit until error bound is reduced
- **System unavailability**: System effectively becomes unavailable

**Correctness impact:**
- **Consistency violations**: External consistency cannot be guaranteed
- **Ordering problems**: Transactions cannot be properly ordered
- **System failure**: Spanner's core consistency guarantees are violated

**Why this breaks the system:**
- Spanner relies on TrueTime for external consistency
- Without time bounds, the system cannot determine transaction ordering
- The system cannot provide its core consistency guarantees

**Key insight**: TrueTime error bounds are crucial for Spanner's operation. The system needs small but non-zero error bounds to provide both performance and correctness guarantees.

**Key Learning Points:**
- Paxos has theoretical limits on the number of unique values that can be proposed
- Paxos prevents problematic execution sequences through its prepare phase
- Raft is vulnerable to Byzantine attacks without additional security mechanisms
- Time synchronization is crucial for systems that rely on time-based ordering
- Understanding system limits and vulnerabilities is essential for system design