# Problem Set 4: Distributed Storage and Consensus Systems

## Learning Objectives
This problem set focuses on advanced distributed systems concepts:
- Understanding register semantics and their implementation in distributed systems
- Analyzing Byzantine fault tolerance mechanisms and their security properties
- Exploring distributed storage systems (BigTable, GFS) and their client interactions
- Understanding cryptocurrency protocols and their incentive mechanisms
- Evaluating protocol optimizations and their correctness implications

---

## Question 1: ABD Multi-Writer, Multi-Reader Register Algorithm Analysis

**Context**: The ABD (Attiya, Bar-Noy, Dolev) algorithm is a fundamental protocol for implementing atomic registers in distributed systems. Understanding how to optimize this protocol while maintaining correctness is crucial for building efficient distributed storage systems.

Answer the following questions about the ABD multi-writer, multi-reader register algorithm, presented in Algorithm 1.

### 1a. Regular Semantics Optimization

**Question**: Recall Lamport's safe, regular, and atomic register semantics. These were defined in terms of single-write registers, but there is a natural way to extend these definitions to multiple writers.

The ABD algorithm guarantees atomicity (linearizability). In this algorithm, both reads and writes consist of two phases, which for the purposes of this question we'll refer to as the query phase and the store phase.

In the multi-write case, regularity implies that for some real-time respecting ordering of writes, reads can return either the value written by the most recently completed write (in the sequential ordering) or any of the values currently being written.

Suppose we only wanted to guarantee regular semantics as defined above. How could the protocol be made more efficient?

**Detailed Hints:**

### Understanding Register Semantics
- **Think about**: What are the differences between safe, regular, and atomic register semantics?
- **Consider**: What guarantees does each semantics provide to clients?
- **Key insight**: Regular semantics is weaker than atomic semantics, allowing more flexibility in read responses
- **Analysis approach**: Identify which operations can be optimized under regular semantics

### Step-by-step Analysis

#### Step 1: Understand the ABD Protocol
**Detailed Hints:**
- **Think about**: What does each phase of the ABD protocol accomplish?
- **Consider**: Why does the protocol need both query and store phases?
- **Key insight**: The query phase gathers information, the store phase ensures consistency

**ABD Protocol Overview:**
- **Query phase**: Client sends QUERY to all servers, waits for majority responses
- **Store phase**: Client sends STORE with new timestamp/value to all servers, waits for majority responses
- **Purpose**: Ensures atomicity by coordinating with majority of servers

#### Step 2: Analyze Regular vs Atomic Semantics
**Detailed Hints:**
- **Think about**: What does regular semantics allow that atomic semantics doesn't?
- **Consider**: How does this affect the protocol's requirements?
- **Key insight**: Regular semantics allows reads to return values from concurrent writes

**Semantics Comparison:**
- **Atomic**: Reads must return values that appear in a linearizable execution
- **Regular**: Reads can return the most recent completed write OR any currently being written
- **Implication**: Regular semantics allows more flexibility in read responses

#### Step 3: Identify Optimization Opportunities
**Detailed Hints:**
- **Think about**: Which phase could be optimized under regular semantics?
- **Consider**: What information does a read operation actually need?
- **Key insight**: Regular semantics might allow skipping the store phase for reads

**Optimization Strategy:**
- **For reads**: Only need the query phase to get the most recent value
- **For writes**: Still need both phases to ensure consistency
- **Key insight**: Regular semantics allows reads to return values without ensuring they're "stored"

### Detailed Answer

**How the protocol could be made more efficient:**

1. **Read optimization**: For read operations, skip the store phase entirely
   - Only perform the query phase to get the most recent value
   - Return the value immediately after receiving majority responses
   - This reduces read latency and network traffic

2. **Write optimization**: For write operations, still perform both phases
   - Query phase to get the current timestamp
   - Store phase to ensure the new value is stored at majority of servers
   - This maintains consistency for writes

3. **Network efficiency**: Reduce the number of messages per read operation
   - Original: 2n messages per read (n queries + n stores)
   - Optimized: n messages per read (n queries only)
   - This significantly reduces network overhead

### 1b. Atomic Semantics with Optimization

**Question**: Now suppose that we do want atomic register semantics, as in the original protocol. What additional checks could you add to the ABD protocol to take advantage of the optimization you made in part (a) in the "common" case?

**Detailed Hints:**

