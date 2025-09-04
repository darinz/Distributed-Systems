# Cryptocurrency: Supplementary Notes

## Decentralized Control

### The Problem
- **PBFT and similar protocols require**: Public-key infrastructure and that servers know who other servers are
- **This must be setup**: By some central authority for the protocol to run
- **Otherwise, these protocols are susceptible**: To Sybil attacks
- **What if you want**: A decentralized system?

### Centralized vs. Decentralized
- **Centralized systems**: Require trusted authorities
- **Sybil attacks**: Attackers can create multiple identities
- **Decentralized systems**: No central authority needed
- **Challenge**: How to achieve consensus without central authority?

## Two Classes of Solutions

### Proof of Work
- **Rate of transaction commitment**: Limited by cryptographically hard problem
- **Nodes called miners**: Solve these problems to commit transactions
- **Assumes that majority of CPU power**: Is controlled by honest* nodes
- **Miners are rewarded**: With transaction fees and mining rewards

### Proof of Stake
- **Transactions are committed**: With votes weighted by amount of stake voters have in system
- **Assume that 2/3rds of money**: Is controlled by honest* nodes
- **Voters sometimes rewarded**: For taking part in protocol (but they also have stake in system)

### Key Differences
- **Proof of Work**: CPU power determines influence
- **Proof of Stake**: Money/stake determines influence
- **Both assume**: Honest majority controls the resource

## Bitcoin

### Overview
- **Bitcoin is a proof-of-work cryptocurrency network**: Started in 2009
- **Goal**: Electronic money without the need for trust
- **Relies on cryptography**: For authentication, proof-of-work for transaction ordering
- **Decentralized**: No central authority

### Key Properties
- **Electronic money**: Digital currency
- **No trust required**: Cryptographic guarantees
- **Decentralized**: No central authority
- **Public ledger**: All transactions visible

## Bitcoin Transactions

### Structure
- **Payment is signed, publicly visible transaction**: Between public/private key pairs
- **Transactions have**: (Potentially multiple) inputs and outputs
- **Transaction inputs**: Are other transactions
- **Transaction outputs**: Are public keys (recipients)

### Transaction Flow
- **Inputs**: Reference to previous transactions
- **Outputs**: New recipients and amounts
- **Signatures**: Prove ownership of inputs
- **Public visibility**: All transactions recorded

## Strawman Proposal

### Simple Approach
- **Lukas just signs transaction**: And gives it to Arvind
- **What could go wrong?**
  - **Arvind couldn't impersonate Lukas**: He doesn't have Lukas's private key
  - **What if sender already spent transaction?**: This is called double-spending
  - **Where does money actually come from?**

### Problems
- **Double-spending**: Same money spent twice
- **No money creation**: Where does initial money come from?
- **No ordering**: Which transaction came first?
- **No consensus**: Who decides what's valid?

## Trusted Third Parties (Not a Strawman)

