package tower

import (
	"log"
	"math"
	"math/rand"

	"github.com/SOMAS2021/SOMAS2021/pkg/infra/messages"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/abm"
)

func (t *Tower) initAgents() {
	t.agents = make(map[string]BaseAgentCore, t.agentCount)
}

func (t *Tower) killAgent(id string) { // this removes the agent from the list of agents in the tower
	log.Printf("Killing agent %s", id)
	t.mx.Lock()
	deadAgentFloor := t.agents[id].floor
	agentType := t.agents[id].agentType
	t.mx.Unlock()
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
		t.mx.Lock()
		newHP := t.agents[id].hp - 20 // TODO: change the function type (exp?parab?)
		t.mx.Unlock()
		if newHP < 0 {
			t.killAgent(id)
		} else {
			t.setHP(id, newHP)
		}
	}

}

func (t *Tower) updateHP(id string, foodTaken float64) {
	t.mx.Lock()
	newHP := math.Min(100, float64(t.agents[id].hp)+foodTaken)
	t.mx.Unlock()
	t.setHP(id, int(newHP))
}

func (t *Tower) FoodRequest(id string, foodRequested float64) float64 {
	t.mx.Lock()
	defer t.mx.Unlock()
	if t.agents[id].floor == t.currPlatFloor {
		foodTaken := math.Min(t.currPlatFood, foodRequested)
		t.mx.Unlock()
		t.updateHP(id, foodTaken)
		t.mx.Lock()
		t.currPlatFood -= foodTaken
		return foodTaken
	}
	return 0.0
}

func (t *Tower) ResetTower() {
	t.currPlatFood = t.maxPlatFood
	t.currPlatFloor = 1
}

func (tower *Tower) SendMessage(direction int, sender abm.Agent, msg messages.Message) {
	tower.mx.Lock()
	var senderFloor int
	for id, agent := range tower.agents {
		//find sender BaseAgentCore
		if id == (sender).ID() {
			senderFloor = agent.floor
		}
	}
	for _, agent := range tower.agents {
		//find reciever and pass them msg
		if agent.floor == senderFloor+direction {
			agent.inbox.PushBack(msg) //<- msg //for some reason channels were causing hanging
		}
	}
	tower.mx.Unlock()
}

func (tower *Tower) ReceiveMessage(reciever abm.Agent) messages.Message {
	tower.mx.Lock()
	defer tower.mx.Unlock()
	for id, agent := range tower.agents {
		//find sender BaseAgentCore
		if id == (reciever).ID() {
			if agent.inbox.Len() > 0 {
				front := agent.inbox.Front()
				msg := (front.Value.(messages.Message)) //<-agent.inbox:
				(agent.inbox).Remove(front)
				return msg
			}
		}
	}
	return nil
}
