# Paxos Made Simple: Supplementary Notes

## The Consensus Problem

### Motivation: State Machine Replication
- **Goal**: Agree on order of operations
- **Think of operations** as a log
- **Challenge**: Multiple proposals can conflict

### Why Multiple Proposals?
- **Consensus is easy** if only one client request at a time
- **Solution**: Select a leader
  - Clients send requests to leader
  - Leader picks what goes first, tells everyone else
- **Problem**: Split brain scenarios
  - Old leader slow → two leaders
  - Multiple slow leaders → three or more leaders
  - Each makes proposals for what to do next

### Primary-Backup Challenges
- **How to tell** if primary failed vs slow?
- **If slow**: Might end up with two primaries
- **View server solution**: What if view server goes down?
- **Replicate view server**: How to tell if replica failed vs slow?

## The Part-Time Parliament Analogy

### The Story
- **Parliament determines laws** by passing numbered decrees
- **Legislators can leave/enter** chamber at arbitrary times
- **No centralized record**: Each legislator carries a ledger

### Government Requirements
- **No contradictory information** in ledgers
- **If majority present** and no one enters/leaves for sufficient time:
  - Any decree proposed will eventually be passed
  - Any passed decree appears on every legislator's ledger

### Key Assumptions
- **Legislators are non-partisan, progressive, well-intentioned**
- **Only care that something is agreed to**, not what is agreed to
- **For Byzantine failures**: See Castro and Liskov, SOSP 99

## Paxos System Model

### Environment
- **Set of processes** that can propose values
- **Processes can crash and recover**
- **Access to stable storage**
- **Asynchronous communication** via messages
- **Messages can be lost/duplicated**, but not corrupted

### The Players
- **Proposers**: Propose values
- **Acceptors**: Accept proposals
- **Learners**: Learn chosen values

### Terminology
- **Value**: Possible operation for next slot in operation log
- **Proposal**: To select a value; uniquely numbered
- **Accept**: Of specific proposal/value
- **Chosen**: Proposal/value accepted by majority
- **Learned**: Fact that proposal is chosen is known

## Why Majorities?

### Key Insight
- **Majorities intersect**: For any two majorities S and S', there is some node in both S and S'
- **Ensures consistency**: No two majorities can choose different values
- **Prevents split decisions**: At least one node knows about both decisions

## Consensus Requirements

### Safety Properties
- **Only proposed values** can be chosen
- **Only single value** is chosen
- **Process never learns** value chosen unless it has been

### Liveness Properties
- **Some proposed value** is eventually chosen
- **If value chosen**, process eventually learns it

## Paxos Algorithm Development

### Approach
- **Start broad**: Eventually choose a value, only choose one value
- **Refine definition** to something implementable
- **At each step**: Argue refinement is valid (e.g., P2a ⇒ P2)

### Basic Requirements

#### P1: Accept First Proposal
- **An acceptor must accept** the first proposal it receives
- **Problem**: What if multiple proposers propose different values?

#### P2: Unique Value Chosen
- **If proposal with value v is chosen**, then every higher-numbered proposal that is chosen has value v
- **Ensures consistency**: All chosen proposals have same value

### Refining P2

#### P2a: Accepted Proposals
- **If proposal with value v is chosen**, then every higher-numbered proposal accepted by any acceptor has value v

#### P2b: Issued Proposals
- **If proposal with value v is chosen**, then every higher-numbered proposal issued by any proposer has value v

#### P2c: Implementation Invariant
- **For any v and n**, if proposal with value v and number n is issued, then there is majority set S such that either:
  - **No acceptor in S** has accepted any proposal numbered less than n, OR
  - **v is the value** of highest-numbered proposal among all proposals numbered less than n accepted by acceptors in S

## Proposal Number Assignment

### Requirements
- **Proposal numbers must be unique and infinite**
- **Proposal number server won't work** (single point of failure)
- **Solution**: Assign each proposer infinite slice
- **Proposer i of N gets**: i, i+N, i+2N, i+3N, ...

