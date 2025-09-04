# Google File System (GFS): Supplementary Notes

## Google Stack Overview

### Components
- **GFS**: Large-scale storage for bulk data
- **Chubby**: Paxos storage for coordination
- **BigTable**: Semi-structured data storage
- **MapReduce**: Big data computation on key-value pairs
- **MegaStore, Spanner**: Transactional storage with geo-replication

### Integration
- **Designed together**: GFS for Google apps, Google apps for GFS
- **Co-optimization**: Each component designed for specific workload
- **Ecosystem approach**: Components work together seamlessly

## GFS Design Motivation

### Problem
- **Needed**: Distributed file system for storing results of web crawl and search index
- **Why not NFS?**: Very different workload characteristics
- **Custom solution**: Design GFS for Google apps, Google apps for GFS

### Requirements
- **Fault tolerance**: Handle node failures gracefully
- **Availability**: System remains accessible
- **Throughput**: High data transfer rates
- **Scale**: Handle massive data volumes
- **Concurrent streaming**: Reads and writes simultaneously

## GFS Workload Characteristics

### Producer/Consumer Pattern
- **Hundreds of web crawling clients**: Continuous data ingestion
- **Periodic batch analytic jobs**: MapReduce processing
- **Throughput focus**: Not latency-sensitive
- **Big data sets**: 1000 servers, 300 TB of data stored

### File Characteristics
- **Few million files**: 100MB+ each
- **Many are huge**: Large file sizes
- **Reads**: Mostly large streaming reads, some sorted random reads
- **Writes**: Most files written once, never updated
- **Appends**: Most writes are appends (e.g., concurrent workers)

### Evolution
- **Later**: BigTable tablet log and SSTables
- **Even later**: Workload now includes small file updates
- **2010**: Incremental updates of Google search index

## GFS Interface

### Design Philosophy
- **App-level library**: Not a kernel file system
- **Not POSIX**: Custom interface for specific needs
- **Operations**: Create, delete, open, close, read, write, append

### Consistency Model
- **Metadata operations**: Linearizable
- **File data**: Eventually consistent (stale reads allowed)
- **Snapshots**: Inexpensive file and directory snapshots

### Key Insight
- **Life without random writes**: Optimized for append-only workloads
- **Chunk-based storage**: Each file stored as 64MB chunks
- **Replication**: Each chunk on 3+ chunkservers
- **Single master**: Stores metadata only

## Single Master Architecture

### Master Responsibilities
- **Metadata storage**:
  - File name space, file name → chunk list
  - Chunk ID → list of chunkservers holding it
  - Metadata stored in memory (~64B/chunk)
- **No file contents**: Master does not store file data
- **Direct access**: All requests for file data go directly to chunkservers

### Fault Tolerance
- **Hot standby replication**: Using shadow masters
- **Faster recovery**: Multiple master replicas
- **Linearizable operations**: All metadata operations are linearizable

## Master Fault Tolerance

### Architecture
- **One master, set of replicas**: Master chosen by Chubby
- **Master logs**: Some metadata operations
  - Changes to namespace, ACLs, file → chunk IDs
  - Not chunk ID → chunkserver (why not?)
- **Replication**: Operations at shadow masters and log to disk, then execute
- **Periodic checkpoint**: Master in-memory data
  - Allows master to truncate log, speed recovery
  - Checkpoint proceeds in parallel with new operations

### Design Choices
- **Selective logging**: Not all metadata changes logged
- **Checkpointing**: Periodic snapshots for recovery
- **Parallel operations**: Checkpoints don't block normal operations

## Handling Write Operations

### Mutation Types
- **Mutation**: Write or append operation
- **Goal**: Minimize master involvement
- **Lease mechanism**: Master picks one replica as primary; gives it a lease
- **Primary responsibility**: Defines serial order of mutations

### Data Flow Design
- **Decoupled**: Data flow decoupled from control flow
- **Efficient**: Reduces master bottleneck
- **Scalable**: Handles concurrent operations

## Write Operations Process

### Step-by-Step
1. **Application originates** write request
2. **GFS client translates** request from (fname, data) → (fname, chunk-index), sends to master
3. **Master responds** with chunk handle and (primary + secondary) replica locations
4. **Client pushes** write data to all locations; data stored in chunkservers' internal buffers
5. **Client sends** write command to primary
6. **Primary determines** serial order for data instances stored in its buffer and writes instances in that order to the chunk
7. **Primary sends** serial order to secondaries and tells them to perform the write
8. **Secondaries respond** to the primary
9. **Primary responds** back to client

### Error Handling
- **If write fails**: At one of the chunkservers, client is informed and retries
- **Stale data**: Another client may read stale data from chunkserver
- **Retry mechanism**: Client handles failures transparently

## At Least Once Append

### Semantics
- **If failure**: At primary or any replica, retry append (at new offset)
- **Append will eventually succeed**: Guaranteed completion
- **May succeed multiple times**: Duplicate appends possible

### Application Responsibility
- **App client library** responsible for:
  - Detecting corrupted copies of appended records
  - Ignoring extra copies (during streaming reads)
- **Why not exactly once?**: Complexity vs. performance trade-off

### BigTable Integration
- **Question**: Does BigTable tablet server use "at least once append" for its operations log?
- **Answer**: Yes, for performance and simplicity

## Caching Strategy