### Understanding the Challenge
- **Think about**: How can we use the regular semantics optimization while maintaining atomicity?
- **Consider**: What additional information do we need to ensure atomicity?
- **Key insight**: We need to detect when the optimization is safe to use

### Step-by-step Analysis

#### Step 1: Identify When Optimization is Safe
**Detailed Hints:**
- **Think about**: Under what conditions can we skip the store phase for reads?
- **Consider**: What guarantees do we need to maintain atomicity?
- **Key insight**: We can skip the store phase if we're confident the value is already stored

#### Step 2: Design Additional Checks
**Detailed Hints:**
- **Think about**: What information can we use to make this decision?
- **Consider**: How can we detect if a value is already stored at a majority?
- **Key insight**: We can check if the value we read is already stored at a majority

### Detailed Answer

**Additional checks for atomic semantics:**

1. **Majority storage check**: After the query phase, check if the value we read is already stored at a majority of servers
   - If yes, we can skip the store phase (optimization applies)
   - If no, we must perform the store phase to ensure atomicity

2. **Timestamp validation**: Verify that the timestamp we read is the highest among all servers
   - This ensures we're reading the most recent value
   - If not, we need to perform the store phase

3. **Concurrent write detection**: Check if there are any concurrent writes in progress
   - If no concurrent writes, the optimization is safe
   - If concurrent writes exist, we need the store phase for atomicity

### 1c. Speculative Write Analysis

**Question**: Suppose the ABD algorithm is modified such that a write is sent speculatively in the query phase (with a new timestamp) and that the store phase is skipped if the speculative write succeeded at a majority. The modified algorithm is shown in Algorithm 2. Does this protocol still provide atomicity? If it does, briefly explain why. If not, provide a counter-example trace that demonstrates the problem.

**Detailed Hints:**

### Understanding the Modification
- **Think about**: What does the speculative write modification do?
- **Consider**: How does this change the protocol's behavior?
- **Key insight**: The modification tries to optimize writes by doing them speculatively

### Step-by-step Analysis

#### Step 1: Analyze the Modified Protocol
**Detailed Hints:**
- **Think about**: What happens in the modified query phase?
- **Consider**: How does the speculative write work?
- **Key insight**: The query phase now includes a speculative write attempt

#### Step 2: Identify Potential Issues
**Detailed Hints:**
- **Think about**: What could go wrong with speculative writes?
- **Consider**: How does this affect the protocol's correctness?
- **Key insight**: Speculative writes might violate atomicity guarantees

### Detailed Answer

**The modified protocol does NOT provide atomicity.**

**Counter-example trace:**

1. **Initial state**: All servers have value v1 with timestamp (1, P1)

2. **Client A starts write(v2)**:
   - Sends QUERY with speculative write (2, P2, v2) to all servers
   - Servers S1, S2 accept the speculative write
   - Server S3's message is delayed

3. **Client B starts read()**:
   - Sends QUERY to all servers
   - Gets responses: S1 says (2, P2, v2), S2 says (2, P2, v2), S3 says (1, P1, v1)
   - Client B returns v2 (from majority)

4. **Client A's speculative write fails**:
   - S3's delayed message arrives, but A has already moved to store phase
   - A sends STORE with (3, P2, v2) to all servers
   - All servers accept (3, P2, v2)

5. **Client C starts read()**:
   - Sends QUERY to all servers
   - Gets responses: All servers say (3, P2, v2)
   - Client C returns v2

**Problem**: Client B read v2 before Client A's write was actually completed, violating atomicity. The read should have returned v1 (the previous value) since A's write wasn't complete when B read.

**Key Learning Points:**
- Speculative optimizations can violate correctness guarantees
- Atomicity requires careful coordination between read and write operations
- The ABD protocol's two-phase design is necessary for maintaining atomicity
- Optimizations must preserve the fundamental correctness properties of the protocol
## Question 2: Byzantine Fault Tolerance (BFT) State Machine Replication Analysis

**Context**: Byzantine Fault Tolerance (BFT) protocols are designed to handle arbitrary failures, including malicious behavior. Understanding how BFT algorithms maintain correctness and security in the presence of Byzantine faults is crucial for building secure distributed systems.

Answer the following questions about the BFT state machine replication algorithm presented in the "Practical Byzantine Fault Tolerance" paper.

### 2a. Client Request Spoofing

**Question**: Can a client spoof a request such that it appears to have been initiated by a different client? Briefly justify your answer.

