## Applications and Code Examples

This folder contains hands-on programming exercises and example implementations of distributed systems concepts. Use these projects to apply ideas from `reference/` and to reinforce understanding through practice.

### Purpose

- Provide concrete implementations of core distributed systems techniques
- Offer incremental exercises that build toward full systems
- Encourage experimentation, benchmarking, and failure-injection testing

### Layout (suggested)

- `rpc/` — Remote Procedure Call basics, idempotence, retries, backoff
- `kv/` — Key-value store with replication and consistency variants
- `raft/` — Raft consensus implementation and tests
- `shard/` — Sharded/partitioned services and rebalancing
- `storage/` — Log-structured storage, snapshots, compaction
- `coordination/` — Leader election, leases, service discovery
- `observability/` — Metrics, logging, tracing across services

You can organize projects differently; the above layout is a suggested starting point.

### How to Use This Folder

1. Read the relevant concept notes in `../reference`.
2. Start with the simplest baseline implementation.
3. Add features incrementally; run tests and measure performance after each step.
4. Perform failure testing (node crash, message loss/reorder, partitions) and verify correctness.

### Suggested Exercises

- Implement a minimal RPC library with at-least-once semantics and idempotent handlers
- Build a single-primary replicated key-value store; add linearizability checks
- Implement Raft: leader election, log replication, snapshotting
- Add sharding with consistent hashing; implement rebalancing with minimal disruption
- Introduce client-side load balancing and retries with exponential backoff and jitter
- Add tracing across microservices; correlate logs with request IDs

### Running Examples

Each subproject should include its own README with:

- Setup instructions and dependencies
- How to run the server(s) and clients
- How to run tests and benchmarks
- Known limitations and future work

### Go Programming Resources

- [Go Language](https://go.dev) - Official Go programming language website
- [Effective Go](https://go.dev/doc/effective_go) - Tips for writing clear, performant, and idiomatic Go code
- [A Tour of Go](https://go.dev/tour/welcome/1) - Interactive introduction to Go programming

### Cross-Links

- Concepts and theory: see `../reference`
- Case studies and analyses: see `../research`