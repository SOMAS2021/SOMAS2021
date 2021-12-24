package team2

func InitPolicies(numStates int, numActions int) [][]float32 {
	var policies = make([][]float32, numStates)
	uniformProb := 1.0 / float32(numActions)
	for i := 0; i < numStates; i++ {
		policies[i] = make([]float32, numActions)
		for j := 0; j < 3; j++ {
			policies[i][j] = uniformProb
		}
	}
	return policies
}

func (a *CustomAgent2) updataPolicies(state int) {
	Delta := float32(0.1) / float32(len(a.actionSpace.actionId)-1)
	bestAction := a.getMaxQ(state).bestAction
	sum := float32(0.0)
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
		sum := float32(0.0)
		for _, prob := range policy {
			sum += prob
		}
		if sum-1.0 != 0.0 {
			policy[0] -= sum - 1.0
		}
	}
}
