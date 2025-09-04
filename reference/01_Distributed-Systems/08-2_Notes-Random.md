# Randomized Consensus: Supplementary Notes

## FLP Impossibility

### The Problem
- **Theorem**: In asynchronous environment with single process crash failure, no protocol solves binary consensus
- **Paxos doesn't save us**: It doesn't guarantee liveness
- **Result assumes**: Deterministic computation model

### The Escape Hatch
- **Let's go random!**: Use randomization to break FLP impossibility
- **Ben-Or's algorithm**: Guarantees consensus for crash failures when f < n/2
- **Variant works**: Even for Byzantine faults

## Intuition

### Basic Idea
- **At first**: Every process proposes their input value
- **After that**: They propose random values
- **When enough processes** propose the same value, the value is chosen
- **Eventually**: That will happen!

### Key Insight
- **Randomization breaks symmetry**: Allows processes to escape from deadlock
- **Probabilistic termination**: Guarantees eventual consensus with probability 1

## Ben-Or Algorithm Setup

### Model
- **Binary consensus**: Values are 0 or 1
- **Asynchronous rounds**: Each round has two phases
- **Message requirements**: Each phase waits for n - f messages
- **Message tagging**: Round and phase numbers
- **Network handling**: Messages can be resent for lossy networks

### Protocol Structure
- **Processes send proposals** for each phase
- **Block and wait** for n - f messages (including their own)
- **Values locked**: Once sent, value is locked for that process/phase/round

## Ben-Or Algorithm

### Phase 1: Preliminary Proposal
- **Processes make preliminary proposal**
- **If receive matching responses** from majority in phase 1:
  - Propose that value in phase 2
- **Otherwise**: Propose NULL (special null value)

### Phase 2: Final Proposal
- **If get enough non-NULL responses** from phase 2:
  - Decide on that value
- **Otherwise**: Continue to next round with random values

### Random Value Selection
- **When no consensus**: Processes randomly choose from {0, 1}
- **Eventually**: All processes will randomly choose same value

## Consensus Properties

### Agreement
- **No two processes** decide different values
- **Ensured by**: Majority intersection and message consistency

### Integrity
- **Every process decides** at most one value
- **If process decides value**: Some process had it as input
- **Trivially satisfied**: If both 0 and 1 are input values

### Termination
- **Every correct process** eventually decides a value
- **Probabilistic guarantee**: With probability 1

## Proof of Correctness

### Integrity I: Same Input Case
- **Suppose all processes** have same input value
- **Round 1, Phase 1**: All send same value
- **Round 1, Phase 2**: All send same value
- **Round 1, End**: All decide that value

### Key Lemma
- **No two processes** receive different non-NULL phase 2 values in same round
- **Proof**: If they did, one process received 0s from majority, another received 1s
- **Contradiction**: Majorities intersect, so impossible

### Agreement + Integrity II: First Decision
- **Let round r** be first round any process decides value (0 w.l.o.g.)
- **If process decided**: Must have received > f 0s in phase 2
- **Every process received** at least one 0 (wait for n - f messages)
- **No process received 1** by previous lemma
- **Round r + 1**: All processes propose 0 and decide 0

### Termination: Probabilistic Guarantee
- **If all processes propose same value**: They all decide that round
- **Worst case probability**: ½^n on any particular round
- **Why**: By lemma, all non-random values are identical
- **Over time**: Probability converges to 1

## Extensions

### Non-Binary Consensus
- **Binary consensus**: Conceptually simple but limited
- **Larger domains**: Algorithm can support larger value sets
- **Unknown domains**: Even when processes don't know domains a priori
- **Missing inputs**: Some processes may not receive input values

### Handling Missing Inputs
- **Processes without input**: Start by proposing NULL
- **Random selection**: Choose from all non-NULL values seen so far
- **NULL as last resort**: Only choose NULL when no other values available

## Key Takeaways

### Randomized Consensus Benefits
- **Breaks FLP impossibility**: Through randomization
- **Guarantees termination**: With probability 1
- **Works for crash failures**: When f < n/2
- **Extends to Byzantine**: With variants

### Algorithm Structure
- **Asynchronous rounds**: Two-phase structure
- **Majority waiting**: n - f messages per phase
- **Random fallback**: When consensus not reached
- **Probabilistic termination**: Eventually succeeds

### Design Principles
- **Use randomization**: To break symmetry
- **Structure with rounds**: Even in asynchronous systems
- **Wait for majorities**: To ensure consistency
- **Handle missing values**: With NULL proposals

### Trade-offs
- **Probabilistic vs Deterministic**: Guarantees termination but not when
- **Performance vs Correctness**: May take many rounds
- **Simplicity vs Robustness**: Simple algorithm, complex analysis
- **Synchrony vs Asynchrony**: Works in fully asynchronous model

### Applications
- **Distributed systems**: When deterministic consensus impossible
- **Byzantine fault tolerance**: With appropriate variants
- **Cryptocurrency**: Blockchain consensus mechanisms
- **Distributed databases**: When strong consistency needed

### Limitations
- **Probabilistic termination**: No bound on number of rounds
- **Performance**: May require many rounds in worst case
- **Analysis complexity**: Probabilistic correctness proofs
- **Implementation challenges**: True randomness requirements