**Detailed Hints:**

### Understanding Client Spoofing
- **Think about**: What mechanisms does BFT use to authenticate client requests?
- **Consider**: How does the protocol ensure that requests come from legitimate clients?
- **Key insight**: BFT protocols typically use cryptographic authentication to prevent spoofing
- **Analysis approach**: Consider what happens if a malicious client tries to impersonate another client

### Step-by-step Analysis

#### Step 1: Understand BFT Client Authentication
**Detailed Hints:**
- **Think about**: How does the BFT protocol authenticate client requests?
- **Consider**: What cryptographic mechanisms are used?
- **Key insight**: BFT protocols use digital signatures or MACs to authenticate requests

#### Step 2: Analyze Spoofing Attempts
**Detailed Hints:**
- **Think about**: What would a malicious client need to do to spoof a request?
- **Consider**: What information does the malicious client need?
- **Key insight**: The malicious client would need the private key of the target client

### Detailed Answer

**No, a client cannot spoof a request to appear as if it was initiated by a different client.**

**Justification:**

1. **Cryptographic authentication**: BFT protocols use digital signatures or message authentication codes (MACs) to authenticate client requests
2. **Private key requirement**: To spoof a request, a malicious client would need the private key of the target client
3. **Key security**: If the target client's private key is compromised, then the client itself is compromised, not just spoofed
4. **Protocol design**: The BFT protocol is designed to prevent such attacks through proper cryptographic authentication

### 2b. Faulty Leader Ignoring Requests

**Question**: Consider a faulty leader that ignores client requests. How does the algorithm make progress as long as there are a sufficient number (2f + 1) of non-faulty replicas (and sufficient network synchrony)?

**Detailed Hints:**

### Understanding Leader Failures
- **Think about**: What happens when the leader fails or becomes unresponsive?
- **Consider**: How does the BFT protocol handle leader failures?
- **Key insight**: BFT protocols have mechanisms to detect and replace faulty leaders
- **Analysis approach**: Consider the view change mechanism in BFT

### Step-by-step Analysis

#### Step 1: Understand View Change Mechanism
**Detailed Hints:**
- **Think about**: How does BFT detect that the leader is faulty?
- **Consider**: What triggers a view change?
- **Key insight**: Timeouts and lack of progress trigger view changes

#### Step 2: Analyze Progress Guarantees
**Detailed Hints:**
- **Think about**: How does the protocol ensure progress after a view change?
- **Consider**: What guarantees does the protocol provide?
- **Key insight**: The protocol ensures that a non-faulty leader will eventually be elected

### Detailed Answer

**The algorithm makes progress through the view change mechanism:**

1. **Timeout detection**: Non-faulty replicas detect that the leader is ignoring requests through timeouts
2. **View change initiation**: Replicas initiate a view change to elect a new leader
3. **New leader election**: The protocol elects a new leader from the remaining non-faulty replicas
4. **Progress guarantee**: As long as there are 2f + 1 non-faulty replicas, a non-faulty leader will eventually be elected
5. **Request processing**: The new leader processes pending client requests and makes progress

### 2c. Faulty Leader Sending Different Pre-prepare Messages

**Question**: Consider a faulty leader that sends different pre-prepare messages to different nodes for the same slot, i.e., it tells some nodes that it has assigned a client command c1 to a given slot and tells other nodes that it has assigned a client command c2 to the same slot. How does the algorithm deal with this issue?

**Detailed Hints:**

### Understanding Pre-prepare Inconsistency
- **Think about**: What happens when the leader sends conflicting pre-prepare messages?
- **Consider**: How do replicas detect this inconsistency?
- **Key insight**: Replicas can detect inconsistencies by comparing pre-prepare messages
- **Analysis approach**: Consider how the prepare phase handles inconsistencies

### Step-by-step Analysis

#### Step 1: Understand Pre-prepare Phase
**Detailed Hints:**
- **Think about**: What does the pre-prepare phase accomplish?
- **Consider**: How do replicas validate pre-prepare messages?
- **Key insight**: Replicas check that pre-prepare messages are consistent

#### Step 2: Analyze Inconsistency Detection
**Detailed Hints:**
- **Think about**: How do replicas detect conflicting pre-prepare messages?
- **Consider**: What happens when inconsistencies are detected?
- **Key insight**: Replicas can detect inconsistencies and trigger view changes

### Detailed Answer

