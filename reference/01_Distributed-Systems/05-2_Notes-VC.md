# Vector Clocks & Distributed Snapshots: Supplementary Notes

## Vector Clocks

### Motivation
- **Lamport clocks limitation**: T(A) < T(B) doesn't imply A → B
- **Need precise causality**: T(A) < T(B) ⟷ happens-before(A, B)
- **Solution**: Vector clocks track causal relationships precisely

### Concept
- **Track events known to each node** on each node
- **Clock is a vector C** with length = number of nodes
- **Used in practice**: Git, Amazon Dynamo, eventual/causal consistency

### Algorithm

#### On Node i:
1. **On each event**: Increment C[i]
2. **On message receipt** with clock Cm:
   - Increment C[i]
   - For each j ≠ i: C[j] = max(C[j], Cm[j])

### Comparison Rules

#### Happens-Before
- **Cx happens before Cy** if Cx[i] ≤ Cy[i] for all i
- **Precise causality**: Captures all transitive relationships

#### Concurrency
- **Cx and Cy are concurrent** if:
  - Vectors are not identical, AND
  - Cx[i] < Cy[i] and Cx[j] > Cy[j] for some i, j
- **Neither happens before the other**

### Example
```
Node 1: [1,0,0] → [2,0,0] → [3,2,0]
Node 2: [0,1,0] → [0,2,0] → [0,3,0]
Node 3: [0,0,1] → [0,0,2] → [0,0,3]
```

## Distributed System Properties

### Key Terms

#### States and Executions
- **State**: Global state S of system (states at all nodes + channels)
- **Execution**: Series of states Si where system can transition from Si to S(i+1)
- **Reachability**: State Sj is reachable from Si if possible to reach Sj starting from Si

#### Property Types
- **Stable property P**: P(Si) → P(Si+1) (once true, stays true)
- **Invariant P**: Holds on all reachable states

### Token Conservation Example
- **Initial state S0**: Node 1 has token, Node 2 doesn't, no messages
- **Operations**: Send token, discard token
- **Invariant**: Token in at most one place
- **Stable property**: No token

## Distributed Snapshots

### Motivation
- **Detect stable properties**: e.g., deadlock detection
- **Distributed garbage collection**: Identify unreachable objects
- **Diagnostics**: Verify invariants at runtime
- **Record global state**: State of every node and channel

### Challenges
- **Physical clock skew**: Can't trust global time
- **No instantaneous snapshot**: State changes continuously
- **Consistency requirement**: Must capture causal dependencies

## Consistent Snapshots

### Definition
- **Consistent global state**: Causal dependencies are captured
- **If snapshot includes event e2**, then all events that happen before e2 must also be included
- **Corresponds to consistent cut** in space-time diagram

### Consistent Cuts
- **Cut C**: Subset of global history H
- **Consistent cut**: If e2 is in cut and e1 → e2, then e1 is also in cut
- **Respects causality**: No "future" events without their causes

### Inconsistent Cuts
- **Problem**: Include effect without cause
- **Example**: Include message receipt without message send
- **Violates causality**: Leads to impossible states

## Physical Time Algorithm (Ideal)

### Concept
- **If clocks were perfect**: Take snapshot at specific time
- **All nodes record state** at same physical time
- **Handle channels**: Messages in transit

### Channel State Recording
- **Timestamp all messages**
- **Receiver records channel state**
- **Channel state**: Messages received after snapshot time but sent before snapshot time

### Example: Token Count
- **Question**: Is there ≤ 1 token in system?
- **Record**: Node states + messages in transit
- **Count**: Total tokens across all states

### Reality Check
- **Physical clocks aren't accurate enough**
- **Need message coordination** for consistency
- **Must ensure**: Node 2 snapshots before receiving messages sent after Node 1 snapshots

## Chandy-Lamport Snapshot Algorithm

### Initiation
- **Any node can initiate** snapshot at any time
- **Multiple nodes can initiate** concurrently
- **Initiating node**:
  - Records current state
  - Sends "marker" message on all outgoing channels

### Marker Processing
- **When node receives marker**:
  - Records current state
  - Sends marker on all outgoing channels
  - Records channel state

### Channel State Recording
- **Recorded by receiver** when marker received
- **Channel state**:
  - **Empty**: If this is first marker received
  - **Messages received since snapshot**: Otherwise

### Multiple Initiators
- **Same rules apply**: Send markers on all channels
- **Concurrent snapshots OK**: As long as messages in flight are accounted for
- **Consistency maintained**: If receive marker before initiating, must snapshot

### Intuition
- **All initiators are concurrent**: No causal ordering between them
- **Messages in flight**: Captured in channel states
- **Consistency**: Respects causality through marker propagation

## Key Takeaways

### Vector Clocks vs Lamport Clocks
- **Vector clocks**: Precise causality (T(A) < T(B) ⟷ A → B)
- **Lamport clocks**: Approximate causality (A → B ⟹ T(A) < T(B))
- **Trade-off**: More storage and computation for precision

### Distributed Snapshots
- **Purpose**: Capture consistent global state
- **Challenge**: No global time, continuous state changes
- **Solution**: Message-based coordination (Chandy-Lamport)

### Consistency Requirements
- **Consistent cuts**: Respect causality
- **Channel states**: Capture messages in transit
- **Marker propagation**: Ensures coordination

### Applications
- **Deadlock detection**: Find stable property "system deadlocked"
- **Garbage collection**: Find unreachable objects
- **Debugging**: Verify system invariants
- **Checkpointing**: Save consistent state for recovery

### Design Principles
- **Causality preservation**: Never include effect without cause
- **Message coordination**: Use markers to ensure consistency
- **Concurrent initiation**: Multiple snapshots can run simultaneously
- **Channel state**: Must capture messages in transit

### Limitations
- **Storage overhead**: Vector clocks grow with number of nodes
- **Message overhead**: Markers add communication cost
- **Complexity**: More complex than simple timestamps
- **Scalability**: Performance degrades with large numbers of nodes