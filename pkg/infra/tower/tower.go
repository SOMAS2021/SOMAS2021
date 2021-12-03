package tower

import (
	"log"
	"math/rand"
	"sync"

	baseagent "github.com/SOMAS2021/SOMAS2021/pkg/agents/default"
	"github.com/divan/goabm/abm"
)

type Tower struct {
	FoodOnPlatform  float64
	FloorOfPlatform uint64
	mx              sync.RWMutex
	agents          []baseagent.BaseAgent // abm.Agent
	AgentCount      int
	// uuid
	// missing agents []abm.Agent
}

func New(foodOnPlatform float64, floorOfPlatform uint64, agentCount int) *Tower {
	// generate uuid randomly
	t := &Tower{
		FoodOnPlatform:  foodOnPlatform,
		FloorOfPlatform: floorOfPlatform,
		AgentCount:      agentCount,
		// uuid: randuuid
	}

	t.initAgents()

	return t
}

func (t *Tower) initAgents() {
	t.agents = make([]abm.Agent, t.AgentCount)
}

func (t *Tower) killAgent( /*agent uuid or index in list(?)*/ ) {
	// this removes the agent from the list of agents in the tower

}

func (t *Tower) checkAgentsHP() {
	totalAgents := len(t.agents)
	for i := 0; i < totalAgents; i++ {
		// check HP:
		if t.agents[i].GetHP() <= 0 {
			t.killAgent()
		}

	}
}

func (t *Tower) replace(agentsPerFloor int) {
	//implementation

	totalAgents := len(t.agents)
	numOfFloors := totalAgents / int(agentsPerFloor)
	numOfAgentsInFloor := make([]int, numOfFloors, numOfFloors)

	for i := 0; i < totalAgents; i++ {
		numOfAgentsInFloor[t.agents[i].floor]++ //to be solved by importing packages
	}

	for floor := 0; floor < numOfFloors; floor++ {
		for numOfAgentsInFloor[floor] < agentsPerFloor {
			//add agent in that floor
			agent, err := baseagent.New(a, floor, 100)
		}
	}

	// doesn't kill people, that is a separate function
	// go through every floor, check how many agents and add them
	// if this isn't the same amount as agents per floor, add agents
}

func (t *Tower) reshuffle(agentsPerFloor int) {
	totalAgents := len(t.agents)
	numOfFloors := totalAgents / int(agentsPerFloor)
	remainingVacanies := make([]int, numOfFloors, numOfFloors)
	for i := 0; i < numOfFloors; i++ {
		remainingVacanies[i] = agentsPerFloor
	}

	for i := 0; i < totalAgents; i++ {
		newFloor := rand.Intn(numOfFloors) // random number in the range 0 - numOfFloors
		for remainingVacanies[newFloor] != 0 {
			newFloor := rand.Intn(numOfFloors)
		}

		//TODO: assign agent to currFloor
		(*agents)[i].(baseagent.SetFloor)(newFloor)
		remainingVacanies[newFloor]--
	}
	// go through list of agents one by one and access the struct
	// access what floor they are on
}

func (t *Tower) Tick() {
	t.mx.RLock()
	defer t.mx.RUnlock()
	log.Printf("A log from the tower!")
	// add all the tower functions here
	// replace(&t)
	// agree upon the reshuffle frequency
	// reshuffle(&tower.Agents, agentsPerFloor)
}

func (t *Tower) SetAgent(index int, agent abm.Agent) {
	t.mx.Lock()
	t.agents[index] = agent
	t.mx.Unlock()
}
