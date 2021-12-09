package agent4_1

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
	log.Printf("Custom agent 1 in team 4 is on floor %d has hp: %d", a.Floor(), a.HP())
	a.TakeFood(10)
}
