package main

import (
	sim "github.com/SOMAS2021/SOMAS2021/pkg/infra/simulation"
)

func main() {
	// can have frontend parameters come go straight into simEnv
	foodOnPlatform := 100.0
	numOfAgents := []int{2, 3}
	agentHP := 100
	iterations := 50
	reshufflePeriod := 20

	simEnv := sim.New(foodOnPlatform, numOfAgents, agentHP, iterations, reshufflePeriod)
	simEnv.Simulate()
}
