package simulation

import (
	"log"

	"github.com/SOMAS2021/SOMAS2021/pkg/agents"
	agent1_1 "github.com/SOMAS2021/SOMAS2021/pkg/agents/team1/agent1"
	agent1_2 "github.com/SOMAS2021/SOMAS2021/pkg/agents/team1/agent2"
	agent2_1 "github.com/SOMAS2021/SOMAS2021/pkg/agents/team2/agent1"
	agent2_2 "github.com/SOMAS2021/SOMAS2021/pkg/agents/team2/agent2"
	agent3_1 "github.com/SOMAS2021/SOMAS2021/pkg/agents/team3/agent1"
	agent3_2 "github.com/SOMAS2021/SOMAS2021/pkg/agents/team3/agent2"
	agent4_1 "github.com/SOMAS2021/SOMAS2021/pkg/agents/team4/agent1"
	agent4_2 "github.com/SOMAS2021/SOMAS2021/pkg/agents/team4/agent2"
	agent5_1 "github.com/SOMAS2021/SOMAS2021/pkg/agents/team5/agent1"
	agent5_2 "github.com/SOMAS2021/SOMAS2021/pkg/agents/team5/agent2"
	agent6_1 "github.com/SOMAS2021/SOMAS2021/pkg/agents/team6/agent1"
	agent6_2 "github.com/SOMAS2021/SOMAS2021/pkg/agents/team6/agent2"
	"github.com/SOMAS2021/SOMAS2021/pkg/infra/tower"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/abm"
	"github.com/google/uuid"
)

type SimEnv struct {
	FoodOnPlatform  float64
	AgentCount      []int
	AgentHP         int
	Iterations      int
	reshufflePeriod int
}

func New(foodOnPlat float64, agentCount []int, agentHP, iterations, reshufflePeriod int) *SimEnv {

	s := &SimEnv{
		FoodOnPlatform:  foodOnPlat,
		AgentCount:      agentCount,
		AgentHP:         agentHP,
		Iterations:      iterations,
		reshufflePeriod: reshufflePeriod,
	}
	// can do other inits here
	return s
}

type AgentNewFunc func(base *agents.Base) (agents.Agent, error)

func (sE *SimEnv) Simulate() {
	a := abm.New()

	totalAgents := sum(sE.AgentCount)
	t := tower.New(sE.FoodOnPlatform, 1, totalAgents, 1, 20, sE.reshufflePeriod)
	a.SetWorld(t)

	agentIndex := 1
	for i := 0; i < len(sE.AgentCount); i++ {
		for j := 0; j < sE.AgentCount[i]; j++ {
			sE.createNewAgent(a, t, i, agentIndex)
			agentIndex++
		}
	}
	a.LimitIterations(sE.Iterations)
	sE.simulationLoop(a, t)
}

func (sE *SimEnv) simulationLoop(a *abm.ABM, t *tower.Tower) {
	// TODO: abstract with an extra function
	for i := 1; i <= sE.Iterations; i++ {
		missingAgentsMap := t.UpdateMissingAgents()
		for floor := range missingAgentsMap {
			for _, agentType := range missingAgentsMap[floor] {
				sE.createNewAgent(a, t, agentType, floor)
			}
		}

		a.SimulationIterate(i)
	}
}

func sum(inputList []int) int {
	totalAgents := 0
	for _, value := range inputList {
		totalAgents += value
	}
	return totalAgents
}

func (sE *SimEnv) createNewAgent(a *abm.ABM, tower *tower.Tower, agentType, floor int) {
	// TODO: clean this looping, make a nice abs map
	abs := []AgentNewFunc{agent1_1.New, agent1_2.New, agent2_1.New, agent2_2.New, agent3_1.New,
		agent3_2.New, agent4_1.New, agent4_2.New, agent5_1.New, agent5_2.New, agent6_1.New, agent6_2.New}

	uuid := uuid.New().String()
	bagent, err := agents.NewBaseAgent(a, uuid)
	if err != nil {
		log.Fatal(err)
	}

	custagent, err := abs[agentType](bagent)
	if err != nil {
		log.Fatal(err)
	}

	a.AddAgent(custagent)
	tower.SetAgent(uuid, 100, floor, agentType) // TODO: edit this call
}
