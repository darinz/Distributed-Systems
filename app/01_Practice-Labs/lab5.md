# Lab 5: Persistence

## Goal
Add persistence to your sharded key/value server so each replica can crash and restart while preserving safety and quickly rejoining its replica group. After restart, the system’s availability should be no worse than if the same servers were temporarily disconnected.

You must persist at least:
- Key/value data owned by the server (per-shard files are fine)
- Paxos state needed to recover/continue the protocol
- Any metadata required for at-most-once semantics and reconfiguration (e.g., client dedup info, current config number, current sequence number)

You do not need a high-performance format. Simple files (one per key) plus a few metadata files are fine.

---

## Files to edit (in `src/diskv/`)
- `server.go`: server logic, persistence helpers, recovery
- `common.go`: shared types, errors, helpers
- `client.go`: client logic (unchanged semantics; uses new types and ops)
- Do NOT copy `test_test.go` from earlier labs; this lab has its own tests.

You will also use `main/diskvd.go` (already provided) as the entrypoint process the tests launch, which calls `StartServer`.

---

## What changed from Lab 4
- Operations: Replace PutHash with Put/Append. If `op == "Append"`, append to existing value; otherwise replace.
- Server bootstrap: `StartServer` now takes `(gid, shardmasters, servers, me, dir, restart)`.
  - `dir`: per-server directory to store its state
  - `restart`: false for the very first boot; true for process restarts
- Tests run each key/value server as a separate UNIX process. Your code must be robust to crashes at arbitrary points, including while writing files.

---

## Persistence and recovery requirements
1. KV data
   - Store each key under its shard directory (`./shard-<id>/key-<encoded>`). Use an atomic write pattern: write to a temp file, then `os.Rename(temp, final)`.
   - On restart, rebuild in-memory KV state by reading only the shards assigned to this group according to the recovered configuration.

2. Client deduplication metadata
   - Persist per-client highest seen sequence and last reply so you can return the same reply if a client retries after a crash.
   - Suggested layout: two subdirectories under `dir/`: `ids/` and `results/`, where each file name is `client-<id>`.

3. Paxos state
   - Persist enough proposer/acceptor state to continue safely after restart (e.g., highest promised, accepted value, decided, etc.).
   - The provided Paxos library accepts parameters to store on disk; use the `paxosdir := dir + "/state/"` directory and pass the `restart` flag.

4. Configuration tracking
   - Persist the current config number so that after restart you can query the shardmaster for the corresponding full `Config` and know which shards to serve.
   - Suggested files: `seq` (current sequence number) and `config` (current config number), written atomically with temp+rename.

5. In-progress operation handling
   - The tester may kill a server mid-operation. To remain consistent:
     - Log the Paxos-chosen operation for `seq` to disk before side effects (e.g., `op-<seq>`), then apply.
     - For Append, if you need to atomically update the file representing a key, consider a helper file recording the final value to write (`append-<seq>`). On recovery, if the helper exists, complete the write.

---

## Disk layout (suggested)
```
<dir>/
  seq                 # current sequence number (int)
  config              # current config number (int)
  state/              # paxos state directory (library-managed)
  shard-<sid>/        # per-shard data (owned shards only)
    key-<base32(key)>
    temp-<base32(key)>
  ids/                # client highest-seen sequence
    client-<cid>
    temp-<cid>
  results/            # client last reply
    client-<cid>
    temp-<cid>
  op-<seq>            # durable log of chosen op before effects
  op-temp<seq>
  append-<seq>        # durable final value for Append, if used
  temp-append-<seq>
```

Use `base32` for filenames derived from keys (Mac filesystems can be case-insensitive).

---

## Crash-safety hints
- Always write to a temporary file and then `os.Rename(temp, final)`. `Rename` is atomic on the same filesystem and helps you avoid partially written files.
- Consider the order:
  1) Persist metadata that you will need to complete or re-apply the operation (e.g., `op-<seq>`, temp append content, client info) 
  2) Apply side effects to the KV store files atomically (temp+rename per key)
  3) Mark forward progress (e.g., bump `seq`, GC old temp files)
- After restart:
  - Reconstruct `seq` and `config` from disk
  - Rebuild `store` from on-disk files of shards you own in that `config`
  - If an `op-<seq>` exists, finish applying it (for Append, check `append-<seq>`), then continue normal operation

---

## Reconfiguration flow
- On each tick, ask the shardmaster for new configs beyond your current `config.Num`.
- For shards moving into your group, contact the old group with `Inquire` to fetch:
  - Key/value data for those shards
  - Client dedup state (highest seq per client, last replies)
