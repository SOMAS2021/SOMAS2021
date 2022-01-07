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
	foodToEat         map[string][]int // the amount of food to eat for various health levels
	intendedFoodToEat int
	waitProbability   map[string][]int // probability the agent will wait to eat on any given day
	ageLastEaten      int              // the age at which the agent last ate
	morality          float64          // the morality of the agent that determines how selfishly or selflessly the agent will act
	craving           int              // the amount of craving the agent has for food which effects the amount of food it is likely to eat

	intendedFoodToTake food.FoodType

	globalTrust              float64       // the overall trust the agent has in other agents in the tower
	coefficients             []float64     // the amount trust score changes by for certain actions
	lastFoodTaken            food.FoodType // food taken on the previous day
	msgToSendBuffer          []messages.Message
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
		// "selfish":  {baseAgent.HealthInfo().HPReqCToW, 61, 41, 0},  // TODO: to optimise more
		// "neutral":  {baseAgent.HealthInfo().HPReqCToW, 82, 68, 44}, // TODO: to optimise more
		// "selfless": {baseAgent.HealthInfo().HPReqCToW, 18, 0, 0},   // TODO: to optimise more
		"selfish":  {baseAgent.HealthInfo().HPReqCToW, 34, 52, 22},
		"neutral":  {baseAgent.HealthInfo().HPReqCToW, 10, 46, 10},
		"selfless": {baseAgent.HealthInfo().HPReqCToW, 3, 14, 0},
	}
	waitProbability := map[string][]int{
		"selfish":  {0, 60, 60, 25}, // TODO: to optimise more int(baseAgent.HealthInfo().MaxDayCritical / 2)
		"neutral":  {0, 21, 16, 64}, // TODO: to optimise more int(baseAgent.HealthInfo().MaxDayCritical / 2)
		"selfless": {0, 11, 35, 75}, // TODO: to optimise more
	}

	// func InitaliseParams(baseAgent *infra.Base) CustomAgentEvoParams {
	// 	foodToEat := map[string][]int{
	// 		// "selfish":  {baseAgent.HealthInfo().HPReqCToW, 61, 41, 0},  // TODO: to optimise more
	// 		// "neutral":  {baseAgent.HealthInfo().HPReqCToW, 82, 68, 44}, // TODO: to optimise more
	// 		// "selfless": {baseAgent.HealthInfo().HPReqCToW, 18, 0, 0},   // TODO: to optimise more
	// 		"selfish":  {baseAgent.HealthInfo().HPReqCToW, 3, 14, 0},
	// 		"neutral":  {baseAgent.HealthInfo().HPReqCToW, 3, 14, 0},
	// 		"selfless": {baseAgent.HealthInfo().HPReqCToW, 3, 14, 0},
	// 	}
	// 	waitProbability := map[string][]int{
	// 		"selfish":  {0, 11, 35, 75}, // TODO: to optimise more int(baseAgent.HealthInfo().MaxDayCritical / 2)
	// 		"neutral":  {0, 11, 35, 75}, // TODO: to optimise more int(baseAgent.HealthInfo().MaxDayCritical / 2)
	// 		"selfless": {0, 11, 35, 75}, // TODO: to optimise more
	// 	}

	return CustomAgentEvoParams{ //initialise the parameters of the agent
		foodToEat:         foodToEat,
		intendedFoodToEat: 0,
		waitProbability:   waitProbability,
		ageLastEaten:      0,
		morality:          100 * rand.Float64(), // TODO: Use this properly
		craving:           0,
		healthStatus:      3,

		globalTrust:        100.0,
		globalTrustLimits:  []int{40, 80},
		coefficients:       []float64{2, 4, 8}, // TODO: maybe train these co-efficients using evolutionary algorithm
		lastFoodTaken:      0,
		intendedFoodToTake: 0,

		messageCounter:           0,
		msgToSendBuffer:          []messages.Message{},
		sentMessages:             []messages.Message{},
		responseMessages:         []messages.Message{},
		requestLeaveFoodMessages: []messages.Message{},
		otherMessageBuffer:       []messages.Message{},

		lastPlatFood:       -1,
		lastTimeFoodSeen:   0,
		maxFoodLimit:       50,
		lastAge:            -1,
		maxFloor:           0,
		currentPersonality: "selfless",
	}
}

