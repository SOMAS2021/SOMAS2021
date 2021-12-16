package simulation

import (
	"sync"
  
	"github.com/SOMAS2021/SOMAS2021/pkg/agents/randomAgent"
	"github.com/SOMAS2021/SOMAS2021/pkg/agents/team1/agent1"
	"github.com/SOMAS2021/SOMAS2021/pkg/agents/team1/agent2"
	"github.com/SOMAS2021/SOMAS2021/pkg/agents/team6"
	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/agent"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/day"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/world"
	log "github.com/sirupsen/logrus"
)

type Fields = log.Fields

type AgentNewFunc func(base *infra.Base) (agent.Agent, error)

type SimEnv struct {
	mx             sync.RWMutex
	FoodOnPlatform float64
	AgentCount     []int
	AgentHP        int
	AgentsPerFloor int
	logger         log.Entry
	dayInfo        *day.DayInfo
	reportFunc     func(*SimEnv)
	world          world.World
	custAgents     map[string]agent.Agent
}

func NewSimEnv(foodOnPlat float64, agentCount []int, agentHP, agentsPerFloor int, dayInfo *day.DayInfo) *SimEnv {
	return &SimEnv{
		FoodOnPlatform: foodOnPlat,
		AgentCount:     agentCount,
		AgentHP:        agentHP,
		dayInfo:        dayInfo,
		AgentsPerFloor: agentsPerFloor,
		logger:         *log.WithFields(log.Fields{"reporter": "simulation"}),
		custAgents:     make(map[string]agent.Agent),
	}
}

func (sE *SimEnv) Simulate() {
	sE.Log("Simulation Initializing")

	totalAgents := sum(sE.AgentCount)
	t := infra.NewTower(sE.FoodOnPlatform, totalAgents, sE.AgentsPerFloor, sE.dayInfo)
	sE.SetWorld(t)

	agentIndex := 1
	for i := 0; i < len(sE.AgentCount); i++ {
		for j := 0; j < sE.AgentCount[i]; j++ {
			sE.createNewAgent(t, i, agentIndex)
			agentIndex++
		}
	}
	sE.Log("Simulation Started")
	sE.simulationLoop(t)
	sE.Log("Simulation Ended")
}

// TODO: move to a general list of functions
func sum(inputList []int) int {

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
