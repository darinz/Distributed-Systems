# Streamlet: A Simple Blockchain Consensus Protocol

## Introduction: Synchrony Models and Consensus

### The Synchrony Spectrum

**Many real-world consensus protocols (e.g., PBFT) live under *weak synchrony***

**Network/node delays grow at a bounded rate (e.g., polynomially)**

**Could easily fix RAFT to be this way--just adapt to network**

**Keep doubling timeout until progress**

Understanding different synchrony models is crucial for designing consensus protocols. The synchrony spectrum ranges from completely asynchronous (no timing assumptions) to completely synchronous (perfect timing guarantees).

#### Weak Synchrony

Weak synchrony assumes that network and node delays grow at a bounded rate, typically polynomially. This is a practical assumption for many real-world systems where:

- **Bounded growth**: Delays don't grow exponentially or unboundedly
- **Adaptive protocols**: Protocols can adapt to changing network conditions
- **Timeout adjustment**: Systems can increase timeouts when delays increase

**Could easily fix RAFT to be this way--just adapt to network**

**Keep doubling timeout until progress**

This adaptive approach works by:
1. **Starting with reasonable timeouts**: Begin with timeouts based on expected network conditions
2. **Doubling on failure**: When timeouts expire without progress, double the timeout
3. **Eventually succeeding**: Since delays are bounded, eventually timeouts will be large enough

#### Partial Synchrony

**Another model: *partial synchrony***

**Initially, things may be very crazy**

**Eventually reach some global stabilization time GST**

**After GST, nodes can communicate within some timebound Delta**

**The catch: don't know GST (so still not quite synchronous)**

Partial synchrony is a more sophisticated model that captures the reality of many distributed systems:

- **Initial chaos**: Before the Global Stabilization Time (GST), the network can be completely unreliable
- **Post-GST stability**: After GST, the network becomes reliable with bounded delays
- **Unknown GST**: The critical challenge is that GST is unknown to the protocol

**Another variation: the unknown Delta model**

**Synchronous, but don't know what delta is**

The unknown Delta model is a variation where:
- **Synchronous after GST**: The network is synchronous after GST
- **Unknown bound**: The delay bound Delta is unknown
- **Adaptive approach**: Protocols must adapt to discover the correct Delta

#### Model Equivalence

**Two models have same power**

**E.g., increasing estimate of Delta is like increasing timeouts**

These models are equivalent in their power because:
- **Adaptive strategies**: Both can be handled with adaptive timeout strategies
- **Timeout scaling**: Increasing estimates of Delta is equivalent to increasing timeouts
- **Practical equivalence**: In practice, both approaches lead to similar protocol designs

## Blockchain Motivation: From Credit Cards to Consensus

### The Digital Payment Problem

**Blockchain digression: Say we replace credit card #s w. public keys?**

**Hold private key on cell phone**

**Sign each transaction with amount, recipient, comment**

**What would this fix?**

**Credit card numbers would be harder to steal**

**Lower fraud could lead to lower (but still non-zero) transaction fees**

**Could also lower fraud-detection false positives**

**Could give merchants protection (eliminate "card not present" transactions)**

The motivation for blockchain technology often starts with the problems of traditional digital payments. Consider replacing credit card numbers with public key cryptography:

#### Advantages of Public Key Payments

**What would this fix?**

- **Security**: Credit card numbers would be harder to steal than physical cards
- **Lower fraud**: Reduced fraud could lead to lower transaction fees
- **Fewer false positives**: Could reduce fraud-detection false positives
- **Merchant protection**: Could eliminate "card not present" transaction risks

#### Comparison with Cash

**How does this compare to cash?**

**Worse privacy, third party trust, bank control, transaction charges**

However, this approach still has significant disadvantages compared to cash:

- **Privacy**: Digital transactions are less private than cash
- **Third party trust**: Still requires trust in payment processors
- **Bank control**: Banks still control the payment system
- **Transaction charges**: Still incurs transaction fees

### Bitcoin's Innovation

**Bitcoin attempts to get rid of the bank...**

**Transfer ownership of digital BTC between different public keys**

Bitcoin's key innovation is eliminating the need for trusted third parties by creating a decentralized digital currency that allows direct peer-to-peer transfers.

#### The Two Fundamental Problems

**Must solve 2 problems for this model to work:**

**1. How do you distribute Bitcoin so people believe it has value?**

**2. How do you prevent someone from double-spending a Bitcoin?**

**Public key A is worth 1 BTC**

**A pays 1 BTC to B, walks away with some goods**

**A tries to pay 1 BTC to C, walk away with more goods**

For a decentralized digital currency to work, two fundamental problems must be solved:

#### Problem 1: Value Distribution
How do you create and distribute a digital currency so that people believe it has value? This involves:
- **Initial distribution**: How are bitcoins initially created and distributed?
- **Value perception**: Why would people accept bitcoins as payment?
- **Scarcity**: How do you ensure the currency has limited supply?

#### Problem 2: Double-Spending Prevention
How do you prevent someone from spending the same bitcoin twice? This is the classic double-spending problem:

**Note that problem #2 is basically consensus--"timestamp server"**

**With timestamp server, everyone agrees A->C message later than A->B**

**Hence A->C is not valid, because A is already spent**

**Keep ordered transaction history as a giant replicated state machine**

The double-spending problem is fundamentally a consensus problem. The solution is to maintain an ordered transaction history as a replicated state machine where:

- **Global ordering**: Everyone agrees on the order of transactions
- **Timestamp server**: A mechanism to determine which transaction happened first
- **State machine**: The blockchain serves as a replicated state machine tracking ownership

## Bitcoin's Consensus: Proof-of-Work Mining

### How Bitcoin Achieves Consensus

**How does Bitcoin do consensus? Proof-of-work mining**

**Flood new transactions**

**History progresses in blocks of new transactions:**

**block = { H(previous-block), nonce, H*(new-transactions) }**

**Note inclusion of previous block makes this a *block chain***

Bitcoin solves the consensus problem through a clever proof-of-work mechanism that creates a blockchain:

#### Block Structure

Each block contains:
- **H(previous-block)**: Hash of the previous block, creating a chain
- **nonce**: A random number that miners adjust to find valid blocks
- **H*(new-transactions)**: Hash of the transactions in this block

The inclusion of the previous block's hash creates a blockchain where each block is cryptographically linked to its predecessor.

#### Proof-of-Work Mechanism

**By convention, only accept block if H(block) < target-value**

**Currently H(block) must start with >=81 0 bits**

**So must try an expected 2^81 nonces+blocks to find good one**

