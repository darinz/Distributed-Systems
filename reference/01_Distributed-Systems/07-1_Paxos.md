# Paxos: The Foundation of Distributed Consensus

## The Fundamental Problem: Consensus in Distributed Systems

In distributed systems, we face a fundamental challenge: **how do we get multiple machines to agree on something?** This might seem simple at first, but it's actually one of the most difficult problems in computer science. The challenge becomes even more complex when we consider that machines can fail, networks can be slow, and messages can be lost.

This document explores Paxos, one of the most important algorithms in distributed systems. Paxos provides a solution to the consensus problem that is both safe and live under the right conditions. Understanding Paxos is essential for anyone working with distributed systems.

### The Challenge: Why Consensus Is Hard

**The Problem**: In a distributed system, we need multiple machines to agree on a single value or decision.

**Why This Is Hard**:
- **Failures**: Machines can crash and recover
- **Network Issues**: Messages can be lost, delayed, or duplicated
- **Concurrency**: Multiple processes might try to propose different values simultaneously
- **Partial Information**: No single machine sees everything that's happening

**The Consequence**: Without proper coordination, machines might make conflicting decisions, leading to inconsistency and system failure.

### The Solution: State Machine Replication

**The Key Insight**: We can solve many distributed systems problems by having all machines execute the same sequence of operations in the same order.

**State Machine Replication**: All machines start in the same initial state and apply the same sequence of operations in the same order.

**The Result**: All machines end up in the same final state, ensuring consistency.

### The Operation Log: A Simple Abstraction

**The Idea**: Think of operations as a log that all machines must agree on.

**The Log**: Op1, Op2, Op3, Op4, Op5, Op6...

**The Challenge**: All machines must agree on what goes in each position of the log.

**The Power**: Once we agree on the log, we can execute the operations in order to maintain consistency.

### A Concrete Example: Distributed Key-Value Store

**The Scenario**: Three servers (S1, S2, S3) need to agree on the order of operations.

**The Operations**: 
- "Put k1 v1" (set key k1 to value v1)
- "Put k2 v2" (set key k2 to value v2)

**The Challenge**: Multiple clients might send requests simultaneously, and we need to agree on the order.

**The Solution**: Use Paxos to agree on the order of operations in the log.

### The Paxos Process: Step by Step

**Step 1**: Clients send requests to the system
- Client A wants to do "Put k1 v1"
- Client B wants to do "Put k2 v2"

**Step 2**: Use Paxos to agree on the first operation (Op1)
- All servers must agree on what goes in position 1 of the log
- Paxos ensures that only one value is chosen for Op1

**Step 3**: Use Paxos to agree on the second operation (Op2)
- All servers must agree on what goes in position 2 of the log
- Paxos ensures that only one value is chosen for Op2

**Step 4**: Continue this process for all operations
- Each position in the log requires a separate consensus decision
- Paxos ensures that all servers agree on each position

### The Result: Consistent State

**The Final State**: All servers have the same log:
- Op1: Put k1 v1
- Op2: Put k2 v2

**The Consistency**: All servers execute the same operations in the same order, ensuring they end up in the same state.

### The Challenge: Multiple Proposals

**The Problem**: What if multiple clients want to propose different operations for the same position in the log?

**The Scenario**: 
- Client A wants to do "Put k1 v1"
- Client B wants to do "Put k2 v2"
- Both want their operation to be Op1

**The Challenge**: We need to ensure that only one operation is chosen for each position.

### The Leader Solution: A Simple Approach

**The Idea**: Select a leader to coordinate decisions.

**The Process**:
- Clients send requests to the leader
- Leader picks what goes first, tells everyone else
- All servers follow the leader's decisions

**The Problem**: What about split brain? (leader failed, or slow)
- If old leader is slow, might have two leaders!
- If old and new leader are slow, might have three!
- Each makes a proposal for what to go next

### The Primary-Backup Problem: A Deeper Challenge

**The Scenario**: Suppose using primary/hot standby replication.

