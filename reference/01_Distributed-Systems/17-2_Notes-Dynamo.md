# Eventual Consistency and Dynamo: Supplementary Notes

## Dynamo Motivation

### Core Requirements
- **Fast, available writes**: Shopping cart - always enable purchases
- **FLP**: Consistency and progress at odds
  - **Paxos**: Must communicate with a quorum
- **Performance**: Strict consistency = "single" copy
  - **Updates serialized**: To single copy
  - **Or, single copy moves**: Between nodes

### Design Philosophy
- **Availability over consistency**: CAP theorem trade-off
- **Fast writes**: Critical for user experience
- **No blocking**: Never wait for unreachable nodes
- **Eventual consistency**: Acceptable for many applications

## Why Fast Available Writes?

### Business Impact
- **Amazon study**: 100ms increase in response time
  - **5% reduction in revenue**
  - **Similar results at other ecommerce sites**
- **99.99% availability**:
  - **Less than an hour outage/year** (total)
  - **Amazon revenue > $300K / minute**

### Economic Justification
- **Response time matters**: Direct impact on revenue
- **Availability critical**: Downtime is expensive
- **User experience**: Fast writes improve satisfaction
- **Business continuity**: Always accept orders

## Dynamo Goals

### Primary Objectives
- **Expose "as much consistency as possible"**: Without sacrificing availability
- **Good latency, 99.9% of the time**: Predictable performance
- **Easy scalability**: Add/remove nodes without disruption

### Design Principles
- **Consistency when possible**: But not at cost of availability
- **Predictable performance**: 99.9% of requests fast
- **Horizontal scaling**: Easy to add capacity
- **Fault tolerance**: Handle node failures gracefully

## Dynamo Consistency

### Eventual Consistency
- **Can have stale reads**: Not all reads see latest write
- **Can have multiple "latest" versions**: Concurrent writes create conflicts
- **Reads can return multiple values**: Application must resolve conflicts
- **Not sequentially consistent**: Can't "defriend and dis"

### Consistency Model
- **Eventual consistency**: System will eventually converge
- **No strong consistency**: Sacrificed for availability
- **Conflict resolution**: Left to application
- **Vector clocks**: Track causal relationships

## External Interface

### Conflict Resolution
- **Applications can choose**: How to handle inconsistency
- **Shopping cart**: Take union of cart versions
- **User sessions**: Take most recent session
- **High score list**: Take maximum score
- **Default**: Highest timestamp wins

### Context and Causal Relationships
- **Context used**: To record causal relationships between gets and puts
- **Once inconsistency resolved**: Should stay resolved
- **Implemented using vector clocks**: Track version history
- **Client responsibility**: Send context with operations

## Dynamo's Vector Clocks

### Vector Clock Structure
- **Each object associated**: With a vector clock
- **Example**: [(node1, 0), (node2, 1)]
- **Each write has coordinator**: And is replicated to multiple other nodes
- **In eventually consistent manner**: Replication is asynchronous

### Vector Clock Operations
- **Nodes in vector clock**: Are coordinators
- **Client sends clock**: With put (as context)
- **Coordinator increments**: Its own index in clock, then replicates across nodes
- **Nodes keep objects**: With conflicting vector clocks
- **These returned**: On subsequent gets

### Conflict Detection
- **If clock(v1) < clock(v2)**: Node deletes v1
- **Vector clock returned**: As context with get
- **Merge of all returned objects' clocks**: Combined context
- **Used to detect inconsistencies**: On write

## Where Does Each Key Live?

### Goals
- **Balance load**: Even as servers join and leave
- **Encourage put/get to see each other**: Locality
- **Avoid conflicting versions**: Consistent routing
- **Solution**: Consistent hashing

### Requirements
- **Clients all have same assignment**: Consistent routing
- **Keys uniformly distributed**: Load balancing
- **Can add/remove nodes**: Without redistributing too many keys
- **Parcel out work**: Of redistributing keys

## Detour: Consistent Hashing

### Basic Concept
- **Node ids hashed**: To many pseudorandom points on a circle
- **Keys hashed onto circle**: Assigned to "next" node
- **Idea used widely**:
  - **Developed for Akamai CDN**
  - **Used in Chord distributed hash table**
  - **Used in Dynamo distributed DB**

