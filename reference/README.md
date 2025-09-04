## Reference

This folder contains comprehensive foundational material and concepts for understanding and implementing distributed systems. It provides detailed explanations, algorithms, and real-world system analyses that complement the `app/` implementations and `research/` case studies.

### Main Content

**[01_Distributed-Systems/](01_Distributed-Systems/)** - Complete distributed systems reference guide covering:

#### **Foundation Concepts**
- Introduction to distributed systems challenges and design principles
- Core definitions, failure modes, and system models

#### **Communication & Coordination** 
- Remote Procedure Calls (RPC) and network communication
- Logical clocks, Lamport timestamps, and causality

#### **Fault Tolerance & Replication**
- Primary-backup replication patterns
- Vector clocks and distributed snapshots
- Consistency models and the CAP theorem

#### **Consensus & Coordination**
- Paxos consensus algorithm (the gold standard)
- Randomized consensus and Byzantine agreement
- Two-phase commit for distributed transactions

#### **Advanced Topics**
- Wait-free data structures and lock-free programming
- Byzantine fault tolerance and security
- Production systems: GFS, BigTable, Spanner, Dynamo
- Distributed caching and performance optimization

### Learning Paths

**Beginner**: Start with introduction → RPC → logical clocks → basic replication

**Intermediate**: Consistency models → Paxos → transactions → real systems (GFS)

**Advanced**: Byzantine faults → cryptography → NoSQL systems → performance optimization

### How to Use This Reference

1. **Start with the comprehensive guide**: `01_Distributed-Systems/README.md` for the complete overview
2. **Follow structured learning paths** based on your experience level
3. **Read main documents first** (e.g., `07-1_Paxos.md`) then review implementation notes
4. **Connect theory to practice** by studying real production systems
5. **Focus on trade-offs** rather than memorizing algorithms

### Cross-References

- **Implementations and exercises**: see `../app/`
- **Case studies and analyses**: see `../research/`

### Key Concepts Covered

- **Fundamental Challenges**: Partial failure, network unreliability, asynchrony, concurrency
- **Core Algorithms**: Paxos, Two-Phase Commit, Vector Clocks, Primary-Backup
- **Consistency Models**: Strong consistency, eventual consistency, CAP theorem
- **Production Systems**: GFS, BigTable, Spanner, Dynamo

### External Resources

- [MIT 6.5840 Distributed Systems](https://pdos.csail.mit.edu/6.824/)
- [Designing Data-Intensive Applications](https://dataintensive.net/) by Martin Kleppmann
- Classic papers: Lamport's time/clock paper, Paxos paper, FLP impossibility result

### Prerequisites

- Basic computer networks knowledge
- Concurrent programming concepts
- Data structures and algorithms
- Some system design experience

---

*This reference provides the essential foundation for understanding and designing distributed systems. Each topic includes both theoretical explanations and practical implementation details, with clear learning paths for different experience levels.*