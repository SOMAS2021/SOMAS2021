package team4EvoAgent

import (
	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
)

func (a *CustomAgentEvo) GetMessage() {
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

func (a *CustomAgentEvo) CallHandleMessage() {
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

func (a *CustomAgentEvo) GenerateMessagesToSend() {
	var msg messages.Message

	// Request Agent above us to leave critical food amount when we have critical HP.
	if a.params.healthStatus == 0 {
		floorToSend := a.Floor() - 1
		msg = messages.NewRequestLeaveFoodMessage(a.ID(), a.Floor(), floorToSend, a.params.foodToEat[a.params.currentPersonality][a.params.healthStatus]) //TODO: need to change how much to request to leave
		a.params.msgToSendBuffer = append(a.params.msgToSendBuffer, msg)
	}

	// Ask HP and Ask Intended Food Intake to figure out if the person is taking excessive food
	msg = messages.NewAskHPMessage(a.ID(), a.Floor(), a.Floor()-1) // send to the agent above
	a.params.msgToSendBuffer = append(a.params.msgToSendBuffer, msg)
	msg = messages.NewAskIntendedFoodIntakeMessage(a.ID(), a.Floor(), a.Floor()-1) // send to the agent above
	a.params.msgToSendBuffer = append(a.params.msgToSendBuffer, msg)

	msg = messages.NewAskHPMessage(a.ID(), a.Floor(), a.Floor()+1) // send to the agent below
	a.params.msgToSendBuffer = append(a.params.msgToSendBuffer, msg)
	msg = messages.NewAskIntendedFoodIntakeMessage(a.ID(), a.Floor(), a.Floor()+1) // send to the agent below
	a.params.msgToSendBuffer = append(a.params.msgToSendBuffer, msg)
}

func (a *CustomAgentEvo) SendingMessage() {

	if len(a.params.msgToSendBuffer) > 0 {
		msg := a.params.msgToSendBuffer[0]
		a.SendMessage(msg)
		a.params.sentMessages = append(a.params.sentMessages, msg)
		a.Log("Team4 agent sent a message", infra.Fields{"message": msg.MessageType()})
		a.params.msgToSendBuffer = remove(a.params.msgToSendBuffer, 0)
	}

	// var msg messages.Message
	// floorToSend := a.Floor() + 1

	// switch a.params.messageCounter {
	// case 0:
	// 	msg = messages.NewAskFoodTakenMessage(a.ID(), a.Floor(), floorToSend)
	// 	a.SendMessage(msg)
	// 	a.params.sentMessages = append(a.params.sentMessages, msg)
	// 	a.Log("Team4 agent sent a message", infra.Fields{"message": msg.MessageType()})
	// case 1:
	// 	msg = messages.NewAskHPMessage(a.ID(), a.Floor(), floorToSend)
	// 	a.SendMessage(msg)
	// 	a.params.sentMessages = append(a.params.sentMessages, msg)
	// 	a.Log("Team4 agent sent a message", infra.Fields{"message": msg.MessageType()})
	// case 2:
	// 	msg = messages.NewAskIntendedFoodIntakeMessage(a.ID(), a.Floor(), floorToSend)
	// 	a.SendMessage(msg)
	// 	a.params.sentMessages = append(a.params.sentMessages, msg)
	// 	a.Log("Team4 agent sent a message", infra.Fields{"message": msg.MessageType()})
	// case 3:
	// 	msg = messages.NewRequestLeaveFoodMessage(a.ID(), a.Floor(), floorToSend, rand.Intn(60)) //TODO: need to change how much to request to leave
	// 	a.SendMessage(msg)
	// 	a.params.sentMessages = append(a.params.sentMessages, msg)
	// 	a.Log("Team4 agent sent a message", infra.Fields{"message": msg.MessageType()})
	// case 4:
	// 	msg = messages.NewRequestTakeFoodMessage(a.ID(), a.Floor(), floorToSend, int(int(a.CurrPlatFood())/a.Floor())) //need to change how much to request to take
	// 	a.SendMessage(msg)
	// 	a.params.sentMessages = append(a.params.sentMessages, msg)
	// 	a.Log("Team4 agent sent a message", infra.Fields{"message": msg.MessageType()})
	// case 5:
	// 	floorToSend = a.Floor() - 1
	// 	msg = messages.NewAskFoodTakenMessage(a.ID(), a.Floor(), floorToSend)
	// 	a.SendMessage(msg)
	// 	a.params.sentMessages = append(a.params.sentMessages, msg)
	// 	a.Log("Team4 agent sent a message", infra.Fields{"message": msg.MessageType()})
	// case 6:
	// 	floorToSend = a.Floor() - 1
	// 	msg = messages.NewAskHPMessage(a.ID(), a.Floor(), floorToSend)
	// 	a.SendMessage(msg)
	// 	a.params.sentMessages = append(a.params.sentMessages, msg)
	// 	a.Log("Team4 agent sent a message", infra.Fields{"message": msg.MessageType()})
	// case 7:
	// 	floorToSend = a.Floor() - 1
	// 	msg = messages.NewAskIntendedFoodIntakeMessage(a.ID(), a.Floor(), floorToSend)
	// 	a.SendMessage(msg)
	// 	a.params.sentMessages = append(a.params.sentMessages, msg)
	// 	a.Log("Team4 agent sent a message", infra.Fields{"message": msg.MessageType()})
	// case 8:
	// 	floorToSend = a.Floor() - 1
	// 	msg = messages.NewRequestLeaveFoodMessage(a.ID(), a.Floor(), floorToSend, a.params.foodToEat[a.params.currentPersonality][a.params.healthStatus]) //need to change how much to request to leave
	// 	a.SendMessage(msg)
	// 	a.params.sentMessages = append(a.params.sentMessages, msg)
	// 	a.Log("Team4 agent sent a message", infra.Fields{"message": msg.MessageType()})
	// case 9:
	// 	floorToSend = a.Floor() - 1
	// 	msg = messages.NewRequestTakeFoodMessage(a.ID(), a.Floor(), floorToSend, 60) //need to change how much to request to take
	// 	a.SendMessage(msg)
	// 	a.params.sentMessages = append(a.params.sentMessages, msg)
	// 	a.Log("Team4 agent sent a message", infra.Fields{"message": msg.MessageType()})
	// default:
	// }

	// if a.HasDayPassed() {
	// 	a.params.messageCounter = 0
	// }
	// a.params.messageCounter++
}

//FOR HONEST AGENTS
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
	//amount on platform - intended amount to take  >= request then respond true
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
		a.AddToGlobalTrust(-a.params.coefficients[1])
	} else {
		a.CheckForResponse(msg)
	}
	a.Log("Team4 agent received a Response message from ", infra.Fields{"floor": msg.SenderFloor(), "response": response, "global_trust": a.params.globalTrust})
}

