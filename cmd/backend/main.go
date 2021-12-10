package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/SOMAS2021/SOMAS2021/pkg/simulation"
	log "github.com/sirupsen/logrus"
)

func main() {
	// logger setup
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		err := os.Mkdir("logs", 0755)
		if err != nil {
			fmt.Println("failed to create logs directory: ", err)
			return
		}
	}

	logfileName := "logs/" + time.Now().Format(time.RFC3339) + ".log"

	f, err := os.OpenFile(logfileName, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("error opening file: ", err)
		return
	}
	defer func() {
		simLog, err := os.OpenFile("simulation.log", os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			fmt.Println("error creating simulation log: ", err)
			return
		}
		defer simLog.Close()
		f, err = os.Open(logfileName)
		if err != nil {
			fmt.Println("error opening file: ", err)
			return
		}
		defer f.Close()
		_, err = io.Copy(simLog, f)
		if err != nil {
			fmt.Println("error copying to simulation log: ", err)
			return
		}
	}()

	log.SetOutput(f)
	log.SetFormatter(&log.JSONFormatter{})

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