**The Challenge**: How can we tell if primary has failed versus is slow?
- If slow, might end up with two primaries!
- This leads to split-brain and inconsistency

**The View Server Solution**: Rely on view server to decide?

**The New Problem**: What if view server goes down? Replicate?

**The Deeper Problem**: How can we tell if view server replica has failed or is slow?

**The Infinite Regress**: We need consensus to elect a leader, but we need a leader to achieve consensus!

### The Fundamental Insight

**The Key Realization**: Consensus is a fundamental problem that cannot be solved by simply adding more layers of coordination.

**The Elegance**: Paxos provides a solution that works even in the presence of failures and network issues.

**The Impact**: Understanding Paxos is essential for building reliable distributed systems.

### The Journey Ahead

This document will take you through the complete story of Paxos:

1. **The Part-Time Parliament**: A beautiful analogy that explains the core concepts
2. **The Consensus Problem**: Formal definition of what we're trying to solve
3. **The Paxos Algorithm**: How to achieve consensus in practice
4. **Phase 1 and Phase 2**: The two-phase protocol that makes Paxos work
5. **Leader Election**: How to make Paxos practical in real systems
6. **Practical Implications**: How to use Paxos in practice

By the end, you'll understand not just how Paxos works, but why it's so important for distributed systems.

### The End of the Beginning

Paxos represents one of the most elegant solutions to one of the most fundamental problems in distributed systems. The algorithm is both simple and powerful, providing a way to achieve consensus even in the face of failures and network issues.

The journey from the simple idea of state machine replication to the sophisticated Paxos algorithm shows how complex problems can be solved with careful thought and elegant design. The challenge is always the same: how do you get multiple machines to agree on something?

Paxos provides an answer to this question, and it continues to influence how we build distributed systems today.
## The Part-Time Parliament: A Beautiful Analogy

Now let's explore the Part-Time Parliament, a beautiful analogy that explains the core concepts of Paxos. This analogy, created by Leslie Lamport, makes the complex consensus problem much easier to understand.

### The Setting: A Part-Time Parliament

**The Parliament**: Determines laws by passing a sequence of numbered decrees.

**The Legislators**: Can leave and enter the chamber at arbitrary times.

**The Challenge**: No centralized record of approved decrees—instead, each legislator carries a ledger.

**The Problem**: How do we ensure that all legislators agree on what decrees have been passed?

### The Core Requirements: Government 101

**Requirement 1**: No two ledgers contain contradictory information.

**What This Means**: If one legislator's ledger says "Decree 5: Raise taxes," then no other legislator's ledger can say "Decree 5: Lower taxes."

**Why This Matters**: Contradictory information would lead to chaos and inconsistency.

**Requirement 2**: If a majority of legislators were in the chamber and no one entered or left the chamber for a sufficiently long time, then:
- Any decree proposed by a legislator would eventually be passed
- Any passed decree would appear on the ledger of every legislator

**What This Means**: The system should make progress when conditions are right.

**Why This Matters**: Without progress, the system would be useless.

### The Paxos Legislature: Government 102

**The Paxos Legislature**: Is non-partisan, progressive, and well-intentioned.

**What This Means**: Legislators only care that something is agreed to, not what is agreed to.

**The Implication**: The focus is on reaching consensus, not on the specific content of the consensus.

**The Byzantine Problem**: To deal with Byzantine legislators (those who might lie or behave maliciously), see Castro and Liskov, SOSP 99.

**The Focus**: For now, we assume legislators are honest but might fail.

### The Analogy to Distributed Systems

**Legislators**: Represent processes in a distributed system.

**The Chamber**: Represents the network where processes can communicate.

**Decrees**: Represent values that processes need to agree on.

**Ledgers**: Represent the state that each process maintains.

**Majority**: Represents the quorum needed to make decisions.

### The Key Insights from the Analogy

**Insight 1**: We need a majority to make decisions.

**Why**: A majority ensures that any two majorities will have at least one legislator in common.

