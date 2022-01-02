package team5

import (
	"math"

	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/health"
	"github.com/google/uuid"
)

type Memory struct {
	foodTaken         food.FoodType //Store last known value of foodTaken by agent
	agentHP           int           //Store last known value of HP of agent
	intentionFood     food.FoodType //Store the last known value of the amount of food intended to take
	favour            int           // e.g. generosity; scale of 0 to 10, with 0 being least favoured and 10 being most favoured
	daysSinceLastSeen int           // days since last interaction
}

type CustomAgent5 struct {
	*infra.Base
	selfishness       int
	lastMeal          food.FoodType
	daysSinceLastMeal int
	hpAfterEating     int
	currentAimHP      int
	attemptFood       food.FoodType
	satisfaction      int
	rememberAge       int
	rememberFloor     int
	messagingCounter  int
	// Social network of other agents
	socialMemory      map[uuid.UUID]Memory
	surroundingAgents map[int]uuid.UUID
}

func New(baseAgent *infra.Base) (infra.Agent, error) {
	return &CustomAgent5{
		Base:        baseAgent,
		selfishness: 10, // of 0 to 10, with 10 being completely selfish, 0 being completely selfless
		lastMeal:    0,  //Stores value of the last amount of food taken
		// TODO: PARAMETRISE TO MAXHP
		hpAfterEating:     100,
		daysSinceLastMeal: 0,   //Count of how many days since last eating
		currentAimHP:      100, //Scale of 0 to 2, 0 being willing to lose health, 1 being maintaining health, 2 being gaining health
		attemptFood:       0,
		satisfaction:      0,                       //Scale of -3 to 3, with 3 being satisfied and unsatisfied
		rememberAge:       -1,                      // To check if a day has passed by our age increasing
		socialMemory:      map[uuid.UUID]Memory{},  // Memory of other agents, key is agent id
		surroundingAgents: make(map[int]uuid.UUID), //Map agent id's of surrounding floors relative to current floor
	}, nil
}

func PercentageHP(a *CustomAgent5) int {
	return int((float64(a.HP()) / float64(a.HealthInfo().MaxHP)) * 100.0)
}

func (a *CustomAgent5) restrictToRange(lowerBound, upperBound, num int) int {
	if num < lowerBound {
		return lowerBound
	}
	if num > upperBound {
		return upperBound
	}
	return num
}

func (a *CustomAgent5) memoryIdExists(id uuid.UUID) bool {
	_, exists := a.socialMemory[id]
	return exists
}

func (a *CustomAgent5) newMemory(id uuid.UUID) {
	a.socialMemory[id] = Memory{
		foodTaken:         100,
		agentHP:           a.HealthInfo().MaxHP,
		intentionFood:     100,
		favour:            0,
		daysSinceLastSeen: 0,
	}
}

func (a *CustomAgent5) updateFoodTakenMemory(id uuid.UUID, foodTaken food.FoodType) {
	mem := a.socialMemory[id]
	mem.foodTaken = foodTaken
	a.socialMemory[id] = mem
}

func (a *CustomAgent5) updateAgentHPMemory(id uuid.UUID, agentHP int) {
	mem := a.socialMemory[id]
	mem.agentHP = agentHP
	a.socialMemory[id] = mem
}

func (a *CustomAgent5) updateIntentionFoodMemory(id uuid.UUID, intentionFood food.FoodType) {
	mem := a.socialMemory[id]
	mem.intentionFood = intentionFood
	a.socialMemory[id] = mem
}

func (a *CustomAgent5) addToSocialFavour(id uuid.UUID, change int) {
	mem := a.socialMemory[id]
	mem.favour = a.restrictToRange(0, 10, mem.favour+change)
	a.socialMemory[id] = mem
	// a.socialMemory[id] = Memory{
	// 	//trust:             a.socialMemory[id].trust,
	// 	foodTaken:         a.socialMemory[id].foodTaken,
	// 	agentHP:           a.socialMemory[id].agentHP,
	// 	intentionFood:     a.socialMemory[id].intentionFood,
	// 	favour:            a.restrictToRange(0, 10, a.socialMemory[id].favour+change),
	// 	daysSinceLastSeen: a.socialMemory[id].daysSinceLastSeen,
	// }
}

func (a *CustomAgent5) incrementDaysSinceLastSeen() {
	for _, mem := range a.socialMemory {
		// mem := a.socialMemory[id]
		mem.daysSinceLastSeen++
		// a.socialMemory[id] = mem
	}
}

