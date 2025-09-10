# Honey Badger BFT: Asynchronous Byzantine Fault Tolerance

## Introduction: The Problem with PBFT

### How to Break PBFT Even with Long Periods of Synchrony?

**How to break PBFT even with long periods of synchrony?**

**One failed replica**

**Network starts working only when failed replica is leader**

The Honey Badger BFT paper addresses a fundamental weakness in Practical Byzantine Fault Tolerance (PBFT) that persists even when the network is synchronous for long periods. The attack scenario is surprisingly simple:

#### The Attack Scenario

**One failed replica**

**Network starts working only when failed replica is leader**

The attack works as follows:

1. **Single failed replica**: Only one replica needs to be Byzantine (failed)
2. **Strategic timing**: The network becomes synchronous only when the failed replica is the leader
3. **Liveness violation**: During these synchronous periods, the failed leader can prevent progress
4. **Censorship**: The failed leader can censor transactions or prevent consensus from proceeding

This demonstrates that PBFT's liveness guarantees are fragile and can be broken by a single malicious node with perfect timing.

#### Why This Matters

This attack is significant because:

- **Minimal requirements**: Only requires one Byzantine node and perfect timing
- **Liveness violation**: Breaks the fundamental liveness property of consensus
- **Practical relevance**: Shows that PBFT is vulnerable even in favorable network conditions
- **Motivation for alternatives**: Demonstrates the need for truly asynchronous consensus protocols

## The Honey Badger Solution: Asynchronous BFT

### Liveness and Censorship Resistance

**What is the liveness/censorship resistant property we want?**

**If N-f nodes have transaction tx in queue, will eventually be output**

Honey Badger BFT provides a stronger liveness guarantee than PBFT:

#### The Censorship Resistance Property

**If N-f nodes have transaction tx in queue, will eventually be output**

The key property is:

- **Input requirement**: If N-f (majority of correct) nodes have a transaction in their input queue
- **Output guarantee**: The transaction will eventually be output (committed)
- **Censorship resistance**: No single node or small group can prevent this from happening
- **Asynchronous guarantee**: This holds even in completely asynchronous networks

This is much stronger than PBFT's liveness guarantees, which can be broken by a single malicious leader.

### The Mathematical Proof

**Why does this hold? Say tx is enqueued with backlog T**

**In any round, proposals must contain at least N/6 Type 1 or N/6 Type 2 nodes**

**Type 1 - correct node with tx in front (first B elements) of queue**

**Commits tx with probability 1-e^{-1/6}**

**So overwhelmingly likely to commit after \lambda rounds**

**Type 2 - correct node with tx not in first B elements of queue**

**Will clear expected >= B(1-e^{-1/6}) tx from backlog**

**So will become Type 1 node eventually**

The proof relies on a clever analysis of node types and their behavior:

#### Type 1 Nodes

**Type 1 - correct node with tx in front (first B elements) of queue**

**Commits tx with probability 1-e^{-1/6}**

**So overwhelmingly likely to commit after \lambda rounds**

Type 1 nodes are those that have the target transaction in the first B elements of their queue:

- **High priority**: The transaction is in the front of their queue
- **Commit probability**: Each Type 1 node commits the transaction with probability 1-e^(-1/6) ≈ 0.15
- **Rapid convergence**: After λ rounds, it's overwhelmingly likely that the transaction will be committed
- **Key insight**: These nodes are the primary drivers of transaction commitment

#### Type 2 Nodes

**Type 2 - correct node with tx not in first B elements of queue**

**Will clear expected >= B(1-e^{-1/6}) tx from backlog**

**So will become Type 1 node eventually**

Type 2 nodes are those that have the target transaction but not in the front of their queue:

- **Lower priority**: The transaction is not in the first B elements
- **Backlog clearing**: Each Type 2 node clears an expected ≥ B(1-e^(-1/6)) transactions from the backlog
- **Type conversion**: Eventually, Type 2 nodes become Type 1 nodes as the backlog clears
- **Progress guarantee**: This ensures that the system makes progress toward committing the transaction

#### The Proof Structure

**In any round, proposals must contain at least N/6 Type 1 or N/6 Type 2 nodes**

The proof works by showing that:

1. **Minimum node requirement**: In any round, at least N/6 nodes must be either Type 1 or Type 2
2. **Progress guarantee**: Type 1 nodes directly commit the transaction with high probability
3. **Backlog reduction**: Type 2 nodes reduce the backlog, eventually becoming Type 1 nodes
4. **Convergence**: The system converges to a state where the transaction is committed

This mathematical analysis provides a rigorous foundation for Honey Badger's liveness guarantees.

## Efficiency Analysis

### Measuring Efficiency in Honey Badger BFT

