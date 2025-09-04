# BigTable: Supplementary Notes

## BigTable Overview

### Motivation
- **Early 2000s**: Google had way more data than anybody else
- **Traditional databases**: Couldn't scale to Google's needs
- **Better than filesystem**: GFS alone wasn't sufficient
- **BigTable optimized for**:
  - Lots of data, large infrastructure
  - Relatively simple queries
- **Dependencies**: Relies on Chubby, GFS

### Design Philosophy
- **Semi-structured data**: More flexible than traditional databases
- **Massive scale**: Handle petabytes of data
- **Simple queries**: Optimized for specific access patterns
- **Ecosystem integration**: Works with GFS and Chubby

## Chubby: Distributed Coordination Service

### Purpose
- **Distributed coordination service**: Allow client applications to synchronize and manage dynamic configuration state
- **Intuition**: Only some parts of an app need consensus!
- **Examples**:
  - Lab 2: Highly available view service
  - Master election in distributed FS (e.g., GFS)
  - Metadata for sharded services

### Implementation
- **Multi-Paxos SMR**: State machine replication using Paxos
- **Consensus as a service**: Provides coordination primitives

### Why Chubby?
- **Many applications need coordination**: Locking, metadata, etc.
- **Paxos is hard**: Every sufficiently complicated distributed system contains an ad-hoc, informally-specified, bug-ridden, slow implementation of Paxos
- **Paxos is known good solution**: But hard to implement and use
- **Consensus as a service**: Chubby provides this

### Chubby API
- **Small files**: Store configuration data
- **Locking**: Distributed locks
- **Sequencers**: Ordering guarantees
- **Filesystem-like API**:
  - Open, Close, Poison
  - GetContents, SetContents, Delete
  - Acquire, TryAcquire, Release
  - GetSequencer, SetSequencer, CheckSequencer

## BigTable Data Model

### Structure
- **Uninterpreted strings**: In rows and columns
- **Mapping**: (r : string) → (c : string) → (t : int64) → string
- **Mostly schema-less**: Column "families" for access control
- **Data sorted by row name**: Lexicographically close names likely to be nearby
- **Versioned data**: Each piece of data versioned via timestamps
  - Either user or server-generated
  - Control garbage-collection

### Key Features
- **Row-based organization**: Data sorted by row name
- **Column families**: Logical grouping of columns
- **Timestamp versioning**: Multiple versions of same data
- **Lexicographic ordering**: Related data stored together

## BigTable Components

### Architecture
- **Client**: Application interface
- **Chubby**: Coordination service
- **Master**: Metadata management
- **Tablet Server**: Data storage and serving
- **GFS**: Persistent storage

### Component Roles
- **Client**: Queries and updates data
- **Chubby**: Provides coordination and locking
- **Master**: Manages tablet assignments and metadata
- **Tablet Server**: Stores and serves tablet data
- **GFS**: Provides persistent storage

## Tablets

### Tablet Structure
- **Each table**: Composed of one or more tablets
- **Starts at one**: Splits once it's big enough
- **Split at row boundaries**: Maintains data locality
- **Tablets size**: ~100MB-200MB
- **Indexed by range**: A tablet is indexed by its range of keys
  - <START> - "c"
  - "c" - <END>

### Tablet Management
- **Each tablet lives on at most one tablet server**: No replication at tablet level
- **Master coordinates**: Assignments of tablets to servers
- **Tablet locations**: Stored in METADATA table
- **Root tablet**: Stores locations of METADATA tablets
- **ROOT tablet location**: Stored in Chubby

### Metadata Hierarchy
- **ROOT tablet**: Location stored in Chubby
- **METADATA tablets**: Store tablet locations
- **Data tablets**: Store actual data
- **Three-level hierarchy**: ROOT → METADATA → DATA

## Tablet Serving

### Data Persistence
- **Tablet data persisted to GFS**: GFS writes replicated to 3 nodes
- **One of these nodes**: Should be the tablet server!
- **Replication**: GFS handles data replication

### Data Structures
- **Memtable**: In-memory map
- **SSTable**: Immutable, on-disk map
- **Commit log**: Operation log used for recovery

### Write Process
- **Writes go to commit log**: Then to the memtable
- **Reads see merged view**: Memtable + SSTables
- **Data could be in memtable or on disk**: Transparent to client

### Read Process
- **Merged view**: Memtable + SSTables
- **Data location**: Could be in memory or on disk
- **Transparent**: Client doesn't need to know location

## Compaction and Compression

### Minor Compaction
- **Memtables spilled to disk**: Once they grow too big
- **Converted to SSTable**: "minor compaction"
- **Size threshold**: When memtable reaches limit

### Major Compaction
- **Periodically**: All SSTables for a tablet compacted
- **Many SSTables → one**: "major compaction"
- **Reduces file count**: Improves read performance

### Compression
- **Each block of SSTable compressed**: Can get enormous ratios with text data
- **Locality helps**: Similar web pages in same block
- **Space efficiency**: Significant storage savings

