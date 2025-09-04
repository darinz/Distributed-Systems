# Byzantine Fault Tolerance (BFT): Supplementary Notes

## Byzantine Faults

### Definition
- **Also called**: "General" or "arbitrary" faults
- **Faulty nodes can take any actions**: They can send any messages, collude with each other, etc.
- **Attempt to "trick"**: Non-faulty nodes and subvert the protocol
- **Why this model?**: Real-world systems face malicious attacks

### Characteristics
- **Arbitrary behavior**: Faulty nodes can do anything
- **Collusion**: Faulty nodes can work together
- **Malicious intent**: Actively trying to break the system
- **Unpredictable**: No assumptions about failure behavior

## What About Paxos?

### Paxos Limitations
- **Paxos tolerates**: A minority of processing failing by crashing
- **What could a malicious replica do** to a Paxos deployment?
  - **Stop processing requests**
  - **A leader could report incorrect results** to a client
  - **A follower could acknowledge a proposal** and then discard it
  - **A follower could respond to prepare messages** without all previously acknowledged commands
  - **A server could continually start new leader elections**

### Security Vulnerabilities
- **Malicious behavior**: Not just crashes
- **Incorrect results**: Leaders can lie
- **Message manipulation**: Followers can misbehave
- **Protocol disruption**: Continuous leader elections
- **Need for stronger guarantees**: Beyond crash tolerance

## Setup

### System Model
- **N = 3f + 1 servers**: f of which can be faulty
- **Unlimited clients**: No limit on client count
- **Public-key infrastructure**: Servers and clients can sign messages and verify signatures
- **Signatures aren't forgeable**: Cryptographic security

### Notation
- **Message m with signature**: <m>
- **Message m signed by p**: <m>p
- **Digest function**: D(m) - cryptographic hash, collision-resistant

### Attacker Model
- **Controls f faulty servers**: Knows the protocol the other servers are running
- **Has control over network**: Can delay and reorder messages to all nodes
- **Full knowledge**: Knows all protocol details

## Goal

### State Machine Replication
- **Goal, as in Paxos**: State-machine replication
- **Safety guarantee**: When there are f or fewer failures (or unlimited crash failures)
- **Liveness guarantee**: During periods of synchrony
- **Easy, right?**: Much harder than crash tolerance

### Requirements
- **Safety**: Correctness despite Byzantine failures
- **Liveness**: Progress during synchrony
- **Fault tolerance**: Up to f Byzantine failures
- **Consensus**: All honest nodes agree

## What About Faulty Clients?

### Client Authentication
- **Existing way**: For clients to authenticate themselves with the system
- **Access controls**: Can be used to restrict what each client is allowed to do
- **System administrators**: Can revoke access for faulty clients
- **System itself**: Can revoke access for faulty clients

### Client Management
- **Authentication**: Verify client identity
- **Authorization**: Control what clients can do
- **Revocation**: Remove access for faulty clients
- **Separation of concerns**: Client faults vs. server faults

## Papers, Please

### Proof Requirements
- **Servers don't take each others' word**: They require proof
- **Verify client's command is legitimate**: Need signed message from client (or proof thereof)
- **All other steps**: Taken only after receiving signed messages from quorum of 2f + 1 servers
- **Certificates**: Servers can collect these messages into certificates to prove legitimacy

### Cryptographic Guarantees
- **Signed messages**: From quorum of 2f + 1 servers
- **Certificates**: Proof of legitimacy for certain steps
- **Verification**: All actions require cryptographic proof
- **Trust**: Only in cryptographic signatures

## Protocol Overview

### Three Sub-Protocols
1. **Normal operations**:
   - Phase 1: Pre-prepare
   - Phase 2: Prepare
   - Phase 3: Commit
2. **View change**
3. **Garbage collection**

### Server State
- **Current view**: Current view number
- **State machine checkpoint**: Periodic checkpoints
- **Current state machine state**: Current execution state
- **Log of all not garbage collected messages**: Message history

### Protocol Structure
- **Three phases**: Pre-prepare, Prepare, Commit
- **View changes**: Handle leader failures
- **Garbage collection**: Manage log size

## Accepting Pre-Prepares

