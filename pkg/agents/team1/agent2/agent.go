package agent2

import (
	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/abm"
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
	a.Log("Custom agent reporting status without using fields")
	a.TakeFood(15)
}
