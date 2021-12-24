package team2

import "math/rand"

func InitActionSpace() actionSpace {
	//TODO: actionID might be removed in further versions
	//index : 0 => disregard food
	//index : 1 => satisfice with food
	//index : 2 => satisfy with food
	initialActionSpace := actionSpace{}
	initialActionSpace.actionId = make([]int, 3)
	for i := 0; i < 3; i++ {
		initialActionSpace.actionId[i] = i
	}
	m := map[int]func(hp int) float64{
		//actions based on the current hp level
		initialActionSpace.actionId[0]: DisFood,
		initialActionSpace.actionId[1]: Satisfice,
		initialActionSpace.actionId[2]: Satisfy,
	}
	initialActionSpace.actionSet = m
	return initialActionSpace
}

//Need to change this func when adding new actions

func DisFood(hp int) float64 {
	return 0.0
}

func Satisfice(hp int) float64 {
	if hp <= 20 { //20 is the daily decay
		return 20.0
	} else {
		return 1.0
	}
}
func Satisfy(hp int) float64 {
	return 100.0 - float64(hp)
}

//select action according to the policies
func (a *CustomAgent2) SelectAction() int {
	//probability density function
	pdf := a.policies[a.CheckState()]
	//convert to cumulative distribution function
	cdf := make([]float32, len(pdf))
	cdf[0] = pdf[0]
	for i := 1; i < len(cdf); i++ {
		cdf[i] = cdf[i-1] + pdf[i]
	}
	//select action with given cdf
	r := rand.Float32()
	action := 0
	for r > cdf[action] {
		action++
	}
	return action
}
