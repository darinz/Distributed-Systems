# Lab 4B: Sharded Key/Value Server

## Overview

In this part you will implement `shardkv`, a sharded, fault-tolerant key/value storage system. Each server participates in a Paxos-based replica group that is responsible for a subset of shards. A single shardmaster service (from Part A) assigns shards to groups and reconfigures assignments over time. Your tasks are to serve client operations, handle dynamic reconfiguration correctly, and maintain sequential consistency.

You will modify these files:
- `src/shardkv/common.go`
- `src/shardkv/client.go`
- `src/shardkv/server.go`

When finished, run the tests:

```bash
$ cd src/shardkv
$ go test
Test: Basic Join/Leave ...
  ... Passed
Test: Shards really move ...
  ... Passed
Test: Reconfiguration with some dead replicas ...
  ... Passed
Test: Concurrent Put/Get/Move ...
  ... Passed
Test: Concurrent Put/Get/Move (unreliable) ...
  ... Passed
PASS
ok      shardkv 62.350s
```

## Requirements

- **Sequential consistency**: All completed client calls to `Clerk.Get`, `Clerk.Put`, and `Clerk.PutHash` must appear in the same total order across all replicas. A `Get` must observe the most recent preceding `Put/PutHash` to the same key.
- **Fault tolerance**: Assume a majority of each Paxos group, of shardmaster servers, and of any other involved group are available and can reach each other. Your system must continue operating (including reconfiguration) if a minority are crashed, slow, or partitioned.
- **No auto-Join**: Your server must not call the shardmaster's `Join()` itself; the tester orchestrates `Join/Leave/Move`.

## Client Behavior

- Use `key2shard(key)` to map a key to a shard.
- Look up the responsible group for that shard in the current shardmaster `Config`.
- Try each server in that group until one answers.
- If the server replies `ErrWrongGroup`, refresh configuration (`Query(-1)`) and retry.
- Support at-most-once semantics via `(Client, Seq)` fields carried with every RPC.

## Server Responsibilities

- Export RPCs: `Get`, `Put` (including PutHash), and any internal RPCs you design for shard transfer (the provided code uses `Inquire`).
- Reject requests for shards not owned by this group (reply `ErrWrongGroup`).
- Log every client operation in Paxos and apply operations in log order only.
- Periodically poll the shardmaster in `tick()` for newer configurations; process them sequentially.
- During reconfiguration, transfer both key/value state for relevant shards and the client deduplication state needed for at-most-once semantics.

## Reconfiguration Protocol (Suggested)

Process reconfigurations as normal Paxos operations to serialize them with client requests:
1. In `tick()`, detect a newer configuration number from the shardmaster.
2. For each shard that moves into your group in the new config:
   - Contact the previous owner group to fetch the shard’s key/value data and client deduplication metadata (e.g., highest seen `(Client -> Seq)` and last replies).
   - If the previous owner reports `ErrNoConfig`, it hasn’t reached that configuration yet; try again next tick.
3. Propose a `Reconfigure` operation in Paxos that installs the new configuration and merges received state into local state.
4. After agreement, apply the state and update `config`.

This ensures a single, agreed order among: client ops, reconfig steps, and resultant state changes.

## At-Most-Once Semantics

- Each client includes `(Client, Seq)` in every RPC.
- The server tracks, per client: `highestSeqSeen` and `lastReply`.
- If a request’s `Seq` ≤ `highestSeqSeen`, return the recorded `lastReply` without re-executing.
- Include this deduplication metadata in shard transfers; the receiver must merge it (taking the maximum `Seq` per client and the corresponding `lastReply`).

## Handling ErrWrongGroup

- Server side: If a `Get/Put` targets a shard not currently owned, return `ErrWrongGroup` and do not update deduplication state.
- Client side: Upon `ErrWrongGroup`, refresh the config and retry the same `(Client, Seq)` (do not bump `Seq` on retries), so duplicates are handled correctly after ownership changes.

## Paxos Usage

- Treat every client `Get/Put/PutHash` and every `Reconfigure` as Paxos operations.
- Only apply state changes by replaying agreed log entries in order (e.g., via a single `execute(op)` path).
- Call `px.Done(seq)` promptly after applying entries so Paxos can discard old log entries.

## State Transfer

- During reconfig, the previous owner of a shard provides:
  - The shard’s key/value pairs.
  - Dedup metadata: for each client, the highest `Seq` seen and the corresponding `lastReply`.
- The receiver merges both maps carefully (deep copy any map values that will be retained) and installs them in a single, agreed `Reconfigure` step.

## Gob and Types

- Ensure all concrete types placed in Paxos `Op.Args` are registered with `gob` on server startup, e.g.:
  - `GetArgs{}`, `PutArgs{}`, and your `ReconfigArgs{}` (or equivalent).
- A missing registration manifests as a non-fatal `gob: type not registered` error that still breaks the lab’s logic; check logs.

## Common Pitfalls

1. Updating the KV store directly in RPC handlers (bypassing the Paxos log). Always propose and apply via the log.
2. Forgetting to propagate dedup state during shard movement.
3. Applying reconfigurations out of order or concurrently. Process one config at a time in config-number order.
4. Failing to return `ErrWrongGroup` when ownership changes mid-flight.
5. Incrementing client `Seq` on client-side retries after `ErrWrongGroup`.
6. Not calling `px.Done()` and leaking Paxos state.
7. Shallow copying maps placed inside Paxos ops; ensure deep copies where needed to avoid later mutation bugs.

## Testing Guidance

- Start by passing:
  - Basic Join/Leave
  - Shards really move
  - Reconfiguration with some dead replicas
- Then focus on concurrent tests; they are non-deterministic—run multiple times.

## Performance and Timing

- `tick()` should poll periodically (e.g., 100–300ms) and process at most one configuration step per tick if needed.
- Reconfiguration and shard transfer may take several ticks while groups converge on the same configuration.
- End-to-end tests typically finish in ~45–70 seconds (implementation dependent).

## Implementation Outline

1. Add `(Client, Seq)` fields to client RPC args and track dedup in the server.
2. Implement a single `execute(op)` path that:
   - Proposes/awaits Paxos decision for the current sequence.
   - Applies `Get` (record last read for dedup), `Put/PutHash` (update store), and `Reconfigure` (merge state and install new config).
3. Implement `tick()` to discover and serialize reconfiguration steps.
4. Implement shard state transfer RPC(s) and merge logic.
5. Ensure all relevant types are `gob.Register`-ed at startup.

## Reference: What the Tests Expect

- `ErrWrongGroup` for keys whose shard is not owned by the contacted group.
- Stable client-observed order under concurrent `Move` and client ops.
- No duplicate effects from re-sent client requests.
- Successful shard migration even with some dead replicas.

---

## Additional Hints
- Keep retry/backoff simple on the client to avoid hot spinning; refresh config only on `ErrWrongGroup`.
- Use a per-client unique ID + monotonic `Seq`; never increment `Seq` on retries.
- Use helpers to deep copy any maps you place into a Paxos `Op`.
- Add minimal debug prints gated by a `Debug` flag to inspect config numbers and shard ownership during development.
