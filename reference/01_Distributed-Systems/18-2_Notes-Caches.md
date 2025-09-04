# Caches & Memcache: Supplementary Notes

## Caching in Distributed Systems

### Basic Assumption
- **Assume that clients use**: Sharded key-value store to coordinate their output
- **Write buffering question**: Can we start to write done1 before we finish write to k1?
- **No, if sharded and want linearizability**: Must serialize writes
- **What if caches can hold**: Out of date data?

### Potential Problems
- **Asia**: done1 = true, cached (old) k1
- **Africa**: done2 = true, cached (old) k1 and k2
- **Africa**: done2 = true, k2 correct, cached k1 (!)
- **Inconsistency**: Different regions see different states

### Key Challenge
- **Caches can hold stale data**: Creates consistency problems
- **Multiple regions**: Can have different cached values
- **Write ordering**: Must be preserved across regions
- **Consistency requirements**: Must be maintained

## Rules for Caches and Shards

### Correct Execution Requirements
- **Operations applied in processor order**: Maintain execution order
- **All operations to single key serialized**: As if to single copy
- **How do we ensure #2?**: Can serialize each memory location in isolation

### Serialization Principles
- **Per-key serialization**: Each key handled independently
- **Processor order**: Operations within processor maintain order
- **Isolation**: Memory locations can be serialized independently
- **Consistency**: Maintains correctness guarantees

## Invalidation vs. Leases

### Invalidations
- **Track where data is cached**: Know all cache locations
- **When doing a write**: Invalidate all (other) locations
- **Data can live in multiple caches**: During reads
- **Immediate consistency**: All caches updated on write

### Leases
- **Permission to serve data**: For some time period
- **Wait until lease expires**: Before update
- **Time-based consistency**: Leases expire naturally
- **Reduced coordination**: Less communication needed

### Trade-offs
- **Invalidations**: Stronger consistency, more communication
- **Leases**: Weaker consistency, less communication
- **Choose based on**: Consistency requirements vs. performance needs

## Write-Through vs. Write-Back

### Write-Through
- **Writes go to server**: Immediately
- **Caches only hold clean data**: No dirty data in cache
- **Strong consistency**: Cache and server always in sync
- **Higher latency**: Every write goes to server

### Write-Back
- **Writes go to cache**: Initially
- **Dirty cache data written to server**: When necessary
- **Better performance**: Writes are faster
- **Weaker consistency**: Cache and server can diverge

### Design Choice
- **Write-through**: When consistency is critical
- **Write-back**: When performance is critical
- **Trade-off**: Consistency vs. performance

## Write-Through Invalidations

### Process
- **Track all caches**: With read copies
- **On a write**:
  - **Send invalidations**: To all caches with a copy
  - **Each cache invalidates**: Responds
  - **Wait for all invalidations**: Do update
  - **Return**: To client

### Read Operations
- **Reads can proceed**:
  - **If there is cached copy**: Return immediately
  - **Or if cache miss**: Read at server
- **Consistency maintained**: Through invalidation protocol

### Coordination
- **All caches notified**: Before write completes
- **Synchronous invalidation**: Wait for all responses
- **Strong consistency**: All caches see same state
- **Higher latency**: Due to coordination overhead

## Consistency Questions

### Key Questions
1. **While write to key k waiting on invalidations**: Can other clients read old values of k from their caches?
2. **While write to key k from client C waiting**: Can C perform another write to different key m?
3. **While write to key k from client C waiting**: Can server perform read from different client D to different key m?
4. **While write to key k from client C waiting**: Can server perform read to k from different client D?
5. **While write to key k from client C waiting**: Can server perform write from client D to same key?

### Consistency Analysis
- **Question 1**: No, other clients blocked from reading k
- **Question 2**: Yes, different keys can be written independently
- **Question 3**: Yes, different keys can be read independently
- **Question 4**: No, reads to k blocked during write
- **Question 5**: No, writes to k serialized

## Facebook's Memcache Service

### Facebook's Scaling Problem
- **Rapidly increasing user base**:
  - **Small initial user base**
  - **2x every 9 months**
  - **2013: 1B users globally**
- **Users read/update many times per day**:
  - **Increasingly intensive app logic per user**
  - **2x I/O every 4-6 months**
- **Infrastructure has to keep pace**: With user growth

### Scaling Challenges
- **Exponential growth**: User base doubling every 9 months
- **Increasing complexity**: More app logic per user
- **I/O growth**: 2x every 4-6 months
- **Infrastructure scaling**: Must keep pace with demand

## Scaling Strategy

### Approach
- **Adapt off the shelf components**: Where possible
- **Fix as you go**: No overarching plan
- **Rule of thumb**: Every order of magnitude requires a rethink

