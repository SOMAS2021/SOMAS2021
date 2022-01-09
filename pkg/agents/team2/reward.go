package team2

import (
	"math"

	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/health"
)

func calcReward(oldHP int, hpInc int, foodIntended int, foodTaken int, DaysAtCritical int, healthInfo *health.HealthInfo) float64 {
	surviveBonus := 0.0
	eatingBonus := 0.0
	wastingBonus := 0.0

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

	return surviveBonus + eatingBonus + wastingBonus
}

func (a *CustomAgent2) updateRTable(oldHP int, hpInc int, foodTaken int, state int, action int) {
	reward := calcReward(oldHP, hpInc, action*5, foodTaken, a.DaysAtCritical(), a.HealthInfo())
	a.rTable[state][action] = reward
	a.cumulativeRewards = a.cumulativeRewards + reward
	a.writeToCSV(a.cumulativeRewards, "cumulative_rewards")
}

func ExpectedHPInc(foodTaken int, healthInfo *health.HealthInfo) int {
	return int(math.Round(healthInfo.Width * (1 - math.Pow(math.E, -float64(foodTaken)/healthInfo.Tau))))
}
