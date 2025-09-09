# Ownership: A Distributed Futures System for Fine-Grained Tasks

## Introduction: The Challenge of Modern Distributed Applications

### The Rise of Large-Scale Distributed Applications

Modern distributed applications have evolved to run across multiple machines to achieve unprecedented performance and scalability. These applications span diverse domains:

- **Model serving**: Real-time inference for machine learning models across multiple servers
- **Online video processing**: Real-time video encoding, transcoding, and streaming
- **Distributed training**: Training large neural networks across GPU clusters
- **Data processing**: Large-scale ETL pipelines and real-time analytics
- **Scientific computing**: High-performance computing simulations and research

### Task and Actor Decomposition

These complex applications are naturally decomposed into smaller, manageable units:

#### Tasks vs. Actors: A Fundamental Distinction

**Tasks** are **stateless, functional units of computation**:
- Execute once and produce a result
- No persistent state between invocations
- Can be executed anywhere in the cluster
- Naturally parallelizable and fault-tolerant

**Actors** are **stateful, long-running entities**:
- Maintain persistent state across multiple interactions
- Execute on specific nodes (stateful)
- Communicate through message passing
- More complex failure recovery due to state

#### The Driver Process: The Orchestrator

A **driver process** serves as the application's coordinator and entry point:
- **Role**: The "head" or root node of the distributed computation
- **Responsibilities**: 
  - Launches the application
  - Coordinates task scheduling
  - Manages the overall execution flow
  - Handles application-level error recovery

### Communication in Distributed Systems

#### Remote Procedure Calls (RPCs)

Tasks and actors communicate primarily through **Remote Procedure Calls (RPCs)**:

**Similarities to local function calls**:
- Same programming model and syntax
- Type-safe parameter passing
- Return value semantics
- Exception handling

**Key differences from local calls**:
- **Network latency**: Orders of magnitude slower than local calls
- **Failure modes**: Network partitions, node failures, timeouts
- **Partial failures**: Remote calls can fail independently
- **Resource management**: Remote resources must be explicitly managed

#### Parameter Passing Strategies

**Pass by value**:
- **Advantage**: Simple, no shared state issues
- **Disadvantage**: Expensive for large objects
- **Use case**: Small objects, immutable data

**Pass by reference**:
- **Advantage**: Efficient for large objects
- **Disadvantage**: Complex lifetime management
- **Use case**: Large objects, mutable shared state

#### Synchronous vs. Asynchronous RPCs

**Synchronous RPCs**:
- **Blocking**: Caller waits for response
- **Simple programming model**: Sequential execution
- **Poor resource utilization**: CPU idle during network wait

**Asynchronous RPCs**:
- **Non-blocking**: Caller continues execution
- **Parallelism**: Multiple operations can overlap
- **Latency hiding**: Computation can overlap with communication
- **Complex programming model**: Requires callback or future-based programming

### The Promise of Distributed Futures

#### What Are Futures and Promises?

**Futures** represent **asynchronous computation results**:
- Placeholder for a value that will be computed later
- Can be passed around, stored, and combined
- Enable composition of asynchronous operations

**Promises** are the **producer side** of futures:
- Used to fulfill a future with a computed value
- Can be fulfilled exactly once
- Enable decoupling of computation from result consumption

#### The Historical Challenge with Distributed Futures

**Previous implementations** required significant coordination overhead:

**Why coordination is necessary**:
- **Failure recovery**: System must be able to recover from node/process failures
- **State consistency**: Multiple processes must maintain consistent view of object state
- **Resource management**: Objects must be properly cleaned up to prevent memory leaks

**The scalability problem**:
- **Coarse-grained tasks**: Coordination overhead is acceptable for long-running tasks
- **Fine-grained tasks**: Overhead dominates execution time, making coordination impractical

#### The Fine-Grained Task Problem

**Why fine-grained tasks are challenging**:

**Overhead dominance**:
- **System overheads** (coordination, metadata management) become a large fraction of total runtime
- **Application compute time** is much smaller relative to system overhead
- **Result**: System spends more time managing tasks than executing them

**Example**: A task that takes 1ms to execute but requires 10ms of coordination overhead is 90% overhead.

**AIFM's approach**: Offload compute to remote nodes to amortize coordination costs, but this doesn't solve the fundamental scalability problem.
## Previous Solutions and Their Limitations

