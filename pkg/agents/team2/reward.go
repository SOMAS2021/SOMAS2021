package team2

import (
	"math"

	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/health"
)

func calcReward(oldHP int, hpInc int, foodIntended int, foodTaken int, DaysAtCritical int, neighbourHP int, healthInfo *health.HealthInfo) float64 {
	surviveBonus := 0.0
	eatingBonus := 0.0
	wastingBonus := 0.0
	savingBonus := 0.0

	//we encourage agent to survive
	if DaysAtCritical == 0 {
		surviveBonus += 1.0
	} else {
		surviveBonus -= 3.0 * float64(DaysAtCritical)
	}
	if oldHP == healthInfo.HPCritical {
		surviveBonus += 5.0 * float64(DaysAtCritical)
	}

	//We encourage agent to eat when weak
	if oldHP <= healthInfo.WeakLevel {
		eatingBonus += 0.01 * float64(foodTaken)
	}

	//We penalise for wanting to waste food
	wastingBonus -= 0.2 * float64(ExpectedHPInc(foodIntended, healthInfo)-hpInc)
	//We penalise for wasting food
	wastingBonus -= 0.2 * float64(ExpectedHPInc(foodTaken, healthInfo)-hpInc)

	//We reward agent when neighbour is not in critical state
	if neighbourHP == healthInfo.HPCritical {
		savingBonus -= 3.0
	} else {
		savingBonus += 1.0
	}
	return surviveBonus + eatingBonus + wastingBonus + savingBonus
}

/*
#another simpler reward function
func calcReward(HP int, foodTaken int, DaysAtCritical int, healthInfo *health.HealthInfo) float64 {
	ret := 0.0
	overEating := 80
	if HP == healthInfo.HPCritical {
		ret = ret - 1.5*float64(DaysAtCritical)
	}

	if HP >= overEating {
		ret = ret - float64(foodTaken)*0.1
	} else {
		ret = 3
	}

	return ret
}
*/

func (a *CustomAgent2) updateRTable(oldHP int, hpInc int, state int, action int) {
	reward := calcReward(oldHP, hpInc, action*5, int(a.lastFoodTaken), a.DaysAtCritical(), a.neiboughHP, a.HealthInfo())
	//reward := calcReward(a.HP(), int(a.lastFoodTaken), a.DaysAtCritical(), a.HealthInfo())
	a.rTable[state][action] = reward
	a.cumulativeRewards = a.cumulativeRewards + reward
	target := []float64{a.cumulativeRewards}
	a.writeToCSV(target, "cumulative_rewards", 1)
}

func ExpectedHPInc(foodTaken int, healthInfo *health.HealthInfo) int {
	return int(math.Round(healthInfo.Width * (1 - math.Pow(math.E, -float64(foodTaken)/healthInfo.Tau))))
}
