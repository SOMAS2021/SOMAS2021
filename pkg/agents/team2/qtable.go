package team2

import (
	"math"
)

type maxQAction struct {
	bestAction int
	maxQ       float64
}

//this function is used to get a pair of (best action, max Q-value),
//based on either past or current state.
func (a *CustomAgent2) getMaxQ(state int) maxQAction {
	ret := maxQAction{}
	ret.maxQ = 0.0
	ret.bestAction = 0
	for _, action := range a.actionSpace.actionId {
		ret.maxQ = math.Max(ret.maxQ, a.rTable[state][action])
		if ret.maxQ-a.rTable[state][action] == 0.0 {
			ret.bestAction = action
		}
	}
	return ret
}
func (a *CustomAgent2) updateQTable(state int, action int) {
	//TODO: take hyperparameters out of function

	// Gamma is the discount factor (0≤γ≤1). It determines how much importance we want to give to future
	// rewards. A high value for the discount factor (close to 1) captures the long-term effective award, whereas,
	// a discount factor of 0 makes our agent consider only immediate reward, hence making it greedy.
	Gamma := float64(0.6)

	// Alpha is the learning rate (0<α≤1). Alpha is the extent to which our Q-values are being updated in every iteration.
	Alpha := float64(0.1)
	a.qTable[state][action] = (1-Alpha)*a.qTable[state][action] +
		Alpha*(a.rTable[state][action]+Gamma*a.getMaxQ(a.CheckState()).maxQ)
}
