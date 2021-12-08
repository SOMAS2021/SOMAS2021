package agent2

import (
	"github.com/SOMAS2021/SOMAS2021/pkg/agents"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/abm"

)

type CustomAgent2 struct {
	*agents.Base
	myNumber int
	// new params
}

func New(baseAgent *agents.Base) (abm.Agent, error) {
	//create other parameters
	return &CustomAgent2{
		Base:     baseAgent,
		myNumber: 0,
	}, nil
}

func (a *CustomAgent2) Run() {
	// log.Printf("Custom agent in team 2 is on floor %d has hp: %d", a.Floor(), a.HP())
	// a.takeFood(15)
	// a.HP()
	// a.IsDead()
}

func (a *CustomAgent2) HP() int {
	return a.Base.HP()
}

func (a *CustomAgent2) IsDead() bool {
	return a.Base.IsDead()
}

func (a *CustomAgent2) takeFood(foodToTake float64) float64 {
	return a.Base.TakeFood(foodToTake)
}

func (a *CustomAgent2) AddToInbox(msg messages.Message) {
	a.Base.AddToInbox(msg)
}
func (a *CustomAgent2) ID() string {
	return a.ID()
}