**Finding these nonces, known as "mining", is computationally very hard**

The proof-of-work mechanism works as follows:

1. **Target difficulty**: Only blocks with hash values below a target are accepted
2. **Computational work**: Finding such blocks requires trying many different nonces
3. **Difficulty adjustment**: The target is adjusted to maintain consistent block times
4. **Mining process**: Miners compete to find valid blocks by trying different nonces

#### Incentive Mechanism

**So incentivize finding the next block with bitcoins!**

**First transaction in block is special "coinbase" transaction**

**Creates new BTC paid to block miner**

**Note: adjust target-value as more miners come on-line**

**In practice, new blocks mined every 7-10 minutes or so**

Bitcoin incentivizes mining through:

- **Block rewards**: New bitcoins are created and paid to the miner who finds each block
- **Coinbase transaction**: The first transaction in each block creates new bitcoins
- **Difficulty adjustment**: The target is adjusted based on network hash rate
- **Consistent timing**: Blocks are mined approximately every 10 minutes

#### Chain Selection Rule

**Given two incompatible block chains, take one with most work (~tallest)**

**Note we've just solved both the coin distribution and consensus problems!**

When there are multiple competing chains, Bitcoin uses the "longest chain" rule:

- **Most work**: Choose the chain with the most computational work (longest chain)
- **Fork resolution**: This resolves forks and ensures consensus
- **Security**: The longest chain represents the majority of computational power

This mechanism solves both fundamental problems:
1. **Coin distribution**: New bitcoins are created and distributed to miners
2. **Consensus**: The longest chain rule ensures everyone agrees on the transaction history
### How Bitcoin Prevents Double-Spending

**How does this solve the double-spending problem?**

**Well-behaved nodes ordinarily shouldn't fork block chain**

**Assume ill-behaved nodes have less aggregate compute than well-behaved ones**

