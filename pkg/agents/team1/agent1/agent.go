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
	if a.Base.Exists() {
		log.Printf("Custom agent %s in team 1 has floor: %d", a.ID(), a.Floor())
	} else {
		log.Printf("Agent %s No longer exists", a.ID())
	}
}

func (a *CustomAgent1) HP() {
	// log.Printf("Custom agent team 1 has floor: %d", a.Floor())
}
