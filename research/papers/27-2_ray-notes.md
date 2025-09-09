# Ray: A Distributed System for Reinforcement Learning

## Introduction: The Reinforcement Learning Challenge

### Ray's Origins and Evolution

**Ray** was originally designed as a distributed system specifically built for **reinforcement learning (RL)** applications. While it has since evolved into a much more general-purpose distributed computing platform, we will focus our discussion on its foundational support for RL workloads, as this reveals the core design principles that make Ray effective for a wide range of distributed applications.

### Understanding Machine Learning Paradigms

#### Supervised Learning

**Supervised learning** is the most common form of machine learning:

**How it works**:
- **Labeled data**: Inputs have corresponding labels (ground truth)
- **Training process**: Model is trained on input-label pairs
- **Model architecture**: Typically uses deep neural networks
- **Objective**: Learn to map inputs to correct labels

**Example**: Image classification where you have photos labeled with their contents (cat, dog, car, etc.)

#### Reinforcement Learning: A Different Paradigm

**Reinforcement learning** represents a fundamentally different approach to machine learning:

**Core characteristics**:
- **Exploration vs. exploitation**: "Not only to exploit the data gathered, but also to explore the space of possible actions"
- **Continuous operation**: "RL deals with learning to operate continuously within an uncertain environment based on delayed and limited feedback"
- **Policy learning**: "The central goal of an RL application is to learn a policy—a mapping from the state of the environment to a choice of action—that yields effective performance over time, e.g., winning a game or piloting a drone"

**Key differences from supervised learning**:
- **No labeled data**: No ground truth labels provided
- **Delayed feedback**: Rewards come after actions are taken
- **Exploration required**: Must try different actions to learn what works
- **Sequential decisions**: Actions affect future states and rewards
## The Three Pillars of Reinforcement Learning

### What Does RL Require?

Reinforcement learning applications have three fundamental requirements:

#### 1. Simulation for Policy Evaluation and Exploration

**Purpose**: Evaluate policies and explore different actions
- **Policy evaluation**: Test how well a current policy performs
- **Exploration**: Try new actions to discover better strategies
- **Environment modeling**: Simulate the environment where the agent will operate

#### 2. Distributed Training for Policy Improvement

**Purpose**: Improve the policy based on data from simulations
- **Data processing**: Process large amounts of trajectory data
- **Model training**: Update neural networks that represent the policy
- **Distributed computation**: Leverage multiple machines for faster training

#### 3. Policy Serving

**Purpose**: Deploy the learned policy for real-world use

**Two serving modes**:
- **Closed loop**: Make choices and observe their outcomes (learning continues)
- **Open loop**: Make choices without observing outcomes (pure inference)

### How Reinforcement Learning Works

#### The RL Framework

**Core components**:
- **Agent**: The learning entity that makes decisions
- **Environment**: The world in which the agent operates
- **Policy**: A mapping from environment state to action choice
- **Reward**: Feedback signal indicating how good an action was

**Objective**: Agent learns a policy that maximizes cumulative reward over time

#### The Two-Step Learning Process

**Reinforcement learning follows an iterative two-step process**:

##### Step 1: Policy Evaluation

**Process**: Agent interacts with environment using current policy
- **Trajectory generation**: Produces a sequence of (state, action, reward) tuples
- **Performance assessment**: Evaluates how well the current policy performs
- **Data collection**: Gathers experience data for policy improvement

**Example**: In a game, the agent plays using its current strategy and records all moves and outcomes

##### Step 2: Policy Improvement

**Process**: Agent uses collected trajectories to improve the policy
- **Data analysis**: Analyzes trajectory data to identify patterns
- **Policy update**: Modifies the policy to increase expected reward
- **Learning direction**: Updates policy in the direction that maximizes reward

**Example**: The agent analyzes its game performance and adjusts its strategy to make better moves

#### The Learning Loop

**Continuous cycle**:
1. **Evaluate** current policy through simulation
2. **Collect** trajectory data from environment interaction
3. **Improve** policy based on collected data
4. **Repeat** with updated policy

This creates a feedback loop where the agent continuously learns and improves its decision-making capabilities.
## Why Not Use Existing Solutions?

### The Integration Challenge

**Question**: There are already existing solutions that handle simulation, distributed training, and policy serving. Why not just combine these existing solutions?

**Answer**: While individual solutions exist for each component, **stitching them together creates significant problems**:

#### Tight Coupling Between Requirements

