package team4EvoAgent

import (
	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
)

// ---------- RECEIVE MESSAGE ---------//

func (a *CustomAgentEvo) getMessage() {
	receivedMsg := a.ReceiveMessage()

	if receivedMsg != nil {
		if receivedMsg.MessageType() == messages.RequestLeaveFood {
			a.params.requestLeaveFoodMessages = append(a.params.requestLeaveFoodMessages, receivedMsg)
		} else {
			a.params.otherMessageBuffer = append(a.params.otherMessageBuffer, receivedMsg)
		}
		a.params.receivedMessagesCount++
	}
}

// ----------- CALL THE RELEVANT MESSAGE HANDLER ---------//

func (a *CustomAgentEvo) callHandleMessage() {
	if a.PlatformOnFloor() && len(a.params.requestLeaveFoodMessages) > 0 {
		a.params.requestLeaveFoodMessages[0].Visit(a)
		a.params.requestLeaveFoodMessages = remove(a.params.requestLeaveFoodMessages, 0)
	} else if len(a.params.otherMessageBuffer) > 0 {
		a.params.otherMessageBuffer[0].Visit(a)
		a.params.otherMessageBuffer = remove(a.params.otherMessageBuffer, 0)
	} else {
		a.Log("I got no messages")
	}
}

/*------------------------HANDLING MESSAGES TO BE SENT ------------------------*/

func (a *CustomAgentEvo) generateMessagesToSend() {
	var msg messages.Message

	// Request Agent above us to leave critical food amount when we have critical HP.
	if a.params.healthStatus == 0 {
		floorToSend := a.Floor() - 1
		msg = messages.NewRequestLeaveFoodMessage(a.ID(), a.Floor(), floorToSend, a.params.foodToEat[a.params.currentPersonality][a.params.healthStatus])
		a.params.msgToSendBuffer = append(a.params.msgToSendBuffer, msg)
	}

	directions := []int{-3, -2, -1, 1, 2, 3}
	// Ask HP and Ask Intended Food Intake to figure out if the person is taking excessive food
	for _, dir := range directions {
		if dir+a.Floor() >= 1 {
			msg = messages.NewAskHPMessage(a.ID(), a.Floor(), a.Floor()+dir) // send to the agent above
			a.params.msgToSendBuffer = append(a.params.msgToSendBuffer, msg)
			msg = messages.NewAskIntendedFoodIntakeMessage(a.ID(), a.Floor(), a.Floor()+dir) // send to the agent above
			a.params.msgToSendBuffer = append(a.params.msgToSendBuffer, msg)
		}
	}
}

func (a *CustomAgentEvo) sendingMessage() {

	if len(a.params.msgToSendBuffer) > 0 {
		msg := a.params.msgToSendBuffer[0]
		a.SendMessage(msg)
		a.params.sentMessages = append(a.params.sentMessages, msg)
		a.Log("Team4 agent sent a message", infra.Fields{"message": msg.MessageType()})
		a.params.msgToSendBuffer = remove(a.params.msgToSendBuffer, 0)
	}
}

/*------------------------HANDLING RECEIVED MESSAGES ------------------------*/

func (a *CustomAgentEvo) HandleAskHP(msg messages.AskHPMessage) {
	reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), a.HP())
	a.SendMessage(reply)
	a.Log("Team4 agent received an askHP message from ", infra.Fields{"floor": msg.SenderFloor(), "hp": a.HP()})
}

func (a *CustomAgentEvo) HandleAskFoodTaken(msg messages.AskFoodTakenMessage) {
	reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), int(a.params.lastFoodTaken))
	a.SendMessage(reply)
	a.Log("Team4 agent received an askFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor(), "food": a.params.lastFoodTaken})
}

func (a *CustomAgentEvo) HandleAskIntendedFoodTaken(msg messages.AskIntendedFoodIntakeMessage) {
	reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), int(a.params.intendedFoodToTake))
	a.SendMessage(reply)
	a.Log("Team4 agent received an askIntendedFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor(), "food": a.params.intendedFoodToTake})
}

func (a *CustomAgentEvo) HandleRequestLeaveFood(msg messages.RequestLeaveFoodMessage) {
	amount := msg.Request()
	response := false
	if a.CurrPlatFood()-a.params.intendedFoodToTake >= food.FoodType(amount) {
		response = true
	}
	reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), response)
	a.SendMessage(reply)
	a.Log("Team4 agent received a requestLeaveFood message from ", infra.Fields{"floor": msg.SenderFloor()})
}

func (a *CustomAgentEvo) HandleRequestTakeFood(msg messages.RequestTakeFoodMessage) {
	amount := msg.Request()
	response := a.params.intendedFoodToTake <= food.FoodType(amount)
	reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), response)
	a.SendMessage(reply)
	a.Log("Team4 agent received a requestTakeFood message from ", infra.Fields{"floor": msg.SenderFloor()})
}

func (a *CustomAgentEvo) HandleResponse(msg messages.BoolResponseMessage) {
	response := msg.Response()
	if !msg.Response() {
		a.addToGlobalTrust(-a.params.trustCoefficients[1])
	} else {
		a.checkForResponse(msg)
	}
	a.Log("Team4 agent received a Response message from ", infra.Fields{"floor": msg.SenderFloor(), "response": response, "global_trust": a.params.globalTrust})
}

func (a *CustomAgentEvo) HandleStateFoodTaken(msg messages.StateFoodTakenMessage) {
	statement := msg.Statement()
	if food.FoodType(statement) > food.FoodType(a.params.foodToEat["selfish"][2]) {
		a.addToGlobalTrust(-a.params.trustCoefficients[1])
	} else {
		a.addToGlobalTrust(a.params.trustCoefficients[1])
	}
	a.Log("Team4 agent received a StateFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor(), "food": statement, "global_trust": a.params.globalTrust})
}

func (a *CustomAgentEvo) HandleStateHP(msg messages.StateHPMessage) {
	statement := msg.Statement()

	_, exists := a.params.agentResponseMemory[msg.SenderID()]
	if !exists {
		a.params.agentResponseMemory[msg.SenderID()] = []int{-1, -1} // initialise the map if sender not in map
	}
	a.params.agentResponseMemory[msg.SenderID()][0] = statement

	a.Log("Team4 agent received a StateHP message from ", infra.Fields{"floor": msg.SenderFloor(), "hp": statement, "global_trust": a.params.globalTrust})
}

func (a *CustomAgentEvo) HandleStateIntendedFoodTaken(msg messages.StateIntendedFoodIntakeMessage) {
	statement := msg.Statement()

	_, exists := a.params.agentResponseMemory[msg.SenderID()]
	if !exists {
		a.params.agentResponseMemory[msg.SenderID()] = []int{-1, -1} // initialise the map if sender not in map
	}
	a.params.agentResponseMemory[msg.SenderID()][1] = statement

	a.Log("Team4 agent received a StateIntendedFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor(), "food": statement, "global_trust": a.params.globalTrust})
}
