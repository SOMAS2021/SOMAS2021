package team3

import (
	"math"
)

//To Do:
//1. Remove int() from platFood once platFood is type int
//2. Change last resort critical state to "if a.HP() == 1 && a.daysInCriticalState == 2" once daysInCriticalState is implemented
//3. Maybe improve foodCalc depending on the performance of the current version (S curve maybe). Provide multiple versions for different cases?

func takeFoodCalculation(a *CustomAgent3) int {

	platFood := int(a.CurrPlatFood())

	switch a.decisions.foodToEat | a.decisions.foodToLeave {
	case (-1 | -1): //uses platFood, HP and morality

		//TO DO: Modify for the new HP functions and requirements
		if a.HP() == 1 { //if critical, calculates to take minimum 2 food if possible. Full range = 2-10
			return int(10.0 - math.Floor((float64(a.vars.morality) / 25.0)) - math.Floor((float64(a.vars.mood))/25.0))
		}

		foodCalc := int(10.0 - math.Floor((float64(a.vars.morality) / 20.0)) - math.Floor((float64(a.vars.mood))/20.0))

		if platFood >= foodCalc {
			return foodCalc
		} else {
			return platFood //if platfood smaller than foodCalc eat all thats left
		}

	case (-1 | a.decisions.foodToLeave): //uses platFood, HP, morality and foodToLeave

		if a.HP() == 1 { //if critical, calculates to take minimum 2 food if possible. Full range = 2-10
			return int(10.0 - math.Floor((float64(a.vars.morality) / 25.0)) - math.Floor((float64(a.vars.mood))/25.0))
		}

		foodCalc := int(10.0 - math.Floor((float64(a.vars.morality) / 20.0)) - math.Floor((float64(a.vars.mood))/20.0)) //adjust to latest values

		if platFood-a.decisions.foodToLeave >= foodCalc {
			return foodCalc
		} else if platFood-a.decisions.foodToLeave >= 0 { //e.g. platFood=15, foodToLeave=10, foodCalc=10 --> if morality=100, take 5, if morality=0, take 10
			return int(math.Floor(float64(platFood-a.decisions.foodToLeave)*float64(a.vars.morality)/100.0 + float64(foodCalc)*(100.0-float64(a.vars.morality))/100.0))
		} //e.g. platFood=300, foodToLeave=310, foodCalc=10
		return int(math.Floor(float64(foodCalc) * (100.0 - float64(a.vars.morality)) / 100.0))

	case (a.decisions.foodToEat | -1): //uses foodToEat
		return a.decisions.foodToEat

	default: //uses foodToEat and foodToLeave
		if platFood-a.decisions.foodToLeave >= a.decisions.foodToEat {
			return a.decisions.foodToEat
		} else if platFood-a.decisions.foodToLeave >= 0 { //e.g. platFood=15, foodToLeave=10, foodCalc=10 --> if morality=100, take 5, if morality=0, take 10
			return int(math.Floor(float64(platFood-a.decisions.foodToLeave)*float64(a.vars.morality)/100.0 + float64(a.decisions.foodToEat)*(100.0-float64(a.vars.morality))/100.0))
		} //e.g. platFood=300, foodToLeave=310, foodCalc=10
		return int(math.Floor(float64(a.decisions.foodToEat) * (100.0 - float64(a.vars.morality)) / 100.0))
	}
}
