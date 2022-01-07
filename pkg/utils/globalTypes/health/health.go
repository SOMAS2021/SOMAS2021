package health

import (
	"math"

	"github.com/SOMAS2021/SOMAS2021/pkg/config"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
)

type HealthInfo struct {
	// Maximum HP
	MaxHP int
	// Border between critical and weak state
	WeakLevel int
	// Parameters of the updateHP function
	Width float64
	Tau   float64
	// HP required to leave the critical level and reach the weak level
	HPReqCToW int
	// HP value attributed in the critical level
	HPCritical int
	// Number of days an agent can stay critical before dying
	MaxDayCritical int
	// HP loss when currentHP = weakLevel
	HPLossBase int
	// HP loss slope w.r.t currentHP - weakLevel
	HPLossSlope float64
	// Maximal amount of food on the platform
	maxPlatFood food.FoodType
	// Maximal amount of food the agents can physically take from the platform
	maxFoodIntake food.FoodType
}

func NewHealthInfo(parameters *config.ConfigParameters) *HealthInfo {
	return &HealthInfo{
		MaxHP:          parameters.MaxHP,
		WeakLevel:      parameters.WeakLevel,
		Width:          parameters.Width,
		Tau:            parameters.Tau,
		HPReqCToW:      parameters.HpReqCToW,
		HPCritical:     parameters.HpCritical,
		MaxDayCritical: parameters.MaxDayCritical,
		HPLossBase:     parameters.HPLossBase,
		HPLossSlope:    parameters.HPLossSlope,
		maxPlatFood:    parameters.FoodOnPlatform,
		maxFoodIntake:  parameters.MaxFoodIntake,
	}
}

// This function takes in an initial HP and a goal HP, and returns how much food you should eat to
// achieve the goal HP, while taking into account HP decay.
// This function should be called only for HP values above WeakLevel.
// If the difference between currentHP and goalHP is too large, returns the maximum food value
func FoodRequired(currentHP int, goalHP int, healthInfo *HealthInfo) food.FoodType {
	slope := healthInfo.HPLossSlope
	base := healthInfo.HPLossBase
	numer := float64(goalHP) - (1-slope)*float64(currentHP) + float64(base) - slope*float64(healthInfo.WeakLevel)
	denom := healthInfo.Width * (1 - slope)
	food := food.FoodType(math.Round(-1 * healthInfo.Tau * math.Log(1-numer/denom)))
	if food < 0 {
		if goalHP < currentHP {
			return 0
		}
		return healthInfo.maxPlatFood
	}
	return food
}
