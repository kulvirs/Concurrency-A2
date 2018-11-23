package main

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

// Records the percentage of runs that are successful for varying values of m when the commander is loyal.
func TestLoyalCommander(t *testing.T) {
	numTrials := 25
	rand.Seed(time.Now().UnixNano())
	// m is the number of traitors.
	for m := 0; m <= 2; m++ {
		// The number of generals must be greater than 3*m.
		n := 3*m + 1
		// Keeps track of the number of successful trials.
		numSuccess := 0
		for r := 0; r < numTrials; r++ {
			// The command that will be sent by the commander, randomly generated.
			command := rand.Intn(2) == 0
			generals := make([]bool, n)
			generals[0] = true // commander is loyal
			for i := 1; i < n; i++ {
				generals[i] = true
			}

			perm := rand.Perm(n - 1) // returns a permutation of the numbers [0, n)
			// For the first m values in the permutation, assign those generals to be traitors.
			for i := 0; i < m; i++ {
				generals[perm[i]+1] = false
			}

			commands := runGenerals(m, generals, command)

			// Verify all loyal lieutenants agreed on the command.
			agreement := true
			for i := 1; i < n; i++ {
				if generals[i] == true && commands[i] != command {
					agreement = false
					break
				}
			}

			if agreement {
				numSuccess++
			}
		}
		fmt.Printf("%0.2f%% trials successful for m = %d, n = %d\n", 100*(float64(numSuccess)/float64(numTrials)), m, n)
	}
}

// Records the percentage of runs that are successful for varying values of m when the commander is a traitor.
func TestTraitorCommander(t *testing.T) {
	numTrials := 25
	rand.Seed(time.Now().UnixNano())
	// m is the number of traitors.
	for m := 0; m <= 2; m++ {
		// The number of generals must be greater than 3*m.
		n := 3*m + 1
		// Keeps track of the number of successful trials.
		numSuccess := 0
		for r := 0; r < numTrials; r++ {
			// The command that will be sent by the commander, randomly generated.
			command := rand.Intn(2) == 0
			generals := make([]bool, n)
			generals[0] = false // commander is a traitor
			for i := 1; i < n; i++ {
				generals[i] = true
			}

			perm := rand.Perm(n - 1) // returns a permutation of the numbers [0, n)
			// For the first m-1 values in the permutation, assign those generals to be traitors.
			for i := 0; i < m-1; i++ {
				generals[perm[i]+1] = false
			}

			commands := runGenerals(m, generals, command)

			// Verify all loyal lieutenants agree on the same value.
			firstLoyal := false
			loyalCommand := false
			agreement := true
			for i := 1; i < n; i++ {
				if generals[i] == true {
					if firstLoyal == false {
						firstLoyal = true
						loyalCommand = commands[i]
					} else {
						if commands[i] != loyalCommand {
							agreement = false
							break
						}
					}
				}
			}

			if agreement {
				numSuccess++
			}
		}
		fmt.Printf("%0.2f%% trials successful for m = %d, n = %d\n", 100*(float64(numSuccess)/float64(numTrials)), m, n)
	}
}
