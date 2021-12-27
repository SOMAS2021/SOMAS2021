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
	a.Log("Team7Agent1 reporting status:", infra.Fields{"floor": a.Floor(), "hp": a.HP()})

	//UserID := a.ID()
	currentHP := a.HP()
	//currentFloor := a.Floor()
	//currentAvailFood := a.CurrPlatFood()

	var foodtotake food.FoodType = food.FoodType(100 - currentHP)
	
	//Only call take food if you want food
	if foodtotake != 0 {
		a.TakeFood(foodtotake)
	}

}