### Centralized Master Approach

**How it works**:
- **Single point of coordination**: All coordination data recorded at a centralized master
- **Simple implementation**: Straightforward to implement and reason about
- **Global view**: Master has complete visibility into system state

**Why it doesn't scale**:
- **Centralized bottleneck**: Master becomes the limiting factor for system throughput
- **Single point of failure**: Master failure brings down entire system
- **Network congestion**: All coordination traffic flows through master
- **Memory limitations**: Master must track state for entire cluster

### Distributed Leases

**How leases work**:
- **Time-limited ownership**: Workers acquire leases for objects they need to access
- **Automatic expiration**: Leases expire if not renewed, enabling failure detection
- **Distributed coordination**: Multiple workers can coordinate through lease mechanisms

**Limitations for fine-grained tasks**:
- **Lease overhead**: Acquiring and renewing leases adds significant overhead
- **Failure detection delay**: Must wait for lease expiration to detect failures
- **Coordination complexity**: Workers must coordinate to determine recovery responsibilities

## The Ownership Solution: Distributed Futures with Horizontal Scaling

### Core Design Philosophy

**Ownership with distributed futures** represents a fundamental shift in approach:

#### Horizontal Scaling Through Distribution
- **Shard coordination**: Distribute coordination work across all nodes instead of centralizing at master
- **Local decision making**: Each node makes local decisions about objects it owns
- **Eliminate bottlenecks**: No single point of coordination or failure

#### Leveraging Application Semantics
- **Task caller as owner**: The task's caller becomes the owner of both the task and its result
- **Local metadata writes**: Owners perform local metadata updates, eliminating remote coordination
- **Simplified failure handling**: Each worker acts as a "centralized master" for its owned objects

### Object Lifetime Management

#### Distributed Reference Counting

**How it works**:
- **Reference tracking**: System tracks references to each object across the cluster
- **Automatic garbage collection**: Objects with reference count of 0 are automatically reclaimed
- **Distributed coordination**: Reference counts maintained across multiple nodes

#### The Reference Cycle Problem

**Challenge**: What happens with reference cycles?

**Example scenario**:
```python
class A:
    def call(self, B):
        self.x_ref = B.foo.remote()  # A holds reference to B's result

class B:
    def foo(self):
        return self.x_ref            # B returns reference to A's result
```

**The problem**: A references B's result, B references A's result → reference cycle

**Ownership solution**: 
- **Fate-sharing**: Objects in reference cycles fate-share with their owners
- **Collective cleanup**: When owner fails, all objects in the cycle are cleaned up together
- **Lineage reconstruction**: Objects can be regenerated from their computation lineage

### Failure Recovery Through Lineage Reconstruction

#### The Lineage Reconstruction Approach

**Core principle**: **Tasks are run again to produce objects again**

**Key requirements**:
- **Idempotent tasks**: Tasks must produce the same result when run multiple times
- **Minimal rerun**: Only the minimal subset of tasks needed for recovery are rerun
- **Fate-sharing**: Tasks fate-share with owners of objects they reference

#### How Lineage Reconstruction Works

**When failure occurs**:
1. **Identify lost objects**: Determine which objects were lost due to failure
2. **Trace dependencies**: Find all tasks that produced or depend on lost objects
3. **Minimal rerun**: Execute only the minimal set of tasks needed to regenerate lost objects
4. **Update references**: Update all references to point to newly generated objects

**Benefits**:
- **Transparent recovery**: Application doesn't need to handle failure explicitly
- **Efficient**: Only necessary computation is repeated
- **Consistent**: System returns to consistent state after recovery

### Key Insight: Application Semantics for Performance

#### Leveraging Application Structure

**The fundamental insight**: **Use application semantics to optimize system performance**

**How this works**:
- **Task ownership**: Task caller naturally becomes object owner
- **Locality**: Owner is likely to access the object most frequently
- **Local operations**: Most metadata operations can be performed locally
- **Reduced coordination**: Eliminates need for distributed coordination in common cases

#### Comparison with AIFM

**AIFM's approach**: Offload computation to remote nodes to amortize coordination costs
**Ownership's approach**: Eliminate coordination costs through better system design

