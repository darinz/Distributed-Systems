# Distributed Systems: Supplementary Notes

## Core Definitions

**Distributed System**: An application that executes protocols to coordinate multiple processes on a network, where all components cooperate to perform related tasks.

**Key Components**:
- **Problem** → Code you write
- **Process** → Running code instance  
- **Message** → Inter-process communication
- **Packet** → Message fragment on wire
- **Protocol** → Message format + exchange rules
- **Network** → Infrastructure linking computers
- **Component** → Process or supporting hardware

## Why Build Distributed Systems?

**Advantages**:
- Connect remote users with remote resources
- **Open**: Components continuously interact
- **Scalable**: Easy to accommodate growth
- Combined capabilities > standalone systems

**Trade-off**: Complexity of simultaneous component interactions

## Essential Properties

### 1. **Fault-Tolerant**
- Recovers from component failures without incorrect actions

### 2. **Highly Available** 
- Restores operations when components fail

### 3. **Recoverable**
- Failed components can restart and rejoin

### 4. **Consistent**
- Coordinates actions despite concurrency and failure
- Enables distributed system to act like non-distributed

### 5. **Scalable**
- Operates correctly as system scales up
- Handles increased network size, users, servers, load

### 6. **Predictable Performance**
- Provides desired responsiveness timely

### 7. **Secure**
- Authenticates access to data and services

## The Failure Reality

> **Ken Arnold (Jini architect)**: "Failure is the defining difference between distributed and local programming"

**Key Insight**: In distributed systems, failure happens "all the time" - not just rarely.

**Partial Failure Problem**: When sending a message fails, you can't distinguish between:
- Message delivered but response lost
- Message never delivered

**Design Implication**: Simplicity becomes crucial - more interactions = more failure recovery scenarios.

## Failure Categories

### Hardware vs Software
- **Hardware**: Improved since 1980s, now mainly network/drive failures
- **Software**: 25-35% of unplanned downtime, even with rigorous testing

### Bug Types
- **Heisenbug**: Disappears when observed (more common in distributed systems)
- **Bohrbug**: Reproducible under defined conditions

### Failure Types
- **Halting**: Component stops (detectable only by timeout)
- **Fail-stop**: Halting with notification to others
- **Omission**: Message send/receive failure (no notification)
- **Network**: Link breaks
- **Network Partition**: Network fragments into disjoint sub-networks
- **Timing**: Temporal property violations (clock sync, delays)
- **Byzantine**: Data corruption, malicious behavior

## The 8 Fallacies of Distributed Computing

**Common (Wrong) Assumptions**:
1. The network is reliable
2. Latency is zero
3. Bandwidth is infinite
4. The network is secure
5. Topology doesn't change
6. There is one administrator
7. Transport cost is zero
8. The network is homogeneous

**Key Terms**:
- **Latency**: Time between request initiation and data transfer start
- **Bandwidth**: Communication channel capacity
- **Topology**: Network configuration (ring, bus, star, mesh)
- **Homogeneous**: Single network protocol
## Client-Server Architecture

### The Challenge
Building reliable systems over unreliable networks requires dealing with uncertainty:
- Processes know their own state + recent state of others
- **No shared memory** between processes
- **No accurate failure detection** or local vs. communication failure distinction

### Solution: Standard Protocols
Focus on client-server model with standard protocols that handle low-level reliable communication details.

### Client-Server Model
- **Server**: Provides service (database queries, stock prices, etc.)
- **Client**: Uses service (displays results, makes recommendations)
- **Requirement**: Reliable communication (no dropped data, correct order)

### Server Types
- **File servers**: Manage disk storage and file systems
- **Database servers**: House and provide access to databases  
- **Network name servers**: Map symbolic names to IP addresses/ports

### Service Concepts
- **Service**: Set of servers of particular type
- **Binding**: Process becomes associated with specific server
- **Binding Policies**:
  - **Locality**: Look for local server first (Unix NIS)
  - **Load balance**: Distribute for uniform responsiveness (CICS)

### Data Management Strategies

