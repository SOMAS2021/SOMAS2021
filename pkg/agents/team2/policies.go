package team2

func InitPolicies(numStates int, numActions int) [][]float64 {
	policies := make([][]float64, numStates)
	uniformProb := 1.0 / float64(numActions)
	for i := 0; i < numStates; i++ {
		policies[i] = make([]float64, numActions)
		for j := 0; j < numActions; j++ {
			policies[i][j] = uniformProb
		}
	}
	return policies
}

func (a *CustomAgent2) updatePolicies(state int) {
	Delta := 0.1 / float64(len(a.actionSpace.actionId)-1)
	bestAction := a.getMaxQ(state).bestAction
	sum := 0.0
	for _, action := range a.actionSpace.actionId {
		if action != bestAction {
			a.policies[state][action] -= Delta
			if a.policies[state][action] < 0 {
				a.policies[state][action] = 0
			}
			sum += a.policies[state][action]
		}
	}
	a.policies[state][bestAction] = 1.0 - sum
	a.adjustPolicies()
}

//fix small errors caused during policy update
func (a *CustomAgent2) adjustPolicies() {
	for _, policy := range a.policies {
		sum := 0.0
		for _, prob := range policy {
			sum += prob
		}
		if sum-1.0 != 0.0 {
			policy[0] -= sum - 1.0
		}
	}
}
