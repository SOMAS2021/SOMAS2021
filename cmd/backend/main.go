package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/SOMAS2021/SOMAS2021/pkg/simulation"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/day"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/health"
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
	logfileName := filepath.Join("logs", time.Now().Format("2006-01-02-15-04-05")+".json")

	// open latest archive
	f, err := os.OpenFile(logfileName, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("error opening file: ", err)
		return
	}

	// copy latest archive to simulation.log
	defer func() {
		// open simulation.log
		simLog, err := os.OpenFile("simulation.json", os.O_CREATE|os.O_RDWR, 0666)
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
	numOfAgents := []int{2, 0, 0, 0, 0} //agent1, agent2, team3, team6, randomAgent
	agentHP := 55
	//numberCriticalDays := 0
	agentsPerFloor := 1 //more than one not currently supported
	numberOfFloors := simulation.Sum(numOfAgents) / agentsPerFloor
	ticksPerFloor := 1

	ticksPerDay := numberOfFloors * ticksPerFloor
	simDays := 8
	reshuffleDays := 1
	dayInfo := day.NewDayInfo(ticksPerFloor, ticksPerDay, simDays, reshuffleDays)

	// define heath parameters
	strongLevel := 55
	healthyLevel := 25
	weakLevel := 5 // larger than maxCriticalDay
	maxDayCritical := 3
	foodReqStrong := 20.0
	foodReqHealthy := 15.0
	foodReqWeak := 10.0
	foodReqCToW := 2.0

	healthInfo := health.NewHealthInfo(strongLevel, healthyLevel, weakLevel, foodReqStrong, foodReqHealthy, foodReqWeak, foodReqCToW, maxDayCritical)

	// TODO: agentParameters - struct

	simEnv := simulation.NewSimEnv(foodOnPlatform, numOfAgents, agentHP, agentsPerFloor, dayInfo, healthInfo)
	simEnv.Simulate()
}
