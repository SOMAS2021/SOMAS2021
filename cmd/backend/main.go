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
	numOfAgents := []int{2, 2, 2, 2, 2} //agent1, agent2, team3, team6, randomAgent
	maxHP := 100
	agentHP := 100
	agentsPerFloor := 1 //more than one not currently supported
	numberOfFloors := simulation.Sum(numOfAgents) / agentsPerFloor
	ticksPerFloor := 10

	ticksPerDay := numberOfFloors * ticksPerFloor
	simDays := 8
	reshuffleDays := 1
	dayInfo := day.NewDayInfo(ticksPerFloor, ticksPerDay, simDays, reshuffleDays)

	// define heath parameters
	strongLevel := 55
	healthyLevel := 25
	weakLevel := 5
	maxDayCritical := 3
	tauStrong := 40.0
	tauHealthy := 16.0
	tauWeak := 4.0
	foodReqCToW := 2.0

	healthInfo := health.NewHealthInfo(maxHP, strongLevel, healthyLevel, weakLevel, tauStrong, tauHealthy, tauWeak, foodReqCToW, maxDayCritical)

	// TODO: agentParameters - struct

	simEnv := simulation.NewSimEnv(foodOnPlatform, numOfAgents, agentHP, agentsPerFloor, dayInfo, healthInfo)
	simEnv.Simulate()
}