**(Though 49% limit is inadequate http://arxiv.org/abs/1311.0243)**

**If bad buys don't immediately fork chain and double spend, can't catch up**

**After receiving payment, before handing over goods,**

**wait for several new blocks (e.g., ~1hr) to cement payment in history**

**And note mining payouts incentivize honest behavior even by greedy miners**

Bitcoin's proof-of-work mechanism prevents double-spending through several mechanisms:

#### The Honest Majority Assumption

**Assume ill-behaved nodes have less aggregate compute than well-behaved ones**

Bitcoin's security relies on the assumption that honest nodes control more than 50% of the network's computational power. This means:

- **Honest majority**: Well-behaved nodes have more aggregate compute than malicious ones
- **Chain selection**: The honest majority will always produce the longest chain
- **Attack resistance**: Malicious nodes cannot outcompete the honest majority

**Note**: The 49% limit mentioned in the reference shows that this assumption has limitations and can be inadequate in certain scenarios.

#### The Catch-Up Problem

**If bad buys don't immediately fork chain and double spend, can't catch up**

If a malicious node tries to double-spend by creating a fork:

1. **Immediate fork**: The malicious node must immediately create a fork when the transaction is first included
2. **Catch-up requirement**: The malicious fork must catch up to and surpass the honest chain
3. **Computational disadvantage**: Since malicious nodes have less compute power, they cannot catch up

#### The Confirmation Mechanism

**After receiving payment, before handing over goods,**

**wait for several new blocks (e.g., ~1hr) to cement payment in history**

The key insight is that recipients should wait for multiple confirmations:

- **Block confirmations**: Wait for several blocks to be mined after the transaction
- **Cementing**: Each additional block makes it exponentially harder to reverse the transaction
- **Practical timing**: Waiting for about 1 hour (6 blocks) provides strong security guarantees

#### Economic Incentives

**And note mining payouts incentivize honest behavior even by greedy miners**

Bitcoin's economic model incentivizes honest behavior:

- **Mining rewards**: Honest mining is more profitable than attacking the network
- **Long-term value**: Attacking the network would destroy the value of bitcoins
- **Rational behavior**: Even greedy miners have incentives to behave honestly

### Bitcoin's Properties: Safety, Liveness, and Fault-Tolerance

**Which of safety, liveness, and fault-tolerance does Bitcoin offer?**

**Liveness and Fault-tolerance (anyone can unilaterally mine a block)**

**It's also randomized (since guessing nonces), so not subject to FLP**

**But Bitcoin not safe w/o synchrony assumption**

**Under network partition, miners create arbitrarily deep forks**

**Also not safe against computationally-powerful attackers**

Bitcoin provides different guarantees depending on the network conditions:

#### Liveness and Fault-Tolerance

**Liveness and Fault-tolerance (anyone can unilaterally mine a block)**

**It's also randomized (since guessing nonces), so not subject to FLP**

Bitcoin provides:

- **Liveness**: The system continues to make progress (new blocks are mined)
- **Fault-tolerance**: Anyone can unilaterally mine a block, so the system doesn't get stuck
- **Randomization**: The proof-of-work mechanism is randomized, so it's not subject to the FLP impossibility result

#### Safety Limitations

**But Bitcoin not safe w/o synchrony assumption**

**Under network partition, miners create arbitrarily deep forks**

**Also not safe against computationally-powerful attackers**

However, Bitcoin has significant safety limitations:

- **Synchrony requirement**: Bitcoin is not safe without synchrony assumptions
- **Network partitions**: Under network partitions, miners can create arbitrarily deep forks
- **Computational attacks**: Not safe against computationally-powerful attackers (e.g., quantum computers)

### The Need for Alternative Approaches

**Another problem: People want actual money, bank-like value guarantees**

**E.g., facebook tried (failed) to created libra/diem currency-backed blockchain**

**Used f-out-of-3f+1 BFT consensus among closed consortium**

**Chose algorithm called HotStuff (chained hotsuff similar to streamlet)**

Bitcoin's limitations have led to the development of alternative approaches:

#### Permissioned Blockchains

**People want actual money, bank-like value guarantees**

Many applications require stronger guarantees than Bitcoin provides:

- **Bank-like guarantees**: Financial applications need stronger safety properties
- **Regulatory compliance**: Many applications must meet regulatory requirements
- **Performance**: Bitcoin's 10-minute block times are too slow for many applications

**Facebook tried (failed) to created libra/diem currency-backed blockchain**

**Used f-out-of-3f+1 BFT consensus among closed consortium**

**Chose algorithm called HotStuff (chained hotsuff similar to streamlet)**

Facebook's Libra/Diem project illustrates this approach:

- **Permissioned system**: Used a closed consortium of trusted validators
- **BFT consensus**: Used Byzantine fault-tolerant consensus (f-out-of-3f+1)
- **HotStuff algorithm**: Chose a chained HotStuff algorithm similar to Streamlet

#### Non-Financial Applications

**Are there non-payment/finance applications for blockchain? Hypothetically yes**

**Timestamping - prove you created a document by a particular time**

**Transparency - commit to some public append-only history, e.g.,**

**could use to improve certificate transparency**

**iphone could refuse to install firmware not posted to blockchain**

Blockchain technology has potential applications beyond cryptocurrency:

- **Timestamping**: Prove you created a document at a particular time
- **Transparency**: Commit to a public append-only history
- **Certificate transparency**: Improve certificate transparency systems
- **Firmware verification**: iPhones could refuse to install firmware not posted to blockchain

#### The Mining Problem

**Problem: mining not appropriate for non-cryptocurrency applications**

**Streamlet motivation: support closed blockchain**

However, proof-of-work mining has significant limitations for non-cryptocurrency applications:

- **Energy consumption**: Mining consumes enormous amounts of energy
- **Economic incentives**: Mining only makes sense when there's a valuable cryptocurrency
- **Performance**: Mining is too slow for many applications
- **Centralization**: Mining tends to centralize over time

**Streamlet motivation: support closed blockchain**

This is where Streamlet comes in - it provides a consensus mechanism suitable for closed blockchains without the need for mining.

## PBFT: The Traditional BFT Consensus Approach

### PBFT vs. Raft

**Original PBFT [Castro] is vaguely like Raft, with a few big differences**

**Requires 3f+1 replicas and uses quorum size is 2f+1**

**Raft: any two quorums must share a node (so f+1 out of 2f+1 good)**

**PBFT: any two quorums must share an *honest* node**

PBFT (Practical Byzantine Fault Tolerance) is a classic BFT consensus protocol that shares some similarities with Raft but has important differences:

#### Replica Requirements

**Requires 3f+1 replicas and uses quorum size is 2f+1**

PBFT requires more replicas than Raft:

- **PBFT**: Requires 3f+1 replicas to tolerate f Byzantine failures
- **Raft**: Requires 2f+1 replicas to tolerate f crash failures
- **Quorum size**: Both use quorum size of 2f+1

#### Quorum Intersection

**Raft: any two quorums must share a node (so f+1 out of 2f+1 good)**

**PBFT: any two quorums must share an *honest* node**

The key difference is in the quorum intersection requirement:

- **Raft**: Any two quorums must share at least one node (which could be faulty)
- **PBFT**: Any two quorums must share at least one honest node
- **Byzantine tolerance**: PBFT's requirement is stronger because it must handle Byzantine failures

#### Communication Rounds

**Need extra communication round to check leader told everyone same thing**

**Leader broadcasts "pre-prepare", becomes prepare when 2f followers agree**

PBFT requires more communication rounds than Raft:

1. **Pre-prepare**: Leader broadcasts a pre-prepare message
2. **Prepare**: Followers send prepare messages when they agree
3. **Commit**: Followers send commit messages to finalize the decision

This extra round is necessary to ensure that the leader told everyone the same thing, which is crucial for Byzantine fault tolerance.

#### Leader Selection

**After timeout, round-robin through nodes for leader**

**With 2f+1 honest nodes, you eventually get an honest leader**

PBFT uses a round-robin approach for leader selection:

- **Timeout-based**: When the current leader fails, the system times out
- **Round-robin**: The system cycles through nodes to select a new leader
- **Honest leader guarantee**: With 2f+1 honest nodes, the system will eventually select an honest leader

### Why Not Use PBFT for Blockchains?

**Why not use PBFT for blockchains?**

**Short answer: you could (e.g., hyperledger fabric supports this model)**

**But PBFT is complicated**

**PBFT optimizes for latency at an irrelevant level for blockchains**

**PBFT may have bursts of requests and periods of inactivity**

**Also with many replicas, PBFT has large message sizes and counts**

**Streamlet paper doesn't address (but see HotStuff)**

While PBFT could theoretically be used for blockchains, there are several reasons why it's not ideal:

#### Complexity Issues

**But PBFT is complicated**

PBFT has several complexity issues:

- **Multiple rounds**: Requires three communication rounds per decision
- **View changes**: Complex leader change protocol
- **Message complexity**: High message overhead and complexity

#### Latency Optimization

**PBFT optimizes for latency at an irrelevant level for blockchains**

PBFT is designed to optimize for low latency, but this optimization is not relevant for blockchains:

- **Block-based**: Blockchains naturally batch transactions into blocks
- **Tolerable latency**: Blockchains can tolerate higher latency since they batch transactions
- **Different priorities**: Blockchains prioritize throughput over individual transaction latency

#### Workload Characteristics

**PBFT may have bursts of requests and periods of inactivity**

PBFT is designed for different workload characteristics:

- **Bursty traffic**: PBFT handles bursts of requests well
- **Inactivity periods**: PBFT can handle periods of inactivity
- **Blockchain workload**: Blockchains have more consistent, continuous workloads

#### Scalability Issues

**Also with many replicas, PBFT has large message sizes and counts**

**Streamlet paper doesn't address (but see HotStuff)**

PBFT has scalability issues:

- **Message complexity**: Message sizes and counts grow with the number of replicas
- **Network overhead**: High network overhead for large numbers of replicas
- **HotStuff solution**: HotStuff addresses some of these issues

### The Blockchain Advantage

**Idea: take advantage of blockchain structure to simplify consensus**

**Batching transactions into blocks anyway, so can tolerate latency**

**Each block securely identifies all prior blocks**

**Constantly producing blocks (e.g., can leverage block n+1 to finalize n)**

Blockchains have structural advantages that can be leveraged to simplify consensus:

#### Natural Batching

**Batching transactions into blocks anyway, so can tolerate latency**

- **Block structure**: Blockchains naturally batch transactions into blocks
- **Latency tolerance**: This batching allows for higher latency tolerance
- **Simplified consensus**: Can use simpler consensus mechanisms

#### Cryptographic Linking

**Each block securely identifies all prior blocks**

- **Hash chains**: Each block contains a hash of the previous block
- **Security**: This creates a secure chain of blocks
- **Simplified verification**: Makes verification of the chain simpler

#### Continuous Production

**Constantly producing blocks (e.g., can leverage block n+1 to finalize n)**

- **Continuous operation**: Blockchains continuously produce new blocks
- **Finalization**: Can use future blocks to finalize previous blocks
- **Simplified protocol**: This allows for simpler consensus protocols

### PBFT's Potential Advantages

**One potential advantage of PBFT still:**

**Stable leader can assemble next block while followers validate current**

**Potentially doubles throughput**

**Also facilitates MEV, but less a concern in closed setting**

PBFT still has some potential advantages for blockchains:

#### Throughput Optimization

**Stable leader can assemble next block while followers validate current**

**Potentially doubles throughput**

- **Parallel processing**: Leader can work on the next block while followers validate the current block
- **Throughput improvement**: This can potentially double throughput
- **Pipeline efficiency**: Creates a pipeline of block production and validation

#### MEV Considerations

**Also facilitates MEV, but less a concern in closed setting**

- **MEV (Maximal Extractable Value)**: PBFT can facilitate MEV extraction
- **Closed setting**: MEV is less of a concern in closed, permissioned blockchains
- **Trade-offs**: Must balance throughput benefits with MEV concerns

## Streamlet: The Protocol Design

### The Streamlet Setting

**The setting for Streamlet:**

**n numbered nodes, f < n/3 Byzantine faulty**

**partially synchronous with \Delta timebound after unknown GST**

**assume all nodes *implicitly echo* all messages received (wasteful)**

**So if a non-faulty node receives a message at time r,**

**all non-faulty nodes will receive it by time max(GST, r+\Delta)**

**proceed through epochs of duration 2*\Delta**

**agree on blocks of the form: (h, e, txs)**

**txs are transactions to append to blockchain history**

**e is epoch number of block - not all epochs included in finalized history**

**h is hash of previous block or special genesis block (\bot,0,\bot)**

***length* of block B, |B|, is distance to genesis block (length of chain)**

**each epoch e has a pseudo-random leader determined by H(e) mod n**

Streamlet operates in a specific setting with well-defined assumptions:

#### Network Model

**n numbered nodes, f < n/3 Byzantine faulty**

**partially synchronous with \Delta timebound after unknown GST**

Streamlet assumes:

- **Node count**: n numbered nodes in the system
- **Byzantine failures**: Up to f < n/3 nodes can be Byzantine faulty
- **Partial synchrony**: Network is partially synchronous with unknown GST
- **Bounded delays**: After GST, messages are delivered within time bound Delta

#### Message Echoing

**assume all nodes *implicitly echo* all messages received (wasteful)**

**So if a non-faulty node receives a message at time r,**

**all non-faulty nodes will receive it by time max(GST, r+\Delta)**

Streamlet assumes that all nodes implicitly echo all messages they receive:

- **Implicit echoing**: Every node automatically forwards all messages it receives
- **Wasteful but simple**: This is wasteful but simplifies the protocol design
- **Delivery guarantee**: If a non-faulty node receives a message at time r, all non-faulty nodes will receive it by time max(GST, r+Delta)

#### Epoch Structure

**proceed through epochs of duration 2*\Delta**

**each epoch e has a pseudo-random leader determined by H(e) mod n**

Streamlet operates in epochs:

- **Epoch duration**: Each epoch lasts for 2*Delta time units
- **Pseudo-random leaders**: Each epoch e has a leader determined by H(e) mod n
- **Deterministic selection**: The leader selection is deterministic but appears random

#### Block Structure

**agree on blocks of the form: (h, e, txs)**

**txs are transactions to append to blockchain history**

**e is epoch number of block - not all epochs included in finalized history**

**h is hash of previous block or special genesis block (\bot,0,\bot)**

***length* of block B, |B|, is distance to genesis block (length of chain)**

Streamlet blocks have a specific structure:

- **h**: Hash of the previous block (or special genesis block)
- **e**: Epoch number of the block
- **txs**: Transactions to append to the blockchain history
- **Length**: The length of a block is its distance from the genesis block
- **Not all epochs**: Not all epochs will be included in the finalized history

### How the Streamlet Protocol Works

**How does the protocol work (p. 7)? For each epoch e:**

**A block signed by (2/3)n nodes is called *notarized***

**Propose: Leader broadcasts signed <(h,e,txs)>**

**h must be one of the longest notarized chains leader has seen before**

**Vote: upon receiving the first proposal <(h,e,txs)> from epoch e's leader**

**If h corresponds to one of the longest notarized chains received**

**then vote for block by signing it, broadcasting signature**

**Can't vote for previous epochs or for multiple blocks in same epoch**

**Finalize: when you have 3 notarized blocks with successive epochs:**

**B0=(h0, e, txs0), B1=(H(B0), e+1, txs1), B2=(H(B1), e+2, txs2)**

**Okay to finalize B1--agreement is guaranteed**

The Streamlet protocol operates in three phases for each epoch:

#### Notarization

**A block signed by (2/3)n nodes is called *notarized***

A block becomes notarized when it receives signatures from at least (2/3)n nodes. This provides Byzantine fault tolerance since:

- **Quorum size**: (2/3)n is greater than (1/2)n, ensuring majority agreement
- **Byzantine tolerance**: Even with f < n/3 Byzantine nodes, honest nodes can still achieve notarization
- **Safety**: Notarization provides a strong safety guarantee

#### Propose Phase

**Propose: Leader broadcasts signed <(h,e,txs)>**

**h must be one of the longest notarized chains leader has seen before**

The leader of epoch e:

1. **Selects parent**: Chooses h as the hash of one of the longest notarized chains it has seen
2. **Creates block**: Creates a block (h, e, txs) with transactions txs
3. **Broadcasts**: Broadcasts the signed block to all nodes

The requirement that h must be from one of the longest notarized chains ensures that the leader builds on the most recent consensus state.

#### Vote Phase

**Vote: upon receiving the first proposal <(h,e,txs)> from epoch e's leader**

**If h corresponds to one of the longest notarized chains received**

**then vote for block by signing it, broadcasting signature**

**Can't vote for previous epochs or for multiple blocks in same epoch**

When a node receives a proposal from the epoch's leader:

1. **First proposal**: Only votes for the first proposal received from the epoch's leader
2. **Longest chain**: Only votes if h corresponds to one of the longest notarized chains
3. **Sign and broadcast**: Signs the block and broadcasts the signature
4. **Restrictions**: Cannot vote for previous epochs or multiple blocks in the same epoch

These restrictions ensure that nodes only vote for valid proposals and prevent double-voting.

#### Finalize Phase

**Finalize: when you have 3 notarized blocks with successive epochs:**

**B0=(h0, e, txs0), B1=(H(B0), e+1, txs1), B2=(H(B1), e+2, txs2)**

**Okay to finalize B1--agreement is guaranteed**

A block can be finalized when there are 3 notarized blocks with successive epochs:

- **B0**: Block at epoch e with hash h0
- **B1**: Block at epoch e+1 with hash H(B0) (pointing to B0)
- **B2**: Block at epoch e+2 with hash H(B1) (pointing to B1)

When this condition is met, B1 can be finalized because agreement is guaranteed. This is the key insight of Streamlet - using future blocks to finalize previous blocks.

### Why Not Finalize Every Notarized Block?

**What goes wrong if you just finalize every notarized block?**

**If signature msgs delayed, nodes might not know a block is notarized**

**Hence, might notarize multiple blocks at the same length**

**Essentially you could notarize blocks like this in perpetuity**

**/-- 1 -- 4 -- 7 -- ...   But notice you *have* to interleave results!**

**0--- 2 -- 5 -- 8 -- ...   A node can't vote 5 once it knows 4 notarized**

**\-- 3 -- 6 -- 9 -- ...**

The naive approach of finalizing every notarized block would lead to serious problems:

#### The Signature Delay Problem

**If signature msgs delayed, nodes might not know a block is notarized**

**Hence, might notarize multiple blocks at the same length**

The core issue is that signature messages can be delayed:

- **Delayed signatures**: Nodes might not receive all signatures immediately
- **Unknown notarization**: A node might not know that a block has been notarized
- **Multiple blocks**: This could lead to multiple blocks being notarized at the same length
- **Inconsistent state**: Different nodes might have different views of what's notarized

#### The Perpetual Forking Problem

**Essentially you could notarize blocks like this in perpetuity**

**/-- 1 -- 4 -- 7 -- ...   But notice you *have* to interleave results!**

**0--- 2 -- 5 -- 8 -- ...   A node can't vote 5 once it knows 4 notarized**

**\-- 3 -- 6 -- 9 -- ...**

This could lead to perpetual forking:

- **Multiple chains**: Multiple chains could be notarized simultaneously
- **No convergence**: The system might never converge on a single chain
- **Interleaving requirement**: The key insight is that nodes must interleave results
- **Voting restrictions**: A node can't vote for block 5 once it knows block 4 is notarized

The interleaving requirement is crucial - it prevents nodes from voting for conflicting blocks at the same length.

### Why Not Finalize After 2 Consecutive Epochs?

**What if you finalize first of *2* notarized blocks in consecutive epochs?**

**B0=(h0, e, txs0), B1=(H(B0), e+1, txs1)**

**Could have multiple blocks B0, B0' notarized at length |B0|**

**Nodes might notarize B1 but the signatures are delayed, so don't realize**

**Later, notarize B1'=(H(B0'), e+2, txs1'), B2=(H(B1'), e+3, txs3)**

