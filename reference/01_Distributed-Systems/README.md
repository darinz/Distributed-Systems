# Distributed Systems Reference Guide

A comprehensive collection of distributed systems concepts, algorithms, and real-world implementations. This folder contains detailed explanations of fundamental distributed systems topics, from basic concepts to advanced consensus algorithms and production systems.

## Table of Contents

### **Foundation Concepts**
- **[01-1_Intro.md](01-1_Intro.md)** - Introduction to distributed systems, challenges, and design principles
- **[01-2_Notes-Intro.md](01-2_Notes-Intro.md)** - Supplementary notes on core definitions and failure modes

### **Communication & Coordination**
- **[02-1_RPC.md](02-1_RPC.md)** - Remote Procedure Calls: extending local procedure calls across networks
- **[02-2_Notes-RPC.md](02-2_Notes-RPC.md)** - RPC implementation details, stubs, and error handling
- **[03-1_Lamport_Clocks.md](03-1_Lamport_Clocks.md)** - Logical clocks and causality in distributed systems
- **[03-2_Notes-Clocks.md](03-2_Notes-Clocks.md)** - Clock synchronization and temporal ordering

### **Fault Tolerance & Replication**
- **[04-1_Primary-Backup.md](04-1_Primary-Backup.md)** - Primary-backup replication for fault tolerance
- **[04-2_Primary-Backup.md](04-2_Primary-Backup.md)** - Advanced primary-backup patterns and optimizations
- **[04-3_Notes-PB.md](04-3_Notes-PB.md)** - Primary-backup implementation notes and trade-offs

### **Consistency & Snapshots**
- **[05-1_VC-Snapshot.md](05-1_VC-Snapshot.md)** - Vector clocks and distributed snapshots
- **[05-2_Notes-VC.md](05-2_Notes-VC.md)** - Vector clock algorithms and snapshot protocols
- **[06-1_Consistency.md](06-1_Consistency.md)** - Consistency models and the CAP theorem
- **[06-2_Notes-Consistency.md](06-2_Notes-Consistency.md)** - Consistency trade-offs and implementation strategies

### **Consensus Algorithms**
- **[07-1_Paxos.md](07-1_Paxos.md)** - The Paxos consensus algorithm: the foundation of distributed consensus
- **[07-2_Notes-Paxos.md](07-2_Notes-Paxos.md)** - Paxos implementation details and variants
- **[08-1_Randomized_Consensus.md](08-1_Randomized_Consensus.md)** - Randomized consensus algorithms
- **[08-2_Notes-Random.md](08-2_Notes-Random.md)** - Probabilistic consensus and Byzantine agreement

### **Advanced Concurrency**
- **[09-1_PMMC.md](09-1_PMMC.md)** - Producer-Multiple Consumer patterns
- **[09-2_Notes-PMMC.md](09-2_Notes-PMMC.md)** - Multi-producer multi-consumer coordination
- **[10-1_Wait-Free_ Registers.md](10-1_Wait-Free_ Registers.md)** - Wait-free data structures and algorithms
- **[10-2_Notes-WFR.md](10-2_Notes-WFR.md)** - Lock-free programming and wait-free implementations

### **Transaction Management**
- **[11-1_Two-Phase_Commit.md](11-1_Two-Phase_Commit.md)** - Two-phase commit protocol for distributed transactions
- **[11-2_Notes-2PC.md](11-2_Notes-2PC.md)** - 2PC implementation, failures, and alternatives

### **Production Systems**
- **[12-1_GFS.md](12-1_GFS.md)** - Google File System: distributed storage at massive scale
- **[12-2_Notes-GFS.md](12-2_Notes-GFS.md)** - GFS design principles and lessons learned
- **[13-1_BigTable.md](13-1_BigTable.md)** - Google BigTable: semi-structured data storage
- **[13-2_Notes-BigTable.md](13-2_Notes-BigTable.md)** - BigTable architecture and implementation
- **[14-1_Spanner.md](14-1_Spanner.md)** - Google Spanner: globally distributed database
- **[14-2_Notes-Spanner.md](14-2_Notes-Spanner.md)** - Spanner's TrueTime and global consistency

### **Byzantine Fault Tolerance**
- **[15-1_BFT.md](15-1_BFT.md)** - Byzantine fault tolerance and malicious failures
- **[15-2_Notes-BFT.md](15-2_Notes-BFT.md)** - BFT algorithms and practical implementations

### **Security & Cryptography**
- **[16-1_Crypto.md](16-1_Crypto.md)** - Cryptographic foundations for distributed systems
- **[16-2_Notes-Crypto.md](16-2_Notes-Crypto.md)** - Digital signatures, hash functions, and secure protocols

### **NoSQL & Key-Value Stores**
- **[17-1_Dynamo.md](17-1_Dynamo.md)** - Amazon Dynamo: eventually consistent key-value store
- **[17-2_Notes-Dynamo.md](17-2_Notes-Dynamo.md)** - Dynamo's design choices and trade-offs

### **Caching & Performance**
- **[18-1_Caches.md](18-1_Caches.md)** - Distributed caching strategies and CDNs
- **[18-2_Notes-Caches.md](18-2_Notes-Caches.md)** - Cache consistency and performance optimization

