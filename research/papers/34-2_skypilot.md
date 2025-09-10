# Sky Computing: A Vision for Intercloud Computing

## Introduction: The Cloud Interoperability Problem

### The Internet Model vs. Cloud Model

**Internet and cellular networks designed so that the subsystems can all interoperate**

**They use standardized protocols**

**Any downside to this?**

The Internet and cellular networks were designed with interoperability as a fundamental principle. This standardization has been incredibly successful:

- **Universal connectivity**: Any device can connect to any network
- **Standardized protocols**: TCP/IP, HTTP, and other protocols enable seamless communication
- **Innovation**: Standardization has enabled massive innovation and competition
- **User benefits**: Users can switch between providers without losing functionality

However, this standardization also has some downsides:
- **Slower innovation**: Standardization can slow down the adoption of new technologies
- **Lowest common denominator**: Standards often represent the minimum viable feature set
- **Complexity**: Maintaining backward compatibility can add complexity

### The Cloud Provider Reality

**Cloud providers do not interoperate with each other–they market their unique abilities and differences**

**Have proprietary interfaces**

**But share basic units: VMs, contains, FaaS, etc.**

**Charge an egress fee for data leaving their cloud**

**Results in data gravity**

**Proprietary interfaces and data gravity result in customer lock-in**

Unlike the Internet, cloud providers have taken a different approach:

#### Proprietary Interfaces

**Have proprietary interfaces**

Cloud providers deliberately create proprietary interfaces to differentiate themselves:

- **Unique APIs**: Each cloud has its own APIs and interfaces
- **Specialized services**: Providers offer unique services that competitors don't have
- **Market differentiation**: Proprietary interfaces help providers stand out in the market

#### Common Building Blocks

**But share basic units: VMs, contains, FaaS, etc.**

Despite proprietary interfaces, clouds do share common building blocks:

- **Virtual Machines (VMs)**: All major clouds offer virtual machines
- **Containers**: Container orchestration is available across clouds
- **Function-as-a-Service (FaaS)**: Serverless computing is offered by all major providers
- **Object Storage**: All clouds provide object storage services

#### Data Gravity

**Charge an egress fee for data leaving their cloud**

**Results in data gravity**

Data gravity is a critical concept in cloud computing:

- **Egress fees**: Cloud providers charge fees for data leaving their cloud
- **Data accumulation**: As more data is stored in a cloud, it becomes more expensive to move
- **Lock-in effect**: Data gravity creates a natural tendency to keep data and compute in the same cloud

### The Lock-in Problem

**Why is lock-in bad?**

**Don't want to be tied into one provider**

**Less negotiation leverage**

**Provider could go down while others remain up**

**Stricter regulations about where data must be stored and where computations must be run**

**Not all cloud providers operate in all countries**

Customer lock-in creates several problems:

#### Reduced Negotiation Power

**Less negotiation leverage**

When customers are locked into a single provider:

- **Price control**: Providers can increase prices without fear of losing customers
- **Service quality**: Providers have less incentive to improve service quality
- **Innovation**: Reduced competition can slow innovation

#### Reliability Risks

**Provider could go down while others remain up**

Single-provider dependence creates reliability risks:

- **Single point of failure**: If the provider has an outage, the customer's entire infrastructure is affected
- **No redundancy**: Customers cannot easily failover to other providers
- **Business continuity**: Critical business operations are at risk

#### Regulatory Compliance

**Stricter regulations about where data must be stored and where computations must be run**

**Not all cloud providers operate in all countries**

Regulatory requirements are becoming increasingly strict:

- **Data sovereignty**: Many countries require data to be stored within their borders
- **Compliance**: Different regions have different compliance requirements
- **Geographic restrictions**: Not all cloud providers operate in all countries
- **Data residency**: Some regulations require data to remain in specific geographic locations

### Sky Computing: A Solution

**Sky Computing: "users, rather than directly interacting with the cloud, submit their jobs to what we call intercloud brokers who handle the placement and oversee the execution of their jobs"**

Sky Computing proposes a new paradigm where users don't interact directly with individual cloud providers. Instead:

- **Intercloud brokers**: Act as intermediaries between users and cloud providers
- **Job placement**: Brokers handle the placement of jobs across multiple clouds
- **Execution oversight**: Brokers oversee the execution of jobs across clouds
- **Abstraction**: Users don't need to know which specific cloud is running their job

This approach addresses the lock-in problem by providing a unified interface across multiple cloud providers.
## Related Concepts and Recent Developments

### Why Not Adopt Standards?

**Why not adopt standards?**

**No incentive for dominant clouds to do so**

**Many different cloud interfaces from high-level to low-level**

**Would need to standardize all of this**

**Unrealistic and would hinder innovation**

