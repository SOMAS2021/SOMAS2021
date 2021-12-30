package team6

import (
	"math/rand"

	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/health"
)

type thresholdData struct {
	maxIntake food.FoodType
}

type levelsData struct { // tiers of HP
	strongLevel  float64
	healthyLevel float64
	weakLevel    float64
	critLevel    float64
}

func (a *CustomAgent6) foodIntake() food.FoodType {
	healthInfo := a.HealthInfo()

	thresholds := thresholdData{
		maxIntake: food.FoodType(80),
	}

	levels := levelsData{
		strongLevel:  0.6 * float64(healthInfo.MaxHP),
		healthyLevel: 0.3 * float64(healthInfo.MaxHP),
		weakLevel:    0.1 * float64(healthInfo.MaxHP),
		critLevel:    0.0,
	}

	currentHP := float64(a.HP())

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
			return health.FoodRequired(int(currentHP), int(currentHP), healthInfo)
		default:
			return health.FoodRequired(int(currentHP), int(levels.healthyLevel), healthInfo)
		}

	case "Narcissist": // Eat max intake (Possible TODO: Stay in strong instead?)
		return thresholds.maxIntake

	default:
		return food.FoodType(0)
	}
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