**The Result**: This prevents contradictory decisions from being made.

**Insight 2**: Legislators can come and go.

**Why**: In distributed systems, processes can fail and recover.

**The Result**: The system must work even when some processes are unavailable.

**Insight 3**: Each legislator maintains their own ledger.

**Why**: In distributed systems, there's no central authority.

**The Result**: Each process must maintain its own state and coordinate with others.

### The Power of the Analogy

**What It Provides**:
- **Intuitive Understanding**: The analogy makes complex concepts easier to grasp
- **Clear Requirements**: The government requirements map directly to consensus requirements
- **Practical Insights**: The analogy reveals why certain design choices are necessary

**Why This Matters**:
- **System Design**: We can use the analogy to guide our system design
- **Correctness**: We can verify that our system satisfies the requirements
- **Debugging**: We can use the analogy to understand what went wrong

### The Fundamental Insight

**The Key Realization**: The Part-Time Parliament analogy reveals the essential requirements for consensus.

**The Elegance**: The analogy makes complex distributed systems concepts accessible and intuitive.

**The Result**: We can use this understanding to build consensus algorithms that work in practice.

### The Journey Forward

Now that we understand the Part-Time Parliament analogy, we can explore the formal consensus problem. The next section will show us how to translate the intuitive requirements into formal mathematical properties.

The key insight is that the Part-Time Parliament analogy provides a powerful way to understand the consensus problem, making complex concepts accessible and intuitive.
## The Consensus Problem: Formal Definition

Now let's translate the intuitive requirements from the Part-Time Parliament analogy into a formal mathematical definition of the consensus problem.

### The System Model

**The Setting**: A set of processes that can propose values.

**The Assumptions**:
- **Processes can crash and recover**: Processes can fail at any time and come back later
- **Processes have access to stable storage**: Processes can persist information across crashes
- **Asynchronous communication via messages**: Processes communicate by sending messages
- **Messages can be lost and duplicated, but not corrupted**: The network is unreliable but not malicious

**The Challenge**: How do we achieve consensus in this environment?

### The Players: Three Types of Processes

**Proposers**: Processes that propose values for consensus.

**What They Do**: Proposers suggest values that the system should agree on.

**Acceptors**: Processes that vote on proposals.

**What They Do**: Acceptors decide whether to accept or reject proposals.

**Learners**: Processes that learn about chosen values.

**What They Do**: Learners discover what values have been chosen by the system.

**The Key Insight**: In practice, a single process can play multiple roles.

### The Terminology: Key Concepts

**Value**: A possible operation to put in the next slot in the operation log (letter values).

**What This Means**: A value is what we want to agree on (e.g., "Put k1 v1").

**Proposal**: To select a value; proposals are uniquely numbered.

**What This Means**: A proposal is a specific attempt to get a value chosen, with a unique identifier.

**Accept**: Of a specific proposal, value.

**What This Means**: An acceptor accepts a proposal when it votes in favor of it.

**Chosen**: Proposal/value accepted by a majority.

**What This Means**: A proposal is chosen when a majority of acceptors have accepted it.

**Learned**: Fact that proposal is chosen is known.

**What This Means**: A process learns a value when it discovers that the value has been chosen.

### The Power of Majorities

**Why does Paxos use majorities?**

**The Key Property**: Majorities intersect: for any two majorities S and S', there is some node in both S and S'.

**What This Means**: Any two majorities will always have at least one process in common.

**Why This Matters**: This property prevents contradictory decisions from being made.

**The Intuition**: If two majorities could be disjoint, they could make conflicting decisions without knowing about each other.

**The Result**: Using majorities ensures that the system can never make contradictory decisions.

### The Game: Consensus Requirements

**SAFETY PROPERTIES** (What should never happen):

**Safety 1**: Only a value that has been proposed can be chosen.

**What This Means**: The system cannot make up values; it can only choose from values that were actually proposed.

**Safety 2**: Only a single value is chosen.

