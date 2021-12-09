package main

import (
	sim "github.com/SOMAS2021/SOMAS2021/pkg/infra/simulation"
)

func main() {
	// can have frontend parameters come go straight into simEnv
	foodOnPlatform := 100.0
	numOfAgents := []int{3}
	agentHP := 100
	days := 4
	reshufflePeriod := 3.0 // in days

	iterationsPerDay := 24.0 * 60.0
	iterations := int(days)// int(days * iterationsPerDay)
	iterationsBetweenReshuffles := int(reshufflePeriod * iterationsPerDay)

	simEnv := sim.New(foodOnPlatform, numOfAgents, agentHP, iterations, iterationsBetweenReshuffles)
	simEnv.Simulate()
}