**Why ownership is superior**:
- **Lower overhead**: No coordination overhead for common operations
- **Better scalability**: Scales with number of nodes, not coordination complexity
- **Simpler programming model**: Leverages natural application structure
## API Design and Programming Model

### The Ownership API

The Ownership system provides a clean, intuitive API for distributed programming:

#### Core API Components

**Task execution**:
```python
# Execute a task remotely, returns a future
future = task.remote(args)

# Get the result (blocks until complete)
result = ray.get(future)

# Check if future is ready
if ray.wait([future], timeout=0)[0]:
    result = ray.get(future)
```

**Actor creation and interaction**:
```python
# Create an actor
actor = ActorClass.remote()

# Call actor method (returns future)
future = actor.method.remote(args)

# Get result
result = ray.get(future)
```

**Object references**:
```python
# Pass futures as arguments to other tasks
future1 = task1.remote()
future2 = task2.remote(future1)  # future1 passed by reference

# Futures can be stored and passed around
futures = [task.remote(i) for i in range(10)]
results = ray.get(futures)  # Wait for all to complete
```

#### Key API Properties

**Transparent distribution**: Code looks like local function calls
**Automatic serialization**: Objects automatically serialized/deserialized
**Fault tolerance**: Built-in failure recovery without application changes
**Composability**: Futures can be combined and nested arbitrarily

## Applications and Use Cases

### Model Serving

**Challenge**: Serve machine learning models with low latency and high throughput

**Ownership solution**:
- **Parallel inference**: Multiple model replicas across nodes
- **Load balancing**: Automatic distribution of inference requests
- **Fault tolerance**: Failed replicas automatically replaced
- **Dynamic scaling**: Add/remove replicas based on load

**Workload graph structure**:
```
Client Request → Load Balancer → Model Replica → Response
                     ↓
              [Multiple Replicas]
```

### Distributed Data Processing

**Challenge**: Process large datasets across multiple nodes efficiently

**Ownership solution**:
- **Pipeline parallelism**: Different stages of processing on different nodes
- **Data parallelism**: Same operation on different data partitions
- **Automatic scheduling**: Tasks scheduled based on data locality
- **Fault recovery**: Failed tasks automatically retried

**Workload graph structure**:
```
Input Data → [Partition 1, Partition 2, Partition 3] → Process → [Result 1, Result 2, Result 3] → Combine → Output
```

### Other Potential Workloads

**Distributed training**: Training neural networks across multiple GPUs
**Real-time analytics**: Processing streaming data with low latency
**Scientific computing**: Parallel simulations and computations
**Microservices**: Coordinating distributed service architectures

**For any workload mentioned by students**: We can sketch the workload graph structure to show how tasks and actors would be organized and how data flows between them.
## System Requirements and Challenges

### Core Requirements

The Ownership system must provide three fundamental capabilities:

#### 1. Automatic Memory Management
- **Garbage collection**: Automatically reclaim objects that are no longer needed
- **Reference counting**: Track references to objects across the distributed system
- **Memory efficiency**: Prevent memory leaks and optimize memory usage

#### 2. Failure Detection
- **Worker failure detection**: Detect when individual workers crash or become unresponsive
- **Node failure detection**: Detect when entire nodes fail
- **Network partition detection**: Detect network connectivity issues

#### 3. Failure Recovery
- **Transparent recovery**: Application should not need to handle failures explicitly
- **Consistent state**: System should return to consistent state after recovery
- **Minimal overhead**: Recovery should be efficient and not impact performance significantly

### The Memory Management Challenge

#### Distributed Garbage Collection

**How it works**:
- **Reference counting**: System tracks references to each object across all nodes
- **Automatic reclamation**: Objects with reference count of 0 are automatically garbage collected
- **Distributed coordination**: Reference counts must be maintained consistently across nodes

**Challenges**:
- **Reference cycles**: Objects that reference each other can create cycles
- **Distributed state**: Reference counts must be kept consistent across multiple nodes
- **Performance**: Garbage collection must not impact application performance

### The Failure Detection Challenge

#### Why Failure Detection Is Complex

**Intuitive expectation**: Shouldn't a worker be able to easily tell when another worker crashes?

**The distributed futures complication**:

**Problem 1: Unknown object locations**
- **Dynamic scheduling**: A worker doesn't know where a value it wants to load will be located
- **Pending objects**: The worker that will generate the value might not be scheduled yet
- **Scheduling updates**: Even if scheduled, the scheduling decision might be updated later

**Problem 2: Metadata management**
- **Object locations**: System must track where each object is stored
- **Task locations**: System must track where each task (pending object) will execute
- **Dynamic updates**: This metadata changes as tasks are scheduled and completed

#### System Metadata Requirements

**The system must maintain**:
- **Object locations**: Where each object is stored (for retrieval by reference holders)
- **Reference status**: Whether objects are still referenced (for garbage collection)
- **Pending object locations**: Where each pending task will execute
- **Object lineage**: How objects were created (for failure recovery)

### The Failure Recovery Challenge

#### Transparency Requirement

**Goal**: Failure recovery should be **transparent to the application**
- **No application changes**: Application code shouldn't need to handle failures
- **Automatic retry**: Failed operations should be automatically retried
- **Consistent semantics**: Application should see consistent behavior despite failures

#### Metadata Consistency

**Critical requirement**: Keep metadata up-to-date during failures

**Metadata includes**:
- **Object locations**: So objects can be retrieved by reference holders
- **Reference status**: For garbage collection decisions
- **Pending object locations**: For task scheduling and execution
- **Object lineage**: For reconstructing lost objects

### Existing Solutions and Their Limitations

#### Centralized Master Approach

**How it works**: Single master maintains all metadata
**Problems**:
- **Scalability bottleneck**: Master becomes limiting factor
- **Single point of failure**: Master failure brings down system
- **Network congestion**: All metadata updates flow through master

#### Distributed Leases

**How distributed leases work**:
- **Time-limited ownership**: Workers acquire leases for objects they need to access
- **Automatic expiration**: Leases expire if not renewed, enabling failure detection
- **Distributed coordination**: Multiple workers coordinate through lease mechanisms

**Why they're insufficient for fine-grained tasks**:

**Slow failure detection**:
- **Lease expiration delay**: Must wait for lease to expire to detect failures
- **Timeout overhead**: Long timeouts for safety, short timeouts for responsiveness
- **No optimal timeout**: Can't find timeout that works for all scenarios

**Complex recovery coordination**:
- **Recovery responsibility**: Upon failure, workers must coordinate to determine who should recover/regenerate objects
- **Coordination overhead**: This coordination adds significant overhead
- **Race conditions**: Multiple workers might try to recover the same object
## The Ownership Solution

### Core Design Principles

**Ownership** represents a fundamental shift in distributed systems design:

#### Distributed Control Plane
- **Eliminate centralization**: Distribute the control plane across all workers instead of centralizing it
- **Local decision making**: Each worker makes local decisions about objects it owns
- **Horizontal scaling**: System scales with number of workers, not coordination complexity

#### Leveraging Application Semantics
- **Task caller as owner**: The task's caller becomes the owner of both the task and its result
- **Natural locality**: Owner is likely to access the object most frequently
- **Local operations**: Most metadata operations can be performed locally

### Why Task Caller as Owner?

#### Performance Benefits

**Task owner is likely to write metadata the most**:
- **Local writes**: Owner can perform metadata updates locally without network coordination
- **Reduced latency**: No network round-trips for common operations
- **Better throughput**: Local operations are orders of magnitude faster than remote operations

**Simplified garbage collection**:
- **Local scope**: If object stays only in owner's scope, garbage collection is much simpler
- **No distributed reference counting**: No need for complex distributed reference counting
- **Lower overhead**: Local reference counting is much more efficient

#### The Two Key Challenges

However, this approach creates two important challenges that must be solved:

##### Challenge 1: First-Class Futures

**The problem**: Futures may leave their owner's scope

**Why this happens**:
- **Futures as values**: Futures can be passed as arguments to other tasks
- **Futures as return values**: Tasks can return futures to their callers
- **Futures as data**: Futures can be stored in data structures and passed around

**The solution**: Distributed reference counting
- **Centralized doesn't scale**: Centralized reference counting becomes a bottleneck
- **Distributed mechanism needed**: Must implement distributed reference counting
- **Reference tracking**: System must track references across multiple nodes

##### Challenge 2: Owner Recovery

