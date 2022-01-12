package team2

import (
	"math"

	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/health"
)

func calcReward(HP int, foodTaken int, DaysAtCritical int, healthInfo *health.HealthInfo) float64 {
	ret := 0.0
	overEating := 80
	if HP == healthInfo.HPCritical {
		ret = ret - 5*float64(DaysAtCritical)
	} else if HP >= overEating {
		ret = ret - float64(foodTaken)*0.2
	} else {
		ret = 3
	}

	return ret
}

func (a *CustomAgent2) updateRTable(oldHP int, hpInc int, foodTaken int, state int, action int) {
	reward := calcReward(a.HP(), foodTaken, a.DaysAtCritical(), a.HealthInfo())
	a.rTable[state][action] = reward
	a.cumulativeRewards = a.cumulativeRewards + reward
	target := []float64{a.cumulativeRewards}
	r := []float64{reward}
	a.writeToCSV(r, "csv_files/reward_vs_days", 1)
	a.writeToCSV(target, "csv_files/cumulative_rewards", 1)
}

func ExpectedHPInc(foodTaken int, healthInfo *health.HealthInfo) int {
	return int(math.Round(healthInfo.Width * (1 - math.Pow(math.E, -float64(foodTaken)/healthInfo.Tau))))
}