**What This Means**: The system cannot choose multiple different values for the same consensus decision.

**Safety 3**: A process never learns that a value has been chosen unless it has been.

**What This Means**: Processes cannot learn about values that were never actually chosen.

**LIVENESS PROPERTIES** (What should eventually happen):

**Liveness 1**: Some proposed value is eventually chosen.

**What This Means**: The system should make progress and eventually reach a decision.

**Liveness 2**: If a value is chosen, a process eventually learns it.

**What This Means**: Once a value is chosen, all processes should eventually discover this fact.

### The Relationship to the Part-Time Parliament

**Safety Properties**: Map to "No two ledgers contain contradictory information."

**Liveness Properties**: Map to "Any decree proposed by a legislator would eventually be passed."

**The Connection**: The formal requirements capture the intuitive requirements from the analogy.

### The Fundamental Insight

**The Key Realization**: The consensus problem can be precisely defined using safety and liveness properties.

**The Elegance**: These properties capture exactly what we need for consensus to work correctly.

**The Result**: We can now design algorithms that satisfy these properties.

### The Journey Forward

Now that we understand the formal consensus problem, we can explore how to solve it. The next section will show us how to build an algorithm that satisfies these properties.

The key insight is that the consensus problem can be precisely defined, and understanding these properties is essential for building correct consensus algorithms.
## The Paxos Algorithm: Building Consensus Step by Step

Now let's explore how to build the Paxos algorithm step by step. We'll start with a broad definition of consensus and refine it into something we can actually implement.

### Our Approach: Refinement by Design

**Start with a broad definition of consensus**:
- We should eventually choose a value
- We should only choose one value

**Refine/narrow definition to something we can implement**:
- At each step, Lamport must argue the refinement is valid, e.g., P2a => P2

**The Goal**: We should only choose one value.

### The First Challenge: Choosing a Value

**The Simple Case**: Use a single acceptor.

**The Values**:
- A = Put k1 v1
- K = PutAppend k2 v2  
- M = Get k3
- Q = Delete k1

**The Problem**: What if the acceptor fails?

**The Result**: M is chosen! But if the acceptor fails, we lose the decision.

**The Solution**: Choose only when a "large enough" set of acceptors accepts.

**The Key Insight**: Using a majority set guarantees that at most one value is chosen.

### The Second Challenge: Accepting a Value

**The Simple Case**: Suppose only one value is proposed by a single proposer.

**The Requirement**: That value should be chosen!

**First requirement**: P1: An acceptor must accept the first proposal that it receives.

**The Problem**: What if we have multiple proposers, each proposing a different value?

**The Scenario**: 
- Proposer 1 proposes A
- Proposer 2 proposes Q  
- Proposer 3 proposes M
- Proposer 4 proposes K

**The Result**: No value is chosen! Each acceptor accepts the first proposal it receives, but they might receive different proposals.

### The Solution: Handling Multiple Proposals

**The Key Insight**: Acceptors must (be able to) accept more than one proposal.

**The Approach**: To keep track of different proposals, assign a natural number to each proposal.

**The Result**: A proposal is then a pair (psn, value).

**The Requirements**:
- Different proposals have different psn
- A proposal is chosen when it has been accepted by a majority of acceptors
- A value is chosen when a single proposal with that value has been chosen

### Assigning Proposal Numbers

**The Challenge**: Proposal numbers must be unique and infinite.

**The Problem**: A proposal number server won't work (it could fail).

**The Solution**: Assign each proposer an infinite slice.

**The Algorithm**: Proposer i of N gets: i, i+N, i+2N, i+3N, …

**The Result**: Each proposer has an infinite sequence of unique proposal numbers.

**Example with 4 proposers**:
- Proposer 0: 0, 4, 8, 12, 16, …
- Proposer 1: 1, 5, 9, 13, 17, …
- Proposer 2: 2, 6, 10, 14, 18, …
- Proposer 3: 3, 7, 11, 15, 19, …