**How do we measure efficiency in HoneyBadger BFT?**

**Expected per-node communication per committed transaction**

**Why does this break down if input queues are not full?**

Honey Badger BFT uses a specific efficiency metric:

#### The Efficiency Metric

**Expected per-node communication per committed transaction**

The efficiency is measured as:

- **Per-node basis**: Communication cost per individual node
- **Per-transaction basis**: Cost normalized by the number of committed transactions
- **Expected value**: Average cost over many transactions
- **Communication focus**: Primarily measures network communication, not computation

#### When the Metric Breaks Down

**Why does this break down if input queues are not full?**

The efficiency metric has limitations:

- **Queue dependency**: The metric assumes input queues are full
- **Batch efficiency**: Honey Badger is most efficient when processing batches of transactions
- **Underutilization**: If queues are not full, the per-transaction cost increases
- **Real-world scenarios**: In practice, queues may not always be full, affecting efficiency

### Asymptotic Communication Cost

**What is the asymptotic communication cost, per node?**

**O(N^2) costs:**

**- each node sends O(N^2) shares for common coin in ABA**

**- RBC ECHO requires echoing O(N^2) hashes**

**- threshold decryption requires O(N^2) decryption shares**

**Costs proportional to B:**

**- Actual echoed data in RBC**

The communication cost has two main components:

#### O(N^2) Costs

**- each node sends O(N^2) shares for common coin in ABA**

**- RBC ECHO requires echoing O(N^2) hashes**

**- threshold decryption requires O(N^2) decryption shares**

The O(N^2) costs come from:

1. **Common coin in ABA**: Each node sends O(N^2) shares for the common coin in Asynchronous Binary Agreement
2. **RBC ECHO**: Reliable Broadcast ECHO requires echoing O(N^2) hashes
3. **Threshold decryption**: Threshold decryption requires O(N^2) decryption shares

These costs are inherent to the cryptographic primitives used in Honey Badger.

#### Costs Proportional to B

**- Actual echoed data in RBC**

The costs proportional to B include:

- **Data transmission**: The actual data being broadcast in Reliable Broadcast
- **Batch size dependency**: These costs scale with the batch size B
- **Efficient component**: This is the most efficient part of the protocol

The total cost is O(N^2) + O(B), where the O(N^2) term dominates for large N.

## Experimental Evaluation

### Figure 5: A Questionable Experiment

**Is Fig. 5 a convincing experiment?**

**No--not an experiment at all, just a graph**

Figure 5 in the Honey Badger paper has significant limitations:

#### What Figure 5 Shows

**No--not an experiment at all, just a graph**

The figure appears to be:

- **Theoretical analysis**: A graph showing theoretical performance, not experimental results
- **No real data**: No actual measurements from running the protocol
- **Misleading presentation**: Presented as if it were experimental data
- **Limited value**: Provides little insight into real-world performance

#### Why This Matters

This is problematic because:

- **Misleading**: Readers might think this represents actual performance
- **No validation**: The theoretical analysis is not validated by experiments
- **Limited insights**: Doesn't provide insights into practical implementation challenges
- **Credibility issues**: Raises questions about the paper's experimental rigor

### Figure 7: Performance Insights

**What does Fig. 7 tell us?**

**Can't saturate throughput without high latency?**

**Why is 104 nodes not scaling?**

**Implementation bottlenecked on CPU because crypto is single threaded**

Figure 7 provides more meaningful insights:

#### Throughput vs. Latency Trade-off

**Can't saturate throughput without high latency?**

The figure shows:

- **Trade-off relationship**: Higher throughput comes at the cost of higher latency
- **Saturation limits**: There's a limit to how much throughput can be achieved
- **Latency penalty**: Achieving high throughput requires accepting higher latency
- **Practical implications**: Real applications must balance these competing goals

#### Scaling Limitations

**Why is 104 nodes not scaling?**

**Implementation bottlenecked on CPU because crypto is single threaded**

The scaling issues are due to:

- **CPU bottleneck**: The implementation is limited by CPU performance
- **Single-threaded crypto**: Cryptographic operations are not parallelized
- **Implementation limitation**: This is a limitation of the specific implementation, not the protocol
- **Optimization opportunity**: Better implementations could improve scaling

### Figure 8: PBFT Performance Analysis

**Fig 8.: Why does PBFT do worse with more cohorts?**

**PBFT leader must send all transactions to other cohorts**

**PBFT cohorts send less data, so save on average, but leader bottlenecks**

**Is this fundamental to async. vs weak synchrony? not really**

**How could you fix leader bottleneck with PBFT?**

**Disseminate each block with P2P protocol (e.g., bittorrent)**

