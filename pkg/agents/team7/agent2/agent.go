package team7agent2

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

	//UserID := a.ID()
	currentHP := a.HP()

	H := 60 //Healthy HP
	G := 40 //Green HP
	W := 10 //Weak HP
	C := 5  //Critial HP

	F := a.Floor() //currentFloor := a.Floor()
	Y := 5
	X := 15

	//currentAvailFood := a.CurrPlatFood()

	if F < Y {
		foodtotake := food.FoodType(G - currentHP)
		if foodtotake != 0 {
			a.TakeFood(foodtotake)
		}
	} else if currentHP < C {
		foodtotake := food.FoodType(G - currentHP)
		if foodtotake != 0 {
			a.TakeFood(foodtotake)
		}
	} else {
		if F < X && (G+(F-Y) < H) {
			foodtotake := food.FoodType((G + (F - Y)) - currentHP)
			if foodtotake != 0 {
				a.TakeFood(foodtotake)
			}
		} else {
			foodtotake := food.FoodType(H - currentHP)
			if foodtotake != 0 {
				a.TakeFood(foodtotake)
			}
		}
	}

}
