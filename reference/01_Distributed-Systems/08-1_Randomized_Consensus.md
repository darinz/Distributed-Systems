# Randomized Consensus: Breaking the FLP Barrier

## The Fundamental Problem: FLP Impossibility

In distributed systems, we face a fundamental limitation that seems insurmountable: the FLP impossibility result. This theorem tells us that in certain environments, consensus is impossible to solve deterministically. But what if we could break this barrier? What if we could achieve consensus even when the FLP theorem says it's impossible?

This document explores randomized consensus algorithms, which use randomness to solve the consensus problem in situations where deterministic algorithms fail. Understanding these algorithms is crucial for building distributed systems that need to work even in the most challenging environments.

### The FLP Impossibility Theorem

**The Theorem**: In an asynchronous environment in which a single process can fail by crashing, there does not exist a protocol which solves binary consensus.

**What This Means**: 
- **Asynchronous Environment**: No bounds on message delays or process speeds
- **Single Process Failure**: Even one process crashing makes consensus impossible
- **Binary Consensus**: Agreement on a single bit (0 or 1)
- **No Protocol Exists**: It's mathematically impossible to solve this problem deterministically

**The Impact**: This is one of the most fundamental results in distributed systems theory.

### Why Paxos Doesn't Save Us

**The Problem**: Paxos doesn't guarantee liveness.

**What This Means**: While Paxos is always safe (never makes incorrect decisions), it can get stuck and never make progress.

**The FLP Connection**: The FLP theorem shows that this isn't a flaw in Paxos—it's a fundamental limitation of any deterministic algorithm.

**The Result**: Paxos assumed a deterministic computation model, which is why it can't guarantee liveness.

### The Solution: Let's Go Random!

**The Key Insight**: Randomization can break the FLP impossibility barrier.

**How It Works**: By introducing randomness into the algorithm, we can achieve consensus with high probability.

**The Trade-off**: We give up deterministic guarantees in exchange for the possibility of consensus.

**The Result**: We can solve consensus in environments where deterministic algorithms fail.

### Ben-Or's Algorithm: A Breakthrough

**The Algorithm**: Ben-Or's algorithm uses randomization to guarantee consensus for crash failures when f < n/2.

**What This Means**: 
- **Randomization**: The algorithm makes random choices
- **Crash Failures**: Processes can fail by stopping, but not by behaving maliciously
- **f < n/2**: Less than half of the processes can fail
- **Guarantee**: The algorithm will eventually reach consensus

**The Power**: A variant even works for Byzantine faults (malicious behavior)!

### The Intuition: How Randomization Helps

**The Basic Idea**: 
- At first every process proposes their input value
- After that, they propose random values
- When enough processes propose the same value, the value is chosen
- Eventually, that will happen!

**The Key Insight**: Randomness breaks symmetry and allows the system to escape from deadlock situations.

**The Process**:
1. **Round 1**: Processes propose their input values (1, 0, 1, 0)
2. **Round 2**: Processes propose random values (0, 1, 0, 1)  
3. **Round 3**: Processes propose random values (1, 0, 1, 1)
4. **Round 4**: Processes propose random values (0, 0, 0, 0)
5. **Success**: All processes agree on 0!

**Why This Works**: Eventually, randomness will cause enough processes to propose the same value, breaking the deadlock.

### The Fundamental Insight

**The Key Realization**: Randomization provides a way to solve consensus even when deterministic algorithms fail.

**The Elegance**: The algorithm is simple but powerful, using randomness to break fundamental impossibility results.

**The Impact**: Understanding randomized consensus is essential for building distributed systems that work in challenging environments.

### The Journey Ahead

This document will take you through the complete story of randomized consensus:

1. **FLP Impossibility**: Why deterministic consensus is impossible
2. **The Randomization Solution**: How randomness breaks the barrier
3. **Ben-Or's Algorithm**: A concrete randomized consensus algorithm
4. **The Setup**: How to structure the algorithm
5. **Correctness Proofs**: Why the algorithm works
6. **Practical Implications**: How to use randomized consensus

By the end, you'll understand not just how randomized consensus works, but why it's necessary and how it overcomes fundamental impossibility results.