**The fundamental issue**: The three RL requirements are **tightly coupled**:

**Simulation ↔ Training coupling**:
- **Data flow**: Simulation produces data that training consumes
- **Latency requirements**: Training needs fresh data from simulations
- **Resource sharing**: Both need access to the same computational resources

**Training ↔ Serving coupling**:
- **Model updates**: New models from training must be deployed to serving
- **Consistency**: Serving must use the latest trained models
- **Performance**: Model updates should not disrupt serving performance

**Simulation ↔ Serving coupling**:
- **Policy evaluation**: Simulations need to test policies currently being served
- **A/B testing**: Compare different policies in simulation before serving
- **Feedback loops**: Serving performance affects simulation strategies

#### Application-System Coupling

**Additional complexity**: In RL applications, **application logic tends to be tightly coupled with the underlying system**:

**Why this matters**:
- **Custom scheduling**: RL workloads have unique scheduling requirements
- **Dynamic resource allocation**: Resource needs change based on learning progress
- **Specialized communication**: Different components need different communication patterns
- **Integrated debugging**: Problems span multiple system components

**Result**: Makes it extremely difficult to stitch multiple systems together effectively

#### Performance Implications

**Stitching existing systems leads to**:
- **High latency**: Data must cross system boundaries
- **Resource inefficiency**: Each system has its own overhead
- **Complexity**: Managing multiple systems increases operational complexity
- **Poor scalability**: Systems don't scale together harmoniously
## What Ray Provides: A Unified Platform

### Core Capabilities

Ray provides a **unified platform** that addresses all three RL requirements in an integrated manner:

#### 1. Support for Fine-Grained Computations

**Why this matters for RL**:
- **Simulation tasks**: Individual simulation steps can be very short (milliseconds)
- **Training updates**: Small gradient updates that need to be applied frequently
- **Policy evaluations**: Quick policy tests that need low latency

**Ray's advantage**: Can handle millions of tasks per second with minimal overhead

#### 2. Support for Heterogeneity

**Temporal heterogeneity**:
- **Simulation duration**: Simulations can take vastly different amounts of time
  - **Short simulations**: Milliseconds (simple games, quick tests)
  - **Long simulations**: Hours (complex environments, detailed scenarios)

**Resource heterogeneity**:
- **Multiple resource types**: CPUs, GPUs, TPUs, and specialized hardware
- **Different resource requirements**: Different tasks need different resources
- **Dynamic allocation**: Resources can be allocated based on task needs

#### 3. Flexible Computation Model

**Ray supports both stateless and stateful computations**:

**Stateless computations (Tasks)**:
- **Functional programming**: Pure functions with no side effects
- **Parallelizable**: Can be executed anywhere in the cluster
- **Fault-tolerant**: Easy to recover from failures

**Stateful computations (Actors)**:
- **Object-oriented**: Maintain state between method calls
- **Long-running**: Can persist across multiple operations
- **Specialized**: Can hold resources like GPU memory

### Why Both Tasks and Actors?

#### Task-Parallel Computations

**Characteristics**:
- **Stateless**: No persistent state between invocations
- **Functional**: Input determines output completely
- **Parallelizable**: Can run on any available worker

**Use cases**:
- **Simulation steps**: Individual environment interactions
- **Data processing**: Transform trajectory data
- **Policy evaluation**: Test policies on specific states

#### Actor-Based Computations

**Characteristics**:
- **Stateful**: Maintain state across multiple method calls
- **Resourceful**: Can hold expensive resources (GPU memory, model weights)
- **Specialized**: Can implement complex behaviors

**Use cases**:
- **Environment simulators**: Maintain simulation state
- **Neural network models**: Hold model parameters and weights
- **Policy servers**: Maintain serving state and handle requests

#### Why You Need Both

**Different requirements**:
- **Tasks**: For parallel, stateless operations
- **Actors**: For stateful, resource-intensive operations

**Complementary roles**:
- **Tasks**: Handle the "compute" part of RL
- **Actors**: Handle the "state" part of RL

**Example**: A simulation might use actors to maintain environment state and tasks to process individual actions

### Dynamic Execution

**Key requirement**: **We don't know what order things will finish in or even which tasks will be invoked through the application lifetime**

**Why this matters for RL**:
- **Exploration**: Don't know which actions will be tried
- **Adaptive learning**: Task creation depends on learning progress
- **Resource availability**: Task scheduling depends on available resources

