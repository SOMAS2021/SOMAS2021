package team6

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
	"github.com/google/uuid"

	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
)

type memory []food.FoodType

type trust map[uuid.UUID]int

type behaviour float64

type neighbours struct {
	above uuid.UUID
	below uuid.UUID
}

type thresholdBehaviourPair struct {
	threshold behaviour
	bType     string
}

type behaviourParameterWeights struct {
	HPWeight    float64
	floorWeight float64
}

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
	// maximum/minimum trust an agent can have of another
	maxTrust int

	behaviourMap []thresholdBehaviourPair
}

type CustomAgent6 struct {
	*infra.Base
	config team6Config
	//keep track of the lowest floor we've been to
	maxFloorGuess      int
	currBehaviour      behaviour
	foodTakeDay        int
	reqLeaveFoodAmount int
	reqTakeFoodAmount  int
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
	// holding proposed treaty not accepted yet
	proposedTreaties map[uuid.UUID]messages.Treaty
	// Mapping of agent id to trust
	trustTeams trust
	// IDs of agents above and below (based on what we've been told)
	neighbours neighbours
	// Previous day age
	prevAge int
	// Previous day social motive
	prevSM string
	// Change in social motive
	changeSM string
}

// var maxBehaviourThreshold behaviour = 10.0

// Defines the initial/base behaviour of our agents
func (a *CustomAgent6) chooseInitialBehaviour() behaviour {
	// 1. Whole spectrum
	// return behaviour(rand.Float64()) * a.config.maxBehaviourThreshold
	// 2. Only Collectivist
	//return behaviour(rand.Float64()) * maxBehaviourThreshold/4+1
	// 3. Only Collectivist and Selfish, in a ratio 1/1
	//return behaviour(rand.Float64())*(maxBehaviourThreshold-2) + 1
	//return behaviour(rand.Float64())*4 + 1
	// 4. Only Collectivist and Selfish, in a ratio 2/1
	// initialScore := behaviour(rand.Float64()) * maxBehaviourThreshold
	// // Selfish
	// if initialScore >= 8 {
	// 	return behaviour(rand.Float64())*4 + 5
	// } else {
	// 	// Collectivist
	// 	return behaviour(rand.Float64())*4 + 1
	// }
	// 4. Genetic replacement
	return a.chooseBehaviourGenetic()
}

func (a *CustomAgent6) chooseBehaviourGenetic() behaviour {
	smctr := a.SMCtr()
	// fmt.Println(smctr)
	smSoftmax := softmax(smctr)
	// fmt.Println(smSoftmax)
	randn := randnFromDistr(smSoftmax)
	// fmt.Println(randn)
	initBehaviour := a.randnToBehaviour(randn)
	return initBehaviour
}

func New(baseAgent *infra.Base) (infra.Agent, error) {
	initialBehaviour := behaviour(0) // Dummy variable, is changed immediately after assignment
	a := &CustomAgent6{
		Base: baseAgent,
		config: team6Config{
			baseBehaviour:         initialBehaviour,
			stubbornness:          0.2,
			maxBehaviourSwing:     2,
			paramWeights:          behaviourParameterWeights{HPWeight: 0.8, floorWeight: 0.2}, //ensure sum of weights = max behaviour enum
			lambda:                3.0,
			maxBehaviourThreshold: 10,
			prevFoodDiscount:      0.6,
			maxTrust:              25,
			behaviourMap: []thresholdBehaviourPair{
				{1, "Altruist"},
				{5, "Collectivist"},
				{9, "Selfish"},
				{10, "Narcissist"},
			},
		},
		currBehaviour:       initialBehaviour,
		maxFloorGuess:       baseAgent.Floor() + 2,
		foodTakeDay:         0,
		reqLeaveFoodAmount:  -1,
		reqTakeFoodAmount:   -1,
		lastFoodTaken:       0,
		averageFoodIntake:   0.0,
		longTermMemory:      make(memory, 0),
		shortTermMemory:     make(memory, 0),
		numReassigned:       0,
		reassignPeriodGuess: 0,
		platOnFloorCtr:      0,
		prevFloor:           -1,
		proposedTreaties:    make(map[uuid.UUID]messages.Treaty),
		trustTeams:          make(trust),
		neighbours:          neighbours{above: uuid.Nil, below: uuid.Nil},
		prevAge:             0,
		prevSM:              initialBehaviour.string(),
	}

	a.config.baseBehaviour = a.chooseInitialBehaviour()
	a.currBehaviour = a.config.baseBehaviour

	return a, nil
}

