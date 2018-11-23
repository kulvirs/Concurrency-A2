package main

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

// Tests that all loyals generals always agree on the value sent by a loyal commander.
func TestLoyalCommander(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	// m is the number of traitors.
	for m := 0; m <= 100; m++ {
		// The number of lieutenants (not including the commander), must be greater than 3*m.
		n := 3*m + 1
		// The command that will be sent by the commander, randomly generated.
		command := rand.Intn(2) == 0

		generals := make([]bool, n+1)
		generals[0] = true // commander is loyal
		for i := 1; i <= n; i++ {
			generals[i] = true
		}

		perm := rand.Perm(n) // returns a permutation of the numbers [0, n)
		// For the first m values in the permutation, assign those generals to be traitors.
		for i := 0; i < m; i++ {
			generals[perm[i]+1] = false
		}

		commands := runGenerals(m, generals, command)

		// Verify all loyal lieutenants agreed on the command.
		for i := 1; i <= n; i++ {
			if generals[i] == true && commands[i] != command {
				t.Errorf("m = %d: Expected loyal general %d to decide command %s, but they decided %s", m, i, convertCommand(command), convertCommand(commands[i]))
			}
		}
	}

}

// Because the algorithm is probabilistic for a traitorous commander, we cannot guarantee a correct result.
// This test just records the percentage of runs that are successful for varying values of m.
func TestTraitorCommanderNumSuccess(t *testing.T) {
	numTrials := 100
	rand.Seed(time.Now().UnixNano())
	// m is the number of traitors.
	for m := 0; m <= 50; m++ {
		// The number of lieutenants (not including the commander), must be greater than 3*m.
		n := 3*m + 1
		// Keeps track of the number of successful trials.
		numSuccess := 0
		for r := 0; r < numTrials; r++ {
			// The command that will be sent by the commander, randomly generated.
			command := rand.Intn(2) == 0
			generals := make([]bool, n+1)
			generals[0] = false // commander is a traitor
			for i := 1; i <= n; i++ {
				generals[i] = true
			}

			perm := rand.Perm(n) // returns a permutation of the numbers [0, n)
			// For the first m values in the permutation, assign those generals to be traitors.
			for i := 0; i < m; i++ {
				generals[perm[i]+1] = false
			}

			commands := runGenerals(m, generals, command)

			// Verify all loyal lieutenants agree on the same value.
			firstLoyal := false
			loyalCommand := false
			agreement := true
			for i := 1; i <= n; i++ {
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
