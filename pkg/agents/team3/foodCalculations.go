package team3

import (
	"math"
)

//To Do:
//1. Remove int() from platFood once platFood is type int

func foodReqCalc(a *CustomAgent3, targetHP int) int {
	decayCurrentHP := int((float64(targetHP+a.HealthInfo().HPLossBase)-float64(a.HealthInfo().WeakLevel)*a.HealthInfo().HPLossSlope)/(1.0-a.HealthInfo().HPLossSlope) + 1)
	return int(-a.HealthInfo().Tau * math.Log(1.0-(float64(decayCurrentHP-a.HP())/a.HealthInfo().Width)))
}

// This function uses current HP of agents. This is a hunger equation.
func targetHPCalc(a *CustomAgent3) int { //10hp -> target is 20hp, 65hp -> target is 65hp, 100hp -> target is 93hp
	return int(12.0 + 9.0*float64(a.HP())/11.0)
}

func foodRange(metric int, min int, max int) int {
	if max > min && max >= 0 && min >= 0 {
		return int(float64(max) - math.Round((float64(metric) / (100.0 / float64(max-min)))))
	}
	return -1
}

func foodScale(number int, metric int, minScale float64, maxScale float64) int {
	if maxScale > minScale && maxScale >= 0 && minScale >= 0 {
		return int(float64(number)*minScale*float64(metric)/100.0 + float64(number)*maxScale*float64(100-metric)/100.0)
	}
	return -1
}

func takeFoodCalculation(a *CustomAgent3) int {

	platFood := int(a.CurrPlatFood())
	hpCtoW := a.HealthInfo().HPReqCToW
	morality := a.vars.morality
	foodToEat := a.decisions.foodToEat
	foodToLeave := a.decisions.foodToLeave

	switch foodToEat | foodToLeave {

	case (-1): // (-1 | -1), uses HP and morality

		switch hp := a.HP(); {
		case hp == a.HealthInfo().HPCritical:
			if a.DaysAtCritical() == 1 { //depend on morality
				return foodRange(morality, 0, hpCtoW+1) //range 0 to hpCtoW+1
			} else if a.DaysAtCritical() == 2 {
				return foodRange(morality, 1, hpCtoW+2) //range 1 to hpCtoW+2
			} else { // a.DaysAtCritical() == 3
				return foodRange(morality, hpCtoW, hpCtoW+3) //range hpCtoW to hpCtoW+3
			}
		default: // 10 <= hp <= 100:
			targetHP := targetHPCalc(a)
			foodRequired := foodReqCalc(a, targetHP)
			return foodScale(foodRequired, morality, 0.0, 2.0) // from foodRequired*0 (morality 100) to foodRequired*2 (morality 0)
		}

	case (-1 | foodToLeave): //uses platFood, HP, morality and foodToLeave

		//Any Hp
		targetHP := targetHPCalc(a)
		foodRequired := foodReqCalc(a, targetHP)
		foodToEat = foodScale(foodRequired, morality, 0.0, 2.0) // from foodRequired*0 (morality 100) to foodRequired*2 (morality 0)

		if platFood-foodToEat >= foodToLeave {
			if a.HP() == a.HealthInfo().HPCritical {
				if a.DaysAtCritical() == 1 { //depend on morality
					foodToEat = foodRange(morality, 0, hpCtoW+1) //range 0 to hpCtoW+1
				} else if a.DaysAtCritical() == 2 {
					foodToEat = foodRange(morality, 1, hpCtoW+2) //range 1 to hpCtoW+2
				} else { // a.DaysAtCritical() == 3
					foodToEat = foodRange(morality, hpCtoW, hpCtoW+3) //range hpCtoW to hpCtoW+3
				}
				if platFood-foodToEat >= foodToLeave {
					return foodToEat
				} else {
					return platFood - foodToLeave
				}
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
