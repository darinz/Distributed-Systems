# Lamport Clocks: Supplementary Notes

## The Problem of Time in Distributed Systems

### Core Challenge
- **Physical time** is fundamental to our thinking
- **Distributed systems** make time ordering complex
- **Key insight**: Sometimes impossible to say which of two events occurred first
- **Result**: "Happened before" relation is only a **partial ordering**

### Why This Matters
- **Specifications** often depend on event ordering
- **Example**: Airline reservation should be granted if made before flight is filled
- **Problem**: In distributed systems, we can't always determine this ordering

### Distributed System Definition
- **Collection of processes** spatially separated
- **Communication** via message exchange
- **Examples**: ARPA network, single computer with separate units
- **Key criterion**: Message transmission delay not negligible compared to inter-event time
## The "Happened Before" Relation

### Definition
The relation "→" on the set of events is the smallest relation satisfying:

1. **Same process**: If a and b are events in the same process, and a comes before b, then a → b
2. **Message sending**: If a is sending a message and b is receiving that message, then a → b  
3. **Transitivity**: If a → b and b → c, then a → c

### Key Properties
- **Irreflexive**: a ↛ a (no event happens before itself)
- **Partial ordering**: Not all events are comparable
- **Concurrent events**: a ↛ b and b ↛ a (neither can causally affect the other)

### Causal Interpretation
- **a → b** means event a can causally affect event b
- **Concurrent events** cannot causally affect each other
- **Example**: Process P cannot know what process Q did until it receives Q's message
## Logical Clocks

### Concept
- **Clock**: Function that assigns a number (timestamp) to each event
- **Logical vs Physical**: No assumption about relation to physical time
- **Implementation**: Can be simple counters with no timing mechanism

### Clock Condition
For any events a, b:
> **If a → b, then C(a) < C(b)**

**Note**: Converse not required (concurrent events can have same timestamp)

### Implementation Rules

#### IR1: Process Clock Increment
- Each process increments its clock between successive events
- Ensures: If a and b are in same process and a comes before b, then C(a) < C(b)

#### IR2: Message Timestamps
- **Sending**: Message contains timestamp Tm = C(a) where a is the send event
- **Receiving**: Process sets C ≥ max(current C, Tm)
- Ensures: If a is sending and b is receiving, then C(a) < C(b)

### Algorithm Summary
```python
# Each process maintains local clock T
def on_event():
    T += 1

def send_message():
    T += 1
    message.timestamp = T
    send(message)

def receive_message(message):
    T = max(T, message.timestamp) + 1
    process(message)
```
## Total Ordering of Events

### Goal
Extend partial ordering to total ordering for practical use

### Method
Define relation "⇒" where a ⇒ b if and only if:
1. **C(a) < C(b)**, OR
2. **C(a) = C(b)** and process(a) < process(b) (arbitrary process ordering)

### Properties
- **Total ordering**: All events are comparable
- **Consistent**: If a → b, then a ⇒ b
- **Breaks ties**: Uses arbitrary process ordering for concurrent events

## Applications

### 1. Primary-Backup Replication
**Problem**: Ensure consistent state across replicas
**Solution**:
- Clients label operations with timestamps
- Send operations to both primary and backup
- Apply events in timestamp order
- Client safe when both acknowledge

### 2. Distributed Make
**Problem**: Determine what needs rebuilding
**Solution**:
- Use timestamps to track file modifications
- If object O depends on source S and O.time < S.time, rebuild O
- **Challenge**: Timestamp correctness across distributed file servers

### 3. Social Media Update Ordering
**Problem**: Ensure consistent view of updates
**Example**: Block boss, then tweet about boss
- Updates sharded across many servers
- Multiple replicas and caches
- **Challenge**: Guarantee no read sees updates in wrong order

### 4. Event Log Merging
**Problem**: Debug distributed systems
**Solution**:
- Each node produces partial event log
- Merge logs using logical timestamps
- Maintain causal ordering for debugging 

## Physical Clocks

### Motivation
- **Logical clocks** don't correspond to physical time
- **User expectations** often based on physical time
- **Anomalous behavior** when logical ordering differs from perceived ordering

### Clock Accuracy Challenges
- **Server clock drift**: ~2 seconds/month
- **Atomic clocks**: Nanosecond accuracy, expensive
- **GPS**: 10ns accuracy, requires antenna
- **Network latency**: Variable, unpredictable

### Synchronization Methods

#### 1. Beacon Approach
- Designate master server with GPS/atomic clock
- Master periodically broadcasts time
- Clients reset clocks on receiving broadcast
- **Problem**: Network latency affects accuracy

#### 2. Client-Driven (NTP, PTP)
- Client queries server
- **Time = server's clock + ½ round trip**
- Average over several servers, throw out outliers
- Adjust for measured clock skew between queries

#### 3. Fine-Grained Physical Clocks
- Timestamps taken in hardware on network interface
- Eliminate samples involving network queueing
- Continually re-estimate clock skew
- Connect all servers in mesh, average all neighbors
- **Accuracy**: ~100ns in worst case

## Mutual Exclusion with Logical Clocks

### Problem
Implement distributed lock using logical clocks

### Assumptions
- In-order point-to-point message delivery
- No failures

### Algorithm

#### Message Types
- **Request**: Broadcast with timestamp
- **Release**: Broadcast with timestamp  
- **Acknowledge**: On receipt of request

#### Node State
- Queue of request messages, ordered by timestamp
- Latest message received from each node

#### Lock Acquisition
1. Send request to everyone (including self)
2. Lock acquired when:
   - My request is at head of my queue, AND
   - I've received same or higher-timestamped messages from everyone

### Key Insight
- **Earliest request wins** due to total ordering
- **Acknowledgment ensures** everyone has seen the request
- **Queue ordering** determines lock ownership

## State Machine Replication

### Generalization
Mutual exclusion generalizes to any state machine:
- **State**: Current state of the system
- **Commands**: Operations that modify state
- **Rule**: Process command iff we've seen all commands with lower timestamp

### Applications
- **Distributed locks**: State = queue of waiting processes
- **Distributed databases**: State = data, commands = transactions
- **Consensus protocols**: State = agreed values, commands = proposals

## Key Takeaways

### Logical vs Physical Time
- **Logical clocks**: Respect causality, no physical meaning
- **Physical clocks**: Approximate real time, subject to drift
- **Choice depends**: On application requirements

### Causal Ordering
- **"Happened before"** captures causality, not time
- **Partial ordering** reflects distributed system reality
- **Total ordering** needed for practical algorithms

### Implementation Principles
- **Local information only**: No global clock required
- **Message timestamps**: Ensure causal ordering across processes
- **Simple rules**: Increment on events, max on message receipt

### Applications
- **Replication**: Consistent state across replicas
- **Synchronization**: Distributed locks and coordination
- **Debugging**: Event log merging and analysis
- **Consistency**: Update ordering in distributed systems

### Limitations
- **No physical time**: Logical clocks don't correspond to real time
- **Concurrent events**: May be ordered arbitrarily
- **Failure handling**: Assumes reliable message delivery
- **Scalability**: Total ordering can be expensive at scale
