# Problem Set 3: Paxos Protocol Analysis

## Learning Objectives
This problem set focuses on deep understanding of the Paxos consensus protocol:
- Analyzing Paxos acceptor states and their validity
- Understanding Paxos safety guarantees and invariants
- Constructing execution scenarios that lead to specific states
- Identifying liveness issues and their causes
- Evaluating alternative protocol implementations for correctness

---

## Question 1: Paxos Acceptor State Validation

**Context**: Understanding what states are possible in Paxos is crucial for reasoning about protocol correctness and debugging distributed systems. This question explores the fundamental invariants that Paxos maintains and how they constrain the possible states of acceptor nodes.

Consider a deployment of single-instance Paxos with three acceptors (A1, A2, A3). For each of the following states, determine whether it is valid and explain your reasoning.

**Key Definitions:**
- **Valid state**: A state is valid if there exists some sequence of message deliveries, message drops, and node failures that leads to this state, assuming correct implementation of proposers and acceptors
- **Acceptor state**: Each acceptor's highest accepted proposal is either (n, v) where n is the proposal number and v is the value, or ⊥ indicating no proposal has been accepted

**Critical Paxos Invariants to Consider:**
1. **Safety Invariant**: If a proposal with value v is chosen, then every higher-numbered proposal that is chosen must also have value v
2. **Majority Requirement**: A value is chosen only when a majority of acceptors accept a proposal with that value
3. **Proposal Numbering**: Proposal numbers are unique and monotonically increasing
4. **Acceptor Behavior**: Acceptors can only accept proposals with numbers higher than any they have previously accepted

### 1a. Initial State

**State**: A1: ⊥, A2: ⊥, A3: ⊥

**Detailed Hints:**
- **Think about**: What does this state represent in the context of Paxos execution?
- **Consider**: Is this the starting state of the system, or could it occur after some execution?
- **Key insight**: This represents the initial state where no proposals have been accepted yet
- **Analysis approach**: This is trivially valid as it represents the system's starting state

**Step-by-step reasoning:**
1. Identify what this state represents (no accepted proposals)
2. Consider whether this could be the initial state of the system
3. Determine if any Paxos invariants are violated
4. Conclude about validity

### 1b. Different Values, Different Numbers

**State**: A1: (1, x), A2: (2, y), A3: ⊥

**Detailed Hints:**
- **Think about**: Can acceptors accept proposals with different values at different proposal numbers?
- **Consider**: What does the Paxos safety invariant say about this scenario?
- **Key insight**: The safety invariant requires that if a value is chosen, all higher-numbered chosen proposals must have the same value
- **Analysis approach**: Check if this state violates any Paxos invariants

**Step-by-step reasoning:**
1. Identify the accepted proposals and their numbers/values
2. Check if any value has been chosen (majority acceptance)
3. If a value is chosen, verify that higher-numbered proposals follow the safety invariant
4. Determine validity based on invariant satisfaction

### 1c. Same Number, Different Values

**State**: A1: (2, x), A2: (2, y), A3: ⊥

**Detailed Hints:**
- **Think about**: Can two acceptors accept proposals with the same number but different values?
- **Consider**: How does Paxos ensure that proposal numbers are unique?
- **Key insight**: Proposal numbers must be unique across all proposers
- **Analysis approach**: This state suggests a violation of proposal number uniqueness

**Step-by-step reasoning:**
1. Notice that both A1 and A2 have accepted proposal number 2
2. Consider how proposal numbers are assigned in Paxos
3. Determine if this violates the uniqueness requirement
4. Conclude about validity

### 1d. Three Different Values

**State**: A1: (1, x), A2: (2, y), A3: (3, z)

**Detailed Hints:**
- **Think about**: What does the Paxos safety invariant require about chosen values?
- **Consider**: If any value has been chosen, what must be true about higher-numbered proposals?
- **Key insight**: The safety invariant prevents different values from being chosen
- **Analysis approach**: Check if this state allows multiple values to be chosen

**Step-by-step reasoning:**
1. Identify all accepted proposals and their numbers/values
2. Determine if any value has been chosen (majority acceptance)
3. If multiple values could be chosen, check if this violates safety
4. Conclude about validity

