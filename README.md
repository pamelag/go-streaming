# go-streaming
A commit log service to handle large amount of small event data coming at very high speed in Go

## Background
RightPrism project's frontend consists of a set of browser-based editors for authoring software requirements, creating wireframes etc. These editors keep sending the delta changes of the edits to these documents continuously to the server. These differentials are json formatted string that are a mix of instructions and data, emitted from the editors as authors type text or resize a shape or drag-drop a stencil.

This project is used to aggregate this stream of delta events. The nature of this data generated from these editors is an ever-growing, infinite data set, also known as unbounded data. As I build this project further, this stream of data would be processed, so that the system can gain certain insights and provide helpful suggestions to the author, who's currently working on the editor.

## Goals of this project:
- Implement a commit log service in Go
- Ingest event stream data using gRPC and Protocol Buffer
- Collect and move log data to data store
