package team2

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

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

func (a *CustomAgent2) updatePolicies(state int, action int) {
	//Update state visit counter
	a.visitCounter[state]++
	//update average policies
	for i := 0; i < len(a.actionSpace); i++ {
		a.averagePolicies[state][i] = a.averagePolicies[state][i] + ((1.0 / float64(a.visitCounter[state])) * (a.policies[state][i] - a.averagePolicies[state][i]))
	}
	//update policies by WoLF
	Delta := a.winOrLose(state) / float64(len(a.actionSpace)-1)
	bestAction := a.getMaxQ(state).bestAction
	sum := 0.0
	for i := 0; i < len(a.actionSpace); i++ {
		if i != bestAction {
			a.policies[state][i] -= Delta
			if a.policies[state][i] < 0 {
				a.policies[state][i] = 0
			}
			sum += a.policies[state][i]
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

func (a *CustomAgent2) winOrLose(state int) float64 {
	sumPolicies := 0.0
	sumAveragePolicies := 0.0
	Delta_win := 0.025
	Delta_lose := 0.1
	for i := 0; i < len(a.actionSpace); i++ {
		sumPolicies += a.policies[state][i] * a.qTable[state][i]
		sumAveragePolicies += a.averagePolicies[state][i] * a.qTable[state][i]
	}
	if sumPolicies > sumAveragePolicies {
		return Delta_win
	}
	return Delta_lose
}

func (a *CustomAgent2) exportPolicies() {
	f, err := os.OpenFile(fmt.Sprintf("%s%s%s", a.ID(), "policies", ".csv"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Println("error opening csv: ", err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	policies := a.policies
	// make string table
	sPolicies := make([][]string, len(policies))
	for i := 0; i < len(policies); i++ {
		sPolicies[i] = make([]string, len(policies[0]))
	}

	// convert float64 policies to string policies
	for i := 0; i < len(policies); i++ {
		for j := 0; j < len(policies[0]); j++ {
			sPolicies[i][j] = fmt.Sprint(policies[i][j])
		}
	}

	err = w.WriteAll(sPolicies)

	if err != nil {
		log.Println("error writing csv:", err)
	}
}