- Submit one Paxos `Reconfigure` op bundling the incoming shard data and dedup maps.
- Apply and persist: write shard files, update dedup files, bump `config`.

Hint: If a source group hasn’t yet reached the target config (it replies `ErrNoConfig`), back off and retry later.

---

## Step-by-step implementation plan
1. Wire persistence directories
   - In `StartServer`, accept `dir` and `restart`. Create subdirs (`ids/`, `results/`, `state/`) when `!restart`.
2. Persist and recover minimal metadata
   - Implement helpers: write/read `seq`, write/read `config` (with temp+rename).
   - On restart, read `seq` and `config` first, query shardmaster for that config, and rebuild in-memory state from files for owned shards.
3. Client at-most-once
   - Implement helpers to persist `ids` and `results` files per client.
   - On applying Get/Put/Append, update both memory and these files first to make replies idempotent.
4. Paxos durability
   - Use `paxos.Make(..., true, paxosdir, restart)` so the library persists and restores its state in `state/`.
   - Ensure `execute` method logs `op-<seq>` before side effects and calls `Done(seq)` after.
5. Append atomicity
   - For Append, write the “final value” for the key to `append-<seq>` before installing it into the KV file, then delete/GC that helper file after you’ve advanced `seq`.
6. Recovery of in-progress op
   - On start, if `op-<seq>` exists, re-apply the effects idempotently:
     - For Get/Put: safe to re-apply using dedup and file+mem update order
     - For Append: if `append-<seq>` exists, ensure that value is installed atomically in the KV store
7. Reconfiguration
   - Implement `Inquire` to return shard data + dedup maps
   - Implement `Reconfigure` to collect, propose via Paxos, then persistently install new shards and dedup info
8. Tick behavior
   - Periodically fetch new configs; if none, submit a `Nop` to drive Paxos progress/GC

---

## Running the tests
- Run only Lab 4-compatible tests: `go test -run Test4`
- Run full Lab 5 tests: `go test`
- Tests may take minutes; do not run multiple test suites concurrently (they clean up each other’s files).

---

## Debugging checklist
- Does every write to persistent files use temp+rename?
- After restart, do you:
  - Restore Paxos state via the library
  - Restore `seq` and `config`
  - Rebuild owned shards only (per current `config`)
  - Finish any in-progress op at `seq`
- Do you update dedup files and return consistent replies on retries?
- Do you garbage-collect helper files for completed sequences (`op-<seq-1>`, `append-<seq-1>`, temps)?
- Do you handle `ErrNoConfig` during reconfig inquiries and back off?

---

## Hints
- Use `gob` to encode/decode structured state (e.g., ops). Remember: exported (capitalized) fields only.
- Use `base32` to encode keys into filenames; avoids issues with `/` and case-insensitive FS.
- Exponential backoff on network/Paxos waits is helpful (see `waitForPaxos`).
- Keep Paxos active (submit `Nop` on idle ticks) so that logs can be GC’d and progress continues.
- Keep memory and disk state in sync. If you write a value to disk, update `store` soon after (or vice versa), but ensure recovery can complete safely if you crash in-between.

---

## Common pitfalls
- Writing files directly without temp+rename (can leave partial data on crash)
- Not persisting client dedup info before replying (can break at-most-once on restart)
- Rebuilding all shards instead of only the shards this group owns
- Forgetting to delete old helper files (disk bloat and test failures)
- Issuing multiple concurrent RPCs without proper locking (race conditions)

---

## Reference: small gob example
```go
package main

import "bytes"
import "encoding/gob"
import "fmt"

func enc(x1 int, x2 string) string {
	w := new(bytes.Buffer)
	e := gob.NewEncoder(w)
	e.Encode(x1)
	e.Encode(x2)
	return string(w.Bytes())
}

func dec(buf string) (int, string) {
	r := bytes.NewBuffer([]byte(buf))
	d := gob.NewDecoder(r)
	var x1 int
	var x2 string
	d.Decode(&x1)
	d.Decode(&x2)
	return x1, x2
}

func main() {
	buf := enc(99, "hello")
	x1, x2 := dec(buf)
	fmt.Printf("%v %v\n", x1, x2)
}
```

---

## Extra credit (optional ideas)
- Batch small writes to reduce I/O
- Write-ahead logging for all KV updates (beyond Append) to unify recovery
- Structured metrics/logging around recovery and reconfiguration