func (a *CustomAgent5) resetSocialKnowledge(id uuid.UUID) {
	mem := a.socialMemory[id]
	mem.foodTaken = 100
	mem.agentHP = a.HealthInfo().MaxHP
	mem.intentionFood = 100
	a.socialMemory[id] = mem
}

func (a *CustomAgent5) resetDaysSinceLastSeen(id uuid.UUID) {
	mem := a.socialMemory[id]
	mem.daysSinceLastSeen = 0
	a.socialMemory[id] = mem
}

func (a *CustomAgent5) updateAimHP() {
	a.currentAimHP = a.HealthInfo().MaxHP - ((10 - a.selfishness) * ((a.HealthInfo().MaxHP - a.HealthInfo().WeakLevel) / 10))
}

func (a *CustomAgent5) updateSelfishness() {
	a.selfishness = 10 - a.calculateAverageFavour()
}

func (a *CustomAgent5) updateSatisfaction() {
	if PercentageHP(a) >= 100 {
		a.satisfaction = 3
	}
	if a.lastMeal == 0 && a.satisfaction > -3 {
		a.satisfaction--
	}
	if PercentageHP(a) < 25 && a.satisfaction > -3 {
		a.satisfaction--
	}
	if PercentageHP(a) > 75 && a.satisfaction < 3 {
		a.satisfaction++
	}
}

func (a *CustomAgent5) GetMessages() {
	receivedMsg := a.ReceiveMessage()
	if receivedMsg != nil {
		receivedMsg.Visit(a)
	}
}

func (a *CustomAgent5) DailyMessages() {
	var msg messages.Message
	targetFloor := a.Floor() + 1
	switch a.messagingCounter {
	case 0:
		msg = messages.NewAskHPMessage(a.ID(), a.Floor(), targetFloor)
	case 1:
		msg = messages.NewAskFoodTakenMessage(a.ID(), a.Floor(), targetFloor)
	case 2:
		msg = messages.NewAskIntendedFoodIntakeMessage(a.ID(), a.Floor(), targetFloor)
	case 3:
		targetFloor = a.Floor() - 1
		msg = messages.NewAskHPMessage(a.ID(), a.Floor(), targetFloor)
	case 4:
		targetFloor = a.Floor() - 1
		msg = messages.NewAskFoodTakenMessage(a.ID(), a.Floor(), targetFloor)
	case 5:
		targetFloor = a.Floor() - 1
		msg = messages.NewAskIntendedFoodIntakeMessage(a.ID(), a.Floor(), targetFloor)
	default:
	}
	a.SendMessage(msg)
	a.Log("Team 5 agent sent message", infra.Fields{"msgType": msg.MessageType().String(), "senderFloor": a.Floor(), "targetFloor": targetFloor})

	a.messagingCounter++
}

func (a *CustomAgent5) ResetSurroundingAgents() {
	for floor := range a.surroundingAgents {
		delete(a.surroundingAgents, floor)
	}
}

func (a *CustomAgent5) calculateAverageFavour() int {
	sum := 0
	count := 0
	for floor := range a.surroundingAgents {
		sum += a.socialMemory[a.surroundingAgents[floor]].favour
		count++
	}
	if count == 0 {
		return 10 - a.selfishness
	}
	return sum / count
}

func (a *CustomAgent5) updateFavour() {
	for id, mem := range a.socialMemory {
		if mem.daysSinceLastSeen < 2 {
			judgement := (a.hpAfterEating - mem.agentHP) + int(a.lastMeal-mem.foodTaken) + int(a.calculateAttemptFood()-mem.intentionFood)
			// a.Log("I have judged an agent", infra.Fields{"judgement": judgement})
			if judgement > 0 {
				a.addToSocialFavour(id, 1)
			}
			if judgement < 0 {
				a.addToSocialFavour(id, int(math.Max(float64(judgement)/20, -3)))
			}
		}
		if mem.daysSinceLastSeen > 5 {
			a.resetSocialKnowledge(id)
		}
	}
}

