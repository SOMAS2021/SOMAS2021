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
	log "github.com/sirupsen/logrus"
)

/*------------------------AGENT STRUCTURE ------------------------*/

type CustomAgentEvoParams struct {
	foodToEat          []int // the amount of food to eat for various health levels
	waitProbability    []int // probability the agent will wait to eat on any given day
	ageLastEaten       int   // the age at which the agent last ate
	craving            int   // the amount of craving the agent has for food which effects the amount of food it is likely to eat
	previousFloor      int
	intendedFoodToTake food.FoodType // amount of food the agent intends to take
	lastFoodTaken      food.FoodType // food taken on the previous day

	msgToSendBuffer          []messages.Message // buffer of messages the agent wants to send
	sentMessages             []messages.Message // list of all sent messages
	responseMessages         []messages.Message // buffer of all response messages recieved
	requestLeaveFoodMessages []messages.Message // buffer of requests to leave food recieved
	otherMessageBuffer       []messages.Message // buffer of all other messages recieved

	lastPlatFood          food.FoodType // last seen amount of food on the platform
	lastTimeFoodSeen      int           // number of days passed since seeing the desired amount of food on the platform
	receivedMessagesCount int           // count of the total messages recieved
	messageCounter        int           // the total number of messages we send in a day
	lastAge               int           // the age of the agent on the previous day
	healthStatus          int           // the current health status of the agent
}

type CustomAgentEvo struct {
	*infra.Base
	// new params
	params CustomAgentEvoParams
}

/*------------------------AGENT INITIALISATION------------------------*/

// Data struct to hold evo param values
type LoadedData struct {
	FoodToEat       []int
	WaitProbability []int
}

func initaliseParams(baseAgent *infra.Base) CustomAgentEvoParams {

	mydir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	file, _ := ioutil.ReadFile(fmt.Sprintf("%s/pkg/agents/team4/trainingAgent/configs/agentConfig.json", mydir))
	var data1 LoadedData
	_ = json.Unmarshal(file, &data1) //parse the config json file to get the coeffs for floor and HP equations

	data1.FoodToEat[0] = baseAgent.HealthInfo().HPReqCToW
	data1.WaitProbability[0] = 0

	return CustomAgentEvoParams{
		foodToEat:                data1.FoodToEat,
		waitProbability:          data1.WaitProbability,
		ageLastEaten:             0,
		craving:                  0,
		healthStatus:             3,
		previousFloor:            -1,
		lastFoodTaken:            0,
		intendedFoodToTake:       0,
		messageCounter:           0,
		msgToSendBuffer:          []messages.Message{},
		sentMessages:             []messages.Message{},
		responseMessages:         []messages.Message{},
		requestLeaveFoodMessages: []messages.Message{},
		otherMessageBuffer:       []messages.Message{},
		receivedMessagesCount:    0,
		lastPlatFood:             -1,
		lastTimeFoodSeen:         0,
		lastAge:                  -1,
	}
}

func New(baseAgent *infra.Base) (infra.Agent, error) {
	//create other parameters
	return &CustomAgentEvo{
		Base:   baseAgent,
		params: initaliseParams(baseAgent),
	}, nil
}

/*------------------------AGENT UTILITIES ------------------------*/

func (a *CustomAgentEvo) hasDayPassed() bool {

	if a.Age() > a.params.lastAge {
		a.params.lastAge = a.Age()
		return true
	}
	return false
}

func (a *CustomAgentEvo) updateLastTimeFoodSeen() {
	if a.CurrPlatFood() >= food.FoodType(a.params.foodToEat[0]) {
		a.params.lastTimeFoodSeen++
	} else {
		a.params.lastTimeFoodSeen = 0
	}
}

func (a *CustomAgentEvo) updateCraving() {

	if a.params.lastTimeFoodSeen > 0 {
		a.params.craving += a.params.lastTimeFoodSeen
	} else {
		a.params.craving -= 2
	}

	if a.params.craving > 100 {
		a.params.craving = 100
	} else if a.params.craving < 0 {
		a.params.craving = 0
	}
}

func (a *CustomAgentEvo) hasReshuffled() bool {
	if a.params.previousFloor != a.Floor() {
		a.params.previousFloor = a.Floor()
		return true
	}
	return false
}

/*------------------------AGENT RUN------------------------*/

func (a *CustomAgentEvo) Run() {

	// agent health band threshold
	healthLevelSeparation := int(0.33 * float64(a.HealthInfo().MaxHP-a.HealthInfo().WeakLevel))

	if a.hasReshuffled() {
		//reset the arrays and maps when we have reshuffled
		a.params.msgToSendBuffer = []messages.Message{}
		a.params.sentMessages = []messages.Message{}
		a.params.responseMessages = []messages.Message{}
		a.params.requestLeaveFoodMessages = []messages.Message{}
		a.params.otherMessageBuffer = []messages.Message{}
		a.params.receivedMessagesCount = 0
	}

	if a.hasDayPassed() {
		//generates messages to send per day altogether and updates params for each new day
		a.generateMessagesToSend()
		a.updateLastTimeFoodSeen()
		a.updateCraving()
	}

	//update last food amount seen on the platform
	if food.FoodType(a.CurrPlatFood()) != a.params.lastPlatFood && a.PlatformOnFloor() {
		a.params.lastPlatFood = a.CurrPlatFood()
	}

	// receive message
	a.getMessage()

	a.params.healthStatus = getHealthStatus(a.HealthInfo(), healthLevelSeparation, a.HP())

	//call the relevant message handler
	a.callHandleMessage()

	a.sendingMessage()

	//----------------------EATING FOOD -------------------------//

	var foodEaten food.FoodType
	var err error

	if rand.Intn(100) >= (a.params.waitProbability[a.params.healthStatus]-a.params.craving) && !a.HasEaten() && a.PlatformOnFloor() {
		a.params.intendedFoodToTake = food.FoodType(a.params.foodToEat[a.params.healthStatus])

		// prioritises treaties over our current food behaviour
		a.handleActiveTreatyConditions()

		// eat and update last food
		foodEaten, err = a.TakeFood(a.params.intendedFoodToTake)
		a.params.lastFoodTaken = foodEaten
		// records whether there was food taken
		if foodEaten > 0 {
			a.params.ageLastEaten = a.Age()
		}

		// setting craving to 0 if we are fulfilled otherwise reduce the craving proportional to foodEaten
		if a.params.foodToEat[a.params.healthStatus] != 0 {
			a.params.craving -= a.params.craving * (int(foodEaten) / a.params.foodToEat[a.params.healthStatus])
		} else {
			a.params.craving = 0
		}

		// error handling
		if err != nil {
			switch err.(type) {
			case *infra.FloorError:
				log.Error("Simulation - team4/agentTraining.go: \t FloorError: is the platform on your floor?")
			case *infra.NegFoodError:
				log.Error("Simulation - team4/agentTraining.go: \t NegFoodError: is calculatedAmountToEat negative?")
			case *infra.AlreadyEatenError:
				log.Error("Simulation - team4/agentTraining.go: \t AlreadyEatenError: Have you already eaten?")
			default:
			}
		}
	}

	//log agent state
	a.Log("team4EvoAgent reporting status:", infra.Fields{"floor": a.Floor(), "hp": a.HP(), "FoodToEat": a.params.intendedFoodToTake, "WaitProbability": a.params.waitProbability, "foodEaten": foodEaten})
}
