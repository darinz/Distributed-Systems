# Primary-Backup Replication: Supplementary Notes

## State Machine Replication

### Concept
- **Replicate state machine** across multiple servers
- **Clients view** all servers as one state machine
- **Simplest form**: Two servers (primary + backup)

### Single-Node State Machine
- **Key/value store** on single node
- **Operations**: Put, Get
- **State**: Current key-value pairs

## Primary-Backup Architecture

### Basic Setup
- **Primary**: Handles client requests
- **Backup**: Replicates primary's state
- **Clients**: Only talk to primary
- **Goal**: Correct and available despite failures

### Operation Flow
1. **Client sends operation** (Put, Get) to primary
2. **Primary decides order** of operations
3. **Primary forwards sequence** to backup
4. **Backup performs operations** in same order
5. **Primary replies to client** after backup confirms

### Backup Types
- **Hot standby**: Backup executes operations immediately
- **Cold standby**: Backup saves log of operations

## The View Service

### Purpose
- **Decides roles**: Who is primary and backup
- **Manages transitions**: When roles change
- **Single point of failure**: Critical component

### View Concept
- **View**: Statement about current roles in system
- **Views form sequence** in time
- **Clients/servers depend** on view server

### View Management
- **Primary fails** → Backup becomes primary, idle server becomes backup
- **Backup fails** → Idle server becomes backup
- **OK to have** primary with no backup (but risky)

## Failure Detection

### Ping Mechanism
- **Each server pings** view server periodically
- **View server considers node "dead"** if misses n pings
- **Node is "live"** after single ping
- **False positives possible**: Network partitions, message loss

### Server Management
- **Any number of servers** can send pings
- **Extras are "idle"** if more than two servers live
- **Idle servers** can be promoted to backup

### View Server Protocol
- **Waits for primary ack** before changing view
- **Must stay with current view** until ack received
- **Even if primary seems failed** - prevents split brain

## Split Brain Prevention

### Rules to Prevent Multiple Primaries
1. **Primary in view i+1** must have been backup or primary in view i
2. **Primary must wait** for backup to accept/execute each op before replying
3. **Backup must accept** forwarded requests only if view is correct
4. **Non-primary must reject** client requests
5. **Every operation** must be before or after state transfer

### Key Insight
- **Only one primary** can respond to clients at a time
- **View server coordination** prevents split brain
- **Backup validation** ensures correct view

## Challenges

### 1. Non-Deterministic Operations
- **Problem**: Operations may behave differently on primary vs backup
- **Solution**: Make operations deterministic or handle non-determinism

### 2. Dropped Messages
- **Problem**: Messages between primary and backup may be lost
- **Solution**: Retry mechanisms, acknowledgments

### 3. State Transfer
- **Problem**: How to initialize backup state
- **Options**: Write log vs write state
- **Challenge**: Must include RPC data

### 4. Single Primary Constraint
- **Problem**: Clients, primary, and backup must agree on roles
- **Solution**: View server coordination

## Progress and Liveness

### Cases Where System Can't Make Progress
- **View server fails**: No role coordination
- **Network fails entirely**: Hard to work around
- **Client can't reach primary** but can ping view server
- **No backup and primary fails**: No failover possible
- **Primary fails before completing state transfer**: Inconsistent state

### State Transfer Requirements
- **Must include RPC data**: Complete state needed
- **Must be atomic**: All or nothing transfer
- **Must handle concurrent operations**: Before or after transfer

## Read Operations

### "Fast" Reads Question
**Does primary need to forward reads to backup?**

### Read vs Write Handling
- **Reads treated as state machine operations** too
- **Can be executed more than once** (idempotent)
- **RPC library can handle differently** from writes
- **Optimization**: Primary might not forward reads to backup

## Deterministic Replay

### Key Idea
- **VM state depends only on input**:
  - Content of all input/output
  - Precise instruction of every interrupt
  - Only few exceptions (e.g., timestamp instruction)

### Implementation
- **Record all hardware events** into log
- **Modern processors**: Instruction counters, precise interrupts
- **Trap and emulate** non-deterministic instructions

## Replicated Virtual Machines

### Concept
- **Whole system replication**: Complete VM replication
- **Transparent to applications**: No changes needed
- **High availability**: For any existing software
- **Restricted to uniprocessor VMs**: Simpler to replicate

### Operation
- **Backup executes events** with lag behind primary
- **Backup stalls** until knows timing of next event
- **Backup doesn't perform** external events
- **Primary stalls** until backup has copy of every event up to output
- **On failure**: Inputs/outputs replayed at backup (idempotent)

### Example Flow
1. **Primary**: Hypervisor forwards interrupt + data to backup
2. **Primary**: Hypervisor delivers network interrupt to OS kernel
3. **Primary**: OS kernel runs, delivers packet to server
4. **Primary**: Server writes response to network card
5. **Primary**: Hypervisor sends response to backup
6. **Primary**: Hypervisor delays sending to client until backup acks
7. **Backup**: Receives log entries, delivers network interrupt
8. **Backup**: Hypervisor does NOT put response on wire
9. **Backup**: Hypervisor ignores local clock interrupts

## Key Takeaways

### Primary-Backup Benefits
- **Simple to understand**: Clear primary/backup roles
- **Transparent to clients**: Only talk to primary
- **High availability**: Backup can take over
- **Deterministic replay**: VM-level replication possible

### Challenges
- **View server**: Single point of failure
- **Split brain prevention**: Complex coordination needed
- **State transfer**: Must be complete and atomic
- **Non-deterministic operations**: Must be handled carefully

### Design Principles
- **Only one primary**: At any given time
- **Backup validation**: Must verify view correctness
- **Atomic state transfer**: All or nothing
- **Deterministic operations**: For reliable replication

### Trade-offs
- **Simplicity vs robustness**: Simple but has single points of failure
- **Performance vs consistency**: Must wait for backup confirmation
- **Transparency vs control**: VM-level replication vs application-level