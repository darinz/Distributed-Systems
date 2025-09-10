# Dynamo: Amazon's Highly Available Key-value Store

## Introduction: A Production-Scale Distributed System

Dynamo represents one of the most influential papers in distributed systems, describing Amazon's production key-value store that powers critical e-commerce operations. This paper is particularly valuable because it describes a real-world system that had to make practical trade-offs between theoretical ideals and operational requirements.

### Learning Goals

The Dynamo paper teaches several crucial concepts:

- **Eventual consistency**: How to build systems that prioritize availability over strong consistency
- **Quorum systems**: How to achieve consistency guarantees with replicated data
- **The importance of tail latency**: Why the 99.9th percentile matters more than average performance
- **Trade-offs in messy real-world systems**: How production systems must balance competing requirements

### Why Dynamo Matters

Dynamo is significant because it:

- **Production validation**: Proves that eventually consistent systems can work at scale
- **Trade-off documentation**: Honestly documents the compromises made for operational needs
- **Influence on industry**: Inspired many subsequent systems including Cassandra, Riak, and others
- **Real-world insights**: Provides insights into building systems that must work in practice, not just theory

## System Design Philosophy

### The Simple Key-Value Store

**Dynamo is a simple key-value store, used in a production environment**

**Amazon's applications do not need the full power of SQL queries**

**But all else equal, traditional consistency would be fine...**

Dynamo's design philosophy is rooted in simplicity and practical needs:

#### Why Key-Value?

**Amazon's applications do not need the full power of SQL queries**

Amazon chose a key-value store because:

- **Application requirements**: Most Amazon applications only need simple read/write operations
- **Performance**: Key-value operations are much faster than complex SQL queries
- **Simplicity**: Simpler data model means simpler implementation and fewer failure modes
- **Scalability**: Key-value stores are easier to scale horizontally

#### The Consistency Trade-off

**But all else equal, traditional consistency would be fine...**

This statement reveals a crucial insight:

- **Ideal world**: In a perfect world, strong consistency would be preferred
- **Real-world constraints**: But real-world systems face constraints that make strong consistency impractical
- **Trade-off necessity**: Amazon had to choose between consistency and other important properties
- **Pragmatic decision**: The choice was made based on business needs, not theoretical purity

### System Goals

**What are the system goals?**

**Incremental scalability**

**Symmetry**

**Decentralization**

**Heterogeneity**

Dynamo was designed with four key goals in mind:

#### Incremental Scalability

**Incremental scalability**

The system must be able to grow gradually:

- **Add nodes**: New nodes can be added without disrupting the entire system
- **Remove nodes**: Nodes can be removed without data loss
- **Linear scaling**: Performance should scale linearly with the number of nodes
- **No downtime**: Scaling operations should not require system downtime

#### Symmetry

**Symmetry**

All nodes should be equal:

- **No special nodes**: No node should have special responsibilities or privileges
- **Uniform software**: All nodes run the same software
- **Load distribution**: Work should be distributed evenly across all nodes
- **Failure handling**: Any node can fail without special recovery procedures

#### Decentralization

**Decentralization**

The system should have no central points of failure:

- **No master nodes**: No single node controls the entire system
- **Distributed decisions**: Decisions are made collectively by the nodes
- **Fault tolerance**: The system should continue operating even if multiple nodes fail
- **Self-healing**: The system should automatically recover from failures

#### Heterogeneity

**Heterogeneity**

The system must work with diverse hardware:

- **Different capabilities**: Nodes may have different CPU, memory, and storage capacities
- **Mixed generations**: New and old hardware must work together
- **Load balancing**: The system should account for different node capabilities
- **Graceful degradation**: The system should work even with heterogeneous performance
## The ACID vs. CAP Trade-off

### Traditional ACID Semantics

**What are traditional "ACID" database semantics (Â§2.1)?**

**Atomicity - All or none of actions happen**

**Consistency - Transaction leaves data in valid state (all invariants hold)**

**Isolation - All actions in a transaction appear to happen before or after other transactions**

**Durability - Effects of transactions will survive reboots**

