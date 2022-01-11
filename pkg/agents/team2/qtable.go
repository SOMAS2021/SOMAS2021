package team2

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	//"time"
)

type maxQAction struct {
	bestAction int
	maxQ       float64
}

//this function is used to get a pair of (best action, max Q-value),
//based on either past or current state.
func (a *CustomAgent2) getMaxQ(state int) maxQAction {
	ret := maxQAction{}
	ret.maxQ = float64(a.qTable[state][0])
	ret.bestAction = 0
	for i := 1; i < len(a.actionSpace); i++ {
		if ret.maxQ < a.qTable[state][i] {
			ret.maxQ = a.qTable[state][i]
			ret.bestAction = i
		}
	}
	return ret
}
func (a *CustomAgent2) updateQTable(state int, action int) {
	//TODO: take hyperparameters out of function

	// Gamma is the discount factor (0≤γ≤1). It determines how much importance we want to give to future
	// rewards. A high value for the discount factor (close to 1) captures the long-term effective award, whereas,
	// a discount factor of 0 makes our agent consider only immediate reward, hence making it greedy.
	Gamma := 0.6

	// Alpha is the learning rate (0<α≤1). Alpha is the extent to which our Q-values are being updated in every iteration.
	Alpha := 0.1
	a.qTable[state][action] = (1-Alpha)*a.qTable[state][action] +
		Alpha*(a.rTable[state][action]+Gamma*a.getMaxQ(a.CheckState()).maxQ)
}

func (a *CustomAgent2) exportQTable() {
	f, err := os.OpenFile(fmt.Sprintf("%s%s%s", a.ID(), "qtable", ".csv"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	qTable := a.qTable
	// make string table
	sQTable := make([][]string, len(qTable))
	for i := 0; i < len(qTable); i++ {
		sQTable[i] = make([]string, len(qTable[0]))
	}

	// convert float64 qtable to string qtable
	for i := 0; i < len(qTable); i++ {
		for j := 0; j < len(qTable[0]); j++ {
			sQTable[i][j] = fmt.Sprint(qTable[i][j])
		}
	}

	err = w.WriteAll(sQTable)

	if err := w.Error(); err != nil {
		log.Fatalln("error writing csv:", err)
	}
}
