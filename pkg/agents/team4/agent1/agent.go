package team4EvoAgent

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

// type Coefficient struct {
// 	Floor []float64
// 	Hp    []float64
// }

// type Equation struct {
// 	coefficients []float64
// }

// // function to evaluate a function f(x) at a specific value of x0
// func (equation *Equation) EvaluateEquation(input int) float64 {
// 	sum := 0.0
// 	for i, coeff := range equation.coefficients {
// 		power := math.Pow(float64(input), float64(i))
// 		sum += coeff * power
// 	}
// 	return sum
// }

// // function to pick the correct co-efficients from the configuration file
// func GenerateEquations() (Equation, Equation) {
// 	mydir, err := os.Getwd()
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	file, _ := ioutil.ReadFile(fmt.Sprintf("%s/pkg/agents/team4/agent1/agentConfig.json", mydir))
// 	var data1 Coefficient
// 	_ = json.Unmarshal(file, &data1) //parse the config json file to get the coeffs for floor and HP equations

// 	var newAgentFloor Equation
// 	var newAgentHp Equation

// 	newAgentFloor.coefficients = data1.Floor
// 	newAgentHp.coefficients = data1.Hp

// 	return newAgentFloor, newAgentHp
// }

type CustomAgentEvoParams struct {
	// these params are based only on singular inputs
	// scalingEquation   Equation
	// currentFloorScore Equation
	// currentHpScore    Equation
	locked       bool
	previousAge  int
	foodToEat    []int
	daysToWait   []int
	ageLastEaten int
	// below params updated based in previous experience of floors
	trustscore        float64
	morality          float64
	traumaScaleFactor float64
}

type CustomAgentEvo struct {
	*infra.Base
	// new params
	params              CustomAgentEvoParams
	globalTrust         float32
	globalTrustAdd      float32
	globalTrustSubtract float32
	coefficients        []float32
	lastFoodTaken       food.FoodType
	intendedFoodTaken   food.FoodType
	sentMessages        []messages.Message //TODO: make it a map hashed by messageIDs
	responseMessages    []messages.Message //TODO: make it a map hashed by messageIDs
	MessageToSend       int
	lastPlatFood        food.FoodType
	maxFoodLimit        food.FoodType
	messageCounter      int
	globalTrustLimit    float32
	lastAge             int
}

type LoadedData struct {
	FoodToEat  []int
	DaysToWait []int
}

func (a *CustomAgentEvo) AppendToMessageMemory(msg messages.Message, msgMemory []messages.Message) {
	msgMemory = append(msgMemory, msg)
}

func (a *CustomAgentEvo) NeighbourFoodEaten() food.FoodType {

	if a.CurrPlatFood() != -1 {
		if !a.PlatformOnFloor() && a.CurrPlatFood() != a.lastPlatFood {
			return a.lastPlatFood - a.CurrPlatFood()
		}
		return 0
	}
	return -1
}
func remove(slice []messages.Message, s int) []messages.Message {
	return append(slice[:s], slice[s+1:]...)
}

func (a *CustomAgentEvo) HasDayPassed() bool {
	if a.Age() != a.lastAge {
		a.lastAge = a.Age()
		return true
	}
	return false
}

func InitaliseParams(baseAgent *infra.Base) CustomAgentEvoParams {
	mydir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	file, _ := ioutil.ReadFile(fmt.Sprintf("%s/pkg/agents/team4/agent1/agentConfig.json", mydir))
	var data1 LoadedData
	_ = json.Unmarshal(file, &data1) //parse the config json file to get the coeffs for floor and HP equations

	data1.FoodToEat[0] = baseAgent.HealthInfo().HPReqCToW
	data1.DaysToWait[0] = int(baseAgent.HealthInfo().MaxDayCritical / 2)

	return CustomAgentEvoParams{ //initialise the parameters of the agent
		// scalingEquation: Equation{
		// 	coefficients: []float64{1, 1, 1, 1},
		// },
		// currentFloorScore: currentFloorScoreEquation,
		// currentHpScore:    currentHpScoreEquation,
		locked:            true,
		previousAge:       0,
		foodToEat:         data1.FoodToEat,
		daysToWait:        data1.DaysToWait,
		ageLastEaten:      0,
		trustscore:        100 * rand.Float64(), //random for now
		morality:          100 * rand.Float64(), //random for now
		traumaScaleFactor: 1,                    //0 for now
	}
}