Traditional databases provide ACID guarantees:

#### Atomicity

**Atomicity - All or none of actions happen**

- **All-or-nothing**: Either all operations in a transaction succeed, or none do
- **Failure handling**: If any operation fails, the entire transaction is rolled back
- **Consistency guarantee**: The database never enters an inconsistent state
- **Example**: Transferring money between accounts - either both accounts are updated or neither

#### Consistency

**Consistency - Transaction leaves data in valid state (all invariants hold)**

- **Invariant preservation**: All database invariants are maintained
- **Valid state**: The database is always in a valid state
- **Constraint enforcement**: All constraints (foreign keys, check constraints, etc.) are enforced
- **Example**: Account balance cannot go negative

#### Isolation

**Isolation - All actions in a transaction appear to happen before or after other transactions**

- **Concurrent execution**: Multiple transactions can run concurrently
- **Serializable execution**: The result is as if transactions ran one after another
- **No interference**: Transactions don't interfere with each other
- **Example**: Two users updating the same account balance

#### Durability

**Durability - Effects of transactions will survive reboots**

- **Persistent storage**: Committed changes are written to persistent storage
- **Crash recovery**: Changes survive system crashes and reboots
- **WAL (Write-Ahead Logging)**: Changes are logged before being applied
- **Example**: Money transfer is not lost if the system crashes

### Why Amazon Abandoned ACID

**Why doesn't Amazon want to pursue traditional ACID semantics?**

**Problem: Consistency vs. availability and response time (Â§2.3)**

**In many scenarios, actually better to get wrong answer than no/slow answer**

**Exacerbating problem is focus on 99.9% latency (a.k.a. TAIL LATENCY)--why?**

**Exceeding latency costs Amazon serious money (blocks purchase)**

**A single end-user operation may require many Dynamo operations**

**Any one of them could blow the whole operation and lose a customer**

**Strict SLAs ensure such failures will be kept to a minimum**

Amazon's decision to abandon ACID was driven by business requirements:

#### The CAP Theorem in Practice

**Problem: Consistency vs. availability and response time (Â§2.3)**

The CAP theorem states that in a distributed system, you can only guarantee two of:
- **Consistency**: All nodes see the same data at the same time
- **Availability**: The system remains operational
- **Partition tolerance**: The system continues to work despite network partitions

Amazon chose **Availability** and **Partition tolerance** over **Consistency**.

#### The Business Case for Availability

**In many scenarios, actually better to get wrong answer than no/slow answer**

This is a crucial business insight:

- **Customer experience**: A slow or unavailable system loses customers
- **Revenue impact**: Every lost customer represents lost revenue
- **Competitive advantage**: Fast, available systems provide competitive advantage
- **User expectations**: Users expect fast responses, even if occasionally inconsistent

#### The Tail Latency Problem

**Exacerbating problem is focus on 99.9% latency (a.k.a. TAIL LATENCY)--why?**

**Exceeding latency costs Amazon serious money (blocks purchase)**

**A single end-user operation may require many Dynamo operations**

**Any one of them could blow the whole operation and lose a customer**

**Strict SLAs ensure such failures will be kept to a minimum**

Tail latency is critical for Amazon:

- **99.9th percentile**: The slowest 0.1% of requests matter most
- **Cascading failures**: One slow operation can block an entire user session
- **Revenue impact**: Slow operations directly impact revenue
- **SLA requirements**: Strict service level agreements must be met
- **Customer retention**: Fast systems keep customers coming back

## Dynamo API Design

### The Simple Interface

**What does Dynamo API look like (Â§4.1)?**

**get (key) -> (context, list of values)**

**put (key, context, value) -> void**

Dynamo provides a remarkably simple API with just two operations:

#### Get Operation

**get (key) -> (context, list of values)**

The get operation:
- **Input**: A key to retrieve
- **Output**: A context and a list of values
- **Multiple values**: Can return multiple values if there are conflicts
- **Context**: Contains metadata needed for future updates

#### Put Operation

**put (key, context, value) -> void**

