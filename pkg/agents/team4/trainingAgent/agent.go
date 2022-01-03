package team4TrainingEvoAgent

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"

	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
)

type CustomAgentEvoParams struct {
	locked            bool    // used to lock the wait time before eating
	foodToEat         []int   // the amount of food to eat for various health levels
	daysToWait        []int   // the days to wait before eating for various health levels
	ageLastEaten      int     // the age at which the agent last ate
	morality          float64 // the morality of the agent that determines how selfishly or selflessly the agent will act
	traumaScaleFactor float64 // the amount of trauma the agent has suffered which effects the amount of food it is likely to eat

	globalTrust      float32            // the overall trust the agent has in other agents in the tower
	coefficients     []float32          // the amount trust score changes by for certain actions
	lastFoodTaken    food.FoodType      // food taken on the previous day
	sentMessages     []messages.Message // TODO: make it a map hashed by messageIDs
	responseMessages []messages.Message // TODO: make it a map hashed by messageIDs
	lastPlatFood     food.FoodType      // last seen food on the platform
	maxFoodLimit     food.FoodType      // maximum food we want to allow others to eat
	messageCounter   int                // the total number of messages we send in a day
	globalTrustLimit float32            // limit to check whether to be selfish or not
	lastAge          int                // the age of the agent on the previous day
	healthStatus     int
}

type CustomAgentEvo struct {
	*infra.Base
	// new params
	params CustomAgentEvoParams
}

// Data struct to hold evo param values
type LoadedData struct {
	FoodToEat  []int
	DaysToWait []int
}

func InitaliseParams(baseAgent *infra.Base) CustomAgentEvoParams {
	mydir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	file, _ := ioutil.ReadFile(fmt.Sprintf("%s/pkg/agents/team4/trainingAgent/agentConfig.json", mydir))
	var data1 LoadedData
	_ = json.Unmarshal(file, &data1) //parse the config json file to get the coeffs for floor and HP equations

	data1.FoodToEat[0] = baseAgent.HealthInfo().HPReqCToW
	data1.DaysToWait[0] = int(baseAgent.HealthInfo().MaxDayCritical / 2)

	return CustomAgentEvoParams{ //initialise the parameters of the agent
		locked:            true,
		foodToEat:         data1.FoodToEat,
		daysToWait:        data1.DaysToWait,
		ageLastEaten:      0,
		morality:          100 * rand.Float64(), // TODO: Use this properly
		traumaScaleFactor: 1,
		healthStatus:      3,
		globalTrust:       0.0,
		coefficients:      []float32{2, 4, 8}, // TODO: maybe train these co-efficients using evolutionary algorithm
		lastFoodTaken:     0,
		sentMessages:      []messages.Message{},
		responseMessages:  []messages.Message{},
		lastPlatFood:      -1,
		maxFoodLimit:      50,
		messageCounter:    0,
		globalTrustLimit:  75,
		lastAge:           0,
	}
}

func New(baseAgent *infra.Base) (infra.Agent, error) {
	//create other parameters
	return &CustomAgentEvo{
		Base:   baseAgent,
		params: InitaliseParams(baseAgent),
	}, nil
}

// Checks if neighbour below has eaten
func (a *CustomAgentEvo) NeighbourFoodEaten() food.FoodType {
	if a.CurrPlatFood() != -1 {
		if !a.PlatformOnFloor() && a.CurrPlatFood() != a.params.lastPlatFood {
			return a.params.lastPlatFood - a.CurrPlatFood()
		}
		return 0
	}
	return -1
}

// removes a specific message from a message array
func remove(slice []messages.Message, s int) []messages.Message {
	return append(slice[:s], slice[s+1:]...)
}

// Check if a day has passed
func (a *CustomAgentEvo) HasDayPassed() bool {
	if a.Age() != a.params.lastAge {
		a.params.lastAge = a.Age()
		return true
	}
	return false
}

func (a *CustomAgentEvo) Run() {

	if food.FoodType(a.CurrPlatFood()) != a.params.lastPlatFood && a.PlatformOnFloor() {
		a.params.lastPlatFood = a.CurrPlatFood()
	}

	receivedMsg := a.ReceiveMessage()
	if receivedMsg != nil {
		receivedMsg.Visit(a)
	} else {
		a.Log("I got nothing")
	}

	//TODO: Define a threshold limit for other agents to respond to our sent message.
	a.SendingMessage()

	// a.intendedFoodTaken = food.FoodType(int(int(a.CurrPlatFood()) * (100 - int(a.globalTrust)) / 100))
	// a.lastFoodTaken, _ = a.TakeFood(a.intendedFoodTaken)

	//dayPass := a.HasDayPassed()

	healthLevelSeparation := int(0.33 * float64(a.HealthInfo().MaxHP-a.HealthInfo().WeakLevel))

	if a.HP() <= a.HealthInfo().WeakLevel { //critical
		a.params.healthStatus = 0
		a.params.locked = false
		// if dayPass {
		// 	a.params.traumaScaleFactor = math.Min(200, a.params.traumaScaleFactor+0.03)
		// }
	}
	if !a.params.locked {
		if a.HP() <= a.HealthInfo().WeakLevel+healthLevelSeparation { //weak
			a.params.healthStatus = 1
			// if dayPass {
			// 	a.params.traumaScaleFactor = math.Min(200, a.params.traumaScaleFactor+0.02)
			// }
		} else if a.HP() <= a.HealthInfo().WeakLevel+2*healthLevelSeparation { //normal
			a.params.healthStatus = 2
			// if dayPass {
			// 	a.params.traumaScaleFactor = math.Max(0, a.params.traumaScaleFactor-0.02)
			// }
		} else { //strong
			a.params.healthStatus = 3
			// if dayPass {
			// 	a.params.traumaScaleFactor = math.Max(0, a.params.traumaScaleFactor-0.03)
			// }
		}
		a.params.locked = true
	}

	var foodEaten food.FoodType
	var err error

	if (a.Age()-a.params.ageLastEaten) >= a.params.daysToWait[a.params.healthStatus] || a.params.healthStatus == 0 {
		a.params.locked = false
		calculatedAmountToEat := a.params.traumaScaleFactor * float64(a.params.foodToEat[a.params.healthStatus])
		foodEaten, err = a.TakeFood(food.FoodType(calculatedAmountToEat))
		a.params.ageLastEaten = a.Age()
		if err != nil {
			switch err.(type) {
			case *infra.FloorError:
			case *infra.NegFoodError:
			case *infra.AlreadyEatenError:
			default:
			}
		}
	}

	a.Log("team4EvoAgent reporting status:", infra.Fields{"floor": a.Floor(), "hp": a.HP(), "FoodToEat": a.params.foodToEat[a.params.healthStatus], "DaysToWait": a.params.daysToWait, "foodEaten": foodEaten})
}
