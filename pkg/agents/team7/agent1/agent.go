package team7agent1

import (
	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/agent"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
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

	currentHP := a.HP()
	//currentFloor := a.Floor()
	//UserID := a.ID()
	//currentAvailFood := a.CurrPlatFood()

	var foodtotake food.FoodType = food.FoodType(100 - currentHP)
	// var food food.FoodType = food.FoodType(foodtotake)
	if foodtotake == 0 {

	} else {
		a.TakeFood(foodtotake)
	}

}
