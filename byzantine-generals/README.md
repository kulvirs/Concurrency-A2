# Byzantine Generals 

The Byzantine Generals problem is a famous consensus problem that was first introduced in a [paper](https://www.microsoft.com/en-us/research/uploads/prod/2016/12/The-Byzantine-Generals-Problem.pdf) by Lamport in 1982.   

In this problem, we have *n* generals and *m* of them are traitors. Among these generals there is one commander, who can also be a traitor. The commander sends out a command `ATTACK` or `RETREAT` to all the other lieutenants.  

In the case of a loyal commander, the goal is for all loyal lieutenants to agree on the command sent by the commander. In the case of a traitor commander, the goal is for all loyal lieutenants to agree on the same command, regardless of what it is. Lamport proposes an algorithm, *OM(m)*, in his paper that solves this problem as long as *m* is less than a third of *n*.   

In terms of applications, the Byzantine Generals problem boils down to solving the consensus problem in a distributed system, where we can have many nodes, some of which are faulty, and we want all the non-faulty nodes to agree on a value. 

## The Process

I originally began implementing a solution to this problem by trying to directly follow Lamport's algorithm. However, after much puzzling, I realized that while Lamport's algorithm works in theory, it does not lend itself well to an actual implementation in a distributed system. For one, the algorithm uses recursion to get the values obtained by the generals at the level below it in order to calculate a majority. However, when we are writing this code for a distributed system, we cannot recurse in this way, all we can do is receive messages and send them out to other nodes.

Of course its possible for each node to maintain a tree of all the values it receives over the course of the message sharing, and then calculate the majority at each level by traversing this tree, but none of this information is mentioned anywhere in Lamport's algorithm.  

At this point, I decided to try a different randomized approach as well that seemed much simpler, but has a probabilistic result (ie: not guaranteed to succeed).

The details and results of both implementations are recorded below.

## Implementation 1 - Lamport

## Implementation 2 - Probabilistic 
This implementation can be found in the folder `probabilistic`. Each general is modelled as a separate goroutine, and generals communicate with each other through channels.   

The algorithm I used closely follows Rabin's randomized global coin algorithm, a description of which can be found [here](https://www.cs.princeton.edu/courses/archive/fall05/cos521/byzantin.pdf). The link also contains a proof of correctness for the algorithm, and states that all loyal nodes can come to a consensus within a constant expected number of rounds.  

However, I modified the implementation of the algorithm slightly, because it does not specify when to terminate. To add a termination condition, I changed it so that after a round is over, all lieutenants send their vote to the commander. If the commander receives at least *2m + 1* votes of the same value, it assumes consensus has been reached among all loyal nodes and the algorithm terminates, otherwise we run another round. However, there is the possiblility that *m* of these *2m+1* votes are from traitor lieutenants sending the wrong value, and in fact all loyal generals have not yet reached a consensus, in which case the algorithm terminates too early. 

For a loyal commander, every loyal lieutenant receives the same initial command, so my modified algorithm is guaranteed to work. This is because if *m* is less than a third of n, each lieutenant will receive the same value from at least *2m + 1* other lieutenants, which means each loyal lieutenant will decide on the same value as the majority vote and send it to the commander, which will then terminate the algorithm. For a traitor commander, loyal lieutenants will receive different initial commands and it's possible the algorithm terminates before consensus is reached. The tests investigate the percentage of successes for a traitor commander with this algorithm. 

The main limitation with this algorithm is that it relies on the notion of a "global coin flip", where every general has access to some global variable that is randomly assigned an ATTACK or RETREAT value each round. In a distributed system this may not be possible. 

### Input
A sample input file `in.txt` in the folder `probabilistic` is provided.  
The first line of the input file contains an integer, *m*, indicating the number of traitorous lieutenants (not including the commander).  
The second line of the file contains a list of generals separated by spaces. Each general contains a value `L` or `T` after it that indicates whether it is loyal or a traitor. The first general in the list is the commander.   
The third line contains the command that will be relayed by the commander to the rest of the lieutenants, it is either `ATTACK` or `RETREAT`

For example, the first line of the following input file indicates that there is one traitorous general. The second line tells us the commander, `G0`, is loyal as well as lieutenants `G2`, `G3`, and `G4` and lieutenant `G1` is a traitor. The last line indicates that the command sent out by the commander should be `ATTACK`.

```
1
G0:L G1:T G2:L G3:L G4:L
ATTACK
```

### Running the Program
To run the program with the input from the sample text file, go into the folder `probabilistic` and use the command `go run bg-prob.go < in.txt`.   
Or alternatively, use the command `go run bg-prob.go` and just enter the input in the terminal.

### Tests
Tests are written in `bg-prob_test.go` in the `probabilistic` folder. There are two types of tests that are run. The first verifies correctness of the algorithm in the case of a loyal commander. It checks that for increasing values of *m* starting at *m = 0* up to *m = 100*, all loyal generals always agree on the value sent out by the commander. Since we want to test the worst case when the number of traitors is maximal, for each value of *m*, we set *n = 3m + 1*.  

The second type of test is for the case of a traitor commander but it doesn't really have a pass/fail value like the first one, since we know there are cases where my algorithm will not work for a traitor general. Instead for each value of *m* from 0 to 50 (with *n = 3m+1* for each *m*), it runs 100 trials of the program and calculates the percentage that succeed. The results from these tests are recorded below. 
```
100.00% trials successful for m = 0, n = 1
66.00% trials successful for m = 1, n = 4
74.00% trials successful for m = 2, n = 7
95.00% trials successful for m = 3, n = 10
90.00% trials successful for m = 4, n = 13
79.00% trials successful for m = 5, n = 16
93.00% trials successful for m = 6, n = 19
100.00% trials successful for m = 7, n = 22
94.00% trials successful for m = 8, n = 25
97.00% trials successful for m = 9, n = 28
97.00% trials successful for m = 10, n = 31
99.00% trials successful for m = 11, n = 34
100.00% trials successful for m = 12, n = 37
100.00% trials successful for m = 13, n = 40
100.00% trials successful for m = 14, n = 43
97.00% trials successful for m = 15, n = 46
100.00% trials successful for m = 16, n = 49
100.00% trials successful for m = 17, n = 52
100.00% trials successful for m = 18, n = 55
98.00% trials successful for m = 19, n = 58
100.00% trials successful for m = 20, n = 61
100.00% trials successful for m = 21, n = 64
100.00% trials successful for m = 22, n = 67
100.00% trials successful for m = 23, n = 70
100.00% trials successful for m = 24, n = 73
100.00% trials successful for m = 25, n = 76
100.00% trials successful for m = 26, n = 79
99.00% trials successful for m = 27, n = 82
100.00% trials successful for m = 28, n = 85
100.00% trials successful for m = 29, n = 88
100.00% trials successful for m = 30, n = 91
100.00% trials successful for m = 31, n = 94
100.00% trials successful for m = 32, n = 97
100.00% trials successful for m = 33, n = 100
100.00% trials successful for m = 34, n = 103
100.00% trials successful for m = 35, n = 106
100.00% trials successful for m = 36, n = 109
100.00% trials successful for m = 37, n = 112
100.00% trials successful for m = 38, n = 115
100.00% trials successful for m = 39, n = 118
100.00% trials successful for m = 40, n = 121
100.00% trials successful for m = 41, n = 124
100.00% trials successful for m = 42, n = 127
100.00% trials successful for m = 43, n = 130
100.00% trials successful for m = 44, n = 133
100.00% trials successful for m = 45, n = 136
100.00% trials successful for m = 46, n = 139
100.00% trials successful for m = 47, n = 142
100.00% trials successful for m = 48, n = 145
100.00% trials successful for m = 49, n = 148
100.00% trials successful for m = 50, n = 151
```

We can see that (aside from the trivial case of m = 0), the percentage of successes seem to increase as *m* grows and seems to be almost consistently at 100% as *m* gets very large. 

### Running the Tests
To run the tests use the following commands: `go test`  
For verbose output, run with the `-v` option.


