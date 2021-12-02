package main

import (
	sim "github.com/SOMAS2021/SOMAS2021/pkg/infra/simulation"
)

func main() {
	// can have frontend parameters come go straight into simEnv
	numOfAgents := []int{2, 3}
	simEnv := sim.New(100, numOfAgents, 100, 2)
	simEnv.Simulate()
}