**Ray's solution**: Dynamic task creation and scheduling based on runtime conditions

### Integration with Existing Systems

**Ray integrates nicely with existing simulators and deep learning frameworks**:
- **Simulators**: Can wrap existing simulation environments
- **Deep learning**: Works with TensorFlow, PyTorch, and other frameworks
- **Legacy systems**: Can gradually migrate from existing solutions
## Comparison with Existing Systems

### Why Not Use MapReduce, Spark, etc.?

**Existing systems** like MapReduce, Spark, and other bulk-synchronous parallel (BSP) systems support distributed computation, so why not use them for RL?

#### Limitations for RL Workloads

**No support for serving**:
- **BSP systems**: Designed for batch processing, not real-time serving
- **RL requirement**: Need to serve policies with low latency
- **Gap**: No built-in support for serving workloads

**No support for fine-grained simulations**:
- **BSP systems**: Designed for coarse-grained tasks (minutes to hours)
- **RL requirement**: Need to handle fine-grained tasks (milliseconds)
- **Overhead**: System overhead dominates for small tasks

#### Complementary Roles

**Important note**: These systems still serve important purposes:

**BSP systems have more expansive APIs**:
- **Rich data processing**: More sophisticated data transformation operations
- **Optimized algorithms**: Highly optimized implementations for batch processing
- **Mature ecosystem**: Well-established tools and libraries

**Different use cases**:
- **BSP systems**: Best for large-scale batch data processing
- **Ray**: Best for interactive, fine-grained, heterogeneous workloads

### The CPU vs. GPU Analogy

**Interesting perspective**: "I kind of view Ray as a CPU and bulk-synchronous parallel systems as a GPU"

**What this means**:
- **Ray (CPU-like)**: General-purpose, flexible, good for diverse workloads
- **BSP systems (GPU-like)**: Specialized, highly optimized for specific patterns

**Implications**:
- **Ray**: Better for RL's diverse, dynamic requirements
- **BSP systems**: Better for large-scale, regular data processing
- **Both needed**: Different tools for different problems
## Programming and Computation Model

### Graph-Based Application Model

**Ray models an application as a graph of dependent tasks**, providing a clear abstraction for understanding and managing distributed computations.

### Tasks: Stateless Computations

#### Task Characteristics

**Tasks** are the fundamental unit of stateless computation in Ray:

**Definition**: Remote function executed on a stateless worker
- **Remote execution**: Function runs on a remote worker, not locally
- **Stateless**: Worker has no persistent state between task executions
- **Future return**: Returns a future that can be dereferenced to get the result

#### Immutable Functions and Idempotency

**Tasks operate on immutable functions**:
- **Immutable inputs**: Function inputs cannot be modified
- **Pure functions**: No side effects, output depends only on input
- **Deterministic**: Same input always produces same output

**Benefits for fault tolerance**:
- **Idempotent tasks**: Tasks can be safely re-executed
- **Simple recovery**: Just re-execute failed tasks (as we saw in the Ownership paper)
- **No state corruption**: No risk of inconsistent state from partial failures

### Actors: Stateful Computations

#### Actor Characteristics

**Actors** provide stateful computation capabilities:

**Definition**: Similar to a class in a program
- **Stateful**: Maintain state between method invocations
- **Object-oriented**: Encapsulate data and behavior
- **Long-running**: Can persist across multiple operations

#### Tasks vs. Actors: Pros and Cons

**See Table 2 in the paper for detailed comparison**:

**Tasks (Pros)**:
- **Fault-tolerant**: Easy to recover from failures
- **Parallelizable**: Can run anywhere in the cluster
- **Scalable**: Can create millions of tasks per second

**Tasks (Cons)**:
- **No state**: Cannot maintain persistent state
- **Limited expressiveness**: Cannot implement stateful algorithms

**Actors (Pros)**:
- **Stateful**: Can maintain persistent state
- **Resourceful**: Can hold expensive resources (GPU memory)
- **Expressive**: Can implement complex, stateful behaviors

**Actors (Cons)**:
- **Complex failure recovery**: State must be recovered separately
- **Less parallelizable**: Tied to specific workers
- **Resource management**: Must manage state and resources carefully

### Ray API

**See Table 1 in the paper for the complete API reference**

**Key API elements**:
- **Task creation**: `@ray.remote` decorator for functions
- **Actor creation**: `@ray.remote` decorator for classes
- **Execution**: `.remote()` method to execute tasks/actor methods
- **Result retrieval**: `ray.get()` to get results from futures