### Centralized Solution
- **Sender could send transaction**: To trusted third party (or system)
- **As long as transaction is valid**: (i.e., input transactions weren't already spent), accepts transaction and puts it in log
- **Log is made publicly visible**: And can be replicated by any number of passive listeners
- **Recipients wait**: Until they see transaction in log

### Benefits
- **Prevents double-spending**: Central authority tracks spent transactions
- **Public log**: Transparent and auditable
- **Replication**: Multiple copies for reliability
- **Commitment**: Once in log, transaction is committed

## Managing the Public Log

### The Challenge
- **We need log to stay consistent**: (i.e., transactions stay in same order in log)
- **We could use Paxos**: But what if replicas aren't trusted?
- **PBFT still requires trusting**: 2f + 1 replicas
- **Need decentralized solution**: No trusted replicas

### Consensus Requirements
- **Consistent ordering**: All nodes see same order
- **No trusted parties**: Decentralized consensus
- **Fault tolerance**: Handle node failures
- **Sybil resistance**: Prevent identity attacks

## Bitcoin Mining

### Mining Process
- **Bitcoin commits transactions**: By having servers called miners solve cryptographic puzzle
- **Transactions are committed in blocks**: Groups of transactions
- **Miners try to find nonce**: Such that hash of entire block is less than some threshold
- **Finding such nonce is difficult**: But miners get compensated

### Compensation
- **Mining rewards**: Bitcoin from nowhere
- **Transaction fees**: Bitcoin from transaction senders
- **Incentive alignment**: Miners profit from securing network

### Block Structure
- **Transactions**: Grouped into blocks
- **Nonce**: Random number to solve puzzle
- **Hash**: Must be below threshold
- **Previous block hash**: Links to previous block

## Finding a Stable Order

### Blockchain Structure
- **Each block has single pointer**: To previous block (except initial block)
- **These blocks form DAG**: Directed Acyclic Graph
- **Honest miners work off longest chain**: If they see two chains of equal length, work off one they saw first
- **Only transactions with unspent inputs**: Are valid

### Confirmation Process
- **Clients wait for transactions**: To be 6 blocks deep
- **6 blocks deep**: In chain 6 blocks longer than any chain without transaction
- **Before considering it confirmed**: Transaction is final
- **As long as honest miners control >50%**: Of hashing power, longest chain can't be overrun

### Security Guarantees
- **Confirmed transactions won't be undone**: By double-spend
- **Longest chain rule**: Prevents chain reorganization
- **Majority honest assumption**: 50%+ hashing power
- **Economic incentives**: Miners profit from honest behavior

## Network Protocol

### Communication
- **Bitcoin uses gossip protocol**: To communicate new blocks and transaction requests
- **Each peer connected**: To set of other peers
- **Peer list bootstrapped**: Usually using DNS by asking hostname that points to known nodes
- **Decentralized network**: No central servers

### Peer-to-Peer
- **Gossip protocol**: Spreads information efficiently
- **Peer connections**: Each node connects to multiple peers
- **DNS bootstrapping**: Initial peer discovery
- **Network resilience**: No single point of failure

## Hash Puzzle Difficulty

### Difficulty Adjustment
- **Threshold for mining puzzle**: Set by difficulty, a 256-bit number
- **If difficulty is 2^254**: There's ½ chance for any given nonce
- **2^253 gives ¼ chance**: Etc.
- **Difficulty adjusted every 2016 blocks**: To keep average throughput at ~1 block/10 mins
- **Average time to confirm transaction**: 1 hour

### Difficulty Mechanics
- **256-bit difficulty**: Very large number
- **Exponential relationship**: Difficulty vs. success probability
- **Automatic adjustment**: Every 2016 blocks
- **Target block time**: 10 minutes average

## Mining Rewards

### Reward Structure
- **Every time block is "mined"**: Miner gets reward
- **Reward started at 50 BTC**: Halved every 210,000 blocks (approximately every 4 years)
- **Since bitcoins aren't infinitely divisible**: Reward will go to 0 at some point
- **Maximum of 21M BTC**: Will ever exist
- **Currently, about 85%**: Of all bitcoins have been mined

### Halving Schedule
- **Initial reward**: 50 BTC per block
- **Halved every 210,000 blocks**: ~4 years
- **Current reward**: 6.25 BTC per block
- **Eventually goes to 0**: When all bitcoins mined

## Transaction Fees

### Fee Structure
- **Transaction senders pay fee**: Claimed by winning miner
- **Higher fee**: More incentivized miners to commit that transaction
- **Once all bitcoins mined**: This will be only mining incentive
- **Currently averaging**: About 0.00050 BTC (= $4 at current prices)

### Fee Economics
- **Fee market**: Higher fees = faster confirmation
- **Miner incentive**: Fees + block reward
- **Future incentive**: Only fees after all bitcoins mined
- **Market-driven**: Fees determined by demand

## Hashing – Like Machine Learning But Less Useful

### Energy Consumption
- **Even with specialized hardware**: Hashing is energy-intensive
- **Currently, overall hashrate**: 55 EH (exahash)/s
- **Bitcoin mining network consumes**: Same amount of energy as Switzerland(!)
- **Environmental concerns**: Massive energy usage

### Hashrate Scale
- **55 exahash per second**: Extremely high computational power
- **Specialized hardware**: ASICs for mining
- **Energy intensive**: Significant environmental impact
- **Switzerland comparison**: Puts energy usage in perspective

## Bitcoin Throughput

### Performance Metrics
- **Currently, average of 2,500 transactions**: In 2MB bitcoin block
- **Network mines block**: Once every 10 minutes on average
- **This gives us ~4 transactions/s**: Very low throughput
- **Scalability challenge**: Limited transaction capacity

### Throughput Limitations
- **2,500 transactions per block**: Limited by block size
- **10-minute block time**: Limited by security requirements
- **4 transactions per second**: Much lower than traditional payment systems
- **Scalability bottleneck**: Block size and block time

## What Did This Get Us?

### Promised Benefits
- **Privacy?**: Well, not really. Your name isn't published, but flow of money from one transaction to another is public
- **Non-repudiation?**: Why couldn't a bank guarantee this?
- **No trusted authority?**: Great, now drug dealers and human traffickers get financial infrastructure, too!
- **No centralized monetary policy?**: You like deflation?

### Reality Check
- **Limited privacy**: Transaction graph is public
- **Banks can provide**: Non-repudiation guarantees
- **Criminal use**: Enables illegal activities
- **Deflationary**: Fixed money supply

### Critical Questions
- **Does this look like a currency?**: Limited utility as currency
- **Why are people putting money in this?**: Speculation vs. utility
- **Energy consumption worth it?**: Environmental cost
- **Decentralization valuable?**: Trade-offs with efficiency

## Other Proof-of-Work Systems

### Bitcoin Alternatives
- **Bitcoin is by no means**: The only popular proof-of-work based system
- **Zerocoin**: Provides better anonymity (which makes it even better for money laundering?)
- **Ethereum**: Allows scripting
- **Ripple**: Tries to maintain stable price
- **...and many others**: Hundreds of cryptocurrencies

### Variations
- **Better anonymity**: Zerocoin, Monero
- **Smart contracts**: Ethereum
- **Stable value**: Ripple, Tether
- **Different algorithms**: Various proof-of-work schemes

## Bitcoin Discussion Questions

### Fundamental Questions
- **Where does value of Bitcoin come from?**: Speculation vs. utility
- **Is energy consumption worth it?**: Environmental impact
- **How valuable is decentralization?**: Trade-offs with efficiency
- **Is Bitcoin useful as currency?**: For small transactions?

### Technical Questions
- **How long will SHA-256 last?**: Cryptographic security
- **How do we make changes to protocol?**: Governance challenges
- **Is Bitcoin actually anonymous?**: Privacy limitations
- **What if miners are rational (greedy) instead of honest?**: Game theory

### Ethical Questions
- **Is Bitcoin ethical**: Given benefits for ransomware, money laundering, etc.?
- **Why do wallets and private exchanges exist?**: Don't they defeat the purpose?
- **What implications does non-reversibility have?**: Irreversible transactions

## Proof-of-Stake – Algorand

### Overview
- **Created in 2017**: Uses proof-of-stake instead of proof-of-work
- **Not the first**: But significant improvement
- **One of approx. 300 billion**: Blockchain startups
- **Committee-based consensus**: More efficient than proof-of-work

### Key Innovation
- **Proof-of-stake**: Instead of proof-of-work
- **Committee-based**: Random selection of validators
- **Cryptographic sortition**: Verifiable random functions
- **More efficient**: Less energy consumption

## Main Ideas

### Stake-Based Consensus
- **Weight users**: By how much money they hold in account
- **Use Byzantine agreement**: But over randomly selected committee
- **Choose committees**: Based on cryptographic sortition
- **Verifiable random functions**: On publicly available data and secret information

### Committee Selection
- **Cryptographic sortition**: Uses VRF on publicly available data and secret information
- **Adversary can't target**: Committee members ahead of time
- **Each committee used**: For single step only
- **As soon as committee member reveals decision**: They're no longer relevant and can't be targeted

### Security Properties
- **Random selection**: Prevents targeting
- **Single-use committees**: Limit exposure
- **Cryptographic proofs**: Verifiable randomness
- **Byzantine agreement**: On selected committees

## Goal

### Requirements
- **Transaction structure same as Bitcoin**: Just need to agree on ordering of transactions (blocks)
- **Want safety (linearizability)**: With high probability
- **Assume at least 2/3 of money**: Is held by honest users running bug-free code
- **System should be reasonably performant**: And scalable

### Design Goals
- **Same transaction model**: As Bitcoin
- **Consensus on ordering**: Of transactions
- **High probability safety**: Linearizability
- **Honest majority**: 2/3 of stake
- **Performance**: Reasonable throughput and latency

## Verifiable Random Functions

### Cryptographic Sortition
- **Bottom line**: Output of VRF determines whether (and how many times) user is chosen for particular role
- **Verifiable randomness**: Can't be predicted or manipulated
- **Public verification**: Anyone can verify selection
- **Secret information**: Prevents targeting

### VRF Properties
- **Unpredictable**: Can't predict output
- **Verifiable**: Can verify correctness
- **Fair**: All participants have equal chance
- **Secure**: Cryptographically sound

## Round Structure

### Consensus Process
- **Each round**: Opportunity to commit some block
- **First, proposers chosen**: And propose blocks
- **Next, system runs BA***: Their main agreement protocol to choose block from among proposed ones
- **BA* proceeds in two phases**: First reduces problem to choosing between two options, then runs BinaryBA* to choose between those
- **BA* can reach TENTATIVE or FINAL agreement**: Block committed if FINAL or if one of block's successors is FINAL

### Two-Phase Protocol
- **Phase 1**: Reduce to binary choice
- **Phase 2**: Binary Byzantine agreement
- **TENTATIVE vs. FINAL**: Different levels of commitment
- **Successor commitment**: Block committed if successor is FINAL

## Algorand Takeaways

### Advantages
- **Algorand doesn't utilize proof-of-work**: Instead weights users based on money they have in system
- **More communication efficient**: Since it is committee-based
- **Lower energy consumption**: No mining required
- **Faster consensus**: Committee-based approach

### Challenges
- **Not clear what incentives users have**: To participate in protocol (their stake in system notwithstanding)
- **Algorand requires money holders**: To be online and broadcasting their address to world
- **Algorand is really complicated**: High complexity
- **Stake concentration**: Rich get richer problem

### Trade-offs
- **Efficiency vs. complexity**: More efficient but more complex
- **Energy vs. stake**: Lower energy but requires stake
- **Decentralization vs. performance**: Better performance but less decentralized
- **Security vs. usability**: More secure but harder to use

## Key Takeaways

### Cryptocurrency Design Principles
- **Decentralized consensus**: Without central authority
- **Proof-of-work vs. proof-of-stake**: Different resource requirements
- **Economic incentives**: Align participant behavior
- **Cryptographic security**: Ensure system integrity
- **Public ledger**: Transparent transaction history

### Bitcoin Characteristics
- **Proof-of-work**: CPU power determines influence
- **Blockchain**: Chain of blocks with transactions
- **Mining**: Solving cryptographic puzzles
- **Limited throughput**: ~4 transactions per second
- **High energy consumption**: Environmental concerns

### Algorand Innovations
- **Proof-of-stake**: Money determines influence
- **Committee-based consensus**: Random selection of validators
- **Verifiable random functions**: Cryptographic sortition
- **More efficient**: Lower energy consumption
- **Higher complexity**: More sophisticated protocol

### Consensus Mechanisms
- **Proof-of-work**: Secure but energy-intensive
- **Proof-of-stake**: Efficient but complex
- **Committee-based**: Balance of efficiency and security
- **Byzantine agreement**: Handle malicious participants
- **Economic incentives**: Align participant behavior

### Trade-offs
- **Decentralization vs. efficiency**: More decentralized = less efficient
- **Security vs. performance**: More secure = slower
- **Energy vs. stake**: Proof-of-work vs. proof-of-stake
- **Complexity vs. simplicity**: More features = more complex
- **Privacy vs. transparency**: More private = less transparent

### Modern Relevance
- **Blockchain technology**: Beyond just cryptocurrency
- **Smart contracts**: Programmable money
- **DeFi**: Decentralized finance
- **NFTs**: Non-fungible tokens
- **Web3**: Decentralized web applications

### Lessons Learned
- **Decentralization is hard**: Trade-offs with efficiency
- **Economic incentives matter**: Align participant behavior
- **Cryptography is essential**: For security and consensus
- **Energy consumption**: Significant environmental impact
- **Scalability challenges**: Limited transaction throughput
- **Governance complexity**: How to make protocol changes