**So actually B0 was eventually excluded from history!**

**Same example as above, just flip parents around**

**/-- 1 -- 6 -- 7 -- ...**

**0--- 2 -- 5 -- 8 -- ...**

**\-- 3 -- 4 -- 9 -- ...**

Finalizing after just 2 consecutive epochs would also be problematic:

#### The Multiple Block Problem

**Could have multiple blocks B0, B0' notarized at length |B0|**

**Nodes might notarize B1 but the signatures are delayed, so don't realize**

The issue is that multiple blocks could be notarized at the same length:

- **Multiple B0 blocks**: Both B0 and B0' could be notarized at length |B0|
- **Delayed signatures**: Nodes might notarize B1 but not realize it due to delayed signatures
- **Conflicting chains**: This creates conflicting chains that could both be considered valid

#### The Exclusion Problem

**Later, notarize B1'=(H(B0'), e+2, txs1'), B2=(H(B1'), e+3, txs3)**

**So actually B0 was eventually excluded from history!**

**Same example as above, just flip parents around**

**/-- 1 -- 6 -- 7 -- ...**

**0--- 2 -- 5 -- 8 -- ...**

**\-- 3 -- 4 -- 9 -- ...**

The problem is that a block that was initially considered finalized could later be excluded:

- **B1' creation**: A new block B1' could be created pointing to B0' instead of B0
- **B2 creation**: B2 could then be created pointing to B1'
- **B0 exclusion**: This would effectively exclude B0 from the final history
- **Safety violation**: This violates the safety property that finalized blocks should remain finalized

