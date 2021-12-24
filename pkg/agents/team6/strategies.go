package team6

import "github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"

type thresholdData struct {
	satisficeThresh food.FoodType
	satisfyThresh   food.FoodType
	maxIntake       food.FoodType
}

func (a *CustomAgent6) foodIntake() food.FoodType {
	thresholds := thresholdData{satisficeThresh: food.FoodType(20), satisfyThresh: food.FoodType(60), maxIntake: food.FoodType(80)}

	switch a.currBehaviour.String() {
	case "Altruist":
		return food.FoodType(0)
	case "Collectivist":
		return thresholds.satisficeThresh
	case "Selfish":
		return thresholds.satisfyThresh
	case "Narcissist":
		return thresholds.maxIntake
	default:
		return food.FoodType(0)
	}
}
