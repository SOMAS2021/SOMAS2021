package tower

import (
	"log"
	"math"
	"math/rand"
)

func (t *Tower) initAgents() {
	t.agents = make(map[string]BaseAgentCore, t.agentCount)
}

func (t *Tower) killAgent(id string) { // this removes the agent from the list of agents in the tower
	log.Printf("Killing agent %s", id)
	t.mx.RLock()
	deadAgentFloor := t.agents[id].floor
	agentType := t.agents[id].agentType
	t.mx.RUnlock()
	t.missingAgents[deadAgentFloor] = append(t.missingAgents[deadAgentFloor], agentType)
	t.mx.Lock()
	delete(t.agents, id)
	t.mx.Unlock()
}

func (t *Tower) reshuffle(numOfFloors int) {
	remainingVacancies := make([]int, numOfFloors)
	log.Printf("Reshuffling alive agents...")

	for i := 0; i < numOfFloors; i++ { // adding a max to each floor
		remainingVacancies[i] = t.agentsPerFloor
	}

	// allocating agents to floors randomly
	// iterate through the uuid strings of each agent
	for id := range t.agents {
		newFloor := rand.Intn(numOfFloors)

		for remainingVacancies[newFloor] == 0 {
			newFloor = rand.Intn(numOfFloors)
		}
		t.setFloor(id, newFloor+1)
		remainingVacancies[newFloor]--
	}
}

func (t *Tower) hpDecay() {
	// TODO: can add a parameter

	for id := range t.agents {
		t.mx.RLock()
		newHP := t.agents[id].hp - 20 // TODO: change the function type (exp?parab?)
		t.mx.RUnlock()
		if newHP < 0 {
			t.killAgent(id)
		} else {
			t.setHP(id, newHP)
		}
	}

}

func (t *Tower) updateHP(id string, foodTaken float64) {
	t.mx.RLock()
	newHP := math.Min(100, float64(t.agents[id].hp)+foodTaken)
	t.mx.RUnlock()
	t.setHP(id, int(newHP))
}

func (t *Tower) FoodRequest(id string, foodRequested float64) float64 {
	t.mx.RLock()
	defer t.mx.RUnlock()
	if t.agents[id].floor == t.currPlatFloor {
		foodTaken := math.Min(t.currPlatFood, foodRequested)
		t.mx.RUnlock()
		t.updateHP(id, foodTaken)
		t.mx.RLock()
		t.currPlatFood -= foodTaken
		return foodTaken
	}
	return 0.0
}

func (t *Tower) ResetTower() {
	t.currPlatFood = t.maxPlatFood
	t.currPlatFloor = 1
}
