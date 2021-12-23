package simulation

import (
	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/day"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/health"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/world"
	log "github.com/sirupsen/logrus"
)

type Fields = log.Fields

type SimEnv struct {
	FoodOnPlatform food.FoodType
	AgentCount     []int
	AgentHP        int
	AgentsPerFloor int
	logger         log.Entry
	dayInfo        *day.DayInfo
	healthInfo     *health.HealthInfo
	world          world.World
<<<<<<< HEAD
=======
	custAgents     map[string]infra.Agent
>>>>>>> 74745ce (chore: move agent interface to baseagent)
}

func NewSimEnv(foodOnPlat food.FoodType, agentCount []int, agentHP, agentsPerFloor int, dayInfo *day.DayInfo, healthInfo *health.HealthInfo) *SimEnv {
	return &SimEnv{
		FoodOnPlatform: foodOnPlat,
		AgentCount:     agentCount,
		AgentHP:        agentHP,
		dayInfo:        dayInfo,
		healthInfo:     healthInfo,
		AgentsPerFloor: agentsPerFloor,
		logger:         *log.WithFields(log.Fields{"reporter": "simulation"}),
	}
}

func (sE *SimEnv) Simulate() {
	sE.Log("Simulation Initializing")

	totalAgents := Sum(sE.AgentCount)
	t := infra.NewTower(sE.FoodOnPlatform, totalAgents, sE.AgentsPerFloor, sE.dayInfo, sE.healthInfo)
	sE.SetWorld(t)

	sE.generateInitialAgents(t)

	sE.Log("Simulation Started")
	sE.simulationLoop(t)
	sE.Log("Simulation Ended")
}

// TODO: move to a general list of functions
func Sum(inputList []int) int {
	totalAgents := 0
	for _, value := range inputList {
		totalAgents += value
	}
	return totalAgents
}

func (s *SimEnv) Log(message string, fields ...Fields) {
	if len(fields) == 0 {
		fields = append(fields, Fields{})
	}
	s.logger.WithFields(fields[0]).Info(message)
}