**The algorithm deals with this issue through inconsistency detection and view change:**

1. **Inconsistency detection**: Replicas detect that the leader sent different pre-prepare messages for the same slot
2. **Prepare phase failure**: The prepare phase fails because replicas cannot agree on the same command for the slot
3. **View change trigger**: The inconsistency triggers a view change to elect a new leader
4. **New leader election**: A new leader is elected from the remaining non-faulty replicas
5. **Consistent processing**: The new leader processes requests consistently across all replicas

### 2d. Faulty Replica Sending Incorrect Prepare Message

**Question**: Assume that you have a faulty replica that received a pre-prepare message assigning sequence number n to client request m. What happens if the replica sends a prepare message associating a different sequence number n' to the client request?

**Detailed Hints:**

### Understanding Prepare Message Inconsistency
- **Think about**: What happens when a replica sends an incorrect prepare message?
- **Consider**: How do other replicas detect this inconsistency?
- **Key insight**: Replicas can detect inconsistencies by comparing prepare messages
- **Analysis approach**: Consider how the protocol handles Byzantine behavior

### Step-by-step Analysis

#### Step 1: Understand Prepare Phase
**Detailed Hints:**
- **Think about**: What does the prepare phase accomplish?
- **Consider**: How do replicas validate prepare messages?
- **Key insight**: Replicas check that prepare messages are consistent with pre-prepare messages

#### Step 2: Analyze Inconsistency Handling
**Detailed Hints:**
- **Think about**: How do replicas handle inconsistent prepare messages?
- **Consider**: What happens when inconsistencies are detected?
- **Key insight**: The protocol can tolerate up to f Byzantine faults

### Detailed Answer

**The algorithm handles this through Byzantine fault tolerance:**

1. **Inconsistency detection**: Other replicas detect that the faulty replica sent an incorrect prepare message
2. **Majority requirement**: The protocol requires 2f + 1 replicas, so at most f can be Byzantine
3. **Correct replicas**: The remaining f + 1 non-faulty replicas send correct prepare messages
4. **Consensus**: The protocol reaches consensus based on the majority of correct messages
5. **Faulty replica isolation**: The faulty replica's incorrect message is ignored by the protocol

### 2e. Faulty Replica Sending Incorrect Response

**Question**: Assume that you have a faulty replica. What happens when the replica sends an incorrect response to the client?

**Detailed Hints:**

### Understanding Client Response Handling
- **Think about**: How does the client handle responses from replicas?
- **Consider**: What mechanisms ensure that clients receive correct responses?
- **Key insight**: Clients typically wait for multiple responses and use majority voting
- **Analysis approach**: Consider how clients can detect and handle incorrect responses

### Step-by-step Analysis

#### Step 1: Understand Client Response Processing
**Detailed Hints:**
- **Think about**: How does the client process responses from replicas?
- **Consider**: What validation does the client perform?
- **Key insight**: Clients can detect inconsistencies in responses

#### Step 2: Analyze Faulty Response Handling
**Detailed Hints:**
- **Think about**: How does the client handle incorrect responses?
- **Consider**: What happens when responses are inconsistent?
- **Key insight**: The client can retry the request or use majority voting

### Detailed Answer

**The algorithm handles this through client-side validation:**

1. **Response collection**: The client collects responses from multiple replicas
2. **Inconsistency detection**: The client detects that the faulty replica sent an incorrect response
3. **Majority voting**: The client uses majority voting to determine the correct response
4. **Faulty replica isolation**: The client ignores the incorrect response from the faulty replica
5. **Correct response**: The client accepts the response that matches the majority

### 2f. Authenticator-based Checkpoint Analysis

**Question**: The paper describes an alternative method to message authentication using what they call authenticators. First, each pair of servers sets up a shared key at the beginning of the execution. A message authentication code (MAC) for a message sent from s1 to s2 is a string generated from the message text and the shared key between s1 and s2 that proves that the sender of the message knew the shared key. An authenticator for a message broadcast by s1 is a vector of MACs, one for each of the message's recipients.

Suppose that checkpoint messages contained authenticators, rather than digital signatures, and that servers only wait for f + 1 matching checkpoints before garbage collection. (f + 1 messages were sufficient when digital signatures were attached to checkpoint messages.) What could go wrong?

**Detailed Hints:**