// Todo: define some sensible values
func newUtilityParams(socialMotive string) utilityParameters {
	switch socialMotive {
	case "Altruist":
		return utilityParameters{
			g: 0.0,
			r: 1.0,
			c: 5.0,
		}
	case "Collectivist":
		return utilityParameters{
			g: 2.0,
			r: 1.5,
			c: 2.0,
		}
	case "Selfish":
		return utilityParameters{
			g: 5.0,
			r: 2.0,
			c: 1.0,
		}
	case "Narcissist":
		return utilityParameters{
			g: 10.0,
			r: 2.5,
			c: 0.0,
		}
	default:
		// error
		return utilityParameters{}
	}
}

// initialise trust based on our social motive - the more narcissistic we are, the less we're willing to initially trust
func (a *CustomAgent6) startingTrust() int {
	switch a.currBehaviour.string() {
	case "Altruist":
		return 10
	case "Collectivist":
		return 5
	case "Narcissist":
		return -5
	default:
		return 0
	}
}

// handle cases where we haven't yet found out who our neighbours are
func (a *CustomAgent6) updateTrust(amount int, agentID uuid.UUID) {
	if amount < 0 {
		switch a.currBehaviour.string() {
		case "Altruist":
			amount *= 0 // don't care about any negative opinion - totally forgive
		case "Selfish":
			amount *= 2
		case "Narcissist":
			amount *= 4 // baseline trust reduction doubled to penalise more
		default:
		}
	} else {
		switch a.currBehaviour.string() {
		case "Altruist":
			amount *= 4 // double trust increase to show our love
		case "Collectivist":
			amount *= 2
		case "Narcissist":
			amount *= 0 // don't improve trust in people when narcissistic - we hate everyone (Nihilist?)
		default:
		}
	}
	// null mapping has implicit value of 0 (https://staticcheck.io/docs/checks#S1036)
	// a.trustTeams[agentID] += amount

	if _, exists := a.trustTeams[agentID]; exists {
		a.trustTeams[agentID] += amount
	} else {
		a.trustTeams[agentID] = a.startingTrust()
	}

	if a.trustTeams[agentID] > a.config.maxTrust {
		a.trustTeams[agentID] = a.config.maxTrust
	} else if a.trustTeams[agentID] < -a.config.maxTrust {
		a.trustTeams[agentID] = -a.config.maxTrust
	}
}

func (a *CustomAgent6) identifyNeighbours(id uuid.UUID, floor int) {
	floorDir := floor - a.Floor()

	if floorDir == 1 {
		a.neighbours.below = id
	} else if floorDir == -1 {
		a.neighbours.above = id
	}
}

func (b behaviour) string() string {
	behaviourMap := [...]thresholdBehaviourPair{{1, "Altruist"}, {5, "Collectivist"}, {9, "Selfish"}, {10, "Narcissist"}}

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

	// Reporting agent state
	a.Log("Reporting agent state:", infra.Fields{"HP:": a.HP(), "Floor:": a.Floor(), "Social motive:": a.currBehaviour.string()})

	// Everything you need to do once a day
	if a.Age() != a.prevAge {
		a.updateBehaviour()

		if a.currBehaviour.string() == "Collectivist" || a.currBehaviour.string() == "Selfish" {
			treaty := a.constructTreaty()
			a.proposeTreaty(treaty)
		}
		a.requestLeaveFood()
		a.regainTrustInNeighbours()
	}

	// Receiving messages and treaties
	receivedMsg := a.ReceiveMessage()
	if receivedMsg != nil {
		receivedMsg.Visit(a)
	}

	// Updates agent's memory
	if a.isReassigned() {
		a.resetShortTermMemory()
		a.updateReassignmentPeriodGuess()
		// Ask HP to discover the ID of neighbour
		a.SendMessage(messages.NewAskHPMessage(a.ID(), a.Floor(), a.Floor()-1))
		a.SendMessage(messages.NewAskHPMessage(a.ID(), a.Floor(), a.Floor()+1))
	} else if a.numReassigned == 0 { // Before any reassignment, reassignment period guess should be days elapsed
		a.reassignPeriodGuess = float64(a.Age())
		// a.Log("Team 6 reassignment number:", infra.Fields{"numReassign": a.numReassigned})
		// a.Log("Team 6 reassignment period guess:", infra.Fields{"guessReassign": a.reassignPeriodGuess})
	}
	a.addToMemory()

	// Eat if needed/wanted
	desiredFood := a.desiredFoodIntake()
	intendedFood := a.intendedFoodIntake()
	foodTaken, err := a.TakeFood(intendedFood)

	if err == nil {
		a.lastFoodTaken = foodTaken
		// Exponential moving average filter to average food taken whilst discounting previous food
		a.updateAverageIntake(foodTaken)
		// Updates trust in the above neighbour based on average food
		if a.averageFoodIntake < 1 {
			a.updateTrust(-1, a.neighbours.above)
		} else {
			a.updateTrust(1, a.neighbours.above)
		}
		// LOG
		a.Log("Team 6 agent took food:", infra.Fields{"floor": a.Floor(), "hp": a.HP(), "social motive": a.currBehaviour.string(), "desiredFood": desiredFood, "intendedFood": intendedFood})
	}

	// Reset the reqLeaveFoodAmount to nothing once the agent has eaten
	if a.HasEaten() {
		a.reqLeaveFoodAmount = -1
		a.reqTakeFoodAmount = -1
	}

	a.prevFloor = a.Floor() // keep at end of Run() function
	if a.Age() != a.prevAge {
		a.updateChangeSMVariable()
		a.prevSM = a.currBehaviour.string()
	}
	a.prevAge = a.Age()

}