**The problem**: What happens when an owner fails?

**The challenge**: Dangling references to objects owned by failed workers

**Ownership's solution**:
- **Fate-sharing**: Objects and their reference holders fate-share with the owner
- **Collective cleanup**: When owner dies, all objects and references are killed together
- **Lineage reconstruction**: System uses lineage reconstruction to regenerate lost objects

**How this works**:
1. **Owner failure detected**: System detects that owner has failed
2. **Reference cleanup**: All references to objects owned by failed worker are invalidated
3. **Object regeneration**: Lost objects are regenerated using lineage reconstruction
4. **Reference update**: New references point to regenerated objects

### Benefits of the Ownership Approach

#### Performance Advantages
- **Local operations**: Most metadata operations are local
- **Reduced coordination**: No need for distributed coordination in common cases
- **Better scalability**: Scales with number of workers, not coordination complexity

#### Simplicity Advantages
- **Natural semantics**: Task caller as owner matches application structure
- **Simplified failure handling**: Each worker acts as "centralized master" for its objects
- **Transparent recovery**: Application doesn't need to handle failures explicitly

#### Reliability Advantages
- **Fault tolerance**: System can recover from worker and node failures
- **Consistency**: System maintains consistent state despite failures
- **Automatic recovery**: Failed operations are automatically retried
## Ownership System Design

### The Ownership Table

#### Core Data Structure

**Each worker maintains an ownership table** that tracks all futures it knows about:

**Owner's view** (complete information):
- **Object metadata**: Complete information about objects it owns
- **Reference tracking**: All references to owned objects
- **Task information**: Details about tasks that produce owned objects
- **Lineage data**: How objects were created (for failure recovery)

**Borrower's view** (subset of information):
- **Object location**: Where to find the object
- **Reference status**: Whether object is still valid
- **Basic metadata**: Essential information for accessing the object

#### Benefits of This Design

**Local operations**: Most metadata operations can be performed locally
**Reduced coordination**: No need for distributed coordination in common cases
**Fault tolerance**: Each worker has sufficient information to handle failures

### Distributed Task Scheduler

#### How Task Scheduling Works

**Ray's distributed scheduler** (detailed in Wednesday's Ray paper):

**Resource allocation process**:
1. **Local resource request**: Owner first requests resources from its local scheduler
2. **Remote fallback**: If no local resources available, scheduler contacts remote schedulers
3. **Lease granting**: Once resources found, scheduler grants owner a lease
4. **Ownership table update**: Owner updates its ownership table with new task information
5. **Resource reuse**: Owner can bypass scheduler and reuse resources if lease still active

#### Benefits of Distributed Scheduling

**Load balancing**: Tasks distributed across available resources
**Locality optimization**: Tasks scheduled close to their data when possible
**Fault tolerance**: Failed schedulers don't bring down entire system

### Distributed Memory Layer

#### Object Storage Strategy

**Distributed object store** manages object storage across the cluster:

**Size-based storage**:
- **Small objects (< 100 KiB)**: Passed by value for efficiency
- **Large objects (≥ 100 KiB)**: Stored in distributed object store

**Memory management**:
- **Primary copy pinning**: Owner's copy is pinned and cannot be evicted
- **LRU eviction**: Non-pinned copies can be evicted under memory pressure
- **Reference-based reclamation**: Objects reclaimed when reference count reaches 0

#### Why Objects Are Immutable

**Question**: Doesn't immutability reduce system utility?

**Answer**: Immutability provides significant benefits:

**Simplified coordination**:
- **No write conflicts**: Multiple readers can access object simultaneously
- **No consistency issues**: Object state cannot change after creation
- **Simplified caching**: Cached copies are always valid

**Better fault tolerance**:
- **No lost updates**: Object state cannot be corrupted by partial failures
- **Simplified recovery**: Object state is deterministic and reproducible
- **Easier debugging**: Object state is predictable and traceable

**Performance benefits**:
- **Lock-free access**: No need for locks or synchronization
- **Efficient caching**: Objects can be cached anywhere without consistency concerns
- **Parallel access**: Multiple tasks can access object simultaneously

### Memory Management Details

#### Object Reclamation Criteria