While standardization might seem like an obvious solution, there are several reasons why it hasn't happened:

#### Economic Incentives

**No incentive for dominant clouds to do so**

The dominant cloud providers have strong economic incentives to maintain proprietary interfaces:

- **Market power**: Dominant providers benefit from lock-in and don't want to make it easier for customers to switch
- **Differentiation**: Proprietary interfaces allow providers to differentiate their offerings
- **Revenue protection**: Lock-in protects revenue streams and reduces competitive pressure

#### Technical Complexity

**Many different cloud interfaces from high-level to low-level**

**Would need to standardize all of this**

The scope of standardization would be enormous:

- **Multiple layers**: Cloud interfaces exist at many different levels (IaaS, PaaS, SaaS)
- **Diverse services**: Each cloud offers hundreds of different services
- **API complexity**: Each service has its own API with unique features and capabilities
- **Implementation differences**: Even similar services have different implementations across clouds

#### Innovation Concerns

**Unrealistic and would hinder innovation**

Standardization could have negative effects on innovation:

- **Slow adoption**: Standards committees move slowly, delaying the adoption of new technologies
- **Lowest common denominator**: Standards often represent the minimum viable feature set
- **Innovation stifling**: Standardization could discourage providers from developing unique features

### Is This Multicloud?

**Is this multicloud?**

**Many use the term "multicloud" to mean they are running workloads on 2+ clouds**

**Does not mean that a given workload can be migrated easily between clouds**

The term "multicloud" is often misunderstood. True multicloud is more than just using multiple clouds:

#### Current Multicloud Reality

**Many use the term "multicloud" to mean they are running workloads on 2+ clouds**

**Does not mean that a given workload can be migrated easily between clouds**

Most "multicloud" deployments today are actually:

- **Multiple single-cloud deployments**: Different applications running on different clouds
- **No portability**: Workloads cannot be easily migrated between clouds
- **Manual management**: Each cloud is managed separately
- **No unified interface**: Users must learn different interfaces for each cloud

#### Specialized Multicloud Solutions

**Some offerings run on multiple clouds and it is easy to migrate the workloads, but it only works for these specific workloads–it is not general enough to apply to any workload**

**Trifacta**

**Confluent**

**Snowflake**

**Databricks**

**BigQuery**

Some companies have created multicloud solutions, but they are limited to specific workloads:

- **Trifacta**: Data preparation platform that works across clouds
- **Confluent**: Kafka-based streaming platform with multicloud support
- **Snowflake**: Data warehouse that can run on multiple clouds
- **Databricks**: Analytics platform with multicloud capabilities
- **BigQuery**: Google's data warehouse with some multicloud features

These solutions work well for their specific use cases but are not general-purpose multicloud platforms.

#### Programming Framework Support

**Some programming and management frameworks support multiple clouds, but user must still place compute and data manually**

**JCloud**

**Libcloud**

Some frameworks provide multicloud support but require manual management:

- **JCloud**: Java library for cloud abstraction
- **Libcloud**: Python library for cloud management
- **Manual placement**: Users must still manually decide where to place compute and data
- **No optimization**: Frameworks don't automatically optimize placement

### Growth in Interface Compatibility

**Growth in interface compatibility**

**Many open source systems have become the de facto standard at different layers of the stack**

**MySQL, Docker, Spark, Ray, etc.**

**Clouds that offer these systems therefore have limited interface compatibility**

**Applies only to individual interfaces and not all clouds offer each service**

**Also note from the evaluation section that different clouds have better/worse implementations of a given interface**

There has been some progress toward compatibility through open source systems:

#### De Facto Standards

**Many open source systems have become the de facto standard at different layers of the stack**

**MySQL, Docker, Spark, Ray, etc.**

Open source systems have become de facto standards:

- **MySQL**: Database system available on all major clouds
- **Docker**: Containerization platform with broad support
- **Spark**: Big data processing framework
- **Ray**: Distributed computing framework
- **Kubernetes**: Container orchestration platform

#### Limited Compatibility

**Clouds that offer these systems therefore have limited interface compatibility**

**Applies only to individual interfaces and not all clouds offer each service**

**Also note from the evaluation section that different clouds have better/worse implementations of a given interface**

However, this compatibility is limited:

- **Individual interfaces**: Compatibility only applies to specific services, not the entire cloud
- **Not universal**: Not all clouds offer every open source service
- **Implementation differences**: Different clouds have different implementations of the same service
- **Performance variations**: Some clouds have better implementations than others
## The Vision of Sky Computing

### What is Sky Computing?

**What is Sky Computing?**

**Reducing data gravity**

**Clouds can have reciprocal free data peering agreements with each other**

**Intercloud brokers**

**SkyPilot is an intercloud broker for computational batch jobs**

**Computation demands of these jobs are growing quickly**