The example shows how the chain could fork and then converge on a different path, excluding previously finalized blocks.

### Why 3 Consecutive Epochs is the Answer

**So why is first 2 of 3 successive notarized epochs the answer?**

**Honest nodes only vote once per epoch, so only one notarized block per epoch**

**Say you notarize B0<--B1<--B2 at epochs e, e+1, e+2 respectively**

**You might notarize multiple blocks of length |B0|**

**But you can't notarize B!=B1 with |B|==|B1|. Why?**

**Proof by contradiction: Say B finalized, let eB be B's epoch**

**eB can't be e, e+1, or e+2 (only one notarized block per epoch)**

**But also can't have eB < e:**

**B's parent is longer than B0's parent by one block**

**So if 2n/3 nodes saw B's parent notarized in epoch eB < e,**

**then 2n/3-f > n/3 barred from voting for B0 in epoch e**

**Also can't have eB > e+2:**

**B's parent is shorter than B2's parent by one block**

**So if 2n/3 nodes saw B2's parent (B1) notarized in epoch e+2**

**then 2n/3-f > n/3 couldn't vote for B in epoch eB > e+2**

**Finalize B1 (+ancestors) when guaranteed only 1 notarized block of len. |B1|**

The key insight is that 3 consecutive epochs provide the necessary safety guarantee:

#### The Uniqueness Guarantee

**Honest nodes only vote once per epoch, so only one notarized block per epoch**

**Say you notarize B0<--B1<--B2 at epochs e, e+1, e+2 respectively**

**You might notarize multiple blocks of length |B0|**

**But you can't notarize B!=B1 with |B|==|B1|. Why?**

The proof shows that B1 is unique at its length:

- **One block per epoch**: Honest nodes only vote once per epoch, so only one notarized block per epoch
- **B0<--B1<--B2**: Three consecutive notarized blocks at epochs e, e+1, e+2
- **Multiple B0 blocks**: There might be multiple blocks of length |B0|
- **Unique B1**: But there cannot be a different block B≠B1 with |B|=|B1|

#### The Proof by Contradiction

**Proof by contradiction: Say B finalized, let eB be B's epoch**

**eB can't be e, e+1, or e+2 (only one notarized block per epoch)**

**But also can't have eB < e:**

**B's parent is longer than B0's parent by one block**

**So if 2n/3 nodes saw B's parent notarized in epoch eB < e,**