func (a *CustomAgent5) dayPassed() {
	a.updateFavour()
	a.updateSelfishness()
	a.updateAimHP()
	a.attemptFood = a.calculateAttemptFood()

	// for id := range a.socialMemory {
	// 	a.Log("Memory at end of day", infra.Fields{"favour": a.socialMemory[id].favour, "agentHP": a.socialMemory[id].agentHP, "foodTaken": a.socialMemory[id].foodTaken, "intent": a.socialMemory[id].intentionFood})
	// }
	// a.Log("Selfishness at end of day", infra.Fields{"selfishness": a.selfishness})
	// a.Log("Aim HP at end of day", infra.Fields{"aim HP": a.currentAimHP})
	// a.Log("Surrounding Agent Knowledge at end of day", infra.Fields{"agent map": a.surroundingAgents})
	a.daysSinceLastMeal++
	a.incrementDaysSinceLastSeen()
	if a.rememberFloor != a.Floor() {
		a.ResetSurroundingAgents()
		a.rememberFloor = a.Floor()
	}
	a.messagingCounter = 0
}

func (a *CustomAgent5) calculateAttemptFood() food.FoodType {
	if a.HP() < a.HealthInfo().WeakLevel {
		// TODO: UPDATE THIS VALUE TO A PARAMETER
		return food.FoodType(3)
	}
	foodAttempt := health.FoodRequired(a.HP(), a.currentAimHP, a.HealthInfo())
	return food.FoodType(math.Min(a.HealthInfo().Tau*3, float64(foodAttempt)))
}

func (a *CustomAgent5) Run() {
	a.Log("Reporting agent state of team 5 agent", infra.Fields{"health": a.HP(), "floor": a.Floor()})

	//Check if a day has passed
	if a.Age() > a.rememberAge {
		a.dayPassed()
		a.rememberAge = a.Age()
	}

	a.GetMessages()
	a.DailyMessages()

	//When platform reaches our floor and we haven't tried to eat, then try to eat
	if a.CurrPlatFood() != -1 && !a.HasEaten() {
		lastMeal, err := a.TakeFood(a.attemptFood)
		if err != nil {
			switch err.(type) {
			case *infra.FloorError:
			case *infra.NegFoodError:
			case *infra.AlreadyEatenError:
			default:
			}
		}
		a.lastMeal = lastMeal
		if a.lastMeal > 0 {
			a.daysSinceLastMeal = 0
		}
		a.updateSatisfaction()
		a.hpAfterEating = a.HP()
		a.messagingCounter = 0
	}
}

//The message handler functions below are for a fully honest agent

func (a *CustomAgent5) HandleAskHP(msg messages.AskHPMessage) {
	reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), a.HP())
	a.SendMessage(reply)
	a.Log("Team 5 agent received an askHP message", infra.Fields{"sender floor": msg.SenderFloor(), "receiver floor": a.Floor()})
	if !a.memoryIdExists(msg.SenderID()) {
		a.newMemory(msg.SenderID())
	}
	a.resetDaysSinceLastSeen(msg.SenderID())
	a.surroundingAgents[msg.SenderFloor()-a.Floor()] = msg.SenderID()
}

func (a *CustomAgent5) HandleAskFoodTaken(msg messages.AskFoodTakenMessage) {
	reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), int(a.lastMeal))
	a.SendMessage(reply)
	a.Log("Team 5 agent received an askFoodTaken message", infra.Fields{"sender floor": msg.SenderFloor(), "receiver floor": a.Floor()})
	if !a.memoryIdExists(msg.SenderID()) {
		a.newMemory(msg.SenderID())
	}
	a.resetDaysSinceLastSeen(msg.SenderID())
	a.surroundingAgents[msg.SenderFloor()-a.Floor()] = msg.SenderID()
}

func (a *CustomAgent5) HandleAskIntendedFoodTaken(msg messages.AskIntendedFoodIntakeMessage) {
	amount := int(a.calculateAttemptFood())
	reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), amount)
	a.SendMessage(reply)
	a.Log("Team 5 agent received an askIntendedFoodTaken message", infra.Fields{"sender floor": msg.SenderFloor(), "receiver floor": a.Floor()})
	if !a.memoryIdExists(msg.SenderID()) {
		a.newMemory(msg.SenderID())
	}
	a.resetDaysSinceLastSeen(msg.SenderID())
	a.surroundingAgents[msg.SenderFloor()-a.Floor()] = msg.SenderID()
}

