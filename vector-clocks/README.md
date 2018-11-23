# Vector Clocks

Vector clocks are used to compare the causal ordering of events across nodes in a distributed system, where a physical clock is not reliable.   
Each node maintains a vector of counters (a vector clock), which has an entry for each node in the distributed system. When a node performs a local event, it updates its own counter in the vector clock.   
When a message is sent between nodes *a* and *b*, the sending node *a* increments its own entry in its vector clock and then sends its message along with its vector clock to *b*.  
Upon receiving the message, node *b* then increments its own entry in its vector clock, and updates its vector clock at every other entry to be the max of its entry and the entry in *a*'s clock. 
Further details of this algorithm are described in [Chapter 3](http://book.mixu.net/distsys/time.html) of *Distributed Systems - for fun and profit*.

## Implementation
The vector clock simulation is implemented in `vc.go`. Each node in the distributed system is modelled as a goroutine called `node` that communicates with other "nodes" through channels. Each node maintains its own vector clock (an array of integers) and updates it as events occur.   
The program reads input from `stdin` to determine the ordering of events at each node as well as when messages should be sent between nodes. The constant, `numProcesses` indicates the number of nodes that are running concurrently in the program. This is currently set to 3 but can be changed if we wish to test with more nodes.  

### Input 
A sample input file, `in.txt` is provided. The input file should contain a line for each node, which indicates the number, order, and type of events that will be run at that node. Events are separated by a space.  
For example, the first line in `in.txt` is: `S1 R1 P0 R1`. This means that the first node, node 0, will run 4 events in the order that they are given. The meaning of each event is defined below:
- `S1`: Send a message to node 1. 
- `R1`: Receive a message from node 1.
- `PO`: A local event, no messages are sent or received between nodes.
- `R1`: Receive a message from node 1.   

### Running the Program 
To run the program with the input from the sample text file, use the following command: `go run vc.go < in.txt`   
Or, alternatively, use the command `go run vc.go` and just enter the input in the terminal.

## Tests
Tests are written in `vc_test.go`. These tests verify the correctness of the algorithm by running the vector clock program with different orderings and types of events for each node.   
According to the Microsoft [blog](https://blogs.msdn.microsoft.com/csliu/2009/05/18/time-and-order-of-events-in-distributed-system/), *Time and Order of Events in Distributed Systems*, an event X happens before an event Y if and only if at least one element in X's vector clock is strictly less than the corresponding element in Y's vector clock, and all other elements in X's vector clock are less than or equal to the corresponding elements in Y's vector clock.   
The tests verify that the vector clocks in the program maintain this property for all events between different nodes that have a strict "happens before" relationship.

### Running the Tests
To run the tests use the following command: `go test`  
For verbose output, run with the `-v` option. 