**then 2n/3-f > n/3 barred from voting for B0 in epoch e**

**Also can't have eB > e+2:**

**B's parent is shorter than B2's parent by one block**

**So if 2n/3 nodes saw B2's parent (B1) notarized in epoch e+2**

**then 2n/3-f > n/3 couldn't vote for B in epoch eB > e+2**

The proof by contradiction shows that no conflicting block B can exist:

1. **eB can't be e, e+1, or e+2**: Because there's only one notarized block per epoch
2. **eB can't be < e**: Because B's parent would be longer than B0's parent, and if 2n/3 nodes saw B's parent notarized in epoch eB < e, then 2n/3-f > n/3 would be barred from voting for B0 in epoch e
3. **eB can't be > e+2**: Because B's parent would be shorter than B2's parent, and if 2n/3 nodes saw B2's parent (B1) notarized in epoch e+2, then 2n/3-f > n/3 couldn't vote for B in epoch eB > e+2

#### The Finalization Rule

**Finalize B1 (+ancestors) when guaranteed only 1 notarized block of len. |B1|**

This leads to the finalization rule:

- **B1 is unique**: When we have B0<--B1<--B2, B1 is guaranteed to be the only notarized block of length |B1|
- **Safe to finalize**: Since B1 is unique, it's safe to finalize B1 and all its ancestors
- **Safety guarantee**: This provides the safety guarantee that finalized blocks will remain finalized

## Post-GST Synchronization

### Epoch Synchronization After GST

**Now say we reach GST, and all communication happens in Delta...**

**How can we ensure nodes synchronize epochs? Paper doesn't say**

**One idea: Set 3\Delta timer for epoch e only when 2n/3 nodes at e**

**If you ever see f+1 nodes at epoch > yours, immediately switch to new epoch**

**Why does this work after GST (assuming bounded clock drift << \Delta/epoch)?**

**Timer for epoch e won't start unless 2n/3-f honest nodes at e**

**Those 2n/3-f messages will propagate to all honest nodes in time \Delta**

**Since f < 1/3, 2n/3-f >= f+1, all honest stragglers will jump to e**

**That leaves 2\Delta rounds with all nodes on e**

After reaching the Global Stabilization Time (GST), the network becomes synchronous and nodes need to synchronize their epochs:

#### The Synchronization Mechanism

**One idea: Set 3\Delta timer for epoch e only when 2n/3 nodes at e**

**If you ever see f+1 nodes at epoch > yours, immediately switch to new epoch**

The proposed synchronization mechanism works as follows:

- **Conditional timer**: Set a 3Δ timer for epoch e only when 2n/3 nodes are at epoch e
- **Immediate switching**: If you see f+1 nodes at a higher epoch, immediately switch to that epoch
- **Consensus-driven**: The mechanism is driven by consensus rather than time alone

#### Why This Works After GST

**Why does this work after GST (assuming bounded clock drift << \Delta/epoch)?**

**Timer for epoch e won't start unless 2n/3-f honest nodes at e**

**Those 2n/3-f messages will propagate to all honest nodes in time \Delta**

**Since f < 1/3, 2n/3-f >= f+1, all honest stragglers will jump to e**

**That leaves 2\Delta rounds with all nodes on e**

The mechanism works because:

1. **Honest majority requirement**: Timer for epoch e won't start unless 2n/3-f honest nodes are at e
2. **Message propagation**: Those 2n/3-f messages will propagate to all honest nodes in time Δ
3. **Sufficient honest nodes**: Since f < 1/3, we have 2n/3-f ≥ f+1, so all honest stragglers will jump to e
4. **Synchronized operation**: This leaves 2Δ rounds with all nodes synchronized on epoch e

#### Clock Drift Assumption

**Why does this work after GST (assuming bounded clock drift << \Delta/epoch)?**

The mechanism assumes that clock drift is much smaller than Δ/epoch:

- **Bounded drift**: Clock drift must be bounded and small
- **Relative to epoch duration**: The drift must be small relative to the epoch duration
- **Synchronization guarantee**: This ensures that nodes can synchronize their epochs reliably

### Post-GST Behavior with Honest Leaders

**What happens (post-GST) after 2 honest leaders (L0,L1) in epochs e, e+1?**

**Proposals guaranteed to increase in length--why?**

**Say L0 proposes B0 in epoch e**

**Everyone will see B0 proposal by time \Delta into epoch e**

**So either notarize |B0|, or vote blocked by some other B with |B|>=|B0|**

**and if blocked by B, everyone will see B notarized by start of epoch e+1**

**(by implicit echo)**

**So L1 will chose B1 parent as either B0 or some B with |B|>=|B0|**

**So either |B1|-1 = |B0| or |B1| - 1 = |B| >= |B0|**

**Big picture: after honest leader epoch, all nodes have same honest node info**

**Epoch e represents 2\Delta time in which nodes won't sign blocks before e**

**By e+\Delta, everyone has all signatures from honest nodes in prior rounds**

**By e+\Delta, everyone has honest leaders's latest proposal**

**So vote/not vote by e+\Delta, at e+2\Delta (start of e+1), all hear votes**

**Of course, faulty nodes can inject last-minute votes**

After GST with honest leaders, the protocol behavior becomes much more predictable:

#### Length Increase Guarantee

**Proposals guaranteed to increase in length--why?**

**Say L0 proposes B0 in epoch e**

**Everyone will see B0 proposal by time \Delta into epoch e**

**So either notarize |B0|, or vote blocked by some other B with |B|>=|B0|**

**and if blocked by B, everyone will see B notarized by start of epoch e+1**

**(by implicit echo)**

**So L1 will chose B1 parent as either B0 or some B with |B|>=|B0|**

**So either |B1|-1 = |B0| or |B1| - 1 = |B| >= |B0|**

The key insight is that proposals are guaranteed to increase in length:

1. **L0 proposes B0**: In epoch e, leader L0 proposes block B0
2. **Everyone sees B0**: By time Δ into epoch e, everyone sees B0's proposal
3. **Two outcomes**: Either B0 gets notarized, or voting is blocked by some other block B with |B|≥|B0|
4. **Block visibility**: If blocked by B, everyone sees B notarized by start of epoch e+1 (by implicit echo)
5. **L1's choice**: Leader L1 will choose B1's parent as either B0 or some B with |B|≥|B0|
6. **Length guarantee**: So either |B1|-1 = |B0| or |B1|-1 = |B| ≥ |B0|

This guarantees that block lengths increase monotonically.

#### Synchronized Information