### The Fundamental Insight

**The Key Realization**: We need a way to handle multiple proposals while ensuring that only one value is chosen.

**The Elegance**: Using proposal numbers allows us to track different proposals and ensure consistency.

**The Result**: We can now handle multiple proposers while maintaining the safety properties.

### The Journey Forward

Now that we understand how to handle multiple proposals, we need to ensure that all chosen proposals result in choosing the same value. The next section will show us how to achieve this.

The key insight is that proposal numbers provide a way to order proposals and ensure that the system can handle multiple proposers while maintaining consistency.
## Choosing a Unique Value: The Core Challenge

Now let's explore the core challenge of ensuring that all chosen proposals result in choosing the same value. This is the heart of the Paxos algorithm.

### The Core Requirement: P2

**The Challenge**: We need to guarantee that all chosen proposals result in choosing the same value.

**The Solution**: We introduce a second requirement (by induction on the proposal number):

**P2**: If a proposal with value v is chosen, then every higher-numbered proposal that is chosen has value v.

**What This Means**: Once a value is chosen, all future chosen values must be the same.

**The Power**: This ensures that the system cannot change its mind about what value was chosen.

### Refining P2: From P2 to P2a

**P2a**: If a proposal with value v is chosen, then every higher-numbered proposal accepted by any acceptor has value v.

**What This Means**: Once a value is chosen, acceptors cannot accept proposals with different values.

**The Connection**: P2a is a stronger requirement that implies P2.

**The Question**: What about P1?

### The Tension Between P1 and P2a

**The Scenario**: 
- Acceptor A has accepted (1,M)
- Proposer wants to propose (2,Q)
- How does A know it should not accept (2,Q)?

**The Problem**: P1 says accept the first proposal, but P2a says don't accept proposals with different values.

**The Resolution**: Do we still need P1?

**The Answer**: YES, to ensure that some proposal is accepted.

**The Challenge**: How well do P1 and P2a play together?

**The Problem**: Asynchrony is a problem...

**The Scenario**: 
- (1,M) is chosen
- Proposer wants to propose (2,K)
- Acceptor A doesn't know that (1,M) was chosen

**The Result**: M is chosen! But the system might not make progress.

### Strengthening P2: From P2a to P2b

**Recall P2a**: If a proposal with value v is chosen, then every higher-numbered proposal accepted by any acceptor has value v.

**The Problem**: P2a is hard to implement because acceptors don't always know what has been chosen.

**The Solution**: We strengthen it to:

**P2b**: If a proposal with value v is chosen, then every higher-numbered proposal issued by any proposer has value v.

**What This Means**: Proposers must ensure that their proposals have the right value.

**The Power**: This shifts the responsibility from acceptors to proposers.

### Implementing P2b: The Three Steps

**Step 1**: Suppose a proposer p wants to issue a proposal numbered n. What value should p propose?

