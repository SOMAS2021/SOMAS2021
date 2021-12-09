package agent1

import (
	"log"

	"github.com/SOMAS2021/SOMAS2021/pkg/agents"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/abm"

)

type CustomAgent1 struct {
	*agents.Base
	myNumber int
}

func New(baseAgent *agents.Base) (abm.Agent, error) {
	return &CustomAgent1{
		Base:     baseAgent,
		myNumber: 0,
	}, nil
}

func (a *CustomAgent1) Run() {
	// log.Printf("Custom agent in team 1 is on floor %d has hp: %d", a.Floor(), a.HP())
	// a.takeFood(10)
	// a.HP()
	// a.IsDead()

	
	if a.Base.Exists() {
		log.Printf("Custom agent %s in team 1 has floor: %d", a.ID(), a.Floor())
		msg := *messages.NewAckMessage(uint(a.Floor()), 3, true)
		log.Printf("Custom agent team 1 has floor: %d", a.Floor())
		if a.Floor() < 9 {a.Base.SendMessage(1, msg)}
		// //make message
		receivedMsg := a.ReceiveMessage()
		if receivedMsg != nil {
			log.Printf("Custom agent team 1 has floor: %d, I got a msg from: %d", a.Floor(), receivedMsg.MessageType())
		}
		a.TakeFood(1)
		
	} else {
		log.Printf("Agent %s No longer exists", a.ID())
	}
	
}

func (a *CustomAgent1) HP() int {
	return a.Base.HP()
}

func (a *CustomAgent1) IsDead() bool {
	return a.Base.IsDead()
}

func (a *CustomAgent1) takeFood(foodToTake float64) float64 {
	return a.Base.TakeFood(foodToTake)
}

func (a *CustomAgent1) ID() string {
	return a.Base.ID()
}
