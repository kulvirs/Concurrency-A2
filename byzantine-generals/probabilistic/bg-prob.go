package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ATTACK represents the command, if it is true, ATTACK, else, RETREAT.
const ATTACK = true

// Calculates the majority of the given array of command values and returns either "ATTACK", "RETREAT", or "TIE" accordingly.
func majority(values []bool) (bool, int) {
	attackCount := 0
	retreatCount := 0
	for _, value := range values {
		if value == ATTACK {
			attackCount++
		} else {
			retreatCount++
		}
	}

	if attackCount > retreatCount {
		return ATTACK, attackCount
	}
	return !ATTACK, retreatCount
}

// Returns true if the value i is in the array of ints, false otherwise.
func in(array []int, i int) bool {
	for _, value := range array {
		if value == i {
			return true
		}
	}
	return false
}

// Converts the boolean command to a string.
func convertCommand(command bool) string {
	if command == true {
		return "ATTACK"
	}
	return "RETREAT"
}

// Sends a command on the given channel. Flips the command if the sender is a traitor and sending to an even-valued general.
func send(channel chan bool, sender int, receiver int, command bool, loyal bool) {
	if loyal == false && receiver%2 == 0 {
		// Traitor general sending to an even-valued general flips the command.
		channel <- !command
	} else {
		channel <- command
	}
}

func commander(n int, m int, id int, loyal bool, command bool, channels []chan bool, wg *sync.WaitGroup) {
	// Generate a global coin flip value.
	defer wg.Done()
	rand.Seed(time.Now().UnixNano())
	coinFlip := rand.Intn(2) == 0

	// Send out initial command to all nodes.
	for i := 1; i < n; i++ {
		send(channels[i], id, i, command, loyal)
		// Send global coin flip value to all nodes.
		channels[i] <- coinFlip
	}

	for {
		// Receive each node's value.
		values := []bool{}
		for i := 1; i < n; i++ {
			value := <-channels[id]
			values = append(values, value)
		}

		// Calculate the number of nodes agreeing on the majority value among the nodes.
		_, tally := majority(values)
		if tally >= (2*m)+1 {
			// No more rounds.
			for i := 1; i < n; i++ {
				close(channels[i])
			}
			break
		} else {
			// Not all loyal nodes are in agreement. Run another round with a new global coin flip value.
			coinFlip = rand.Intn(2) == 0
			for i := 1; i < n; i++ {
				channels[i] <- coinFlip
			}
		}
	}
}

func lieutenant(n int, m int, id int, loyal bool, commChannels []chan bool, channels []chan bool, commands []bool, wg *sync.WaitGroup) {
	defer wg.Done()

	// Get initial command from commander.
	command := <-commChannels[id]
	//fmt.Printf("Lieutenant %d received message from commander with command %v\n", id, convertCommand(command))

	for {
		coinFlip, more := <-commChannels[id]
		if more == false {
			// End of algorithm
			break
		}

		// Send command to all other lieutenants.
		for i := 1; i < n; i++ {
			send(channels[i], id, i, command, loyal)
		}

		// Receive commands from all other lieutenants.
		values := []bool{}
		for i := 1; i < n; i++ {
			value := <-channels[id]
			values = append(values, value)
		}

		// Compute the majority.
		majority, tally := majority(values)
		if tally >= (2*m)+1 {
			command = majority
		} else {
			command = coinFlip
		}

		// Send majority value back to commander.
		send(commChannels[0], id, 0, command, loyal)
		// Update the entry for this node in the array of commands.
		commands[id] = command
	}
}

// Runs the byzantine generals simulation with the given inputs.
// m is the number of traitors (not including the commander)
// generals is an array of generals with index 0 being the commander. The value at each index, i, is true if general i is loyal, false otherwise.
// commOrder is the order that the commander will relay to the lieutenants, true = ATTACK, false = RETREAT.
// Returns an array containing the final command made by each lieutenant i at index i.
func runGenerals(m int, generals []bool, commOrder bool) []bool {
	var wg sync.WaitGroup
	n := len(generals)
	wg.Add(n)

	// Create channels used to communicate between lieutenants.
	channels := []chan bool{}
	for range generals {
		channels = append(channels, make(chan bool, n))
	}

	//Create channels used to communicate between commander and lieutenants.
	commChannels := []chan bool{}
	for range generals {
		commChannels = append(commChannels, make(chan bool, n))
	}

	// This stores the final command at each node by the end of the algorithm.
	commands := make([]bool, n)

	for i, loyal := range generals {
		if i == 0 {
			// Get the commander to send out initial commands.
			go commander(n, m, i, loyal, commOrder, commChannels, &wg)
		} else {
			// Create a goroutine for each lieutenant.
			go lieutenant(n, m, i, loyal, commChannels, channels, commands, &wg)
		}
	}
	wg.Wait()
	return commands
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var wg sync.WaitGroup

	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	m, err := strconv.Atoi(lines[0])
	if err != nil {
		// The first line should be an integer indicating what m is.
		log.Fatal(err)
	}

	generals := strings.Split(lines[1], " ")

	// Number of generals
	n := len(generals)
	wg.Add(n)

	// Order sent out by the commander.
	var cOrder bool
	if lines[2] == "ATTACK" {
		cOrder = true
	} else {
		cOrder = false
	}

	// Create channels to communicate between generals.
	channels := []chan bool{}
	for range generals {
		channels = append(channels, make(chan bool, n))
	}

	//Create channels used to communicate between commander and lieutenants.
	commChannels := []chan bool{}
	for range generals {
		commChannels = append(commChannels, make(chan bool, n))
	}

	// After the algorithm terminates, this stores the final command made at each node.
	commands := make([]bool, n)

	for i, general := range generals {
		generalInfo := strings.Split(general, ":")
		var loyal bool
		if generalInfo[1] == "L" {
			loyal = true
		} else {
			loyal = false
		}

		if i == 0 {
			// Get the commander to send out initial commands.
			go commander(n, m, i, loyal, cOrder, commChannels, &wg)
		} else {
			// Create a goroutine for the lieutenants.
			go lieutenant(n, m, i, loyal, commChannels, channels, commands, &wg)
		}
	}
	wg.Wait()
	for i, command := range commands[1:] {
		fmt.Printf("Lieutenant %d: %s\n", i+1, convertCommand(command))
	}
}
