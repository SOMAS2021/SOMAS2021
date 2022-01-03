package team3

import (
	"math"

	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
)

//To Do:
//1. Remove int() from platFood once platFood is type int

// Calculates the food required to go from the current HP of an agent to the specified targetHP argument.
func (a *CustomAgent3) foodReqCalc(targetHP int) int {
	decayCurrentHP := int((float64(targetHP+a.HealthInfo().HPLossBase)-float64(a.HealthInfo().WeakLevel)*a.HealthInfo().HPLossSlope)/(1.0-a.HealthInfo().HPLossSlope) + 1)
	return int(-a.HealthInfo().Tau * math.Log(1.0-(float64(decayCurrentHP-a.HP())/a.HealthInfo().Width)))
}

// This function uses current HP of agents. This is a hunger equation.
// 10hp -> target is 20hp, 65hp -> target is 65hp, 100hp -> target is 93hp.
func (a *CustomAgent3) targetHPCalc() int {
	grad := 9.0 / 11.0
	c := 12.0
	return int(grad*float64(a.HP()) + c)
}

// Returns a range from min to max, depending on the value of the metric
func foodRange(metric int, min int, max int) int {
	if max > min && max >= 0 && min >= 0 {
		return int(float64(max) - math.Round((float64(metric) / (100.0 / float64(max-min)))))
	}
	return -1
}

// Returns a range from number*minScale to number*maxScale, depending on the value of the metric
func foodScale(number int, metric int, minScale float64, maxScale float64) int {
	if maxScale > minScale && maxScale >= 0 && minScale >= 0 {
		return int(float64(number)*minScale*float64(metric)/100.0 + float64(number)*maxScale*float64(100-metric)/100.0)
	}
	return -1
}

func OpSolver(r int, l int, op messages.Op) bool {
	switch op {
	case (messages.GT):
		return r > l
	case (messages.GE):
		return r >= l
	case (messages.EQ):
		return r == l
	case (messages.LE):
		return r <= l
	default: //LT
		return r <= l
	}
}

func (a *CustomAgent3) conditionMet(tr messages.Treaty) bool {
	conditionMet := false
	switch tr.Condition() {
	case (messages.HP):
		conditionMet = OpSolver(a.HP(), tr.ConditionValue(), tr.ConditionOp())
	case (messages.Floor):
		conditionMet = OpSolver(a.BaseAgent().Floor(), tr.ConditionValue(), tr.ConditionOp())
	case (messages.AvailableFood):
		conditionMet = OpSolver(int(a.Base.CurrPlatFood()), tr.ConditionValue(), tr.ConditionOp())
	}
	return conditionMet
}

func (a *CustomAgent3) handleTreaties() {

	for _, tr := range a.BaseAgent().ActiveTreaties() {
		condMet := a.conditionMet(tr)

		if condMet {
			switch tr.Request() {
			case (messages.LeaveAmountFood):
				a.decisions.foodToLeave = tr.RequestValue() //would sb request for us to leave less than certain food?
			case (messages.LeavePercentFood):
				//Is leave Percent Food
			case (messages.Inform):
				//Send the message they require
			}
		}

	}
}

func (a *CustomAgent3) takeFoodCalculation() int {

	platFood := int(a.CurrPlatFood())
	morality := a.vars.morality
	foodToEat := a.decisions.foodToEat
	foodToLeave := a.decisions.foodToLeave

	switch foodToEat | foodToLeave {

	case (-1): // (-1 | -1), uses HP and morality

		switch hp := a.HP(); {
		case hp == a.HealthInfo().HPCritical:
			survivalFood := a.foodReqCalc(a.HealthInfo().HPReqCToW)
			if a.DaysAtCritical() < 3 { //depend on morality
				return foodRange(morality, a.DaysAtCritical()-1, survivalFood+a.DaysAtCritical()) // adaptive range
			} else { // a.DaysAtCritical() == 3
				return foodRange(morality, survivalFood, survivalFood+a.DaysAtCritical()) //range ensures survival if possible with foodToLeave
			}
		default: // 10 <= hp <= 100:
			targetHP := a.targetHPCalc()
			foodRequired := a.foodReqCalc(targetHP)
			return foodScale(foodRequired, morality, 0.0, 2.0) // from foodRequired*0 (morality 100) to foodRequired*2 (morality 0)
		}

	case (-1 | foodToLeave): //uses platFood, HP, morality and foodToLeave

		//Any Hp
		targetHP := a.targetHPCalc()
		foodRequired := a.foodReqCalc(targetHP)
		foodToEat = foodScale(foodRequired, morality, 0.0, 2.0) // from foodRequired*0 (morality 100) to foodRequired*2 (morality 0)

		if platFood-foodToEat >= foodToLeave {
			if a.HP() == a.HealthInfo().HPCritical {
				survivalFood := a.foodReqCalc(a.HealthInfo().HPReqCToW)
				if a.DaysAtCritical() < 3 { //depend on morality
					foodToEat = foodRange(morality, a.DaysAtCritical()-1, survivalFood+a.DaysAtCritical()) // adaptive range
				} else { // a.DaysAtCritical() == 3
					foodToEat = foodRange(morality, survivalFood, survivalFood+a.DaysAtCritical()) //range ensures survival if possible with foodToLeave
				}
				if platFood-foodToEat >= foodToLeave {
					return foodToEat
				}
				return platFood - foodToLeave
			}
			return foodToEat

		} else if platFood >= foodToLeave {
			return platFood - foodToLeave
		} // Cannot satisfy platform requirement
		return 0

	case (foodToEat | -1): //uses foodToEat
		return foodToEat

	default: //uses foodToEat and foodToLeave

		if platFood-foodToLeave >= foodToEat {
			return foodToEat
		} else if platFood >= foodToLeave { //e.g. platFood=15, foodToLeave=10, foodToEat=10 --> if morality=100, take 5, if morality=0, take 10
			return foodRange(morality, platFood-foodToLeave, foodToEat)
		} //e.g. platFood=300, foodToLeave=310, foodToEat=10, take less if moral since you "should" leave a high amount.
		return foodRange(morality, 0, foodToEat) // Cannot satisfy platform requirement

	}
}