The put operation:
- **Input**: A key, context from previous get, and a new value
- **Output**: Void (no return value)
- **Context requirement**: Must provide context from previous get operation
- **Atomic**: Either succeeds completely or fails completely

### The Context Mechanism

**What is this weird context value from the client's point of view**

**Just an opaque string of bytes needing to be sent back in a put**

**But idea is to help server resolve conflicts (more on this later)**

The context is a crucial mechanism for conflict resolution:

#### Client Perspective

**Just an opaque string of bytes needing to be sent back in a put**

From the client's perspective:
- **Opaque data**: The client doesn't need to understand the context
- **Pass-through**: The client simply passes the context from get to put
- **No interpretation**: The client doesn't parse or modify the context
- **Required**: The context must be provided for put operations

#### Server Perspective

**But idea is to help server resolve conflicts (more on this later)**

The server uses the context to:
- **Track versions**: The context contains version information
- **Detect conflicts**: The server can detect when conflicts occur
- **Resolve conflicts**: The server can resolve conflicts using the context
- **Maintain causality**: The context helps maintain causal ordering

## Data Distribution Challenges

### Key Distribution Problems

**What are issues in deciding how to spread keys amongst the servers?**

**Replicate values on at least N machines (in case of failure)**

**Handle heterogeneous nodes (new machines more powerful than old)**

**Minimize churn when machines are added or removed**

Distributing data across servers involves several complex challenges:

#### Replication Requirements

**Replicate values on at least N machines (in case of failure)**

The system must:
- **Fault tolerance**: Store each value on at least N machines to survive failures
- **Load distribution**: Spread the load evenly across all machines
- **Consistency**: Ensure all replicas are kept in sync
- **Recovery**: Handle the case when machines fail and come back online

#### Heterogeneity Challenges

**Handle heterogeneous nodes (new machines more powerful than old)**

The system must account for:
- **Different capabilities**: New machines may be more powerful than old ones
- **Load balancing**: Distribute work based on machine capabilities
- **Capacity planning**: Account for different storage and processing capacities
- **Performance optimization**: Route requests to the most capable machines

#### Membership Changes

**Minimize churn when machines are added or removed**

The system must handle:
- **Node addition**: Adding new machines without disrupting existing operations
- **Node removal**: Removing machines without losing data
- **Data migration**: Moving data between machines efficiently
- **Minimal disruption**: Keep the system operational during changes

## Consistent Hashing vs. Fixed Partitions

### Why Not Consistent Hashing?

**Does dynamo place data using consistent hashing? Not really.**

**Why is consistent hashing not great for Dynamo?**

**Partitions determined when nodes join and leave system**

**So nodes must scan their data to re-partition and transfer state**

**Makes reconciliation harder (must recompute Merkle trees)**

**Makes snapshots harder**

Dynamo initially tried consistent hashing but found it problematic:

#### Dynamic Partitioning Problems

**Partitions determined when nodes join and leave system**

**So nodes must scan their data to re-partition and transfer state**

Consistent hashing creates several issues:
- **Dynamic boundaries**: Partition boundaries change when nodes join/leave
- **Data scanning**: Nodes must scan all their data to determine what to transfer
- **Expensive operations**: Re-partitioning is computationally expensive
- **Disruption**: Membership changes cause significant system disruption

#### Operational Complexity

**Makes reconciliation harder (must recompute Merkle trees)**

**Makes snapshots harder**

Consistent hashing increases operational complexity:
- **Merkle tree recomputation**: Must recompute Merkle trees after re-partitioning
- **Snapshot complexity**: Taking consistent snapshots becomes more difficult
- **Debugging difficulty**: Harder to debug issues with dynamic partitions
- **Monitoring challenges**: More difficult to monitor and manage the system

### Dynamo's Solution: Fixed Partitions

**What does Dynamo actually do instead of consistent hashing (Â§6.2)?**

**Split ring into fixed, equal size arcs/segments (Figure 7 strategy 3)**

**Use many more segments than there are nodes**

**Divide these segments up amongst nodes (each segment replicated N times)**

