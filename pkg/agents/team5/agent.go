package team5

import (
	"math"

	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
	"github.com/google/uuid"
)

type Memory struct {
	trust             int // scale of -5 to 5, with -5 being least trustworthy and 5 being most trustworthy, 0 is neutral
	favour            int // e.g. generosity; scale of -5 to 5, with -5 being least favoured and 5 being most favoured, 0 is neutral
	daysSinceLastSeen int // days since last interaction
}

type CustomAgent5 struct {
	*infra.Base
	selfishness       int
	lastMeal          food.FoodType
	daysSinceLastMeal int
	// TODO: Change this to an enum
	currentAim   int
	satisfaction int
	// TODO: Check difference between this and HasEaten()
	// If true, then agent will attempt to eat
	attemptToEat     bool
	rememberAge      int
	rememberFloor    int
	messagingCounter int
	// Social network of other agents
	socialMemory      map[uuid.UUID]Memory
	surroundingAgents map[int]uuid.UUID
}

func New(baseAgent *infra.Base) (infra.Agent, error) {
	return &CustomAgent5{
		Base:              baseAgent,
		selfishness:       3,                       // of 0 to 3, with 3 being completely selfish, 0 being completely selfless
		lastMeal:          0,                       //Stores value of the last amount of food taken
		daysSinceLastMeal: 0,                       //Count of how many days since last eating
		currentAim:        0,                       //Scale of 0 to 2, 0 being willing to lose health, 1 being maintaining health, 2 being gaining health
		satisfaction:      0,                       //Scale of -3 to 3, with 3 being satisfied and unsatisfied
		attemptToEat:      true,                    //Variable needed to check if we have already attempted to eat on a day
		rememberAge:       0,                       // To check if a day has passed by our age increasing
		socialMemory:      map[uuid.UUID]Memory{},  // Memory of other agents, key is agent id
		surroundingAgents: make(map[int]uuid.UUID), //Map agent id's of surrounding floors relative to current floor
	}, nil
}

func PercentageHP(a *CustomAgent5) int {
	return int((float64(a.HP()) / float64(a.HealthInfo().MaxHP)) * 100.0)
}

// force number to be in range because min and max only takes float64
func restrictToRange(lowerBound, upperBound, num int) int {
	if num < lowerBound {
		return lowerBound
	}
	if num > upperBound {
		return upperBound
	}
	return num
}

// Checks if agent id exists in socialMemory.
func (a *CustomAgent5) memoryIdExists(id uuid.UUID) bool {
	_, exists := a.socialMemory[id]
	return exists
}

// Initialises memory of new agent
func (a *CustomAgent5) newMemory(id uuid.UUID) {
	a.socialMemory[id] = Memory{
		trust:             0,
		favour:            0,
		daysSinceLastSeen: 0,
	}
}

// Changes trust in socialMemory, change can be negative
func (a *CustomAgent5) addToSocialTrust(id uuid.UUID, change int) {
	a.socialMemory[id] = Memory{
		trust:             restrictToRange(-5, 5, a.socialMemory[id].trust+change),
		favour:            a.socialMemory[id].favour,
		daysSinceLastSeen: a.socialMemory[id].daysSinceLastSeen,
	}
}

// Changes favour in socialMemory, change can be negative
func (a *CustomAgent5) addToSocialFavour(id uuid.UUID, change int) {
	a.socialMemory[id] = Memory{
		trust:             a.socialMemory[id].trust,
		favour:            restrictToRange(-5, 5, a.socialMemory[id].favour+change),
		daysSinceLastSeen: a.socialMemory[id].daysSinceLastSeen,
	}
}

// Increments all daysSinceLastSeen by 1
func (a *CustomAgent5) incrementDaysSinceLastSeen() {
	for id := range a.socialMemory {
		a.socialMemory[id] = Memory{
			trust:             a.socialMemory[id].trust,
			favour:            a.socialMemory[id].favour,
			daysSinceLastSeen: a.socialMemory[id].daysSinceLastSeen + 1,
		}
	}
}

// Sets daysSinceLastSeen to 0 of given agent id
func (a *CustomAgent5) resetDaysSinceLastSeen(id uuid.UUID) {
	a.socialMemory[id] = Memory{
		trust:             a.socialMemory[id].trust,
		favour:            a.socialMemory[id].favour,
		daysSinceLastSeen: 0,
	}
}