## The Two-Phase Protocol

### Phase 1: Prepare
1. **Proposer chooses new proposal number n**
2. **Sends prepare request** to majority of acceptors
3. **Acceptor responds** with:
   - Promise never to accept proposals numbered less than n
   - Highest numbered proposal (if any) it has accepted

### Phase 2: Accept
1. **If proposer receives majority responses**:
   - Sends accept request with proposal (n, v)
   - v is value of highest numbered proposal from responses, or any value if no proposals
2. **Acceptor accepts proposal** unless it has responded to prepare request with number > n

### Acceptor Protocol
- **P1a**: Acceptor can accept proposal numbered n iff it has not responded to prepare request having number greater than n
- **Can ignore requests** without affecting safety
- **Must store on stable storage**: Highest numbered proposal accepted, highest numbered prepare request responded to

## Learning Chosen Values

### Strategies
1. **Each acceptor informs each learner** whenever it accepts proposal
2. **Acceptors inform distinguished learner**, who informs others
3. **Something in between**: Set of not-quite-as-distinguished learners

### Failure Handling
- **Because of failures**: Learner may not learn value has been chosen
- **Need retry mechanisms** and multiple learners

## Liveness and Progress

### Progress Not Guaranteed
- **Paxos is always safe** despite asynchrony
- **Once leader elected**, Paxos is live
- **FLP Impossibility**: Leader election impossible in asynchronous system
- **Paxos is next best thing**: Always safe, live during periods of synchrony

## State Machine Replication

### Implementation
- **Sequence of separate consensus instances**
- **i-th instance chooses i-th message** in sequence
- **Each server assumes all three roles** in each instance
- **Set of servers is fixed**

### Role of Leader
- **Elect single server as leader**
- **Leader acts as distinguished proposer** in all instances
- **Clients send commands to leader**
- **Leader decides where** in sequence each command should appear

### New Leader Election
- **New leader λ is elected**
- **λ knows most commands** already chosen (e.g., 1-10, 13, 15)
- **Executes phase 1** of instances 11, 12, 14, and all instances 16+
- **This leaves 14, 16 constrained** and 11, 12, 16+ unconstrained
- **Executes phase 2** of 14 and 16, choosing those commands

### Stop-Gap Measures
- **All replicas execute commands 1-10**
- **Cannot execute 13-16** because 11, 12 haven't been chosen
- **Leader can**:
  - Take next two client commands as 11, 12, OR
  - Propose no-op commands for 11, 12
- **Run phase 2** for instances 11, 12
- **Once consensus achieved**, all replicas can execute through 16

### Efficient Phase 1
- **Leader can efficiently execute phase 1** for infinitely many instances
- **Send message with sufficiently high proposal number** for all instances
- **Acceptor replies non-trivially** only for instances it has already accepted a value

## Key Takeaways

### Paxos Properties
- **Always safe**: Despite asynchrony and failures
- **Live during synchrony**: When leader is elected
- **Uses majorities**: To ensure consistency
- **Two-phase protocol**: Prepare then Accept

### Design Principles
- **Start with safety**: Ensure correctness first
- **Add liveness**: Through leader election
- **Use majorities**: For consistency guarantees
- **Handle failures**: Through stable storage and retries

### Trade-offs
- **Safety vs Liveness**: Always safe, live only with leader
- **Performance vs Correctness**: Two-phase protocol adds latency
- **Simplicity vs Robustness**: Complex but handles all failure modes

### Applications
- **State machine replication**: Ordering operations
- **Distributed databases**: Consensus on transactions
- **Configuration management**: Agreeing on system state
- **Leader election**: Choosing distinguished proposer

### Limitations
- **Requires leader**: For liveness
- **Two-phase protocol**: Adds latency
- **Majority requirement**: Cannot tolerate >50% failures
- **Complex implementation**: Many edge cases to handle
