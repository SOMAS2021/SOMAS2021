package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/SOMAS2021/SOMAS2021/pkg/simulation"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/day"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/health"
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
	foodOnPlatform := food.FoodType(100)
	numOfAgents := []int{0, 0, 0, 0, 10, 0, 0, 0} //agent1, agent2, team3, team6, team7agent1, team7agent2, team7agent3, randomAgent
	agentHP := 100
	agentsPerFloor := 1 //more than one not currently supported
	numberOfFloors := simulation.Sum(numOfAgents) / agentsPerFloor
	ticksPerFloor := 1

	ticksPerDay := numberOfFloors * ticksPerFloor
	simDays := 8
	reshuffleDays := 1
	dayInfo := day.NewDayInfo(ticksPerFloor, ticksPerDay, simDays, reshuffleDays)

	// define health parameters
	maxHP := 100
	weakLevel := 10
	width := 45.0
	tau := 10.0
	hpReqCToW := 2
	hpCritical := 5
	maxDayCritical := 3
	HPLossBase := 10
	HPLossSlope := 0.25

	healthInfo := health.NewHealthInfo(maxHP, weakLevel, width, tau, hpReqCToW, hpCritical, maxDayCritical, HPLossBase, HPLossSlope)

	// TODO: agentParameters - struct

	simEnv := simulation.NewSimEnv(foodOnPlatform, numOfAgents, agentHP, agentsPerFloor, dayInfo, healthInfo)
	simEnv.Simulate()
}