#### Replication
- Multiple copies of data across locations
- Enables local access + increased availability
- Used when server processes may crash

#### Caching  
- Local copy for quick access
- **Cache hit**: Request satisfied from cache vs. primary service
- **Staleness issue**: Cached data may become outdated
- **Validation policy**: Check data freshness before use
- **Active refresh**: Cache identical to replication when actively updated
## TCP/IP Protocol Suite

### Overview
- **IP Suite**: Communication protocols for Internet and commercial networks
- **TCP**: Core protocol providing reliable, in-order data delivery
- **Protocol Stack**: Layered implementation (hardware/software mix)

### Four-Layer Architecture

#### 1. Application Layer
- Programs requiring network communication
- Examples: HTTP, FTP, Telnet
- Data passed down in application-specific format

#### 2. Transport Layer  
- **End-to-end message transfer** independent of network
- **Error control, fragmentation, flow control**
- **Two protocols**:
  - **TCP** (Connection-oriented): Reliable delivery
    - 3-packet handshake for connection establishment
    - Automatic retransmission (typically 3x)
    - Packet splitting for large data
    - Duplicate detection and reordering
  - **UDP** (Connectionless): Best-effort delivery

#### 3. Network Layer
- **Single network**: Get packets across one network
- **Internetworking**: Route across network of networks (Internet)
- **IP**: Basic packet routing from source to destination

#### 4. Link Layer
- **Physical transmission** of data
- Frame headers/trailers for physical network
- Physical component handling

### TCP Limitations
- **Disconnection timeout**: 30-90 seconds max
- **False failures**: May signal failure when endpoints are fine
- **Long outages**: Cannot overcome extended communication failures
## Remote Procedure Calls (RPC)

### Concept
- **RPC**: Extends local procedure calling across address spaces
- **Processes**: May be on same system or different systems (network-connected)
- **Similar to function call**: Arguments passed, caller waits for response

### RPC Flow
1. Client makes procedure call → sends request to server
2. Client waits for reply or timeout
3. Server calls dispatch routine → performs service → sends reply
4. Client process continues after RPC completion

### Threading Model
- **Server**: Each request spawns new thread
- **Client**: Thread issues RPC → blocks → resumes on reply

### RPC Development Process
1. **Specify protocol** for client-server communication
2. **Develop client program**
3. **Develop server program**

### Stubs
- **Generated by protocol compiler**
- **Minimal code**: Declares itself + parameters
- **Enables compilation/linking**
- **Client**: Uses stub-generated classes to execute RPC
- **Server**: Provides logic classes for handling requests

### RPC Error Cases
**New errors not present in local programming**:
- **Binding error**: Server not running when client starts
- **Version mismatch**: Client compiled against different server version
- **Timeout**: Server crash, network problem, client issue

### Error Handling Strategies
- **Unrecoverable**: Some applications treat errors as fatal
- **Fault-tolerant**: Alternate services + fail-over to backup servers

### Challenging Case: Partial Failure
**Example**: Ticket-selling server
- Client requests seat availability
- Server records sale if available
- Request times out
- **Problem**: Was seat available? Was sale recorded?
- **Risk**: Duplicate sales if retry to backup server

### Common Error Conditions
- **Network data loss**: Retransmit with "at most once" semantics
- **Server crash during operation**: Client retry after recovery
- **Server crash after completion**: Duplicate requests from client retries
- **Client crash before response**: Server discards response data
## Distributed Design Principles

### Core Philosophy
> **Ken Arnold**: "You have to design distributed systems with the expectation of failure"

### Key Principles

#### 1. **Design for Failure**
- **No assumptions** about component state
- **Anticipate failures** in all interactions
- **Example**: Don't assume second machine is ready to receive after processing

#### 2. **Explicit Failure Handling**
- **Define failure scenarios** explicitly
- **Identify likelihood** of each scenario
- **Cover most likely failures** thoroughly in code

#### 3. **Handle Unresponsive Components**
- **Both clients and servers** must deal with unresponsive counterparts
- **Timeout mechanisms** and retry logic essential

