## Applications and Code Examples

This folder contains hands-on programming exercises and example implementations of distributed systems concepts. Use these projects to apply ideas from `reference/` and to reinforce understanding through practice.

### What’s Inside

- `01_Practice-Labs/` — Guided labs (MapReduce, shardmaster, sharded KV, persistence)
  - Clear step-by-step instructions and hints in: `lab4a.md`, `lab4b.md`, `lab5.md`
  - Updated to work with Go 1.25.1; see the README inside for run/test commands
- Additional app folders (optional scaffolds):
  - `rpc/` — Remote Procedure Call basics, idempotence, retries, backoff
  - `kv/` — Key-value store with replication and consistency variants
  - `raft/` — Raft consensus implementation and tests
  - `shard/` — Sharded/partitioned services and rebalancing
  - `storage/` — Log-structured storage, snapshots, compaction
  - `coordination/` — Leader election, leases, service discovery
  - `observability/` — Metrics, logging, tracing across services

You can organize projects differently; the above layout is a suggested starting point.

### Quickstart

1. Ensure Go 1.25.1 is installed:
   ```bash
   go version
   # go version go1.25.1 ...
   ```
2. Open the Practice Labs README for per-lab instructions:
   ```bash
   cd app/01_Practice-Labs
   ```
3. Run example MapReduce word count:
   ```bash
   cd app/01_Practice-Labs/src/main
   go run wc.go master kjv12.txt sequential
   ```
4. Run Shardmaster tests (Lab 4A):
   ```bash
   cd app/01_Practice-Labs/src/shardmaster
   go test
   ```
5. Run Sharded KV tests (Lab 4B):
   ```bash
   cd app/01_Practice-Labs/src/shardkv
   go test
   ```
6. Run Persistent KV tests (Lab 5):
   ```bash
   cd app/01_Practice-Labs/src/diskv
   go test -run Test4   # Lab 4 compatibility subset
   go test              # Full Lab 5 suite
   ```

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

### Environment

- Go 1.25.1 or later is recommended for all examples in this repository
- Standard tooling: `go test`, `go build`, `go run`

### Go Programming Resources

- [Go Language](https://go.dev) - Official Go programming language website
- [Effective Go](https://go.dev/doc/effective_go) - Tips for writing clear, performant, and idiomatic Go code
- [A Tour of Go](https://go.dev/tour/welcome/1) - Interactive introduction to Go programming

### Cross-Links

- Concepts and theory: see `../reference`
- Case studies and analyses: see `../research`