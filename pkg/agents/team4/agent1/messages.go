package team4EvoAgent

import (
	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
)

func (a *CustomAgentEvo) SendingMessage() {

	var msg messages.Message
	floorToSend := a.Floor() + 1
	switch a.messageCounter {
	case 0:
		msg = messages.NewAskFoodTakenMessage(a.ID(), a.Floor(), floorToSend)
		a.SendMessage(msg)
		a.AppendToMessageMemory(msg, a.sentMessages)
		a.Log("Team4 agent sent a message", infra.Fields{"message": msg.MessageType()})
	case 1:
		msg = messages.NewAskHPMessage(a.ID(), a.Floor(), floorToSend)
		a.SendMessage(msg)
		a.AppendToMessageMemory(msg, a.sentMessages)
		a.Log("Team4 agent sent a message", infra.Fields{"message": msg.MessageType()})
	case 2:
		msg = messages.NewAskIntendedFoodIntakeMessage(a.ID(), a.Floor(), floorToSend)
		a.SendMessage(msg)
		a.AppendToMessageMemory(msg, a.sentMessages)
		a.Log("Team4 agent sent a message", infra.Fields{"message": msg.MessageType()})
	case 3:
		msg = messages.NewRequestLeaveFoodMessage(a.ID(), a.Floor(), floorToSend, 10) //need to change how much to request to leave
		a.SendMessage(msg)
		a.AppendToMessageMemory(msg, a.sentMessages)
		a.Log("Team4 agent sent a message", infra.Fields{"message": msg.MessageType()})
	case 4:
		msg = messages.NewRequestTakeFoodMessage(a.ID(), a.Floor(), floorToSend, 20) //need to change how much to request to take
		a.SendMessage(msg)
		a.AppendToMessageMemory(msg, a.sentMessages)
		a.Log("Team4 agent sent a message", infra.Fields{"message": msg.MessageType()})
	case 5:
		floorToSend = a.Floor() - 1
		msg = messages.NewAskFoodTakenMessage(a.ID(), a.Floor(), floorToSend)
		a.SendMessage(msg)
		a.AppendToMessageMemory(msg, a.sentMessages)
		a.Log("Team4 agent sent a message", infra.Fields{"message": msg.MessageType()})
	case 6:
		floorToSend = a.Floor() - 1
		msg = messages.NewAskHPMessage(a.ID(), a.Floor(), floorToSend)
		a.SendMessage(msg)
		a.AppendToMessageMemory(msg, a.sentMessages)
		a.Log("Team4 agent sent a message", infra.Fields{"message": msg.MessageType()})
	case 7:
		floorToSend = a.Floor() - 1
		msg = messages.NewAskIntendedFoodIntakeMessage(a.ID(), a.Floor(), floorToSend)
		a.SendMessage(msg)
		a.AppendToMessageMemory(msg, a.sentMessages)
		a.Log("Team4 agent sent a message", infra.Fields{"message": msg.MessageType()})
	case 8:
		floorToSend = a.Floor() - 1
		msg = messages.NewRequestLeaveFoodMessage(a.ID(), a.Floor(), floorToSend, 10) //need to change how much to request to leave
		a.SendMessage(msg)
		a.AppendToMessageMemory(msg, a.sentMessages)
		a.Log("Team4 agent sent a message", infra.Fields{"message": msg.MessageType()})
	case 9:
		floorToSend = a.Floor() - 1
		msg = messages.NewRequestTakeFoodMessage(a.ID(), a.Floor(), floorToSend, 20) //need to change how much to request to take
		a.SendMessage(msg)
		a.AppendToMessageMemory(msg, a.sentMessages)
		a.Log("Team4 agent sent a message", infra.Fields{"message": msg.MessageType()})
	default:
	}

	if a.HasDayPassed() {
		a.messageCounter = 0
	}
	a.messageCounter++
}

//FOR HONEST AGENTS
func (a *CustomAgentEvo) HandleAskHP(msg messages.AskHPMessage) {
	reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), a.HP())
	a.SendMessage(reply)
	a.Log("Team4 agent received an askHP message from ", infra.Fields{"floor": msg.SenderFloor(), "hp": a.HP()})
}

func (a *CustomAgentEvo) HandleAskFoodTaken(msg messages.AskFoodTakenMessage) {
	reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), int(a.lastFoodTaken))
	a.SendMessage(reply)
	a.Log("Team4 agent received an askFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor(), "food": a.lastFoodTaken})
}

func (a *CustomAgentEvo) HandleAskIntendedFoodTaken(msg messages.AskIntendedFoodIntakeMessage) {
	reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), int(a.intendedFoodTaken))
	a.SendMessage(reply)
	a.Log("Team4 agent received an askIntendedFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor(), "food": a.intendedFoodTaken})
}

func (a *CustomAgentEvo) HandleRequestLeaveFood(msg messages.RequestLeaveFoodMessage) {
	response := true
	if a.globalTrust < a.globalTrustLimit {
		response = false
	}
	reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), response) // TODO: Change for later dependent on circumstance
	a.SendMessage(reply)
	a.Log("Team4 agent received a requestLeaveFood message from ", infra.Fields{"floor": msg.SenderFloor()})
}

func (a *CustomAgentEvo) HandleRequestTakeFood(msg messages.RequestTakeFoodMessage) {
	amount := msg.Request()
	response := true
	if a.intendedFoodTaken > food.FoodType(amount) {
		response = false
	}
	reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), response) // TODO: Change for later dependent on circumstance
	a.SendMessage(reply)
	a.Log("Team4 agent received a requestTakeFood message from ", infra.Fields{"floor": msg.SenderFloor()})
}

func (a *CustomAgentEvo) HandleResponse(msg messages.BoolResponseMessage) {
	response := msg.Response() // TODO: Change for later dependent on circumstance
	if !msg.Response() {
		a.SubFromGlobalTrust(a.coefficients[1]) // TODO: adapt for other conditions
	} else {
		a.CheckForResponse(msg)
	}
	a.Log("Team4 agent received a Response message from ", infra.Fields{"floor": msg.SenderFloor(), "response": response, "global_trust": a.globalTrust})

}

func (a *CustomAgentEvo) HandleStateFoodTaken(msg messages.StateFoodTakenMessage) {
	statement := msg.Statement()
	if food.FoodType(statement) > a.maxFoodLimit {
		a.SubFromGlobalTrust(a.coefficients[1])
	} else {
		a.AddToGlobalTrust(a.coefficients[1])
	}
	a.Log("Team4 agent received a StateFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor(), "food": statement, "global_trust": a.globalTrust})
}

func (a *CustomAgentEvo) HandleStateHP(msg messages.StateHPMessage) {
	statement := msg.Statement()
	a.AddToGlobalTrust(a.coefficients[0])
	a.Log("Team4 agent received a StateHP message from ", infra.Fields{"floor": msg.SenderFloor(), "hp": statement, "global_trust": a.globalTrust})
}

func (a *CustomAgentEvo) HandleStateIntendedFoodTaken(msg messages.StateIntendedFoodIntakeMessage) {
	statement := msg.Statement()
	if food.FoodType(statement) > a.maxFoodLimit {
		a.SubFromGlobalTrust(a.coefficients[1])
	} else {
		a.AddToGlobalTrust(a.coefficients[1])
	}
	a.Log("Team4 agent received a StateIntendedFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor(), "food": statement, "global_trust": a.globalTrust})
}
