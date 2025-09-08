# Lab 4A: Sharded Key/Value Service

## Introduction

In this lab, you'll build a **sharded key/value storage system** that distributes data across multiple replica groups for improved performance and scalability. This is a fundamental distributed systems pattern used in real-world systems like Google's BigTable, Apache HBase, and many others.

### What is Sharding?

**Sharding** is the process of partitioning data across multiple servers. Instead of storing all key/value pairs on a single server, we divide them into "shards" and distribute these shards across different replica groups. For example:
- Shard 0: keys starting with "a" through "c"
- Shard 1: keys starting with "d" through "f"
- Shard 2: keys starting with "g" through "i"
- And so on...

### Why Sharding?

**Performance**: Each replica group handles only a subset of the data, allowing multiple groups to process requests in parallel. This increases total system throughput proportionally to the number of groups.

**Scalability**: As your system grows, you can add more replica groups to handle increased load.

### System Architecture

Your sharded system has **two main components**:

1. **Replica Groups**: Each group is responsible for a subset of shards. A replica group consists of multiple servers that use Paxos to replicate data within the group.

2. **Shard Master**: A central coordinator that decides which replica group should serve each shard. This information is called the "configuration" and changes over time as groups join/leave or shards are moved.

### The Challenge: Dynamic Reconfiguration

The main challenge is handling **reconfiguration** - when shards need to be moved between groups. This happens when:
- New replica groups join the system
- Existing groups leave the system  
- Load needs to be rebalanced

**Key Requirement**: All servers in a replica group must agree on the order of operations. For example, if a Put request arrives at the same time as a reconfiguration that moves the shard to another group, all servers must agree whether the Put happened before or after the reconfiguration.

### Reconfiguration Process

When a shard moves from Group A to Group B:
1. Group A stops accepting new requests for that shard
2. Group A sends the shard's data to Group B
3. Group B starts accepting requests for that shard

This ensures **at most one group** serves each shard at any time.

### Communication Rules

- **Only RPC** may be used for all communication
- No shared Go variables or files between server instances
- All interaction must go through the network

### Real-World Inspiration

This lab's architecture is inspired by production systems like:
- **Google BigTable**: Distributed storage system
- **Apache HBase**: Open-source BigTable implementation
- **Google Spanner**: Globally distributed database
- **FAWN**: Fast Array of Wimpy Nodes

While simplified, this lab captures the essential challenges of building scalable, fault-tolerant distributed storage systems.

## Software

You should already have the skeleton code for this lab in `src/shardmaster` and `src/shardkv`.

## Part A: The Shard Master

### Overview

The **Shard Master** is the central coordinator that manages which replica groups are responsible for which shards. It maintains a sequence of numbered configurations, where each configuration describes:
- Which replica groups exist
- Which shards each group is responsible for

### Your Task

Implement the shard master in `shardmaster/server.go`. When complete, you should pass all tests:

```bash
$ cd src/shardmaster
$ go test
Test: Basic leave/join ...
  ... Passed
Test: Historical queries ...
  ... Passed
Test: Move ...
  ... Passed
Test: Concurrent leave/join ...
  ... Passed
Test: Min advances after joins ...
  ... Passed
Test: Minimal transfers after joins ...
  ... Passed
Test: Minimal transfers after leaves ...
  ... Passed
Test: Query() returns latest configuration ...
  ... Passed
Test: Concurrent leave/join, failure ...
  ... Passed
PASS
ok      shardmaster     11.200s
```

### RPC Interface

You must implement **four RPC methods** (defined in `shardmaster/common.go`):

#### 1. Join(gid, servers)
- **Purpose**: Add a new replica group to the system
- **Arguments**: 
  - `gid`: Unique group identifier (must be > 0)
  - `servers`: List of server addresses in the group
- **Behavior**: Create a new configuration that includes the new group and redistributes shards evenly

#### 2. Leave(gid)
- **Purpose**: Remove a replica group from the system
- **Arguments**: `gid`: Group identifier to remove
- **Behavior**: Create a new configuration that excludes the group and redistributes its shards to remaining groups

#### 3. Move(shard, gid)
- **Purpose**: Manually assign a specific shard to a specific group
- **Arguments**: 
  - `shard`: Shard number (0 to NShards-1)
  - `gid`: Group identifier to assign the shard to
- **Behavior**: Create a new configuration with the shard moved to the specified group

#### 4. Query(num)
- **Purpose**: Retrieve a configuration
- **Arguments**: `num`: Configuration number (-1 for latest)
- **Behavior**: Return the requested configuration

### Configuration Management

- **Configuration 0**: Initial state with no groups, all shards assigned to GID 0 (invalid)
- **Configuration 1+**: Created in response to Join/Leave/Move operations
- **Numbering**: Each new configuration gets the next sequential number

### Load Balancing Requirements