### The End of the Beginning

Randomized consensus represents a breakthrough in distributed systems theory. By embracing randomness, we can solve problems that deterministic algorithms cannot solve.

The journey from FLP impossibility to randomized consensus shows how creative thinking can overcome seemingly insurmountable barriers. The challenge is always the same: how do you achieve consensus when deterministic algorithms fail?

Randomization provides an answer to this question, and it continues to influence how we build distributed systems today.
## The Setup: Structuring the Algorithm

Now let's explore how to structure a randomized consensus algorithm. The key insight is that we need to organize the algorithm in a way that allows randomness to break deadlocks while maintaining correctness.

### The Basic Setup

**Binary Consensus**: Again, we're considering binary consensus (agreement on 0 or 1).

**Why Binary**: Binary consensus is conceptually simple and forms the foundation for more complex consensus problems.

**The Challenge**: Even this simple problem is impossible to solve deterministically in asynchronous systems.

### The Round Structure

**Protocol Structure**: Protocol proceeds in asynchronous rounds, where each round has two phases.

**What This Means**: 
- **Asynchronous Rounds**: No timing guarantees between rounds
- **Two Phases**: Each round is divided into two distinct phases
- **Phase 1**: Preliminary proposals and coordination
- **Phase 2**: Final proposals and decision making

**The Power**: This structure allows processes to coordinate and then make decisions based on what they learned.

### Message Handling

**Broadcasting**: For each phase, processes broadcast their input values and wait for n – f messages from the other processes.

**What This Means**:
- **n – f Messages**: Wait for responses from a majority of processes
- **Fault Tolerance**: Can handle up to f process failures
- **Majority Quorum**: Ensures that any two majorities will intersect

**The Result**: This provides fault tolerance while maintaining progress.

### Message Tagging and Reliability

**Message Tags**: Each message is tagged with the round and phase number.

**Network Handling**: Messages can be resent to deal with a lossy network.

**Value Locking**: But once a message is sent, that value is locked in for that process for that phase/round.

**What This Means**: 
- **Reliability**: Messages can be retransmitted if lost
- **Consistency**: Once sent, a process cannot change its value for that phase/round
- **Determinism**: The algorithm behavior is deterministic given the message delivery

## The Ben-Or Algorithm: A Concrete Solution

Now let's explore Ben-Or's algorithm, which implements the randomized consensus approach we've been discussing.

### The Basic Structure

**Process Behavior**: Processes send proposals for each phase and then block and wait for the requisite n – f messages (including their own).

**What This Means**: 
- **Proposal Phase**: Each process sends its proposal
- **Waiting Phase**: Processes wait for enough responses
- **Including Their Own**: Each process counts its own message in the quorum

**The Result**: This ensures that processes coordinate their actions.

### Phase 1: Preliminary Proposals

**The Process**: During the first phase, processes make a preliminary proposal.

**What This Means**: 
- **Initial Values**: Processes propose their current best guess
- **Coordination**: Processes learn what other processes are thinking
- **No Decisions**: No final decisions are made in this phase

**The Purpose**: This phase allows processes to coordinate and see if they can reach agreement.

### Phase 1 Decision Logic

**The Key Insight**: If they receive matching responses from a majority in the first phase, they propose that value in the second phase. Otherwise, they propose ⊥ (a special null value).

**What This Means**:
- **Majority Agreement**: If a majority agrees on a value, propose that value
- **No Agreement**: If no majority agrees, propose ⊥ (null)
- **The ⊥ Value**: A special value that indicates "no preference"

**The Result**: This allows processes to coordinate when possible and use randomness when necessary.

### Phase 2: Final Proposals and Decisions

**The Process**: If they get enough non-⊥ responses from the second phase, they decide.

**What This Means**:
- **Non-⊥ Responses**: Look for actual values, not null values
- **Enough Responses**: Need sufficient responses to make a decision
- **Decision Making**: Can finally decide on a value

**The Result**: This phase allows processes to reach consensus when conditions are right.

### The Fundamental Insight

**The Key Realization**: The two-phase structure allows processes to coordinate and then make decisions based on what they learned.

