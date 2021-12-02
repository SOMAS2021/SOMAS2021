package agent

import (
	"log"

	. "github.com/SOMAS2021/SOMAS2021/pkg/agents/default"
)

type CustomAgent1 struct {
	*BaseAgent
	myNumber int
}

func New(baseAgent *BaseAgent) (Agent, error) {
	return &CustomAgent1{
		BaseAgent: baseAgent,
		myNumber:  0,
	}, nil
}

func (a *CustomAgent1) Run() {
	log.Printf("Custom agent team 1 has floor: %d", a.GetFloor())
}
