package simulation

import (
	"github.com/SOMAS2021/SOMAS2021/pkg/agents/team1/agent1"
	"github.com/SOMAS2021/SOMAS2021/pkg/agents/team1/agent2"
	"github.com/SOMAS2021/SOMAS2021/pkg/agents/team6"
	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/abm"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/agent"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/day"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type Fields = log.Fields

type SimEnv struct {
	FoodOnPlatform float64
	AgentCount     []int
	AgentHP        int
	AgentsPerFloor int
	logger         log.Entry
	dayInfo        *day.DayInfo
}

func (s *SimEnv) Log(message string, fields ...Fields) {
	if len(fields) == 0 {
		fields = append(fields, Fields{})
	}
	s.logger.WithFields(fields[0]).Info(message)
}

func NewSimEnv(foodOnPlat float64, agentCount []int, agentHP, agentsPerFloor int, dayInfo *day.DayInfo) *SimEnv {
	s := &SimEnv{
		FoodOnPlatform: foodOnPlat,
		AgentCount:     agentCount,
		AgentHP:        agentHP,
		dayInfo:        dayInfo,
		AgentsPerFloor: agentsPerFloor,
		logger:         *log.WithFields(log.Fields{"reporter": "simulation"}),
	}
	// can do other inits here
	return s
}

type AgentNewFunc func(base *infra.Base) (agent.Agent, error)

func (sE *SimEnv) Simulate() {
	sE.Log("Simulation Initializing")
	a := abm.New()

	totalAgents := sum(sE.AgentCount)
	t := infra.NewTower(sE.FoodOnPlatform, totalAgents, 1, sE.dayInfo)
	a.SetWorld(t)

	agentIndex := 1
	for i := 0; i < len(sE.AgentCount); i++ {
		for j := 0; j < sE.AgentCount[i]; j++ {
			sE.createNewAgent(a, t, i, agentIndex)
			agentIndex++
		}
	}
	sE.Log("Simulation Started")
	sE.simulationLoop(a, t)
	sE.Log("Simulation Ended")
}

func (sE *SimEnv) simulationLoop(a *abm.ABM, t *infra.Tower) {
	for sE.dayInfo.CurrTick <= sE.dayInfo.TotalTicks {
		missingAgentsMap := t.UpdateMissingAgents()
		for floor := range missingAgentsMap {
			for _, agentType := range missingAgentsMap[floor] {
				sE.createNewAgent(a, t, agentType, floor)
			}
		}
		a.SimulationIterate(sE.dayInfo.CurrTick)
		sE.dayInfo.CurrTick++
	}
}

func sum(inputList []int) int {
	totalAgents := 0
	for _, value := range inputList {
		totalAgents += value
	}
	return totalAgents
}

func (sE *SimEnv) createNewAgent(a *abm.ABM, tower *infra.Tower, i, floor int) {
	// TODO: clean this looping, make a nice abs map
	sE.Log("Creating new agent")
	abs := []AgentNewFunc{agent1.New, agent2.New, team6.New}

	uuid := uuid.New().String()
	bagent, err := infra.NewBaseAgent(a, i, sE.AgentHP, floor, uuid)
	if err != nil {
		log.Fatal(err)
	}

	custagent, err := abs[i](bagent)
	if err != nil {
		log.Fatal(err)
	}

	a.AddAgent(custagent)
	tower.AddAgent(bagent) // TODO: edit this call
}
