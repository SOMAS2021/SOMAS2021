package team6

import (
	"math"
	"math/rand"

	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/health"
)

type thresholdData struct {
	maxIntake float64
}

type levelsData struct {
	strongLevel  float64
	healthyLevel float64
	weakLevel    float64
	critLevel    float64
}

// type widthData struct {
// 	strong2max     float64
// 	healthy2strong float64
// 	weak2healthy   float64
// 	crit2weak      float64
// }

var ctr int = 0
var foodTakeDay int

func (a *CustomAgent6) foodIntake() float64 {

	// towerInfo := a.Tower()
	healthInfo := a.HealthInfo()

	thresholds := thresholdData{
		maxIntake: 80.0,
	}

	levels := levelsData{
		strongLevel:  0.6 * float64(healthInfo.MaxHP),
		healthyLevel: 0.3 * float64(healthInfo.MaxHP),
		weakLevel:    0.1 * float64(healthInfo.MaxHP),
		critLevel:    0.0,
	}

	// widths := widthData{
	// 	strong2max:     float64(healthInfo.MaxHP) - levels.strongLevel,
	// 	healthy2strong: levels.strongLevel - levels.healthyLevel,
	// 	weak2healthy:   levels.healthyLevel - levels.weakLevel,
	// 	crit2weak:      levels.weakLevel,
	// }

	hp := float64(a.HP())

	switch a.currBehaviour.String() {
	case "Altruist": // Never eat
		return 0.0

	case "Collectivist": // Only eat when in critical zone on Day 3
		switch {
		case hp >= levels.weakLevel:
			ctr = 0
			foodTakeDay = rand.Intn(healthInfo.MaxDayCritical) + 1
			return 0.0
		case hp >= levels.critLevel:
			ctr = ctr + 1
			if ctr == foodTakeDay {
				return float64(healthInfo.HPReqCToW)
			}
			return 0.0
		default:
			return 0.0
		}

	case "Selfish": // Stay in Healthy zone
		switch {
		case hp >= levels.strongLevel:
			return 0.0
		case hp >= levels.healthyLevel:
			return foodRequired(hp, hp, healthInfo)
		default:
			return foodRequired(hp, levels.healthyLevel, healthInfo)
		}

	case "Narcissist": // Eat max intake (Later development: stay in Strong zone?)
		return thresholds.maxIntake

	default:
		return 0.0
	}
}

func foodRequired(currentHP, goalHP float64, healthInfo *health.HealthInfo) float64 {
	return healthInfo.Tau*math.Log(healthInfo.Width) - healthInfo.Tau*math.Log(healthInfo.Width-goalHP+0.75*currentHP-10+0.25*float64(healthInfo.WeakLevel))
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
