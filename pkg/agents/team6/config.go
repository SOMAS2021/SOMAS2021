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
// z - desired food (maximum of the function)
func Utility(x, z float64, socialMotive string) float64 {

	params := NewUtilityParams(socialMotive)
	// calculate the function scaling parameter a

	if socialMotive == "Altruist" /*|| socialMotive == "Nacissist"*/ {
		// Don't scale depending on desired food
		return params.g*math.Pow(x, 1/params.r) - params.c*x
	} else {
		a := (1 / z) * math.Pow((params.c*params.r)/params.g, params.r/(1-params.r))
		return params.g*math.Pow(a*x, 1/params.r) - params.c*a*x
	}
}

func min(x, y food.FoodType) food.FoodType {
	if x < y {
		return x
	}
	return y
}

// Evaluate our agents current utility based on the current desired food
func (a *CustomAgent6) evaluateUtility(mem memory) float64 {
	sum := food.FoodType(0)
	for _, foodAvailable := range mem {
		sum += food.FoodType(Utility(float64(min(foodAvailable, a.desiredFoodIntake())), float64(a.desiredFoodIntake()), a.currBehaviour.String()))
	}

	return float64(sum) / math.Max(float64(len(mem)), 1.0)
}

// Decides if to accept or reject a treaty
func (a *CustomAgent6) considerTreaty(t *messages.Treaty) bool {

	// 1. Estimate the food intake of the proposed treaty

	// Calculate how much food has been available on average on this floor
	sum := food.FoodType(0)
	for _, food := range a.shortTermMemory {
		sum += food
	}
	averageFoodAvailable := float64(sum) / math.Max(float64(len(a.shortTermMemory)), 1.0)

	// convert different treaty "Request types" to a "food intake"
	estimatedTakeFood := 0.0
	switch t.Request() {
	case messages.LeaveAmountFood:
		estimatedTakeFood = averageFoodAvailable - float64(t.RequestValue())
	case messages.LeavePercentFood:
		estimatedTakeFood = averageFoodAvailable * (1.0 - float64(t.RequestValue()))
	case messages.Inform:
		// TODO: Handel this form of treaty
		return false
	default:
		return false
	}

	// check the exact request condition
	if t.ConditionOp() == messages.GE || t.ConditionOp() == messages.GT /*t.ConditionOp() == messages.EQ || */ {
		// The treaty is of the form "Take X (or less) food"

		// 2. Calculate the agent's utility given different outcomes (accept or reject treaty)
		
		// a. estimated utility when accepting the treaty
		treatyTrustFactor := 0.8 // TODO: Include population value
		treatyUtility := treatyTrustFactor * Utility(estimatedTakeFood, float64(a.desiredFoodIntake()), a.currBehaviour.String())
		
		// b. estimated utility when rejecting the treaty
		currentShortTermUtility := a.evaluateUtility(a.shortTermMemory) // estimated utility on the current floor
		currentLongTermUtility := a.evaluateUtility(a.longTermMemory) // estimated utility over the entire time in the tower
		
		// benefit of signing the treaty (> 0 == beneficial)
		shortTermBenefit := treatyUtility - currentShortTermUtility
		longTermBenefit := treatyUtility - currentLongTermUtility

		estmatedPeriod := len(a.longTermMemory)
		if a.numReassigned != 0 {
			estmatedPeriod /= a.numReassigned
		}
		estimatedTimeLeft := estmatedPeriod - len(a.shortTermMemory)
		benefit := 0.0

		// only think in the short term if
		// - the duration is shorter than the time left in estimated reassignment period OR
		// - we only have short term experience OR
		// - HP is critical (survival instincts take over)
		if t.Duration() < estimatedTimeLeft ||
			len(a.shortTermMemory) == len(a.longTermMemory) ||
			a.HP() <= a.HealthInfo().HPCritical {
			benefit = shortTermBenefit

		} else {
			// → The longer the duration of a treaty, the more important is it's long term benefit
			// E.g.
			// (1) time left on this level = 4, duration = 5 --> shortTermFocus 80%
			// 		--> The treaty mostly matters in the short term
			// (2) time left on this level = 4, duration = 8 --> shortTermFocus 50%
			// 		--> The treaty counts as much on this level as for the time after
			// (3) time left on this level = 1, duration = 100 --> shortTermFocus 1%
			// 		--> The treaty mostly matters in the long term
			
			shortTermFocus := float64(estimatedTimeLeft) / float64(t.Duration())
			benefit = shortTermFocus*shortTermBenefit + (1.0-shortTermFocus)*longTermBenefit
		}

		// 3. If we benefit from the treaty, accept it
		return benefit > 0.0

	} else {
		// The treaty is of the form "Take X (or mire) food"

		// - collectivist --> accept treaty if estimatedTakeFood less than (hoping to get others to eat at least the critical level)
		if a.currBehaviour.String() == "Collectivist" {
			return estimatedTakeFood <= float64(a.desiredFoodIntake()) // TODO: desired food will be 0 some of the days (instead 2??)
		}
		// All other social motives will allways reject
		// - altruist wants to avoid eating anyhing to save more for others
		// - selfish / narcissist doesn't want others to eat more than they want
		return false
	}
}
