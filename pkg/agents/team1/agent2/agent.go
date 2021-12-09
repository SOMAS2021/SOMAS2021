package agent2

import (
	"log"

	infra "github.com/SOMAS2021/SOMAS2021/pkg/infra"
	abm "github.com/SOMAS2021/SOMAS2021/pkg/utils/abm"
)

type CustomAgent2 struct {
	*infra.Base
	myNumber int
	// new params
}

func New(baseAgent *infra.Base) (abm.Agent, error) {
	//create other parameters
	return &CustomAgent2{
		Base:     baseAgent,
		myNumber: 0,
	}, nil
}

func (a *CustomAgent2) Run() {
	log.Printf("Custom agent in team 2 is on floor %d has hp: %d", a.Floor(), a.HP())
	a.TakeFood(15)
}
