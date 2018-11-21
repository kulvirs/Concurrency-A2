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

// ATTACK represents the command, if it is true, ATTACK, else, RETREAT.
const ATTACK = true

// Message represents the structure of messages sent between generals.
type Message struct {
	// The sender of the message.
	Sender int
	// The list of previous generals that this message has already been sent by (in chronological order).
	Prev []int
	// The value being sent, true = ATTACK, false = RETREAT
	Value bool
	// The round of recursion, starting at m down to 0.
	Round int
}

// Calculates the majority of the given array of command values and returns either "ATTACK", "RETREAT", or "TIE" accordingly.
func majority(values []bool) string {
	attackCount := 0
	retreatCount := 0
	for _, value := range values {
		if value {
			attackCount++
		} else {
			retreatCount++
		}
	}

	if attackCount > retreatCount {
		return "ATTACK"
	} else if attackCount < retreatCount {
		return "RETREAT"
	} else {
		return "TIE"
	}
}

// Returns true if the value i is in the array of ints, false otherwise.
func in(array []int, i int) bool{
	for _, value := range(array){
		if value == i {
			return true
		}
	}
	return false
}

func commander(n int, m int, id int, loyal bool, command bool, channels []chan Message) {
	for i := 1; i < n; i++ {
		var msg Message
		if loyal == false && i%2 == 0 {
			// Traitor commander sending to even-valued general flips the command.
			msg = Message{id, []int{id}, !command, m}
		} else {
			msg = Message{id, []int{id}, command, m}
		}
		channels[i] <- msg
	}
}

func lieutenant(n int, m int, id int, loyal bool, channels []chan Message, wg *sync.WaitGroup) {
	msg := <-channels[id]
	fmt.Printf("Lieutenant %d received message from commander with command %v\n", id, msg.Value)

	prev := msg.Prev
	prev = append(prev, id)
	order := msg.Value

	for i := 0; i < n; i++ {
		if in(prev, i) == false {
			// Node i has not yet received this message.
			var newMsg Message
			if loyal == false && i%2 == 0 {
				// Traitor lieutenant sending to even-valued lieutenant flips the command.
				newMsg = Message{id, prev, !msg.Value, msg.Round - 1}
			} else {
				newMsg = Message{id, prev, msg.Value, msg.Round - 1}
			}
			channels[i] <- newMsg
		}
	}

	values := []bool{}
	for i := 0; i < n-2; i++ {
		msg = <-channels[id]
		values = append(values, msg.Value)
	}

	values = append(values, order)
	majorityValue := majority(values)
	fmt.Printf("Lieutenant %d: %s\n", id, majorityValue)

	wg.Done()
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var wg sync.WaitGroup

	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
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
	wg.Add(n - 1)

	// Order sent out by the commander.
	var cOrder bool
	if lines[2] == "ATTACK" {
		cOrder = true
	} else {
		cOrder = false
	}

	channels := []chan Message{}
	for range generals {
		channels = append(channels, make(chan Message, n*n*m))
	}

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
			commander(n, m, i, loyal, cOrder, channels)
		} else {
			// Create a goroutine for the lieutenant
			go lieutenant(n, m, i, loyal, channels, &wg)
		}
	}
	wg.Wait()
}
