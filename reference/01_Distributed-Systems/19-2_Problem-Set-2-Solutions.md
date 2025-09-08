
# Distributed Systems Problem Set

## Question 1: True/False Questions

Answer true or false, explain, and (where appropriate) give an example (if true) or counterexample (if false).
### 1a. Exactly-once RPC requires both sender and receiver to write to stable storage

**Answer: True**

**Explanation:** In an asynchronous network with arbitrary message delays, exactly-once semantics requires both the sender and receiver to maintain state in stable storage to handle:
- **Sender side:** Must track which requests have been sent and their responses to handle retransmissions and detect duplicates
- **Receiver side:** Must track which requests have been processed to avoid duplicate execution

**Example:** If a client sends a "transfer $100" RPC and the network delays the response, the client might retransmit. Without stable storage on both sides, the receiver might process the duplicate request, causing the transfer to happen twice.

### 1b. Two events can have the same Lamport clock time value but occur minutes apart in real-time

**Answer: True**

**Explanation:** Lamport clocks only capture causal ordering, not real-time ordering. Events that are not causally related can have the same Lamport timestamp even if they occur far apart in real-time.

**Example:** Consider two processes that never communicate with each other. Process A increments its counter at time T, and Process B increments its counter at time T+5 minutes. Both events will have the same Lamport timestamp (assuming they started with the same initial value) even though they occurred 5 minutes apart.

### 1c. Write-back caching never does more communication than write-through caching

**Answer: False**

**Explanation:** While write-back caching typically reduces communication by batching writes, there are scenarios where it can require more communication than write-through caching.

**Counterexample:** In a system with frequent cache evictions, write-back caching might need to:
1. Write dirty data back to storage when evicting
2. Then fetch the new data from storage
This results in two network operations compared to write-through caching's single write operation. Additionally, write-back caching requires more complex cache coherence protocols that can generate additional communication overhead.

### 1d. Two-phase commit with locks released before commit reaches stable storage is serializable

**Answer: False**

**Explanation:** Releasing locks before the commit reaches stable storage violates the durability property of ACID transactions and can lead to non-serializable behavior.

**Counterexample:** Consider two transactions T1 and T2:
1. T1 acquires lock on item X, modifies it, releases lock, but crashes before commit reaches stable storage
2. T2 acquires lock on item X, reads the uncommitted value, commits successfully
3. T1's changes are lost, but T2 has already committed based on those changes
This creates a non-serializable execution where T2's commit depends on T1's uncommitted changes.

### 1e. Updates in Dynamo are not serializable

**Answer: True**

**Explanation:** Dynamo uses eventual consistency and allows concurrent writes to the same key, which can result in conflicts that are resolved using vector clocks or last-write-wins semantics, not serializable ordering.

**Example:** Two clients simultaneously update the same shopping cart item. Client A sets quantity to 5, Client B sets quantity to 3. Dynamo will store both versions and require application-level conflict resolution, rather than ensuring a serializable ordering of these updates.

### 1f. GFS may contain duplicate records even if application writes each record once

**Answer: True**

**Explanation:** GFS uses append-only semantics and does not guarantee exactly-once delivery. Network retries, client crashes, and other failures can cause the same record to be appended multiple times.

**Example:** A client appends a record to a GFS file, but the network response is delayed. The client times out and retries the append operation. GFS will append the record twice, resulting in duplicate records in the file, even though the application intended to write it only once.
## Question 2: Distributed Systems Features

Consider the following systems: git, Facebook's use of memcache, GFS, BigTable, Spanner, Dynamo, and Bitcoin. For each of the following, find one example of that feature in one of the systems, and sketch its role in the system. A few sentences are sufficient for each example, but use each system as an answer at most once.

### 2g. RPC - **BigTable**

BigTable uses RPC extensively for communication between clients and tablet servers, and between tablet servers and the master. The master uses RPC to assign tablets to servers, monitor server health, and handle load balancing. Tablet servers use RPC to serve read/write requests from clients and to communicate with other tablet servers for tablet splits and merges. This RPC layer provides the abstraction that makes BigTable appear as a single, unified storage system despite being distributed across many machines.

### 2h. Caching - **Facebook's memcache**

Facebook uses memcache as a distributed caching layer between web servers and persistent storage. When rendering a web page, front-end servers first check memcache for user data, friend lists, and posts. If data is not cached, the front-end fetches it from the database and stores it in memcache for future requests. This dramatically reduces database load and improves response times for frequently accessed data. The cache uses a "lookaside" pattern where the application explicitly manages cache invalidation.

### 2i. Eventual Consistency - **Dynamo**

