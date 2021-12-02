package agent

import (
	"log"

	. "github.com/SOMAS2021/SOMAS2021/pkg/agents/default"
)

type CustomAgent2 struct {
	*BaseAgent
	myNumber int
}

func New(baseAgent *BaseAgent) (Agent, error) {
	return &CustomAgent2{
		BaseAgent: baseAgent,
		myNumber:  0,
	}, nil
}

func (a *CustomAgent2) Run() {
	log.Printf("Custom agent team 2 has floor: %d", a.GetFloor())
}
