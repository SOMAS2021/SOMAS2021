package team6

// Adds food available to agent's long-term and short-term memory
// Only runs on the first tick after platform arrives
func (a *CustomAgent6) addToMemory() {

	if a.PlatformOnFloor() && a.platOnFloorCtr == 0 {

		currPlatFood := a.CurrPlatFood()
		// Infra returns -1 if a.CurrPlatFood() is called and the platform is not at the agent's floor
		if currPlatFood == -1 {
			return
		}
		a.longTermMemory = append(a.longTermMemory, currPlatFood)
		a.shortTermMemory = append(a.shortTermMemory, currPlatFood)

		a.platOnFloorCtr++
	} else if !a.PlatformOnFloor() {
		a.platOnFloorCtr = 0
	}
}

// Returns whether agent finds themselves reassigned
// Could be the case that reassignment has occurred, but agent remains on the same floor, so in their POV they have not been reassigned
func (a *CustomAgent6) isReassigned() bool {
	if a.prevFloor != a.Floor() && a.prevFloor != -1 {
		return true
	}
	return false
}

// Clears short-term memory
// Called if isReassigned() returns true
func (a *CustomAgent6) resetShortTermMemory() {
	a.shortTermMemory = memory{}
}

// Updates agent's reassignment period guess by dividing agent age (in days) and number of times agent has been reassigned
func (a *CustomAgent6) updateReassignmentPeriodGuess() {
	a.numReassigned++
	a.reassignPeriodGuess = float64(a.Age()) / float64(a.numReassigned)
}
