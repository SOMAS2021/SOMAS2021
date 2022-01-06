package team6

import (
	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
)

func (a *CustomAgent6) addToMemory() {

	if a.PlatformOnFloor() && a.platOnFloorCtr == 0 {

		currPlatFood := a.CurrPlatFood()
		if currPlatFood == -1 { // Infra returns -1 if a.CurrPlatFood() is called and the platform is not at the agent's floor
			return
		}
		a.longTermMemory = append(a.longTermMemory, currPlatFood)
		a.shortTermMemory = append(a.shortTermMemory, currPlatFood)

		a.Log("Team 6 age:", infra.Fields{"Age": a.Age()})
		a.Log("Team 6 food available:", infra.Fields{"CurrPlatFood": a.CurrPlatFood()})
		a.Log("Team 6 long term memory:", infra.Fields{"Lmemory": a.longTermMemory})
		a.Log("Team 6 short term memory:", infra.Fields{"Smemory": a.shortTermMemory})

		a.platOnFloorCtr++
	} else if !a.PlatformOnFloor() {
		a.platOnFloorCtr = 0
	}

}

func (a *CustomAgent6) isReassigned() bool {
	if a.prevFloor != a.Floor() && a.prevFloor != -1 {
		return true
	}
	return false
}

func (a *CustomAgent6) resetShortTermMemory() {
	a.shortTermMemory = Memory{}
}

func (a *CustomAgent6) updateReassignmentPeriodGuess() {
	a.Log("Team 6 previous floor:", infra.Fields{"prevFloor": a.prevFloor})
	a.Log("Team 6 current floor:", infra.Fields{"currFloor": a.Floor()})
	a.reassignNum++
	a.reassignPeriodGuess = float64(a.Age()) / float64(a.reassignNum)
	a.Log("Team 6 reassignment number:", infra.Fields{"numReassign": a.reassignNum})
	a.Log("Team 6 reassignment period guess:", infra.Fields{"guessReassign": a.reassignPeriodGuess})
}
