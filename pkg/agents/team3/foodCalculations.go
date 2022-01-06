package team3

import (
	"math"

	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
)

// Calculates the food required to go from the current HP of an agent to the specified targetHP argument.
func (a *CustomAgent3) foodReqCalc(initialHP int, targetHP int) int {
	decayCurrentHP := int((float64(targetHP+a.HealthInfo().HPLossBase)-float64(a.HealthInfo().WeakLevel)*a.HealthInfo().HPLossSlope)/(1.0-a.HealthInfo().HPLossSlope) + 1)
	return int(-a.HealthInfo().Tau * math.Log(1.0-(float64(decayCurrentHP-initialHP)/a.HealthInfo().Width)))
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
		return int(float64(max) - math.Round(float64(metric)/(100.0/float64(max-min))))
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
	case (messages.LT):
		return r < l
	default:
		return false
	}
}

//Checks wether a treaty condition is currently met
func (a *CustomAgent3) conditionMet(tr messages.Treaty) bool {
	switch tr.Condition() {
	case (messages.HP):
		return OpSolver(a.HP(), tr.ConditionValue(), tr.ConditionOp())
	case (messages.Floor):
		return OpSolver(a.BaseAgent().Floor(), tr.ConditionValue(), tr.ConditionOp())
	case (messages.AvailableFood):
		return OpSolver(int(a.Base.CurrPlatFood()), tr.ConditionValue(), tr.ConditionOp())
	default:
		return false
	}

}

func (a *CustomAgent3) handleTreaties() {
	for _, tr := range a.BaseAgent().ActiveTreaties() { //check all treaties in our memory bank

		if a.conditionMet(tr) { //if we meet the treaties condition
			switch tr.Request() { //do the requuest
			case (messages.LeaveAmountFood):
				a.decisions.foodToLeave = tr.RequestValue() //would sb request for us to leave less than certain food?
			case (messages.LeavePercentFood):
				if a.BaseAgent().CurrPlatFood() != -1 {
					a.decisions.foodToLeave = int((float64(tr.RequestValue()) / 100.0) * float64(a.BaseAgent().CurrPlatFood()))
				}
			case (messages.Inform):
				if tr.Condition() == messages.HP {
					msg := messages.NewStateHPMessage(a.BaseAgent().ID(), a.Floor(), a.Floor()-1, a.BaseAgent().HP())
					a.SendMessage(msg)
				}
			}
		}
	}
}

func (a *CustomAgent3) takeFoodCalculation() int {
	a.handleTreaties() //take into account treaties.

	platFood := int(a.CurrPlatFood())
	morality := a.vars.morality
	foodToEat := a.decisions.foodToEat
	foodToLeave := a.decisions.foodToLeave

	switch foodToEat | foodToLeave {

	case (-1): // (-1 | -1), uses HP and morality

		switch hp := a.HP(); {
		case hp == a.HealthInfo().HPCritical:

			survivalFood := a.foodReqCalc(a.HP(), a.HealthInfo().HPReqCToW)

			if a.DaysAtCritical() < 2 { //depend on morality
				//a.Log("HP is critical", infra.Fields{"hp: ": a.HP(), "Days at Crit:": a.DaysAtCritical(), "survivalFood: ": survivalFood, "food Intended: ": foodRange(morality, a.DaysAtCritical(), survivalFood+a.DaysAtCritical()+3)})
				return foodRange(morality, a.DaysAtCritical(), survivalFood+a.DaysAtCritical()) // adaptive range
			} else { // a.DaysAtCritical() == 2
				//a.Log("HP is critical", infra.Fields{"hp: ": a.HP(), "Days at Crit:": a.DaysAtCritical(), "survivalFood: ": survivalFood, "food Intended: ": foodRange(morality, survivalFood, survivalFood+a.DaysAtCritical())})
				return foodRange(morality, survivalFood, survivalFood+a.DaysAtCritical()) //range ensures survival if possible with foodToLeave
			}
		default: // 10 <= hp <= 100:
			targetHP := a.targetHPCalc()
			foodRequired := a.foodReqCalc(a.HP(), targetHP)
			//a.Log("Case -1, HP is NOT critical", infra.Fields{"HP: ": a.HP(), "targetHP: ": targetHP, "Morality: ": morality, "foodRequired: ": foodRequired, "Scaled foodRequired: ": foodScale(foodRequired, morality, 0.0, 2.0)})
			return foodScale(foodRequired, morality, 0.0, 2.0) // from foodRequired*0 (morality 100) to foodRequired*2 (morality 0)
		}

	case (-1 | foodToLeave): //uses platFood, HP, morality and foodToLeave
		//Any Hp
		targetHP := a.targetHPCalc()
		foodRequired := a.foodReqCalc(a.HP(), targetHP)
		foodToEat = foodScale(foodRequired, morality, 0.0, 2.0) // from foodRequired*0 (morality 100) to foodRequired*2 (morality 0)

		if platFood-foodToEat >= foodToLeave {
			if a.HP() == a.HealthInfo().HPCritical {
				survivalFood := a.foodReqCalc(a.HP(), a.HealthInfo().HPReqCToW)
				if a.DaysAtCritical() < 2 { //depend on morality
					foodToEat = foodRange(morality, 0, survivalFood+3) // adaptive range
				} else { // a.DaysAtCritical() == 3
					foodToEat = foodRange(morality, survivalFood, survivalFood*2) //range ensures survival if possible with foodToLeave
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