**Objects are reclaimed when**:
1. **No owner tasks using object**: No tasks on the owner are using the object
2. **No dependent tasks**: No dependent tasks are using or borrowing the object
3. **Reference count zero**: All references to the object have been released

#### Memory Pressure Handling

**Under memory pressure**:
- **LRU eviction**: Non-pinned copies of objects are evicted using LRU policy
- **Owner protection**: Owner's copy is always pinned and protected from eviction
- **Reference tracking**: System tracks which copies exist and where they are located

**Benefits**:
- **Automatic management**: System automatically manages memory without application intervention
- **Efficient utilization**: Memory is used efficiently across the cluster
- **Fault tolerance**: Lost copies can be regenerated from owner or through lineage reconstruction
## Failure Recovery Mechanisms

### Failure Detection

#### Worker vs. Node Failures

**Worker failure detection**:
- **Local scheduler notification**: When a worker fails, its local scheduler publishes this information to other workers and nodes
- **Fast detection**: Worker failures are detected quickly through local monitoring
- **Granular recovery**: Only the specific worker's tasks need to be recovered

**Node failure detection**:
- **Heartbeat mechanism**: Nodes exchange heartbeats to detect when entire nodes fail
- **Slower detection**: Node failures take longer to detect due to heartbeat timeouts
- **Broader impact**: All workers on failed node need to be recovered

#### Why This Distinction Matters

**Worker failures**: Common, fast recovery, minimal impact
**Node failures**: Less common, slower recovery, broader impact
**Recovery strategy**: Different recovery strategies for different failure types

### Lineage Reconstruction

#### How Lineage Reconstruction Works

**The owner performs lineage reconstruction**:

**Process**:
1. **Ownership table scan**: Owner scans its ownership table to determine what needs to be recovered
2. **Dependency analysis**: Identifies the minimal set of tasks that need to be re-run
3. **Task execution**: Re-executes only the necessary tasks
4. **State update**: Updates ownership table with new object locations

#### Optimization Through Ownership Table

**Why consult the ownership table?**
- **Minimal rerun**: Only re-run tasks that are actually needed for recovery
- **Avoid unnecessary overhead**: Could always re-run all tasks, but this would be inefficient
- **Selective recovery**: Ownership table provides precise information about what needs to be recovered

**Benefits**:
- **Efficient recovery**: Only necessary computation is repeated
- **Faster recovery**: Less work means faster recovery time
- **Resource efficiency**: Doesn't waste resources on unnecessary computation

### Object Recovery

#### Simple Object Recovery

**Object recovery is straightforward**:
- **Re-run tasks**: Basically just run the tasks again to produce the objects
- **Idempotent tasks**: Tasks must be idempotent to ensure correct recovery
- **Deterministic results**: Same inputs should produce same outputs

#### Owner Recovery

**All reference holders fate-share with the owner**:

**Fate-sharing includes**:
- **Children**: Tasks that were spawned by the owner
- **Ancestors**: Tasks that spawned the owner
- **Borrowers**: Any task that holds references to owner's objects

**How fate-sharing works**:
- **Owner can pass futures to children**: `DFut` or `SharedDFut` can be passed to child tasks
- **Owner can return values to ancestors**: Owner can return values to parent tasks
- **Any borrower can be child or ancestor**: Reference holders can be in either direction of the call graph

**Recovery implications**:
- **Collective failure**: When owner fails, all fate-shared tasks also fail
- **Collective recovery**: All failed tasks are recovered together
- **Consistent state**: System maintains consistent state after recovery

### Actor Recovery

#### Scope and Limitations

**Actor recovery is outside the scope of this paper**:

**Why it's different**:
- **Stateful entities**: Actors maintain persistent state that must be recovered
- **State recovery**: Local actor state cannot be recovered using the techniques in this paper
- **Different mechanisms**: Actor recovery requires different mechanisms than task recovery

**How it could work**:
- **Reuse ownership mechanism**: Could reuse the same ownership and lineage reconstruction mechanisms
- **State recovery**: Would need additional mechanisms for recovering actor state
- **Checkpointing**: Might require checkpointing or other state persistence mechanisms

#### Future Work

**Actor recovery represents an important area for future research**:
- **State management**: How to efficiently manage and recover actor state
- **Consistency**: How to maintain consistency during actor recovery
- **Performance**: How to minimize the performance impact of actor recovery