# Distributed Systems

[![Status: Active Development](https://img.shields.io/badge/Status-Active%20Development-brightgreen.svg)](https://github.com/yourusername/Distributed-Systems)

A comprehensive repository covering the theory, implementation, and practical applications of distributed systems. This repository serves as both a learning resource and a hands-on implementation guide for building robust, scalable distributed systems.

## Overview

Distributed systems power modern computing: web applications, e-commerce, content delivery, and cloud infrastructure. This repository provides a structured path to learn and build distributed systems through:

- **Theoretical foundations** and abstractions
- **Practical implementation** techniques
- **Real-world case studies** and research
- **Hands-on programming** applications

## Quickstart

1. Install Go 1.25.1 or later:
   ```bash
   go version
   # go version go1.25.1 ...
   ```
2. Explore the Practice Labs (guided exercises with code and tests):
   ```bash
   cd app/01_Practice-Labs
   ```
3. Run a quick example (MapReduce word count):
   ```bash
   cd app/01_Practice-Labs/src/main
   go run wc.go master kjv12.txt sequential
   ```
4. Run shardmaster tests (Lab 4A):
   ```bash
   cd app/01_Practice-Labs/src/shardmaster
   go test
   ```
5. Run sharded KV tests (Lab 4B):
   ```bash
   cd app/01_Practice-Labs/src/shardkv
   go test
   ```
6. Run persistence tests (Lab 5):
   ```bash
   cd app/01_Practice-Labs/src/diskv
   go test -run Test4   # Lab 4 subset
   go test              # Full Lab 5
   ```

## Repository Structure

```
Distributed-Systems/
├── app/                    # Applications and code examples
│   └── README.md             # Implementation guides and exercises
├── reference/             # Reference material and concepts
│   └── README.md             # Theoretical foundations and resources
├── research/              # Research and case studies
│   └── README.md             # Real-world system analyses
└── README.md                 # This file
```

### Directory Descriptions

- **`app/`** – Hands-on programming applications and code examples
  - Implementation exercises and working samples
  - Updated `01_Practice-Labs` with clearer instructions for Lab 4A/4B/5
- **`reference/`** – Theoretical foundations and learning resources
  - Core concepts, abstractions, and implementation techniques
- **`research/`** – Case studies and research in distributed systems
  - Real-world system analyses, papers, and industry implementations

## Learning Objectives

By progressing through this repository, you will gain experience with:

- **Client-server computing** and service architectures
- **Remote procedure calls (RPC)** and inter-service communication
- **Distributed storage** and consistency models
- **Consensus** (Paxos/Raft) and replication
- **Fault tolerance** and high availability
- **Scaling** and performance optimization
- **Correctness** under failures, partitions, and retries

## Recommended Learning Path

1. Read core concepts: `reference/`
2. Study real systems: `research/`
3. Implement and test: `app/` (start with `01_Practice-Labs`)

## Why This Repository?

- **Hands-on focus** – Build working systems, not just read about them
- **Balanced approach** – Theory meets practical engineering
- **Modern relevance** – Techniques aligned with current industry practices

## External Resources

- **MIT 6.5840 Distributed Systems**: [Course Materials](https://pdos.csail.mit.edu/6.824/)
- Additional resources and references are in `reference/`

## Contributing

Contributions are welcome:

- Improve implementations and tests
- Add case studies and research notes
- Enhance documentation and learning guides

---

**Note**: This repository is under active development. New content, examples, and improvements are added regularly.