### Understanding Authenticator vs Digital Signature
- **Think about**: What are the differences between authenticators and digital signatures?
- **Consider**: How do these differences affect security properties?
- **Key insight**: Authenticators provide different security guarantees than digital signatures
- **Analysis approach**: Consider what attacks are possible with authenticators

### Step-by-step Analysis

#### Step 1: Understand Authenticator Properties
**Detailed Hints:**
- **Think about**: What security properties do authenticators provide?
- **Consider**: How do authenticators differ from digital signatures?
- **Key insight**: Authenticators provide authentication but not non-repudiation

#### Step 2: Analyze Security Implications
**Detailed Hints:**
- **Think about**: What attacks are possible with authenticators?
- **Consider**: How could a Byzantine replica exploit authenticators?
- **Key insight**: Authenticators can be forged by any replica that knows the shared key

### Detailed Answer

**Several problems could occur with authenticator-based checkpoints:**

1. **Authenticator forgery**: Any replica that knows the shared key can forge authenticators, making it impossible to determine which replica actually created the checkpoint
2. **Non-repudiation loss**: Unlike digital signatures, authenticators don't provide non-repudiation - any replica can claim to have created a checkpoint
3. **Byzantine attack**: A Byzantine replica could create fake checkpoints with valid authenticators, causing other replicas to perform garbage collection prematurely
4. **Consensus violation**: The protocol might reach consensus on fake checkpoints, violating the safety properties of the system
5. **State inconsistency**: Premature garbage collection based on fake checkpoints could lead to state inconsistencies across replicas

**Key Learning Points:**
- BFT protocols use cryptographic authentication to prevent spoofing attacks
- View change mechanisms ensure progress even with faulty leaders
- Inconsistency detection and majority voting handle Byzantine behavior
- Authenticators provide different security guarantees than digital signatures
- Protocol design must consider the specific security properties of cryptographic primitives

---

## Algorithm 1: ABD MRMW Atomic Register Algorithm
Server local state:
t ←(0, ⊥) v ←⊥ 1: upon receiving ⟨QUERY⟩
2: Send reply ⟨QUERY-REPLY,t, v⟩ 3: end upon
4: upon receiving ⟨STORE,t′
, v′⟩
5: if t < t′then
6: t ←t′
7: v ←v′
8: end if
9: Send reply ⟨STORE-REPLY⟩ 10: end upon
◃ Current timestamp, initially unique minimum value; lexicographically ordered
◃ Current value, initially special null value
◃ Messaging infrastructure correctly associates messages with replies
◃ Messaging infrastructure correctly associates messages with replies
Client local state:
p ◃ Unique process ID, immutable
11: procedure READ
12: COMMUNICATE(⊥)
13: end procedure
14: procedure WRITE(v)
15: COMMUNICATE(v)
16: end procedure
17: function COMMUNICATE(v)
18: Send ⟨QUERY⟩to all servers
19: Wait for n
2 + 1 replies, stored in R
20: t ←max{m.t : m ∈R} 21: if v = ⊥then
22: v ←m.v : m ∈R ∧m.t = t 23: else
24: t ←(t[0] + 1, p)
25: end if
26: Send ⟨STORE,t, v⟩to all servers
27: Wait for n
2 + 1 replies
28: end function
◃ Maximum timestamp out of replies seen
◃ Value associated with maximum timestamp
Page 2 of 4
CSE 452 – Winter 2020 Problem Set 3 DUE: 11:59pm March 11th
Algorithm 2 Modified MRMW Register Algorithm
Server local state:
t ←(0, ⊥)
v ←⊥
1: upon receiving ⟨QUERY,t′
, v′⟩
2: if v′̸= ⊥∧t < t′then
3: t ←t′
4: v ←v′
5: end if
6: Send reply ⟨QUERY-REPLY,t, v⟩
7: end upon
8: upon receiving ⟨STORE,t′
, v′⟩
9: if t < t′then
10: t ←t′
11: v ←v′
12: end if
13: Send reply ⟨STORE-REPLY⟩
14: end upon
Client local state:
p
t ←(0, ⊥) 15: procedure READ
16: COMMUNICATE(⊥)
17: end procedure
18: procedure WRITE(v)
19: COMMUNICATE(v)
20: end procedure
21: function COMMUNICATE(v)
22: if v ̸= ⊥then
23: t ←(t[0] + 1, p)
24: end if
25: Send ⟨QUERY,t, v⟩to all servers
26: Wait for n
2 + 1 replies, stored in R
27: t′←max{m.t : m ∈R}
28: if v ̸= ⊥∧t′
= t then 29: return
30: else if v ̸= ⊥then
31: t ←(t′[0] + 1, p)
32: t′←t
33: else
34: v ←m.v : m ∈R ∧m.t = t′
35: end if
36: Send ⟨STORE,t′
, v⟩to all servers
37: Wait for n
2 + 1 replies
38: end function
◃ Local timestamp
◃ Speculative write succeeded at a majority
Page 3 of 4
CSE 452 – Winter 2020 Problem Set 3 DUE: 11:59pm March 11th
(e) (2 points) Assume that you have a faulty replica. What happens when the replica sends an
incorrect response to the client?
(f) (6 points) The paper describes an alternative method to message authentication using what
they call authenticators. First, each pair of servers sets up a shared key at the beginning
of the execution. A message authentication code (MAC) for a message sent from s1 to
s2 is a string generated from the message text and the shared key between s1 and s2 that
proves that the sender of the message knew the shared key. An authenticator for a message
broadcast by s1 is a vector of MACs, one for each of the message’s recipients.
Suppose that checkpoint messages contained authenticators, rather than digital signatures,
and that servers only wait for f + 1 matching checkpoints before garbage collection. ( f + 1
messages were sufficient when digital signatures were attached to checkpoint messages.)
What could go wrong?
## Question 3: BigTable Client Request Analysis