### 1e. Same Value, Different Numbers

**State**: A1: (1, x), A2: (2, x), A3: (3, x)

**Detailed Hints:**
- **Think about**: Can acceptors accept proposals with the same value at different numbers?
- **Consider**: What does the Paxos safety invariant say about this scenario?
- **Key insight**: The safety invariant is satisfied when all chosen proposals have the same value
- **Analysis approach**: Verify that this state maintains Paxos safety guarantees

**Step-by-step reasoning:**
1. Identify all accepted proposals and their numbers/values
2. Check if a value has been chosen (majority acceptance)
3. Verify that the safety invariant is satisfied
4. Determine validity based on invariant satisfaction

**General Analysis Strategy:**
1. For each state, identify all accepted proposals
2. Check if any value has been chosen (accepted by majority)
3. If a value is chosen, verify that higher-numbered proposals satisfy the safety invariant
4. Consider whether the state could result from a valid Paxos execution
5. Remember that proposal numbers must be unique and monotonically increasing

## Question 2: Acceptor States in a Larger System

**Context**: Analyzing Paxos states in larger systems helps understand how the protocol scales and how complex execution scenarios can lead to specific acceptor configurations. This question explores the relationship between proposal numbers, values, and the safety guarantees in a five-acceptor system.

Consider a deployment with five acceptors (A1, A2, A3, A4, A5). Analyze the following state for validity and provide detailed reasoning.

**State to Analyze**: A1: (20, x), A2: ⊥, A3: (22, y), A4: (20, x), A5: (18, x)

**Detailed Hints:**
- **Think about**: What does this state tell us about the execution history?
- **Consider**: Which values have been chosen (accepted by majority)?
- **Key insight**: With 5 acceptors, a majority is 3 acceptors
- **Analysis approach**: Check if any value has been chosen and verify safety invariants

**Step-by-step Analysis:**

### Step 1: Identify Accepted Proposals
- A1: (20, x) - accepted proposal 20 with value x
- A2: ⊥ - no proposal accepted
- A3: (22, y) - accepted proposal 22 with value y  
- A4: (20, x) - accepted proposal 20 with value x
- A5: (18, x) - accepted proposal 18 with value x

### Step 2: Check for Chosen Values
- **Value x**: Accepted by A1, A4, A5 (3 acceptors = majority) ✓
- **Value y**: Accepted by A3 only (1 acceptor < majority) ✗

### Step 3: Verify Safety Invariant
- **Question**: If value x is chosen at proposal 18, what must be true about higher-numbered proposals?
- **Safety requirement**: All higher-numbered chosen proposals must have value x
- **Check**: Proposal 20 has value x ✓, but proposal 22 has value y ✗

### Step 4: Determine Validity
**Detailed Reasoning:**
- The state shows that value x was chosen (accepted by majority at proposal 18)
- However, proposal 22 with value y was also accepted
- This violates the Paxos safety invariant: if a value is chosen, all higher-numbered chosen proposals must have the same value
- Therefore, this state is **not valid**

**Alternative Analysis Approach:**
If you believe the state might be valid, consider:
1. **Execution scenario construction**: Try to construct a sequence of events that leads to this state
2. **Message ordering**: Consider what order of message deliveries would be required
3. **Proposal numbering**: Verify that proposal numbers are unique and increasing
4. **Safety violation**: Identify the specific safety invariant that would be violated

**Key Learning Points:**
- Paxos safety requires that once a value is chosen, all future chosen values must be the same
- The safety invariant applies even when proposals are accepted by different numbers of acceptors
- A single acceptor accepting a higher-numbered proposal with a different value can violate safety
- Understanding the relationship between proposal numbers and values is crucial for Paxos analysis

## Question 3: Analyzing a Dubious Execution Scenario

**Context**: This question explores how Paxos safety guarantees prevent certain problematic execution scenarios. Understanding why certain sequences of events cannot occur is crucial for appreciating Paxos's correctness properties and the mechanisms that ensure safety.

Consider a Paxos deployment with acceptors A1, A2, and A3, proposers P1, P2, and a distinguished learner L. According to Paxos, a value is chosen when a majority of acceptors accept a proposal with that value, and only a single value can be chosen.

