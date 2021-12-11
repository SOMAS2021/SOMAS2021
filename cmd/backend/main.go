package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/SOMAS2021/SOMAS2021/pkg/simulation"
	log "github.com/sirupsen/logrus"
)

func main() {
	// logger setup
	// TODO: clean up logger initialisation and closing code
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		err := os.Mkdir("logs", 0755)
		if err != nil {
			fmt.Println("failed to create logs directory: ", err)
			return
		}
	}

	// archive logs by default
	logfileName := filepath.Join("logs", time.Now().Format("2006-01-02-15-04-05")+".log")

	// open latest archive
	f, err := os.OpenFile(logfileName, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("error opening file: ", err)
		return
	}

	// copy latest archive to simulation.log
	defer func() {
		// open simulation.log
		simLog, err := os.OpenFile("simulation.log", os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			fmt.Println("error creating simulation log: ", err)
			return
		}
		// close simulation.log
		defer simLog.Close()
		f, err = os.Open(logfileName)
		if err != nil {
			fmt.Println("error opening file: ", err)
			return
		}
		// close latest archive
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
	numOfAgents := []int{2, 2, 2}
	agentHP := 100
	days := 10
	daysPerReshuffle := 3

	iterationsPerDay := 24 * 60
	iterations := days * iterationsPerDay
	reshufflePeriod := daysPerReshuffle * iterationsPerDay

	simEnv := simulation.NewSimEnv(foodOnPlatform, numOfAgents, agentHP, iterations, reshufflePeriod)
	simEnv.Simulate()
}