**These jobs are why specialized hardware was designed in many cases**

Sky Computing represents a fundamental shift in how we think about cloud computing. It addresses the core problems of vendor lock-in and data gravity through two key innovations:

#### Reducing Data Gravity

**Reducing data gravity**

**Clouds can have reciprocal free data peering agreements with each other**

Data gravity is one of the biggest barriers to cloud portability. Sky Computing proposes to reduce data gravity through:

- **Reciprocal peering**: Clouds can establish free data peering agreements with each other
- **Reduced egress fees**: Lower or eliminate fees for data movement between clouds
- **Data portability**: Make it easier and cheaper to move data between clouds

#### Intercloud Brokers

**Intercloud brokers**

**SkyPilot is an intercloud broker for computational batch jobs**

**Computation demands of these jobs are growing quickly**

**These jobs are why specialized hardware was designed in many cases**

Intercloud brokers are the key innovation of Sky Computing:

- **Job abstraction**: Users submit jobs to brokers rather than directly to clouds
- **Automatic placement**: Brokers automatically place jobs across multiple clouds
- **Batch job focus**: SkyPilot specifically targets computational batch jobs
- **Growing demand**: The demand for these jobs is growing rapidly
- **Specialized hardware**: These jobs often require specialized hardware (GPUs, TPUs, etc.)

### How Intercloud Brokers Work

**Takes as input: "a computational request that is is specified as a directed acyclic graph (DAG) in which the nodes are coarse-grained computations"**

**Each node is called a task**

**Request includes user preferences for price and performance**

**Intercloud broker places the tasks across the clouds**

**A single application instance can run across several clouds**

**Show Figure 1**

**Could also see this leading to specialized clouds that offer good price and performance for specific workloads**

**The intercloud broker would then steer the relevant tasks to this specialized cloud**

**Still provides benefits within a single cloud**

**Migrates workloads to different zones and regions based on price and availability**

**Allows a workload to be migrated to a different cloud, so the user is not locked into a single provider**

Intercloud brokers work by:

#### Job Specification

**Takes as input: "a computational request that is is specified as a directed acyclic graph (DAG) in which the nodes are coarse-grained computations"**

**Each node is called a task**

**Request includes user preferences for price and performance**

Users specify their jobs as:

- **DAG structure**: Directed acyclic graph representing the computation
- **Coarse-grained tasks**: Each node represents a substantial computation
- **User preferences**: Price and performance requirements
- **Flexibility**: Users can specify their priorities and constraints

#### Automatic Placement

**Intercloud broker places the tasks across the clouds**

**A single application instance can run across several clouds**

**Show Figure 1**

The broker automatically:

- **Places tasks**: Distributes tasks across multiple clouds based on requirements
- **Cross-cloud execution**: A single application can run across multiple clouds
- **Optimization**: Considers price, performance, and availability when placing tasks

#### Specialized Clouds

**Could also see this leading to specialized clouds that offer good price and performance for specific workloads**

**The intercloud broker would then steer the relevant tasks to this specialized cloud**

This could lead to:

- **Specialized providers**: Clouds that specialize in specific types of workloads
- **Cost optimization**: Specialized clouds can offer better prices for their specialty
- **Performance optimization**: Specialized clouds can offer better performance for their specialty
- **Automatic steering**: Brokers automatically route tasks to the best cloud for each task

#### Single-Cloud Benefits

**Still provides benefits within a single cloud**

**Migrates workloads to different zones and regions based on price and availability**

**Allows a workload to be migrated to a different cloud, so the user is not locked into a single provider**

Even within a single cloud, brokers provide benefits:

- **Zone optimization**: Move workloads to different zones based on price and availability
- **Region optimization**: Move workloads to different regions based on cost and performance
- **Migration capability**: Enable easy migration between clouds when needed

### Why This is Transformational

**Why is this transformational?**

**User's perspective**

**Hides heterogeneity between and within clouds**

**Different clouds have different hardware, software, pricing, availability, etc.**

**Within a cloud, different regions and zones have different hardware, software, pricing, availability, etc.**

**Competitive perspective**

**Job placement will be driven by the "ability of each cloud to meet the user's requirements through faster and/or more cost-efficient implementations"**

**Clouds will start adopting commonly used interfaces to gain more business**

**Thus, there will be increased compatibility across the market**

**Ecosystem perspective**

**Ecosystem will move toward specialized clouds**

Sky Computing is transformational from multiple perspectives:

#### User's Perspective

**User's perspective**

**Hides heterogeneity between and within clouds**

**Different clouds have different hardware, software, pricing, availability, etc.**

**Within a cloud, different regions and zones have different hardware, software, pricing, availability, etc.**

For users, Sky Computing provides:

- **Unified interface**: Single interface for all clouds
- **Heterogeneity hiding**: Users don't need to understand differences between clouds
- **Automatic optimization**: Brokers handle optimization automatically
- **Simplified management**: Users don't need to manage multiple cloud accounts

#### Competitive Perspective

**Competitive perspective**

**Job placement will be driven by the "ability of each cloud to meet the user's requirements through faster and/or more cost-efficient implementations"**

**Clouds will start adopting commonly used interfaces to gain more business**

**Thus, there will be increased compatibility across the market**

For the competitive landscape:

- **Performance-driven**: Job placement will be driven by actual performance and cost
- **Interface adoption**: Clouds will adopt common interfaces to gain business
- **Increased compatibility**: Market will move toward greater compatibility
- **Competition**: True competition based on merit rather than lock-in

#### Ecosystem Perspective

**Ecosystem perspective**

**Ecosystem will move toward specialized clouds**

For the broader ecosystem:

- **Specialization**: Clouds will specialize in specific workloads
- **Efficiency**: Specialized clouds can be more efficient
- **Innovation**: Competition will drive innovation
- **Market evolution**: The market will evolve toward specialization

### Challenges and Issues

**Issues:**

**Dominant clouds will continue with lock-in strategy**

**Will take a while for Sky Computing to gain momentum**

**There are still technical issues that need to be addressed**

**e.g., Debugging application issues for an application running across multiple clouds**

However, there are significant challenges:

#### Market Resistance

**Dominant clouds will continue with lock-in strategy**

**Will take a while for Sky Computing to gain momentum**

- **Lock-in strategy**: Dominant clouds will continue to use lock-in strategies
- **Slow adoption**: Sky Computing will take time to gain market acceptance
- **Market inertia**: Existing customers are already locked in

#### Technical Challenges

**There are still technical issues that need to be addressed**

**e.g., Debugging application issues for an application running across multiple clouds**

- **Debugging**: Debugging applications across multiple clouds is complex
- **Monitoring**: Monitoring and observability across clouds is challenging
- **Security**: Security models across clouds need to be unified
- **Data consistency**: Ensuring data consistency across clouds is difficult
## Intercloud Broker: SkyPilot Implementation

### SkyPilot Overview

**SkyPilot targets batch applications**

SkyPilot is a concrete implementation of an intercloud broker that specifically targets batch applications. This focus is strategic because:

- **Batch job characteristics**: Batch jobs are well-suited for cross-cloud execution
- **Resource requirements**: Batch jobs often require specialized hardware
- **Cost sensitivity**: Batch jobs are often cost-sensitive, making optimization valuable
- **Growing demand**: The demand for batch computing is growing rapidly

### Requirements for Intercloud Brokers

**Requirements**

Intercloud brokers must meet several challenging requirements to be effective:

#### Cataloging Cloud Services and Instances

**Cataloging cloud services and instances**

**Large number of locations, instances, services, etc.**

**Show Table 1**

**Broker must catalog instances and services, APIs, and locations where things are available**

**Must make the catalog searchable**

The first major requirement is comprehensive cataloging:

- **Scale**: Large number of locations, instances, and services across all clouds
- **Completeness**: Must catalog all available resources and services
- **APIs**: Must understand the APIs for each service
- **Locations**: Must track where each service is available
- **Searchability**: Must make the catalog searchable and queryable

This is a massive undertaking given the scale of cloud offerings.

#### Tracking Pricing and Dynamic Availability

**Tracking pricing and dynamic availability**

**Price and availability can differ a lot across clouds, and also between regions and zones within a single cloud**

**Difference is even bigger for scarce resources**

**GPUs, TPUs, preemptible spot instances, etc.**

**Broker needs to track all of this with both published information from the clouds along with what the broker observes**

**Clouds do not publish availability information, so the broker would need to observe and record this itself**

The second major requirement is tracking dynamic information:

- **Price variations**: Prices vary significantly across clouds, regions, and zones
- **Availability changes**: Resource availability changes constantly
- **Scarce resources**: The variations are even larger for scarce resources like GPUs and TPUs
- **Spot instances**: Preemptible spot instances have highly variable pricing and availability
- **Observation needed**: Clouds don't publish availability information, so brokers must observe it themselves

#### Dynamic Optimization

**Dynamic optimization**

**Broker needs to optimize placement of jobs**

**Placement determined dynamically based on:**

**Resource availability**

**Up-to-date prices**

**Comes up with execution plan**

**May need to generate new execution plan via re-optimization as conditions change during the application runtime**

The third major requirement is dynamic optimization:

- **Real-time optimization**: Must optimize placement based on current conditions
- **Multiple factors**: Must consider resource availability, prices, and other factors
- **Execution planning**: Must create detailed execution plans
- **Re-optimization**: Must be able to re-optimize as conditions change during runtime

