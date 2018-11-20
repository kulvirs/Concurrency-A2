package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

// Number of processes in the vector clock simulation.
const numProcesses = 3

// VectorClock represents the vector clock that is kept at each node to maintain ordering.
type VectorClock struct {
	// The process id that the vector clock belongs to.
	ID int
	// Stores the counter value for each process.
	Counters [numProcesses]int
}

// Stores channels for communicating between processes.
var channels [numProcesses][numProcesses]chan VectorClock

// Records all clock values for all processes.
var allClockValues [numProcesses][][numProcesses]int

// Increments the process's counter.
func (clock *VectorClock) inc() {
	clock.Counters[clock.ID]++
}

// Increments the process's counter and sends a message to another process.
func (clock *VectorClock) send(dest int) {
	clock.inc()
	channels[clock.ID][dest] <- *clock
}

// Receives a message from another process and updates the vector clock.
func (clock *VectorClock) recv(source int) {
	recvClock := <-channels[source][clock.ID]
	for i := 0; i < numProcesses; i++ {
		if i == clock.ID {
			clock.inc()
		} else if clock.Counters[i] <= recvClock.Counters[i] {
			clock.Counters[i] = recvClock.Counters[i]
			if i == recvClock.ID {
				clock.Counters[i]++
			}
		}
	}
}

// Populates the 2D channels array with all the channels needed to communicate between processes.
func createChannels() {
	for i := 0; i < numProcesses; i++ {
		for j := 0; j < numProcesses; j++ {
			channels[i][j] = make(chan VectorClock, numProcesses)
		}
	}
}

// A goroutine representing a node in the vector clock simulation.
func node(id int, events []string, wg *sync.WaitGroup) {
	clock := VectorClock{}
	clock.ID = id
	clockValues := [][numProcesses]int{} // records all clock values for the process
	for _, command := range events {
		switch command[0] {
		case 'P':
			// local process
			clock.inc()
		case 'S':
			// send to another process
			dest, err := strconv.Atoi(command[1:])
			if err != nil || dest >= numProcesses || dest < 0 || dest == clock.ID {
				// The value after the first character should be an integer in the range [0,numProcesses) and not equal to the current process id.
				log.Fatal("invalid command: " + command)
			}
			clock.send(dest)
		case 'R':
			// receive from another process
			src, err := strconv.Atoi(command[1:])
			if err != nil || src >= numProcesses || src < 0 || src == clock.ID {
				// The value after the first character should be an integer in the range [0, numProcesses) and not equal to the current process id.
				log.Fatal("invalid command: " + command)
			}
			clock.recv(src)
		default:
			// invalid character, should only be P, S, or R
			log.Fatal("invalid command: " + command)
		}
		clockValues = append(clockValues, clock.Counters)
	}
	allClockValues[clock.ID] = clockValues
	wg.Done()
}

// Prints the clock values for each process.
func printClockValues() {
	for i, clockValues := range allClockValues {
		fmt.Printf("Process %d timeline: ", i)
		for j, clockValue := range clockValues {
			if j != 0 {
				fmt.Printf("-> ")
			}
			fmt.Printf("%v ", clockValue)
		}
		fmt.Printf("\n")
	}
}

// RunVectorClock runs the vector clock by reading the events from the given event array (used for testing).
func RunVectorClock(events [numProcesses][]string) {
	createChannels()
	var wg sync.WaitGroup
	wg.Add(numProcesses)

	for i, processEvents := range events {
		go node(i, processEvents, &wg)
	}
	wg.Wait()
}

// Runs the vector clock by reading the events from stdin.
func main() {
	scanner := bufio.NewScanner(os.Stdin)
	createChannels()
	count := 0
	var wg sync.WaitGroup
	wg.Add(numProcesses)

	for scanner.Scan() {
		line := scanner.Text()
		go node(count, strings.Split(line, " "), &wg)
		count++
		if count > numProcesses {
			// The file should only contain lines equal to the number of processes.
			log.Fatal("The file should only contain", numProcesses, "lines")
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	wg.Wait()
	printClockValues()
}
