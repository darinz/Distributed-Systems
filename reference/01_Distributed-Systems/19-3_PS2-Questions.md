# Problem Set 2: Distributed Systems Timing and Consistency

## Learning Objectives
This problem set focuses on fundamental concepts in distributed systems:
- Understanding causality and concurrency in distributed executions
- Implementing and reasoning about logical clocks and vector clocks
- Analyzing primary-backup replication protocols and their consistency guarantees
- Understanding global state snapshots and their limitations

---

## Question 1: Space-Time Diagrams and Clock Implementations

**Context**: Understanding the temporal relationships between events in a distributed system is crucial for reasoning about consistency, ordering, and causality. This question explores how different clock implementations capture these relationships.

Suppose we have the following space-time diagram describing an execution of a distributed system (time advances downwards). The diagram shows three processes (p1, p2, p3) with labeled events and message passing between them.

<img src="./img/ps_space-time.png" width="200px">

### 1a. Causal Ordering Analysis

For event F, partition the other events (A, B, C, D, E, G, H, and I) into those that happen before F, those that happen after F, and those that are concurrent with F.

**Detailed Hints:**
- **Before F**: Events that causally precede F. Look for events that F could have been influenced by through message chains or local process ordering.
- **After F**: Events that F causally precedes. These are events that could have been influenced by F.
- **Concurrent with F**: Events that are neither before nor after F in the causal ordering. These events could have occurred in any order relative to F without affecting the system's behavior.
- **Key insight**: Two events are concurrent if neither can influence the other through any sequence of message passing or local process execution.

**Step-by-step approach:**
1. Identify all message chains that could lead to F
2. Identify all events that F could influence through message chains
3. Events that don't fall into either category are concurrent with F

### 1b. Logical Clock Implementation

Assume that each process maintains a logical clock. Each clock starts at 0 and is updated at each labeled event, at each message send, and at each message receive. Give the clock value corresponding to each event.

**Given hints**: D has timestamp 1 and G has timestamp 4.

**Detailed Hints:**
- **Logical clock rules**: 
  - Each process increments its clock for every local event
  - When sending a message, include the current clock value
  - When receiving a message, set clock to max(current_clock, received_timestamp) + 1
- **Process identification**: Determine which process each event belongs to by following the vertical lines
- **Message identification**: Arrows represent messages; the timestamp sent is the clock value at the send event
- **Verification**: Use the given hints to check your work - if D=1 and G=4, work backwards and forwards from these known values

**Step-by-step approach:**
1. Assign process identifiers to each vertical line
2. Start with the given hints (D=1, G=4) and work systematically through the diagram
3. For each event, apply the logical clock rules based on whether it's a local event, send, or receive
4. Verify consistency across all events

### 1c. Vector Clock Implementation

Assume instead that each process maintains a vector clock. Give the clock values corresponding to each event.

**Given hint**: G has timestamp {p1 : 0, p2 : 2, p3 : 2}.

**Detailed Hints:**
- **Vector clock rules**:
  - Each process maintains a vector with an entry for each process in the system
  - For local events: increment own entry
  - For message send: include current vector
  - For message receive: set each entry to max(local_vector[i], received_vector[i]), then increment own entry
- **Advantage over logical clocks**: Vector clocks can distinguish between concurrent events that have the same logical timestamp
- **Process identification**: Use the same process identification as in part (b)
- **Verification**: The given hint tells you that at event G, p1 has seen 0 events from itself, p2 has seen 2 events from itself, and p3 has seen 2 events from itself

**Step-by-step approach:**
1. Initialize vector clocks for each process as {p1: 0, p2: 0, p3: 0}
2. Process events in chronological order, applying vector clock rules
3. Use the given hint to verify your calculations
4. Check that concurrent events have incomparable vector timestamps

## Question 2: Primary-Backup Replication Protocol Constraints

**Context**: Primary-backup replication is a fundamental technique for providing fault tolerance in distributed systems. The primary server handles all client requests and maintains a backup server that can take over if the primary fails. However, ensuring consistency between primary and backup while maintaining availability requires careful protocol design with specific constraints.

For each of the following constraints, provide a detailed explanation of why the constraint is necessary and what problems would arise if it were violated. Consider both correctness and consistency implications.

### 2a. State Transfer Metadata Requirement

**Constraint**: State transfer from primary to backup must include metadata on which requests have received replies, and what the response was.

