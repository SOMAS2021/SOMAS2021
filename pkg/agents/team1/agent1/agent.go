package agent1

import (
	infra "github.com/SOMAS2021/SOMAS2021/pkg/infra"
)

type CustomAgent1 struct {
	*infra.Base
	myNumber int
}

func New(baseAgent *infra.Base) (infra.Agent, error) {
	return &CustomAgent1{
		Base:     baseAgent,
		myNumber: 0,
	}, nil
}

func (a *CustomAgent1) Run() {
	// log.Printf("Custom agent in team 1 is on floor %d has hp: %d", a.Floor(), a.HP())
	// a.takeFood(10)
	// a.HP()
}

func (a *CustomAgent1) HP() int {
	return a.Base.HP()
}

func (a *CustomAgent1) Die() {
	a.Base.Die()
}

func (a *CustomAgent1) IsAlive() bool {
	return a.Base.IsAlive()
}

func (a *CustomAgent1) takeFood(foodToTake float64) float64 {
	return a.Base.TakeFood(foodToTake)
}
