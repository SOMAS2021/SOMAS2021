package agent

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
	log.Printf("Custom agent team 1 has floor: %d", a.Floor())
}