**The Dubious Execution Sequence:**
1. P1 prepares proposal number 1, and gets responses from A1, A2, and A3
2. P1 sends (1, x) to A1 and A3 and gets responses from both. However, P1's proposal to A2 was dropped. Because a majority accepted, P1 informs L that x has been chosen. P1 then crashes
3. P2 prepares proposal number 2, and gets responses from A2 and A3
4. P2 sends (2, y) messages to A2 and A3 gets responses from both, so P2 informs L that y has been chosen

**Analysis Questions:**
1. How does Paxos ensure that this sequence cannot happen?
2. What actually happens instead?
3. Which value is ultimately chosen?

**Detailed Hints:**

### Understanding the Scenario
- **Think about**: What does this scenario claim happened?
- **Consider**: What are the key events that seem problematic?
- **Key insight**: The scenario suggests that both x and y could be chosen, violating Paxos safety

### Step-by-step Analysis

#### Step 1: Analyze the First Phase (P1's Proposal)
**Detailed Hints:**
- **Think about**: What does P1 learn from the prepare responses?
- **Consider**: What information do acceptors include in their prepare responses?
- **Key insight**: Prepare responses include the highest-numbered proposal each acceptor has accepted

**What P1 learns from prepare responses:**
- A1 responds: "I haven't accepted any proposal" (or reports highest accepted)
- A2 responds: "I haven't accepted any proposal" (or reports highest accepted)  
- A3 responds: "I haven't accepted any proposal" (or reports highest accepted)

#### Step 2: Analyze P1's Accept Phase
**Detailed Hints:**
- **Think about**: What value should P1 propose in the accept phase?
- **Consider**: What does the Paxos algorithm require when no previous proposal has been accepted?
- **Key insight**: P1 can propose its own value x since no previous proposal was accepted

**P1's accept phase:**
- P1 sends (1, x) to A1, A2, A3
- A1 and A3 accept (1, x)
- A2's message is dropped
- P1 concludes x is chosen (majority acceptance)

#### Step 3: Analyze P2's Prepare Phase
**Detailed Hints:**
- **Think about**: What does P2 learn from A2 and A3's prepare responses?
- **Consider**: What information do A2 and A3 have about previously accepted proposals?
- **Key insight**: A3 accepted (1, x), so it must report this to P2

**What P2 learns from prepare responses:**
- A2 responds: "I haven't accepted any proposal" (message to A2 was dropped)
- A3 responds: "I have accepted proposal (1, x)"

#### Step 4: Analyze P2's Accept Phase
**Detailed Hints:**
- **Think about**: What value must P2 propose in the accept phase?
- **Consider**: What does the Paxos algorithm require when a previous proposal has been accepted?
- **Key insight**: P2 must propose the value from the highest-numbered proposal it learned about

**P2's accept phase:**
- P2 learns that A3 accepted (1, x)
- P2 must propose x (not y) to maintain safety
- P2 sends (2, x) to A2 and A3
- Both accept (2, x)

### Why the Dubious Sequence Cannot Happen

**Detailed Explanation:**
1. **Prepare phase information**: When P2 prepares, it learns about previously accepted proposals
2. **Safety requirement**: P2 must propose the value from the highest-numbered proposal it learned about
3. **Value constraint**: Since A3 accepted (1, x), P2 must propose x, not y
4. **Safety guarantee**: This ensures that once a value is chosen, all future chosen values must be the same

### What Actually Happens

