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

func (a *CustomAgent6) foodRange() (int, int) {
	mini := 0 //initial minimum value
	maxi := 0 //maximum value to indicate no maximum
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
			eqVal := value.RequestValue() //if Equal then find the value as there can only be one equal operator unless the values in both treaties are the same
			mini, maxi = eqVal, eqVal
		default:
			mini, maxi = -1, -1 //unknown op code
		}
	}

	return mini, maxi
}

func (a *CustomAgent6) treatyValid(treaty messages.Treaty) bool {
	chkTrtyVal := treaty.RequestValue()
	mini, maxi := a.foodRange()

	switch treaty.RequestOp() {
	case 0:
		if chkTrtyVal > mini && chkTrtyVal < maxi {
			return true
		}
		return false
	case 1:
		if chkTrtyVal >= mini && chkTrtyVal <= maxi {
			return true
		}
		return false
	case 3:
		if chkTrtyVal >= mini && chkTrtyVal <= maxi {
			return true
		}
		return false
	case 4:
		if chkTrtyVal > mini && chkTrtyVal < maxi {
			return true
		}
		return false
	case 2:
		if chkTrtyVal != mini {
			return false
		}
		return true

	default:
		return false
	}
}

// func (a *CustomAgent6) incomingTreatyValidityAndFoodRange(trty messages.Treaty) (int, int, bool) {
// 	chkTrtyVal := trty.RequestValue()
// 	mini := 0    //initial minimum value
// 	maxi := 0    //maximum value to indicate no maximum
// 	out := false //bool to indicate validity
// 	eqFound := 0 //is there an equal requestOp
// 	eqVal := 0   //equal requestOp requestValue
// 	noMaxi := 0  //Indicator for no maximum
// 	for _, value := range a.ActiveTreaties() {
// 		switch value.RequestOp() {
// 		case 0:
// 			if value.RequestValue() > mini {
// 				mini = value.RequestValue() + 1 //If Greater than, then the max value is one larger
// 			}
// 		case 1:
// 			if value.RequestValue() > mini {
// 				mini = value.RequestValue() //If Greater or Equal, then the max value is itself
// 			}
// 		case 3:
// 			if value.RequestValue() < maxi || maxi == 0 {
// 				maxi = value.RequestValue() //If Less than or Equal, then the max value is itself
// 			}
// 		case 4:
// 			if value.RequestValue() < maxi || maxi == 0 {
// 				maxi = value.RequestValue() - 1 //If Less than, then the max value is one smaller
// 			}
// 		case 2:
// 			eqFound = 1
// 			eqVal = value.RequestValue() //if Equal then find the value as there can only be one equal operator unless the values in both treaties are the same
// 		}
// 	}
// 	if maxi == 0 {
// 		maxi = chkTrtyVal
// 		noMaxi = 1
// 	}
// 	switch trty.RequestOp() {
// 	case 0:
// 		if chkTrtyVal > mini && chkTrtyVal < maxi {
// 			out = true
// 		}
// 	case 1:
// 		if chkTrtyVal >= mini && chkTrtyVal <= maxi {
// 			out = true
// 		}
// 	case 3:
// 		if chkTrtyVal >= mini && chkTrtyVal <= maxi {
// 			out = true
// 		}
// 	case 4:
// 		if chkTrtyVal > mini && chkTrtyVal < maxi {
// 			out = true
// 		}
// 	case 2:
// 		if (eqFound == 1) && (chkTrtyVal != eqVal) {
// 			out = false
// 		} else {
// 			if chkTrtyVal >= mini && chkTrtyVal <= maxi {
// 				out = true
// 			}
// 		}
// 	}
// 	if noMaxi == 1 {
// 		maxi = 0
// 	}
// 	return mini, maxi, out
// }

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

func (a *CustomAgent6) updateAverageIntake(foodTaken food.FoodType) {
	a.averageFoodIntake = (a.config.prevFoodDiscount * float64(foodTaken)) + (1.0-a.config.prevFoodDiscount)*a.averageFoodIntake
}
