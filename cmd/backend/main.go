package main

import (
	simulation "github.com/SOMAS2021/SOMAS2021/pkg/simulation"
)

func main() {
	// can have frontend parameters come go straight into simEnv
	foodOnPlatform := 100.0
	numOfAgents := []int{5, 5}
	agentHP := 100
	days := 10.0
	reshufflePeriod := 3.0 // in days

	iterationsPerDay := 24.0 * 60.0
	iterations := int(days * iterationsPerDay)
	iterationsBetweenReshuffles := int(reshufflePeriod * iterationsPerDay)

	simEnv := simulation.NewSimEnv(foodOnPlatform, numOfAgents, agentHP, iterations, iterationsBetweenReshuffles)
	simEnv.Simulate()
}