### Pre-Prepare Phase
- **Leader receives client request**: Validates and assigns sequence number
- **Broadcasts PRE-PREPARE**: To all followers
- **Followers validate**: Check signature, sequence number, view
- **Accept or reject**: Based on validation

### Validation Rules
- **Signature verification**: Client request must be signed
- **Sequence number**: Must be valid for current view
- **View number**: Must match current view
- **Client authentication**: Must be legitimate client

## Prepare Certificates

### Prepare Phase
- **Followers accept PRE-PREPARE**: Broadcast (signed) PREPARE messages
- **Server receives 2f matching PREPAREs**: Plus associated PRE-PREPARE
- **Prepare Certificate**: Proof of agreement on command

### Certificate Properties
- **Quorums intersect**: At at least one honest server
- **Honest servers don't prepare different commands**: In same slot
- **No two prepare certificates**: For same view, sequence number, different commands
- **Single server not enough**: What about view changes?

### View Change Considerations
- **New leader might not get**: Prepare Certificate
- **Might not have enough information**: To pick correct command in new view
- **Need stronger guarantees**: For view changes

## Commit Certificates

### Commit Phase
- **Server has Prepare Certificate**: Broadcasts COMMIT message
- **Server receives 2f + 1 matching COMMITs**: Plus associated client message
- **Commit Certificate**: Proof of commitment

### Certificate Guarantees
- **Proves every quorum**: Of 2f + 1 servers has at least one non-faulty node with Prepare Certificate
- **Command is now stable**: Will be fixed in same slot future view changes
- **Server can execute**: Command (provided it executed all previous commands)
- **Reply to client**: Send response

### Execution
- **Execute command**: After receiving Commit Certificate
- **Execute all previous commands**: Maintain order
- **Reply to client**: Send response
- **Update state**: Apply command to state machine

## View Change

### Triggering View Change
- **Followers monitor leader**: If leader stops responding to pings or does anything shady
- **Start view change**: Send VIEW-CHANGE messages
- **Stop accepting messages**: For old view

### View Change Process
1. **Follower sends**: <VIEW-CHANGE, v+1, P>p to leader of view v+1
2. **Follower sends**: <VIEW-CHANGE, v+1>p to other followers
3. **P is set**: Of all Prepare Certificates (or Commit Certificates) follower has received
4. **Other followers join**: When they receive f+1 VIEW-CHANGE messages

### View Change Messages
- **VIEW-CHANGE message**: Contains view number and certificates
- **Certificates**: Proof of previous agreements
- **Quorum requirement**: f+1 VIEW-CHANGE messages to start new view

## Starting a New View

### New View Process
- **New leader selected**: Based on view change
- **Collect certificates**: From VIEW-CHANGE messages
- **Determine state**: What commands to execute
- **Broadcast NEW-VIEW**: With new view information

### State Recovery
- **Collect all certificates**: From view change messages
- **Determine committed commands**: Based on certificates
- **Execute missing commands**: In correct order
- **Start normal operation**: In new view

## Garbage Collection

### Normal Case
- **Servers save log**: Of commands and all messages received
- **Periodic compaction**: In non-Byzantine case
- **State transfer**: Bring out-of-date servers up-to-date

### Byzantine Case
- **Server can't just accept**: State transfer from another node
- **Needs proof**: Cryptographic verification
- **Checkpoint mechanism**: Required for garbage collection

### Checkpoint Process
1. **Server decides**: To take a checkpoint
2. **Hashes state**: Of its state machine
3. **Broadcasts**: <CHECKPOINT, n, D(S)>p
   - **n**: Sequence number of last executed command
   - **D(S)**: Hash of the state
4. **Server receives f + 1 CHECKPOINT messages**: Can compact log and discard old protocol messages
5. **Checkpoint Certificate**: Proves validity of state

### Garbage Collection Benefits
- **Log compaction**: Reduce memory usage
- **Discard old messages**: Free up space
- **State verification**: Cryptographic proof of correctness
- **Recovery**: Bring new servers up-to-date

## But What Did That Buy Us?

### Before vs. After
- **Before**: Could only tolerate crash failures
- **PBFT tolerates**: Any failures, as long as only less than a third of servers are faulty
- **What happens if more are faulty?**: System can be compromised
- **However**: PBFT and friends haven't seen wide adoption