### Design Philosophy
- **Pragmatic approach**: Use existing solutions when possible
- **Iterative improvement**: Fix problems as they arise
- **Scale-driven redesign**: Major changes at each order of magnitude
- **No perfect plan**: Adapt to changing requirements

## Facebook Three Layer Architecture

### Application Front End
- **Stateless**: Rapidly changing program logic
- **If app server fails**: Redirect client to new app server
- **Horizontally scalable**: Easy to add/remove servers
- **No persistent state**: All state in backend

### Memcache
- **Lookaside key-value cache**: High-performance cache layer
- **Keys defined by app logic**: Can be computed results
- **Independent of backend**: Works with any storage
- **High volume, low latency**: Optimized for performance

### Fault Tolerant Storage Backend
- **Stateful**: Persistent data storage
- **Careful engineering**: To provide safety and performance
- **Both SQL and NoSQL**: Multiple storage options
- **Reliable**: Handles failures gracefully

## Workload Characteristics

### User Page Uniqueness
- **Each user's page is unique**: Personalized content
- **Draws on events**: Posted by other users
- **Users not in cliques**: For the most part
- **Complex dependencies**: Between user data

### Popularity Distribution
- **User popularity is zipf**: Some users much more popular
- **Some user posts affect**: Very large #'s of other pages
- **Most affect**: Much smaller number
- **Heavy tail distribution**: Few very popular users

### Scaling Implications
- **Personalized content**: Each page unique
- **Social connections**: Complex data dependencies
- **Popularity skew**: Some content much more accessed
- **Cache effectiveness**: Depends on access patterns

## Scale By Caching: Memcache

### Sharded In-Memory Cache
- **Key, values assigned**: By application code
- **Values can be data**: Result of computation
- **Independent of backend**: Storage architecture (SQL, NoSQL) or format
- **Design for high volume**: Low latency

### Key Features
- **Application-defined keys**: Flexible key naming
- **Computed values**: Can cache results of expensive operations
- **Storage agnostic**: Works with any backend
- **Performance optimized**: High throughput, low latency

## Lookaside Architecture

### Lookaside Operations (Read)
- **Webserver needs key value**: Requests from memcache
- **Memcache**: If in cache, return it
- **If not in cache**:
  - **Return error**
  - **Webserver gets data**: From storage server
  - **Possibly SQL query**: Or complex computation
  - **Webserver stores result**: Back into memcache

### Read Process
- **Cache hit**: Return immediately
- **Cache miss**: Application fetches from storage
- **Application responsibility**: To populate cache
- **Asynchronous**: Cache population doesn't block

### Question: Cache Miss Storms
- **What if swarm of users**: Read same key at same time?
- **Multiple cache misses**: For same key
- **Storage server overload**: From concurrent requests
- **Need solution**: To prevent thundering herd

## Lookaside Operation (Write)

### Write Process
- **Webserver changes value**: That would invalidate memcache entry
  - **Could be update to key**
  - **Could be update to value**: Used to derive some key value
- **Client puts new data**: On storage server
- **Client invalidates entry**: In memcache

### Write Strategy
- **Update storage first**: Ensure data persisted
- **Then invalidate cache**: Remove stale data
- **Application responsibility**: To maintain consistency
- **Cache invalidation**: Removes potentially stale data

## Memcache Consistency

### Consistency Questions
- **Is memcache linearizable?**: No, not guaranteed
- **Is lookaside protocol eventually consistent?**: Yes, with caveats

### Consistency Model
- **Not linearizable**: No global ordering guarantees
- **Eventually consistent**: System will converge
- **Application responsibility**: To handle inconsistencies
- **Cache invalidation**: Helps maintain consistency

## Lookaside With Leases

### Goals
- **Reduce (eliminate?) per-key inconsistencies**: Better consistency
- **Reduce cache miss swarms**: Prevent thundering herd

### Lease Mechanism
- **On a read miss**:
  - **Leave marker in cache**: (fetch in progress)
  - **Return timestamp**
  - **Check timestamp when filling cache**
  - **If changed means value has (likely) changed**: Don't overwrite
- **If another thread read misses**:
  - **Find marker and wait**: For update (retry later)

### Lease Benefits
- **Prevents thundering herd**: Only one fetch per key
- **Reduces inconsistencies**: Timestamp-based coordination
- **Better performance**: Fewer redundant fetches
- **Coordination**: Between concurrent requests

### Lease Questions
- **What if web server crashes**: While holding lease?
- **Is lookaside with leases linearizable?**: Still no
- **Is lookaside with leases eventually consistent?**: Yes, improved
- **Would this be made "more correct"?**:
  - **Read misses obtain lease**
  - **Writes obtain lease**: (prevent reads during update)

