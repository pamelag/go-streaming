# go-streaming
A simple streaming engine written in Go

The modern programming landscape has exploded with big data frameworks and technologies. But having the ability to process large quantities of data in bulk (batch processing) isn’t always enough. Increasingly, organizations are finding that they need to process data as it becomes available (stream processing). 


## Streaming fundamentals

### DataFlow Programming and DataFlow Graphs

Dataflow programs describe how data flows between operations. these programs are a form of directed graph in whichh data flows in certain 
diections. In this graph nodes are called operators and they represent computations and edges represent data dependencies. Operators are basic functional units of dataflow application. They consume data from inputs, perform a computation on them, and produce data to outputs for further processing. 

### Data parallelism and task parallelism

In order to achieve parallelism data can be partitioned and have tasks of the same operation execute on the data subsets in parallel. This is called data parallelism. Second, you can have tasks from different operators performing computations on the same or different data in parallel. This type of parallelism is called task parallelism.


### Latency and Throughput

Latency indicates how long it takes for an event to be processed. It is 
the time interval between receiving an event and seeing the effect of the processing in the output.

Throughput is a measure of the system’s processing capacity. This is the rate of processing.

Latency and Throughput are interdependent. If we have high latency throughput decreases. So the challenge is to achieve both high latence and throughput.


### Operations on Data Streams
Operations on the Data Stream can be either stateful or stateless. Stateful operations maintain internal state.

Statleless operations are easy to parallelize. In case of failure the 
stateless operations can be simply restarted and continue processingfrom where they left off.

Stateful stream processing is more difficult to handle in terms of failure scenarios and parallelize.

### Data Ingestion and Data Egress
Data ingestion is the operation of fetching raw data from external sources and converting it into a format that is suitable for processing. Operators that implement data ingestion logic are called data sources.

### Transformation Operations
These operations consume one event after the other and apply some transformation to the event data, producing a new output stream.

### Rolling Aggregation
A rolling aggregation is a continiously running process on the Data Stream and it calculates some kind of aggregate like Sum, Minimum, Maximum etc.
Aggregation operations are stateful.

### Windows
Transforms and rolling aggregations process an event at a time to produce the output. Operations that require a holistic aggregate must collect buffer records to compute their results. These kind of operations
is interested on a data coming in a particular time window.

Window operations continuously create finite sets of events called buckets from an unbounded event stream and let us perform computations on these finite sets. Events are usually assigned to buckets based on data properties or based on time.

#### Common window types
Tumbling windows assign events into non-overlapping buckets of fixed size. When the window border is passed, all the events are sent to an evaluation function for processing. Count-based tumbling windows define how many events are collected before triggering evaluation. 

#### Sliding
Sliding windows assign events into overlapping buckets of fixed size. Thus, an event might belong to multiple buckets.

#### Session
Session windows are useful in a common real-world scenario where neither tumbling nor sliding windows can be applied. Consider an application that analyzes online user behavior. In such applications, we would like to group together events that origin from the same period of user activity or session. Sessions comprise of a series of events happening in adjacent times followed by a period of inactivity. For example, user interactions with a series of news articles one after the other could be considered a session. Since the length of a session is not defined beforehand but depends on the actual data, tumbling and sliding windows cannot be applied in this scenario. Instead, we need a window operation that assigns events belonging to the same session in the same bucket. Session windows group events in session based on a session gap value that defines the time of inactivity to consider a session closed. 

### Event Time
Event-time is the time when an event in the stream actually happened. Event time is based on a timestamp that is attached on the events of the stream. Timestamps usually exist inside the event data before they enter the processing pipeline (e.g. event creation time).

### Watermark
how long do we have to wait before we can be certain that we have received all events that happened before a certain point of time?
A watermark is a global progress metric that indicates a certain point in time when we are confident that no more delayed events will arrive. In essence, watermarks provide a logical clock which informs the system about the current event time. When an operator receives a watermark with time T, it can assume that no further events with timestamp less than T will be received.

### State and consistency model
A UDF accumulates state over a period or number of events, e.g. to compute an aggregation or detect a pattern. Stateful operators use both incoming events and internal state to compute their output. n streaming we have durable state across events and we can expose it as a first-class citizen in the programming model. 


## Streaming System Architecture
A streaming system can consist of multiple processes that run distributed across multiple machines. Common challenges that distributed systems need to address are allocation and management of compute resources in a cluster, process coordination, durable and available data storage, and failure recovery.

### Components
A stream processing application setup consists of 4 different components to execute streaming application.

### Timestamps

### Watermarks

### Steam Execution Environment
The execution environment determines whether the program is running on a local machine or on a cluster. The stream execution environment provides methods to create a stream source that ingests data streams into the application. Data streams can be ingested from sources such as message queues, files, or also be generated on the fly.











