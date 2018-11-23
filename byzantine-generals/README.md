# Byzantine Generals 

The Byzantine Generals problem is a famous consensus problem that was first introduced in a [paper](https://www.microsoft.com/en-us/research/uploads/prod/2016/12/The-Byzantine-Generals-Problem.pdf) by Lamport in 1982.   

In this problem, we have *n* generals and *m* of them are traitors. Among these generals there is one commander, who can also be a traitor. The commander sends out a command `ATTACK` or `RETREAT` to all the other lieutenants.  

In the case of a loyal commander, the goal is for all loyal lieutenants to agree on the command sent by the commander. In the case of a traitor commander, the goal is for all loyal lieutenants to agree on the same command, regardless of what it is. Lamport proposes an algorithm, *OM(m)*, in his paper that solves this problem as long as *m* is less than a third of *n*.   

In terms of applications, the Byzantine Generals problem boils down to solving the consensus problem in a distributed system, where we can have many nodes, some of which are faulty, and we want all the non-faulty nodes to agree on a value. 

## The Process

I originally began implementing a solution to this problem by trying to directly follow Lamport's algorithm, which can be found in the `lamport` sub-directory. However, after much puzzling, I realized that while Lamport's algorithm works in theory, it does not lend itself well to an actual implementation in a distributed system. For one, the algorithm uses recursion to get the values obtained by the generals at the level below it in order to calculate a majority. However, when we are writing this code for a distributed system, we cannot recurse in this way, all we can do is receive messages and send them out to other nodes.

Of course its possible for each node to maintain a tree of all the values it receives over the course of the message sharing, and then calculate the majority at each level by traversing this tree, but none of this information is mentioned anywhere in Lamport's algorithm. At this point, I decided to try a different randomized approach as well that seemed much simpler, but has a probabilistic result (ie: not guaranteed to succeed). This implementation can be found in the `probabilistic` sub-directory.

Overall I think the probabilistic algorithm is simpler to understand and implement. Since I had to modify my Lamport algorithm because I did not maintain a tree at each node, it no longer guaranteed correctness, and as a result, my probabilistic algorithm ended up performing much better than the Lamport one. It also used considerably less messages and space.