### Computation Graph Representation

#### Two Kinds of Nodes

**Ray represents applications with computation graphs containing two types of nodes**:

**Data objects**:
- **Immutable values**: Results produced by tasks
- **Shared state**: Can be referenced by multiple tasks
- **Futures**: Placeholders for values not yet computed

**Remote function invocations**:
- **Task executions**: Individual task runs
- **Actor method calls**: Method invocations on actors
- **Computation units**: The actual work being performed

#### Three Kinds of Edges

**The graph contains three types of edges that capture different dependencies**:

##### Data Edges

**Purpose**: Capture dependencies between objects and tasks
- **Input dependencies**: Task depends on specific data objects
- **Output relationships**: Task produces specific data objects
- **Data flow**: Shows how data flows through the computation

##### Control Edges

**Purpose**: Capture computation dependencies from nested remote functions
- **Nested calls**: When a task calls another task
- **Execution order**: Shows which tasks must complete before others can start
- **Control flow**: Shows the logical execution sequence

##### Stateful Edges

**Purpose**: Capture state dependencies between multiple method invocations on the same actor
- **Actor state**: Shows dependencies on internal actor state
- **Method sequencing**: Shows order of method calls on the same actor
- **State evolution**: Shows how actor state changes over time

**Useful for**:
1. **Capturing implicit data dependencies**: Dependencies on internal actor state between successive invocations
2. **Maintaining lineage**: Tracking how actor state evolves for failure recovery

**Example**: If an actor method modifies internal state, subsequent method calls depend on that modified state, creating a stateful edge.
## Ray Architecture

### Two-Layer Design

**Ray's architecture consists of two main layers**:

#### Application Layer
**Purpose**: Implements the Ray API and provides the programming interface

**Components**:
- **Driver**: Process that executes the main program
- **Worker**: Executes stateless functions (tasks)
- **Actor**: A stateful process that maintains state

#### System Layer
**Purpose**: Provides high scalability and fault tolerance

**Key responsibility**: **Tracking the status of futures** (as we discussed in the Ownership paper)

### Global Control Store (GCS)

#### Core Functionality

**The GCS maintains control state of the system**:
- **Object locations**: Where each object is stored in the cluster
- **Object metadata**: Size, type, and other properties of objects
- **Task status**: Which tasks are running, completed, or failed
- **Actor state**: Location and status of actors

#### Implementation Details

**GCS is a key-value store backed by Redis**:
- **Redis backend**: Leverages Redis for performance and reliability
- **Key-value interface**: Simple, efficient data model
- **Persistence**: Redis provides durability guarantees

#### Scalability and Fault Tolerance

**Achieves scale with sharding**:
- **Horizontal scaling**: Distribute data across multiple Redis instances
- **Load distribution**: Spread metadata across multiple shards
- **Independent scaling**: Each shard can scale independently

**Provides fault tolerance with per-shard chain replication**:
- **Replication**: Each shard is replicated across multiple nodes
- **Chain replication**: Provides strong consistency guarantees
- **Failure recovery**: Failed shards can be recovered from replicas

#### Critical Design Insight

**Importantly, GCS stores object metadata that other systems store in the scheduler**

**Why this matters**:
- **Low latency**: Object metadata access doesn't require scheduler interaction
- **High throughput**: Can handle many concurrent metadata requests
- **Independent scaling**: GCS and scheduler can scale independently

**Performance benefits**:
- **Scheduler not on critical path**: Task dispatch doesn't require scheduler for object metadata
- **Fast object location**: Can quickly determine where objects are located
- **Efficient resource allocation**: Can make scheduling decisions without blocking on metadata

**Task dispatch process**:
1. **Check object locations**: Query GCS for object locations and sizes
2. **Make scheduling decision**: Choose where to run the task
3. **Launch task**: Execute task on chosen worker
4. **Update metadata**: Update GCS with new object information

This design ensures that the scheduler is not a bottleneck for task dispatch, enabling high throughput and low latency.
### Bottom-up Distributed Scheduler

#### Two-Level Hierarchical Design

**Ray uses a two-level hierarchical scheduler**:

**Global scheduler**:
- **Cluster-wide view**: Has visibility into all nodes in the cluster
- **Resource allocation**: Makes decisions about where to place tasks
- **Load balancing**: Distributes tasks across the cluster

