package team5

import (
	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
)

func (a *CustomAgent5) getMessages() {
	receivedMsg := a.ReceiveMessage()
	if receivedMsg != nil {
		receivedMsg.Visit(a)
	}
}

func (a *CustomAgent5) dailyMessages() {
	var msg messages.Message
	targetFloor := a.Floor() + 1
	switch a.messagingCounter {
	case 0:
		msg = messages.NewAskHPMessage(a.ID(), a.Floor(), targetFloor)
		a.SendMessage(msg)
		a.Log("Team 5 agent sent message", infra.Fields{"msgType": msg.MessageType().String(), "senderFloor": a.Floor(), "targetFloor": targetFloor})
	case 1:
		msg = messages.NewAskFoodTakenMessage(a.ID(), a.Floor(), targetFloor)
		a.SendMessage(msg)
		a.Log("Team 5 agent sent message", infra.Fields{"msgType": msg.MessageType().String(), "senderFloor": a.Floor(), "targetFloor": targetFloor})
	case 2:
		msg = messages.NewAskIntendedFoodIntakeMessage(a.ID(), a.Floor(), targetFloor)
		a.SendMessage(msg)
		a.Log("Team 5 agent sent message", infra.Fields{"msgType": msg.MessageType().String(), "senderFloor": a.Floor(), "targetFloor": targetFloor})
	case 3:
		targetFloor = a.Floor() - 1
		msg = messages.NewAskHPMessage(a.ID(), a.Floor(), targetFloor)
		a.SendMessage(msg)
		a.Log("Team 5 agent sent message", infra.Fields{"msgType": msg.MessageType().String(), "senderFloor": a.Floor(), "targetFloor": targetFloor})
	case 4:
		targetFloor = a.Floor() - 1
		msg = messages.NewAskFoodTakenMessage(a.ID(), a.Floor(), targetFloor)
		a.SendMessage(msg)
		a.Log("Team 5 agent sent message", infra.Fields{"msgType": msg.MessageType().String(), "senderFloor": a.Floor(), "targetFloor": targetFloor})
	case 5:
		targetFloor = a.Floor() - 1
		msg = messages.NewAskIntendedFoodIntakeMessage(a.ID(), a.Floor(), targetFloor)
		a.SendMessage(msg)
		a.Log("Team 5 agent sent message", infra.Fields{"msgType": msg.MessageType().String(), "senderFloor": a.Floor(), "targetFloor": targetFloor})
	default:
	}
	a.messagingCounter++
}

//Override baseAgent's HandlePropogate, and process any StateMessages

func (a *CustomAgent5) HandlePropogate(msg messages.Message) {
	if stateFoodTakenMsg, ok := msg.(*messages.StateFoodTakenMessage); ok {
		a.HandleStateFoodTaken(*stateFoodTakenMsg)
	} else if stateHPMsg, ok := msg.(*messages.StateHPMessage); ok {
		a.HandleStateHP(*stateHPMsg)
	} else if stateIntendedFoodIntakeMsg, ok := msg.(*messages.StateIntendedFoodIntakeMessage); ok {
		a.HandleStateIntendedFoodTaken(*stateIntendedFoodIntakeMsg)
	}

	a.SendMessage(msg)
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
