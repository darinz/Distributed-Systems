# Distributed Systems Mini-Project

## Overview

This mini-project is designed to give you hands-on experience with distributed systems concepts through practical implementation. You will choose a project that interests you and build a working distributed system that demonstrates key principles learned throughout the course.

## Project Requirements

### Scope and Timeline
- **Duration**: This should be a substantial project that demonstrates deep understanding of distributed systems concepts
- **Scope**: Choose a project that is challenging but achievable within the given timeframe
- **Documentation**: Include comprehensive documentation explaining your design decisions and implementation

### Programming Language
You can use whatever programming language you want for the project. Two excellent choices are:

- **C++**: Offers fine-grained control over system resources and performance
- **Go**: Excellent built-in concurrency primitives and networking libraries

Other good options include:
- **Rust**: Memory safety with high performance
- **Java**: Rich ecosystem of distributed systems libraries
- **Python**: Rapid prototyping and extensive libraries

### Evaluation Criteria
Your project will be evaluated based on:
- **Technical depth**: Demonstration of distributed systems concepts
- **Code quality**: Clean, well-documented, and maintainable code
- **Innovation**: Creative solutions to distributed systems challenges
- **Testing**: Comprehensive testing including failure scenarios
- **Documentation**: Clear explanation of design and implementation

## Project Ideas

The following list provides inspiration for your project. These ideas are organized by category and include detailed descriptions of the concepts involved and challenges to address.

### 1. Communication and Collaboration Systems

#### Real-time Communication Tools
**Build better tools for remote collaboration (text, voice, or video chat).**
- **Focus**: Real-time messaging, voice/video streaming, presence systems
- **Concepts**: WebRTC, peer-to-peer networking, NAT traversal, real-time protocols
- **Challenge**: Handle network partitions, ensure message delivery, manage bandwidth

#### Network Object Systems
**Build a network object system for C++.**
- **Focus**: Distributed object-oriented programming, remote method invocation
- **Concepts**: RPC, object serialization, distributed garbage collection
- **Challenge**: Handle object lifecycle, manage distributed references, ensure consistency

### 2. Consensus and Replication Protocols

#### Protocol Testing and Verification
**Build a simple system that takes an implementation of these protocols and systematically explores their behavior in the face of crashes and network partitioning.**
- **Focus**: 2PC, Paxos, and other consensus protocols
- **Concepts**: Failure injection, systematic testing, protocol verification
- **Challenge**: Create comprehensive test scenarios, detect protocol violations
- **Reference**: See here for an example of how to do this for file systems

#### Raft Implementation Testing
**Build a checking infrastructure that can plug into the many different RAFT implementations and find protocol errors.**
- **Focus**: Raft protocol correctness, implementation comparison
- **Concepts**: Consensus protocols, distributed testing, correctness verification
- **Challenge**: Define correctness criteria, create comprehensive test suite
- **Reference**: Look at what Kyle Kingsbury has done with Jepsen

#### View Stamped Replication
**Build a clean, simple implementation of view stamped replication based on the updated Liskov paper.**
- **Focus**: View stamped replication, distributed consensus
- **Concepts**: View changes, replica coordination, failure handling
- **Challenge**: Implement clean API, handle view changes, ensure correctness

#### Byzantine Fault Tolerance
**Design and implement a Byzantine-fault-tolerant version of Raft.**
- **Focus**: Byzantine fault tolerance, consensus with malicious nodes
- **Concepts**: Cryptographic signatures, Byzantine agreement, fault tolerance
- **Challenge**: Handle malicious behavior, ensure safety despite Byzantine failures

**Design and implement a Byzantine-fault-tolerant state machine replication system that uses witnesses.**
- **Focus**: Witness-based replication, reducing storage requirements
- **Concepts**: Witness nodes, 3f+1 servers with fewer than 3f+1 state copies
- **Challenge**: Design witness protocol, ensure consistency with reduced storage

### 3. Storage Systems

#### Distributed File Systems
**Build a large file store, like GFS, and possibly using RAID like Zebra.**
- **Focus**: Distributed file systems, fault tolerance, scalability
- **Concepts**: Chunk servers, master coordination, replication, load balancing
- **Challenge**: Handle large files, ensure availability, manage metadata

#### Virtual Disk Systems
**Build a scalable virtual disk like Petal.**
- **Focus**: Virtual disk abstraction, distributed storage
- **Concepts**: Block-level storage, virtualization, distributed RAID
- **Challenge**: Provide consistent block interface, handle failures
- **Reference**: Maybe built using the Intel Open Storage Toolkit

