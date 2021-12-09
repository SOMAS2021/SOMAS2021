package main

import (
	sim "github.com/SOMAS2021/SOMAS2021/pkg/infra/simulation"
)

func main() {
	// can have frontend parameters come go straight into simEnv
	foodOnPlatform := 100.0
	numOfAgents := []int{3}
	agentHP := 100
<<<<<<< HEAD
	days := 10
	daysPerReshuffle := 3

	iterationsPerDay := 24 * 60
	iterations := days * iterationsPerDay
	reshufflePeriod := daysPerReshuffle * iterationsPerDay
=======
	days := 4
	reshufflePeriod := 3.0 // in days

	iterationsPerDay := 24.0 * 60.0
	iterations := int(days)// int(days * iterationsPerDay)
	iterationsBetweenReshuffles := int(reshufflePeriod * iterationsPerDay)
>>>>>>> removed all concurrency, each agents runs one after another, need to separate run into phases, send msg, recieve msg, eat food

	simEnv := sim.New(foodOnPlatform, numOfAgents, agentHP, iterations, reshufflePeriod)
	simEnv.Simulate()
}
