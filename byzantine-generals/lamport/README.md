# Lamport
An attempt to implement Lamport's algorithm is in `bg_test.go`. Each general is modelled as a separate goroutine and generals communicate with each other through channels.

Since there is no way to do the recursion when each general is a separate goroutine without maintaining a tree data structure at each node, which I did not have time to implement, I decided to simplify Lamport's algorithm by just appending each order a node receives to an array, and then taking the majority of the values in the array at the end. This obviously will not return ideal results, but I was curious to see how it would perform.

For the cases where *m* is less than 2, the algorithm should produce the correct solution since it runs exactly as Lamport's algorithm indicates it should. It's only when *m* is greater than 2 that my algorithm starts to deviate from Lamport's, so changes will probably be seen at this stage.

## Input
A sample file `in.txt` is provided. The first line of the input file contains an integer, *m*, indicating the number of traitorous generals (including the commander). 
The second line of the file contains a list of generals separated by spaces. Each general contains a value `L` or `T` after it that indicates whether it is loyal or a traitor. The first general in the list is the commander.   
The third line contains the command that will be relayed by the commander to the rest of the lieutenants, it is either `ATTACK` or `RETREAT`

For example, the first line of the following input file indicates that there is one traitorous general. The second line tells us the commander, `G0`, is loyal as well as lieutenants `G1` and `G2`, and `G3` is a traitor. The last line indicates that the command sent out by the commander should be `ATTACK`.

```
1
G0:L G1:L G2:L G3:T
ATTACK
```

## Running the Program
To run the program with the input from the sample text file, use the command `go run bg.go < in.txt`.   
Or alternatively, use the command `go run bg.go` and just enter the input in the terminal.

## Tests
Since the message sharing part of Lamport's algorithm sends at most *(n-1)(n-2)...(n-m)* messages per round, I had to make the buffer on the channels used to communicate this size so they could send without blocking during a round. Since the channels send Message objects, which are each 48 bytes, for any value of *m* greater than 2, with a corresponding value of *n = 3m + 1* causes an out of memory exception when allocating memory for the channel, so I was only able to test with values of *m* below 3.  

I ran two types of tests. Each one calculates the percentage of trials that are successful for values of *m* ranging from 0 to 2, where the corresponding value of *n* is *3m+1* so as to maximize the number of traitors. 

All trials in the first test have a loyal commander. The following results were gathered:
```
100.00% trials successful for m = 0, n = 1
100.00% trials successful for m = 1, n = 4
64.00% trials successful for m = 2, n = 7
```

All trials in the second test have a traitor commander. The following results were gathered:
```
100.00% trials successful for m = 0, n = 1
100.00% trials successful for m = 1, n = 4
72.00% trials successful for m = 2, n = 7
```

As predicted, we can see that for values of m less than 2, we are able to get the correct result everytime, but this is not guaranteed when m is greater than 2. It would be interesting to see what the percentages are for higher values of m, but due to the size allocation error, we can only experiment this far (with goroutines anyway).

## Running the Tests
To run the tests, use the following command: `go test`  
For verbose output, run with the `-v` option.