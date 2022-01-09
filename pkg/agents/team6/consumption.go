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

// Implicitly defined tiers based on agent HP
type levelsData struct {
	strongLevel  int
	healthyLevel int
	weakLevel    int
	critLevel    int
}

// Desired food intake without any constraint coming from messages/treaties
func (a *CustomAgent6) desiredFoodIntake() food.FoodType {
	healthInfo := a.HealthInfo()

	thresholds := thresholdData{
		maxIntake: food.FoodType(80),
	}

	// healthInfo.maxIntake

	levels := levelsData{
		strongLevel:  healthInfo.MaxHP * 3 / 5,
		healthyLevel: healthInfo.MaxHP * 3 / 10,
		weakLevel:    healthInfo.WeakLevel,
		critLevel:    0,
	}

	currentHP := a.HP()

	switch a.currBehaviour.string() {
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
			return a.foodRequired(currentHP, healthInfo)
		default:
			return a.foodRequired(levels.healthyLevel, healthInfo)
		}

	case "Narcissist": // Eat max intake
		return thresholds.maxIntake

	default:
		return food.FoodType(0)
	}
}

// Calculates food needed to get from current HP to a goal HP after HP decay. Based on health function
// Should be provided in health.go
func (a *CustomAgent6) foodRequired(goalHP int, healthInfo *health.HealthInfo) food.FoodType {
	denom := healthInfo.Width - float64(goalHP) + (1-healthInfo.HPLossSlope)*float64(a.HP()) - float64(healthInfo.HPLossBase) + healthInfo.HPLossSlope*float64(healthInfo.WeakLevel)
	return food.FoodType(healthInfo.Tau * math.Log(healthInfo.Width/denom))
}

// Computes maximum allowable food based on currently active treaties and messages
func (a *CustomAgent6) maxAllowedFood() food.FoodType {
	max := a.CurrPlatFood() //maximum value to indicate no maximum

	// Iterate through ActiveTreaties
	for _, treaty := range a.ActiveTreaties() {
		// Only consider treaties who's condition applies to our state
		if a.conditionApplies(&treaty) {
			// Convert LeaveFoodAmount and LeavePercentFood to an equivalent takeFood value
			takeFoodAmount := a.convertToTakeFoodAmount(float64(a.CurrPlatFood()), treaty.Request(), treaty.RequestValue()) - 1 // -1 to make sure GT is fulfilled

			if takeFoodAmount <= max {
				max = takeFoodAmount
			}
		}

	}

	// Check the RequestLeaveFood message
	if a.reqLeaveFoodAmount != -1 {
		takeFoodAmount := a.convertToTakeFoodAmount(float64(a.CurrPlatFood()), messages.LeaveAmountFood, a.reqLeaveFoodAmount) - 1
		if takeFoodAmount <= max {
			max = takeFoodAmount
		}
	}

	return max
}

// Returns how much food the agent intends on consuming
// Takes minimum between desired and maximum allowed food intake
func (a *CustomAgent6) intendedFoodIntake() food.FoodType {

	desiredFoodIntake := a.desiredFoodIntake()
	maxAllowedFood := a.maxAllowedFood()

	if maxAllowedFood >= desiredFoodIntake {
		return desiredFoodIntake
	} else {
		return maxAllowedFood
	}
}

// Rolling average of food taken
// Significance of past vs. present represented in prevFoodDiscount
func (a *CustomAgent6) updateAverageIntake(foodTaken food.FoodType) {
	a.averageFoodIntake = (a.config.prevFoodDiscount * float64(foodTaken)) + (1.0-a.config.prevFoodDiscount)*a.averageFoodIntake
}
