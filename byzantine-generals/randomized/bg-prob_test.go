package main

import (
	"fmt"
	"testing"
)

func TestSample(t *testing.T) {
	commands := runGenerals(1, []bool{true, true, true, false, true}, true)
	fmt.Println(commands)
}
