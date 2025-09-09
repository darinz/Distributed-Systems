# Distributed Systems Research Papers

This directory contains foundational and influential research papers in distributed systems, organized chronologically and by topic. These papers form the theoretical and practical foundation for understanding modern distributed systems.

## Table of Contents

- [Reading Guide](#reading-guide)
- [Foundational Papers](#foundational-papers)
- [Consensus and Coordination](#consensus-and-coordination)
- [Distributed Storage Systems](#distributed-storage-systems)
- [Fault Tolerance and Byzantine Systems](#fault-tolerance-and-byzantine-systems)
- [Cryptocurrency and Blockchain](#cryptocurrency-and-blockchain)
- [Performance and Optimization](#performance-and-optimization)
- [Reading Order Recommendations](#reading-order-recommendations)

## Reading Guide

### How to Read a Paper
- **01_how-to-read-a-paper.pdf** - Essential guide for effective paper reading techniques
- **02_Intro-to-Distributed-System_Google.pdf** - Google's comprehensive introduction to distributed systems concepts

## Foundational Papers

### Knowledge and Reasoning in Distributed Systems
- **03_knowledge-in-distributed-enviornment.pdf** - Fundamental work on knowledge and reasoning in distributed environments
- **04_clock-lamport.pdf** - Lamport's seminal work on logical clocks and causality in distributed systems

### Fault Tolerance and Virtual Machines
- **05_design-fault-tolerant-vm.pdf** - Design principles for fault-tolerant virtual machines
- **06_consistent-global-states_chapt4.pdf** - Theory and algorithms for consistent global states in distributed systems

### Interprocess Communication
- **07_interprocess.pdf** - Interprocess communication mechanisms and protocols

## Consensus and Coordination

### Wait-Free Synchronization
- **08_p463-herlihy.pdf** - Herlihy's work on wait-free synchronization and consensus protocols

### Safety and Liveness Properties
- **09_RecSafeLive.pdf** - Formal treatment of safety and liveness properties in distributed systems

### Paxos Consensus Protocol
- **10_paxos-simple.pdf** - Lamport's "Paxos Made Simple" - the definitive guide to the Paxos consensus algorithm
- **12_paxos-made-complex.pdf** - More detailed treatment of Paxos with additional complexity and variations

### Two-Phase Commit
- **14_two-phase-commit_chapter7.pdf** - Comprehensive coverage of two-phase commit protocol

## Distributed Storage Systems

### Shared Memory and Registers
- **13_sharing-mem-abd.pdf** - Attiya, Bar-Noy, Dolev (ABD) algorithm for implementing shared memory in distributed systems

### Google's Distributed Systems
- **15_gfs-sosp2003.pdf** - Google File System (GFS) - large-scale distributed file system
- **16_bigtable-osdi06.pdf** - BigTable - distributed storage system for structured data
- **17_spanner-osdi2012.pdf** - Spanner - globally distributed database with external consistency

### Amazon's Dynamo
- **22_dynamo.pdf** - Dynamo - highly available key-value storage system

### Caching Systems
- **23_memcache_nsdi13-final170_update.pdf** - Facebook's memcache system for distributed caching

## Fault Tolerance and Byzantine Systems

### Byzantine Fault Tolerance
- **18_byzantine_osdi99.pdf** - Practical Byzantine Fault Tolerance (PBFT) algorithm

### Eventual Consistency
- **21_bayou.pdf** - Bayou system for managing replicated, weakly consistent data

## Cryptocurrency and Blockchain

### Bitcoin
- **19_bitcoin.pdf** - Satoshi Nakamoto's original Bitcoin paper

### Algorand
- **20_gilad-algorand-eprint.pdf** - Algorand consensus protocol for blockchain systems

## Performance and Optimization

### Free Choice in Distributed Systems
- **11_advantage-of-free-choice.pdf** - Advantages of free choice in distributed systems

## Reading Order Recommendations

### Beginner Path (Start Here)
1. **01_how-to-read-a-paper.pdf** - Learn how to read research papers effectively
2. **02_Intro-to-Distributed-System_Google.pdf** - Get comprehensive overview of distributed systems
3. **04_clock-lamport.pdf** - Understand logical clocks and causality
4. **10_paxos-simple.pdf** - Learn the basics of consensus with Paxos
5. **15_gfs-sosp2003.pdf** - Study a real distributed file system

### Intermediate Path
1. **13_sharing-mem-abd.pdf** - Understand shared memory in distributed systems
2. **16_bigtable-osdi06.pdf** - Learn about distributed storage systems
3. **22_dynamo.pdf** - Study eventual consistency and high availability
4. **18_byzantine_osdi99.pdf** - Understand Byzantine fault tolerance
5. **17_spanner-osdi2012.pdf** - Learn about global consistency

### Advanced Path
1. **12_paxos-made-complex.pdf** - Deep dive into Paxos variations
2. **19_bitcoin.pdf** - Understand blockchain and cryptocurrency
3. **20_gilad-algorand-eprint.pdf** - Study modern consensus protocols
4. **23_memcache_nsdi13-final170_update.pdf** - Learn about large-scale caching
5. **21_bayou.pdf** - Understand eventual consistency systems

### Specialized Topics

#### Consensus and Coordination
- 08_p463-herlihy.pdf
- 10_paxos-simple.pdf
- 12_paxos-made-complex.pdf
- 14_two-phase-commit_chapter7.pdf
- 18_byzantine_osdi99.pdf
- 20_gilad-algorand-eprint.pdf

#### Distributed Storage
- 13_sharing-mem-abd.pdf
- 15_gfs-sosp2003.pdf
- 16_bigtable-osdi06.pdf
- 17_spanner-osdi2012.pdf
- 22_dynamo.pdf
- 23_memcache_nsdi13-final170_update.pdf

#### Fault Tolerance
- 05_design-fault-tolerant-vm.pdf
- 06_consistent-global-states_chapt4.pdf
- 18_byzantine_osdi99.pdf
- 21_bayou.pdf

#### Blockchain and Cryptocurrency
- 19_bitcoin.pdf
- 20_gilad-algorand-eprint.pdf

## Key Concepts Covered

### Fundamental Concepts
- **Logical Clocks**: Lamport clocks, vector clocks, causality
- **Consensus**: Paxos, Raft, Byzantine fault tolerance
- **Consistency Models**: Strong consistency, eventual consistency, causal consistency
- **Fault Tolerance**: Crash failures, Byzantine failures, recovery

### System Design
- **Distributed File Systems**: GFS, HDFS
- **Distributed Databases**: BigTable, Spanner, Dynamo
- **Caching Systems**: Memcache, CDNs
- **Blockchain Systems**: Bitcoin, Algorand

### Performance and Scalability
- **Partitioning**: Horizontal and vertical scaling
- **Replication**: Primary-backup, chain replication
- **Load Balancing**: Consistent hashing, sharding
- **Optimization**: Caching, compression, indexing

## Paper Categories by Impact

### Highly Influential (Must Read)
- 04_clock-lamport.pdf - Logical clocks
- 10_paxos-simple.pdf - Paxos consensus
- 15_gfs-sosp2003.pdf - Google File System
- 19_bitcoin.pdf - Bitcoin blockchain

### Industry Defining
- 16_bigtable-osdi06.pdf - BigTable
- 17_spanner-osdi2012.pdf - Spanner
- 22_dynamo.pdf - Dynamo
- 23_memcache_nsdi13-final170_update.pdf - Memcache

### Theoretically Important
- 08_p463-herlihy.pdf - Wait-free synchronization
- 13_sharing-mem-abd.pdf - ABD algorithm
- 18_byzantine_osdi99.pdf - Byzantine fault tolerance

## Notes

- Papers are organized by filename number for easy reference
- Each paper includes the original publication venue and year when available
- Reading order recommendations are based on complexity and dependencies
- Specialized topic groupings help focus study on specific areas
- Impact ratings help prioritize reading based on influence and relevance

## Contributing

When adding new papers to this collection:
1. Use the format: `XX_paper-name.pdf` where XX is the next sequential number
2. Update this README with the paper's information
3. Include the paper in appropriate topic sections
4. Update reading order recommendations if needed

---

*This collection represents the foundational knowledge required for understanding modern distributed systems. Each paper contributes unique insights to the field and together they form a comprehensive body of knowledge for distributed systems practitioners and researchers.*
