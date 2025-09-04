# Safety, Liveness, and Consistency: Supplementary Notes

## System Specification Framework

### Core Concepts
- **Execution**: Sequence of events (steps taken by system), potentially infinite
- **Property**: Predicate on executions
- **Safety property**: Specifies "bad things" that shouldn't happen in any execution
- **Liveness property**: Specifies "good things" that should happen in every execution

### Key Theorem
> **Every property is expressible as the conjunction of a safety property and a liveness property**

### Example Properties
- **Safety**: "The system never deadlocks"
- **Liveness**: "Every client that sends a request eventually gets a reply"
- **Both**: "Both generals attack simultaneously"

## Consensus Properties

### Problem Definition
- **N processes**, each with input value from some domain
- **Output**: Processes call `decide(v)`
- **Fault tolerance**: f faulty processes, others continue correctly

### Required Properties

#### 1. Agreement
- **No two correct processes** decide different values
- **Safety property**: Prevents disagreement

#### 2. Integrity
- **Every correct process** decides at most one value
- **If process decides v**, some process had v as input
- **Safety property**: Prevents invalid decisions

#### 3. Termination
- **Every correct process** eventually decides a value
- **Liveness property**: Ensures progress

## Consistency Models

### Definition
- **Consistency**: Allowed semantics (return values) of operations to data store/shared object
- **Specifies interface**, not implementation
- **Anomaly**: Violation of consistency semantics

### Terminology
- **Strong consistency**: System behaves as if single copy of data
- **Weak consistency**: Allows behaviors significantly different from single store model
- **Eventual consistency**: Aberrant behaviors are only temporary

### Why Different Models?

#### Performance
- **Consistency requires synchronization** when data is replicated
- **Often slower** to ensure correct answers

#### Availability
- **What if client offline** or network not working?
- **Weak/eventual consistency** may be only option

#### Programmability
- **Weaker models** are harder to reason about

## Lamport's Register Semantics

### Single-Writer Registers
- **Operations**: Write and read
- **Semantics defined** in terms of real-time beginnings and ends of operations

### Consistency Levels

#### Safe
- **Read not concurrent** with any write obtains previously written value
- **Weakest guarantee**

#### Regular
- **Safe +** read that overlaps write obtains either old or new value
- **No intermediate values**

#### Atomic
- **Safe +** reads and writes behave as if they occur in some definite order
- **Strongest guarantee**

## Sequential Consistency

### Definition
- **Applies to arbitrary shared objects**
- **History equivalent** to legal sequential history
- **Legal sequential history**: Respects local ordering at each node
- **Called serializability** when applied to transactions

### Key Property
- **Operations can appear out of real-time order**
- **Most common**: Read-only operations return stale data

### Example Analysis
```
Process 1: W(x,1) → R(x) → R(y)
Process 2: W(y,1) → R(y) → R(x)
```
**Question**: Is this sequentially consistent?

## Linearizability

### Definition
- **Linearizability = Sequential consistency + respects real-time ordering**
- **If e1 ends before e2 begins**, then e1 appears before e2 in sequential history
- **Behaves as if single, correct copy** of data
- **Atomic registers are linearizable**

### Real-Time Ordering
- **Stronger than sequential consistency**
- **Must respect actual time** of operations
- **No stale reads** allowed

### Example Analysis
**Question**: Is this linearizable?

## Causal Consistency

### Definition
- **Writes not concurrent** (related by happens-before) must be seen in that order
- **Concurrent writes** can be seen in different orders on different nodes
- **Linearizability implies causal consistency**

### Key Insight
> **Causal consistency is the strongest form of consistency that can be provided in an always-available convergent system**

**Translation**: If you want to process writes even with network partitions and failures, causal consistency is the best you can do.

## Weaker Consistency Models

### FIFO Consistency
- **Writes by same process** seen in that order
- **Writes to different processes** can be seen in different orders
- **Equivalent to PRAM model**

### Eventual Consistency
- **If writes to object stop**, eventually all processes read same value
- **Not even a safety property**: "Eventual consistency is no consistency"
- **Only guarantees convergence**, not ordering

## Consistency Hierarchy

### From Strongest to Weakest
1. **Linearizability**: Real-time ordering + sequential consistency
2. **Sequential Consistency**: Legal sequential history
3. **Causal Consistency**: Happens-before ordering
4. **FIFO Consistency**: Per-process ordering
5. **Eventual Consistency**: Eventual convergence

### Safety Properties
- **All except eventual consistency** are safety properties
- **Eventual consistency** is a liveness property

## Practical Considerations

### Using Consistency Guarantees

#### Example Program
```python
# Initially a = 0, b = 0
Process 1: a = 1; print(b)
Process 2: b = 1; print(a)
```

**Question**: What are possible outputs?

#### Analysis
- **If both print 0**: Cycle in happens-before graph → Not sequential!
- **Possible outputs**: (0,1), (1,0), (1,1)

### Java's Memory Model
- **Java is not sequentially consistent!**
- **Guarantees sequential consistency** only when program is data-race free
- **Data-race**: Two threads access same memory location concurrently, one is write, not protected by locks

### Using Weak Consistency
- **Separate operations** with stronger semantics
- **Weak consistency** (and high performance) by default
- **Application-level protocols** using:
  - Separate communication
  - Extra synchronization variables in data store

## Key Takeaways

### Trade-offs
- **Weaker consistency**: Harder to program against, easier to enforce
- **Stronger consistency**: Easier to program against, harder to enforce

### Design Principles
- **Choose consistency model** based on application requirements
- **Consider performance vs correctness** trade-offs
- **Plan for availability** in face of failures
- **Use stronger consistency** only when necessary

### Real-World Applications
- **Strong consistency**: Banking, critical systems
- **Causal consistency**: Social media, collaborative editing
- **Eventual consistency**: DNS, caching systems
- **FIFO consistency**: Message queues, logs

### Implementation Challenges
- **Strong consistency**: Requires coordination, can impact performance
- **Weak consistency**: Requires careful application design
- **Network partitions**: May force weaker consistency models
- **Scalability**: Stronger models harder to scale