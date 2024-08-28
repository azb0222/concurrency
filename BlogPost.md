# Concurrency in Go  

This post and my package is based off of the book *Concurrency in Go*. I fucking loved it. 

## Sync Package: 
//TODO: summary

## Goroutines and Channels
//TODO: a quick summary. show blocking functionality. use for select and done channel 

## Key Principles
- concurrency != parallelism 
  - **concurrency**: a property of the code 
  - **parallelism**: a property of the runtime environment
- goroutines are CHEAP! //TODO: show demo 
- Each channel's lifecycle should be maintained end to end by a parent goroutine
  - Every time a pipeline stage is restarted within a new goroutine, a new channel will be created (//TODO: move to pipelines section?)
- //TODO:confinement 
- if a parent goroutine creates a child goroutine, it should also be able to stop the child goroutine (via done channel)
- errors are first class citizens (//TODO: idk if this is the right word): errors produced in a goroutine should be propogated to another part of the program with more detailed error handling to maintain seperation of concerns

## Concurrency Helper Package 
In my opinion, the logic behind the concurrency design patterns should be decoupled from your core program's logic. 

Therefore, in order to use in my everyday Go coding from now on, I created a helper package with abstractions for every design pattern mentioned. Again, all credit for the core logic in these design patterns goes towards the book. 

The package also has helper functions for error handling. 

### Example Usage 

//TODO: include explanation for each design pattern, not just code example 
#### Pipelines

#### Queues

## Conclusion 
This package is still a work in progress. 

I am very bad at coding, so please message me if there may be a better to structure some of this. Would very much appreciate it!!

