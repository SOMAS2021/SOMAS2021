package team6

import (
	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
)

// Request another agent to leave food on the platform
func (a *CustomAgent6) RequestLeaveFood() {
	healthInfo := a.HealthInfo()
	currentHP := a.HP()

	levels := levelsData{
		strongLevel:  healthInfo.MaxHP * 3 / 5,
		healthyLevel: healthInfo.MaxHP * 3 / 10,
		weakLevel:    healthInfo.MaxHP * 1 / 10,
	}

	reqAmount := -1
	switch a.currBehaviour.String() {
	case "Altruist":
		reqAmount = -1

	case "Collectivist":
		if currentHP >= levels.weakLevel {
			reqAmount = -1
		} else {
			reqAmount = int(FoodRequired(currentHP, a.HealthInfo().HPCritical+a.HealthInfo().HPReqCToW, a.HealthInfo()))
		}

	case "Selfish":
		if currentHP >= levels.strongLevel {
			reqAmount = -1
		} else {
			reqAmount = int(FoodRequired(currentHP, levels.healthyLevel, a.HealthInfo()))
		}

	case "Narcissist":
		reqAmount = 4 * int(a.HealthInfo().Tau)

	default:
		reqAmount = -1
	}

	if reqAmount != -1 {
		msg := messages.NewRequestLeaveFoodMessage(a.ID(), a.Floor(), a.Floor()-1, reqAmount)
		a.SendMessage(msg)
		a.Log("I sent a message", infra.Fields{"message": "RequestLeaveFood", "floor": a.Floor()})
	}
}

// Request another agent to take a precise amount of food
func (a *CustomAgent6) RequestTakeFood() {

	reqAmount := -1
	switch a.currBehaviour.String() {
	case "Altruist":
		reqAmount = -1

	case "Collectivist":
		reqAmount = -1

	case "Selfish":
		reqAmount = 0

	case "Narcissist":
		reqAmount = 0

	default:
		reqAmount = -1
	}

	if reqAmount != -1 {
		msg := messages.NewRequestTakeFoodMessage(a.ID(), a.Floor(), a.Floor()-1, reqAmount)
		a.SendMessage(msg)
		a.Log("I sent a message", infra.Fields{"message": "RequestLeaveFood", "floor": a.Floor()})
	}
}

func (a *CustomAgent6) HandleRequestLeaveFood(msg messages.RequestLeaveFoodMessage) {
	healthInfo := a.HealthInfo()
	currentHP := a.HP()

	levels := levelsData{
		strongLevel:  healthInfo.MaxHP * 3 / 5,
		healthyLevel: healthInfo.MaxHP * 3 / 10,
		weakLevel:    healthInfo.MaxHP * 1 / 10,
	}

	reply := true
	switch a.currBehaviour.String() {
	case "Altruist":
		reply = true

	case "Collectivist":
		if currentHP >= levels.weakLevel {
			reply = true
		} else {
			reply = false
		}

	case "Selfish":
		if currentHP >= levels.strongLevel {
			reply = true
		} else {
			reply = false
		}

	case "Narcissist":
		reply = false

	default:
		reply = true
	}

	replyMessage := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), reply)
	a.SendMessage(replyMessage)

	if reply == true {
		a.reqLeaveFoodAmount = msg.Request()
		a.Log("I received a requestLeaveFood message and my response was true")
	} else {
		a.reqLeaveFoodAmount = -1
		a.Log("I received a requestLeaveFood message and my response was false")
	}

}

func (a *CustomAgent6) HandleRequestTakeFood(msg messages.RequestTakeFoodMessage) {
	reply := false
	replyMessage := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), reply)
	a.SendMessage(replyMessage)

	if reply == true {
		a.reqLeaveFoodAmount = msg.Request()
		a.Log("I received a requestTakeFood message and my response was true")
	} else {
		a.reqLeaveFoodAmount = -1
		a.Log("I received a requestTakeFood message and my response was false")
	}
}

func (a *CustomAgent6) HandleAskHP(msg messages.AskHPMessage) {
	if a.currBehaviour.String() != "Narcissist" {
		reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), a.HP())
		a.SendMessage(reply)
		a.Log("I recieved an askHP message from ", infra.Fields{"senderFloor": msg.SenderFloor(), "myFloor": a.Floor()})
	}
}

func (a *CustomAgent6) HandleAskFoodTaken(msg messages.AskFoodTakenMessage) {
	reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), int(a.lastFoodTaken))
	a.SendMessage(reply)
	a.Log("I recieved an askFoodTaken message from ", infra.Fields{"senderFloor": msg.SenderFloor(), "myFloor": a.Floor()})
}

func (a *CustomAgent6) HandleAskIntendedFoodTaken(msg messages.AskIntendedFoodIntakeMessage) {
	reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), int(a.intendedFoodIntake()))
	a.SendMessage(reply)
	a.Log("I recieved an askIntendedFoodTaken message from ", infra.Fields{"senderFloor": msg.SenderFloor(), "myFloor": a.Floor()})
}
