
# A New Presumed Commit Optimization for Two Phase Commit

## Introduction: Understanding Consistency Models

### Serializability vs. Linearizability

Before diving into Two-Phase Commit optimizations, it's crucial to understand the fundamental consistency models that distributed systems must provide:

#### Serializability

**Definition**: Transactions appear to execute as if they ran one at a time (in some serial order)

**Key characteristics**:
- **Global ordering**: All transactions must be ordered consistently across the entire system
- **Atomicity**: Each transaction either completes entirely or has no effect
- **Isolation**: Concurrent transactions don't interfere with each other

**Example**: If transaction T1 reads account A and writes to account B, and transaction T2 reads account B and writes to account A, the system must ensure that both transactions see a consistent view of the data.

#### Strict Serializability

**Definition**: Serializability with the additional constraint that non-overlapping transactions maintain their temporal order

**Key difference**: Not only must transactions be serializable, but if transaction T1 completes before transaction T2 starts, then T1 must appear to execute before T2 in the serial order.

#### Linearizability

**Definition**: Strict serializability for transactions limited to one operation on one object

**Key characteristics**:
- **Purely local property**: Defined on the history of single objects
- **Non-blocking implementations**: Unlike serializability, linearizability admits non-blocking implementations
- **Always responsive**: A linearizable response always exists regardless of pending operations

**Why this matters**: Linearizability is easier to implement efficiently because it doesn't require global coordination for every operation.

### The Fundamental Challenge

**The core problem**: How do we provide these consistency guarantees in distributed systems where data is spread across multiple nodes?

**The trade-off**: Stronger consistency guarantees (like strict serializability) require more coordination and can impact performance, while weaker guarantees (like linearizability) are more efficient but may not be sufficient for all applications.

## The Banking Example: A Concrete Motivation

### The Transfer Transaction

**Scenario**: You want to move $100 from account A to account B at a bank

**Why this transaction could fail**:
- **Insufficient funds**: Account A doesn't have $100
- **Account deletion**: Either account A or B has been deleted
- **System constraints**: Other business rules prevent the transfer

### Why We Need Serializability

**Without serializability, we could violate critical invariants**:

#### Concurrent Transaction Interference

**The problem**: Concurrent transactions on accounts A and B must not interfere with each other

**Example**: The operation "read B, write B+100" must be equivalent to the atomic operation "B += 100"

**Why this matters**: If two transactions try to modify the same account simultaneously, they could interfere with each other, leading to incorrect results.

#### Global Agreement on Transaction Order

**Critical requirement**: Everyone must agree on exactly what transactions preceded this transfer

**Example 1 - Double spending**:
- Account A has only $100
- Two concurrent transactions: A→B($100) and A→C($100)
- If neither transaction sees the other, both will succeed
- **Result**: A ends up with -$100, violating the "no negative balance" invariant

**Example 2 - Circular transfers**:
- Accounts A and B each try to move $1M to the other (neither has $1M)
- If each transaction thinks the other happened first, both transfers could succeed
- **Result**: Both accounts end up with $1M they didn't have

### Why We Need Recoverability

**Recoverability requirement**: Transaction must be atomic (all or nothing)

**The guarantee**: Debit A if and only if credit B, even with power failures, crashes, etc.

**Why this matters**: If the system crashes after debiting A but before crediting B, the money would disappear. Recoverability ensures that either both operations happen or neither happens.

## Single Database Solution: RDBMS Approach

### How to Achieve These Properties in a Single RDBMS

**Scenario**: Bank ledger is stored in a single relational database management system (RDBMS)

#### Achieving Serializability

**Solution**: Take out locks on both accounts before moving money

**How it works**:
1. **Acquire locks**: Lock account A and account B before performing the transfer
2. **Perform operations**: Debit A and credit B while holding locks
3. **Release locks**: Release locks after transaction completes

**Why this works**: Locks ensure that no other transaction can modify the accounts while the transfer is in progress, preventing interference.

#### Achieving Recoverability

**Solution**: Use a write-ahead log (WAL)