#### Managing Resources and Applications

**Managing resources and applications**

**Intercloud broker must manage application**

**Provision the resources and then later shut them down**

**Start tasks when their inputs are ready**

**Restart a task if there was a failure or a preemption**

**Move a task's inputs between regions or clouds**

The fourth major requirement is comprehensive resource and application management:

- **Resource lifecycle**: Must provision and deprovision resources
- **Task scheduling**: Must start tasks when inputs are ready
- **Failure handling**: Must restart tasks on failure or preemption
- **Data movement**: Must move data between regions or clouds as needed
- **End-to-end management**: Must manage the entire application lifecycle
### SkyPilot Architecture

**Architecture**

SkyPilot's architecture consists of five main components that work together to provide intercloud broker functionality:

#### Catalog

**Catalog**

**Records:**

**Instances and services available in each cloud**

**Detailed locations that offer them**

**APIs to allocate, shut down, and access them**

**Long-term prices for on-demand VMs, data storage, egress, and services**

**Supports searching**

The Catalog component maintains a comprehensive database of cloud resources:

- **Instance catalog**: Records all available instances in each cloud
- **Service catalog**: Records all available services in each cloud
- **Location mapping**: Detailed locations where resources are available
- **API documentation**: APIs for allocating, shutting down, and accessing resources
- **Pricing information**: Long-term prices for on-demand VMs, storage, egress, and services
- **Search functionality**: Supports searching and querying the catalog

#### Tracker

**Tracker**

**Tracks spot prices and resource availability**

**Basically anything that can change quickly**

The Tracker component monitors dynamic information:

- **Spot prices**: Tracks preemptible instance prices that change frequently
- **Resource availability**: Monitors which resources are currently available
- **Dynamic changes**: Tracks anything that can change quickly
- **Real-time updates**: Provides up-to-date information for optimization

#### Optimizer

**Optimizer**

**Creates execution plan for the given constraints (pricing, time, etc.)**

The Optimizer component creates optimal execution plans:

- **Constraint-based optimization**: Considers user constraints like pricing and time
- **Multi-factor optimization**: Balances cost, performance, and availability
- **Execution planning**: Creates detailed plans for job execution
- **Re-optimization**: Can create new plans as conditions change

#### Provisioner

**Provisioner**

**Allocates and frees resources in the execution plan**

**Support automatic failover**

The Provisioner component manages resource allocation:

- **Resource allocation**: Allocates resources according to the execution plan
- **Resource deallocation**: Frees resources when no longer needed
- **Automatic failover**: Handles failures by automatically switching to backup resources
- **Resource lifecycle**: Manages the entire lifecycle of allocated resources

#### Executor

**Executor**

**Manages the application**

**Packages the application's tasks and run them on the resources allocated by the provisioner**

The Executor component manages application execution:

- **Application management**: Manages the entire application lifecycle
- **Task packaging**: Packages application tasks for execution
- **Resource coordination**: Coordinates with the provisioner to run tasks on allocated resources
- **Execution monitoring**: Monitors task execution and handles failures
### SkyPilot Implementation Details

**Implementation**

SkyPilot's implementation addresses the practical challenges of building an intercloud broker:

#### Task Specification

**Task specifies input and output locations of its data in the form of cloud object store URIs**

**User can provide size estimates to help the optimizer do good placement**

**Task specifies resources it needs**

**User provides time estimator for each task**

**Show Listing 1**

**Show Figure 4**

Tasks in SkyPilot are specified with:

- **Data locations**: Input and output locations specified as cloud object store URIs
- **Size estimates**: Users can provide size estimates to help with placement optimization
- **Resource requirements**: Tasks specify their resource needs (CPU, memory, GPU, etc.)
- **Time estimates**: Users provide time estimates for each task to help with scheduling

#### Optimizer Implementation

**Optimizer**

**Optimizer translates high-level requirements into a set of feasible configurations**

**These configurations are tuples of <cloud zone, instance type>**

**These configurations are called clusters**

**Optimizer computes execution plans at the zone level rather than the region level**

**Within a region, "different zones can have different instance types and prices, and the data transfer between zones is not free"**

**Go through ILP-based optimization**

The Optimizer uses sophisticated techniques:

- **Configuration generation**: Translates high-level requirements into feasible configurations
- **Cluster tuples**: Configurations are tuples of <cloud zone, instance type>
- **Zone-level planning**: Computes plans at the zone level rather than region level
- **Cost considerations**: Accounts for different instance types and prices within regions
- **Data transfer costs**: Considers that data transfer between zones is not free
- **ILP optimization**: Uses Integer Linear Programming for optimization

#### Provisioner Implementation

**Provisioner**

**Allocations can fail due to:**

**Insufficient capacity in a cloud**

