package agent1

import (
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