func (a *CustomAgent5) HandleRequestLeaveFood(msg messages.RequestLeaveFoodMessage) {
	amount := msg.Request()
	// Always set to false for now to prevent deception, needs some calculations to determine whether we will leave the requested amount
	reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), false)
	a.SendMessage(reply)
	a.Log("Team 5 agent received a requestLeaveFood message", infra.Fields{"sender floor": msg.SenderFloor(), "receiver floor": a.Floor(), "request amount": amount})
	if !a.memoryIdExists(msg.SenderID()) {
		a.newMemory(msg.SenderID())
	}
	a.resetDaysSinceLastSeen(msg.SenderID())
	a.surroundingAgents[msg.SenderFloor()-a.Floor()] = msg.SenderID()
}

func (a *CustomAgent5) HandleRequestTakeFood(msg messages.RequestTakeFoodMessage) {
	amount := food.FoodType(msg.Request())
	reponse := true
	if a.calculateAttemptFood() > amount {
		reponse = false
	}
	reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), reponse)
	a.SendMessage(reply)
	a.Log("Team 5 agent received a requestTakeFood message", infra.Fields{"sender floor": msg.SenderFloor(), "receiver floor": a.Floor(), "request amount": amount})
	if !a.memoryIdExists(msg.SenderID()) {
		a.newMemory(msg.SenderID())
	}
	a.resetDaysSinceLastSeen(msg.SenderID())
	a.surroundingAgents[msg.SenderFloor()-a.Floor()] = msg.SenderID()
}

func (a *CustomAgent5) HandleResponse(msg messages.BoolResponseMessage) {
	response := msg.Response()
	a.Log("Team 5 agent received a Response message", infra.Fields{"sender floor": msg.SenderFloor(), "receiver floor": a.Floor(), "response": response})
	if !a.memoryIdExists(msg.SenderID()) {
		a.newMemory(msg.SenderID())
	}
	a.resetDaysSinceLastSeen(msg.SenderID())
	a.surroundingAgents[msg.SenderFloor()-a.Floor()] = msg.SenderID()
}

func (a *CustomAgent5) HandleStateFoodTaken(msg messages.StateFoodTakenMessage) {
	statement := food.FoodType(msg.Statement())
	a.Log("Team 5 agent received a StateFoodTaken message", infra.Fields{"sender floor": msg.SenderFloor(), "receiver floor": a.Floor(), "statement": statement})
	if !a.memoryIdExists(msg.SenderID()) {
		a.newMemory(msg.SenderID())
	}
	a.resetDaysSinceLastSeen(msg.SenderID())
	a.surroundingAgents[msg.SenderFloor()-a.Floor()] = msg.SenderID()
	a.updateFoodTakenMemory(msg.SenderID(), statement)
	a.Log("New value of foodTaken", infra.Fields{"statement": statement, "memory": a.socialMemory[msg.SenderID()].foodTaken})
}

func (a *CustomAgent5) HandleStateHP(msg messages.StateHPMessage) {
	statement := msg.Statement()
	a.Log("Team 5 agent received a StateHP message", infra.Fields{"sender floor": msg.SenderFloor(), "receiver floor": a.Floor(), "statement": statement})
	if !a.memoryIdExists(msg.SenderID()) {
		a.newMemory(msg.SenderID())
	}
	a.resetDaysSinceLastSeen(msg.SenderID())
	a.surroundingAgents[msg.SenderFloor()-a.Floor()] = msg.SenderID()
	a.updateAgentHPMemory(msg.SenderID(), statement)
	a.Log("New value of agentHP", infra.Fields{"statement": statement, "memory": a.socialMemory[msg.SenderID()].agentHP})
}

func (a *CustomAgent5) HandleStateIntendedFoodTaken(msg messages.StateIntendedFoodIntakeMessage) {
	statement := food.FoodType(msg.Statement())
	a.Log("Team 5 agent received a StateIntendedFoodTaken message", infra.Fields{"sender floor": msg.SenderFloor(), "receiver floor": a.Floor(), "statement": statement})
	if !a.memoryIdExists(msg.SenderID()) {
		a.newMemory(msg.SenderID())
	}
	a.resetDaysSinceLastSeen(msg.SenderID())
	a.surroundingAgents[msg.SenderFloor()-a.Floor()] = msg.SenderID()
	a.updateIntentionFoodMemory(msg.SenderID(), statement)
	a.Log("New value of intendedFood", infra.Fields{"statement": statement, "memory": a.socialMemory[msg.SenderID()].intentionFood})
}
