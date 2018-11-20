package main

import (
	"testing"
)

// Verifies that the event referenced by indices 0,1 comes before the event referenced by indices 2,3 by comparing the values of their vector clocks.
// According to https://blogs.msdn.microsoft.com/csliu/2009/05/18/time-and-order-of-events-in-distributed-system/, an event X happens before Y if
// and only if at least one element in X is strictly less than the corresponding element in Y, and all other elements in X are less than or equal
// to the corresponding elements in Y.
func verifyOrder(order [4]int) bool {
	firstClock := allClockValues[order[0]][order[1]]
	secondClock := allClockValues[order[2]][order[3]]

	lessThan := false
	for i := 0; i < numProcesses; i++ {
		if firstClock[i] > secondClock[i] {
			return false
		} else if firstClock[i] < secondClock[i] {
			lessThan = true
		}
	}
	return lessThan
}

// Tests the vector clock implementation maintains correct ordering when just one message is sent between two processes.
func TestVectorClockSingleMessage(t *testing.T) {
	p0Events := []string{"P0", "S1", "P0"}
	p1Events := []string{"P1", "P1", "R0", "P1"}
	p2Events := []string{"P2", "P2"}
	events := [3][]string{p0Events, p1Events, p2Events}

	// Each expected order is represented with an array of 4 integers
	// Example: {a, b, c, d} symbolizes that Process a Event b should come before Process c Event d
	// Here Process 0 Event 1 happens before Process 1 Event 2.
	expectedOrders := [][4]int{[4]int{0, 1, 1, 2}}

	RunVectorClock(events)
	// Check for correct ordering.
	for _, expectedOrder := range expectedOrders {
		if verifyOrder(expectedOrder) == false {
			t.Errorf("Expected Process %d Event %d to come before Process %d Event %d, but it did not.", expectedOrder[0], expectedOrder[1], expectedOrder[2], expectedOrder[3])
		}
	}
}

// Tests the vector clock implementation maintains the correct ordering when two messages are sent between two different pairs of processes.
func TestVectorClockTwoMessages(t *testing.T) {
	p0Events := []string{"S1", "P0", "P0"}
	p1Events := []string{"P1", "R0", "S2"}
	p2Events := []string{"P2", "R1", "P2"}
	events := [3][]string{p0Events, p1Events, p2Events}

	// Process 0 Event 0 happens before Process 1 Event 1.
	// Process 1 Event 2 happens before Process 2 Event 1.
	// Transitively, Process 0 Event 0 should come before Process 2 Event 1.
	expectedOrders := [][4]int{[4]int{0, 0, 1, 1}, [4]int{1, 2, 2, 1}, [4]int{0, 0, 2, 1}}

	RunVectorClock(events)
	// Check for correct ordering.
	for _, expectedOrder := range expectedOrders {
		if verifyOrder(expectedOrder) == false {
			t.Errorf("Expected Process %d Event %d to come before Process %d Event %d, but it did not.", expectedOrder[0], expectedOrder[1], expectedOrder[2], expectedOrder[3])
		}
	}
}

// Tests the vector clock implementation maintains the correct ordering when a single messages is sent between each pair of processes.
func TestVectorClockSingleMessageAllPairs(t *testing.T) {
	p0Events := []string{"R1", "S2"}
	p1Events := []string{"S0", "P1", "R2"}
	p2Events := []string{"S1", "P2", "R0"}
	events := [3][]string{p0Events, p1Events, p2Events}

	// Process 1 Event 0 happens before Process 0 Event 0.
	// Process 0 Event 1 happens before Process 2 Event 2.
	// Process 2 Event 0 happens before Process 1 Event 2.
	// Transitively, Process 1 Event 0 happens before Process 2 Event 2.
	expectedOrders := [][4]int{[4]int{1, 0, 0, 0}, [4]int{0, 1, 2, 2}, [4]int{2, 0, 1, 2}, [4]int{1, 0, 2, 2}}

	RunVectorClock(events)
	// Check for correct ordering.
	for _, expectedOrder := range expectedOrders {
		if verifyOrder(expectedOrder) == false {
			t.Errorf("Expected Process %d Event %d to come before Process %d Event %d, but it did not.", expectedOrder[0], expectedOrder[1], expectedOrder[2], expectedOrder[3])
		}
	}
}

// Tests the vector clock implementation maintains the correct ordering when a message is sent both ways between each pair of processes.
func TestVectorClockMessageAllPairsBothDirections(t *testing.T) {
	p0Events := []string{"S1", "R1", "R2", "S2"}
	p1Events := []string{"S0", "R0", "S2", "R2"}
	p2Events := []string{"S0", "R1", "S1", "R0"}
	events := [3][]string{p0Events, p1Events, p2Events}

	// Process 0 Event 0 happens before Process 1 Event 1.
	// Process 1 Event 0 happens before Process 0 Event 1.
	// Process 2 Event 0 happens before Process 0 Event 2.
	// Process 1 Event 2 happens before Process 2 Event 1.
	// Process 2 Event 2 happens before Process 1 Event 3.
	// Process 0 Event 3 happens before Process 2 Event 3.
	// Transitively, Process 0 Event 0 happens before Process 2 Event 1.
	// Transitively, Process 1 Event 0 happens before Process 2 Event 3.
	expectedOrders := [][4]int{[4]int{0, 0, 1, 1}, [4]int{1, 0, 0, 1}, [4]int{2, 0, 0, 2}, [4]int{1, 2, 2, 1}, [4]int{2, 2, 1, 3}, [4]int{0, 3, 2, 3}, [4]int{0, 0, 2, 1}, [4]int{1, 0, 2, 3}}

	RunVectorClock(events)
	// Check for correct ordering.
	for _, expectedOrder := range expectedOrders {
		if verifyOrder(expectedOrder) == false {
			t.Errorf("Expected Process %d Event %d to come before Process %d Event %d, but it did not.", expectedOrder[0], expectedOrder[1], expectedOrder[2], expectedOrder[3])
		}
	}
}