**Context**: BigTable is a distributed storage system that provides a sparse, multidimensional, sorted map. Understanding how BigTable clients interact with the system to read data is crucial for understanding distributed storage system design and performance characteristics.

**Question**: Consider a BigTable client whose cache is empty, i.e., it has never talked to a given BigTable deployment before. Enumerate the requests it will make to read a single row with key k from table T.

**Detailed Hints:**

### Understanding BigTable Architecture
- **Think about**: What are the main components of BigTable?
- **Consider**: How does a client locate data in BigTable?
- **Key insight**: BigTable uses a hierarchical structure with tablets, tablet servers, and metadata
- **Analysis approach**: Trace the client's journey from initial request to data retrieval

### Step-by-step Analysis

#### Step 1: Understand BigTable Components
**Detailed Hints:**
- **Think about**: What are the main components of BigTable?
- **Consider**: How is data organized in BigTable?
- **Key insight**: BigTable uses tablets (ranges of rows) stored on tablet servers

**BigTable Components:**
- **Master**: Manages metadata and tablet assignments
- **Tablet Servers**: Store and serve tablet data
- **Tablets**: Ranges of rows stored on tablet servers
- **Metadata**: Information about tablet locations and assignments

#### Step 2: Trace Client Request Flow
**Detailed Hints:**
- **Think about**: What does the client need to know to read a row?
- **Consider**: How does the client find the right tablet server?
- **Key insight**: The client needs to locate the tablet containing the desired row

### Detailed Answer

**The BigTable client will make the following requests to read a single row with key k from table T:**

1. **Initial metadata request**: 
   - Client contacts the BigTable master to get metadata about table T
   - This includes information about tablet locations and assignments

2. **Tablet location request**:
   - Client requests the location of the tablet that contains row key k
   - The master responds with the tablet server that hosts the relevant tablet

3. **Tablet server contact**:
   - Client contacts the identified tablet server directly
   - Requests the specific row with key k from the tablet

4. **Data retrieval**:
   - The tablet server returns the requested row data
   - Client caches the tablet location for future requests

**Additional considerations:**
- If the tablet server is unavailable, the client may need to retry with the master
- The client may need to handle tablet splits or reassignments
- Caching reduces the number of metadata requests for subsequent reads

**Key Learning Points:**
- BigTable uses a hierarchical structure for data organization
- Clients need to locate the correct tablet server before reading data
- Metadata requests are necessary for initial data location
- Caching improves performance for subsequent requests
## Question 4: GFS Record Append Duplicate Analysis

**Context**: GFS (Google File System) provides a "record append" operation that allows multiple clients to append data to the same file concurrently. Understanding when and why duplicate records can occur is crucial for understanding the trade-offs between performance and consistency in distributed storage systems.

**Question**: Give one scenario where GFS's "record append" would insert duplicate records at the end of a file.

**Detailed Hints:**