Dynamo provides eventual consistency by allowing multiple versions of the same data to coexist temporarily. When concurrent writes occur, Dynamo stores multiple versions with vector clocks to track causality. The system eventually converges to a consistent state through read-repair (where clients resolve conflicts during reads) and anti-entropy processes. This design prioritizes availability and partition tolerance over strong consistency, making it suitable for applications that can tolerate temporary inconsistencies.

### 2j. Serializability - **Spanner**

Spanner provides serializable transactions across globally distributed data using TrueTime and two-phase commit. When a transaction spans multiple shards, Spanner uses two-phase commit with a coordinator to ensure all-or-nothing semantics. The TrueTime API allows Spanner to assign globally meaningful timestamps to transactions, ensuring that concurrent transactions are ordered consistently across all replicas. This enables applications to maintain ACID properties even when data is distributed across multiple data centers.

### 2k. Logging - **GFS**

GFS uses extensive logging for fault tolerance and recovery. The master maintains operation logs that record all metadata changes, including file creation, deletion, and chunk assignments. These logs are replicated to multiple machines and flushed to stable storage before responding to clients. When the master restarts, it replays the operation log to reconstruct its in-memory state. Additionally, GFS logs all chunk operations for debugging and performance analysis, enabling the system to recover from failures and maintain data consistency.

### 2l. State Machine Replication - **Bitcoin**

Bitcoin uses state machine replication through its blockchain mechanism. All nodes maintain identical copies of the blockchain state (account balances, transaction history). When a new block is mined, it contains a set of transactions that represent state transitions. All nodes validate and apply these transactions in the same order, ensuring that all honest nodes maintain identical state. The proof-of-work consensus mechanism ensures that only valid state transitions are accepted, and the longest chain rule resolves conflicts between competing state transitions.

### 2m. Hint - **Git**

Git uses "hint" mechanisms in several ways to optimize performance. The most notable example is the packfile format, where Git stores hints about object locations to speed up lookups. When Git needs to find an object, it can use these hints to quickly locate the object in the packfile without scanning the entire file. Additionally, Git maintains hints about branch tips and commit relationships to optimize operations like `git log` and `git merge`. These hints allow Git to perform complex operations efficiently even on large repositories with millions of objects.	
## Question 3: Facebook Three-Tier System

Facebook uses a three-tier system for implementing its website. An array of front-end servers interacts with web clients (each client is hashed into exactly one front-end server); these front-end servers gather the information needed to render the client web page from an array of cache servers and a separate array of storage servers. Hashing is used to locate which cache and storage server might have a particular object (e.g., a friend list, or set of postings). The number of front-end servers, cache servers, and storage servers is not identical (the numbers are chosen to balance the workload), so in general, all front-ends talk to all cache servers and all storage servers.

The cache servers (called memcache servers) are managed as a "lookaside" cache. When rendering an object on a page, the front-end first sends a message to the relevant memcache server; if the data is not available, the front-end (not the cache) then retrieves the data from the relevant storage server. The front-end then stores the fetched data into the memcache server. On update, the front-end invalidates the cached copy (if any) and updates the storage server.

### 3a. What semantics would occur if the front-end first invalidates the cache, and then updates the storage server?

**Answer: Inconsistent semantics**

**Explanation:** This approach creates a race condition window where the system can be in an inconsistent state. If another front-end server reads the data between cache invalidation and storage update, it will:
1. Find the cache entry missing (due to invalidation)
2. Read stale data from storage (before the update)
3. Store the stale data back in cache

This results in the cache containing stale data even after the storage has been updated, leading to inconsistent reads until the cache is eventually invalidated again.

### 3b. What semantics would occur if the front-end updates the storage server and then invalidates the cache?

**Answer: Eventual consistency**

**Explanation:** This approach provides eventual consistency. Initially, there may be a period where:
1. Storage contains the new value
2. Cache still contains the old value
3. Reads from cache return stale data

However, once the cache invalidation completes, subsequent reads will fetch the updated data from storage and cache it. The system eventually converges to a consistent state where both storage and cache contain the same updated value. This is eventual consistency because there's a temporary period of inconsistency, but the system eventually becomes consistent.

### 3c. What semantics would occur if the front-end invalidates the cache, updates the storage server, and then re-invalidates the cache?

**Answer: Eventual consistency (same as 3b)**

**Explanation:** The second cache invalidation is redundant and doesn't improve the semantics. The system still exhibits eventual consistency with the same race condition window as in 3b. The second invalidation doesn't eliminate the possibility that another front-end might have read stale data from storage and cached it between the storage update and the second invalidation. The semantics remain eventual consistency, just with additional overhead from the redundant invalidation.

### 3d. What semantics would occur with the write-token algorithm?

**Answer: Serializable semantics**

**Explanation:** The write-token mechanism provides serializable semantics by ensuring mutual exclusion during updates. The algorithm works as follows:
1. When a front-end wants to update, it acquires the write-token (atomically invalidating cache and blocking reads)
2. Other front-ends attempting to read the same data will stall until the token is released
3. The front-end updates storage and then releases the token
4. Stalled reads can now proceed and will get the updated data

