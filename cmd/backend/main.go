package main

import (
	"github.com/SOMAS2021/SOMAS2021/pkg/simulation"
)

func main() {
	// can have frontend parameters come go straight into simEnv
	foodOnPlatform := 100.0
	numOfAgents := []int{5, 5}
	agentHP := 100
	days := 10
	daysPerReshuffle := 3

	iterationsPerDay := 24 * 60
	iterations := days * iterationsPerDay
	reshufflePeriod := daysPerReshuffle * iterationsPerDay

	simEnv := simulation.NewSimEnv(foodOnPlatform, numOfAgents, agentHP, iterations, reshufflePeriod)
	simEnv.Simulate()
}
