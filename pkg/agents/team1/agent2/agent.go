package agent2

import (
	"log"

	"github.com/SOMAS2021/SOMAS2021/pkg/agents"
)

type CustomAgent2 struct {
	*agents.Base
	myNumber int
	// new params
}

func New(baseAgent *agents.Base) (agents.Agent, error) {
	//create other parameters
	return &CustomAgent2{
		Base:     baseAgent,
		myNumber: 0,
	}, nil
}

func (a *CustomAgent2) Run() {
	if a.Base.Exists() {
		log.Printf("Custom agent %s in team 2 has floor: %d", a.ID(), a.Floor())
	} else {
		log.Printf("Agent %s No longer exists", a.ID())
	}
}

func (a *CustomAgent2) HP() int {
	// test
	return a.Base.HP()
}
