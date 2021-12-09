package main

import (
	sim "github.com/SOMAS2021/SOMAS2021/pkg/infra/simulation"
)

func main() {
	// can have frontend parameters come go straight into simEnv
	foodOnPlatform := 100.0
	numOfAgents := []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	//creating agents: 1_1, 1_2, 2_1, 2_2, 3_1, 3_2, 4_1, 4_2, 5_1, 5_2, 6_1, 6_2, 7_1, 7_2
	agentHP := 100
	days := 10
	daysPerReshuffle := 3

	iterationsPerDay := 24 * 60
	iterations := days * iterationsPerDay
	reshufflePeriod := daysPerReshuffle * iterationsPerDay

	simEnv := sim.New(foodOnPlatform, numOfAgents, agentHP, iterations, reshufflePeriod)
	simEnv.Simulate()
}