func (a *CustomAgent6) Behaviour() string {
	return a.currBehaviour.string()
}

func (a *CustomAgent6) BehaviourChange() string {
	return a.changeSM
}

func (a *CustomAgent6) updateChangeSMVariable() {
	switch a.prevSM {
	case "Altruist":
		switch a.currBehaviour.string() {
		case "Altruist":
			a.changeSM = "A2A"
		case "Collectivist":
			a.changeSM = "A2C"
		case "Selfish":
			a.changeSM = "A2S"
		case "Narcissist":
			a.changeSM = "A2N"
		default:
			a.changeSM = "N/A"
		}
	case "Collectivist":
		switch a.currBehaviour.string() {
		case "Altruist":
			a.changeSM = "C2A"
		case "Collectivist":
			a.changeSM = "C2C"
		case "Selfish":
			a.changeSM = "C2S"
		case "Narcissist":
			a.changeSM = "C2N"
		default:
			a.changeSM = "N/A"
		}
	case "Selfish":
		switch a.currBehaviour.string() {
		case "Altruist":
			a.changeSM = "S2A"
		case "Collectivist":
			a.changeSM = "S2C"
		case "Selfish":
			a.changeSM = "S2S"
		case "Narcissist":
			a.changeSM = "S2N"
		default:
			a.changeSM = "N/A"
		}
	case "Narcissist":
		switch a.currBehaviour.string() {
		case "Altruist":
			a.changeSM = "N2A"
		case "Collectivist":
			a.changeSM = "N2C"
		case "Selfish":
			a.changeSM = "N2S"
		case "Narcissist":
			a.changeSM = "N2N"
		default:
			a.changeSM = "N/A"
		}
	default:
		a.changeSM = "N/A"
	}
	// a.Log("SM Change summary", infra.Fields{"prevSM": a.prevSM, "currSM": a.currBehaviour.string(), "functionOutput": a.changeSM})
}

func (a *CustomAgent6) randnToBehaviour(randn int) behaviour {
	// behaviourMap := [...]thresholdBehaviourPair{{1, "Altruist"}, {5, "Collectivist"}, {9, "Selfish"}, {10, "Narcissist"}}
	switch randn {
	case 0: // Altruist
		return a.config.behaviourMap[0].threshold / 2
	case 1: // Collectivist
		return (a.config.behaviourMap[0].threshold + a.config.behaviourMap[1].threshold) / 2
	case 2: // Selfish
		return (a.config.behaviourMap[1].threshold + a.config.behaviourMap[2].threshold) / 2
	case 3: // Narcissist
		return (a.config.behaviourMap[2].threshold + a.config.behaviourMap[3].threshold) / 2
	default: // Random SM
		return behaviour(rand.Float64()) * a.config.maxBehaviourThreshold
	}
}

// Returns int between 0 and number of bins - 1 according to given distribution
func randnFromDistr(distr []float64) int {

	partitions := []float64{distr[0]}
	for i := 1; i < 4; i++ {
		partitions = append(partitions, distr[i]+partitions[i-1])
	}

	randFloat := rand.Float64()
	// fmt.Println(randFloat)
	for i, val := range partitions {
		if randFloat <= val {
			return i
		}
	}
	return -1
}

func softmax(inputArr []int) []float64 {
	sumExp := 0.0
	for _, val := range inputArr {
		sumExp += math.Exp(float64(val))
	}

	softmaxOutput := []float64{}

	for _, val := range inputArr {
		softmaxOutput = append(softmaxOutput, math.Exp(float64(val))/sumExp)
	}

	return softmaxOutput
}