### Advantages
- **Uniform distribution**: Keys spread evenly
- **Minimal redistribution**: When nodes join/leave
- **Deterministic**: Same key always goes to same node
- **Scalable**: Works with any number of nodes

## Scaling Systems: Shards

### Sharding Concept
- **Distribute portions**: Of dataset to various groups of nodes
- **Question**: How do we allocate a data item to a shard?
- **Shard master decides**: Which group has which keys
- **Shards operate independently**: No cross-shard coordination

### Client Discovery Problem
- **How do clients know**: Who has what keys?
- **Ask shard master?**: Becomes the bottleneck!
- **Avoid shard master communication**: If possible
- **Can clients predict**: Which group has which keys?

### Recurring Problem
- **Client needs to access**: Some resource
- **Sharded for scalability**: Distribute load
- **How does client find**: Specific server to use?
- **Central redirection won't scale**: Bottleneck problem

## Other Examples

### Sharded Services
- **Scalable shopping cart service**
- **Scalable email service**
- **Scalable cache layer** (Memcache)
- **Scalable network path allocation**
- **Scalable network function virtualization** (NFV)

### Common Pattern
- **Want to assign keys**: To servers without communication
- **Requirement 1**: Clients all have same assignment
- **Requirement 2**: Keys uniformly distributed
- **Requirement 3**: Can add/remove nodes without redistributing too many keys
- **Requirement 4**: Parcel out work of redistributing keys

## Proposal 2: Hashing

### Simple Hashing
- **For n nodes**: Key k goes to hash(k) mod n
- **Hash distributes keys uniformly**: Good load balancing
- **But, new problem**: What if we add a node?
- **Redistribute a lot of keys**: (on average, all but K/n)

### Problems
- **Massive redistribution**: When nodes join/leave
- **Load imbalance**: During redistribution
- **Complexity**: Managing key movement
- **Performance impact**: During rebalancing

## Proposal 3: Consistent Hashing

### Basic Consistent Hashing
- **First, hash the node ids**: Place nodes on circle
- **Keys are hashed**: Go to the "next" node
- **Minimal redistribution**: Only keys between old and new node move
- **Load balancing**: Keys distributed around circle

### Benefits
- **Minimal key movement**: When nodes join/leave
- **Deterministic**: Same key always goes to same node
- **Scalable**: Works with any number of nodes
- **Simple**: Easy to implement and understand

## Proposal 4: Virtual Nodes

### Virtual Node Concept
- **First, hash the node ids**: To multiple locations
- **Each physical node**: Appears as multiple virtual nodes
- **Better load balancing**: More uniform distribution
- **Hash functions come in families**: Members are independent

### Advantages
- **Better load balancing**: More uniform distribution
- **Fault tolerance**: Node failure affects multiple virtual nodes
- **Flexibility**: Can adjust virtual node count per physical node
- **Smooth rebalancing**: Gradual load redistribution

## Load Balancing At Scale

### Virtual Node Benefits
- **Suppose you have N servers**: Using consistent hashing with virtual nodes
- **Heaviest server has x% more load**: Than the average
- **Lightest server has x% less load**: Than the average
- **What is peak load of system?**: N * load of average machine? No!
- **Need to minimize x**: For better load balancing

### Load Distribution
- **Virtual nodes help**: Distribute load more evenly
- **Reduces variance**: In load distribution
- **Improves utilization**: Of all servers
- **Reduces hotspots**: Where single server overloaded

## Key Popularity

### Popularity Problem
- **What if some keys**: Are more popular than others
- **Consistent hashing**: Is no longer load balanced!
- **One model for popularity**: Is the Zipf distribution
- **Popularity of kth most popular item**: 1 < c < 2
  - **1 / k^c**
- **Example**: 1, ½, 1/3, … 1/100 … 1/1000 … 1/10000

### Zipf Distribution
- **"Heavy Tail" Distribution**: Few very popular items, many unpopular ones
- **Examples**:
  - **Web pages**
  - **Movies**
  - **Library books**
  - **Words in text**
  - **Salaries**
  - **City population**
  - **Twitter followers**
- **Whenever popularity is self-reinforcing**: Zipf distribution emerges

