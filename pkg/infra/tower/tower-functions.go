package tower

import (
	"log"
	"math/rand"
)

func (t *Tower) initAgents() {
	t.agents = make(map[string]BaseAgentCore, t.AgentCount)
}

func (t *Tower) killAgent(id string) { // this removes the agent from the list of agents in the tower
	log.Printf("Killing agent %s", id)
	deadAgentFloor := t.agents[id].floor
	agentType := t.agents[id].agentType
	t.missingAgents[deadAgentFloor] = append(t.missingAgents[deadAgentFloor], agentType)
	delete(t.agents, id)
}

func (t *Tower) reshuffle(agentsPerFloor int) {
	numOfFloors := t.AgentCount / int(agentsPerFloor)
	remainingVacanies := make([]int, numOfFloors)
	log.Printf("Reshuffling alive agents...")

	for i := 0; i < numOfFloors; i++ { // adding a max to each floor
		remainingVacanies[i] = agentsPerFloor
	}

	// allocating agents to floors randomly
	// iterate through the uuid strings of each agent
	for id := range t.agents {
		newFloor := rand.Intn(numOfFloors)

		for remainingVacanies[newFloor] == 0 {
			newFloor = rand.Intn(numOfFloors)
		}

		t.setFloor(id, newFloor+1)
		remainingVacanies[newFloor]--
	}
}

func (t *Tower) hpDecay() {
	// TODO: can add a parameter
	for id := range t.agents {
		newHP := t.agents[id].hp - 3 // TODO: change the function type (exp?parab?)

		if newHP < 0 {
			t.killAgent(id)
		} else {
			t.setHP(id, newHP)
		}
	}

}
