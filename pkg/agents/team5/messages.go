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
	case 6:
		targetFloor = a.Floor() - 2
		msg = messages.NewAskHPMessage(a.ID(), a.Floor(), targetFloor)
	case 7:
		targetFloor = a.Floor() - 2
		msg = messages.NewAskFoodTakenMessage(a.ID(), a.Floor(), targetFloor)
	case 8:
		targetFloor = a.Floor() - 2
		msg = messages.NewAskIntendedFoodIntakeMessage(a.ID(), a.Floor(), targetFloor)
	case 9:
		targetFloor = a.Floor() + 2
		msg = messages.NewAskHPMessage(a.ID(), a.Floor(), targetFloor)
	case 10:
		targetFloor = a.Floor() + 2
		msg = messages.NewAskFoodTakenMessage(a.ID(), a.Floor(), targetFloor)
	case 11:
		targetFloor = a.Floor() + 2
		msg = messages.NewAskIntendedFoodIntakeMessage(a.ID(), a.Floor(), targetFloor)
	default:
	}
	if msg != nil {
		a.SendMessage(msg)
		a.Log("Team 5 agent sent message", infra.Fields{"msgType": msg.MessageType().String(), "senderFloor": a.Floor(), "targetFloor": targetFloor})
	}
	a.messagingCounter++
	// This needs to be at least twice the number of cases in the switch statement.
	if a.messagingCounter > 25 {
		a.messagingCounter = 0
	}
}

//Override baseAgent's HandlePropogate, and process any StateMessages
func (a *CustomAgent5) HandlePropogate(msg messages.Message) {
	switch stateMsg := msg.(type) {
	case *messages.StateFoodTakenMessage:
		a.HandleStateFoodTaken(*stateMsg)
	case *messages.StateHPMessage:
		a.HandleStateHP(*stateMsg)
	case *messages.StateIntendedFoodIntakeMessage:
		a.HandleStateIntendedFoodTaken(*stateMsg)
	}
	a.SendMessage(msg)
}

func (a *CustomAgent5) postHandleAskMessage(msg *messages.BaseMessage, reply messages.StateMessage) {
	a.SendMessage(reply)
	a.Log("Team 5 agent received a message", infra.Fields{"senderFloor": msg.SenderFloor(), "messageType": msg.MessageType().String()})
	a.updateSocialMemory(msg.SenderID(), msg.SenderFloor())
}

func (a *CustomAgent5) postHandleRequestMessage(msg *messages.BaseMessage, reply messages.ResponseMessage, amount int) {
	a.SendMessage(reply)
	a.Log("Team 5 agent received a message", infra.Fields{"senderFloor": msg.SenderFloor(), "messageType": msg.MessageType().String(), "requestAmount": amount})
	a.updateSocialMemory(msg.SenderID(), msg.SenderFloor())
}

func (a *CustomAgent5) handleStateMessage(msg *messages.BaseMessage, statement int) {
	a.Log("Team 5 agent received a message", infra.Fields{"senderFloor": msg.SenderFloor(), "messageType": msg.MessageType().String(), "statement": statement})
	a.updateSocialMemory(msg.SenderID(), msg.SenderFloor())
}

//The message handler functions below are for a fully honest agent

func (a *CustomAgent5) HandleAskHP(msg messages.AskHPMessage) {
	reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), a.HP())
	a.postHandleAskMessage(msg.BaseMessage, reply)
}

func (a *CustomAgent5) HandleAskFoodTaken(msg messages.AskFoodTakenMessage) {
	reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), int(a.lastMeal))
	a.postHandleAskMessage(msg.BaseMessage, reply)
}

func (a *CustomAgent5) HandleAskIntendedFoodTaken(msg messages.AskIntendedFoodIntakeMessage) {
	amount := int(a.calculateAttemptFood())
	reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), amount)
	a.postHandleAskMessage(msg.BaseMessage, reply)
}

func (a *CustomAgent5) HandleRequestLeaveFood(msg messages.RequestLeaveFoodMessage) {
	// Always set to false for now to prevent deception, needs some calculations to determine whether we will leave the requested amount
	reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), false)
	a.postHandleRequestMessage(msg.BaseMessage, reply, msg.Request())
}

func (a *CustomAgent5) HandleRequestTakeFood(msg messages.RequestTakeFoodMessage) {
	response := a.calculateAttemptFood() <= food.FoodType(msg.Request())
	reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), response)
	a.postHandleRequestMessage(msg.BaseMessage, reply, msg.Request())
}

func (a *CustomAgent5) HandleResponse(msg messages.BoolResponseMessage) {
	response := msg.Response()
	a.Log("Team 5 agent received a message", infra.Fields{"senderFloor": msg.SenderFloor(), "messageType": msg.MessageType().String(), "response": response})
	a.updateSocialMemory(msg.SenderID(), msg.SenderFloor())
}

func (a *CustomAgent5) HandleStateFoodTaken(msg messages.StateFoodTakenMessage) {
	statement := food.FoodType(msg.Statement())
	a.handleStateMessage(msg.BaseMessage, msg.Statement())
	a.updateFoodTakenMemory(msg.SenderID(), statement)
	//a.Log("New value of foodTaken", infra.Fields{"statement": statement, "memory": a.socialMemory[msg.SenderID()].foodTaken})
}

func (a *CustomAgent5) HandleStateHP(msg messages.StateHPMessage) {
	statement := msg.Statement()
	a.handleStateMessage(msg.BaseMessage, statement)
	a.updateAgentHPMemory(msg.SenderID(), statement)
	//a.Log("New value of agentHP", infra.Fields{"statement": statement, "memory": a.socialMemory[msg.SenderID()].agentHP})
}

func (a *CustomAgent5) HandleStateIntendedFoodTaken(msg messages.StateIntendedFoodIntakeMessage) {
	statement := food.FoodType(msg.Statement())
	a.handleStateMessage(msg.BaseMessage, msg.Statement())
	a.updateIntentionFoodMemory(msg.SenderID(), statement)
	//a.Log("New value of intendedFood", infra.Fields{"statement": statement, "memory": a.socialMemory[msg.SenderID()].intentionFood})
}