### Understanding GFS Record Append
- **Think about**: How does GFS's record append operation work?
- **Consider**: What happens when multiple clients append to the same file?
- **Key insight**: GFS prioritizes performance over strict consistency
- **Analysis approach**: Consider what happens when network issues or failures occur

### Step-by-step Analysis

#### Step 1: Understand GFS Record Append Operation
**Detailed Hints:**
- **Think about**: What does the record append operation do?
- **Consider**: How does GFS handle concurrent appends?
- **Key insight**: GFS allows multiple clients to append to the same file simultaneously

**GFS Record Append Characteristics:**
- **Concurrent appends**: Multiple clients can append to the same file
- **Atomicity**: Each append is atomic at the record level
- **Consistency**: GFS provides eventual consistency, not strong consistency
- **Performance**: Optimized for high throughput and low latency

#### Step 2: Identify Duplicate Scenarios
**Detailed Hints:**
- **Think about**: What could cause duplicate records?
- **Consider**: How does GFS handle network failures or timeouts?
- **Key insight**: Network issues can cause clients to retry operations

### Detailed Answer

**Scenario: Network timeout and client retry**

1. **Initial append request**: Client A sends a record append request to GFS
2. **Network timeout**: The request times out due to network issues
3. **Client retry**: Client A retries the append operation, sending the same data again
4. **First request succeeds**: The original request eventually reaches GFS and is processed
5. **Second request succeeds**: The retry request also reaches GFS and is processed
6. **Duplicate records**: The same record appears twice at the end of the file

**Detailed sequence of events:**

1. **Client A initiates append**: Client A sends record "data123" to GFS
2. **Network delay**: The request is delayed in the network
3. **Client timeout**: Client A times out waiting for a response
4. **Client retry**: Client A sends the same record "data123" again
5. **First request arrives**: The original request reaches GFS and appends "data123"
6. **Second request arrives**: The retry request reaches GFS and appends "data123" again
7. **Result**: The file now contains "data123" twice

**Why this happens:**
- GFS prioritizes performance over strict consistency
- The system doesn't prevent duplicate appends from the same client
- Network issues can cause clients to retry operations
- GFS doesn't maintain state about which records have already been appended

**Key Learning Points:**
- GFS's record append is designed for high performance, not strict consistency
- Network issues can cause duplicate records in distributed systems
- Clients may retry operations when they don't receive responses
- Understanding system trade-offs is crucial for application design
## Question 5: Bitcoin Protocol Security and Incentive Analysis

**Context**: Bitcoin is a decentralized cryptocurrency that relies on cryptographic proofs and economic incentives to maintain security and consensus. Understanding the security properties and incentive mechanisms of Bitcoin is crucial for understanding how decentralized systems can achieve consensus without a central authority.

Answer the following questions about the Bitcoin protocol.

### 5a. 51% Attack Analysis

**Question**: Suppose that an attacker gains control of > 50% of the compute power currently operating on the Bitcoin network. Could this attacker convince non-faulty observers that improperly signed transactions are actually valid?

**Detailed Hints:**

### Understanding 51% Attacks
- **Think about**: What can an attacker with majority compute power do?
- **Consider**: How does Bitcoin's consensus mechanism work?
- **Key insight**: Majority compute power allows control over block creation
- **Analysis approach**: Consider what the attacker can and cannot do

### Step-by-step Analysis

#### Step 1: Understand Bitcoin Consensus
**Detailed Hints:**
- **Think about**: How does Bitcoin reach consensus on transactions?
- **Consider**: What role does compute power play in consensus?
- **Key insight**: Bitcoin uses proof-of-work for consensus

#### Step 2: Analyze Attack Capabilities
**Detailed Hints:**
- **Think about**: What can an attacker with majority compute power do?
- **Consider**: What are the limitations of such an attack?
- **Key insight**: The attacker can control block creation but not transaction validation

### Detailed Answer

**No, the attacker cannot convince non-faulty observers that improperly signed transactions are valid.**

**Reasoning:**

1. **Transaction validation**: Bitcoin nodes validate transactions independently using cryptographic signatures
2. **Signature verification**: Improperly signed transactions will be rejected by all honest nodes
3. **Consensus limitation**: While the attacker can control block creation, they cannot change the validation rules
4. **Network rejection**: Honest nodes will reject blocks containing invalid transactions
5. **Fork creation**: The attacker's invalid blocks will create a fork that honest nodes won't follow

**What the attacker CAN do:**
- Double-spend transactions
- Censor transactions
- Reorganize the blockchain
- Control block creation timing