#### Distributed Caching
**Build a scalable web cache using consistent hashing or CARP.**
- **Focus**: Distributed caching, load balancing, cache coherence
- **Concepts**: Consistent hashing, cache replacement, distributed coordination
- **Challenge**: Handle cache misses, ensure consistency, balance load

### 4. Coordination and Synchronization

#### Distributed Locking Service
**Build a simplified version of a synchronization service like Google's Chubby.**
- **Focus**: Distributed locking, configuration management, leader election
- **Concepts**: Paxos, leases, session management, failure detection
- **Challenge**: Handle client failures, ensure lock safety, provide high availability

#### File Synchronization
**Build a file synchronization tool like tra.**
- **Focus**: File synchronization, conflict resolution, version control
- **Concepts**: Eventual consistency, conflict-free replicated data types
- **Challenge**: Handle conflicts, ensure convergence, optimize bandwidth

### 5. Development and Debugging Tools

#### Distributed Build Systems
**Build a simple, automatic distributed-parallel make implementation.**
- **Focus**: Distributed compilation, dependency analysis, parallel execution
- **Concepts**: Dependency graphs, parallel execution, load balancing
- **Challenge**: Infer true dependencies, handle dynamic dependencies, optimize build time
- **Technique**: Intercept "open()" system calls to see file dependencies

#### Distributed Debugging
**Build a parallel debugger that allows you to debug distributed systems.**
- **Focus**: Distributed debugging, execution tracing, message flow
- **Concepts**: Distributed breakpoints, message tracing, execution replay
- **Challenge**: Follow execution across nodes, handle timing issues
- **Reference**: Ideally using some modification of GDB

#### Distributed Profiling
**Build a distributed profiler that allows you to observe where time really goes in a distributed system.**
- **Focus**: Performance analysis, bottleneck identification, distributed tracing
- **Concepts**: Profiling, tracing, performance measurement
- **Challenge**: Correlate events across nodes, identify bottlenecks
- **Requirement**: Use it to spot bottlenecks in at least one existing distributed system

### 6. Protocol Visualization and Verification

#### Protocol Visualization
**Build a raftscope-like visualization tool for a different protocol.**
- **Focus**: Protocol visualization, interactive debugging, educational tools
- **Concepts**: Protocol state machines, message flow, failure scenarios
- **Challenge**: Create intuitive visualizations, handle complex protocols
- **Reference**: Similar to raftscope for Raft

#### Formal Verification
**Formally model and verify a consensus protocol using TLA+ or IVy.**
- **Focus**: Formal methods, protocol verification, correctness proofs
- **Concepts**: Temporal logic, model checking, formal specification
- **Challenge**: Create accurate models, prove correctness properties
- **Tools**: TLA+, IVy, or similar formal verification tools

### 7. Security and Privacy

#### Secure Communication
**Build a message-level interposition library that adds security to existing networked services.**
- **Focus**: Security middleware, transparent encryption, authentication
- **Concepts**: Encryption, authentication, nonces, secure checksums
- **Challenge**: Transparent integration, performance impact, key management
- **Reference**: Similar to VPNs, but more comprehensive

#### Privacy-Preserving Systems
**Build a mobile-phone based privacy-preserving contact-tracing system.**
- **Focus**: Privacy-preserving protocols, contact tracing, epidemic modeling
- **Concepts**: Differential privacy, cryptographic protocols, distributed computation
- **Challenge**: Balance privacy and utility, handle scale, ensure security

#### Blockchain-Based Security
**Build a blockchain-based key management server that is more secure than current PGP key servers.**
- **Focus**: Key management, blockchain security, identity verification
- **Concepts**: Blockchain, cryptographic keys, distributed trust
- **Challenge**: Ensure key authenticity, handle key revocation, provide privacy
- **Enhancement**: Optionally provide increased privacy via techniques like CONIKS

### 8. Modern Distributed Systems

#### Conflict-Free Replicated Data Types
**Build a replicated system that leverages CRDTs to achieve eventual state convergence.**
- **Focus**: CRDTs, eventual consistency, conflict resolution
- **Concepts**: Commutative operations, idempotent updates, convergence
- **Challenge**: Design appropriate CRDTs, handle complex data types

