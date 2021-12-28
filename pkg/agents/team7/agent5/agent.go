package team7agent5

import (
	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
)

type CustomAgent1 struct {
	*infra.Base
	// new params
	prevFloors      []int
	daysElapsed     int
	trustOfAgents   map[string]int
	agentTrustLevel int
}

func New(baseAgent *infra.Base) (infra.Agent, error) {
	//create other parameters
	return &CustomAgent1{
		Base: baseAgent,
	}, nil
}

func (a *CustomAgent1) Run() {
	a.Log("Team 7 Agent 5 reporting status:", infra.Fields{"floor": a.Floor(), "hp": a.HP()})

	// UserID := a.ID()
	// currentHP := a.HP()
	// currentAvailFood := a.CurrPlatFood()

	currentFloor := a.Floor()
	a.prevFloors = append(a.prevFloors, currentFloor)

}