This ensures that all operations on the same data are serialized - either a read sees the old value (before the update) or the new value (after the update), but never a partially updated or inconsistent state. The mutual exclusion provided by the write-token guarantees serializable execution.
## Question 4: Paxos Maximum Values

**Question:** What is the maximum number of unique values that can be proposed to a group of k Paxos acceptors (for a single instance of the protocol)? Briefly explain.

**Answer:** The maximum number of unique values that can be proposed is **k** (the number of acceptors).

**Explanation:** In Paxos, each acceptor can accept at most one value per instance of the protocol. The key insight is that once an acceptor accepts a value, it cannot accept a different value for the same instance. This is enforced by the Paxos algorithm's Phase 2 rules:

1. An acceptor can only accept a proposal if it hasn't already accepted a proposal with a higher sequence number
2. Once an acceptor accepts a value, it must reject any subsequent proposals with lower sequence numbers

Since there are k acceptors, and each can accept at most one unique value, the theoretical maximum is k unique values. However, in practice, Paxos is designed to converge to a single chosen value. The algorithm ensures that once a majority of acceptors accept the same value, that value becomes the chosen value, and no other value can be chosen for that instance.

The k unique values scenario would only occur in a pathological case where:
- Each acceptor accepts a different value
- No majority is formed for any single value
- The protocol fails to converge

This is why Paxos requires a majority (⌊k/2⌋ + 1) of acceptors to choose a value, ensuring that at most one value can be chosen per instance.
## Question 5: Paxos Sequence Scenario

**Question:** In Paxos, suppose that the acceptors are A, B, and C. A and B are also proposers, and there is a distinguished learner L. According to the Paxos paper, a value is chosen when a majority of acceptors accept it, and only a single value is chosen. How does Paxos ensure that the following sequence of events cannot happen? What actually happens, and which value is ultimately chosen?

**Given sequence:**
- a) A proposes sequence number 1, and gets responses from A, B, and C.
- b) A sends accept(1, "foo") messages to A and C and gets responses from both. Because a majority accepted, A tells L that "foo" has been chosen. However, A crashes before sending an accept to B.
- c) B proposes sequence number 2, and gets responses from B and C.
- d) B sends accept(2, "bar") messages to B and C and gets responses from both, so B tells L that "bar" has been chosen.

**Answer:** This sequence **cannot happen** in Paxos due to the algorithm's safety guarantees.

**Explanation:** The key issue is in step (c). When B proposes sequence number 2, it must follow the Paxos Phase 1 protocol, which requires B to:

1. **Send prepare(2) messages** to all acceptors (A, B, C)
2. **Wait for responses** from a majority of acceptors
3. **Learn about any previously accepted values** from the responses

**What actually happens:**