**How it works**:
1. **Log before change**: Write an atomic, idempotent description to the log before changing the B-trees
2. **Apply changes**: Modify the actual data structures
3. **Crash recovery**: Post-crash, log replay has the same effect regardless of B-tree state

**Alternative approaches**:
- **Undo logging**: Keep undo information (used by SQLite)
- **DO-UNDO-REDO logging**: More complex but flexible approach
- **Shadow paging**: Maintain shadow copies of modified pages

#### Handling Lock Failures

**Important note**: Locks can fail or be revoked, in which case the RDBMS aborts the transaction

**What happens on abort**:
- **No effect**: Aborted transaction has no effect on database contents
- **Client notification**: Must be reported to client
- **Retry or give up**: Client can try again or give up

**Why locks can fail**:
- **Deadlock detection**: System detects potential deadlocks and aborts one transaction
- **Timeout**: Locks can timeout if held too long
- **Resource pressure**: System may revoke locks under memory pressure

## The Distributed Challenge: Sharded Databases

### The Problem with Multiple Databases

**Scenario**: Bank ledger is large and must be sharded, with accounts A and B in separate RDBMSes

**The naive approach**: Run two separate transactions on two separate RDBMSes

**What goes wrong?** Both serializability and recoverability are violated.

### Serializability Problems

#### Example: Concurrent Payments

**Scenario**: Two concurrent payments
- **T1**: A→B($10) 
- **T2**: C→A($100)

**Data distribution**: A in DB1, B and C in DB2

#### The Ordering Problem

**Execution sequence**:
1. **T2**: Lock C (in DB2)
2. **T1**: Lock B (in DB2), lock A (in DB1), commit, release locks
3. **T2**: Lock A (in DB1), commit, release locks

**The problem**: Different databases see the transactions in different orders

**DB1's view**: T1 before T2 (because T1 acquired lock on A first)
**DB2's view**: T2 before T1 (because T2 acquired lock on C first)

**Why this is bad**: DB2 can order T2 before concurrent T1 (still linearizable locally) but this creates global inconsistency.

#### When Operations Don't Commute

**The real problem**: This ordering inconsistency is bad when operations do not commute

**Example**: Account A assessed overdraft fee even though C debited before B credited
- **DB1 thinks**: T1 (A→B) happened before T2 (C→A)
- **DB2 thinks**: T2 (C→A) happened before T1 (A→B)
- **Result**: A might be charged an overdraft fee even though C's payment should have covered it

### Recoverability Problems

**The critical question**: What if T1 commits on DB2 but aborts on DB1?

**The problem**: 
- **DB2**: T1 successfully credits B with $10
- **DB1**: T1 fails to debit A with $10
- **Result**: $10 appears out of nowhere, violating the fundamental accounting principle

**Why this happens**: Without coordination between databases, there's no way to ensure that both operations succeed or both fail.

## Two-Phase Commit: The Solution

### The 2PC Protocol

**Solution**: Use Two-Phase Commit (2PC) to ensure *all* databases either commit or abort

**The process**:
1. **Acquire locks**: Acquire a bunch of locks, finish your transaction, decide to commit
2. **Phase 1**: Preparing
3. **Phase 2**: Committing/aborting

### Phase 1: Preparing

**Coordinator broadcasts PREPARE to all cohorts** (databases in the transaction)

**Key points**:
- **Up to this point**: Each cohort can abort (e.g., if it revoked a lock)
- **If cohort aborted**: Responds with ABORT-VOTE
- **Otherwise**: Respond with COMMIT-VOTE: now can't abort or touch locks!

**The critical transition**: Once a cohort votes COMMIT-VOTE, it has made an irrevocable promise to commit the transaction if the coordinator decides to commit.

### Phase 2: Committing/Aborting

**Decision making**:
- **If any cohort voted ABORT-VOTE**: Coordinator broadcasts ABORT message
- **If all voted COMMIT-VOTE**: Coordinator broadcasts COMMIT message

**Cohort responsibilities**:
- **If COMMIT**: Cohorts *must* commit the transaction
- **Lock release**: Only now can cohorts release any locks associated with the transaction
- **Acknowledgment**: Cohorts reply with ACK so coordinator knows they received the outcome

