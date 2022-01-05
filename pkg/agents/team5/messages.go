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

func (a *CustomAgent5) treatyProposal() {
	sendToFloor := a.Floor()
	switch a.treatySendCounter {
	case 1:
		sendToFloor = a.Floor() + 1
	case 2:
		sendToFloor = a.Floor() - 1
	case 3:
		sendToFloor = a.Floor() + 2
	case 4:
		sendToFloor = a.Floor() - 2
	case 5:
		sendToFloor = a.Floor() + 3
	case 6:
		sendToFloor = a.Floor() - 3
	}

	msg := messages.NewProposalMessage(a.ID(), a.Floor(), sendToFloor, *a.currentProposal)
	a.SendMessage(msg)
	a.treatySendCounter++
	if a.treatySendCounter > 6 {
		a.treatySendCounter = 0
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
	case 6:
		targetFloor = a.Floor() - 2
		msg = messages.NewAskHPMessage(a.ID(), a.Floor(), targetFloor)
		a.SendMessage(msg)
		a.Log("Team 5 agent sent message", infra.Fields{"msgType": msg.MessageType().String(), "senderFloor": a.Floor(), "targetFloor": targetFloor})
	case 7:
		targetFloor = a.Floor() - 2
		msg = messages.NewAskFoodTakenMessage(a.ID(), a.Floor(), targetFloor)
		a.SendMessage(msg)
		a.Log("Team 5 agent sent message", infra.Fields{"msgType": msg.MessageType().String(), "senderFloor": a.Floor(), "targetFloor": targetFloor})
	case 8:
		targetFloor = a.Floor() - 2
		msg = messages.NewAskIntendedFoodIntakeMessage(a.ID(), a.Floor(), targetFloor)
		a.SendMessage(msg)
		a.Log("Team 5 agent sent message", infra.Fields{"msgType": msg.MessageType().String(), "senderFloor": a.Floor(), "targetFloor": targetFloor})
	case 9:
		targetFloor = a.Floor() + 2
		msg = messages.NewAskHPMessage(a.ID(), a.Floor(), targetFloor)
		a.SendMessage(msg)
		a.Log("Team 5 agent sent message", infra.Fields{"msgType": msg.MessageType().String(), "senderFloor": a.Floor(), "targetFloor": targetFloor})
	case 10:
		targetFloor = a.Floor() + 2
		msg = messages.NewAskFoodTakenMessage(a.ID(), a.Floor(), targetFloor)
		a.SendMessage(msg)
		a.Log("Team 5 agent sent message", infra.Fields{"msgType": msg.MessageType().String(), "senderFloor": a.Floor(), "targetFloor": targetFloor})
	case 11:
		targetFloor = a.Floor() + 2
		msg = messages.NewAskIntendedFoodIntakeMessage(a.ID(), a.Floor(), targetFloor)
		a.SendMessage(msg)
		a.Log("Team 5 agent sent message", infra.Fields{"msgType": msg.MessageType().String(), "senderFloor": a.Floor(), "targetFloor": targetFloor})
	default:
	}
	a.messagingCounter++
	if a.messagingCounter > 25 {
		a.messagingCounter = 0
	}
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
	//a.Log("New value of foodTaken", infra.Fields{"statement": statement, "memory": a.socialMemory[msg.SenderID()].foodTaken})
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
	//a.Log("New value of agentHP", infra.Fields{"statement": statement, "memory": a.socialMemory[msg.SenderID()].agentHP})
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
	//a.Log("New value of intendedFood", infra.Fields{"statement": statement, "memory": a.socialMemory[msg.SenderID()].intentionFood})
}

func (a *CustomAgent5) RejectTreaty(msg messages.ProposeTreatyMessage) {
	reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), false)
	a.SendMessage(reply)
	a.Log("Rejected treaty", infra.Fields{"proposerID": msg.SenderID(), "proposerFloor": msg.SenderFloor(), "treatyID": msg.TreatyID()})
}