func (a *CustomAgentEvo) HandleStateFoodTaken(msg messages.StateFoodTakenMessage) {
	statement := msg.Statement()
	if food.FoodType(statement) > a.params.maxFoodLimit {
		a.AddToGlobalTrust(-a.params.coefficients[1])
	} else {
		a.AddToGlobalTrust(a.params.coefficients[1])
	}
	a.Log("Team4 agent received a StateFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor(), "food": statement, "global_trust": a.params.globalTrust})
}

func (a *CustomAgentEvo) HandleStateHP(msg messages.StateHPMessage) {
	statement := msg.Statement()

	_, exists := a.params.agentResponseMemory[msg.SenderID()]
	if !exists {
		a.params.agentResponseMemory[msg.SenderID()] = []int{-1, -1}
	}
	a.params.agentResponseMemory[msg.SenderID()][0] = statement

	//a.AddToGlobalTrust(a.params.coefficients[0])
	a.Log("Team4 agent received a StateHP message from ", infra.Fields{"floor": msg.SenderFloor(), "hp": statement, "global_trust": a.params.globalTrust})
}

func (a *CustomAgentEvo) HandleStateIntendedFoodTaken(msg messages.StateIntendedFoodIntakeMessage) {
	statement := msg.Statement()

	_, exists := a.params.agentResponseMemory[msg.SenderID()]
	if !exists {
		a.params.agentResponseMemory[msg.SenderID()] = []int{-1, -1}
	}
	a.params.agentResponseMemory[msg.SenderID()][1] = statement

	// if food.FoodType(statement) > a.params.maxFoodLimit {
	// 	a.AddToGlobalTrust(-a.params.coefficients[1])
	// } else {
	// 	a.AddToGlobalTrust(a.params.coefficients[1])
	// }
	a.Log("Team4 agent received a StateIntendedFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor(), "food": statement, "global_trust": a.params.globalTrust})
}
