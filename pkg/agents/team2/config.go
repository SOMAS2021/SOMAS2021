package team2

import (
	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
)

/*
    Description:
		An agent observes the status of itself and the other agents. Based on the observations,
		the agents decides how much food to take in order to increase individual and collective
		utilities.
    Observation:
        Observation               Min                     Max
        TODO

		Other observations should come from communication with other agents
        savedAgents               0                       number of agent per florr
		Note: the particular combination of the observations correspond to a particular state of
		the agent
	State:
		TODO
    Action:
		We separate the TakeFood action into several actions

        Action					Comment
        disregard food 			takeFood(0), sparing food for agents below me
		satisfice with food 	takeFood(min), eat only food required for survival
		satisfy with food   	takeFood(max), eat as manu as possible (get max hp increment)
	Policy:
		We associate a probability, or a policy, to a certain action in a certain state.
    Reward:
        For MVP, the reward is only related with to hp increment of the agent itself.
	Q-table:
		We associate a Q-value to each action under a particular state.
    Starting State:
        Initialise policies to a uniform distribution
		Initialise Q-table to 0

    Episode Termination:
        Agent gets killed
		or
		Simulation terminated
*/

type CustomAgent2 struct {
	*infra.Base
	stateSpace            [][][]int
	actionSpace           []food.FoodType
	policies              [][]float64
	rTable                [][]float64
	qTable                [][]float64
	daysAtCriticalCounter int
	PreviousdayAtCritical int
	lastAge               int
	neiboughHP            int
	lastAction            int
	lastFoodTaken         food.FoodType
	cumulativeRewards     float64
}

func InitTable(numStates int, numActions int) [][]float64 {
	table := make([][]float64, numStates)
	for i := 0; i < numStates; i++ {
		table[i] = make([]float64, numActions)
	}
	return table
}

func New(baseAgent *infra.Base) (infra.Agent, error) {
	actionDim := 6
	daysAtCriticalDim := baseAgent.HealthInfo().MaxDayCritical + 1

	stateSpace := InitStateSpace(10, 10, daysAtCriticalDim)
	actionSpace := InitActionSpace(actionDim)
	policies := InitPolicies(10*10*daysAtCriticalDim, actionDim)
	rTable := InitTable(10*10*daysAtCriticalDim, actionDim)
	qTable := InitTable(10*10*daysAtCriticalDim, actionDim)

	return &CustomAgent2{
		Base:                  baseAgent,
		stateSpace:            stateSpace,
		actionSpace:           actionSpace,
		policies:              policies,
		rTable:                rTable,
		qTable:                qTable,
		daysAtCriticalCounter: 0,
		PreviousdayAtCritical: 0,
		lastAge:               0,
		neiboughHP:            -1,
		lastAction:            0,
		lastFoodTaken:         0,
		cumulativeRewards:     0,
	}, nil
}

func (a *CustomAgent2) Run() {
	//Communication & Observation
	//communicate before platform arrives to my floor

	/*
		msg := *messages.NewBaseMessage(a.Floor())
		a.SendMessage(1, msg)
	*/

	//Perform the following only once per day when platform arrives
	if a.PlatformOnFloor() && a.isNewDay() {
		oldState := a.CheckState()
		oldHP := a.HP()
		a.Log("Agent team2 before action:", infra.Fields{"floor": a.Floor(), "hp": oldHP, "food": a.CurrPlatFood(), "state": oldState})
		action := a.SelectAction()
		//fmt.Printf("************************\n")
		//fmt.Printf("It's day %d!\n", a.Age())
		//fmt.Printf("Food on Platform: %d\n", a.CurrPlatFood())
		//fmt.Printf("HP before action: %d\n", oldHP)

		foodTaken, err := a.TakeFood(food.FoodType(a.actionSpace[action])) //perform selected action
		if err != nil {
			//if there's error, cease updating tables
			return
		}
		a.Log("Agent team2:", infra.Fields{"selected and performed action": action})
		a.Log("Agent team2 after action:", infra.Fields{"floor": a.Floor(), "hp": a.HP(), "food": a.CurrPlatFood(), "state": a.CheckState()})
		if a.DaysAtCritical() > 0 && a.PreviousdayAtCritical != a.DaysAtCritical() {
			a.PreviousdayAtCritical = a.DaysAtCritical()
			a.daysAtCriticalCounter += 1
			if a.DaysAtCritical() >= (a.HealthInfo().MaxDayCritical - 1) {
				a.Log("Agent team2 at critical state", infra.Fields{"daysAtCriticalCounter": a.daysAtCriticalCounter, "floor": a.Floor(), "hp": a.HP(), "food": a.CurrPlatFood(), "state": a.CheckState()})
			}
		}
		//fmt.Printf("Intended action: %d\n", action*5)
		//fmt.Printf("Actual eaten food: %d\n", foodTaken)
		//fmt.Printf("HP after action: %d\n", a.HP())
		//fmt.Printf("************************\n")
		hpInc := a.HP() - oldHP
		a.updateRTable(oldHP, hpInc, int(foodTaken), oldState, action)
		a.updateQTable(oldState, action)
		a.updatePolicies(oldState)
		a.lastAge = a.Age()

		//write hp vs days
		target := []float64{float64(a.HP())}
		a.writeToCSV(target, "csv_files/HP_by_days", 1)
		//write foodtaken vs hp
		target = []float64{float64(foodTaken), float64(a.HP())}
		a.writeToCSV(target, "csv_files/foodtaken_vs_hp", 2)
		// table logging
		a.exportQTable()
		a.exportPolicies()
	}

}

func (a *CustomAgent2) isNewDay() bool {
	if a.Age() < a.lastAge {
		a.lastAge = -1
		return true
	}
	if a.Age() == a.lastAge {
		return false
	}
	return true
}

func (a *CustomAgent2) CustomLogs() {
	a.exportQTable()
	a.exportPolicies()
}