In step (c), when B sends prepare(2) messages:
- **Acceptor A:** Will respond with "no promise" (since A has already accepted value "foo" with sequence number 1, and 2 > 1, so A can promise to B)
- **Acceptor B:** Will respond with "no promise" (B hasn't accepted anything yet)
- **Acceptor C:** Will respond with "no promise" (C has already accepted value "foo" with sequence number 1, and 2 > 1, so C can promise to B)

However, when B receives responses from A and C, it will learn that both A and C have already accepted value "foo" with sequence number 1. According to Paxos rules, B must propose the same value "foo" that was already accepted by a majority, not a new value "bar".

**What B actually does:**
- B receives prepare(2) responses from A and C
- B learns that "foo" was already accepted by a majority (A and C)
- B sends accept(2, "foo") messages to A, B, and C (not "bar")
- All acceptors accept "foo" with sequence number 2
- B tells L that "foo" has been chosen

**Final result:** The value "foo" is ultimately chosen, not "bar". Paxos ensures that once a value is chosen by a majority, no other value can be chosen for the same instance, even if the original proposer crashes.
## Question 6: Raft Byzantine Node Attacks

**Question:** For the Raft algorithm described in the reading list, outline how a Byzantine node would be able to cause each of the correctness constraints to be violated.

**Answer:** A Byzantine node can violate Raft's correctness constraints through various malicious behaviors:

### 1. **Election Safety** (At most one leader can be elected in a given term)

**Violation:** A Byzantine node can cause multiple leaders to be elected in the same term.

**Attack:** The Byzantine node can:
- Send conflicting `RequestVote` responses to different candidates
- Vote for multiple candidates in the same term
- Send fake `AppendEntries` messages claiming to be the leader
- This can lead to split-brain scenarios where multiple nodes believe they are the leader

### 2. **Leader Append-Only** (A leader never overwrites or deletes entries in its log)

**Violation:** A Byzantine leader can modify or delete log entries.

**Attack:** The Byzantine leader can:
- Send `AppendEntries` RPCs with modified log entries
- Delete entries from its log and send truncated logs to followers
- Send entries with incorrect terms or indices
- This violates the integrity of the replicated log

### 3. **Log Matching** (If two logs contain an entry with the same index and term, then the logs are identical in all preceding entries)

**Violation:** A Byzantine node can create inconsistent log states.

**Attack:** The Byzantine node can:
- Send `AppendEntries` with mismatched log entries
- Claim to have entries that don't exist in its log
- Send entries with incorrect previous log terms
- This breaks the log matching property and can cause followers to accept inconsistent logs

### 4. **Leader Completeness** (If a log entry is committed in a given term, then that entry will be present in the logs of the leaders for all higher-numbered terms)

**Violation:** A Byzantine node can cause committed entries to be lost.

**Attack:** The Byzantine node can:
- Become a leader and refuse to replicate committed entries
- Send `AppendEntries` that overwrite committed entries
- Claim that committed entries were never committed
- This can cause data loss and violate the durability guarantee

### 5. **State Machine Safety** (If a server has applied a log entry at a given index to its state machine, no other server will ever apply a different log entry at the same index)

**Violation:** A Byzantine node can cause different state machines to apply different entries at the same index.

**Attack:** The Byzantine node can:
- Send different log entries to different followers
- Cause followers to apply entries that were never committed
- Send fake commit messages to some followers but not others
- This leads to inconsistent state machine states across the cluster

### **Key Insight:**
Raft assumes a **crash-fault model** where nodes can fail but cannot behave maliciously. Byzantine nodes can exploit this assumption by:
- Sending arbitrary messages
- Lying about their state
- Coordinating attacks with other Byzantine nodes
- Violating protocol rules while appearing to follow them

This is why Byzantine fault-tolerant consensus algorithms (like PBFT) require more complex mechanisms, including cryptographic signatures, message authentication, and higher fault thresholds (typically 3f+1 nodes to tolerate f Byzantine failures).
## Question 7: Spanner TrueTime Error Bounds

**Question:** In Spanner, explain what would happen to the system performance/correctness if the error bound with true time is either zero or infinite.

**Answer:** The TrueTime error bound (ε) is crucial for Spanner's correctness and performance. Different values have significant implications:

### **Case 1: Error Bound = Zero (ε = 0)**

**What it means:** Perfect clock synchronization with no uncertainty.

**Performance Impact:**
- **Excellent performance:** No waiting time required for commit timestamps
- **Immediate commits:** Transactions can commit immediately without waiting for clock uncertainty to pass
- **High throughput:** No artificial delays, maximum transaction processing speed

**Correctness Impact:**
- **Perfect correctness:** No risk of timestamp ordering violations
- **Ideal serializability:** Global transaction ordering is guaranteed to be correct
- **No conflicts:** No need for retry mechanisms due to timestamp ordering issues

**Reality:** This is impossible to achieve in practice due to:
- Network latency variations
- Clock drift between machines
- Physical limitations of clock synchronization

### **Case 2: Error Bound = Infinite (ε = ∞)**

**What it means:** No clock synchronization, complete uncertainty about time.

**Performance Impact:**
- **Severe performance degradation:** Every transaction must wait for the maximum possible clock uncertainty
- **Very low throughput:** Artificial delays make the system extremely slow
- **Practical unusability:** System becomes too slow for real-world applications

**Correctness Impact:**
- **Correctness maintained:** Spanner can still provide serializability, but at a huge performance cost
- **Conservative approach:** System errs on the side of safety by waiting for maximum uncertainty
- **No timestamp ordering violations:** But at the expense of making the system practically unusable

**Example:** If ε = ∞, every commit might need to wait for hours or days to ensure no timestamp ordering violations, making the system unusable.

### **Spanner's Actual Approach: ε ≈ 1-10ms**

**Performance Impact:**
- **Good performance:** Small, manageable delays (1-10ms) for commit timestamps
- **Reasonable throughput:** System remains fast enough for practical use
- **Balanced trade-off:** Acceptable performance with strong consistency guarantees

**Correctness Impact:**
- **Strong consistency:** Serializability maintained with high confidence
- **Rare conflicts:** Occasional retries needed, but infrequent enough to not impact overall performance
- **Practical correctness:** System provides the consistency guarantees needed for real applications

### **Key Insight:**
The TrueTime error bound represents a fundamental trade-off in distributed systems:
- **Smaller ε:** Better performance, but harder to achieve in practice
- **Larger ε:** Easier to achieve, but worse performance
- **Spanner's choice:** A small but realistic ε that provides both good performance and strong consistency guarantees

This is why Spanner invests heavily in clock synchronization infrastructure (GPS receivers, atomic clocks) to minimize ε while keeping it practical.