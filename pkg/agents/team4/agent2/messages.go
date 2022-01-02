package agentTrust

import (
	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
)

func (a *CustomAgent4) SendingMessage(direction int) {

	var msg messages.Message
	floorToSend := a.Floor() + 1
	switch a.message_counter {
	case 0:
		msg = messages.NewAskFoodTakenMessage(a.ID(), a.Floor(), floorToSend)
	case 1:
		msg = messages.NewAskHPMessage(a.ID(), a.Floor(), floorToSend)
	case 2:
		msg = messages.NewAskIntendedFoodIntakeMessage(a.ID(), a.Floor(), floorToSend)
	case 3:
		msg = messages.NewRequestLeaveFoodMessage(a.ID(), a.Floor(), floorToSend, 10) //need to change how much to request to leave
	case 4:
		msg = messages.NewRequestTakeFoodMessage(a.ID(), a.Floor(), floorToSend, 20) //need to change how much to request to take
	case 5:
		floorToSend := a.Floor() - 1
		msg = messages.NewAskFoodTakenMessage(a.ID(), a.Floor(), floorToSend)
	case 6:
		floorToSend := a.Floor() - 1
		msg = messages.NewAskHPMessage(a.ID(), a.Floor(), floorToSend)
	case 7:
		floorToSend := a.Floor() - 1
		msg = messages.NewAskIntendedFoodIntakeMessage(a.ID(), a.Floor(), floorToSend)
	case 8:
		floorToSend := a.Floor() - 1
		msg = messages.NewRequestLeaveFoodMessage(a.ID(), a.Floor(), floorToSend, 10) //need to change how much to request to leave
	case 9:
		floorToSend := a.Floor() - 1
		msg = messages.NewRequestTakeFoodMessage(a.ID(), a.Floor(), floorToSend, 20) //need to change how much to request to take
	}
	a.SendMessage(msg)
	a.AppendToMessageMemory(direction, msg, a.sentMessages)
	a.Log("Team4 agent sent a message", infra.Fields{"message": msg.MessageType()})

	a.message_counter++
}

//FOR HONEST AGENTS
func (a *CustomAgent4) HandleAskHP(msg messages.AskHPMessage) {
	reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), a.HP())
	a.SendMessage(reply)
	a.Log("Team4 agent received an askHP message from ", infra.Fields{"floor": msg.SenderFloor(), "hp": a.HP()})
}

func (a *CustomAgent4) HandleAskFoodTaken(msg messages.AskFoodTakenMessage) {
	reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), int(a.lastFoodTaken))
	a.SendMessage(reply)
	a.Log("Team4 agent received an askFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor(), "food": a.lastFoodTaken})
}

func (a *CustomAgent4) HandleAskIntendedFoodTaken(msg messages.AskIntendedFoodIntakeMessage) {
	reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), int(a.IntendedFoodTaken))
	a.SendMessage(reply)
	a.Log("Team4 agent received an askIntendedFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor(), "food": a.IntendedFoodTaken})
}

func (a *CustomAgent4) HandleRequestLeaveFood(msg messages.RequestLeaveFoodMessage) {
	amount := msg.Request()
	response := true
	if a.PlatformOnFloor() && a.CurrPlatFood()-a.IntendedFoodTaken >= food.FoodType(amount) {
		response = true
	} else {
		response = false
	}
	reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), response) // TODO: Change for later dependent on circumstance
	a.SendMessage(reply)
	a.Log("Team4 agent received a requestLeaveFood message from ", infra.Fields{"floor": msg.SenderFloor()})
}

func (a *CustomAgent4) HandleRequestTakeFood(msg messages.RequestTakeFoodMessage) {
	amount := msg.Request()
	response := true
	if a.IntendedFoodTaken > food.FoodType(amount) {
		response = false
	}
	reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), response) // TODO: Change for later dependent on circumstance
	a.SendMessage(reply)
	a.Log("Team4 agent received a requestTakeFood message from ", infra.Fields{"floor": msg.SenderFloor()})
}

func (a *CustomAgent4) HandleResponse(msg messages.BoolResponseMessage) {
	response := msg.Response() // TODO: Change for later dependent on circumstance
	if !msg.Response() {
		a.globalTrust += a.globalTrustSubtract * a.coefficients[0] // TODO: adapt for other conditions
	} else {
		a.CheckForResponse(msg)
	}
	a.Log("Team4 agent received a Response message from ", infra.Fields{"floor": msg.SenderFloor(), "response": response})
}

func (a *CustomAgent4) HandleStateFoodTaken(msg messages.StateFoodTakenMessage) {
	statement := msg.Statement()
	if food.FoodType(statement) > a.maxFoodLimit {
		a.globalTrust += a.globalTrustSubtract * a.coefficients[3]
	} else {
		a.globalTrust += a.globalTrustAdd * a.coefficients[3]
	}
	a.Log("Team4 agent received a StateFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor(), "food": statement})
}

func (a *CustomAgent4) HandleStateHP(msg messages.StateHPMessage) {
	statement := msg.Statement()
	a.globalTrust += a.globalTrustAdd * a.coefficients[0]
	a.Log("Team4 agent received a StateHP message from ", infra.Fields{"floor": msg.SenderFloor(), "hp": statement})
}

func (a *CustomAgent4) HandleStateIntendedFoodTaken(msg messages.StateIntendedFoodIntakeMessage) {
	statement := msg.Statement()
	if food.FoodType(statement) > a.maxFoodLimit {
		a.globalTrust += a.globalTrustSubtract * a.coefficients[3]
	} else {
		a.globalTrust += a.globalTrustAdd * a.coefficients[3]
	}
	a.Log("Team4 agent received a StateIntendedFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor(), "food": statement})
}