**Detailed Hints:**
- **Think about**: What happens when a backup takes over after a primary failure?
- **Consider**: What information does the backup need to maintain consistency with clients?
- **Key insight**: Clients may have received responses that the backup doesn't know about
- **Problem without constraint**: The backup might re-execute requests that have already been completed, leading to duplicate operations or inconsistent responses
- **Example scenario**: Primary processes a write request, sends response to client, then fails. Backup takes over without knowing this request was completed.

**Step-by-step reasoning:**
1. Identify what happens during normal operation (primary processes requests, sends responses)
2. Consider the failure scenario (primary fails after sending response but before backup knows)
3. Analyze what the backup needs to know to maintain consistency
4. Explain the specific problems that arise without this metadata

### 2b. View Synchronization Requirement

**Constraint**: The backup must accept a request forwarded by the primary if and only if the request and the backup have the same notion of the current view.

**Detailed Hints:**
- **Think about**: What is a "view" in the context of primary-backup systems?
- **Consider**: How do view changes affect the relationship between primary and backup?
- **Key insight**: View changes can occur due to failures, network partitions, or reconfigurations
- **Problem without constraint**: The backup might accept requests from a primary that is no longer the legitimate primary in the current view
- **Example scenario**: View changes due to network partition, but old primary continues operating and forwarding requests to backup

**Step-by-step reasoning:**
1. Define what a "view" represents in the system
2. Identify when and why view changes occur
3. Analyze the relationship between view changes and primary-backup coordination
4. Explain the consistency problems that arise when views are not synchronized

### 2c. Read-Only Request Coordination

**Constraint**: Even on a read-only request, the primary must wait for the backup to accept the request before the primary can reply to the client.

**Detailed Hints:**
- **Think about**: Why would read-only requests need backup coordination?
- **Consider**: What guarantees does the system need to provide about data consistency?
- **Key insight**: Read-only requests still need to reflect a consistent state across primary and backup
- **Problem without constraint**: Clients might read stale data that doesn't reflect the most recent updates
- **Example scenario**: Primary has newer data than backup, but allows reads without backup coordination

**Step-by-step reasoning:**
1. Identify why read-only requests might seem like they don't need coordination
2. Consider the consistency guarantees the system must provide
3. Analyze what happens when primary and backup have different data states
4. Explain the specific problems that arise without this coordination requirement

**Additional considerations for all parts:**
- Think about the trade-offs between consistency and performance
- Consider how these constraints affect system availability
- Analyze the impact on client-perceived system behavior

## Question 3: Global State Analysis in Primary-Backup Systems

**Context**: Understanding what global states are possible in a distributed system is crucial for reasoning about system behavior, debugging, and ensuring correctness. This question explores the difference between consistent global states (which could actually occur in some execution) and inconsistent global states (which represent impossible combinations of local states).

**System Setup**: We have a set of servers, clients, and a view server all running a correct version of the primary/backup protocol. There are exactly two clients, both of which send one command `Append("foo", "x")` and then halt. The network is completely asynchronous (messages can be delayed arbitrarily but are not lost).

### 3a. Consistent Global State Analysis

For each of the following predicates, indicate whether they could be true of a **consistent global state** in any possible execution. A consistent global state is one that could actually occur at some point during a valid execution of the system.

**Detailed Hints for Analysis:**
- **Consistent global state**: A global state is consistent if there exists some execution where all the local states could occur simultaneously
- **Key insight**: Think about the causal relationships between events - if event A causally precedes event B, then A must have occurred before B in any consistent state
- **Primary-backup protocol**: Remember that only one server can be primary at a time, and view changes are coordinated through the view server
- **Asynchronous network**: Messages can be delayed, so local states might be temporarily inconsistent

#### 3a.i. Multiple Primary Servers

**Predicate**: Two different servers report currently being primary.

**Detailed Hints:**
- **Think about**: How does the primary-backup protocol ensure that only one server is primary at a time?
- **Consider**: What role does the view server play in coordinating primary selection?
- **Key insight**: The view server maintains the authoritative view of which server is primary
- **Analysis approach**: Can two servers legitimately believe they are primary simultaneously?

**Step-by-step reasoning:**
1. Understand how primary selection works in the protocol
2. Consider what happens during view changes and failures
3. Analyze whether the protocol allows multiple primaries
4. Determine if this could occur in a consistent state

#### 3a.ii. Backup Accepting Before Primary View Entry

**Predicate**: The backup for view v reports having accepted a request from the primary in view v, while the primary has not yet entered view v (or any later view).

**Detailed Hints:**
- **Think about**: What does it mean for a server to "enter" a view?
- **Consider**: How does the primary-backup protocol coordinate view changes?
- **Key insight**: The primary must be aware of the view before it can forward requests to the backup in that view
- **Analysis approach**: Can a backup accept a request from a primary that doesn't know about the view?