### Does 2PC Solve Our Problems?

#### Recoverability: Yes (if everything is logged)

**Why it works**: The two-phase protocol ensures that either all databases commit or all abort, preventing the "money appears out of nowhere" problem.

**Requirements**: All participants must log their decisions to ensure they can recover from crashes.

#### Serializability: Possibly

**The challenge**: Even with 2PC, serializability is not automatically guaranteed.

**Example problem**:
- **DB2 must commit T1** before learning whether or not T2 committed
- **But T2 could still have a lower timestamp** than T1 on DB2
- **Result**: Global ordering is still not guaranteed

**Solution**: Assign timestamps in Phase 2 (e.g., on receipt of PREPARE)
- **Why this works**: At that point, all associated locks are held on all cohorts
- **Result**: Timestamp assignment happens when all participants are locked and committed to the transaction

## Coordinator State Management

### In-Memory Protocol Database

**What must the coordinator store?** (See Figure 1 in the paper)

#### Transaction-Level State

**Tid**: Transaction ID
- **Purpose**: Unique identifier for the transaction
- **Scope**: Global across all participants

**Stable**: Yes if existence of transaction has been persistently logged
- **Purpose**: Indicates whether the transaction's existence is durably recorded
- **Importance**: Critical for crash recovery

**State**: Initiated|Preparing|Aborted|Committed
- **Initiated**: Transaction has started but not yet in prepare phase
- **Preparing**: Currently in Phase 1, waiting for votes
- **Aborted**: Transaction has been aborted
- **Committed**: Transaction has been committed

#### Per-Cohort State

**Per-cohort information**:
- **Cohort-id**: Identifier for each participating database
- **Vote**: None|Abort|Read-only|Commit
  - **None**: No vote received yet
  - **Abort**: Cohort voted to abort
  - **Read-only**: Cohort is read-only (optimization)
  - **Commit**: Cohort voted to commit
- **Ack**: Whether acknowledgment has been received

#### State Cleanup

**Memory management**:
- **Per-cohort state**: Can delete per-cohort state on receiving ACK
- **Transaction state**: Can delete Tid entirely when all ACKs received

**Why this matters**: The coordinator needs to track the state of each transaction and each participant to ensure proper protocol execution and crash recovery.

## When Must Cohorts Force-Write to Disk?

### Before Sending COMMIT-VOTE: Always

**Why this is mandatory**:
- **COMMIT-VOTE is a promise**: Not to abort transaction or release locks
- **Crash recovery**: If cohort reboots with no record of promise, can't possibly keep it

**Example scenario**: Before hearing from coordinator of transaction it forgot about, might agree to commit conflicting transaction with other coordinator

**The problem**: Without logging the COMMIT-VOTE, a crashed cohort could:
1. Reboot and forget about the transaction
2. Receive a conflicting transaction from another coordinator
3. Commit the conflicting transaction, violating the promise made in the COMMIT-VOTE

### Before Sending ACK: Depends on Variant

**Naive approach**: Yes, because it should permanently commit/abort transaction

**But there are alternatives**:
- **Crash recovery**: If it crashes, will know transaction existed from COMMIT-VOTE record
- **Coordinator memory**: If coordinator remembers what happened, cohort can just ask it

**Whether an ACK message is even required** (and associated write) depends on:
- **What information coordinator retains**: Does coordinator remember the transaction outcome?
- **Whether transaction committed or aborted**: Different variants handle these differently

### Performance Implications

**How much do we care about these forced writes?**

#### COMMIT-VOTE: Critical Path

**COMMIT-VOTE is on the critical path for transaction latency**:
- **Blocking**: Transaction cannot proceed until COMMIT-VOTE is logged
- **Performance impact**: Disk latency directly affects transaction latency
- **Optimization priority**: This is the most important write to optimize

#### ACK: Non-Critical Path

**ACK is not on the critical path**:
- **Coordinator already knows**: Transaction committed
- **No latency impact**: Disk latency won't affect transaction latency
- **Optimization opportunities**: 
  - **Delay**: Maybe delay the write
  - **Piggyback**: Piggyback on another log write for better throughput
  - **Batching**: Batch multiple ACK writes together

