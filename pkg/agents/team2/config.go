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
        hp			              0.0                     100.0
        floor	                  1                       inf
        foodOnPlat                0    	                  100

		Other observations should come from communication with other agents
        savedAgents               0                       number of agent per florr
		Note: the particular combination of the observations correspond to a particular state of
		the agent
	State:
		We initially only define 3x3x3 (27) states for testing purposes

		Num				hp			floor			foodOnPlat
		0				61-100		1-30				61-100
		1				61-100		1-30				31-60
		2				61-100		1-30				0-30
		3				61-100		31-60				61-100
		...
		26				0-30		>60					0-30
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

type actionSpace struct {
	//TODO: actionId is redundant and may be removed in further version
	actionId  []int
	actionSet map[int]func(hp int) food.FoodType
}
type CustomAgent2 struct {
	*infra.Base
	stateSpace            [][][][]int
	actionSpace           actionSpace
	policies              [][]float64
	rTable                [][]float64
	qTable                [][]float64
	isPlatArrived         bool
	daysAtCriticalCounter int
	PreviousdayAtCritical int
}

func InitTable(numStates int, numActions int) [][]float64 {
	table := make([][]float64, numStates)
	for i := 0; i < numStates; i++ {
		table[i] = make([]float64, numActions)
	}
	return table
}

func New(baseAgent *infra.Base) (infra.Agent, error) {
	hpStatesDim := baseAgent.HealthInfo().MaxHP + 1
	actionDim := baseAgent.HealthInfo().MaxHP + 1
	daysAtCriticalDim := baseAgent.HealthInfo().MaxDayCritical + 1

	stateSpace := InitStateSpace(hpStatesDim, 3, 3, daysAtCriticalDim)
	actionSpace := InitActionSpace(actionDim)
	policies := InitPolicies(hpStatesDim*3*3*daysAtCriticalDim, actionDim)
	rTable := InitTable(hpStatesDim*3*3*daysAtCriticalDim, actionDim)
	qTable := InitTable(hpStatesDim*3*3*daysAtCriticalDim, actionDim)

	return &CustomAgent2{
		Base:                  baseAgent,
		stateSpace:            stateSpace,
		actionSpace:           actionSpace,
		policies:              policies,
		rTable:                rTable,
		qTable:                qTable,
		isPlatArrived:         false,
		daysAtCriticalCounter: 0,
		PreviousdayAtCritical: 0,
	}, nil
}

func (a *CustomAgent2) Run() {
	//Communication & Observation
	//communicate before platform arrives to my floor

	/*
		msg := *messages.NewBaseMessage(a.Floor())
		a.SendMessage(1, msg)
	*/

	//Perform the following only when platform arrives
	//NOTE: should let infra team add a func to see whether the plaftfrom has arrived or not
	if a.CheckState() != -1 {
		a.isPlatArrived = true
	} else {
		a.isPlatArrived = false
	}
	if a.isPlatArrived {
		oldState := a.CheckState()
		oldHP := a.HP()
		a.Log("Agent team2 before action:", infra.Fields{"floor": a.Floor(), "hp": oldHP, "food": a.CurrPlatFood(), "state": oldState})
		action := a.SelectAction()
		_, err := a.TakeFood(food.FoodType(a.actionSpace.actionId[action])) //perform selected action
		if err != nil {
			switch err.(type) {
			case *infra.FloorError:
			case *infra.NegFoodError:
			case *infra.AlreadyEatenError:
			default:
			}
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
		hpInc := a.HP() - oldHP
		a.updateRTable(hpInc, oldState, action)
		a.updateQTable(oldState, action)
		a.updatePolicies(oldState)
	}

}
