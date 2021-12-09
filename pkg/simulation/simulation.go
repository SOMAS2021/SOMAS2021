package simulation

import (
	"log"

	"github.com/SOMAS2021/SOMAS2021/pkg/agents/team1/agent1"
	"github.com/SOMAS2021/SOMAS2021/pkg/agents/team1/agent2"
	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
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

func NewSimEnv(foodOnPlat float64, agentCount []int, agentHP, iterations, reshufflePeriod int) *SimEnv {

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

type AgentNewFunc func(base *infra.Base) (abm.Agent, error)

func (sE *SimEnv) Simulate() {
	a := abm.New()

	totalAgents := sum(sE.AgentCount)
	t := infra.NewTower(sE.FoodOnPlatform, 1, totalAgents, 1, 20, sE.reshufflePeriod)
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

func (sE *SimEnv) simulationLoop(a *abm.ABM, t *infra.Tower) {
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

func (sE *SimEnv) createNewAgent(a *abm.ABM, tower *infra.Tower, i, floor int) {
	// TODO: clean this looping, make a nice abs map
	log.Printf("Creating new agent")
	abs := []AgentNewFunc{agent1.New, agent2.New}

	uuid := uuid.New().String()
	bagent, err := infra.NewBaseAgent(a, i, sE.AgentHP-3*floor, floor, uuid)
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