When redistributing shards (Join/Leave operations):
1. **Even Distribution**: Shards should be divided as evenly as possible among groups
2. **Minimal Movement**: Move as few shards as possible to achieve even distribution
3. **Example**: With 10 shards and 3 groups, ideal distribution is [4, 3, 3] or [3, 4, 3] or [3, 3, 4]

### Fault Tolerance

- **Use Paxos**: Your shard master must be fault-tolerant using the Paxos library from Lab 3
- **No Duplicate Detection**: You don't need to detect duplicate client requests for the shard master (unlike Part B)

### Important Constraints

- **Don't modify**: `common.go` or `client.go`
- **RPC only**: All communication must use RPC
- **No shared state**: No shared Go variables or files between server instances

## Implementation Hints

### Getting Started

1. **Start with KVPaxos**: Begin with a stripped-down copy of your KVPaxos server from Lab 3B
2. **Expected Size**: Part A should take around 200 lines of code
3. **Key Files**: Focus on `shardmaster/server.go` - don't modify `common.go` or `client.go`

### High-Level Design

Your shard master should follow this pattern for **all operations**:

1. **Get Paxos Agreement**: Use Paxos to agree on a sequence number for the operation
2. **Execute Operation**: Apply the operation to create a new configuration
3. **Update State**: Store the new configuration and notify Paxos that you're done

### Detailed Implementation Steps

#### Step 1: Operation Sequencing
```go
// For each RPC (Join, Leave, Move, Query):
// 1. Create an Op struct with unique operation ID
// 2. Use Paxos to agree on a sequence number
// 3. Wait for Paxos to decide on that sequence
```

#### Step 2: Operation Execution
```go
// Execute operations in sequence order:
// 1. Check if operation is a duplicate (optional for Part A)
// 2. Wait for all previous operations to complete
// 3. Apply the operation to create new configuration
// 4. Call px.Done(seq) to free Paxos memory
```

#### Step 3: Load Balancing Algorithm
```go
// For Join/Leave operations:
// 1. Calculate ideal shards per group
// 2. Identify groups with too many/few shards
// 3. Move shards to achieve even distribution
// 4. Minimize the number of shard movements
```

### Critical Implementation Hints

#### Hint 1: Go Map References
```go
// WRONG - both variables point to the same map:
newConfig := oldConfig
newConfig.Groups[gid] = servers  // Modifies oldConfig too!

// CORRECT - create a new map:
newConfig := Config{
    Num: oldConfig.Num + 1,
    Shards: oldConfig.Shards,  // Arrays are copied by value
    Groups: make(map[int64][]string),  // New map
}
// Copy each group individually
for gid, servers := range oldConfig.Groups {
    newConfig.Groups[gid] = servers
}
```

#### Hint 2: Load Balancing Algorithm
```go
// Calculate ideal distribution:
totalShards := NShards
numGroups := len(config.Groups)
shardsPerGroup := totalShards / numGroups
extraShards := totalShards % numGroups

// Some groups get shardsPerGroup + 1 shards
// Others get shardsPerGroup shards
```

#### Hint 3: Paxos Integration
```go
// Use Paxos like in KVPaxos:
// 1. Generate unique operation ID
// 2. Try different sequence numbers until one succeeds
// 3. Wait for Paxos to decide
// 4. Execute all operations up to that sequence
// 5. Call px.Done() to free memory
```

#### Hint 4: Configuration Management
```go
// Store configurations in a slice:
configs []Config

// Configuration 0: Initial state (no groups, all shards to GID 0)
// Configuration 1+: Created by Join/Leave/Move operations
```

### Common Pitfalls to Avoid

1. **Map Reference Bug**: Don't share map references between configurations
2. **Load Balancing**: Ensure shards are distributed as evenly as possible
3. **Sequence Ordering**: Operations must be applied in Paxos sequence order
4. **Memory Management**: Call `px.Done()` to prevent Paxos memory leaks
5. **Configuration Numbering**: Each new configuration gets the next sequential number

### Testing Strategy

1. **Start Simple**: Implement Join first, then Leave, then Move, then Query
2. **Test Incrementally**: Test each RPC method individually
3. **Check Load Balancing**: Verify shards are distributed evenly
4. **Test Concurrency**: Ensure concurrent operations work correctly
5. **Test Fault Tolerance**: Verify system works with server failures

### Debugging Tips

1. **Use DPrintf**: Add debug output to trace operation flow
2. **Check Configurations**: Print configurations after each operation
3. **Verify Paxos**: Ensure Paxos operations complete successfully
4. **Test Load Balancing**: Manually verify shard distribution is correct

---

## Step-by-Step Checklist (Quickstart)
- Implement `Join` creating a new deep-copied config and rebalancing minimally
- Implement `Leave` removing group and rebalancing minimally
- Implement `Move` for a single shard
- Implement `Query` to return past or latest config
- Ensure each op is sequenced through Paxos and `px.Done` is called

## Extra Hints
- Keep a helper to deep-copy `Config` safely (arrays vs maps)
- Prefer deterministic shard movement to ease test debugging
- Add small, focused unit prints (group->shards) during development and remove later
