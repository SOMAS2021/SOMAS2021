package team6

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/messages"

	//"github.com/SOMAS2021/SOMAS2021/pkg/messages"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
)

type memory []food.FoodType

type behaviour float64

// const (
//  altruist behaviour = iota
//  collectivist
//  selfish
//  narcissist
// )

type utilityParameters struct {
	// Greediness
	g float64
	// Risk aversion
	r float64
	// community cost
	c float64
}

type team6Config struct {
	baseBehaviour behaviour
	//the scaling factor which limits the change in agent behaviour
	stubbornness float64
	//the largest jump in behaviour an agent can take
	maxBehaviourSwing float64
	//weights used to assess score for behaviour update
	paramWeights behaviourParameterWeights
	//floor scaling discount factor
	lambda float64
	//maximum behaviour score an agent can reach
	maxBehaviourThreshold behaviour
	//discount previous food intakes for EMA filter
	prevFoodDiscount float64
}

type CustomAgent6 struct {
	*infra.Base
	config team6Config
	//keep track of the lowest floor we've been to
	maxFloorGuess      int
	currBehaviour      behaviour
	foodTakeDay        int
	reqLeaveFoodAmount int
	lastFoodTaken      food.FoodType
	averageFoodIntake  float64
	// Memory of food available throughout agent's lifetime
	longTermMemory memory
	// Memory of food available while agent is at a particular floor
	shortTermMemory memory
	// Number of times the agent has been reassigned
	numReassigned int
	// What the agent thinks the reassignment period is
	reassignPeriodGuess float64
	// Counts how many ticks the platform is at the agent's floor for. Used to call functions only once when the platform arrives
	platOnFloorCtr int
	// Keeps track of previous floor to see if agent has been reassigned
	prevFloor int
}

type thresholdBehaviourPair struct {
	threshold behaviour
	bType     string
}

type behaviourParameterWeights struct {
	HPWeight    float64
	floorWeight float64
}

var maxBehaviourThreshold behaviour = 10.0

func chooseInitialBehaviour() behaviour {
	return behaviour(rand.Float64()) * maxBehaviourThreshold
}

func New(baseAgent *infra.Base) (infra.Agent, error) {
	initialBehaviour := chooseInitialBehaviour()
	return &CustomAgent6{
		Base: baseAgent,
		config: team6Config{
			baseBehaviour:         initialBehaviour,
			stubbornness:          0.0,
			maxBehaviourSwing:     8,
			paramWeights:          behaviourParameterWeights{HPWeight: 0.7, floorWeight: 0.3}, //ensure sum of weights = max behaviour enum
			lambda:                3.0,
			maxBehaviourThreshold: maxBehaviourThreshold,
			prevFoodDiscount:      0.6,
		},
		currBehaviour:       initialBehaviour,
		maxFloorGuess:       baseAgent.Floor() + 2,
		foodTakeDay:         0,
		reqLeaveFoodAmount:  -1,
		lastFoodTaken:       0,
		averageFoodIntake:   0.0,
		longTermMemory:      memory{},
		shortTermMemory:     memory{},
		numReassigned:       0,
		reassignPeriodGuess: 0,
		platOnFloorCtr:      0,
		prevFloor:           -1,
	}, nil
}

// Todo: define some sensible values
func NewUtilityParams(socialMotive string) utilityParameters {
	switch socialMotive {
	case "Altruist":
		return utilityParameters{
			g: 1.0,
			r: 2.0,
			c: 3.0,
		}
	case "Collectivist":
		return utilityParameters{
			g: 1.0,
			r: 2.0,
			c: 3.0,
		}
	case "Selfish":
		return utilityParameters{
			g: 1.0,
			r: 2.0,
			c: 3.0,
		}
	case "Narcissist":
		return utilityParameters{
			g: 1.0,
			r: 2.0,
			c: 3.0,
		}
	default:
		// error
		return utilityParameters{}
	}
}

func (b behaviour) String() string {
	behaviourMap := [...]thresholdBehaviourPair{{2, "Altruist"}, {7, "Collectivist"}, {9, "Selfish"}, {10, "Narcissist"}}

	if b >= 0 {
		for _, v := range behaviourMap {
			if b <= v.threshold {
				return v.bType
			}
		}
	}

	return fmt.Sprintf("UNKNOWN Behaviour '%v'", int(b))
}

func (a *CustomAgent6) Run() {

	// a.Log("Custom agent 6 before update:", infra.Fields{"floor": a.Floor(), "hp": a.HP(), "behaviour": a.currBehaviour.String(), "maxFloorGuess": a.maxFloorGuess})

	a.updateBehaviour()

	// Sending messages
	a.RequestLeaveFood()

	// Receiving messages
	// receivedMsg := a.ReceiveMessage()
	// if receivedMsg != nil {
	// 	receivedMsg.Visit(a)
	// } else {
	// 	a.Log("I got no thing")
	// }

	// MEMORY STUFF
	if a.isReassigned() {
		a.resetShortTermMemory()
		a.updateReassignmentPeriodGuess()
	} else if a.numReassigned == 0 { // Before any reassignment, reassignment period guess should be days elapsed
		a.reassignPeriodGuess = float64(a.Age())
		a.Log("Team 6 reassignment number:", infra.Fields{"numReassign": a.numReassigned})
		a.Log("Team 6 reassignment period guess:", infra.Fields{"guessReassign": a.reassignPeriodGuess})
	}
	a.addToMemory()

	foodTaken, err := a.TakeFood(a.intendedFoodIntake())
	if err != nil {
		switch err.(type) {
		case *infra.FloorError:
		case *infra.NegFoodError:
		case *infra.AlreadyEatenError:
		default:
		}
	} else {
		a.lastFoodTaken = foodTaken
	}

	//exponential moving average filter to average food taken whilst discounting previous food
	a.updateAverageIntake(foodTaken)

	// a.Log("Team 6 took:", infra.Fields{"foodTaken": foodTaken, "bType": a.currBehaviour.String()})
	// a.Log("Team 6 agent has HP:", infra.Fields{"hp": a.HP()})

	a.updateBehaviourWeights()

	//fmt.Println(a.ActiveTreaties())

	treaty := messages.NewTreaty(1, 1, 1, 1, 1, 1, 5, a.ID())
	min, max := a.foodRange()
	valid := a.treatyValid(*treaty)

	a.Log("Team 6 processed treaty:", infra.Fields{"treaty": treaty, "range": max - min, "isValid:": valid})
	// // treatyMsg := messages.NewProposalMessage(a.ID(), a.Floor()+1, *treaty)

	// treatyMsg.Visit(a).

	a.prevFloor = a.Floor() // keep at end of Run() function

}

// The utility function with
// x - food input
// g - greediness
// r - risk aversion
// c - community cost
// z - desired food (maximum of the function)
func Utility(x, z float64, params utilityParameters) float64 {
	// calculate the function scaling parameter a
	a := (1 / z) * math.Pow((params.c*params.r)/params.g, params.r/(1-params.r))
	return params.g*math.Pow(a*x, 1/params.r) - params.c*a*x
}