func New(baseAgent *infra.Base) (infra.Agent, error) {
	// currentFloorScoreEquation, currentHpScoreEquation := GenerateEquations()
	//create other parameters
	return &CustomAgentEvo{
		Base:                baseAgent,
		params:              InitaliseParams(baseAgent),
		globalTrust:         0.0,                           // TODO: Amend values for correct agent behaviour
		globalTrustAdd:      9.0,                           // TODO: Amend values for correct agent behaviour
		globalTrustSubtract: -9.0,                          // TODO: Amend values for correct agent behaviour
		coefficients:        []float32{0.1, 0.2, 0.4, 0.5}, // TODO: Amend values for correct agent behaviour

		// Initialise the amount of food our agent intends to eat.
		intendedFoodTaken: 0,
		// Initialise the actual food taken on the previous run.
		lastFoodTaken: 0,

		// Initialise Agents individual message memory
		sentMessages:     []messages.Message{},
		responseMessages: []messages.Message{},
		// Define what message to send during a run.
		MessageToSend:    rand.Intn(8),
		lastPlatFood:     -1,
		maxFoodLimit:     50,
		messageCounter:   0,
		globalTrustLimit: 75,
		lastAge:          0,
	}, nil
}

func (a *CustomAgentEvo) Run() {

	if food.FoodType(a.CurrPlatFood()) != a.lastPlatFood && a.PlatformOnFloor() {
		a.lastPlatFood = a.CurrPlatFood()
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

	dayPass := false
	if a.Age() != a.params.previousAge {
		a.params.previousAge = a.Age()
		dayPass = true
	}
	dayPass = dayPass
	var healthStatus int
	var healthLevelSeparation = int(0.33 * float64(a.HealthInfo().MaxHP-a.HealthInfo().WeakLevel))

	if a.HP() <= a.HealthInfo().WeakLevel { //critical
		healthStatus = 0
		a.params.locked = false
		// if dayPass {
		// 	a.params.traumaScaleFactor = math.Min(200, a.params.traumaScaleFactor+0.03)
		// }
	}
	if !a.params.locked {
		if a.HP() <= a.HealthInfo().WeakLevel+healthLevelSeparation { //weak
			healthStatus = 1
			// if dayPass {
			// 	a.params.traumaScaleFactor = math.Min(200, a.params.traumaScaleFactor+0.02)
			// }
		} else if a.HP() <= a.HealthInfo().WeakLevel+2*healthLevelSeparation { //normal
			healthStatus = 2
			// if dayPass {
			// 	a.params.traumaScaleFactor = math.Max(0, a.params.traumaScaleFactor-0.02)
			// }
		} else { //strong
			healthStatus = 3
			// if dayPass {
			// 	a.params.traumaScaleFactor = math.Max(0, a.params.traumaScaleFactor-0.03)
			// }
		}
		a.params.locked = true
	}

	var foodEaten food.FoodType
	var err error

	if (a.Age()-a.params.ageLastEaten) >= a.params.daysToWait[healthStatus] || healthStatus == 0 {
		a.params.locked = false
		calculatedAmountToEat := a.params.traumaScaleFactor * float64(a.params.foodToEat[healthStatus])
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

	// // Using custom params to decide how much food to eat
	// scaledFloorScore := a.params.currentFloorScore.EvaluateEquation(a.Floor()) / a.params.scalingEquation.EvaluateEquation(a.Floor())
	// scaledHpScore := a.params.currentHpScore.EvaluateEquation(a.HP()) / a.params.scalingEquation.EvaluateEquation(a.HP())
	// foodToEat := food.FoodType(math.Max(0.0, 50*scaledFloorScore+50*scaledHpScore))
	// foodEaten := a.TakeFood(foodToEat)
	a.Log("team4EvoAgent reporting status:", infra.Fields{"floor": a.Floor(), "hp": a.HP(), "FoodToEat": a.params.foodToEat[healthStatus], "DaysToWait": a.params.daysToWait, "foodEaten": foodEaten})

	// // Trauma score for future use
	// if a.HP() < 25 {
	// 	a.params.trauma = math.Min(100, a.params.trauma+1)
	// } else if a.HP() > 75 {
	// 	a.params.trauma = math.Max(0, a.params.trauma-1)
	// }

}
