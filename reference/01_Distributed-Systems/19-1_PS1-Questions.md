# Distributed Systems Problem Set 1

## Question 1: True/False Questions

Answer true or false, explain, and (where appropriate) give an example (if true) or counterexample (if false).

### 1a. Exactly-once RPC requires both sender and receiver to write to stable storage

### 1b. Two events can have the same Lamport clock time value but occur minutes apart in real-time

### 1c. Write-back caching never does more communication than write-through caching

### 1d. Two-phase commit with locks released before commit reaches stable storage is serializable

### 1e. Updates in Dynamo are not serializable

### 1f. GFS may contain duplicate records even if application writes each record once

## Question 2: Distributed Systems Features

Consider the following systems: git, Facebook's use of memcache, GFS, BigTable, Spanner, Dynamo, and Bitcoin. For each of the following, find one example of that feature in one of the systems, and sketch its role in the system. A few sentences are sufficient for each example, but use each system as an answer at most once.

### 2g. RPC

### 2h. Caching

### 2i. Eventual Consistency

### 2j. Serializability

### 2k. Logging

### 2l. State Machine Replication

### 2m. Hint

## Question 3: Facebook Three-Tier System

Facebook uses a three-tier system for implementing its website. An array of front-end servers interacts with web clients (each client is hashed into exactly one front-end server); these front-end servers gather the information needed to render the client web page from an array of cache servers and a separate array of storage servers. Hashing is used to locate which cache and storage server might have a particular object (e.g., a friend list, or set of postings). The number of front-end servers, cache servers, and storage servers is not identical (the numbers are chosen to balance the workload), so in general, all front-ends talk to all cache servers and all storage servers.

The cache servers (called memcache servers) are managed as a "lookaside" cache. When rendering an object on a page, the front-end first sends a message to the relevant memcache server; if the data is not available, the front-end (not the cache) then retrieves the data from the relevant storage server. The front-end then stores the fetched data into the memcache server. On update, the front-end invalidates the cached copy (if any) and updates the storage server.

### 3a. What semantics would occur if the front-end first invalidates the cache, and then updates the storage server?

### 3b. What semantics would occur if the front-end updates the storage server and then invalidates the cache?

### 3c. What semantics would occur if the front-end invalidates the cache, updates the storage server, and then re-invalidates the cache?

### 3d. What semantics would occur with the write-token algorithm?

An employee at Facebook suggests adding a write-token to the memcache server. When a front-end wants to change a value, it sends a message to memcache to atomically invalidate the entry and set the write-token; subsequent accesses to the server stall. The front-end releases the write-token when the data is updated at the server, allowing stalled accesses to proceed. What semantics would occur in this algorithm?

## Question 4: Paxos Maximum Values

What is the maximum number of unique values that can be proposed to a group of k Paxos acceptors (for a single instance of the protocol)?

## Question 5: Paxos Sequence Scenario

In Paxos, suppose that the acceptors are A, B, and C. A and B are also proposers, and there is a distinguished learner L. According to the Paxos paper, a value is chosen when a majority of acceptors accept it, and only a single value is chosen. How does Paxos ensure that the following sequence of events cannot happen? What actually happens, and which value is ultimately chosen?

**Given sequence:**
- a) A proposes sequence number 1, and gets responses from A, B, and C.
- b) A sends accept(1, "foo") messages to A and C and gets responses from both. Because a majority accepted, A tells L that "foo" has been chosen. However, A crashes before sending an accept to B.
- c) B proposes sequence number 2, and gets responses from B and C.
- d) B sends accept(2, "bar") messages to B and C and gets responses from both, so B tells L that "bar" has been chosen.

## Question 6: Raft Byzantine Node Attacks

For the Raft algorithm described in the reading list, outline how a Byzantine node would be able to cause each of the correctness constraints to be violated.

## Question 7: Spanner TrueTime Error Bounds

In Spanner, explain what would happen to the system performance/correctness if the error bound with true time is either zero or infinite.