## Proposal 5: Table Indirection

### Consistent Hashing Limitations
- **Consistent hashing is (mostly) stateless**: Given list of servers and # of virtual nodes, client can locate key
- **Worst case unbalanced**: Especially with zipf
- **Popular keys**: Can overload single server

### Table Indirection Solution
- **Add small table on each client**: Table maps virtual node → server
- **Shard master reassigns**: Table entries to balance load
- **Dynamic load balancing**: Adjust assignments based on actual load
- **Handles popularity skew**: Popular keys can be moved

## Consistent Hashing in Dynamo

### Dynamo Implementation
- **Each key has "preference list"**: Next nodes around the circle
- **Skip duplicate virtual nodes**: Ensure diversity
- **Ensure list spans data centers**: Geographic distribution
- **Slightly more complex**: Dynamo ensures keys evenly distributed

### Token-Based Routing
- **Nodes choose "tokens"**: (positions in ring) when joining system
- **Tokens used to route requests**: Deterministic routing
- **Each token = equal fraction**: Of the keyspace
- **Load balancing**: Through token selection

## Replication in Dynamo

### Three Parameters: N, R, W
- **N**: Number of nodes each key replicated on
- **R**: Number of nodes participating in each read
- **W**: Number of nodes participating in each write
- **Data replicated**: Into first N live nodes in preference list
- **But respond to client**: After contacting W
- **Reads see values**: From R nodes
- **Common configuration**: (3, 2, 2)

### Quorum Requirements
- **Quorum**: R + W > N
- **Don't wait for all N**: Only need R or W responses
- **R and W will (usually) overlap**: Ensures consistency
- **Flexible configuration**: Can tune for different needs

## Sloppy Quorum

### Sloppy Quorum Concept
- **Never block waiting**: For unreachable nodes
- **Try next node in list**: Keep going until find enough nodes
- **Want get to see most recent put**: (as often as possible)
- **Nodes ping each other**: Each has independent opinion of up/down
- **"Sloppy" quorum**: Nodes can disagree about which nodes are running

### Benefits
- **High availability**: Never block on failed nodes
- **Fast writes**: Don't wait for slow nodes
- **Fault tolerance**: Handle network partitions
- **Eventual consistency**: System will converge

## Replication in Dynamo

### Request Processing
- **Coordinator (or client) sends**: Each request (put or get) to first N reachable nodes in preference list
- **Wait for R replies**: (for read) or W replies (for write)
- **Normal operation**: Gets see all recent versions
- **Failures/delays**:
  - **Writes still complete quickly**
  - **Reads eventually see**: Latest values

### Ensuring Eventual Consistency
- **What if puts end up far away**: From first N?
- **Could happen if some nodes**: Temporarily unreachable
- **Server remembers "hint"**: About proper location
- **Once reachability restored**: Forwards data
- **Nodes periodically synchronize**: Whole DB
- **Fast comparisons**: Using Merkle trees

## Dynamo Deployments

### Deployment Model
- **~100 nodes each**: Per service
- **One for each service**: (parameters global)
- **How to extend to multiple apps?**: Different apps use different (N, R, W)
- **Different configurations**:
  - **Pretty fast, pretty durable**: (3, 2, 2)
  - **Many reads, few writes**: (3, 1, 3) or (N, 1, N)
  - **(3, 3, 3)?**: Strong Consistency
  - **(3, 1, 1)?**: Eventual Consistency

### Configuration Flexibility
- **Tunable parameters**: N, R, W for different needs
- **Read-heavy workloads**: Lower R, higher W
- **Write-heavy workloads**: Higher R, lower W
- **Strong consistency**: R + W > N, R = W = N
- **Eventual consistency**: R + W ≤ N

## Dynamo Results

### Performance
- **Average much faster**: Than 99.9%
- **But, 99.9% acceptable**: Meets requirements
- **Inconsistencies rate in practice**: Allow inconsistency, but minimize it
- **Good latency**: 99.9% of the time

### Consistency Trade-offs
- **Accept some inconsistency**: For better availability
- **Minimize inconsistency**: Through good design
- **Application-level resolution**: Handle conflicts in application
- **Eventual convergence**: System will become consistent

## Dynamo Revisited