func (a *CustomAgent5) updateAim() {
	switch {
	case a.selfishness >= 3:
		// If fully selfish always try to gain health
		a.currentAim = 2
	case PercentageHP(a) > 80 && a.selfishness == 2:
		// Try to maintain health if near max health if mostly selfish
		a.currentAim = 1
	case PercentageHP(a) > 80 && a.selfishness < 2:
		// Willing to lose health near max health if mostly or completely selfless
		a.currentAim = 0
	case PercentageHP(a) > 50 && a.selfishness == 2:
		// Try to gain health if mostly selfish when above half health
		a.currentAim = 2
	case PercentageHP(a) > 50 && a.selfishness == 1:
		// Try to maintain half health even if being mostly selfless
		a.currentAim = 1
	case PercentageHP(a) > 50 && a.selfishness == 0:
		// Willing to lose health if being completely selfless
		a.currentAim = 0
	case PercentageHP(a) > 10 && a.selfishness >= 1:
		// Try to gain health if less than half health and being anything but completely selfless
		a.currentAim = 2
	default:
		// Default to maintain health if being completely selfless at less than half health
		a.currentAim = 1
	}
}

func (a *CustomAgent5) updateSelfishness() {
	if a.satisfaction == 3 {
		a.selfishness--
	}
	if a.satisfaction < 0 || a.daysSinceLastMeal > 2 {
		a.selfishness++
	}
	//The above is a basic implementation for now while messaging is not functional
	//Once messages are implemented this function will be dependent on our social network and treaties etc
}

// This should probably be done inside health.
func (a *CustomAgent5) foodGain() food.FoodType {
	return food.FoodType(a.HealthInfo().Tau * 3)
}

// This should probably be done inside health.
func (a *CustomAgent5) foodMaintain() food.FoodType {
	return food.FoodType(a.HealthInfo().Tau * math.Log(1-(float64((a.HP()+30))/(3*a.HealthInfo().Width))) * -1)
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

// func (a *CustomAgent5) SendMessages() {
// 	//function that will send all messages we need to the other agents
// }

func (a *CustomAgent5) DailyMessages() {
	sendingToFloor := 0
	switch a.messagingCounter {
	case 0:
		msg := messages.NewAskHPMessage(a.ID(), a.Floor())
		sendingToFloor = 1
		a.Log("Team 5 agent is sending an Ask HP Message", infra.Fields{"sender floor": a.Floor(), "sending to floor": a.Floor() + sendingToFloor})
		a.SendMessage(sendingToFloor, msg)
		a.messagingCounter++
	case 1:
		msg := messages.NewAskFoodTakenMessage(a.ID(), a.Floor())
		sendingToFloor = 1
		a.Log("Team 5 agent is sending an Ask Food Taken Message", infra.Fields{"sender floor": a.Floor(), "sending to floor": a.Floor() + sendingToFloor})
		a.SendMessage(sendingToFloor, msg)
		a.messagingCounter++
	case 2:
		msg := messages.NewAskIntendedFoodIntakeMessage(a.ID(), a.Floor())
		sendingToFloor = 1
		a.Log("Team 5 agent is sending an Ask Intention Message", infra.Fields{"sender floor": a.Floor(), "sending to floor": a.Floor() + sendingToFloor})
		a.SendMessage(sendingToFloor, msg)
		a.messagingCounter++
	case 3:
		msg := messages.NewAskHPMessage(a.ID(), a.Floor())
		sendingToFloor = -1
		a.Log("Team 5 agent is sending an Ask HP Message", infra.Fields{"sender floor": a.Floor(), "sending to floor": a.Floor() + sendingToFloor})
		a.SendMessage(sendingToFloor, msg)
		a.messagingCounter++
	case 4:
		msg := messages.NewAskFoodTakenMessage(a.ID(), a.Floor())
		sendingToFloor = -1
		a.Log("Team 5 agent is sending an Ask Food Taken Message", infra.Fields{"sender floor": a.Floor(), "sending to floor": a.Floor() + sendingToFloor})
		a.SendMessage(sendingToFloor, msg)
		a.messagingCounter++
	case 5:
		msg := messages.NewAskIntendedFoodIntakeMessage(a.ID(), a.Floor())
		sendingToFloor = -1
		a.Log("Team 5 agent is sending an Ask Intention Message", infra.Fields{"sender floor": a.Floor(), "sending to floor": a.Floor() + sendingToFloor})
		a.SendMessage(sendingToFloor, msg)
		a.messagingCounter++
	}
}

func (a *CustomAgent5) ResetSurroundingAgents() {
	a.surroundingAgents = make(map[int]uuid.UUID)
}

func (a *CustomAgent5) dayPassed() {
	a.daysSinceLastMeal++
	a.incrementDaysSinceLastSeen()
	if a.rememberFloor != a.Floor() {
		a.ResetSurroundingAgents()
		a.rememberFloor = a.Floor()
	}
	//Also update daySinceLastSeen for memory here
	a.messagingCounter = 0
}

func (a *CustomAgent5) Run() {
	a.Log("Reporting agent state of team 5 agent", infra.Fields{"health": a.HP(), "floor": a.Floor()})
	a.GetMessages()
	a.DailyMessages()
	a.updateSelfishness()
	a.updateAim()
	attemptFood := food.FoodType(0)
	if a.HP() < 10 {
		//No point taking more food than 1 if in critical state as we will only reach 10hp with any amount of food > 0
		attemptFood = 1
	} else {
		if a.currentAim == 1 {
			attemptFood = a.foodMaintain()
		} else if a.currentAim == 2 {
			attemptFood = a.foodGain()
		}
	}

	// placeholder for memory functions
	agentAbove := uuid.New() // random name
	agentBelow := uuid.New() // random name
	if !a.memoryIdExists(agentAbove) {
		a.newMemory(agentAbove)
	}

	// placeholder for memory functions
	expectedFood := 55    // random number
	if expectedFood > 0 { // if we have an expectation
		a.resetDaysSinceLastSeen(agentAbove)
		trustChange := (int(a.CurrPlatFood()) - expectedFood + 20) / 5
		a.addToSocialTrust(agentAbove, trustChange)
	}

	// placeholder for memory functions
	if a.CurrPlatFood() != -1 {
		// rage is high when platform has little food compared to what agent wants
		rageToAgentAbove := 3 * int((float64(attemptFood)/float64(a.CurrPlatFood()))-2.5)
		a.addToSocialFavour(agentAbove, rageToAgentAbove)

		agentBelowRequestedFood := 30 // random number
		// e.g. if agent below requested food
		if agentBelowRequestedFood > 0 {
			a.resetDaysSinceLastSeen(agentBelow)
			// rage is low when agent below requests less food than food on platform minus our agent's required food
			rageToAgentBelow := 3 * int((float64(agentBelowRequestedFood)/float64(a.CurrPlatFood()-attemptFood))-1.5)
			a.addToSocialFavour(agentBelow, rageToAgentBelow)
		}
	}

	//When platform reaches our floor and we haven't tried to eat, then try to eat
	if a.CurrPlatFood() != -1 && a.attemptToEat {
		lastMeal, err := a.TakeFood(attemptFood)
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
			a.Log("Team 5 agent has taken food", infra.Fields{"amount": a.lastMeal})
			a.daysSinceLastMeal = 0
		}
		a.updateSatisfaction()
		a.attemptToEat = false
	}

	if a.Age() > a.rememberAge {
		//Check for if a day has passed
		a.dayPassed()
		a.attemptToEat = true
		a.rememberAge = a.Age()
	}

}