**Big picture: after honest leader epoch, all nodes have same honest node info**

**Epoch e represents 2\Delta time in which nodes won't sign blocks before e**

**By e+\Delta, everyone has all signatures from honest nodes in prior rounds**

**By e+\Delta, everyone has honest leaders's latest proposal**

**So vote/not vote by e+\Delta, at e+2\Delta (start of e+1), all hear votes**

**Of course, faulty nodes can inject last-minute votes**

After an honest leader epoch, all nodes have synchronized information:

- **Same information**: All nodes have the same information from honest nodes
- **Epoch duration**: Epoch e represents 2Δ time in which nodes won't sign blocks before e
- **Signature synchronization**: By e+Δ, everyone has all signatures from honest nodes in prior rounds
- **Proposal synchronization**: By e+Δ, everyone has the honest leader's latest proposal
- **Voting synchronization**: So vote/not vote by e+Δ, at e+2Δ (start of e+1), all hear votes
- **Faulty node behavior**: Of course, faulty nodes can inject last-minute votes

This synchronization is crucial for the protocol's safety and liveness properties.

### Three Honest Leaders: Uniqueness Guarantee

**What happens following 3 honest leaders at e, e+1, e+2 after GST?**

**Leader of e+2 will propose some block B**

**B will be notarized and all honest nodes will see it at start of e+3**

**No B' != B with |B'| == |B| will ever be notarized**

**Why?**

**Let l0, l1, l2 be height of blocks proposed at e, e+1, e+2**

**Since assuming honest leaders and after GST, l0 < l1 < l2**

**No honest node could ever vote for conflicting B' with |B'| == l2**

**Can't vote for it in e or e+1, since not proposed by honest leader**

**But if voted before e, saw notarized B'' with |B''| == l2-1**

**So all nodes would see B'' notarized by e+1 and reject l1 proposal**

**So all honest *can* vote for B**

**All *will* since by e+2+\Delta will see B and parent notarized chain**

With three consecutive honest leaders, we get a strong uniqueness guarantee:

#### The Uniqueness Proof

**Let l0, l1, l2 be height of blocks proposed at e, e+1, e+2**

**Since assuming honest leaders and after GST, l0 < l1 < l2**

**No honest node could ever vote for conflicting B' with |B'| == l2**

**Can't vote for it in e or e+1, since not proposed by honest leader**

**But if voted before e, saw notarized B'' with |B''| == l2-1**

**So all nodes would see B'' notarized by e+1 and reject l1 proposal**

**So all honest *can* vote for B**

**All *will* since by e+2+\Delta will see B and parent notarized chain**

The proof shows that no conflicting block can be notarized:

1. **Monotonic heights**: Since we have honest leaders and are after GST, l0 < l1 < l2
2. **No honest votes for B'**: No honest node could ever vote for conflicting B' with |B'| = l2
3. **Can't vote in e or e+1**: Can't vote for B' in epochs e or e+1, since it wasn't proposed by the honest leader
4. **Can't vote before e**: If voted before e, the node saw notarized B'' with |B''| = l2-1
5. **B'' visibility**: So all nodes would see B'' notarized by e+1 and reject the l1 proposal
6. **All honest can vote**: So all honest nodes can vote for B
7. **All honest will vote**: All honest nodes will vote since by e+2+Δ they will see B and its parent notarized chain

#### The Guarantee

**Leader of e+2 will propose some block B**

**B will be notarized and all honest nodes will see it at start of e+3**

**No B' != B with |B'| == |B| will ever be notarized**

This provides a strong guarantee:

- **B will be notarized**: The block B proposed by the leader of e+2 will be notarized
- **All nodes see it**: All honest nodes will see B at the start of e+3
- **Uniqueness**: No different block B' with |B'| = |B| will ever be notarized
- **Safety**: This provides the safety guarantee that finalized blocks remain unique

### Liveness: When New Blocks Are Guaranteed

**What condition guarantees new block? 5 successive non-faulty leaders after GST**

**First three leaders produce a notarized |B3| unique for its height**

**Next two leaders produce notarized B3<-B4<-B5 w. successive epoch numbers**

**0--- 1 -- 3 -- 5 -- 6 -- 7**

**\-- 2 -- 4**

For liveness, we need to understand when new blocks are guaranteed to be finalized:

#### The 5-Leader Condition

**What condition guarantees new block? 5 successive non-faulty leaders after GST**

**First three leaders produce a notarized |B3| unique for its height**

**Next two leaders produce notarized B3<-B4<-B5 w. successive epoch numbers**

The condition for guaranteed new block finalization is 5 successive non-faulty leaders after GST:

1. **First three leaders**: Produce a notarized block B3 that is unique for its height
2. **Next two leaders**: Produce notarized blocks B3←B4←B5 with successive epoch numbers
3. **Finalization**: This allows B3 to be finalized (as the first of 3 consecutive blocks)

#### The Chain Structure

**0--- 1 -- 3 -- 5 -- 6 -- 7**

**\-- 2 -- 4**

The chain structure shows:

- **Main chain**: 0→1→3→5→6→7
- **Fork**: 0→2→4
- **Finalization**: Block 3 can be finalized once we have blocks 3, 4, 5 in consecutive epochs

This demonstrates how the protocol ensures both safety and liveness through the careful sequencing of honest leaders.

## Synchronous Variant

### Synchronous Streamlet for f < n/2

**How does synchronous variant work for n nodes, f < n/2?**

**Notarization requires n/2 nodes, not 2n/3**

**What's to stop two partitioned of size n/2 from disagreeing?**

**Synchrony assumption rules out network partition**

