package team7agent4

import (
	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
)

type CustomAgent1 struct {
	*infra.Base
	// new params
	prevFloors        map[int]int
	orderPrevFloors   []int
	currentDayonFloor int
	reshuffle         int
	trustOfAgents     map[string]int
	agentTrustLevel   int
	lastHP            int
}

func New(baseAgent *infra.Base) (infra.Agent, error) {
	//create other parameters
	return &CustomAgent1{
		Base:              baseAgent,
		currentDayonFloor: 1,
		lastHP:            100,
	}, nil
}

func (a *CustomAgent1) Run() {
	a.Log("Team 7 Agent 4 reporting status:", infra.Fields{"floor": a.Floor(), "hp": a.HP()})

	// UserID := a.ID()
	currentHP := a.HP()
	currentAvailFood := a.CurrPlatFood()

	if currentHP < a.lastHP {
		a.lastHP = currentHP
		a.currentDayonFloor++
	}

	currentFloor := a.Floor()
	if a.orderPrevFloors[len(a.orderPrevFloors)-1] != currentFloor {
		a.orderPrevFloors = append(a.orderPrevFloors, currentFloor)
	}

	if a.currentDayonFloor == 0 {

	}

	//BUG: if on reshuffle end up on same floor code breaks, reshuffle day calculation wrong, average wrong
	if a.prevFloors[currentFloor] == 0 {
		a.prevFloors[currentFloor] = int(currentAvailFood)
	} else {
		temp := (a.currentDayonFloor - 1) * a.prevFloors[currentFloor]
		temp += int(currentAvailFood)
		temp = temp / a.currentDayonFloor
		a.prevFloors[currentFloor] = temp
	}

}
