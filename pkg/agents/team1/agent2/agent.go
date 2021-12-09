package agent2

import (
	"log"
	"github.com/SOMAS2021/SOMAS2021/pkg/agents"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/abm"

)

type CustomAgent2 struct {
	*agents.Base
	myString string
}

func New(baseAgent *agents.Base) (abm.Agent, error) {
	//create other parameters
	return &CustomAgent2{
		Base:     baseAgent,
		myString: "hello world",
	}, nil
}

func (a *CustomAgent2) Run() {
	log.Printf("Custom agent in team 2 is on floor %d has hp: %d", a.Floor(), a.HP())
	a.TakeFood(15)
}


func (a *CustomAgent2) ID() string {
	return a.Base.ID()
}
