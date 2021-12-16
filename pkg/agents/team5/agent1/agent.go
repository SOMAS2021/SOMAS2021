package team5

import (
	"math"

	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/abm"
)

type CustomAgent5 struct {
	*infra.Base
	desperation  float64
	selflishness float64
	lastMeal     float64
}

func New(baseAgent *infra.Base) (abm.Agent, error) {
	return &CustomAgent5{
		Base:         baseAgent,
		desperation:  3.0, //Scale of 1 to 4, with 1 being near max health, 4 being about to die and 1 & 2 in between
		selflishness: 2.0, // of 0 to 3, with 3 being completely selflish, 0 being completely selfless
		lastMeal:     0,   //Stores value of the last amount of food taken
	}, nil
}

func (a *CustomAgent5) UpdateDesperation() {
	switch {
	case a.HP() <= 20:
		a.desperation = 4.0
	case a.HP() <= 50:
		a.desperation = 3.0
	case a.HP() < 80:
		a.desperation = 2.0
	case a.HP() >= 80:
		a.desperation = 1.0
	}
}

func (a *CustomAgent5) UpdateSelfishness() {
	//Some function that returns a new value of selfishness based on social network
	//Will also be based on the desperation value
}

func (a *CustomAgent5) FoodAmount() float64 {
	provisional := 20.0 * a.desperation //This determines the base level of food based on desperation
	selfWeight := a.selflishness / 3.0  //This scales the taken amound by selfishness of agent
	return math.Round(provisional * selfWeight)
}

func (a *CustomAgent5) Run() {
	a.Log("Reporting agent state of team 5 agent", infra.Fields{"health": a.HP(), "floor": a.Floor()})
	a.UpdateDesperation()
	a.UpdateSelfishness()
	attemptFood := a.FoodAmount()
	if !a.HasEaten() {
		a.lastMeal = a.TakeFood(attemptFood)
		if a.lastMeal > 0 {
			a.Log("Team 5 agent has taken food", infra.Fields{"amount": a.lastMeal})
		}
	}
}
