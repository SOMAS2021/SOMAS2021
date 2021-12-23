package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/SOMAS2021/SOMAS2021/pkg/simulation"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/day"
	log "github.com/sirupsen/logrus"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		err := os.Mkdir("logs", 0755)
		if err != nil {
			fmt.Println("failed to create logs directory: ", err)
			return
		}
	}

	logfileName := filepath.Join("logs", time.Now().Format("2006-01-02-15-04-05")+".json")
	f, err := os.OpenFile(logfileName, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("error opening file: ", err)
		return
	}
	defer f.Close()

	log.SetOutput(f)
	log.SetFormatter(&log.JSONFormatter{})

	// can have frontend parameters come go straight into simEnv
	foodOnPlatform := 100.0
	numOfAgents := []int{2, 2, 2, 2, 2}
	agentHP := 100
	agentsPerFloor := 1 //more than one not currently supported
	numberOfFloors := simulation.Sum(numOfAgents) / agentsPerFloor
	ticksPerFloor := 10

	ticksPerDay := numberOfFloors * ticksPerFloor
	simDays := 3
	reshuffleDays := 1
	dayInfo := day.NewDayInfo(ticksPerFloor, ticksPerDay, simDays, reshuffleDays)

	// TODO: agentParameters - struct

	simEnv := simulation.NewSimEnv(foodOnPlatform, numOfAgents, agentHP, agentsPerFloor, dayInfo)
	simEnv.Simulate()
}
