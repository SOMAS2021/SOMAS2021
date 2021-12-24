package team6

import "math"

type thresholdData struct {
	satisficeThresh float64
	satisfyThresh   float64
	maxIntake       float64
}

type levelsData struct {
	strongLevel  float64
	healthyLevel float64
	weakLevel    float64
	critLevel    float64
}

type widthData struct {
	strong2max     float64
	healthy2strong float64
	weak2healthy   float64
	crit2weak      float64
}

var tau float64 = 10.0
var hp float64
var ctr int = 0

func (a *CustomAgent6) foodIntake() float64 {
	thresholds := thresholdData{
		satisficeThresh: 20.0,
		satisfyThresh:   60.0,
		maxIntake:       80.0,
	}

	levels := levelsData{
		strongLevel:  60.0,
		healthyLevel: 30.0,
		weakLevel:    10.0,
		critLevel:    0.0,
	}

	widths := widthData{
		strong2max:     40.0,
		healthy2strong: 30.0,
		weak2healthy:   20.0,
		crit2weak:      10.0,
	}

	hp = float64(a.HP())

	switch a.currBehaviour.String() {
	case "Altruist": // Never eat
		return 0.0

	case "Collectivist": // Only eat when in critical zone on Day 3
		switch {
		case hp >= levels.weakLevel:
			ctr = 0
			return 0.0
		case hp >= levels.critLevel:
			ctr = ctr + 1
			if ctr == 3 {
				return foodRequiredToAscend(hp, levels.critLevel, widths.crit2weak, tau)
				// return 2.0;
				// return a.tower.healthInfo.FoodReqCToW;
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
			return foodRequiredToStay(hp, levels.healthyLevel, widths.healthy2strong, tau)
		case hp >= levels.weakLevel:
			return foodRequiredToAscend(hp, levels.weakLevel, widths.weak2healthy, tau)
		default:
			return foodRequiredToAscend(hp, levels.critLevel, widths.crit2weak, tau)
			// return 2.0
		}

	case "Narcissist": // Eat max intake (Later development: stay in Strong zone?)
		return thresholds.maxIntake

	default:
		return 0.0
	}
}

func foodRequiredToStay(currentHP float64, currentLevel float64, levelWidth float64, tau float64) float64 {
	return tau*math.Log(currentLevel+levelWidth-currentHP) - tau*math.Log(math.Pow(math.E, -1)*levelWidth)
}

func foodRequiredToAscend(currentHP float64, currentLevel float64, levelWidth float64, tau float64) float64 {
	return tau*math.Log(currentLevel+levelWidth-currentHP) - tau*math.Log(math.Pow(math.E, -3)*levelWidth)
}

// func hpFunc(x float64, currentHP float64, currentLevel float64, levelWidth float64, tau float64) float64 {
// 	return currentLevel + levelWidth - math.Exp(-x/tau)*(levelWidth+currentHP-currentLevel)
// }
