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
	remainingVacanies := make([]int, numOfFloors)
	log.Printf("Reshuffling alive agents...")

	for i := 0; i < numOfFloors; i++ { // adding a max to each floor
		remainingVacanies[i] = t.agentsPerFloor
	}

	// allocating agents to floors randomly
	// iterate through the uuid strings of each agent
	for id := range t.agents {
		newFloor := rand.Intn(numOfFloors)

		for remainingVacanies[newFloor] == 0 {
			newFloor = rand.Intn(numOfFloors)
		}
		t.mx.RLock()
		currHP := t.agents[id].hp
		aType := t.agents[id].agentType
		t.mx.RUnlock()
		// currHP++
		// aType++
		t.setFloor(id, currHP, newFloor+1, aType)
		remainingVacanies[newFloor]--
	}
}

func (t *Tower) hpDecay() {
	// TODO: can add a parameter

	for id := range t.agents {
		t.mx.RLock()
		newHP := t.agents[id].hp - 20 // TODO: change the function type (exp?parab?)
		floor := t.agents[id].floor
		aType := t.agents[id].agentType
		t.mx.RUnlock()
		if newHP < 0 {
			t.killAgent(id)
		} else {
			t.setHP(id, newHP, floor, aType)
		}
	}

}

func (t *Tower) updateHP(id string, foodTaken float64) {
	t.mx.RLock()
	newHP := math.Min(100, float64(t.agents[id].hp)+foodTaken)
	floor := t.agents[id].floor
	aType := t.agents[id].agentType
	t.mx.RUnlock()
	t.setHP(id, int(newHP), floor, aType)
}

func (t *Tower) FoodRequest(id string, foodRequested float64) float64 {
	t.mx.RLock()
	defer t.mx.RUnlock()
	if t.agents[id].floor == int(t.currPlatFloor) {
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

func (tower *Tower) SendMessage(direction int, sender abm.Agent , msg messages.Message){
	tower.mx.RLock()
	defer tower.mx.RUnlock()
	var senderFloor int
	for _, agent := range tower.agents {
		//find sender BaseAgentCore
		if (agent.cust.ID() == sender.ID()){
			senderFloor = agent.floor
		}
	}
	for _, agent := range tower.agents {
		//find reciever and pass them msg
		if (agent.floor == senderFloor + direction){
			reciever := agent.cust
			reciever.AddToInbox(msg)
		}
	}
}