**(Recall Byz generals--sync. model doesn't limit # of faults tolerated)**

**To finalize block at epoch e:**

**Need chain of *6* notarized blocks in consecutive epochs (e,...,e+5)**

**AND Must see no conflicting notarized blocks at the same lengths**

The synchronous variant of Streamlet works differently:

#### Different Fault Tolerance

**Notarization requires n/2 nodes, not 2n/3**

**What's to stop two partitioned of size n/2 from disagreeing?**

**Synchrony assumption rules out network partition**

**(Recall Byz generals--sync. model doesn't limit # of faults tolerated)**

In the synchronous model:

- **Lower threshold**: Notarization requires n/2 nodes instead of 2n/3
- **No network partitions**: The synchrony assumption rules out network partitions
- **Higher fault tolerance**: The synchronous model doesn't limit the number of faults tolerated (unlike the Byzantine generals problem)
- **No partition problem**: Two partitions of size n/2 can't disagree because partitions are impossible

#### Different Finalization Rule

**To finalize block at epoch e:**

**Need chain of *6* notarized blocks in consecutive epochs (e,...,e+5)**

**AND Must see no conflicting notarized blocks at the same lengths**

The synchronous variant requires:

- **6 consecutive blocks**: Need a chain of 6 notarized blocks in consecutive epochs (e,...,e+5)
- **No conflicts**: Must see no conflicting notarized blocks at the same lengths
- **Stronger requirement**: This is a stronger requirement than the 3-block rule in the partially synchronous version

### Multiple Blocks in Same Epoch

**Can you notarize two blocks in the same epoch? Yes**

**But say non-faulty node v sees B is notarized in epoch e ("Fact 4")**

**All nodes will see any messages v saw (and agree B notarized) by e+2**

**And say B0 <-- B1 notarized at e, e+1 respectively ("Lemma 7")**

**Means some honest node saw B0 notarized in e+1**

**So no B with |B|<=|B0| can be notarized in e+3 or later**

**Suppose two honest nodes finalize chain (ending B0...B5) and chain'**

**Must be that B0 appears in history of B0'. Proof by contradiction**

**Assume w.l.o.g. chain no longer than chain'**

**Let B0,...,B5 have length l-5,...,l and epoch e-5,...,e**

**Let B0' = chain'[l-5], B1' = chain'[l-4]; let e' be epoch of B1'**

**By assumption, B0' != B0**

**By lemma 7, epoch of B1' must satisfy e' < (e-4)+3, or e' <= e-2**

**By lemma 4, by e'+2, all nodes will see B1' notarized**

**But e'+2 <= e, meaning by e, all nodes see B1', preventing finalization**

The synchronous variant allows multiple blocks in the same epoch but with strong guarantees:

#### Fact 4: Message Propagation

**But say non-faulty node v sees B is notarized in epoch e ("Fact 4")**

**All nodes will see any messages v saw (and agree B notarized) by e+2**

- **Message propagation**: If a non-faulty node sees B notarized in epoch e, all nodes will see this by e+2
- **Consensus propagation**: This ensures that notarization information propagates quickly

#### Lemma 7: Block Ordering

**And say B0 <-- B1 notarized at e, e+1 respectively ("Lemma 7")**

**Means some honest node saw B0 notarized in e+1**

**So no B with |B|<=|B0| can be notarized in e+3 or later**

- **Block ordering**: If B0←B1 are notarized at epochs e, e+1, then some honest node saw B0 notarized in e+1
- **Future restriction**: No block B with |B|≤|B0| can be notarized in e+3 or later
- **Consistency guarantee**: This provides consistency guarantees across epochs

#### Safety Proof

**Suppose two honest nodes finalize chain (ending B0...B5) and chain'**

**Must be that B0 appears in history of B0'. Proof by contradiction**

The safety proof shows that if two honest nodes finalize different chains, they must share common blocks:

- **Contradiction assumption**: Assume two honest nodes finalize different chains
- **Common block requirement**: B0 must appear in the history of B0'
- **Proof by contradiction**: This leads to a contradiction, proving safety

## Implementation Considerations and Optimizations

### Minor Issues and Improvements

**Nits to pick:**

**Notarization should wait for >2n/3 rather than >=2n/3**

**Remember quorum size T = f_S + f_L + 1**

**So if N = 3k and T = 2k, best to choose f_S = k, f_L = k - 1**

**Instead, authors chose f_L > f_s, which is not useful**

**With multiple longest notarized chains, have leader pick largest e?**

**More likely to be prior epoch, maybe lead to faster finalization**

There are several minor issues and potential improvements:

#### Quorum Size Optimization

**Notarization should wait for >2n/3 rather than >=2n/3**

**Remember quorum size T = f_S + f_L + 1**

**So if N = 3k and T = 2k, best to choose f_S = k, f_L = k - 1**

**Instead, authors chose f_L > f_s, which is not useful**

The quorum size could be optimized:

- **Strict inequality**: Notarization should wait for >2n/3 rather than ≥2n/3
- **Quorum formula**: Remember that quorum size T = f_S + f_L + 1
- **Optimal choice**: If N = 3k and T = 2k, best to choose f_S = k, f_L = k-1
- **Suboptimal choice**: Instead, authors chose f_L > f_S, which is not useful

#### Chain Selection Optimization

**With multiple longest notarized chains, have leader pick largest e?**

**More likely to be prior epoch, maybe lead to faster finalization**

When there are multiple longest notarized chains:

- **Epoch selection**: Have the leader pick the chain with the largest epoch number
- **Prior epoch preference**: This is more likely to be from a prior epoch
- **Faster finalization**: This might lead to faster finalization

### Eliminating Implicit Echo

**How might you get rid of implicit echo?**

**Like HotStuff, coordinate all messages through the leader**

**Leader broadcast proposal, all send votes, leader broadcasts votes**

**Threshold cryptography can make this even better**

**Instead of sending 2n/3 signatures, combine into single threshold signature**

The implicit echo assumption can be eliminated through better message coordination:

#### Leader-Coordinated Communication

**Like HotStuff, coordinate all messages through the leader**

**Leader broadcast proposal, all send votes, leader broadcasts votes**

Similar to HotStuff, messages can be coordinated through the leader:

- **Leader broadcasts proposal**: The leader broadcasts the proposal to all nodes
- **Nodes send votes**: All nodes send their votes to the leader
- **Leader broadcasts votes**: The leader broadcasts all votes to all nodes
- **Eliminates echo**: This eliminates the need for implicit echo

#### Threshold Cryptography

**Threshold cryptography can make this even better**

**Instead of sending 2n/3 signatures, combine into single threshold signature**

Threshold cryptography can further improve efficiency:

- **Single signature**: Instead of sending 2n/3 individual signatures, combine them into a single threshold signature
- **Reduced communication**: This significantly reduces communication overhead
- **Same security**: Provides the same security guarantees as individual signatures
- **Better efficiency**: Much more efficient than collecting and broadcasting individual signatures

## Conclusion

Streamlet represents a significant advancement in blockchain consensus protocols by:

1. **Simplicity**: Providing a much simpler alternative to complex protocols like PBFT
2. **Blockchain optimization**: Leveraging the natural structure of blockchains to simplify consensus
3. **Safety and liveness**: Providing strong safety and liveness guarantees
4. **Practical applicability**: Being suitable for both permissioned and permissionless blockchains
5. **Optimization potential**: Having clear paths for further optimization and improvement

The protocol demonstrates that blockchain consensus can be both simple and secure, opening new possibilities for practical blockchain applications.