func (a *CustomAgent5) HandleProposeTreaty(msg messages.ProposeTreatyMessage) {
	treaty := msg.Treaty()
	if treaty.Request() == messages.Inform {
		// Reject any inform treaties as I'm not sure what they are or how to handle them
		a.RejectTreaty(msg)
		return
	}

	if a.treatyConflicts(treaty) {
		// Reject treaty if it conflicts with any other active treaty
		a.RejectTreaty(msg)
		return
	}

	if (treaty.Request() == messages.LeaveAmountFood && treaty.RequestOp() == messages.EQ) || treaty.RequestOp() == messages.LE || treaty.RequestOp() == messages.LT {
		//Reject all treaties that ask you to leave less food, don't see why you would do this
		//Reject any treaty that asks you to leave a specific amount of food as this would lead to multiple people just not eating as only 1 agent can eat and fufill that criteria
		a.RejectTreaty(msg)
		return
	}

	if (treaty.ConditionOp() == messages.LE || treaty.ConditionOp() == messages.LT) && (treaty.Condition() == messages.HP || treaty.Condition() == messages.AvailableFood) {
		//Reject treaties that have a condition of being less than something
		// Why agree to treaty which is bounded by you having lower HP
		// Why agree to treaty which is bounded by the amount of food on platform being low
		// All of these are conditions where you will be desperate for food so shouldn't be limitting yourself if you do get the opurrtunity to eat
		a.RejectTreaty(msg)
		return
	}

	if (treaty.ConditionOp() == messages.GE || treaty.ConditionOp() == messages.GT) && treaty.Condition() == messages.Floor {
		// Why agree to treaty which is bounded by you being lower in tower
		a.RejectTreaty(msg)
		return
	}

	if treaty.Condition() == messages.HP && treaty.ConditionValue() < a.HealthInfo().WeakLevel*2 {
		//Reject any HP condition based treaty if the condition is too strict, and you would be put into critical by not eating
		a.RejectTreaty(msg)
		return
	}

	if treaty.Condition() == messages.Floor && treaty.ConditionValue() != 1 && treaty.Duration() >= a.HealthInfo().MaxDayCritical {
		//Reject any floor condition that involves us not being on the top floor and lasts for more than days you can survive in critical
		//This is because there is a risk of signing your own death if you agree as you may be forced to eat no food with no get out condition.
		a.RejectTreaty(msg)
		return
	}

	if treaty.Condition() == messages.AvailableFood {
		if treaty.Request() == messages.LeaveAmountFood && treaty.ConditionValue()-treaty.RequestValue() < 3 {
			// Reject treaty in which you would not be able to eat enough food to avoid critical level
			a.RejectTreaty(msg)
			return
		}
		if treaty.Request() == messages.LeavePercentFood && treaty.ConditionValue()*int(float64(100-treaty.RequestValue())/100) < 3 {
			// Reject treaty in which you would not be able to eat enough food to avoid critical level
			a.RejectTreaty(msg)
			return
		}
	}

	//Now we have ruled treaties that are always unacceptable, now to decide if agent will agree to an acceptable treaty
	//TODO: Develop the decision calculation.
	//For now: Middle range of selfishness - selfishness + opinion of agent proposing treaty - duration of treaty/4
	decision := 5 - a.selfishness + a.socialMemory[msg.SenderID()].favour - treaty.Duration()/4
	if decision > 0 {
		treaty.SignTreaty()
		a.AddTreaty(treaty)
		reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), true)
		a.SendMessage(reply)
		a.Log("Accepted treaty", infra.Fields{"proposerID": msg.SenderID(), "proposerFloor": msg.SenderFloor(), "treatyID": msg.TreatyID()})
		passingOnTo := a.Floor() + 1
		if msg.SenderFloor() > a.Floor() {
			passingOnTo = a.Floor() - 1
		}
		passItOn := messages.NewProposalMessage(a.ID(), a.Floor(), passingOnTo, treaty)
		a.SendMessage(passItOn)
	} else {
		a.RejectTreaty(msg)
	}

	// The code below can be used to accept all treaties by default.
	// treaty := msg.Treaty()
	// treaty.SignTreaty()
	// a.activeTreaties[msg.TreatyID()] = treaty
	// reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), true)
	// a.SendMessage(reply)
	// a.Log("Accepted treaty", Fields{"proposerID": msg.SenderID(), "proposerFloor": msg.SenderFloor(),
	// 	"treatyID": msg.TreatyID()})
}

func (a *CustomAgent5) HandleTreatyResponse(msg messages.TreatyResponseMessage) {
	if msg.Response() {
		treaty := a.ActiveTreaties()[msg.TreatyID()]
		treaty.SignTreaty()
		a.ActiveTreaties()[msg.TreatyID()] = treaty
		a.addToSocialFavour(msg.SenderID(), a.socialMemory[msg.SenderID()].favour+2)
	}
}