func New(baseAgent *infra.Base) (infra.Agent, error) {
	//create other parameters
	return &CustomAgentEvo{
		Base:   baseAgent,
		params: InitaliseParams(baseAgent),
	}, nil
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

func (a *CustomAgentEvo) updateLastTimeFoodSeen() {
	if a.CurrPlatFood() >= food.FoodType(a.params.foodToEat[a.params.currentPersonality][0]) {
		a.params.lastTimeFoodSeen++
	} else {
		a.params.lastTimeFoodSeen = 0
	}
}

func (a *CustomAgentEvo) updateHealthStatus(healthLevelSeparation int) {

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

func (a *CustomAgentEvo) setPersonality() {
	//prev_personality := a.params.currentPersonality
	if a.params.globalTrust < float64(a.params.globalTrustLimits[0]) {
		a.params.currentPersonality = "selfish"
	} else if a.params.globalTrust < float64(a.params.globalTrustLimits[1]) {
		a.params.currentPersonality = "neutral"
	} else {
		a.params.currentPersonality = "selfless"
	}
	// if prev_personality != a.params.currentPersonality {
	// 	identifier := a.ID().String()
	// 	fmt.Printf("%s : Personality changed to %s \n", identifier[0:2], a.params.currentPersonality)
	// }
}

func getMatchingSentMessage(a *CustomAgentEvo, resMsg messages.ResponseMessage) messages.Message {
	var retMsg messages.Message
	for _, sentMsg := range a.params.sentMessages { // Iterate through each sent message
		if resMsg.RequestID() == sentMsg.ID() {
			retMsg = sentMsg
			break
		}
	}
	return retMsg
}

func removeMatchingSentMessage(a *CustomAgentEvo, resMsg messages.ResponseMessage) {
	for i, sentMsg := range a.params.sentMessages { // Iterate through each sent message
		if resMsg.RequestID() == sentMsg.ID() {
			a.params.sentMessages = remove(a.params.sentMessages, i)
			break
		}
	}
}

func (a *CustomAgentEvo) verifyResponses() {
	if len(a.params.responseMessages) > 0 {
		for i, respMsg := range a.params.responseMessages { // Iterate through each response message
			resMsg, err := a.typeAssertResponseMessage(respMsg)
			if err != nil {
				log.Error(err)
			} else {
				sentMsg := getMatchingSentMessage(a, resMsg)
				isHandled := false
				if a.PlatformOnFloor() && sentMsg.MessageType() == messages.RequestLeaveFood && a.Floor()-resMsg.SenderFloor() == 1 { // Check if there are any responses messages.
					a.UpdateGlobalTrustReqLeaveFood(resMsg, sentMsg)
				} else if a.params.lastFoodTaken+a.CurrPlatFood() != a.params.lastPlatFood && sentMsg.MessageType() == messages.RequestTakeFood && a.Floor()-resMsg.SenderFloor() == -1 {
					a.UpdateGlobalTrustReqTakeFood(resMsg, sentMsg)
				}
				if isHandled {
					removeMatchingSentMessage(a, resMsg)
					a.params.responseMessages = remove(a.params.responseMessages, i)
				}
			}
		}
		return
	}
}

func (a *CustomAgentEvo) Run() {

	// fmt.Println("Food Required ", health.FoodRequired(11, 44, a.HealthInfo()))
	// update agents memory of the last time it saw food on the platform
	// fmt.Println("The day passed value is ", check)
	if a.HasDayPassed() {
		a.GenerateMessagesToSend()
		a.updateLastTimeFoodSeen()
		a.params.globalTrust *= 1.0
		a.updateCraving()
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

	a.updateHealthStatus(healthLevelSeparation)
	a.SendingMessage()
	a.CallHandleMessage() //call the relevant message handler

	var foodEaten food.FoodType
	var err error
	var calculatedAmountToEat food.FoodType

	a.setPersonality()

	if a.PlatformOnFloor() || (!a.PlatformOnFloor() && a.CurrPlatFood() != -1) {
		a.verifyResponses()
	}

	a.params.intendedFoodToEat = 0
	if rand.Intn(100) >= (a.params.waitProbability[a.params.currentPersonality][a.params.healthStatus]-a.params.craving) && !a.HasEaten() && a.PlatformOnFloor() {
		a.params.intendedFoodToTake = food.FoodType(a.params.foodToEat[a.params.currentPersonality][a.params.healthStatus]) // TODO: add floor

		if a.params.healthStatus != 0 {
			a.handleActiveTreatyConditions()
		}
		foodEaten, err = a.TakeFood(a.params.intendedFoodToTake)
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

	a.Log("team4EvoAgent reporting status:", infra.Fields{"floor": a.Floor(), "hp": a.HP(), "FoodToEat": calculatedAmountToEat, "WaitProbability": a.params.waitProbability, "foodEaten": foodEaten})
}