### Client-Side Caching
- **GFS caches file metadata** on clients
- **Example**: Chunk ID → chunkservers
- **Used as hint**: Invalidate on use
- **Scale**: TB file → 16K chunks

### No Data Caching
- **GFS does not cache file data** on clients
- **Reason**: Large files, streaming access patterns
- **Trade-off**: Memory usage vs. performance

## Garbage Collection

### Deletion Process
- **File delete** → rename to hidden file
- **Background task** at master:
  - Deletes hidden files
  - Deletes any unreferenced chunks
- **Simpler than foreground deletion**: What if chunk server is partitioned during delete?

### Benefits
- **Need background GC anyway**: Stale/orphan chunks
- **Fault tolerance**: Handles network partitions
- **Efficiency**: Batch operations

## Data Correction

### Corruption Sources
- **Linux bugs**: Sometimes silent corruptions
- **Disk failures**: Disks are not fail-stop
- **Stored blocks**: Can become corrupted over time
- **Example**: Writes to sectors on nearby tracks
- **Scale effect**: Rare events become common at scale

### Detection and Prevention
- **Chunkservers maintain per-chunk CRCs** (64KB)
- **Local log**: CRC updates
- **Verification**: CRCs before returning read data
- **Periodic revalidation**: Detect background failures

## Design Discussion

### Questions
- **Is this a good design?**: Trade-offs and benefits
- **Can we improve on it?**: Areas for optimization
- **Will it scale?**: To even larger workloads

### 15 Years Later
- **Scale is much bigger**:
  - Now 10k servers instead of 1K
  - Now 100 PB instead of 100TB
- **Bigger workload change**: Updates to small files!
- **Around 2010**: Incremental updates of Google search index

## GFS → Colossus Evolution

### GFS Limitations
- **GFS scaled to**: ~50 million files, ~10 PB
- **Developer constraints**: Had to organize apps around large append-only files
- **Latency-sensitive applications**: Suffered from GFS design
- **Eventually replaced**: With new design, Colossus

### Metadata Scalability
- **Main scalability limit**: Single master stores all metadata
- **HDFS has same problem**: Single NameNode
- **Solution**: Partition metadata among multiple masters
- **New system supports**: ~100M files per master and smaller chunk sizes (1MB instead of 64MB)

## Reducing Storage Overhead

### Replication vs. Erasure Coding
- **Replication**: 3x storage to handle two copies
- **Erasure coding**: More flexible - m pieces, n check pieces
- **Example**: RAID-5 - 2 disks, 1 parity disk (XOR of other two)
- **Result**: 1 failure with only 1.5x storage

### Trade-offs
- **Sub-chunk writes**: More expensive (read-modify-write)
- **Recovery**: Get all other pieces, generate missing one
- **Complexity**: Higher implementation complexity

## Erasure Coding Examples

### Comparison
- **3-way replication**:
  - 3x overhead, 2 failures tolerated, easy recovery
- **Google Colossus**: (6, 3) Reed-Solomon code
  - 1.5x overhead, 3 failures
- **Facebook HDFS**: (10, 4) Reed-Solomon
  - 1.4x overhead, 4 failures, expensive recovery
- **Azure**: More advanced code (12, 4)
  - 1.33x overhead, 4 failures, same recovery cost as Colossus

### Design Principles
- **Storage efficiency**: Lower overhead than replication
- **Fault tolerance**: Multiple failure handling
- **Recovery cost**: Balance between efficiency and complexity

## Key Takeaways

### GFS Design Principles
- **Append-only optimization**: Designed for streaming writes
- **Single master**: Simple but limited scalability
- **Chunk-based storage**: 64MB chunks for efficiency
- **Replication**: 3x replication for fault tolerance
- **At-least-once semantics**: Simpler than exactly-once

### Workload Characteristics
- **Large files**: 100MB+ files
- **Streaming access**: Sequential reads and writes
- **Producer-consumer**: Web crawling and MapReduce
- **Throughput focus**: Not latency-sensitive
- **Append-heavy**: Most writes are appends

### Architecture Benefits
- **Fault tolerance**: Handle node failures
- **High throughput**: Optimized for bulk operations
- **Simple interface**: App-level library
- **Efficient caching**: Metadata caching on clients
- **Background operations**: Garbage collection and data correction

### Limitations
- **Single master**: Scalability bottleneck
- **Metadata limits**: ~50 million files maximum
- **Latency**: Not suitable for latency-sensitive applications
- **Small file updates**: Inefficient for random writes
- **Storage overhead**: 3x replication cost

### Evolution to Colossus
- **Metadata partitioning**: Multiple masters
- **Smaller chunks**: 1MB instead of 64MB
- **Better small file support**: Improved random access
- **Higher scale**: ~100M files per master
- **Erasure coding**: Reduced storage overhead

### Lessons Learned
- **Workload matters**: Design for specific access patterns
- **Scale changes everything**: What works at 1K servers may not work at 10K
- **Trade-offs are inevitable**: Consistency vs. performance vs. complexity
- **Evolution is necessary**: Systems must adapt to changing requirements
- **Ecosystem approach**: Components designed together work better

### Modern Relevance
- **HDFS**: Based on GFS principles
- **Cloud storage**: Similar design patterns
- **Big data systems**: Append-only optimization
- **Distributed file systems**: GFS influence on modern systems
- **Storage efficiency**: Erasure coding adoption