### Lease Limitations
- **FB replicates popular keys**: Need lease on every copy?
- **Memcache server might fail**: Or appear to fail by being slow
- **Complexity**: More complex than simple invalidation
- **Failure handling**: What happens when servers fail?

## Latency Optimizations

### Concurrent Lookups
- **Issue many lookups concurrently**: Parallel requests
- **Prioritize those**: That have chained dependencies
- **Reduce latency**: Through parallelism
- **Dependency awareness**: Order requests by dependencies

### Batching
- **Batch multiple requests**: (e.g., for different end users) to same memcache server
- **Reduce network overhead**: Fewer round trips
- **Improve throughput**: More efficient communication
- **Server efficiency**: Better resource utilization

### Incast Control
- **Limit concurrency**: To avoid collisions among RPC responses
- **Prevent network congestion**: From too many concurrent requests
- **Smooth traffic**: Avoid bursty patterns
- **Better performance**: More predictable latency

## More Optimizations

### Return Stale Data
- **Return stale data to web server**: If lease is held
- **No guarantee**: That concurrent requests returning stale data will be consistent with each other
- **Frequently accessed, cheap to recompute**: Acceptable trade-off
- **If mixed, frequent accesses**: Will evict all others

### Replicate Keys
- **Replicate keys**: If access rate is too high
- **Implication for consistency?**: More complex consistency model
- **Load distribution**: Spread popular keys across servers
- **Consistency trade-off**: Weaker consistency for better performance

## Gutter Cache

### Problem
- **When memcache server fails**: Flood of requests to fetch data from storage layer
- **Slows users needing any key**: On failed server
- **Slows other users**: Due to storage server contention
- **Cascading failures**: One failure affects many users

### Solution: Backup Cache
- **Backup (gutter) cache**: Secondary cache for failed servers
- **Time-to-live invalidation**: Ok if clients disagree as to whether memcache server is still alive
- **TTL is eventually consistent**: Will converge over time
- **Reduces load**: On storage servers during failures

### Gutter Cache Benefits
- **Fault tolerance**: Handle server failures gracefully
- **Reduced load**: On storage servers
- **Better user experience**: Less impact from failures
- **Eventually consistent**: TTL-based invalidation

## Key Takeaways

### Caching Design Principles
- **Consistency vs. performance**: Fundamental trade-off
- **Invalidation vs. leases**: Different consistency models
- **Write-through vs. write-back**: Different performance characteristics
- **Application responsibility**: For maintaining consistency
- **Fault tolerance**: Handle cache server failures

### Memcache Architecture
- **Lookaside cache**: Application-managed caching
- **Sharded design**: Distribute load across servers
- **High performance**: Optimized for low latency
- **Storage agnostic**: Works with any backend
- **Application-defined keys**: Flexible caching strategy

### Consistency Models
- **Not linearizable**: No global ordering guarantees
- **Eventually consistent**: System will converge
- **Lease-based coordination**: Reduce inconsistencies
- **Timestamp-based**: Conflict detection and resolution
- **Application-level**: Consistency handling

### Performance Optimizations
- **Concurrent lookups**: Parallel request processing
- **Batching**: Reduce network overhead
- **Incast control**: Prevent network congestion
- **Stale data serving**: Acceptable for some use cases
- **Key replication**: For high-access keys

### Fault Tolerance
- **Gutter cache**: Handle server failures
- **TTL-based invalidation**: Eventually consistent
- **Graceful degradation**: Continue serving during failures
- **Load distribution**: Reduce impact of failures
- **Recovery mechanisms**: Restore service after failures

### Scaling Strategies
- **Horizontal scaling**: Add more cache servers
- **Sharding**: Distribute keys across servers
- **Replication**: For popular keys
- **Load balancing**: Distribute requests evenly
- **Capacity planning**: Handle growth

### Trade-offs
- **Consistency vs. performance**: Choose based on requirements
- **Complexity vs. functionality**: More features = more complex
- **Latency vs. throughput**: Optimize for different metrics
- **Memory vs. computation**: Cache vs. recompute
- **Availability vs. consistency**: CAP theorem trade-offs

### Modern Relevance
- **Distributed caching**: Still widely used
- **CDN systems**: Global content distribution
- **Microservices**: Service-to-service caching
- **Database caching**: Reduce database load
- **Application caching**: Improve user experience

### Lessons Learned
- **Caching is complex**: Many trade-offs to consider
- **Consistency is hard**: In distributed systems
- **Performance matters**: Users notice latency
- **Fault tolerance**: Essential for production systems
- **Application awareness**: Caching strategy depends on use case
- **Monitoring**: Essential for cache effectiveness