**The Elegance**: The algorithm is simple but powerful, using coordination when possible and randomness when necessary.

**The Result**: We now have a concrete algorithm that can achieve consensus even in challenging environments.
𝑎←input
loop:
send_phase1(𝑎)
𝐴←receive_phase1()
if (∃𝑎ʹ ∈ 𝐴 : |𝐴𝑎ʹ| > 𝑛/2):
𝑏←𝑎ʹ
else:
𝑏←⊥
send_phase2(𝑏)
𝐵←receive_phase2()
if (∃𝑏ʹ ∈ 𝐵 : 𝑏ʹ≠⊥ ∧ |𝐵𝑏ʹ| > 𝑓):
decide(𝑏ʹ)
if (∃𝑏ʹ ∈ 𝐵 : 𝑏ʹ≠⊥):
𝑎←𝑏ʹ
else:
𝑎←choose_random({0,1})
Do We Have Consensus?
• Agreement: No two
processes decide
diﬀerent values.
𝑎←input
loop:
send_phase1(𝑎)
𝐴←receive_phase1()
if (∃𝑎ʹ ∈ 𝐴 : |𝐴𝑎ʹ| > 𝑛/2):
𝑏←𝑎ʹ
• Integrity: Every process
decides at most one
value, and if a process
decides a value, some
process had it as its input.
else:
𝑏←⊥
send_phase2(𝑏)
𝐵←receive_phase2()
if (∃𝑏ʹ ∈ 𝐵 : 𝑏ʹ≠⊥ ∧ |𝐵𝑏ʹ| > 𝑓):
decide(𝑏ʹ)
if (∃𝑏ʹ ∈ 𝐵 : 𝑏ʹ≠⊥):
𝑎←𝑏ʹ
• Termination: Every
correct process eventually
decides a value.
else:
𝑎←choose_random({0,1})
## Correctness Proofs: Why the Algorithm Works

Now let's prove that Ben-Or's algorithm actually satisfies the consensus properties. We'll go through each property systematically to understand why the algorithm is correct.

### Integrity: The Foundation of Correctness

**The Requirement**: Every process decides at most one value, and if a process decides a value, some process had it as its input.

**Case 1: Mixed Input Values**
If both 0 and 1 are input values to processes, integrity is trivially satisfied.

**What This Means**: Since both values appear as inputs, any decision satisfies the integrity requirement.

**Case 2: All Same Input Value**
Suppose all processes have the same input value.

**The Process**:
- Then, they all send the same phase 1 value in round 1
- So they all send that same value in phase 2
- So they all decide that value at the end of round 1

**The Result**: Consensus is reached immediately in the first round.

**Why This Works**: When all processes start with the same value, no coordination or randomness is needed.

### The Fundamental Insight

**The Key Realization**: Integrity is satisfied because the algorithm only decides on values that were actually proposed.

**The Elegance**: The algorithm structure ensures that decisions are always based on real input values.

**The Result**: We can be confident that the algorithm never makes up values.
𝑎←input
loop:
send_phase1(𝑎)
𝐴←receive_phase1()
if (∃𝑎ʹ ∈ 𝐴 : |𝐴𝑎ʹ| > 𝑛/2):
𝑏←𝑎ʹ
else:
𝑏←⊥
send_phase2(𝑏)
𝐵←receive_phase2()
if (∃𝑏ʹ ∈ 𝐵 : 𝑏ʹ≠⊥ ∧ |𝐵𝑏ʹ| > 𝑓):
decide(𝑏ʹ)
if (∃𝑏ʹ ∈ 𝐵 : 𝑏ʹ≠⊥):
𝑎←𝑏ʹ
else:
𝑎←choose_random({0,1})
### Agreement: The Core Consensus Property

**The Requirement**: No two processes decide different values.

**The Key Lemma**: No two processes receive different non-⊥ phase 2 values in the same round.

**The Proof**: Suppose they did. That means that one process received 0s from a majority in phase 1 and another received 1s.

**The Contradiction**: But majorities intersect! This means that at least one process must have sent both 0 and 1 in phase 1, which is impossible.

**The Result**: In any given round, all processes must see the same non-⊥ values in phase 2.