### Implementation Model
- **Implemented as library**: Not as a service
- **Each service**: (e.g. shopping cart) instantiated a Dynamo instance
- **When inconsistency happens**:
  - **Is it a problem in Dynamo?**
  - **Is it intended side effect**: Of Dynamo's design?
- **Every service runs its own ops**: Every service needs to be expert at sloppy quorum

### Library vs. Service
- **Library approach**: Each service manages its own Dynamo
- **Expertise required**: Every service needs to understand sloppy quorum
- **Flexibility**: Each service can configure differently
- **Complexity**: More complex to manage

## DynamoDB

### Evolution
- **Replaced Dynamo the library**: With DynamoDB the service
- **DynamoDB**: Strictly consistent key value store
- **Validated with TLA**: And model checking
- **Eventually consistent**: As an option
- **(afaik) no multikey transactions?**: Limited transaction support

### Service Model
- **Managed service**: Amazon handles operations
- **Strict consistency**: By default
- **Eventually consistent**: As option
- **Model checking**: Formal verification
- **Simplified interface**: Easier to use

### Consistency Evolution
- **Dynamo is eventually consistent**: Original design
- **Amazon is eventually strictly consistent**: Evolution to stronger guarantees
- **Service model**: Easier to provide strong consistency
- **Formal verification**: TLA+ model checking

## Discussion

### Design Questions
- **Why is symmetry valuable?**: Do seeds break it?
- **Dynamo and SOA**: What about malicious/buggy clients?
- **Issues with hot keys?**: Popularity skew problems
- **Transactions and strict consistency**: Why were transactions implemented at Google and not at Amazon?
- **Do Amazon's programmers not want strict consistency?**: Trade-off decisions

### Key Insights
- **Symmetry**: Important for load balancing
- **Security**: Malicious clients can cause problems
- **Hot keys**: Popularity skew affects load balancing
- **Transactions**: Complex to implement in distributed systems
- **Consistency preferences**: Different companies make different trade-offs

## Key Takeaways

### Dynamo Design Principles
- **Availability over consistency**: CAP theorem trade-off
- **Fast writes**: Critical for user experience
- **Eventual consistency**: Acceptable for many applications
- **Application-level conflict resolution**: Let application handle conflicts
- **Sloppy quorum**: Never block on failed nodes

### Consistency Model
- **Eventual consistency**: System will eventually converge
- **Vector clocks**: Track causal relationships
- **Conflict resolution**: Application responsibility
- **No strong consistency**: Sacrificed for availability
- **Context tracking**: Maintain causal relationships

### Consistent Hashing
- **Minimal redistribution**: When nodes join/leave
- **Load balancing**: Keys distributed around circle
- **Virtual nodes**: Better load distribution
- **Token-based routing**: Deterministic key placement
- **Popularity skew**: Can cause load imbalance

### Replication Strategy
- **N, R, W parameters**: Tunable for different needs
- **Sloppy quorum**: High availability
- **Preference lists**: Geographic distribution
- **Hint-based recovery**: Handle temporary failures
- **Merkle trees**: Efficient synchronization

### Performance Characteristics
- **Fast writes**: Don't wait for all replicas
- **Good latency**: 99.9% of requests fast
- **High availability**: Never block on failures
- **Eventual consistency**: Accept some inconsistency
- **Scalable**: Easy to add/remove nodes

### Trade-offs
- **Availability vs. consistency**: Choose availability
- **Performance vs. consistency**: Choose performance
- **Simplicity vs. functionality**: Choose simplicity
- **Library vs. service**: Choose library for flexibility
- **Eventual vs. strong consistency**: Choose eventual

### Modern Relevance
- **NoSQL databases**: Many use eventual consistency
- **Distributed systems**: Dynamo influenced many systems
- **CAP theorem**: Classic example of trade-offs
- **Consistent hashing**: Widely used technique
- **Vector clocks**: Still used for conflict detection

### Lessons Learned
- **Availability matters**: More than perfect consistency
- **Application-level resolution**: Can handle many conflicts
- **Sloppy quorum**: Enables high availability
- **Consistent hashing**: Enables scalable key distribution
- **Trade-offs are inevitable**: Choose based on requirements
- **Formal verification**: Important for complex systems
