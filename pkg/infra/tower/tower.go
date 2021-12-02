package Tower

import (
	"log"
	"sync"
	messages "github.com/SOMAS2021/SOMAS2021/pkg/infra/messages"
	"github.com/divan/goabm/abm"
)

type Tower struct {
	FoodOnPlatform  float64
	FloorOfPlatform uint64
	mx              sync.RWMutex
	agents          []abm.Agent
	AgentCount      int
}

func New(foodOnPlatform float64, floorOfPlatform uint64, agentCount int) *Tower {
	t := &Tower{
		FoodOnPlatform:  foodOnPlatform,
		FloorOfPlatform: floorOfPlatform,
		AgentCount:      agentCount,
	}

	t.initAgents()

	return t
}

func (t *Tower) initAgents() {
	t.agents = make([]abm.Agent, t.AgentCount)
}

func (t *Tower) Tick() {
	t.mx.RLock()
	defer t.mx.RUnlock()
	log.Printf("A log from the tower!")
}

func (t *Tower) SetAgent(index int, agent abm.Agent) {
	t.mx.Lock()
	t.agents[index] = agent
	t.mx.Unlock()
}

func (tower *Tower) sendMessage(direction uint64, sender abm.Agent, msg messages.Message){
    if abs(direction) >= 1 {break}
    senderIndex = tower.agents.find(sender)
    reciever = tower.agents[senderIndex + direction]
    
	go func(reciever baseAgent) {
		reciever.inbox <- msg
	}(reciever)
    

}