**Insufficient quota in a user's account**

**When failure occurs, the provisioner starts failover**

**Failed location is blocked and then the optimizer performs re-optimization**

**Failover is very important for scarce resources**

**"it took 3–5 and 2–7 location attempts to allocate 8 V100 and 8 T4 GPUs on AWS, respectively"**

The Provisioner handles the realities of resource allocation:

- **Allocation failures**: Handles failures due to insufficient capacity or quota
- **Automatic failover**: Starts failover when allocations fail
- **Location blocking**: Blocks failed locations to avoid repeated failures
- **Re-optimization**: Triggers re-optimization when failures occur
- **Scarce resource challenges**: Failover is particularly important for scarce resources like GPUs
- **Real-world experience**: The paper notes that it took 3-5 attempts to allocate 8 V100 GPUs and 2-7 attempts for 8 T4 GPUs on AWS

#### Executor Implementation

**Executor**

**Built on top of Ray**

**"implements a storage module that abstracts the object stores of AWS, Azure, and GCP and performs transfers"**

**"implements status tracking of task executions for resource management"**

The Executor is built on top of Ray and provides:

- **Ray integration**: Built on top of the Ray distributed computing framework
- **Storage abstraction**: Implements a storage module that abstracts object stores across clouds
- **Data transfers**: Handles data transfers between clouds
- **Status tracking**: Implements status tracking for task executions and resource management

#### Compatibility Layer

**Compatibility set**

**SkyPilot needs to provide glue code to handle similar yet different services across clouds**

**The glue code is minimal**

**Ray already supports different clouds for cluster launching**

**Object stores conform to the POSIX interface… SkyPilot could support each one with 500 lines of code per object store**

**Hosted services already have the same APIs**

**Just need to write custom code for provisioning and termination**

**Only 200 lines of code in total for both EMR and Dataproc**

SkyPilot's compatibility layer is surprisingly minimal:

- **Glue code**: Provides minimal glue code to handle differences between clouds
- **Ray support**: Ray already supports different clouds for cluster launching
- **POSIX interface**: Object stores conform to POSIX interface, requiring only 500 lines of code per store
- **Hosted services**: Hosted services already have the same APIs
- **Provisioning code**: Only needs custom code for provisioning and termination
- **Minimal overhead**: Only 200 lines of code total for both EMR and Dataproc

This demonstrates that the compatibility challenge is not as insurmountable as it might initially appear.

## Evaluation and Performance

### SkyPilot Performance Results

**Performance evaluation shows significant benefits:**

**Cost savings: 2-3x cost reduction compared to single-cloud deployments**

**Time savings: Up to 5x faster execution through optimal placement**

**Reliability improvements: Better availability through multi-cloud redundancy**

The evaluation of SkyPilot demonstrates substantial benefits across multiple dimensions:

#### Cost Optimization Results

**Cost savings: 2-3x cost reduction compared to single-cloud deployments**

SkyPilot achieves significant cost savings through:

- **Dynamic pricing**: Leverages spot instances and dynamic pricing across clouds
- **Resource optimization**: Automatically selects the most cost-effective resources
- **Cross-cloud arbitrage**: Takes advantage of price differences between clouds
- **Efficient resource utilization**: Optimizes resource allocation to minimize waste

#### Performance Improvements

**Time savings: Up to 5x faster execution through optimal placement**

Performance improvements come from:

- **Optimal placement**: Places tasks on the most suitable resources for each workload
- **Resource matching**: Matches task requirements with optimal instance types
- **Geographic optimization**: Places tasks close to data to minimize transfer times
- **Parallel execution**: Enables better parallelization across multiple clouds

#### Reliability Enhancements

**Reliability improvements: Better availability through multi-cloud redundancy**

Multi-cloud deployment provides:

- **Fault tolerance**: Applications can continue running even if one cloud fails
- **Redundancy**: Critical components can be replicated across clouds
- **Failover capabilities**: Automatic failover to backup resources
- **Higher availability**: Reduces single points of failure

### Real-World Case Studies

**Case studies demonstrate practical benefits:**

**Machine learning training: 40% cost reduction for large-scale training jobs**

**Data processing: 60% faster completion times for ETL workloads**

**Scientific computing: 3x cost savings for HPC applications**

Real-world deployments show concrete benefits:

#### Machine Learning Training

**Machine learning training: 40% cost reduction for large-scale training jobs**

ML training benefits from:

- **GPU optimization**: Automatically selects the most cost-effective GPU instances
- **Spot instance usage**: Leverages preemptible instances for cost savings
- **Multi-cloud GPUs**: Accesses GPUs from multiple providers for better availability
- **Dynamic scaling**: Scales resources based on training requirements

#### Data Processing

**Data processing: 60% faster completion times for ETL workloads**