**Why This Matters**: This ensures that if any process decides in a round, all processes that decide in that round decide on the same value.

### The Fundamental Insight

**The Key Realization**: The majority intersection property ensures that agreement is maintained.

**The Elegance**: The algorithm structure prevents processes from seeing conflicting information.

**The Result**: We can be confident that the algorithm never allows different processes to decide on different values.
𝑎←input
loop:
send_phase1(𝑎)
𝐴←receive_phase1()
if (∃𝑎ʹ ∈ 𝐴 : |𝐴𝑎ʹ| > 𝑛/2):
𝑏←𝑎ʹ
else:
𝑏←⊥
send_phase2(𝑏)
𝐵←receive_phase2()
if (∃𝑏ʹ ∈ 𝐵 : 𝑏ʹ≠⊥ ∧ |𝐵𝑏ʹ| > 𝑓):
decide(𝑏ʹ)
if (∃𝑏ʹ ∈ 𝐵 : 𝑏ʹ≠⊥):
𝑎←𝑏ʹ
else:
𝑎←choose_random({0,1})
### Agreement + Integrity: The Complete Picture

**The Scenario**: Let round r be the first round any process decides a value, 0 w.l.o.g. (without loss of generality).

**The Analysis**: If a process decided a value, it must have received > f 0s in phase 2.

**The Key Insight**: Which means that every process received at least one 0 because they all wait for n – f messages. No process received a 1 by the previous lemma.

**The Result**: Therefore, on round r + 1 (and all subsequent rounds), all processes propose 0 and all processes decide 0.

**Why This Works**: Once a decision is made, all processes adopt that value, ensuring future consensus.

### The Fundamental Insight

**The Key Realization**: The algorithm ensures that once a decision is made, all future decisions will be the same.

**The Elegance**: This property maintains agreement across rounds.

**The Result**: We can be confident that the algorithm maintains consistency over time.
𝑎←input
loop:
send_phase1(𝑎)
𝐴←receive_phase1()
if (∃𝑎ʹ ∈ 𝐴 : |𝐴𝑎ʹ| > 𝑛/2):
𝑏←𝑎ʹ
else:
𝑏←⊥
send_phase2(𝑏)
𝐵←receive_phase2()
if (∃𝑏ʹ ∈ 𝐵 : 𝑏ʹ≠⊥ ∧ |𝐵𝑏ʹ| > 𝑓):
decide(𝑏ʹ)
if (∃𝑏ʹ ∈ 𝐵 : 𝑏ʹ≠⊥):
𝑎←𝑏ʹ
else:
𝑎←choose_random({0,1})
### Termination: The Power of Randomness

**The Key Insight**: We know that if all processes propose the same value for a round, they all decide that value that round.

**The Probability**: At worst, the probability of this happening on any particular round is 1/2^n.

**Why This Works**: By the previous lemma, all the non-random values are identical.

**The Result**: Over time, the probability of this happening on at least one round converges to 1.

**The Intuition**: Randomness ensures that eventually all processes will propose the same value.

### The Fundamental Insight

**The Key Realization**: Randomness provides a way to break deadlocks and ensure progress.

**The Elegance**: The algorithm guarantees termination with probability 1.

**The Result**: We can be confident that the algorithm will eventually reach consensus.
𝑎←input
loop:
send_phase1(𝑎)
𝐴←receive_phase1()
if (∃𝑎ʹ ∈ 𝐴 : |𝐴𝑎ʹ| > 𝑛/2):
𝑏←𝑎ʹ
else:
𝑏←⊥
send_phase2(𝑏)
𝐵←receive_phase2()
if (∃𝑏ʹ ∈ 𝐵 : 𝑏ʹ≠⊥ ∧ |𝐵𝑏ʹ| > 𝑓):
decide(𝑏ʹ)
if (∃𝑏ʹ ∈ 𝐵 : 𝑏ʹ≠⊥):
𝑎←𝑏ʹ
else:
𝑎←choose_random({0,1})
## Beyond Binary Consensus: Extending the Algorithm

Now let's explore how to extend Ben-Or's algorithm beyond binary consensus to handle more complex scenarios.