#### Asynchronous Programming
**Design an asynchronous RPC library using C++20 coroutines.**
- **Focus**: Asynchronous programming, coroutines, high-concurrency I/O
- **Concepts**: Coroutines, async/await, stackless concurrency
- **Challenge**: Hide stack ripping, provide intuitive API, ensure performance

#### Blockchain Integration
**Modify an open-source database to use a public blockchain as a two-phase commit coordinator.**
- **Focus**: Blockchain integration, distributed transactions, cross-system consistency
- **Concepts**: Two-phase commit, blockchain consensus, distributed transactions
- **Challenge**: Integrate with existing database, handle blockchain delays
- **Reference**: See this paper for inspiration

#### Batch Transaction Processing
**Build a replicated database that handles transactions in batches with commutative semantics.**
- **Focus**: Batch processing, commutative transactions, parallel execution
- **Concepts**: Transaction batching, commutative operations, parallel processing
- **Challenge**: Design commutative semantics, ensure replicability, maintain consistency

### 9. Hardware and Embedded Systems

#### Raspberry Pi Distributed Systems
**Build a distributed system using Raspberry Pi nodes and interesting cheap hardware.**
- **Focus**: Embedded distributed systems, IoT, hardware integration
- **Concepts**: Embedded systems, sensor networks, resource constraints
- **Challenge**: Handle resource limitations, ensure reliability, manage power
- **Advanced**: Build a clean, simple "bare-metal" toolkit for easy system building

### 10. Advanced Research Projects

#### Porcupine Improvements
**Build something like Porcupine that addresses some of the paper's shortcomings.**
- **Focus**: Linearizability testing, distributed system verification
- **Concepts**: Linearizability, testing frameworks, correctness verification
- **Challenge**: Identify limitations, propose improvements, implement solutions

#### System Interposition
**Build a system-call or message-level interposition library for transparent service replication.**
- **Focus**: Transparent replication, system interposition, fault tolerance
- **Concepts**: System call interception, message interception, transparent replication
- **Challenge**: Handle various protocols, ensure transparency, manage state
- **Reference**: Something similar but more complicated than what you would build: parrot

## Getting Started

### Project Selection
1. **Choose your focus**: Select a project that aligns with your interests and demonstrates key distributed systems concepts
2. **Define scope**: Ensure your project is challenging but achievable within the timeframe
3. **Plan milestones**: Break your project into manageable milestones with clear deliverables

### Development Process
1. **Research**: Study existing systems and papers related to your chosen project
2. **Design**: Create a detailed design document explaining your approach
3. **Implement**: Build your system with clean, well-documented code
4. **Test**: Create comprehensive tests including failure scenarios
5. **Document**: Write clear documentation explaining your design and implementation

### Resources
- **Papers**: Study relevant research papers for your chosen domain
- **Open source**: Examine existing implementations for inspiration
- **Tools**: Use appropriate development and testing tools
- **Community**: Engage with distributed systems communities for guidance

## Submission Requirements

Your project submission should include:
- **Source code**: Complete, well-documented implementation
- **Design document**: Detailed explanation of your approach and design decisions
- **Testing report**: Results of your testing including failure scenarios
- **Demo**: Working demonstration of your system
- **Reflection**: Analysis of what you learned and how you would improve the system

Remember: The goal is not just to build a working system, but to demonstrate deep understanding of distributed systems concepts through practical implementation.

## Additional Considerations

### Technical Challenges
When selecting and implementing your project, consider these technical challenges:

- **Fault tolerance**: How will your system handle node failures, network partitions, and other failures?
- **Consistency**: What consistency model will you use, and how will you ensure it?
- **Scalability**: How will your system scale as the number of nodes or load increases?
- **Performance**: What are the performance characteristics of your system?
- **Security**: What security considerations are relevant to your project?

### Implementation Tips
- **Start simple**: Begin with a basic version and add complexity gradually
- **Test thoroughly**: Include unit tests, integration tests, and failure scenario tests
- **Document everything**: Keep detailed documentation of your design decisions
- **Use existing tools**: Leverage existing libraries and frameworks where appropriate
- **Get feedback**: Share your progress with peers and instructors

### Common Pitfalls to Avoid
- **Over-engineering**: Don't make your system more complex than necessary
- **Under-testing**: Ensure you test failure scenarios, not just happy paths
- **Poor documentation**: Document your design decisions and implementation details
- **Ignoring performance**: Consider the performance implications of your design choices
- **Security oversights**: Don't forget about security considerations in your design

Good luck with your distributed systems mini-project!
