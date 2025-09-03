## Reference

This folder contains foundational material and concepts for understanding and implementing distributed systems. It complements the `app/` implementations and `research/` case studies by providing clear explanations, definitions, and links to further reading.

### Purpose

- Establish common terminology and abstractions
- Explain core protocols and algorithms
- Provide concise notes that bridge theory and practice
- Curate high-quality external resources for deeper study

### How to Use This Folder

1. Start here to build a theoretical foundation before diving into code in `app/`.
2. Refer back while working through implementations to clarify concepts.
3. Use the resources to explore topics in greater depth as needed.

### Core Topics (suggested coverage)

- System models: timing, failure assumptions, partial synchrony
- Communication: RPC, message passing, idempotence, retries, backoff
- Time and ordering: clocks, Lamport timestamps, vector clocks
- Consistency models: linearizability, sequential consistency, eventual consistency
- Replication: primary/backup, chain replication, quorum systems
- Consensus and coordination: Paxos, Raft, leases, leader election
- Fault tolerance: replication, checksums, retries, circuit breakers
- Partitioning and scaling: sharding, consistent hashing, load balancing
- Storage systems: logs, LSM-trees, snapshots, compaction
- Transactions: 2PC/3PC, Sagas, isolation levels
- Observability: logging, metrics, tracing in distributed contexts

### Suggested Reading Order

1. System models and communication basics
2. Time, ordering, and consistency models
3. Replication strategies and consensus algorithms
4. Storage and transaction primitives
5. Scalability patterns and operational concerns

### Cross-Links

- Implementations and exercises: see `../app`
- Case studies and analyses: see `../research`

### External Resources

- MIT 6.5840 Distributed Systems: https://pdos.csail.mit.edu/6.824/
- Google's Introduction to Systems Design: https://www.hpcs.cs.tsukuba.ac.jp/~tatebe/lecture/h23/dsys/dsd-tutorial.html
- Additional curated references may be added alongside topic notes in this folder.