**The Key Insight**: If (n',v) with n' < n is chosen, then in every majority set S of acceptors at least one acceptor has accepted (n',v).

**The Result**: If there is a majority set S where no acceptor has accepted (or will accept) a proposal with number less than n, then p can propose any value.

**Step 2**: What if for all S some acceptor ends up accepting a pair (n',v) with n' < n?

**The Claim**: p should propose the value of the highest numbered proposal among all accepted proposals numbered less than n.

**The Proof**: By induction on the number of proposals issued after a proposal is chosen.

**Step 3**: Achieved by enforcing the following invariant:

**P2c**: For any v and n, if a proposal with value v and number n is issued, then there is a set S consisting of a majority of acceptors such that either:
- No acceptor in S has accepted any proposal numbered less than n, or
- v is the value of the highest-numbered proposal among all proposals numbered less than n accepted by the acceptors in S

### The Fundamental Insight

**The Key Realization**: P2c provides a way to implement P2b, which ensures that all chosen proposals have the same value.

**The Elegance**: The invariant captures exactly what proposers need to do to maintain consistency.

**The Result**: We now have a concrete way to implement the core requirement of consensus.

### The Journey Forward

Now that we understand the core requirements, we can explore how to implement them in practice. The next section will show us how to build the actual Paxos algorithm.

The key insight is that P2c provides a concrete invariant that ensures consensus properties while being implementable in practice.
## Implementing P2c: The Two-Phase Protocol

Now let's explore how to implement P2c in practice. This leads us to the famous two-phase protocol that makes Paxos work.

### Understanding P2c: A Concrete Example

**The Scenario**: We have three acceptors with the following state:
- Acceptor 1: (1,A)
- Acceptor 2: (2,K)
- Acceptor 3: ?

**The Question**: What do we know about the third acceptor?

**The Analysis**:
- Could it have accepted (1,A)? No.
- Could it have accepted (2,K)? Yes.

**The Key Insight**: Proposal with highest number is the only proposal that could have been chosen!

### The Challenge: How Many Nodes to Consult?

**The Scenario**: (1,A) (2,K) nil

**The Question**: How many nodes do we need to consult?

**Option 1**: Consult all 3?
- **Result**: We know nothing was chosen!
- **Problem**: Want to be non-blocking if a majority are up

**Option 2**: Consult different pairs?
- **Consult 1 and 2?** Safe to propose (4,K)
- **Consult 1 and 3?** Safe to propose (4,A)
- **Consult 2 and 3?** Safe to propose (4,K)

**The Key Insight**: We only need to consult a majority, not all acceptors.

### P2c in Action: Three Scenarios

**Scenario 1**: No acceptor in S has accepted any proposal numbered less than n
- **State**: (4,K), (2,A), (1,A), nil
- **Result**: Safe to propose any value

**Scenario 2**: v is the value of the highest-numbered proposal among all proposals numbered less than n and accepted by the acceptors in S
- **State**: (4,K), (18,Q), (3,Q), (5,Q)
- **Result**: Must propose Q (the highest-numbered proposal)

**Scenario 3**: Mixed state with some acceptors having no proposals
- **State**: (18,Q), (2,K), nil, (4,Q)
- **Result**: Must propose Q (the highest-numbered proposal)

**Scenario 4**: Multiple acceptors with the same highest proposal
- **State**: (18,Q), (2,K), (5,K), (5,K), nil, (4,Q)
- **Result**: Must propose K (the highest-numbered proposal)

### The Future Telling Problem

**The Challenge**: To maintain P2c, a proposer that wishes to propose a proposal numbered n must learn the highest-numbered proposal with number less than n, if any, that has been or will be accepted by each acceptor in some majority of acceptors.

**The Problem**: We need to predict the future! We need to know what acceptors will accept, not just what they have accepted.

**The Solution**: Avoid predicting the future by extracting a promise from a majority of acceptors not to subsequently accept any proposals numbered less than n.

**The Power**: This promise allows us to implement P2c without needing to predict the future.

### The Two-Phase Protocol

**Phase 1**: The proposer asks acceptors to promise not to accept proposals with lower numbers.

**Phase 2**: The proposer asks acceptors to accept its proposal.

**The Result**: This two-phase approach ensures that P2c is satisfied.

### The Fundamental Insight

**The Key Realization**: The two-phase protocol provides a way to implement P2c without needing to predict the future.

**The Elegance**: The protocol ensures that proposers can safely propose values while maintaining consistency.

**The Result**: We now have a concrete algorithm that satisfies the consensus properties.

### The Journey Forward

Now that we understand how to implement P2c, we can explore the complete Paxos algorithm. The next section will show us the detailed protocol for both proposers and acceptors.

The key insight is that the two-phase protocol provides a practical way to implement the theoretical requirements of consensus.
## The Complete Paxos Protocol

Now let's explore the complete Paxos protocol, including the detailed algorithms for proposers and acceptors.

### The Proposer's Protocol: Phase 1

**The Process**: A proposer chooses a new proposal number n and sends a request to each member of some (majority) set of acceptors, asking it to respond with:

**a. A promise never again to accept a proposal numbered less than n, and**

**b. The accepted proposal with highest number less than n if any.**

**The Name**: Call this a prepare request with number n.

**What This Means**: The proposer is asking acceptors to promise not to accept older proposals and to tell it about any proposals they've already accepted.

### The Proposer's Protocol: Phase 2

**The Condition**: If the proposer receives a response from a majority of acceptors, then it can issue a proposal with number n and value v, where v is:

**a. The value of the highest-numbered proposal among the responses, or**

**b. Any value selected by the proposer if responders returned no proposals**

**The Process**: A proposer issues a proposal by sending, to some set of acceptors, a request that the proposal be accepted.

**The Name**: Call this an accept request.

**What This Means**: The proposer is asking acceptors to accept its proposal.

### The Acceptor's Protocol

**The Input**: An acceptor receives prepare and accept requests from proposers.

**The Flexibility**: It can ignore these without affecting safety.

**The Rules**:
- It can always respond to a prepare request
- It can respond to an accept request, accepting the proposal, iff it has not promised not to

**P1a**: An acceptor can accept a proposal numbered n iff it has not responded to a prepare request having number greater than n.

**The Connection**: This subsumes P1.

### Small Optimizations

**Optimization 1**: If an acceptor receives a prepare request r numbered n when it has already responded to a prepare request for n' > n, then the acceptor can simply ignore r.

**Optimization 2**: An acceptor can also ignore prepare requests for proposals it has already accepted.

**The Result**: An acceptor needs only remember:
- The highest numbered proposal it has accepted
- The number of the highest-numbered prepare request to which it has responded

**The Requirement**: This information needs to be stored on stable storage to allow restarts.
Choosing a value:
Phase 1
A proposer chooses a new n and sends <prepare,n>
to a majority of acceptors
If an acceptor a receives <prepare,n’>, where n’ > n
of any <prepare,n> to which it has responded, then it
responds to <prepare, n’ > with
a promise not to accept any more proposals
numbered less than n’
the highest numbered proposal (if any) that it has
accepted
Choosing a value:
Phase 2
If the proposer receives a response to <prepare,n>
from a majority of acceptors, then it sends to each
<accept,n,v>, where v is either
the value of the highest numbered proposal
among the responses
any value if the responses reported no proposals
If an acceptor receives <accept,n,v>, it accepts the
proposal unless it has in the meantime responded to
<prepare,n’> , where n’ > n
### Learning Chosen Values

**The Challenge**: Once a value is chosen, learners should find out about it.

**Strategy 1**: Each acceptor informs each learner whenever it accepts a proposal.

**Strategy 2**: Acceptors inform a distinguished learner, who informs the other learners.

**Strategy 3**: Something in between (a set of not-quite-as-distinguished learners).

**The Problem**: Because of failures (message loss and acceptor crashes) a learner may not learn that a value has been chosen.

**The Scenario**: 
- (4,K) was chosen
- Was M chosen? ☠
- (7,M) - Propose something!

### Liveness: The Progress Problem

**The Challenge**: Progress is not guaranteed.

**The Scenario**: n1 < n2 < n3 < n4 < …

**The Problem**: 
- p1 proposes n1, gets accepted
- p2 proposes n2, gets accepted  
- p1 proposes n3, gets accepted
- p2 proposes n4, gets accepted

**The Result**: The system can get stuck in a cycle where no value is ever chosen.

**The Solution**: We need a leader to coordinate proposals.
## Implementing State Machine Replication

Now let's explore how to use Paxos to implement state machine replication, which is the practical application of consensus.

### The Basic Approach

**The Idea**: Implement a sequence of separate instances of consensus, where the value chosen by the ith instance is the ith message in the sequence.

**The Roles**: Each server assumes all three roles in each instance of the algorithm.

**The Assumption**: Assume that the set of servers is fixed.

**The Result**: We can build a replicated state machine that all servers execute in the same order.

### The Role of the Leader

**In Normal Operation**: Elect a single server to be a leader. The leader acts as the distinguished proposer in all instances of the consensus algorithm.

**The Process**: Clients send commands to the leader, which decides where in the sequence each command should appear.

**The Example**: If the leader, for example, decides that a client command is the kth command, it tries to have the command chosen as the value in the kth instance of consensus.

**The Power**: This ensures that all servers execute commands in the same order.

### Paxos and FLP

**The FLP Theorem**: In an asynchronous system, it's impossible to solve consensus deterministically in the presence of even one failure.

**Paxos's Response**: Paxos is always safe–despite asynchrony.

**The Liveness**: Once a leader is elected, Paxos is live.

**The Catch**: "Ciao ciao" FLP? To be live, Paxos requires a single leader.

**The Paradox**: "Leader election" is impossible in an asynchronous system (gotcha!)

**The Resolution**: Given FLP, Paxos is the next best thing: always safe, and live during periods of synchrony.

### A New Leader is Elected

**The Scenario**: λ is elected as the new leader.

**The Challenge**: Since λ is a learner in all instances of consensus, it should know most of the commands that have already been chosen. For example, it might know commands 1-10, 13, and 15.

**The Process**: It executes phase 1 of instances 11, 12, and 14 and of all instances 16 and larger.

**The Result**: This might leave, say, 14 and 16 constrained and 11, 12 and all commands after 16 unconstrained.

**The Resolution**: λ then executes phase 2 of 14 and 16, thereby choosing the commands numbered 14 and 16.

### Stop-gap Measures

**The Problem**: All replicas can execute commands 1-10, but not 13-16 because 11 and 12 haven't yet been chosen.

**The Solution**: λ can either take the next two commands requested by clients to be commands 11 and 12, or can propose immediately that 11 and 12 be no-op commands.

**The Process**: λ runs phase 2 of consensus for instance numbers 11 and 12.

**The Result**: Once consensus is achieved, all replicas can execute all commands through 16.

### To Infinity, and Beyond

**The Efficiency**: λ can efficiently execute phase 1 for infinitely many instances of consensus! (e.g. command 16 and higher)

**The Method**: λ just sends a message with a sufficiently high proposal number for all instances.

**The Optimization**: An acceptor replies non trivially only for instances for which it has already accepted a value.

**The Power**: This allows the system to handle an infinite sequence of commands efficiently.

### The Fundamental Insight

**The Key Realization**: Paxos provides a complete solution to consensus that can be used to build practical distributed systems.

**The Elegance**: The algorithm ensures safety even in the presence of failures and provides liveness when conditions are right.

**The Result**: We can build reliable distributed systems using Paxos as the foundation.

### The Journey Complete: Understanding Paxos

**What We've Learned**:
1. **The Part-Time Parliament**: A beautiful analogy that explains the core concepts
2. **The Consensus Problem**: Formal definition of what we're trying to solve
3. **The Paxos Algorithm**: How to achieve consensus in practice
4. **Phase 1 and Phase 2**: The two-phase protocol that makes Paxos work
5. **State Machine Replication**: How to use Paxos in practice
6. **Leader Election**: How to make Paxos practical in real systems

**The Fundamental Insight**: Paxos provides a complete solution to the consensus problem that is both safe and live under the right conditions.

**The Impact**: Understanding Paxos is essential for building reliable distributed systems.

**The Legacy**: Paxos continues to influence how we build distributed systems today.

### The End of the Journey

Paxos represents one of the most elegant solutions to one of the most fundamental problems in distributed systems. The algorithm is both simple and powerful, providing a way to achieve consensus even in the face of failures and network issues.

The journey from the simple idea of state machine replication to the sophisticated Paxos algorithm shows how complex problems can be solved with careful thought and elegant design. The challenge is always the same: how do you get multiple machines to agree on something?

Paxos provides an answer to this question, and it continues to influence how we build distributed systems today.