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
	log.Printf("Custom agent %s in team 1 is on floor %d has hp: %d", a.ID(), a.Floor(), a.HP())
}

func (a *CustomAgent1) HP() int {
	return a.Base.HP()
	// log.Printf("Custom agent team 1 has floor: %d", a.Floor())
}

func (a *CustomAgent1) IsDead() bool {
	return a.Base.IsDead()
}

func (a *CustomAgent1) takeFood(foodToTake float64) float64 {
	return a.Base.TakeFood(foodToTake)
}