//The message handler functions below are for a fully honest agent

func (a *CustomAgent5) HandleAskHP(msg messages.AskHPMessage) {
	reply := msg.Reply(a.ID(), a.Floor(), a.HP())
	a.SendMessage(msg.SenderFloor()-a.Floor(), reply)
	a.Log("Team 5 agent recieved an askHP message", infra.Fields{"sender floor": msg.SenderFloor(), "reciever floor": a.Floor()})
	if a.memoryIdExists(msg.SenderID()) {
		a.newMemory(msg.SenderID())
	}
	a.resetDaysSinceLastSeen(msg.SenderID())
	a.surroundingAgents[msg.SenderFloor()-a.Floor()] = msg.SenderID()
}

func (a *CustomAgent5) HandleAskFoodTaken(msg messages.AskFoodTakenMessage) {
	reply := msg.Reply(a.ID(), a.Floor(), int(a.lastMeal))
	a.SendMessage(msg.SenderFloor()-a.Floor(), reply)
	a.Log("Team 5 agent recieved an askFoodTaken message", infra.Fields{"sender floor": msg.SenderFloor(), "reciever floor": a.Floor()})
	if a.memoryIdExists(msg.SenderID()) {
		a.newMemory(msg.SenderID())
	}
	a.resetDaysSinceLastSeen(msg.SenderID())
	a.surroundingAgents[msg.SenderFloor()-a.Floor()] = msg.SenderID()
}