### Trade-offs
- **Stronger fault tolerance**: But higher complexity
- **Cryptographic overhead**: But better security
- **More messages**: But stronger guarantees
- **Limited adoption**: But important for security-critical systems

## Performance

### Overhead
- **Extra round of communication**: Adds latency (can be avoided with speculative execution)
- **Committing single operation**: Requires O(n^2) messages (can be improved, though at cost of added latency)
- **Cryptography operations**: Are slow! (though paper describes strategies to speed up using MACs)

### Optimizations
- **Speculative execution**: Avoid extra round
- **Message optimization**: Reduce O(n^2) to lower complexity
- **MACs instead of signatures**: Speed up cryptography
- **Batching**: Multiple operations in one round

### Performance Characteristics
- **Higher latency**: Due to extra rounds
- **More messages**: Due to Byzantine requirements
- **Cryptographic overhead**: Signatures and verification
- **Scalability challenges**: O(n^2) message complexity

## How to Use BFT?

### When to Use
- **Reason to believe**: Number of Byzantine failures will be limited
- **Failures will be independent**: And separated in time
- **Hardware failures**: Probably hold true
- **Security-critical systems**: Where Byzantine tolerance is essential

### Challenges
- **Security flaws**: What about security flaws and software bugs?
- **One possible solution**: n-version programming
- **Independent failures**: Hard to guarantee
- **Cost vs. benefit**: High overhead for uncertain benefit

### Applications
- **Blockchain systems**: Where Byzantine tolerance is critical
- **Security-critical systems**: Where malicious behavior is a concern
- **Distributed ledgers**: Where trust is limited
- **Military systems**: Where adversaries are expected

## Key Takeaways

### BFT Design Principles
- **3f + 1 servers**: Tolerate f Byzantine failures
- **Cryptographic signatures**: All messages must be signed
- **Three-phase protocol**: Pre-prepare, Prepare, Commit
- **View changes**: Handle leader failures
- **Garbage collection**: Manage log size with checkpoints

### Byzantine Fault Model
- **Arbitrary behavior**: Faulty nodes can do anything
- **Collusion**: Faulty nodes can work together
- **Malicious intent**: Actively trying to break system
- **Stronger than crash faults**: Much more challenging

### Protocol Structure
- **Normal operations**: Three-phase commit protocol
- **View changes**: Handle leader failures
- **Garbage collection**: Manage log size
- **Certificates**: Proof of agreement at each phase

### Cryptographic Requirements
- **Public-key infrastructure**: For signatures
- **Message signing**: All messages must be signed
- **Certificate verification**: Proof of agreement
- **Collision-resistant hashes**: For state verification

### Performance Characteristics
- **Higher latency**: Extra communication rounds
- **More messages**: O(n^2) message complexity
- **Cryptographic overhead**: Signatures and verification
- **Scalability challenges**: Due to message complexity

### Trade-offs
- **Stronger fault tolerance**: But higher complexity
- **Better security**: But more overhead
- **Byzantine tolerance**: But limited adoption
- **Cryptographic guarantees**: But performance cost

### When to Use BFT
- **Security-critical systems**: Where Byzantine tolerance is essential
- **Limited Byzantine failures**: Independent and separated in time
- **Hardware failures**: Where Byzantine model applies
- **Blockchain systems**: Where trust is limited

### Limitations
- **Limited adoption**: High complexity and overhead
- **Performance cost**: Higher latency and message complexity
- **Cryptographic overhead**: Signatures and verification
- **Scalability challenges**: O(n^2) message complexity

### Modern Relevance
- **Blockchain systems**: Bitcoin, Ethereum use BFT concepts
- **Distributed ledgers**: Where Byzantine tolerance is critical
- **Security-critical systems**: Military, financial systems
- **Consensus protocols**: Many modern systems use BFT principles

### Lessons Learned
- **Byzantine faults**: Much harder than crash faults
- **Cryptographic proofs**: Essential for Byzantine tolerance
- **Performance trade-offs**: Security vs. performance
- **Limited adoption**: Due to complexity and overhead
- **Specialized applications**: Where Byzantine tolerance is essential
