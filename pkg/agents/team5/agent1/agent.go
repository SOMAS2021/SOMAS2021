package team5

import (
	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/abm"
)

type CustomAgent1 struct {
	*infra.Base
}

func New(baseAgent *infra.Base) (abm.Agent, error) {
	return &CustomAgent1{
		Base: baseAgent,
	}, nil
}

func (a *CustomAgent1) Run() {
	a.Log("Reporting agent state", infra.Fields{"health": a.HP(), "floor": a.Floor()})

	a.TakeFood(20)
}