**Step-by-step reasoning:**
1. Understand the view change protocol and how servers enter new views
2. Analyze the sequence of events required for a backup to accept a request
3. Consider the causal relationship between view entry and request forwarding
4. Determine if this scenario violates protocol invariants

#### 3a.iii. Asymmetric Client Response

**Predicate**: One client has received a reply to its command, while the other has not.

**Detailed Hints:**
- **Think about**: What factors could cause one client to receive a response while another doesn't?
- **Consider**: How does the asynchronous network affect message delivery?
- **Key insight**: Network delays and server processing times can cause responses to arrive at different times
- **Analysis approach**: Is it possible for one request to complete while another is still pending?

**Step-by-step reasoning:**
1. Consider the normal request processing flow in the primary-backup protocol
2. Analyze how network asynchrony affects response timing
3. Think about scenarios where one request completes before another
4. Determine if this represents a consistent system state

#### 3a.iv. Both Clients Receive Same Response

**Predicate**: Both clients report receiving `AppendReply("x")`.

**Detailed Hints:**
- **Think about**: What does `AppendReply("x")` mean in the context of the append operation?
- **Consider**: How does the primary-backup protocol handle duplicate requests?
- **Key insight**: The response value indicates what was actually appended to the key
- **Analysis approach**: Can both clients legitimately receive the same response value?

**Step-by-step reasoning:**
1. Understand what the append operation does and what the response represents
2. Analyze how the protocol handles concurrent append operations
3. Consider the consistency guarantees provided by the primary-backup protocol
4. Determine if this response pattern is possible in a correct execution

### 3b. Snapshot-Based Global State Analysis

Now consider a global state gathered by a monitor using the following procedure:
- The monitor node sends a SNAPSHOT message to all other nodes
- Upon receiving SNAPSHOT, each node sends its state to the monitor
- After the monitor receives the states of all nodes, it combines them to form a global state

**Key Difference**: This snapshot procedure captures the state of each node at the moment it receives the SNAPSHOT message, which may not represent a consistent global state.

**Detailed Hints for Snapshot Analysis:**
- **Snapshot timing**: Each node reports its state when it receives the SNAPSHOT message, not at a coordinated time
- **Inconsistent states possible**: The snapshot may capture states that could never occur simultaneously in a real execution
- **Message in transit**: Messages sent before the snapshot but received after it can create inconsistencies
- **Analysis approach**: Consider whether the snapshot timing could create artificial inconsistencies

For each of the following predicates, indicate whether they could be true of a global state gathered in this way.

#### 3b.i. Multiple Primary Servers (Snapshot)

**Predicate**: Two different servers report currently being primary.

**Detailed Hints:**
- **Think about**: How could the snapshot timing create this inconsistency?
- **Consider**: What happens if the SNAPSHOT message arrives at different servers at different times during a view change?
- **Key insight**: Snapshot timing can capture intermediate states during protocol transitions

#### 3b.ii. Backup Accepting Before View Server Acknowledgment (Snapshot)

**Predicate**: The backup for view v reports having accepted a request from the primary in view v, while the view server has not yet received an acknowledgement for view v.

**Detailed Hints:**
- **Think about**: How does the snapshot timing affect the view server's state?
- **Consider**: What happens if the SNAPSHOT arrives at the backup after it accepts the request but before the view server processes the acknowledgment?
- **Key insight**: Snapshot timing can capture states where acknowledgments are in transit

#### 3b.iii. Asymmetric Client Response (Snapshot)

**Predicate**: One client has received a reply to its command, while the other has not.

**Detailed Hints:**
- **Think about**: How could snapshot timing create this pattern?
- **Consider**: What happens if the SNAPSHOT arrives at one client after it receives a response but at the other client before it receives a response?
- **Key insight**: Snapshot timing can capture responses that are in transit

#### 3b.iv. Both Clients Receive Same Response (Snapshot)

**Predicate**: Both clients report receiving `AppendReply("x")`.

**Detailed Hints:**
- **Think about**: How does snapshot timing affect this analysis?
- **Consider**: Could the snapshot capture both clients after they've received responses?
- **Key insight**: Snapshot timing doesn't change the fundamental consistency of the responses themselves

**General Analysis Strategy for Part (b):**
1. For each predicate, consider how the snapshot timing could create the described state
2. Think about message delivery timing and how it affects what each node reports
3. Distinguish between states that are impossible due to protocol invariants vs. states that are impossible due to causal relationships
4. Remember that snapshots can capture inconsistent states that would never occur in a real execution