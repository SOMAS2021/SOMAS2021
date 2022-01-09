package team2

import (
	"math/rand"

	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
)

func InitActionSpace(actionDim int) []food.FoodType {
	actionSpace := make([]food.FoodType, actionDim)
	actionSpace[0] = 5
	for i := 1; i < actionDim; i++ {
		actionSpace[i] = actionSpace[i-1] + 5
	}
	return actionSpace
}

//select action according to the policies
func (a *CustomAgent2) SelectAction() int {
	//probability density function
	pdf := a.policies[a.CheckState()]
	//convert to cumulative distribution function
	cdf := make([]float64, len(pdf))
	cdf[0] = pdf[0]
	for i := 1; i < len(cdf); i++ {
		cdf[i] = cdf[i-1] + pdf[i]
	}
	//select action with given cdf
	r := rand.Float64()
	action := 0
	for r > cdf[action] {
		action++
	}
	return action
}
