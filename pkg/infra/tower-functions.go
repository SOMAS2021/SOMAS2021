package infra

import (
	"log"
	"math/rand"
)

func (t *Tower) initAgents() {
	t.agents = make([]Base, 0)
}

func (t *Tower) reshuffle(numOfFloors int) {
	remainingVacanies := make([]int, numOfFloors)
	log.Printf("Reshuffling alive agents...")
	for i := 0; i < numOfFloors; i++ { // adding a max to each floor
		remainingVacanies[i] = t.agentsPerFloor
	}
	// allocating agents to floors randomly
	// iterate through the uuid strings of each agent
	for _, agent := range t.agents {
		newFloor := rand.Intn(numOfFloors)
		for remainingVacanies[newFloor] == 0 {
			newFloor = rand.Intn(numOfFloors)
		}
		agent.setFloor(newFloor + 1)
		remainingVacanies[newFloor]--
	}
}

func (t *Tower) hpDecay() {
	// TODO: can add a parameter

	for _, agent := range t.agents {
		newHP := agent.HP() - 20 // TODO: change the function type (exp?parab?)
		if newHP < 0 {
			agent.Die()
			t.missingAgents[agent.Floor()] = append(t.missingAgents[agent.Floor()], agent.agentType)
		} else {
			agent.setHP(newHP)
		}
	}

}

func (t *Tower) ResetTower() {
	t.currPlatFood = t.maxPlatFood
	t.currPlatFloor = 1
}
