## Research

This folder contains case studies, papers, and analyses of real-world distributed systems. It complements `reference/` (theory and concepts) and `app/` (implementations) by connecting ideas to influential systems and empirical lessons.

### Purpose

- Study seminal and modern systems papers
- Extract design principles and trade-offs from real deployments
- Compare alternative approaches to replication, consensus, storage, and scaling
- Bridge theory to practice with concrete examples

### How to Use This Folder

1. Read a paper or case study end-to-end.
2. Summarize goals, assumptions, design, and evaluation in your own words.
3. Identify trade-offs and failure modes; relate them to concepts in `reference/`.
4. Recreate simplified components or experiments in `app/` when possible.

### Content Types

- System design papers and conference publications
- Industry case studies and postmortems
- Comparative analyses and survey notes
- Reproductions and simplified re-implementations

### Suggested Study Path

1. Start with classic consensus and replication systems
2. Move to storage architectures and distributed databases
3. Explore coordination, scheduling, and cloud-scale systems
4. Study operational lessons and debugging in production

### Case Studies (examples)

- Paxos: The Part-Time Parliament — Leslie Lamport  
  https://lamport.azurewebsites.net/pubs/lamport-paxos.pdf

(Add additional papers here, e.g., Raft, Dynamo, Spanner, GFS, Bigtable, Chubby, Zookeeper, Kafka, Kubernetes, etc.)

### Cross-Links

- Concepts and theory: see `../reference`
- Implementations and exercises: see `../app`

### External Resources

- MIT 6.5840 Distributed Systems: https://pdos.csail.mit.edu/6.824/
- Systems papers collections (e.g., curated lists by universities and labs)
- Reproduction guides or labs associated with the above papers