**Corrected Execution:**
1. P1 prepares proposal 1, gets responses from A1, A2, A3
2. P1 sends (1, x) to A1, A3 (A2's message dropped)
3. A1, A3 accept (1, x), P1 informs L that x is chosen
4. P1 crashes
5. P2 prepares proposal 2, gets responses from A2, A3
6. **A3 reports that it accepted (1, x)**
7. P2 sends (2, x) to A2, A3 (not (2, y))
8. A2, A3 accept (2, x)
9. P2 informs L that x is chosen

### Which Value is Ultimately Chosen

**Answer**: Value x is ultimately chosen.

**Reasoning:**
- P1's proposal (1, x) was accepted by majority (A1, A3)
- P2's proposal (2, x) was also accepted by majority (A2, A3)
- Both proposals have the same value x, satisfying safety
- The learner L receives notifications that x is chosen (from both P1 and P2)

**Key Learning Points:**
- Paxos safety is maintained through the prepare phase, which informs proposers about previously accepted proposals
- Proposers must propose values from previously accepted proposals to maintain safety
- The prepare phase acts as a "safety check" that prevents conflicting values from being chosen
- Even with message drops and crashes, Paxos ensures that only one value can be chosen

## Question 4: Paxos Liveness Issues

**Context**: While Paxos guarantees safety (only one value can be chosen), it does not guarantee liveness (progress) in all scenarios. Understanding liveness issues is crucial for designing practical consensus systems and implementing solutions like distinguished proposers or leader election.

In the absence of a distinguished proposer, it is possible for Paxos to fail to make progress even if no messages are dropped and no nodes fail. Describe in detail how this can happen in a system with two proposers and three acceptors.

**System Setup:**
- **Proposers**: P1, P2
- **Acceptors**: A1, A2, A3
- **Network**: No message drops, no node failures
- **Problem**: System fails to make progress (no value is chosen)

**Detailed Hints:**

### Understanding Liveness Issues
- **Think about**: What could prevent Paxos from making progress?
- **Consider**: How do competing proposers interact in Paxos?
- **Key insight**: Multiple proposers can interfere with each other's progress
- **Analysis approach**: Construct a scenario where proposers continuously interfere

### Step-by-step Scenario Construction

#### Step 1: Initial State
**Detailed Hints:**
- **Think about**: What is the starting state of the system?
- **Consider**: What do acceptors know initially?
- **Key insight**: All acceptors start with no accepted proposals

**Initial state:**
- A1: ⊥ (no accepted proposal)
- A2: ⊥ (no accepted proposal)  
- A3: ⊥ (no accepted proposal)

#### Step 2: First Proposer Attempts
**Detailed Hints:**
- **Think about**: What happens when P1 tries to propose?
- **Consider**: What messages does P1 send and what responses does it get?
- **Key insight**: P1 needs majority responses to proceed

**P1's first attempt:**
1. P1 sends prepare(1) to A1, A2, A3
2. All acceptors respond with "no previous proposal accepted"
3. P1 sends accept(1, x) to A1, A2, A3
4. All acceptors accept (1, x)
5. P1 concludes x is chosen

#### Step 3: Second Proposer Interferes
**Detailed Hints:**
- **Think about**: What happens when P2 tries to propose?
- **Consider**: What does P2 learn from its prepare phase?
- **Key insight**: P2 learns about P1's accepted proposal

**P2's first attempt:**
1. P2 sends prepare(2) to A1, A2, A3
2. All acceptors respond with "I accepted (1, x)"
3. P2 must propose x (not its own value)
4. P2 sends accept(2, x) to A1, A2, A3
5. All acceptors accept (2, x)
6. P2 concludes x is chosen

#### Step 4: The Liveness Problem
**Detailed Hints:**
- **Think about**: What happens if P1 tries to propose again?
- **Consider**: How can proposers interfere with each other?
- **Key insight**: Proposers can continuously interfere by proposing higher numbers

**The problematic sequence:**
1. P1 sends prepare(3) to A1, A2, A3
2. All acceptors respond with "I accepted (2, x)"
3. P1 must propose x, sends accept(3, x) to A1, A2, A3
4. All acceptors accept (3, x)
5. P2 sends prepare(4) to A1, A2, A3
6. All acceptors respond with "I accepted (3, x)"
7. P2 must propose x, sends accept(4, x) to A1, A2, A3
8. All acceptors accept (4, x)
9. **This pattern continues indefinitely...**

### Detailed Scenario: Continuous Interference

**The Complete Liveness Failure Scenario:**

**Round 1:**
- P1: prepare(1) → all acceptors → accept(1, x) → all acceptors
- P2: prepare(2) → all acceptors → accept(2, x) → all acceptors

**Round 2:**
- P1: prepare(3) → all acceptors → accept(3, x) → all acceptors
- P2: prepare(4) → all acceptors → accept(4, x) → all acceptors

**Round 3:**
- P1: prepare(5) → all acceptors → accept(5, x) → all acceptors
- P2: prepare(6) → all acceptors → accept(6, x) → all acceptors

**And so on...**

### Why This Causes Liveness Failure

**Detailed Explanation:**
1. **Continuous interference**: Each proposer's prepare phase invalidates the previous proposer's work
2. **No progress**: Although proposals are accepted, no proposer can complete its full protocol
3. **Infinite loop**: The pattern repeats indefinitely with higher proposal numbers
4. **Safety maintained**: All proposals have the same value x, so safety is preserved
5. **Liveness violated**: No value is ever "chosen" in the sense that no proposer completes the protocol

### Message Ordering Requirements

**Specific Message Delivery Order:**
1. P1's prepare(1) messages arrive at all acceptors
2. All acceptors respond to P1
3. P1's accept(1, x) messages arrive at all acceptors
4. All acceptors accept (1, x)
5. P2's prepare(2) messages arrive at all acceptors
6. All acceptors respond to P2
7. P2's accept(2, x) messages arrive at all acceptors
8. All acceptors accept (2, x)
9. **Pattern repeats with higher proposal numbers**

### Solutions to Liveness Issues

**Common approaches:**
1. **Distinguished proposer**: Only one proposer is active at a time
2. **Leader election**: Use a separate protocol to elect a leader
3. **Randomized backoff**: Proposers wait random amounts of time before retrying
4. **Exponential backoff**: Increase wait time between retries

**Key Learning Points:**
- Paxos safety is guaranteed even with multiple proposers
- Paxos liveness requires additional mechanisms (like distinguished proposer)
- Multiple proposers can interfere with each other's progress
- The prepare phase is crucial for safety but can cause liveness issues
- Practical Paxos implementations need leader election or similar mechanisms

## Question 5: Alternative Paxos Implementation Analysis

**Context**: This question explores the subtle but crucial differences in how we define "chosen" values in consensus protocols. Understanding why the original Paxos definition is necessary helps appreciate the precision required in distributed systems protocols and the potential safety violations that can arise from seemingly minor changes.

**Original Paxos Definition (from Paxos Made Simple, page 3):**
> "A value is chosen when a single proposal with that value has been accepted by a majority of the acceptors."

**Alternative Definition:**
> "A value is chosen when proposals with that value have been accepted by a majority of the acceptors."

**Analysis Question:** Would the resulting implementation be correct? Justify your answer with detailed reasoning.

**Detailed Hints:**

### Understanding the Key Difference
- **Think about**: What is the difference between "a single proposal" and "proposals"?
- **Consider**: Can multiple proposals with the same value be accepted by a majority?
- **Key insight**: The alternative definition allows multiple proposals with the same value to be considered "chosen"
- **Analysis approach**: Construct a scenario where this difference matters

### Step-by-step Analysis

#### Step 1: Understand the Original Definition
**Detailed Hints:**
- **Think about**: What does "a single proposal" mean in practice?
- **Consider**: How does this constrain the system's behavior?
- **Key insight**: Only one specific proposal (with a unique number) can be chosen

**Original definition implications:**
- A value is chosen only when one specific proposal with that value is accepted by majority
- This ensures that the chosen value comes from a single, well-defined proposal
- The proposal number provides a unique identifier for the chosen proposal

#### Step 2: Understand the Alternative Definition
**Detailed Hints:**
- **Think about**: What does "proposals with that value" mean?
- **Consider**: How does this differ from the original definition?
- **Key insight**: Multiple proposals with the same value could be considered "chosen"

**Alternative definition implications:**
- A value is chosen when any proposals with that value are accepted by majority
- This allows multiple proposals with the same value to be considered "chosen"
- The specific proposal numbers become less important

#### Step 3: Construct a Safety Violation Scenario
**Detailed Hints:**
- **Think about**: How could the alternative definition lead to safety violations?
- **Consider**: What happens when multiple proposals with the same value are accepted?
- **Key insight**: The alternative definition could allow conflicting values to be chosen

**Safety Violation Scenario:**

**System Setup:**
- Acceptors: A1, A2, A3
- Proposers: P1, P2

**Execution:**
1. P1 sends prepare(1) to A1, A2, A3
2. All acceptors respond with "no previous proposal"
3. P1 sends accept(1, x) to A1, A2, A3
4. A1 and A2 accept (1, x), A3's message is delayed
5. P2 sends prepare(2) to A1, A2, A3
6. A1 and A2 respond with "I accepted (1, x)"
7. A3 responds with "no previous proposal" (message delayed)
8. P2 sends accept(2, y) to A1, A2, A3
9. A2 and A3 accept (2, y), A1's message is delayed
10. A3's delayed accept(1, x) message arrives and is accepted
11. A1's delayed accept(2, y) message arrives and is accepted

**Final State:**
- A1: (2, y) - accepted proposal 2 with value y
- A2: (2, y) - accepted proposal 2 with value y  
- A3: (1, x) - accepted proposal 1 with value x

**Analysis with Alternative Definition:**
- Value x: Accepted by A3 (1 acceptor < majority) - NOT chosen
- Value y: Accepted by A1, A2 (2 acceptors < majority) - NOT chosen
- **No value is chosen according to the alternative definition**

**Analysis with Original Definition:**
- Proposal (1, x): Accepted by A3 only (1 acceptor < majority) - NOT chosen
- Proposal (2, y): Accepted by A1, A2 (2 acceptors < majority) - NOT chosen
- **No proposal is chosen according to the original definition**

#### Step 4: Construct a More Problematic Scenario
**Detailed Hints:**
- **Think about**: Can we construct a scenario where the alternative definition allows conflicting values?
- **Consider**: What happens with network partitions and message reordering?
- **Key insight**: The alternative definition could allow multiple values to be "chosen" simultaneously

**More Problematic Scenario:**

**Execution with Message Reordering:**
1. P1 sends prepare(1) to A1, A2, A3
2. All acceptors respond with "no previous proposal"
3. P1 sends accept(1, x) to A1, A2, A3
4. A1 and A2 accept (1, x), A3's message is delayed
5. P2 sends prepare(2) to A1, A2, A3
6. A1 and A2 respond with "I accepted (1, x)"
7. A3 responds with "no previous proposal" (message delayed)
8. P2 sends accept(2, y) to A1, A2, A3
9. A2 and A3 accept (2, y), A1's message is delayed
10. A3's delayed accept(1, x) message arrives and is accepted
11. A1's delayed accept(2, y) message arrives and is accepted

**Final State:**
- A1: (2, y) - accepted proposal 2 with value y
- A2: (2, y) - accepted proposal 2 with value y
- A3: (1, x) - accepted proposal 1 with value x

**Critical Analysis:**
- **Original definition**: No single proposal is accepted by majority - no value chosen
- **Alternative definition**: No value is accepted by majority - no value chosen

#### Step 5: The Real Safety Issue
**Detailed Hints:**
- **Think about**: What happens when we consider the system's evolution over time?
- **Consider**: How does the alternative definition affect the safety invariant?
- **Key insight**: The alternative definition could allow the safety invariant to be violated

**The Real Problem:**

The alternative definition creates ambiguity about which specific proposal is chosen, which can lead to:

1. **Inconsistent learning**: Different learners might consider different proposals as "chosen"
2. **Safety invariant violation**: The system might not maintain the property that only one value can be chosen
3. **Implementation complexity**: It becomes unclear which proposal to use for future decisions

### Conclusion

**Answer**: The alternative implementation would **not be correct**.

**Detailed Justification:**

1. **Safety violation**: The alternative definition allows multiple proposals with the same value to be considered "chosen," which can lead to inconsistent system behavior

2. **Ambiguity**: Without specifying which specific proposal is chosen, the system loses the ability to make consistent decisions about future proposals

3. **Invariant violation**: The alternative definition can violate the fundamental Paxos safety invariant that only one value can be chosen

4. **Implementation issues**: The alternative definition makes it unclear how to handle future proposals and maintain consistency

**Key Learning Points:**
- The precision of definitions in distributed systems protocols is crucial for correctness
- Seemingly minor changes to protocol definitions can have significant safety implications
- The original Paxos definition ensures that only one specific proposal is chosen, maintaining safety
- Alternative definitions that appear equivalent can actually violate fundamental safety properties
- Understanding the subtle differences between protocol definitions is essential for system design