**This is a good technique to know about (# partitions >> # servers)!**

Dynamo uses a much simpler approach:

#### Fixed Partition Strategy

**Split ring into fixed, equal size arcs/segments (Figure 7 strategy 3)**

**Use many more segments than there are nodes**

**Divide these segments up amongst nodes (each segment replicated N times)**

The fixed partition approach:
- **Fixed boundaries**: Partition boundaries never change
- **Many partitions**: Use many more partitions than nodes
- **Node assignment**: Assign partitions to nodes, not the other way around
- **Replication**: Each partition is replicated N times across different nodes

#### Benefits of Fixed Partitions

**This is a good technique to know about (# partitions >> # servers)!**

Fixed partitions provide several advantages:
- **Simpler operations**: Adding/removing nodes is much simpler
- **Better load balancing**: More partitions allow better load distribution
- **Easier reconciliation**: Merkle trees are easier to maintain
- **Simpler snapshots**: Taking snapshots is more straightforward

### Historical Context

**Note consistent hashing / DHTs were a hot topic in 2007**

**Probably why authors tried it first, even though simpler technique better**

**Maybe this wasn't the best way to write the paper, though**

This reveals an important lesson about system design:
- **Trend following**: Consistent hashing was popular in 2007
- **Pragmatic choice**: The simpler technique was actually better
- **Paper presentation**: The paper might have been clearer if it started with the final solution
- **Engineering wisdom**: Sometimes the simpler solution is the better solution

Fig. 6: What is fairness metric/balance here?  mean node / max load
      out of balance node has load >= 1.15 * mean-load
  Why is this a good metric?  most heavily loaded nodes affect tail latency
  Why does imbalance go up as load goes down--is this a problem?
    The more randomly assigned requests, the more even the load
    But under low load doesn't matter, because no latency anyway

Explain Figure 8 (p. 217)?  (#1 is consistent hashing, #3 fixed buckets)
  S nodes in system, T tokens per node in #1, Q keyspace partitions in #3
  First, what is state that must be stored (for 1-hop lookups)?
    #1,2: Need S*T token -> node mappings (every token of every node)
    #3: Need to store Q token (partition) -> node mappings
  For same fairness, #1 uses "3 orders of magnitude" more state than #3.  Why?
    In #3, the tokens are spaced perfectly evenly, will even out load
    Tokens might also be smaller (make least significant bits all 0, truncate)
  Any disadvantages to #3?
    "changing the node membership requires coordination" (p. 217)
  Why is #2 so much worse?  Quantizes unfairness, which should exacerbate
    What is advantage of #2 over #1?  State reconciliation still easier

How does Dynamo achieve geographic replication?
  Â§4.6: "In essence, the preference list of a key is constructed such
         that the storage nodes are spread across multiple data centers."
  Maybe they ensure this while doing partition assignment to servers?
    Reserve even partitions for one data center, odd partitions for another?
  Might be clearer if paper not written in terms of consistent hashing

## Quorum Systems: The Foundation of Consistency

### Understanding Quorum Techniques

**What is a "quorum technique"?**

**Say you store a value on N servers and care about linearizability...**

**Write must reach some write-set of W servers before deemed complete**

**Read must hear from some read-set of R servers before being deemed complete**

**If every write-set and read-set have a non-empty intersection:**

**Guarantees a reader will see effects of a previously completed write**

**An easy way to guarantee this is to ensure R + W > N**

**With N fixed nodes, read will always hear one copy of successful write**

**In Dynamo, each instance configured with its own N, R, W**

Quorum systems are the mathematical foundation that allows Dynamo to provide consistency guarantees:

#### The Basic Quorum Model

**Say you store a value on N servers and care about linearizability...**

**Write must reach some write-set of W servers before deemed complete**

**Read must hear from some read-set of R servers before being deemed complete**

The quorum model works as follows:
- **N servers**: Store replicas of each value on N servers
- **Write quorum (W)**: A write must reach W servers before being considered complete
- **Read quorum (R)**: A read must hear from R servers before being considered complete
- **Overlap requirement**: W and R must be chosen so that they always overlap

#### The Intersection Guarantee

**If every write-set and read-set have a non-empty intersection:**

**Guarantees a reader will see effects of a previously completed write**

This is the key insight:
- **Non-empty intersection**: Every write quorum and read quorum must share at least one server
- **Consistency guarantee**: This ensures that reads will see the effects of completed writes
- **Linearizability**: This provides the linearizability guarantee
- **Mathematical foundation**: This is the mathematical basis for consistency

#### The Simple Rule

**An easy way to guarantee this is to ensure R + W > N**

**With N fixed nodes, read will always hear one copy of successful write**

The simple rule for ensuring intersection:
- **R + W > N**: This guarantees that any read quorum and write quorum will overlap
- **Proof**: If R + W > N, then R + W - N > 0, meaning there's at least one server in common
- **Practical implication**: The read will always hear from at least one server that participated in the write
- **Consistency**: This ensures that completed writes are visible to subsequent reads

#### Dynamo's Configuration

**In Dynamo, each instance configured with its own N, R, W**

Dynamo allows flexible configuration:
- **Per-instance configuration**: Each Dynamo instance can have different N, R, W values
- **Application-specific tuning**: Different applications can choose different consistency levels
- **Trade-off flexibility**: Allows trading off consistency for performance or availability
- **Operational control**: Operators can tune these parameters based on requirements

In general, quorums useful with many independently accessed objects
  Unlike, e.g., Raft, no need for nodes to agree on complete state of system
    linearizability still possible, because local property of each object
  But to ensure object linearizability, must worry about conflicts
    E.g., imagine 3 nodes, A, B, C storing replicated register V with R=W=2
      Client 1 writes x to A; simultaneously client 2 writes y to B; C fails
      Now you can't read V, because no majority
  Solution?  Assign version numbers to values (a bit like Paxos ballots)
    Put client ID in version numbers to make them unique
    To write:
      - Ask R replicas for current version number
      - Pick new version number greater than any previous version
      - Send new version to all replicas
      - Replicas accept only if new version number higher than previous
      - Write completes at client only if/when W replicas acknowledge success
    To read:
      - Ask all replicas for current (value, version) pair
      - Receive R matching replies?  Done
      - Otherwise, clean up the mess:  Re-broadcast most recent
        (value, version) to finish write started by failed/slow client.

Digression: Can a quorum system survive f Byzantine failures? [Malkhi & Reiter]
      https://link.springer.com/article/10.1007/s004460050050
  To write:
    - Wait for R version numbers, pick one greater than all seen/sent
    - Broadcast write, wait for W acknowledgments
  To read:
    - Wait for R read replies *such that f+1 are matching*
      (otherwise, if only f, could all be lies)
    - If multiple f+1 matching replies, take one with highest version number
    - What if don't get f+1 matching reads?
      Might have concurrent writes, so just ask again
        (Could rebroadcast failed client writes yourself before retry)
  How many servers do you need for f failures?
    - Safety:  R + W - N >= 2f + 1
      When read follows write, minimum overlap of two quorums is R + W - N
      That overlap must contain at least f+1 honest nodes for read to succeed
    - Liveness:  R <= N - f  (a read quorum must exist)
    - Also:  W <= N (otherwise doesn't make any sense)
  Note this works if N == 4f+1 and R == W == 3f+1
    Safety: (3f+1) + (3f+1) - (4f+1) = 2f+1
    Liveness: 2f+1 <= 4f+1 - f

What happens when Dynamo client issues a request (Â§4.5)?
  A key's N successors around Chord ring are its *preference list* (Â§4.3)
  Option 1:  Send to generic load balancer, which sends to Dynamo server
  Option 2:  Clients have server map and know where to send (Â§6.4)
    Clients poll random server every 10 seconds for membership change
  Then at Dynamo server:
    If not in preference list, forward to (ideally first) in pref list
      Will usually only be necessary for Option 1 above
      Forwarding is required for writes, optional for reads
        (Why required for writes?  To keep vector timestamps small)
    Wait for W or R servers (depending on op), then reply
    Do "syntactic" read repair if necessary (Â§5)
    Keep waiting for stragglers to update them and simplify anti-entropy
  Which of Options 1 & 2 is better?  Client-driven (Table 2, p. 218)  Why?
    Avoids one hop through load balancer
      Also avoids relaying writes to coordinator node in preference list
    Why is difference more pronounced for 99.9%ile latency than average?
      Could get unlucky at load balancer or at dynamo node (cf. IPFS gateway)

What happens when a server is down or partitioned (Â§4.6)?
  Hinted handoff:  Write values to next server around ring
    E.g., in Fig. 2 N=3, if B down, store node K on E
    But tell E that the write was really intended for B
  Server (E) stores hinted values in separate file,
    Makes it easy to transfer whole file to target (B) when back up
    (Note with partitioning scheme #3 the separate file thing automatic)
  What does this do to our quorum algorithm?!!
    Basically breaks it--they call this a "Sloppy quorum"
    But fits with philosophy of prioritizing availability over correctness

What happens if server is permanently down?
  Administrator must manually decide this and configure system
  Fortunately, new writes have been going to the right place
    But may now have old values replicated at only N-1 servers
  Administrator may also add new node(s) to replace old
After add/remove/failure--how to ensure state properly replicated
  Transfer needed ring arcs from other servers
  Confirmation round (Â§4.9) ensures you don't get something you don't need
In many cases, don't need to transfer a whole arc
  Just want to fix up a few missing entries
  Merkle trees make it easy to compare state when only small parts are missing

How does every node know the composition of the whole Chord ring?
  Only human administrator can change ring membership
  Gossip protocol (Â§4.8) - every sec exchange membership info with random peer
When adding servers, what's a seed?  (Â§4.8.2)
  Server known to all--want to make sure you don't get divergent parallel rings
  Does this break symmetry?  Maybe mildly
    All servers run the same software, just some may be listed in, e.g., DNS

## Conflict Resolution in Dynamo

### Is Dynamo Linearizable?

**Is Dynamo Linearizable? No. API even exposes conflicts**

Dynamo explicitly chooses not to be linearizable:
- **Eventual consistency**: Dynamo provides eventual consistency, not strong consistency
- **Conflict exposure**: The API exposes conflicts to the application
- **Application responsibility**: Applications must handle conflicts themselves
- **Trade-off choice**: This is a deliberate trade-off for availability and performance

### Sources of Conflicts

**What can lead to an update conflict?**

**Servers partitioned, multiple coordinators**

**Concurrent client updates**

**Unlike usual quorum system, write coordinator doesn't read versions**

**Client does because of "context," but can have concurrent clients**

Conflicts can arise from several sources:

#### Network Partitions

**Servers partitioned, multiple coordinators**

When servers are partitioned:
- **Multiple coordinators**: Different partitions may have different coordinators
- **Independent updates**: Each partition can process updates independently
- **Reconciliation needed**: When partitions heal, conflicts must be resolved
- **Split-brain scenario**: This is a classic split-brain problem

#### Concurrent Updates

**Concurrent client updates**

**Unlike usual quorum system, write coordinator doesn't read versions**

**Client does because of "context," but can have concurrent clients**

Concurrent updates create conflicts:
- **No coordinator coordination**: Unlike traditional quorum systems, coordinators don't coordinate
- **Client-driven versioning**: Clients handle versioning through the context mechanism
- **Concurrent clients**: Multiple clients can update the same key simultaneously
- **Race conditions**: This creates race conditions that can lead to conflicts

### Conflict Resolution Mechanisms

**How are such conflicts resolved?**

**Servers may resolve "syntactically" using vector clocks. Review:**

**List of serverId-versioNo pairs. E.g., <A-1, B-2> (assume 0 for others)**

**Each server bumps its version number when making an update**

**Partial order: vector clock V1 <= V2 if forall s V1[s] <= V2[s]**

**If V1 <= V2, replace version V1 with V2 (cf. git fast-forward merge)**

**If V1 </= V2 and V2 </= V1, have update conflict**

Dynamo uses vector clocks for conflict detection:

#### Vector Clock Mechanics

**List of serverId-versioNo pairs. E.g., <A-1, B-2> (assume 0 for others)**

**Each server bumps its version number when making an update**

Vector clocks work as follows:
- **Server-version pairs**: Each vector clock contains (serverId, version) pairs
- **Version incrementing**: Each server increments its version when making updates
- **Causal tracking**: Vector clocks track causal relationships between updates
- **Partial ordering**: They provide a partial order on updates

#### Conflict Detection

**Partial order: vector clock V1 <= V2 if forall s V1[s] <= V2[s]**

**If V1 <= V2, replace version V1 with V2 (cf. git fast-forward merge)**

**If V1 </= V2 and V2 </= V1, have update conflict**

Conflict detection uses partial ordering:
- **Partial order**: V1 ≤ V2 if every component of V1 is ≤ the corresponding component of V2
- **Fast-forward merge**: If V1 ≤ V2, V1 can be replaced with V2 (like git fast-forward)
- **True conflict**: If neither V1 ≤ V2 nor V2 ≤ V1, there's a true conflict
- **Conflict resolution**: True conflicts require application-level resolution

### Semantic Conflict Resolution

**How does Dynamo handle update conflicts - resolve "semantically"?**

**Option 1: expose to user (that's why get returns multiple values)**

**Requires application willing to resolve conflicts**

**Option 2: just pick one (last writer wins, discards updates)**

Dynamo provides two approaches to semantic conflict resolution:

#### Application-Level Resolution

**Option 1: expose to user (that's why get returns multiple values)**

**Requires application willing to resolve conflicts**

This approach:
- **Multiple values**: The get operation returns multiple conflicting values
- **Application responsibility**: The application must resolve conflicts
- **Flexible resolution**: Applications can implement custom conflict resolution logic
- **Semantic awareness**: Applications understand the semantics of their data

#### Automatic Resolution

**Option 2: just pick one (last writer wins, discards updates)**

This approach:
- **Last writer wins**: Simply pick one value (usually the most recent)
- **Data loss**: This can result in data loss
- **Simple implementation**: Much simpler to implement
- **Acceptable for some use cases**: May be acceptable for certain applications

### Real-World Implications

**Can a "remove from cart" operation be lost?**

**Probably--semantics might just union the carts; cost-benefit ok for Amazon**

**Note Dynamo truncates vector clocks at 10 newest entries--implication?**

**Can make non-conflicting vals conflict <A-1,B-1,C-1> vs. <[A-1,]B-1,C-1,D-1>**

**Could it make a conflict look like a non-conflict? Unlikely**

These design choices have real-world implications:

#### Shopping Cart Example

**Can a "remove from cart" operation be lost?**

**Probably--semantics might just union the carts; cost-benefit ok for Amazon**

For shopping carts:
- **Union semantics**: Conflicting cart updates might be resolved by unioning the carts
- **Lost removals**: A "remove from cart" operation might be lost
- **Business acceptable**: Amazon found this acceptable for their use case
- **Cost-benefit analysis**: The benefits of availability outweighed the costs of occasional data loss

#### Vector Clock Truncation

**Note Dynamo truncates vector clocks at 10 newest entries--implication?**

**Can make non-conflicting vals conflict <A-1,B-1,C-1> vs. <[A-1,]B-1,C-1,D-1>**

**Could it make a conflict look like a non-conflict? Unlikely**

Vector clock truncation:
- **Space efficiency**: Truncation saves space by limiting vector clock size
- **False conflicts**: Can create false conflicts when clocks are truncated
- **No false non-conflicts**: Unlikely to make real conflicts look like non-conflicts
- **Practical trade-off**: This is a practical trade-off for space efficiency

No global failure state (Â§4.8.3)--why?
  Node join and leave operations are explicitly triggered by administrator
  While node is unavailable, other nodes can just detect this for themselves
  Much simpler

How do durability compromises work to decrease latency (Â§6.1)?
  Nodes use background thread to write--reply before data stably on disk
  Coordinator asks one out of N nodes to write synchronously
  But usually W < N, so data won't be on disk before replying to client

What are background tasks (Â§6.5)?
  Replica synchronization and data handoff (from hinting or node churn)
  Use feedback-based mechanism to ensure these don't disrupt foreground put/get
  Feedback is a neat technique to use in system design (PLL analogy)

What types of task might you tune N, R, W for?
  High performance read (product catalog) has R=1, W=N

What are scaling bottlenecks?
  1-hop routing will require more and more state as # nodes S grows (Q ~ S)
  That also increases size of gossiping protocol messages
  Constant reliability demands replication factor N ~ log(S)

## Conclusion: Lessons from Dynamo

### Key Insights

Dynamo represents a landmark in distributed systems design, demonstrating how to build production-scale systems that prioritize availability over strong consistency. The key insights from Dynamo include:

#### Trade-offs Are Inevitable

**Would design be different if they had used C++ instead of Java?**

**W/o garbage collection, maybe could have best node coordinate all updates**

**GC cycles might force sloppier consistency than necessary**

Even implementation language choices affect system design:
- **Garbage collection impact**: Java's GC cycles forced sloppier consistency than necessary
- **C++ advantages**: Without GC, Dynamo could have had better consistency guarantees
- **Engineering trade-offs**: Every design decision involves trade-offs
- **Real-world constraints**: Production systems must work within real-world constraints

#### Availability Over Consistency

Dynamo's most important lesson is that availability often matters more than consistency:
- **Business requirements**: Amazon's business required high availability
- **Customer experience**: Slow or unavailable systems lose customers
- **Revenue impact**: Every lost customer represents lost revenue
- **Pragmatic choice**: The choice was made based on business needs, not theoretical purity

#### Simplicity Wins

Dynamo demonstrates that simpler solutions often work better:
- **Fixed partitions**: Simpler than consistent hashing
- **Simple API**: Just get and put operations
- **Pragmatic choices**: Sometimes the simpler solution is the better solution
- **Engineering wisdom**: Don't over-engineer solutions

### Legacy and Influence

Dynamo has had enormous influence on the distributed systems field:

#### Industry Impact

- **NoSQL movement**: Dynamo inspired the NoSQL movement
- **Eventual consistency**: Made eventual consistency acceptable for production systems
- **Open source systems**: Inspired systems like Cassandra, Riak, and others
- **Cloud computing**: Influenced the design of cloud storage systems

#### Research Impact

- **CAP theorem**: Provided a real-world example of the CAP theorem in practice
- **Quorum systems**: Demonstrated practical applications of quorum systems
- **Conflict resolution**: Advanced the understanding of conflict resolution in distributed systems
- **System design**: Influenced how we think about designing distributed systems

### Modern Relevance

Dynamo's lessons remain relevant today:

#### Microservices Architecture

- **Service independence**: Each service can choose its own consistency model
- **Eventual consistency**: Many microservices use eventual consistency
- **Conflict resolution**: Applications must handle conflicts themselves
- **Trade-off awareness**: Developers must understand the trade-offs involved

#### Cloud Computing

- **Multi-region deployment**: Cloud systems often span multiple regions
- **Availability zones**: Systems must handle failures of availability zones
- **Eventual consistency**: Many cloud services use eventual consistency
- **Conflict resolution**: Cloud applications must handle conflicts

### Final Thoughts

Dynamo represents a masterclass in practical distributed systems design. It shows how to:

1. **Make informed trade-offs**: Choose availability over consistency when business requires it
2. **Keep things simple**: Use the simplest solution that meets requirements
3. **Design for operations**: Consider operational complexity in design decisions
4. **Document honestly**: Be honest about the limitations and trade-offs of your system
5. **Learn from experience**: Use real-world experience to guide design decisions

The paper's lasting value lies not just in its technical contributions, but in its honest documentation of the messy reality of building production distributed systems. It shows that sometimes the "wrong" answer from a theoretical perspective is the right answer from a practical perspective.

Dynamo proves that distributed systems can be both simple and powerful, and that sometimes the best engineering is knowing when to stop engineering and start shipping.