#### Paper's Analysis

**Possibly unfair of paper to lump these together in single "n" value**:
- **Different importance**: COMMIT-VOTE and ACK have very different performance implications
- **Different optimization strategies**: What works for one might not work for the other
- **Misleading metrics**: Combining them might hide important performance characteristics

## Coordinator Force-Write Requirements: PrN Variant

### When Must Coordinator Force-Write in PrN (Presume Nothing)?

#### Before Sending PREPARE: No (only if pedantically presuming nothing)

**Why not required**:
- **Unilateral abort**: Until it sends COMMIT, coordinator can unilaterally ABORT a transaction
- **Presume abort**: Cohort inquires about unknown transaction? Coordinator "presumes" abort

**The logic**: Since the coordinator can always abort the transaction before sending COMMIT, it doesn't need to log the PREPARE message. If a cohort asks about an unknown transaction, the coordinator can safely assume it was aborted.

#### Before Sending COMMIT: Yes

**Why this is required**:
- **Transaction happened**: At this point transaction happened, so need to record it durably
- **Crash recovery**: Otherwise, couldn't properly respond to cohort inquiries after crash

**The problem**: Once the coordinator sends COMMIT, the transaction is committed. If the coordinator crashes after sending COMMIT but before logging it, it won't know the transaction was committed when cohorts ask about it.

#### Before Sending ABORT: Yes

**Why this is required**:
- **Presume nothing**: This is the part about presume nothing
- **Crash recovery**: Need to remember that the transaction was aborted

**The logic**: If the coordinator doesn't log the ABORT, it might forget that the transaction was aborted and make incorrect decisions about related transactions.

#### Upon Receiving ACK: No (non-forced write)

**Why not required**:
- **Non-critical**: ACK is not on the critical path
- **Optimization opportunity**: Can be delayed or batched for better performance

**The reasoning**: Since the coordinator already knows the transaction outcome, receiving ACKs is just for cleanup purposes and doesn't need to be forced to disk immediately.

## 2PC Variants Comparison Table

### Force-Write Requirements by Variant

**Fill in this table during lecture**:

| Operation | PrN | PrA | PrC | NPrC |
|-----------|-----|-----|-----|------|
| log before PREPARE | N* | N | Y | N |
| log before COMMIT | Y | Y | Y | Y |
| log before ABORT | Y* | N | N | N |
| ACK after COMMIT | Y | Y | N | N |
| ACK after ABORT | Y | N | Y | Y |

**Note**: Cohorts log before COMMIT-VOTE in all schemes

