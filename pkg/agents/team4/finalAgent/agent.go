package team4EvoAgent

import (
	"math"
	"math/rand"

	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
	log "github.com/sirupsen/logrus"
)

type CustomAgentEvoParams struct {
	foodToEat       map[string][]int // the amount of food to eat for various health levels
	waitProbability map[string][]int // probability the agent will wait to eat on any given day
	ageLastEaten    int              // the age at which the agent last ate
	morality        float64          // the morality of the agent that determines how selfishly or selflessly the agent will act
	craving         int              // the amount of craving the agent has for food which effects the amount of food it is likely to eat

	globalTrust              float32            // the overall trust the agent has in other agents in the tower
	coefficients             []float32          // the amount trust score changes by for certain actions
	lastFoodTaken            food.FoodType      // food taken on the previous day
	sentMessages             []messages.Message // TODO: make it a map hashed by messageIDs
	responseMessages         []messages.Message // TODO: make it a map hashed by messageIDs
	requestLeaveFoodMessages []messages.Message
	otherMessageBuffer       []messages.Message
	lastPlatFood             food.FoodType // last seen food on the platform
	lastTimeFoodSeen         int           // number of days passed since seeing the desired amount of food on the platform
	maxFoodLimit             food.FoodType // maximum food we want to allow others to eat
	messageCounter           int           // the total number of messages we send in a day
	globalTrustLimits        []int         // limit to check what personality to choose
	lastAge                  int           // the age of the agent on the previous day
	healthStatus             int
	maxFloor                 int
	currentPersonality       string
}

type CustomAgentEvo struct {
	*infra.Base
	// new params
	params CustomAgentEvoParams
}

