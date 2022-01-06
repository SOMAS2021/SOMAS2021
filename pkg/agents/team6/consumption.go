package team6

import (
	"math"
	"math/rand"

	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/health"
)

type thresholdData struct {
	maxIntake food.FoodType
}

type levelsData struct { // tiers of HP
	strongLevel  int
	healthyLevel int
	weakLevel    int
	critLevel    int
}

// desired food intake without any constraint coming from messages/treaties
func (a *CustomAgent6) desiredFoodIntake() food.FoodType {
	healthInfo := a.HealthInfo()

	thresholds := thresholdData{
		maxIntake: food.FoodType(80),
	}

	levels := levelsData{
		strongLevel:  healthInfo.MaxHP * 3 / 5,
		healthyLevel: healthInfo.MaxHP * 3 / 10,
		weakLevel:    healthInfo.MaxHP * 1 / 10,
		critLevel:    0,
	}

	currentHP := a.HP()

	switch a.currBehaviour.String() {
	case "Altruist": // Never eat
		return food.FoodType(0)

	case "Collectivist": // Only eat when in critical zone randomly before expiry
		switch {
		case currentHP >= levels.weakLevel:
			a.foodTakeDay = rand.Intn(healthInfo.MaxDayCritical) // Stagger the days when agents return to weak
			return food.FoodType(0)
		case currentHP >= levels.critLevel:
			if a.DaysAtCritical() == a.foodTakeDay {
				return food.FoodType(healthInfo.HPReqCToW)
			}
			return food.FoodType(0)
		default:
			return food.FoodType(0)
		}

	case "Selfish": // Stay in Healthy zone
		switch {
		case currentHP >= levels.strongLevel:
			return food.FoodType(0)
		case currentHP >= levels.healthyLevel:
			return FoodRequired(currentHP, currentHP, healthInfo)
		default:
			return FoodRequired(currentHP, levels.healthyLevel, healthInfo)
		}

	case "Narcissist": // Eat max intake (Possible TODO: Stay in strong instead?)
		return thresholds.maxIntake

	default:
		return food.FoodType(0)
	}
}

func FoodRequired(currentHP int, goalHP int, healthInfo *health.HealthInfo) food.FoodType {
	denom := healthInfo.Width - float64(goalHP) + (1-healthInfo.HPLossSlope)*float64(currentHP) - float64(healthInfo.HPLossBase) + healthInfo.HPLossSlope*float64(healthInfo.WeakLevel)
	return food.FoodType(healthInfo.Tau * math.Log(healthInfo.Width/denom))
}

func (a *CustomAgent6) incomingTreatyValidityAndFoodRange(trty messages.Treaty) (int, int, bool) {
	chkTrtyVal := trty.RequestValue()
	mini := 0    //initial minimum value
	maxi := 0    //maximum value to indicate no maximum
	out := false //bool to indicate validity
	eqFound := 0 //is there an equal requestOp
	eqVal := 0   //equal requestOp requestValue
	noMaxi := 0  //Indicator for no maximum
	for _, value := range a.ActiveTreaties() {
		switch value.RequestOp() {
		case 0:
			if value.RequestValue() > mini {
				mini = value.RequestValue() + 1 //If Greater than, then the max value is one larger
			}
		case 1:
			if value.RequestValue() > mini {
				mini = value.RequestValue() //If Greater or Equal, then the max value is itself
			}
		case 3:
			if value.RequestValue() < maxi || maxi == 0 {
				maxi = value.RequestValue() //If Less than or Equal, then the max value is itself
			}
		case 4:
			if value.RequestValue() < maxi || maxi == 0 {
				maxi = value.RequestValue() - 1 //If Less than, then the max value is one smaller
			}
		case 2:
			eqFound = 1
			eqVal = value.RequestValue() //if Equal then find the value as there can only be one equal operator unless the values in both treaties are the same
		}
	}
	if maxi == 0 {
		maxi = chkTrtyVal
		noMaxi = 1
	}
	switch trty.RequestOp() {
	case 0:
		if chkTrtyVal > mini && chkTrtyVal < maxi {
			out = true
		}
	case 1:
		if chkTrtyVal >= mini && chkTrtyVal <= maxi {
			out = true
		}
	case 3:
		if chkTrtyVal >= mini && chkTrtyVal <= maxi {
			out = true
		}
	case 4:
		if chkTrtyVal > mini && chkTrtyVal < maxi {
			out = true
		}
	case 2:
		if (eqFound == 1) && (chkTrtyVal != eqVal) {
			out = false
		} else {
			if chkTrtyVal >= mini && chkTrtyVal <= maxi {
				out = true
			}
		}
	}
	if noMaxi == 1 {
		maxi = 0
	}
	return mini, maxi, out
}

func (a *CustomAgent6) intendedFoodIntake() food.FoodType {

	//if a.tower
	//foodOnPlatform := foodOnPlatform
	// else
	//foodOnPlatform = lastFoodOnPlatform

	intendedFoodIntake := a.desiredFoodIntake()
	if a.reqLeaveFoodAmount != -1 {
		intendedFoodIntake = food.FoodType(math.Min(float64(a.CurrPlatFood())-float64(a.reqLeaveFoodAmount), float64(intendedFoodIntake)))

		//intendedFoodIntake = food.FoodType(a.reqLeaveFoodAmount) // to correct
	}
	return intendedFoodIntake
}

// func foodRequiredToStay(currentHP float64, currentLevel float64, levelWidth float64, tau float64) float64 {
// 	return tau*math.Log(currentLevel+levelWidth-currentHP) - tau*math.Log(math.Pow(math.E, -1)*levelWidth)
// }

// func foodRequiredToAscend(currentHP float64, currentLevel float64, levelWidth float64, tau float64) float64 {
// 	return tau*math.Log(currentLevel+levelWidth-currentHP) - tau*math.Log(math.Pow(math.E, -3)*levelWidth)
// }

// func hpFunc(x float64, currentHP float64, currentLevel float64, levelWidth float64, tau float64) float64 {
// 	return currentLevel + levelWidth - math.Exp(-x/tau)*(levelWidth+currentHP-currentLevel)
// }
