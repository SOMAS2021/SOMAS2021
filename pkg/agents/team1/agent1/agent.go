package agent1

import (
	"log"

	"github.com/SOMAS2021/SOMAS2021/pkg/agents"
)

type CustomAgent1 struct {
	*agents.Base
	myNumber int
}

func New(baseAgent *agents.Base) (agents.Agent, error) {
	return &CustomAgent1{
		Base:     baseAgent,
		myNumber: 0,
	}, nil
}

func (a *CustomAgent1) Run() {
	log.Printf("Custom agent in team 1 is on floor %d has hp: %d", a.Floor(), a.HP())
	a.TakeFood(10)
}

func (a *CustomAgent1) AddToInbox(msg messages.Message) {
	a.Base.AddToInbox(msg)
}