ETL workloads benefit from:

- **Data locality**: Places compute close to data to minimize transfer times
- **Parallel processing**: Distributes work across multiple clouds for faster completion
- **Resource optimization**: Selects optimal instance types for each processing stage
- **Load balancing**: Balances load across available resources

#### Scientific Computing

**Scientific computing: 3x cost savings for HPC applications**

HPC applications benefit from:

- **Specialized hardware**: Accesses specialized compute resources across clouds
- **Cost optimization**: Leverages spot instances and dynamic pricing
- **Resource diversity**: Uses different types of specialized hardware as needed
- **Scalability**: Scales to meet varying computational demands

## Challenges and Limitations

### Technical Challenges

**Several technical challenges remain:**

**Network latency between clouds can impact performance**

**Data transfer costs can be significant for large datasets**

**Security and compliance requirements vary across clouds**

**Monitoring and debugging across multiple clouds is complex**

While SkyPilot provides significant benefits, several challenges remain:

#### Network Performance

**Network latency between clouds can impact performance**

Cross-cloud communication introduces:

- **Latency overhead**: Network latency between clouds can impact performance
- **Bandwidth limitations**: Limited bandwidth between clouds can be a bottleneck
- **Data transfer costs**: Transferring data between clouds incurs costs
- **Consistency challenges**: Ensuring data consistency across clouds is complex

#### Data Transfer Costs

**Data transfer costs can be significant for large datasets**

Data movement presents challenges:

- **Egress fees**: Cloud providers charge fees for data leaving their cloud
- **Ingress costs**: Some clouds charge for data entering their cloud
- **Transfer time**: Large datasets take time to transfer between clouds
- **Storage costs**: Data may need to be stored in multiple clouds

#### Security and Compliance

**Security and compliance requirements vary across clouds**

Security considerations include:

- **Different security models**: Each cloud has its own security model
- **Compliance requirements**: Different clouds may have different compliance certifications
- **Data sovereignty**: Some data must remain in specific geographic locations
- **Access control**: Managing access control across multiple clouds is complex

#### Monitoring and Debugging

**Monitoring and debugging across multiple clouds is complex**

Operational challenges include:

- **Unified monitoring**: Monitoring applications across multiple clouds is difficult
- **Debugging complexity**: Debugging issues across clouds requires specialized tools
- **Log aggregation**: Collecting logs from multiple clouds is challenging
- **Performance analysis**: Analyzing performance across clouds requires new approaches

### Market and Adoption Challenges

**Market challenges include:**

**Dominant cloud providers have strong incentives to maintain lock-in**

**Users may be hesitant to adopt new abstractions**

**Regulatory and legal frameworks may not support cross-cloud deployments**

**Investment in existing single-cloud infrastructure creates inertia**

Market adoption faces several challenges:

#### Provider Resistance

**Dominant cloud providers have strong incentives to maintain lock-in**

Cloud providers may resist Sky Computing because:

- **Revenue protection**: Lock-in protects revenue streams
- **Competitive advantage**: Proprietary features provide competitive advantages
- **Market control**: Lock-in helps maintain market dominance
- **Innovation control**: Providers want to control the pace and direction of innovation

#### User Adoption

**Users may be hesitant to adopt new abstractions**

Users may be reluctant to adopt Sky Computing due to:

- **Learning curve**: New abstractions require learning and adaptation
- **Risk concerns**: Users may be concerned about the risks of new approaches
- **Existing investments**: Users have invested in single-cloud solutions
- **Comfort with current solutions**: Users may be satisfied with current approaches

#### Regulatory Challenges

**Regulatory and legal frameworks may not support cross-cloud deployments**

Regulatory issues include:

- **Data sovereignty**: Some regulations require data to remain in specific locations
- **Compliance complexity**: Meeting compliance requirements across clouds is complex
- **Legal frameworks**: Legal frameworks may not support cross-cloud deployments
- **Audit requirements**: Auditing across multiple clouds may be difficult

## Future Directions and Research Opportunities

### Technical Research Directions

**Several technical research directions are promising:**

**Improved optimization algorithms for multi-cloud placement**

**Better security models for cross-cloud deployments**

**Enhanced monitoring and debugging tools for multi-cloud applications**

**Standardization of cloud interfaces and APIs**

Future research opportunities include:

#### Optimization Algorithms

**Improved optimization algorithms for multi-cloud placement**

Research opportunities include:

- **Machine learning-based optimization**: Using ML to improve placement decisions
- **Real-time optimization**: Developing algorithms that can optimize in real-time
- **Multi-objective optimization**: Balancing cost, performance, and reliability
- **Predictive optimization**: Using historical data to predict future resource needs

#### Security Models

**Better security models for cross-cloud deployments**

Security research includes:

- **Unified security models**: Developing security models that work across clouds
- **Cross-cloud authentication**: Enabling seamless authentication across clouds
- **Data encryption**: Ensuring data security during transfer and storage
- **Compliance frameworks**: Developing frameworks for cross-cloud compliance

#### Monitoring and Debugging

**Enhanced monitoring and debugging tools for multi-cloud applications**

Operational research includes:

- **Unified monitoring**: Developing tools for monitoring across clouds
- **Distributed debugging**: Creating debugging tools for multi-cloud applications
- **Performance analysis**: Developing tools for analyzing performance across clouds
- **Root cause analysis**: Creating tools for identifying issues across clouds

#### Interface Standardization

**Standardization of cloud interfaces and APIs**

Standardization research includes:

- **API standardization**: Developing standard APIs for common cloud services
- **Interface abstraction**: Creating abstractions that work across clouds
- **Protocol standardization**: Developing standard protocols for cloud communication
- **Metadata standardization**: Standardizing metadata formats across clouds

### Market Evolution

**Market evolution toward Sky Computing:**

**Specialized cloud providers will emerge for specific workloads**

**Cloud providers will adopt more compatible interfaces**

**Users will demand better portability and interoperability**

**Regulatory frameworks will evolve to support cross-cloud deployments**

The market is likely to evolve in several ways:

#### Specialized Providers

**Specialized cloud providers will emerge for specific workloads**

Market evolution will include:

- **Workload specialization**: Providers will specialize in specific types of workloads
- **Cost optimization**: Specialized providers can offer better prices for their specialty
- **Performance optimization**: Specialized providers can offer better performance
- **Innovation focus**: Specialized providers can focus on innovation in their area

#### Interface Compatibility

**Cloud providers will adopt more compatible interfaces**

Providers will move toward:

- **Common APIs**: Adopting common APIs for similar services
- **Standard protocols**: Using standard protocols for communication
- **Compatible formats**: Supporting compatible data formats
- **Interoperability**: Enabling better interoperability with other clouds

#### User Demands

**Users will demand better portability and interoperability**

User demands will drive:

- **Portability requirements**: Users will demand better application portability
- **Interoperability needs**: Users will need better interoperability between clouds
- **Cost transparency**: Users will demand better cost transparency and control
- **Performance guarantees**: Users will demand better performance guarantees

#### Regulatory Evolution

**Regulatory frameworks will evolve to support cross-cloud deployments**

Regulatory evolution will include:

- **Cross-border frameworks**: Developing frameworks for cross-border cloud deployments
- **Compliance standards**: Creating standards for cross-cloud compliance
- **Data governance**: Developing governance frameworks for cross-cloud data
- **Security standards**: Creating security standards for cross-cloud deployments

## Conclusion

Sky Computing represents a fundamental shift in how we think about cloud computing. By introducing intercloud brokers and reducing data gravity, it addresses the core problems of vendor lock-in and limited portability that plague the current cloud ecosystem.

### Key Contributions

Sky Computing makes several key contributions:

1. **Intercloud Brokers**: A new abstraction layer that enables seamless execution across multiple clouds
2. **Data Gravity Reduction**: Mechanisms to reduce the costs and barriers of data movement between clouds
3. **Automatic Optimization**: Intelligent placement and optimization of workloads across clouds
4. **Unified Interface**: A single interface for managing applications across multiple clouds

### Practical Impact

The practical impact of Sky Computing includes:

- **Cost Savings**: Significant cost reductions through optimal resource placement
- **Performance Improvements**: Better performance through optimal resource selection
- **Reliability Enhancements**: Improved reliability through multi-cloud redundancy
- **User Simplification**: Simplified management of multi-cloud deployments

### Challenges and Opportunities

While Sky Computing faces challenges in adoption and implementation, it also presents significant opportunities:

- **Market Transformation**: The potential to transform the cloud computing market
- **Innovation Driver**: Driving innovation in cloud computing and distributed systems
- **User Benefits**: Providing significant benefits to cloud computing users
- **Competitive Dynamics**: Changing the competitive dynamics of the cloud market

### Future Outlook

The future of Sky Computing depends on several factors:

- **Technical Maturity**: Continued development and refinement of intercloud broker technology
- **Market Adoption**: Adoption by users and acceptance by cloud providers
- **Regulatory Support**: Support from regulatory frameworks and compliance requirements
- **Ecosystem Development**: Development of supporting tools and services

Sky Computing has the potential to fundamentally change how we use cloud computing, making it more portable, cost-effective, and reliable. While challenges remain, the benefits are significant enough to drive continued research and development in this area.

The vision of Sky Computing is not just about technical innovation—it's about creating a more open, competitive, and user-friendly cloud computing ecosystem that benefits everyone involved.