**Per-node local schedulers**:
- **Local resource management**: Manages resources on individual nodes
- **Fast local decisions**: Can make quick decisions for local tasks
- **Resource optimization**: Optimizes resource usage on the local node

#### Bottom-up Scheduling Process

**Submitted tasks are first sent to the local scheduler**:

**Local scheduling attempt**:
- **Check local resources**: See if local node has available resources
- **Local execution**: If resources available, execute task locally
- **Fast path**: Avoids global scheduler for local tasks

**Global scheduler involvement**:
- **When needed**: If no local resources are available
- **Resource constraints**: Machine is overloaded or lacks required resources (e.g., GPU)
- **Fallback mechanism**: Global scheduler finds alternative nodes

#### Preventing Global Scheduler Bottleneck

**By going to the local scheduler first, we can prevent the global scheduler from becoming a bottleneck**:

**Benefits**:
- **Reduced load**: Global scheduler only handles tasks that can't be scheduled locally
- **Lower latency**: Local tasks can be scheduled immediately
- **Better scalability**: System scales with number of nodes, not global scheduler capacity

#### Global Scheduler Decision Making

**The global scheduler chooses a machine based on two criteria**:

**1. Resource availability**: Machine must be able to provide the required resources
**2. Lowest estimated waiting time**: Choose the machine with the shortest expected wait

#### Estimated Waiting Time Calculation

**Formula**: `Estimated waiting time = (1) estimated queuing time of task at node + (2) estimated transfer time of task's remote inputs`

**Components**:
- **Queuing time**: How long the task will wait in the node's queue
- **Transfer time**: How long it will take to transfer the task's input data

#### Question About Transfer Time

**Interesting observation**: "I'm not sure why the estimated transfer time of the inputs is considered. Presumably this could be overlapped with the queuing time."

**The concern**: Transfer time could potentially be overlapped with queuing time, making the current formula potentially suboptimal.

**Alternative approach**: The estimated waiting time should be:
```
max(estimated queuing time, estimated transfer time)
```

**Reasoning**: If transfer and queuing can happen in parallel, the total waiting time should be the maximum of the two, not their sum.

**This suggests**: The current implementation might be conservative, potentially leading to suboptimal scheduling decisions.
### In-Memory Distributed Object Store

#### Efficient Object Transfer

**Transfer objects on the same node via shared memory**:
- **Zero-copy transfers**: Objects are passed by reference, not copied
- **Performance benefit**: Eliminates expensive memory copying operations
- **Low latency**: Shared memory access is extremely fast

#### Object Replication Strategy

**Replicate objects across nodes so that tasks/actors on remote nodes have local access to objects**:
- **Locality optimization**: Tasks can access objects locally when possible
- **Reduced network traffic**: Minimize data transfer across the network
- **Performance improvement**: Local access is much faster than remote access

#### Immutability Benefits

**Objects are immutable, which significantly simplifies the consistency protocol and fault tolerance**:

**Consistency benefits**:
- **No write conflicts**: Multiple readers can access objects simultaneously
- **No synchronization needed**: No locks or coordination required
- **Simplified caching**: Cached copies are always valid

**Fault tolerance benefits**:
- **No lost updates**: Object state cannot be corrupted by partial failures
- **Deterministic recovery**: Object state is predictable and reproducible
- **Easier debugging**: Object state is traceable and consistent

#### Failure Recovery

**On failure, objects are recovered through lineage re-execution** (as discussed in the Ownership paper):
- **Lineage reconstruction**: Re-execute tasks to regenerate lost objects
- **Automatic recovery**: System automatically recovers from failures
- **Transparent to applications**: Applications don't need to handle failures explicitly

#### Lineage Storage

**GCS stores the lineage for both tasks (stateless) and actors (stateful)**:
- **Task lineage**: How stateless tasks were created and what they depend on
- **Actor lineage**: How stateful actors evolved and what state they maintain
- **Recovery information**: Provides all information needed for failure recovery

#### Object Size Limitations

**Only support objects that can fit on a single node**:
- **Single-node constraint**: Objects must fit within one node's memory
- **Large object handling**: Large objects (matrices, trees) require application-level support
- **Ray limitation**: Ray does not provide native support for distributed large objects

**Implications**:
- **Application responsibility**: Applications must handle large objects themselves
- **Design trade-off**: Simplicity vs. support for very large objects
- **Common case optimization**: Optimized for typical object sizes in RL workloads