**Legend**:
- **PrN**: Presume Nothing
- **PrA**: Presume Abort  
- **PrC**: Presume Commit
- **NPrC**: New Presume Commit (the paper's contribution)

**Key observations**:
- **Y**: Yes, force-write required
- **N**: No, force-write not required
- **N***: Depends how pedantic we want to be about presuming nothing

**The table shows**: Different 2PC variants have different logging requirements, which directly impact performance characteristics.

## Presumed Abort (PrA) vs. Presume Nothing (PrN)

### Key Difference: Don't Write to Disk Before Sending ABORT

**PrA optimization**: Coordinator never writes aborted transactions to disk

#### Benefits of Not Logging Aborts

**No garbage to collect**:
- **Clean state**: No matter when coordinator crashes, nothing on disk
- **Consistent responses**: Cohorts inquiring about transaction will always get same answer
- **Simplified recovery**: No need to clean up aborted transaction records

**Reduced message overhead**:
- **No ACK required**: Cohorts don't even need to send ACK message to ABORT
- **No force-write**: Obviously no force-write before (non-existent) ACK message
- **Lower latency**: Fewer messages and disk writes for aborted transactions

#### Performance Characteristics

**COMMITs (common case) are exactly same cost as PrN**:
- **No optimization for commits**: PrA doesn't improve the performance of committed transactions
- **Focus on aborts**: The optimization specifically targets the less common abort case
- **Trade-off**: Better abort performance, same commit performance

**Why this matters**: In most systems, commits are much more common than aborts, so optimizing the abort path has limited overall impact.

## Traditional Presumed Commit (PrC) Costs

### Coordinator Write Requirements

#### Before PREPARE: Yes

**Why this is required**:
- **Crash recovery**: Otherwise, after crash would presume committed if cohort inquired
- **Asymmetry**: Why the difference from PrA? Asymmetry in coordinator capabilities

**The asymmetry**:
- **Without all votes**: Coordinator can unilaterally ABORT but not COMMIT
- **Commit requires consensus**: All cohorts must vote COMMIT for transaction to commit
- **Abort is unilateral**: Coordinator can abort without waiting for all votes

#### Before COMMIT: Yes

**Why this is required**:
- **Undo effect**: Have to log to "undo" effect of previously written PREPARE log record
- **Crash recovery**: Otherwise, would see transaction after crash and abort it
- **Consistency**: Need to ensure coordinator remembers the commit decision

**But cohort doesn't have to ACK COMMIT**:
- **Optimization**: Since coordinator logged the commit, cohorts don't need to acknowledge
- **Reduced overhead**: Fewer messages and disk writes

#### Before ABORT: No

**Why not required**:
- **Presume commit**: System presumes transactions are committed by default
- **Abort is exception**: Aborted transactions are the exception, not the rule

**But cohorts must ACK ABORT**:
- **Cleanup required**: Need to ensure all cohorts know the transaction was aborted
- **Additional overhead**: Coordinator has one more non-forced cleanup write when all ACKs are in

### Performance Characteristics

**PrC is optimized for the common case (commits)**:
- **Commit optimization**: Fewer writes and messages for committed transactions
- **Abort penalty**: More overhead for aborted transactions
- **Trade-off**: Better performance for commits, worse performance for aborts

## The Read-Only Optimization

### When Transactions Are Read-Only at a Cohort

**Scenario**: Transaction might be read-only at a cohort

**What this means**:
- **Only effect**: Hold locks for duration of transaction
- **No modifications**: Cohort doesn't modify any data
- **Indifference**: So cohort doesn't care whether transaction commits or aborts

### How Read-Only Optimization Works

#### Cohort Response

**Cohort replies to PREPARE with READ-ONLY-VOTE**:
- **Special vote**: Indicates the cohort is read-only
- **Lock release**: Also releases all locks when it sends READ-ONLY-VOTE
- **Safety**: Any locked data unmodified since transaction read-only at cohort

#### Coordinator Behavior

**Coordinator doesn't send COMMIT/ABORT message to read-only cohort**:
- **No need**: Read-only cohorts don't need to know the final outcome
- **Reduced messages**: Fewer messages in the system
- **Faster completion**: Read-only cohorts can finish immediately

### All Read-Only Transaction

**If all cohorts reply READ-ONLY-VOTE, then whole transaction read-only**:

#### Logging Optimization

**If coordinator didn't log PREPARE message, doesn't need to log anything**:
- **No persistence needed**: Cohorts don't care whether committed or aborted
- **Maximum optimization**: No disk writes required at all

#### PrC Special Case

**With PrC, must write (non-forced) record to delete logged PREPARE message**:
- **Cleanup required**: Since PrC logs PREPARE, need to clean it up
- **Non-forced write**: Can be delayed or batched for better performance
- **Overhead**: Additional write compared to other variants

### Performance Benefits

**Read-only optimization provides**:
- **Faster completion**: Read-only cohorts finish immediately
- **Reduced messages**: No COMMIT/ABORT messages to read-only cohorts
- **Lock efficiency**: Locks released as soon as possible
- **Logging optimization**: Minimal or no logging required

## The NPrC (New Presumed Commit) Optimization

### Core Idea: Trade Garbage Collection for Performance

**The fundamental trade-off**: Trade garbage collection after crash to reduce messages and writes

### The Transaction ID Window

**Window of recent transaction ids**: REC=(tid_l,...,tid_h) presumed aborted
- **All other transactions**: Presumed committed
- **Cohorts act exactly like normal presumed committed**

### Window Parameter Management

#### Must Stably Log Window Parameters (tid_l,tid_h)

**tid_h (high water mark)**:
- **Can be implicit**: Based on highest stable transaction
- **Or explicit**: Can be explicitly logged but amortized over many transactions
- **Constraint**: Must guarantee no tids >= tid_h ever used

**tid_l (low water mark)**:
- **Piggybacked**: Onto other log writes
- **Purpose**: Points to oldest "undocumented" transaction
- **Definition**: Oldest that is Initiated or Preparing or Aborted with missing ACKs

### Performance Optimizations

#### No Need to Log Before PREPARE (Common Case)

**When in window**: No need to log before sending PREPARE when in window (common case)

**Flexibility**: Can always log PREPARE of slow ("recalcitrant") transaction later
- **Logging PREPARE**: Allows tid_l to advance despite a few stragglers
- **Adaptive**: System can handle both fast and slow transactions efficiently

#### No Need for Cohorts to ACK COMMIT Messages

**Why this works**: Because tid_l suffices for coordinator to "remember" an arbitrary number of committed transactions cohorts might ask about

**The insight**: The window mechanism provides enough information for the coordinator to answer inquiries about committed transactions without needing explicit ACKs.

### Result: Minimal Forced Writes

**Only common-case forced log write is before sending COMMIT**:
- **Maximum optimization**: Minimal disk writes for the common case
- **Performance benefit**: Reduced latency and improved throughput
- **Trade-off**: More complex garbage collection after crashes

## NPrC Crash Recovery Scenarios

### Coordinator Crashes While Committing

**What happens if NPrC coordinator crashes while committing tid?**

#### Case 1: tid >= tid_l

**If tid >= tid_l, then disk will contain COMMIT log record**:
- **Logged**: The COMMIT was logged before being sent
- **Recovery**: Coordinator can respond correctly to cohort inquiries
- **Consistency**: System maintains correct state

#### Case 2: tid < tid_l

**If tid < tid_l, then will presume committed, so no need for record**:
- **Presumption**: System presumes the transaction was committed
- **No logging needed**: No COMMIT record required on disk
- **Efficiency**: Reduces logging overhead for fast transactions

### Coordinator Crashes While Aborting

**What happens if NPrC coordinator crashes while aborting tid?**

**Won't advance tid_l > tid until all cohorts ACK the abort**:
- **Constraint**: tid_l cannot advance past tid until abort is fully acknowledged
- **Safety**: Ensures system can handle inquiries about the aborted transaction

**Two possible outcomes**:
1. **tid > tid_l and presumed aborted**: Transaction is in the abort window
2. **No cohorts will inquire**: All cohorts have already acknowledged the abort

**Why this works**: The system ensures that either the transaction is in the abort window (and will be presumed aborted) or all cohorts have already been notified of the abort (so no inquiries will come).

## The Garbage Collection Issue

### The Problem with Never Logging PREPARE

**What is the garbage collection issue?**

**Never logged PREPARE? Won't know what cohorts involved in transaction**:
- **Missing information**: Can't collect ACKs for ABORT
- **Inquiry handling**: Have to be prepared for inquiries
- **No time-bound**: No time-bound on unknown cohorts inquiring about unknown transactions

**The core problem**: Without logging PREPARE, the coordinator doesn't know which cohorts participated in the transaction, making it impossible to properly clean up aborted transactions.

### The Solution: Permanent Records

**Solution is keep permanent record of presumed abort ranges after each crash**:
- **Persistent state**: Maintain records of which transaction ranges were presumed aborted
- **Inquiry handling**: Use these records to answer cohort inquiries
- **Cleanup**: Eventually clean up these records when safe

### Alternative Approaches

**What if you really don't want "forever garbage" in your system?**

#### Track Active Cohorts

**Could alternatively keep track of all possibly active cohorts**:
- **Cohort registry**: Maintain a list of all active cohorts
- **Crash notification**: Inform all cohorts of crashes with CRASH messages
- **Active management**: Keep track of which cohorts are still alive

#### Handle Dead Cohorts

**But what if old cohort is permanently dead? Manually remove dead cohorts**:
- **Manual intervention**: System administrator must manually remove dead cohorts
- **Operational overhead**: Requires human intervention
- **Risk**: Risk of removing cohorts that are actually alive but temporarily unreachable

### Trade-offs

**The fundamental trade-off**:
- **NPrC**: Better performance, but more complex garbage collection
- **Traditional approaches**: Simpler garbage collection, but worse performance
- **Choice**: Depends on system requirements and operational capabilities

## Constraint Checking and Read-Only Cohorts

### The Constraint Checking Problem

**With simple transactions, tx done and all locks held when PREPARE sent**

**Consider the following scenario**:
- **Database constraint**: SUM(colX) < 1000
- **Deferred checking**: Want to defer constraint checking to when cohort receives PREPARE
- **Additional locks**: Checking constraint requires additional read locks on colX

### The Read-Only Cohort Problem

**What can go wrong with constraints if a cohort has READ-ONLY-VOTE?**

#### The Issue

**RO cohort never gets COMMIT/ABORT from coordinator**:
- **Immediate release**: So releases locks immediately on sending READ-ONLY-VOTE
- **Timing problem**: But that might be before other cohort acquires lock on colX
- **Result**: Transactions no longer guaranteed serializable!

#### The Race Condition

**Example scenario**:
1. **Cohort A**: Read-only, releases locks immediately after READ-ONLY-VOTE
2. **Cohort B**: Needs to check constraint, requires locks on colX
3. **Race**: Cohort A releases locks before Cohort B acquires them
4. **Problem**: Constraint checking happens without proper locking

### Solution: Timestamp Transactions

**One solution: timestamp transactions**:
- **Order by timestamps**: Ensure same order on all cohorts
- **ABORT if inconsistent**: ABORT if no timestamp reflects order on all cohorts

**How it works**:
1. **Assign timestamps**: Each transaction gets a timestamp
2. **Global ordering**: All cohorts must see transactions in timestamp order
3. **Consistency check**: If timestamps don't reflect the same order on all cohorts, abort
4. **Serializability**: Ensures global serializability despite read-only optimizations

**Benefits**:
- **Maintains serializability**: Despite read-only optimizations
- **Handles constraints**: Allows proper constraint checking
- **Performance**: Still benefits from read-only optimizations

## NPrC and User-Visible Timestamps

### Why NPrC is Bad for User-Visible Timestamps

**Why is NPrC bad for transactions with user-visible timestamps?**

#### The Problem

**Coordinator chooses a timestamp, sends it to cohorts in COMMIT message**:
- **Inquiry handling**: A cohort inquires after missing COMMIT must be told the timestamp
- **Force-write required**: Since coordinator must force-write timestamp, might as well use PrA

**The issue**: NPrC's main optimization is avoiding force-writes, but user-visible timestamps require force-writes anyway, negating the benefit.

#### Why This Defeats NPrC

**NPrC's advantage**: Minimal force-writes for the common case
**User-visible timestamps**: Require force-writes for every commit
**Result**: NPrC provides no benefit over PrA when timestamps are user-visible

### Internal Timestamps for Read-Only Optimization

**But timestamps may be just for RO optimization, not user-visible (see [5])**

#### Is This Compatible with NPrC? Yes.

**Cohort's COMMIT-VOTE specifies range of permissible timestamps**:
- **Range specification**: Guarantees no overlapping ranges for conflicting transactions
- **Sufficient information**: So if commit, cohort knows timestamp in valid range, which is good enough

#### How This Works

**The process**:
1. **Cohort specifies range**: In COMMIT-VOTE, cohort specifies valid timestamp range
2. **No conflicts**: System ensures no overlapping ranges for conflicting transactions
3. **Commit decision**: If transaction commits, cohort knows timestamp is in valid range
4. **No force-write needed**: Coordinator doesn't need to force-write the exact timestamp

**Benefits**:
- **NPrC compatibility**: Works with NPrC's minimal force-write approach
- **Read-only optimization**: Still enables read-only optimizations
- **Performance**: Maintains NPrC's performance benefits

**The key insight**: Internal timestamps for optimization don't need to be user-visible, so they don't require the same force-write guarantees as user-visible timestamps.