func (a *CustomAgent5) HandleAskIntendedFoodTaken(msg messages.AskIntendedFoodIntakeMessage) {
	amount := 0
	if a.currentAim == 2 {
		amount = int(a.foodGain())
	}
	if a.currentAim == 1 {
		amount = int(a.foodMaintain())
	}
	reply := msg.Reply(a.ID(), a.Floor(), amount)
	a.SendMessage(msg.SenderFloor()-a.Floor(), reply)
	a.Log("Team 5 agent recieved an askIntendedFoodTaken message", infra.Fields{"sender floor": msg.SenderFloor(), "reciever floor": a.Floor()})
	if a.memoryIdExists(msg.SenderID()) {
		a.newMemory(msg.SenderID())
	}
	a.resetDaysSinceLastSeen(msg.SenderID())
	a.surroundingAgents[msg.SenderFloor()-a.Floor()] = msg.SenderID()
}

func (a *CustomAgent5) HandleRequestLeaveFood(msg messages.RequestLeaveFoodMessage) {
	amount := msg.Request()
	reply := msg.Reply(a.ID(), a.Floor(), false) //Always set to false for now to prevent deception, needs some calculations to determine whether we will leave the requested amount
	a.SendMessage(msg.SenderFloor()-a.Floor(), reply)
	a.Log("Team 5 agent recieved a requestLeaveFood message", infra.Fields{"sender floor": msg.SenderFloor(), "reciever floor": a.Floor(), "request amount": amount})
	if a.memoryIdExists(msg.SenderID()) {
		a.newMemory(msg.SenderID())
	}
	a.resetDaysSinceLastSeen(msg.SenderID())
	a.surroundingAgents[msg.SenderFloor()-a.Floor()] = msg.SenderID()
}

func (a *CustomAgent5) HandleRequestTakeFood(msg messages.RequestTakeFoodMessage) {
	amount := msg.Request()
	reponse := true
	if (a.currentAim == 2 && amount > int(a.foodGain())) || (a.currentAim == 1 && amount > int(a.foodMaintain())) {
		reponse = false
	}
	reply := msg.Reply(a.ID(), a.Floor(), reponse)
	a.SendMessage(msg.SenderFloor()-a.Floor(), reply)
	a.Log("Team 5 agent recieved a requestTakeFood message", infra.Fields{"sender floor": msg.SenderFloor(), "reciever floor": a.Floor(), "request amount": amount})
	if a.memoryIdExists(msg.SenderID()) {
		a.newMemory(msg.SenderID())
	}
	a.resetDaysSinceLastSeen(msg.SenderID())
	a.surroundingAgents[msg.SenderFloor()-a.Floor()] = msg.SenderID()
}

func (a *CustomAgent5) HandleResponse(msg messages.BoolResponseMessage) {
	response := msg.Response()
	a.Log("Team 5 agent recieved a Response message", infra.Fields{"sender floor": msg.SenderFloor(), "reciever floor": a.Floor(), "response": response})
	if a.memoryIdExists(msg.SenderID()) {
		a.newMemory(msg.SenderID())
	}
	a.resetDaysSinceLastSeen(msg.SenderID())
	a.surroundingAgents[msg.SenderFloor()-a.Floor()] = msg.SenderID()
}

func (a *CustomAgent5) HandleStateFoodTaken(msg messages.StateFoodTakenMessage) {
	statement := msg.Statement()
	a.Log("Team 5 agent recieved a StateFoodTaken message", infra.Fields{"sender floor": msg.SenderFloor(), "reciever floor": a.Floor(), "statement": statement})
	if a.memoryIdExists(msg.SenderID()) {
		a.newMemory(msg.SenderID())
	}
	a.resetDaysSinceLastSeen(msg.SenderID())
	a.surroundingAgents[msg.SenderFloor()-a.Floor()] = msg.SenderID()
}

func (a *CustomAgent5) HandleStateHP(msg messages.StateHPMessage) {
	statement := msg.Statement()
	a.Log("Team 5 agent recieved a StateHP message", infra.Fields{"sender floor": msg.SenderFloor(), "reciever floor": a.Floor(), "statement": statement})
	if a.memoryIdExists(msg.SenderID()) {
		a.newMemory(msg.SenderID())
	}
	a.resetDaysSinceLastSeen(msg.SenderID())
	a.surroundingAgents[msg.SenderFloor()-a.Floor()] = msg.SenderID()
}

func (a *CustomAgent5) HandleStateIntendedFoodTaken(msg messages.StateIntendedFoodIntakeMessage) {
	statement := msg.Statement()
	a.Log("Team 5 agent recieved a StateIntendedFoodTaken message", infra.Fields{"sender floor": msg.SenderFloor(), "reciever floor": a.Floor(), "statement": statement})
	if a.memoryIdExists(msg.SenderID()) {
		a.newMemory(msg.SenderID())
	}
	a.resetDaysSinceLastSeen(msg.SenderID())
	a.surroundingAgents[msg.SenderFloor()-a.Floor()] = msg.SenderID()
}
