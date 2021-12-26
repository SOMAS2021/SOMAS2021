package team7agent3

import (
	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/agent"
)

type CustomAgent1 struct {
	*infra.Base
	// new params
}

func New(baseAgent *infra.Base) (agent.Agent, error) {
	//create other parameters
	return &CustomAgent1{
		Base: baseAgent,
	}, nil
}

func (a *CustomAgent1) Run() {
	a.Log("Agent7 reporting status:", infra.Fields{"floor": a.Floor(), "hp": a.HP()})

	// UserID := a.ID()
	// currentHP := a.HP()

	//currentAvailFood := a.CurrPlatFood()

}
