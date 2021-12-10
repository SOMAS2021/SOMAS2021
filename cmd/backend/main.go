package main

import (
	"fmt"
	"os"

	"github.com/SOMAS2021/SOMAS2021/pkg/simulation"
	"github.com/sirupsen/logrus"
)

func main() {
	// logger setup
	f, err := os.OpenFile("test.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	defer f.Close()
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}
	logrus.SetOutput(f)
	// can have frontend parameters come go straight into simEnv
	foodOnPlatform := 100.0
	numOfAgents := []int{5, 5}
	agentHP := 100
	days := 1
	daysPerReshuffle := 3

	iterationsPerDay := 24 * 60
	iterations := days * iterationsPerDay
	reshufflePeriod := daysPerReshuffle * iterationsPerDay

	simEnv := simulation.NewSimEnv(foodOnPlatform, numOfAgents, agentHP, iterations, reshufflePeriod)
	simEnv.Simulate()
}