#### 4. **Minimize Network Traffic**
- **Think carefully** about data volume sent over network
- **Optimize for bandwidth** efficiency

#### 5. **Optimize for Latency**
- **Latency**: Time from request initiation to data transfer start
- **Trade-off**: Many small calls vs. one big call
- **Solution**: Experiment with small tests to find optimal compromise

#### 6. **Data Integrity**
- **Don't assume** data unchanged during transmission
- **Use checksums/validation** to verify data integrity
- **Applies to**: Network transmission, disk-to-disk transfers

#### 7. **State Management**
- **Minimize stateful components** (challenging but important)
- **State**: Information held in one place for another process
- **Cannot be reconstructed** by other components
- **If reconstructible** → it's a cache, not state

#### 8. **Caching vs Replication**
- **Caching**: Mitigates state risks but data can become stale
- **Replication**: Reduces single point of failure
- **Challenges**:
  - Consistency across replicas
  - Network partition handling
  - Coordination complexity

#### 9. **Performance Optimization**
- **Identify bottlenecks** and their causes
- **Small tests** to evaluate alternatives
- **Profile and measure** for data-driven decisions
- **Collaborate** on solution selection

#### 10. **Minimize Expensive Operations**
- **Acknowledgments**: Expensive, avoid when possible
- **Retransmissions**: Costly, tune delay parameters optimally
## Key Concepts Summary

### The Two Generals Problem
- **Scenario**: Two armies must coordinate attack timing
- **Communication**: Only via messengers through enemy territory
- **Problem**: Messengers may be captured/lost
- **Result**: **No solution possible** - common knowledge cannot be achieved through unreliable channels

### Distributed Systems Evolution
- **Past**: "Can't get work done because some unknown machine is broken"
- **Present**: Work gets done (almost always) despite failures
- **Requirements**: Available wherever, whenever, even with failures, at scale

### Why Distributed Systems?
- **Geographic separation**: 2.3B smartphone users need locality
- **Availability**: System shouldn't fail when one computer does
- **Scale**: Cycles, memory, disks, network bandwidth
- **Specialization**: Custom computers for specific tasks

### Scaling Challenges
- **End of Dennard Scaling**: Power increases with transistor density
- **Solution**: Scale out for performance
- **Reality**: All large-scale computing is distributed

### Facebook Scaling Example
- **2004**: Single server (web + database)
- **2008**: 100M users
- **2010**: 500M users  
- **2012**: 1B users
- **Evolution**: Two-tier → Three-tier → Worldwide distribution

### Data Center Reality
**Typical yearly failures**:
- ~0.5 overheating events
- ~1 PDU failure (500-1000 machines)
- ~1 rack-move (500-1000 machines)
- ~1 network rewiring
- ~20 rack failures
- ~1000 individual machine failures
- ~thousands of hard drive failures

## Study Questions

1. **Heisenbugs**: Have you encountered one? How did you isolate/fix it?
2. **Failure Types**: What makes each failure type difficult to guard against? What processing can help?
3. **8 Fallacies**: Explain why each is actually a fallacy
4. **TCP vs UDP**: Contrast protocols. When would you choose each?
5. **Caching vs Replication**: What's the difference?
6. **RPC Stubs**: What are stubs in RPC implementation?
7. **Distributed Errors**: What error conditions exist in distributed but not local environments?
8. **Pointers in RPC**: Why aren't pointers usually passed as RPC parameters?
9. **Partial Connectivity**: How would you add diagnostics for the A-B-C communication failure scenario?
10. **Leader Election**: What is it and how is it used in distributed systems?
11. **Byzantine Generals**: How can generals coordinate attack timing with unreliable messengers?

## References
- Birman, Kenneth. *Reliable Distributed Systems: Technologies, Web Services and Applications*. Springer-Verlag, 2005.
- Interview with Ken Arnold
- The Eight Fallacies
- Wikipedia: Internet Protocol Suite
- Gray, J. and Reuter, A. *Transaction Processing: Concepts and Techniques*. Morgan Kaufmann, 1993.
- Bohrbugs and Heisenbugs 