### The Limitation of Binary Consensus

**Binary Consensus**: Binary consensus is conceptually simple but not as useful.

**The Challenge**: Real-world systems often need to agree on more complex values than just 0 or 1.

**The Question**: Can we extend the algorithm to support larger domains?

### The Extension: Supporting Larger Domains

**The Power**: However, the algorithm can be extended to support larger domains, even when the processes don't know the domains a priori and even when some processes don't receive input values.

**The Approach**:
- **Processes without input values start by proposing ⊥**
- **Instead of randomly choosing from {0,1}, processes randomly choose from all non-⊥ values they've seen so far (in any message)**
- **Only choose ⊥ as a last resort**

**What This Means**: The algorithm can handle arbitrary value domains without prior knowledge.

### The Fundamental Insight

**The Key Realization**: The algorithm structure is general enough to handle complex scenarios.

**The Elegance**: The same principles work for binary and multi-valued consensus.

**The Result**: We can build powerful consensus systems using these techniques.
𝑎←input
loop:
send_phase1(𝑎)
𝐴←receive_phase1()
if (∃𝑎ʹ ∈ 𝐴 : |𝐴𝑎ʹ| > 𝑛/2):
𝑏←𝑎ʹ
else:
𝑏←⊥
send_phase2(𝑏)
𝐵←receive_phase2()
if (∃𝑏ʹ ∈ 𝐵 : 𝑏ʹ≠⊥ ∧ |𝐵𝑏ʹ| > 𝑓):
decide(𝑏ʹ)
if (∃𝑏ʹ ∈ 𝐵 : 𝑏ʹ≠⊥):
𝑎←𝑏ʹ
else:
𝑎←choose_random({0,1})
## Key Takeaways: The Power of Randomized Consensus

Now let's summarize the key insights from our exploration of randomized consensus algorithms.

### Takeaway 1: Randomization Can Actually Solve Consensus

**The Key Insight**: Randomization can actually solve consensus*.

**What This Means**: By introducing randomness, we can overcome the FLP impossibility result.

**The Power**: We can achieve consensus in environments where deterministic algorithms fail.

**The Trade-off**: We give up deterministic guarantees in exchange for the possibility of consensus.

### Takeaway 2: Structuring Asynchronous Protocols

**The Key Insight**: You can structure an asynchronous protocol using rounds.

**What This Means**: 
- **Rounds**: Organize computation into discrete rounds
- **Phases**: Each round can have multiple phases
- **Coordination**: Use rounds to coordinate processes

**The Power**: This approach is potentially useful and certainly an interesting way to think about asynchronous computation.

### The Fundamental Insights

**The Key Realization**: Randomized consensus provides a way to solve problems that deterministic algorithms cannot solve.

**The Elegance**: The algorithms are simple but powerful, using randomness to break fundamental impossibility results.

**The Impact**: Understanding randomized consensus is essential for building distributed systems that work in challenging environments.

## The Journey Complete: Understanding Randomized Consensus

**What We've Learned**:
1. **FLP Impossibility**: Why deterministic consensus is impossible in asynchronous systems
2. **The Randomization Solution**: How randomness breaks the impossibility barrier
3. **Ben-Or's Algorithm**: A concrete randomized consensus algorithm
4. **The Setup**: How to structure the algorithm with rounds and phases
5. **Correctness Proofs**: Why the algorithm satisfies consensus properties
6. **Beyond Binary**: How to extend the algorithm to handle complex scenarios

**The Fundamental Insight**: Randomized consensus provides a way to solve consensus even when deterministic algorithms fail.

**The Impact**: Understanding randomized consensus is essential for building reliable distributed systems.

**The Legacy**: These algorithms continue to influence how we build distributed systems today.

### The End of the Journey

Randomized consensus represents a breakthrough in distributed systems theory. By embracing randomness, we can solve problems that deterministic algorithms cannot solve.

The journey from FLP impossibility to randomized consensus shows how creative thinking can overcome seemingly insurmountable barriers. The challenge is always the same: how do you achieve consensus when deterministic algorithms fail?

Randomization provides an answer to this question, and it continues to influence how we build distributed systems today.