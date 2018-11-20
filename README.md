# Overview
This repository examines two different problems that often come up in distributed systems: timing and consensus.   

The timing problem is explored using [vector clocks](http://book.mixu.net/distsys/time.html). An implementation of vector clocks is written in Go along with some tests to verify its correctness.  

The [Byzantine Generals](https://www.microsoft.com/en-us/research/publication/byzantine-generals-problem/) problem is a popular consensus problem. An implementation of an algorithm to achieve consensus is written in Go, along with some tests to verify its correctness.

Note that the nodes in the distributed system are simulated with goroutines in these implementations. Channels are used to communicate between the "nodes". 
