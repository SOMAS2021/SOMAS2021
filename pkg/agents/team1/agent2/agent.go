package agent

import (
	"log"

	"github.com/SOMAS2021/SOMAS2021/pkg/agents"
)

type CustomAgent2 struct {
	*agents.Base
	myNumber int
}

func New(baseAgent *agents.Base) (agents.Agent, error) {
	return &CustomAgent2{
		Base:     baseAgent,
		myNumber: 0,
	}, nil
}

func (a *CustomAgent2) Run() {
	log.Printf("Custom agent team 2 has floor: %d", a.Floor())
}