**Or better, leader encodes blocks with erasure code**

**disseminate coded blocks to cohorts who share with each other**

**rateless erasure code (t,\infty) makes particularly easy**

Figure 8 reveals important insights about PBFT's performance:

#### The Leader Bottleneck

**PBFT leader must send all transactions to other cohorts**

**PBFT cohorts send less data, so save on average, but leader bottlenecks**

The problem is:

- **Leader responsibility**: The leader must send all transactions to all cohorts
- **Asymmetric load**: Cohorts send less data, but the leader is overloaded
- **Bottleneck effect**: The leader becomes the bottleneck as the number of cohorts increases
- **Scaling limitation**: This limits PBFT's ability to scale to large numbers of nodes

#### Is This Fundamental?

**Is this fundamental to async. vs weak synchrony? not really**

The bottleneck is not fundamental to the synchrony model:

- **Implementation issue**: This is primarily an implementation issue, not a fundamental protocol limitation
- **Synchrony independent**: The bottleneck exists regardless of the synchrony model
- **Optimization possible**: The bottleneck can be addressed through better implementation techniques

#### Potential Solutions

**How could you fix leader bottleneck with PBFT?**

**Disseminate each block with P2P protocol (e.g., bittorrent)**

**Or better, leader encodes blocks with erasure code**

**disseminate coded blocks to cohorts who share with each other**

**rateless erasure code (t,\infty) makes particularly easy**

Several solutions are possible:

1. **P2P dissemination**: Use peer-to-peer protocols like BitTorrent to distribute blocks
2. **Erasure coding**: Encode blocks with erasure codes to reduce the leader's load
3. **Coded block sharing**: Have cohorts share coded blocks with each other
4. **Rateless codes**: Use rateless erasure codes (t,∞) for particularly easy implementation

These solutions could significantly improve PBFT's scalability.

## Comparison with PBFT

### Is Honey Badger Better Than PBFT?

**Is HoneyBadger better than PBFT over ToR?**

**Probably with high node count, but we don't know for sure**

**Paper should have benchmarked--all we know is less bandwidth with many nodes**

The comparison between Honey Badger and PBFT is complex:

#### Advantages of Honey Badger

**Probably with high node count, but we don't know for sure**

Honey Badger has several advantages:

- **Asynchronous guarantee**: Works in completely asynchronous networks
- **Censorship resistance**: Stronger liveness guarantees than PBFT
- **No leader bottleneck**: No single point of failure like PBFT's leader
- **Better scaling**: Potentially better performance with high node counts

#### Limitations of the Comparison

**Paper should have benchmarked--all we know is less bandwidth with many nodes**

The paper's comparison is limited:

- **Insufficient benchmarking**: The paper should have included more comprehensive benchmarks
- **Limited metrics**: Only bandwidth is compared, not latency or other important metrics
- **Uncertain conclusions**: It's unclear whether Honey Badger is actually better in practice
- **Implementation gaps**: The comparison may be limited by implementation quality

#### What We Know

**all we know is less bandwidth with many nodes**

The only clear advantage demonstrated is:

- **Bandwidth efficiency**: Honey Badger uses less bandwidth with many nodes
- **Scaling benefit**: The advantage becomes more pronounced with larger node counts
- **Limited scope**: This is only one aspect of performance

## Conclusion

Honey Badger BFT represents a significant advancement in asynchronous Byzantine fault tolerance. It addresses fundamental weaknesses in PBFT, particularly the vulnerability to liveness attacks and the leader bottleneck problem.

### Key Contributions

1. **Asynchronous BFT**: Provides Byzantine fault tolerance in completely asynchronous networks
2. **Censorship resistance**: Stronger liveness guarantees than previous protocols
3. **No leader bottleneck**: Eliminates the single point of failure present in PBFT
4. **Mathematical rigor**: Provides rigorous mathematical analysis of liveness properties

### Limitations and Future Work

1. **Experimental validation**: More comprehensive experimental evaluation is needed
2. **Implementation optimization**: Better implementations could improve performance
3. **Comparison studies**: More thorough comparison with other BFT protocols is needed
4. **Real-world deployment**: Practical deployment challenges need to be addressed

### Impact and Significance

Honey Badger BFT has had significant impact on the field of distributed systems:

- **Theoretical contribution**: Advanced the theoretical understanding of asynchronous BFT
- **Practical relevance**: Addressed real-world problems with existing BFT protocols
- **Research direction**: Influenced subsequent research in asynchronous consensus
- **Implementation basis**: Served as a foundation for practical BFT implementations

The protocol demonstrates that it's possible to achieve Byzantine fault tolerance in asynchronous networks while maintaining strong liveness guarantees, opening new possibilities for robust distributed systems.