func InitaliseParams(baseAgent *infra.Base) CustomAgentEvoParams {
	foodToEat := map[string][]int{
		"selfish":  {baseAgent.HealthInfo().HPReqCToW, 61, 41, 0},  // TODO: to optimise more
		"neutral":  {baseAgent.HealthInfo().HPReqCToW, 82, 68, 44}, // TODO: to optimise more
		"selfless": {baseAgent.HealthInfo().HPReqCToW, 18, 0, 0},   // TODO: to optimise more
	}
	waitProbability := map[string][]int{
		"selfish":  {0, 95, 72, 18}, // TODO: to optimise more int(baseAgent.HealthInfo().MaxDayCritical / 2)
		"neutral":  {0, 63, 74, 37}, // TODO: to optimise more int(baseAgent.HealthInfo().MaxDayCritical / 2)
		"selfless": {0, 98, 25, 47}, // TODO: to optimise more
	}

	return CustomAgentEvoParams{ //initialise the parameters of the agent
		foodToEat:       foodToEat,
		waitProbability: waitProbability,
		ageLastEaten:    0,
		morality:        100 * rand.Float64(), // TODO: Use this properly
		craving:         0,
		healthStatus:    3,

		globalTrust:       60.0,
		globalTrustLimits: []int{40, 80},
		coefficients:      []float32{2, 4, 8}, // TODO: maybe train these co-efficients using evolutionary algorithm
		lastFoodTaken:     0,

		messageCounter:           0,
		sentMessages:             []messages.Message{},
		responseMessages:         []messages.Message{},
		requestLeaveFoodMessages: []messages.Message{},
		otherMessageBuffer:       []messages.Message{},

		lastPlatFood:       -1,
		lastTimeFoodSeen:   0,
		maxFoodLimit:       50,
		lastAge:            -1,
		maxFloor:           0,
		currentPersonality: "neutral",
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

	if a.Age() > a.params.lastAge {
		a.params.lastAge = a.Age()
		return true
	}
	return false
}

func (a *CustomAgentEvo) UpdateLastTimeFoodSeen() {
	if a.CurrPlatFood() >= food.FoodType(a.params.foodToEat[a.params.currentPersonality][0]) {
		a.params.lastTimeFoodSeen++
	} else {
		a.params.lastTimeFoodSeen = 0
	}
}

////--------------------Message Parsing--------------------////

func (a *CustomAgentEvo) GetMessage() { //move this function to messages.go
	receivedMsg := a.ReceiveMessage()

	if receivedMsg != nil {
		if receivedMsg.MessageType() == messages.RequestLeaveFood {
			a.params.requestLeaveFoodMessages = append(a.params.requestLeaveFoodMessages, receivedMsg)
		} else {
			a.params.otherMessageBuffer = append(a.params.otherMessageBuffer, receivedMsg)
		}
	}
}

func (a *CustomAgentEvo) CallHandleMessage() { //move this function to messages.go
	if a.PlatformOnFloor() && len(a.params.requestLeaveFoodMessages) > 0 {
		a.params.requestLeaveFoodMessages[0].Visit(a)
		remove(a.params.requestLeaveFoodMessages, 0)
	} else if len(a.params.otherMessageBuffer) > 0 {
		a.params.otherMessageBuffer[0].Visit(a)
		remove(a.params.otherMessageBuffer, 0)
	} else {
		a.Log("I got no messages")
	}
}

////--------------------Message Parsing--------------------////

func (a *CustomAgentEvo) UpdateHealthStatus(healthLevelSeparation int) {

	if a.HP() <= a.HealthInfo().WeakLevel { //critical
		a.params.healthStatus = 0
	} else if a.HP() <= a.HealthInfo().WeakLevel+healthLevelSeparation { //weak
		a.params.healthStatus = 1
	} else if a.HP() <= a.HealthInfo().WeakLevel+2*healthLevelSeparation { //normal
		a.params.healthStatus = 2
	} else { //strong
		a.params.healthStatus = 3
	}
}

func (a *CustomAgentEvo) UpdateCraving() {

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

func (a *CustomAgentEvo) SetPersonality() {
	// prev_personality := a.params.currentPersonality
	if a.params.globalTrust < float32(a.params.globalTrustLimits[0]) {
		a.params.currentPersonality = "selfish"
	} else if a.params.globalTrust < float32(a.params.globalTrustLimits[1]) {
		a.params.currentPersonality = "neutral"
	} else {
		a.params.currentPersonality = "selfless"
	}
	// if prev_personality != a.params.currentPersonality {
	// 	identifier := a.ID().String()
	// 	fmt.Printf("%s : Personality changed to %s \n", identifier[0:2], a.params.currentPersonality)
	// }
}

func (a *CustomAgentEvo) Run() {

	// update agents memory of the last time it saw food on the platform
	// fmt.Println("The day passed value is ", check)
	if a.HasDayPassed() {
		a.UpdateLastTimeFoodSeen()
		a.params.globalTrust *= 0.85
	}

	// update agent's perception of maxFloor
	a.params.maxFloor = int(math.Max(float64(a.params.maxFloor), float64(a.Floor()))) // math.Max only take in floats while we require an integer floor.

	if food.FoodType(a.CurrPlatFood()) != a.params.lastPlatFood && a.PlatformOnFloor() {
		a.params.lastPlatFood = a.CurrPlatFood()
	}
	healthLevelSeparation := int(0.33 * float64(a.HealthInfo().MaxHP-a.HealthInfo().WeakLevel))

	// receive message
	a.GetMessage()

	// TODO: Define a threshold limit for other agents to respond to our sent message.
	a.SendingMessage()

	a.UpdateHealthStatus(healthLevelSeparation)
	a.CallHandleMessage() //call the relevant message handler

	var foodEaten food.FoodType
	var err error
	var calculatedAmountToEat food.FoodType

	a.SetPersonality()

	if rand.Intn(100) >= (a.params.waitProbability[a.params.currentPersonality][a.params.healthStatus]-a.params.craving) && !a.HasEaten() && a.PlatformOnFloor() {
		calculatedAmountToEat = food.FoodType(a.params.foodToEat[a.params.currentPersonality][a.params.healthStatus]) // TODO: add floor

		foodEaten, err = a.TakeFood(calculatedAmountToEat)
		if foodEaten > 0 {
			a.params.ageLastEaten = a.Age()
		}

		// setting craving to 0 if we are fulfilled otherwise reduce the craving proportional to foodEaten
		if a.params.foodToEat[a.params.currentPersonality][a.params.healthStatus] != 0 {
			a.params.craving -= a.params.craving * (int(foodEaten) / a.params.foodToEat[a.params.currentPersonality][a.params.healthStatus])
		} else {
			a.params.craving = 0
		}

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
	// identifier := a.ID().String()
	// fmt.Printf("%s : Trust Score is  %f \n", identifier[0:2], a.params.globalTrust)

	a.Log("team4EvoAgent reporting status:", infra.Fields{"floor": a.Floor(), "hp": a.HP(), "FoodToEat": calculatedAmountToEat, "WaitProbability": a.params.waitProbability, "foodEaten": foodEaten})
}