**What the attacker CANNOT do:**
- Create valid transactions without proper signatures
- Change the validation rules
- Force honest nodes to accept invalid transactions

### 5b. Difficulty Adjustment Analysis

**Question**: Suppose the difficulty were decreased dramatically so that the average rate of block discovery was much higher than once every ten minutes. Further suppose that clients proportionally increased the number of blocks deep they waited for a transaction to be considered confirmed. What could go wrong?

**Detailed Hints:**

### Understanding Difficulty Adjustment
- **Think about**: What happens when block discovery is much faster?
- **Consider**: How does this affect network synchronization?
- **Key insight**: Faster block creation can cause network issues
- **Analysis approach**: Consider the implications of rapid block creation

### Step-by-step Analysis

#### Step 1: Understand Block Creation Rate
**Detailed Hints:**
- **Think about**: What happens when blocks are created very quickly?
- **Consider**: How does this affect network propagation?
- **Key insight**: Faster block creation can cause synchronization issues

#### Step 2: Analyze Network Implications
**Detailed Hints:**
- **Think about**: What problems can arise from rapid block creation?
- **Consider**: How does this affect consensus and security?
- **Key insight**: Network delays can cause forks and consensus issues

### Detailed Answer

**Several problems could occur:**

1. **Network synchronization issues**: 
   - Blocks may not propagate fast enough across the network
   - Some nodes may not receive new blocks before the next block is created
   - This can cause temporary forks and consensus issues

2. **Increased fork probability**:
   - Faster block creation increases the chance of multiple blocks being created simultaneously
   - Network delays can cause honest miners to work on different chains
   - This reduces the security of the consensus mechanism

3. **Resource consumption**:
   - Faster block creation increases bandwidth and storage requirements
   - Nodes need to process and store blocks more frequently
   - This can make running a full node more expensive

4. **Security implications**:
   - Shorter block times reduce the security of the proof-of-work mechanism
   - Attackers need less time to perform certain attacks
   - The network becomes more vulnerable to various attack vectors

### 5c. Miner Incentive Analysis

**Question**: Miners refuse to add invalid transactions (including transactions whose inputs have already been spent) to their proposed block. What incentive does a miner have to follow this procedure? Why doesn't it just add the transaction to the block and let consumers of the transaction log ignore invalid transactions at playback time?

**Detailed Hints:**

### Understanding Miner Incentives
- **Think about**: What motivates miners to follow the rules?
- **Consider**: What happens if miners include invalid transactions?
- **Key insight**: Miners have economic incentives to maintain network integrity
- **Analysis approach**: Consider the economic consequences of including invalid transactions

### Step-by-step Analysis

#### Step 1: Understand Miner Economics
**Detailed Hints:**
- **Think about**: How do miners earn rewards?
- **Consider**: What happens if their blocks are rejected?
- **Key insight**: Miners only earn rewards if their blocks are accepted by the network

#### Step 2: Analyze Invalid Transaction Consequences
**Detailed Hints:**
- **Think about**: What happens if a miner includes invalid transactions?
- **Consider**: How does this affect block acceptance?
- **Key insight**: Invalid transactions cause blocks to be rejected

### Detailed Answer

**Miners have strong economic incentives to follow the rules:**

1. **Block rejection**: If a miner includes invalid transactions, their block will be rejected by the network
2. **Lost rewards**: Rejected blocks don't earn the miner any rewards (block reward + transaction fees)
3. **Wasted resources**: The miner has wasted computational power and electricity on a rejected block
4. **Reputation damage**: Consistently including invalid transactions can damage a miner's reputation
5. **Network participation**: Miners need the network to accept their blocks to earn rewards

**Why miners don't include invalid transactions:**

1. **Economic loss**: Including invalid transactions guarantees that the block will be rejected
2. **No benefit**: There's no economic benefit to including invalid transactions
3. **Network integrity**: Maintaining network integrity is in the miner's long-term interest
4. **Competition**: Other miners will create valid blocks and earn the rewards instead

**Key Learning Points:**
- Bitcoin's security relies on economic incentives, not just cryptographic proofs
- Miners are motivated by profit to maintain network integrity
- Invalid transactions cause blocks to be rejected, resulting in lost rewards
- The protocol design aligns miner incentives with network security
- Understanding economic incentives is crucial for designing secure decentralized systems
Page 4 of 4