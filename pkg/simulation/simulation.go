package simulation

import (
	"github.com/SOMAS2021/SOMAS2021/pkg/config"
	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/day"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/health"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/world"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/utilFunctions"
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
}

func NewSimEnv(parameters *config.ConfigParameters, healthInfo *health.HealthInfo) *SimEnv {
	return &SimEnv{
		FoodOnPlatform: parameters.FoodOnPlatform,
		AgentCount:     parameters.NumOfAgents,
		AgentHP:        parameters.AgentHP,
		dayInfo:        parameters.DayInfo,
		healthInfo:     healthInfo,
		AgentsPerFloor: parameters.AgentsPerFloor,
		logger:         *log.WithFields(log.Fields{"reporter": "simulation"}),
	}
}

func (sE *SimEnv) Simulate() {
	sE.Log("Simulation Initializing")

	totalAgents := utilFunctions.Sum(sE.AgentCount)
	t := infra.NewTower(sE.FoodOnPlatform, totalAgents, sE.AgentsPerFloor, sE.dayInfo, sE.healthInfo)
	sE.SetWorld(t)

	sE.generateInitialAgents(t)

	sE.Log("Simulation Started")
	sE.simulationLoop(t)
	sE.Log("Simulation Ended")
	sE.Log("Summary of dead agents", infra.Fields{"Agent Type and number that died": t.DeadAgents()})
}

func (s *SimEnv) Log(message string, fields ...Fields) {
	if len(fields) == 0 {
		fields = append(fields, Fields{})
	}
	s.logger.WithFields(fields[0]).Info(message)
}