## Learning Path

### **Beginner Path** (Start Here)
1. **Introduction** - Read `01-1_Intro.md` to understand the fundamental challenges
2. **Communication** - Study `02-1_RPC.md` to learn how distributed systems communicate
3. **Time & Ordering** - Explore `03-1_Lamport_Clocks.md` for logical time concepts
4. **Fault Tolerance** - Understand `04-1_Primary-Backup.md` for basic replication

### **Intermediate Path**
5. **Consistency** - Dive into `06-1_Consistency.md` for the CAP theorem and consistency models
6. **Consensus** - Study `07-1_Paxos.md` for the most important consensus algorithm
7. **Transactions** - Learn `11-1_Two-Phase_Commit.md` for distributed transaction management
8. **Real Systems** - Explore `12-1_GFS.md` to see concepts applied in production

### **Advanced Path**
9. **Byzantine Faults** - Study `15-1_BFT.md` for handling malicious failures
10. **Cryptography** - Learn `16-1_Crypto.md` for security foundations
11. **NoSQL Systems** - Explore `17-1_Dynamo.md` for eventually consistent systems
12. **Performance** - Study `18-1_Caches.md` for optimization strategies

## Key Concepts Covered

### **Fundamental Challenges**
- **Partial Failure**: Components fail independently in distributed systems
- **Network Unreliability**: Messages can be lost, delayed, or duplicated
- **Asynchrony**: No global clock or synchronized execution
- **Concurrency**: Multiple processes operating simultaneously

### **Core Algorithms**
- **Paxos**: The gold standard for consensus in distributed systems
- **Two-Phase Commit**: Classic protocol for distributed transactions
- **Vector Clocks**: Logical time for ordering events in distributed systems
- **Primary-Backup**: Basic replication for fault tolerance

### **Consistency Models**
- **Strong Consistency**: All nodes see the same data simultaneously
- **Eventual Consistency**: System converges to consistent state over time
- **CAP Theorem**: Trade-offs between Consistency, Availability, and Partition tolerance

### **Production Systems**
- **Google File System (GFS)**: Petabyte-scale distributed storage
- **BigTable**: Semi-structured data storage for Google services
- **Spanner**: Globally distributed database with external consistency
- **Dynamo**: Amazon's eventually consistent key-value store

## Practical Applications

### **When to Use Each Pattern**
- **RPC**: When you need synchronous communication between services
- **Paxos**: When you need strong consistency and can tolerate some latency
- **Eventual Consistency**: When availability is more important than perfect consistency
- **Primary-Backup**: When you need simple fault tolerance for critical services

### **Common Trade-offs**
- **Consistency vs. Availability**: Strong consistency reduces availability during partitions
- **Latency vs. Throughput**: Optimizing for one often hurts the other
- **Complexity vs. Reliability**: More reliable systems are usually more complex
- **Cost vs. Performance**: Better performance typically means higher costs

## Study Approach

### **For Each Topic**
1. **Read the main document** (e.g., `07-1_Paxos.md`) for comprehensive understanding
2. **Review the notes** (e.g., `07-2_Notes-Paxos.md`) for implementation details
3. **Understand the trade-offs** and when to use each approach
4. **Connect to real systems** that implement these concepts

### **Key Questions to Ask**
- What problem does this solve?
- What are the assumptions and limitations?
- How does it handle failures?
- What are the performance characteristics?
- When would you use this in practice?

## Prerequisites

### **Required Knowledge**
- Basic understanding of computer networks
- Familiarity with concurrent programming concepts
- Understanding of basic data structures and algorithms
- Some experience with system design

### **Helpful Background**
- Database systems (ACID properties, transactions)
- Operating systems (processes, threads, synchronization)
- Computer networks (TCP/IP, latency, bandwidth)
- Probability theory (for randomized algorithms)

## Related Resources

### **Classic Papers**
- Lamport, Leslie. "Time, clocks, and the ordering of events in a distributed system."
- Lamport, Leslie. "The part-time parliament." (Paxos)
- Fischer, Michael J., Nancy A. Lynch, and Michael S. Paterson. "Impossibility of distributed consensus with one faulty process."
- Ghemawat, Sanjay, et al. "The Google file system."

### **Books**
- *Designing Data-Intensive Applications* by Martin Kleppmann
- *Distributed Systems: Concepts and Design* by George Coulouris
- *Introduction to Reliable and Secure Distributed Programming* by Christian Cachin

## Getting Started

1. **Start with the introduction** (`01-1_Intro.md`) to build intuition
2. **Follow the learning path** based on your experience level
3. **Focus on understanding trade-offs** rather than memorizing algorithms
4. **Connect concepts to real systems** you use or have heard of
5. **Practice with system design problems** to apply your knowledge

Remember: Distributed systems are complex, but understanding the fundamental principles will help you design better systems and make informed decisions about trade-offs in your own projects.

---

*This reference guide covers the essential concepts needed to understand and design distributed systems. Each document builds upon previous concepts, so following the suggested learning path will provide the most comprehensive understanding.*