## Master

### Responsibilities
- **Tracks tablet servers**: Using Chubby
- **Assigns tablets to servers**: Load balancing
- **Handles tablet server failures**: Recovery and reassignment

### Master Startup
1. **Acquire master lock**: In Chubby
2. **Find live tablet servers**: Each tablet server writes its identity to directory in Chubby
3. **Communicate with live servers**: Find out who has which tablet
4. **Scan METADATA tablets**: Find unassigned tablets

### Master Operation
- **Detect tablet server failures**: Assign tablets to other servers
- **Merge tablets**: If they fall below size threshold
- **Handle split tablets**: Splits initiated by tablet servers, master responsible for assigning new tablet
- **Clients never read from master**: Direct access to tablet servers

## Client Operation

### Query Process
1. **Client queries Chubby**: Where is the root tablet?
2. **Chubby returns**: Tablet server 2
3. **Client queries Tablet Server 2**: Where is the METADATA tablet for table T row R?
4. **Tablet Server returns**: Tablet server 1
5. **Client queries Tablet server 1**: Where is table T Row R?
6. **Tablet server 1 returns**: Tablet Server 3
7. **Client reads table T row R**: In Tablet server 3
8. **Tablet Server 3 returns**: Row

### Three-Level Lookup
- **ROOT tablet**: Location in Chubby
- **METADATA tablet**: Location from ROOT tablet
- **Data tablet**: Location from METADATA tablet
- **Direct access**: Client reads directly from data tablet

## Optimizations

### Client-Side Optimizations
- **Clients cache tablet locations**: Reduce lookup overhead
- **Tablet servers only respond**: If Chubby session active, so this is safe
- **Location caching**: Avoid repeated lookups

### Server-Side Optimizations
- **Locality groups**: Put column families that are infrequently accessed together in separate SSTables
- **Smart caching**: On tablet servers
- **Bloom filters**: On SSTables for efficient lookups

### Performance Optimizations
- **Locality groups**: Separate SSTables for different access patterns
- **Smart caching**: Reduce disk I/O
- **Bloom filters**: Avoid unnecessary disk reads
- **Compression**: Reduce storage and I/O

## Key Takeaways

### BigTable Design Principles
- **Semi-structured data**: More flexible than traditional databases
- **Row-based organization**: Data sorted by row name
- **Tablet-based partitioning**: Automatic splitting and merging
- **Three-level metadata**: ROOT → METADATA → DATA
- **GFS integration**: Persistent storage with replication

### Data Model
- **Uninterpreted strings**: Rows, columns, timestamps
- **Column families**: Logical grouping and access control
- **Timestamp versioning**: Multiple versions of same data
- **Lexicographic ordering**: Related data stored together
- **Schema-less**: Flexible data structure

### Architecture Benefits
- **Scalability**: Handle petabytes of data
- **Fault tolerance**: GFS replication and Chubby coordination
- **Performance**: Optimized for specific access patterns
- **Simplicity**: Relatively simple queries
- **Ecosystem integration**: Works with GFS and Chubby

### Tablet Management
- **Automatic splitting**: When tablets grow too large
- **Automatic merging**: When tablets become too small
- **Load balancing**: Master assigns tablets to servers
- **Failure handling**: Automatic reassignment on server failure
- **Metadata hierarchy**: Efficient location lookup

### Storage Layer
- **Memtable**: In-memory writes for performance
- **SSTable**: Immutable on-disk storage
- **Commit log**: Recovery and durability
- **Compaction**: Minor and major compaction
- **Compression**: Space efficiency

### Coordination
- **Chubby integration**: Distributed coordination
- **Master election**: Using Chubby locks
- **Metadata management**: Centralized but fault-tolerant
- **Consensus as a service**: Chubby provides coordination primitives

### Performance Optimizations
- **Client caching**: Tablet location caching
- **Locality groups**: Separate SSTables for different access patterns
- **Smart caching**: Server-side caching
- **Bloom filters**: Efficient lookups
- **Compression**: Storage and I/O efficiency

### Limitations
- **Simple queries**: Not optimized for complex queries
- **Eventual consistency**: Not strongly consistent
- **Single master**: Potential bottleneck
- **GFS dependency**: Relies on GFS for storage
- **Chubby dependency**: Relies on Chubby for coordination

### Modern Relevance
- **NoSQL databases**: Influence on modern NoSQL systems
- **Wide-column stores**: Cassandra, HBase based on BigTable
- **Cloud databases**: Google Cloud Bigtable, Amazon DynamoDB
- **Big data systems**: Foundation for many big data technologies
- **Distributed systems**: Lessons for distributed database design

### Lessons Learned
- **Workload matters**: Design for specific access patterns
- **Ecosystem approach**: Components designed together work better
- **Coordination is hard**: Consensus as a service is valuable
- **Metadata management**: Critical for scalability
- **Storage efficiency**: Compression and compaction are important
- **Fault tolerance**: Replication and coordination are essential