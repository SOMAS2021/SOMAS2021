package simulation

import (
	"log"

	"github.com/SOMAS2021/SOMAS2021/pkg/agents"
	"github.com/SOMAS2021/SOMAS2021/pkg/agents/team1/agent1"
	"github.com/SOMAS2021/SOMAS2021/pkg/agents/team1/agent2"
	"github.com/SOMAS2021/SOMAS2021/pkg/infra/tower"
	"github.com/divan/goabm/abm"
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
	tower := tower.New(sE.FoodOnPlatform, 1, totalAgents, 1, 20, sE.reshufflePeriod)
	a.SetWorld(tower)

	// TODO: clean this looping, make a nice abs map
	abs := []AgentNewFunc{agent1.New, agent2.New}

	agentIndex := 0
	for i := 0; i < len(sE.AgentCount); i++ {
		for j := 0; j < sE.AgentCount[i]; j++ {
			// generates the UUID
			uuid := uuid.New().String()
			bagent, err := agents.NewBaseAgent(a, uuid)
			if err != nil {
				log.Fatal(err)
			}
			// generates a custom agent on the current base agent
			custagent, err := abs[i%len(abs)](bagent)
			if err != nil {
				log.Fatal(err)
			}
			// adds custom agent to the world & cof checking ontroller
			a.AddAgent(custagent)
			// we pass in the agentIndex as the floor but right now it just consecutively declares agents per floor
			tower.SetAgent(sE.AgentHP, agentIndex+1, uuid, custagent) // TODO: edit this call
			agentIndex++
		}
	}

	a.LimitIterations(sE.Iterations)
	a.StartSimulation()

}

func sum(inputList []int) int {
	totalAgents := 0
	for _, value := range inputList {
		totalAgents += value
	}
	return totalAgents
}
