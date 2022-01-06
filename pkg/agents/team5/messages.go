package team5

import (
	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
	"github.com/google/uuid"
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

func (a *CustomAgent5) updateSocialMemory(senderID uuid.UUID, senderFloor int) {
	if !a.memoryIdExists(senderID) {
		a.newMemory(senderID)
	}
	a.resetDaysSinceLastSeen(senderID)
	a.surroundingAgents[senderFloor-a.Floor()] = senderID
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

func (a *CustomAgent5) RejectTreaty(msg messages.ProposeTreatyMessage) {
	reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), false)
	a.SendMessage(reply)
	a.Log("Rejected treaty", infra.Fields{"proposerID": msg.SenderID(), "proposerFloor": msg.SenderFloor(), "treatyID": msg.TreatyID()})
}

func (a *CustomAgent5) HandleProposeTreaty(msg messages.ProposeTreatyMessage) {
	treaty := msg.Treaty()

	switch {
	case treaty.Request() == messages.Inform:
		// Reject any inform treaties as I'm not sure what they are or how to handle them
		fallthrough
	case a.treatyConflicts(treaty):
		fallthrough
	case treaty.Request() == messages.LeavePercentFood && (treaty.RequestValue() > 100 || treaty.RequestValue() < 0):
		fallthrough
	case (treaty.Request() == messages.LeaveAmountFood && treaty.RequestOp() == messages.EQ) || treaty.RequestOp() == messages.LE || treaty.RequestOp() == messages.LT:
		//Reject all treaties that ask you to leave less food, don't see why you would do this
		//Reject any treaty that asks you to leave a specific amount of food as this would lead to multiple people just not eating as only 1 agent can eat and fufill that criteria
		fallthrough
	case (treaty.ConditionOp() == messages.LE || treaty.ConditionOp() == messages.LT) && (treaty.Condition() == messages.HP || treaty.Condition() == messages.AvailableFood):
		// Don't agree to treaty which is bounded by you having lower HP
		// Don't agree to treaty which is bounded by the amount of food on platform being low
		// Both of these are conditions where you will be desperate for food so shouldn't be limitting yourself if you do get the opurrtunity to eat
		fallthrough
	case (treaty.ConditionOp() == messages.GE || treaty.ConditionOp() == messages.GT) && treaty.Condition() == messages.Floor:
		// Don't agree to treaty which is bounded by you being lower down in the tower for same as reasons above
		fallthrough
	case treaty.Condition() == messages.HP && treaty.ConditionValue() < a.HealthInfo().WeakLevel*2:
		//Reject any HP condition based treaty if the condition is too strict, and you would be put into critical by not eating
		fallthrough
	case treaty.Condition() == messages.Floor && treaty.ConditionValue() != 1 && treaty.Duration() >= a.HealthInfo().MaxDayCritical:
		//Reject any floor condition that involves us not being on the top floor and lasts for more than days you can survive in critical
		//This is because there is a risk of signing your own death if you agree as you may be forced to eat no food with no get out condition
		fallthrough
	case treaty.Condition() == messages.AvailableFood && treaty.Request() == messages.LeaveAmountFood && treaty.ConditionValue()-treaty.RequestValue() < 3:
		// Reject treaty in which you would not be able to eat enough food to avoid critical level
		fallthrough
	case treaty.Condition() == messages.AvailableFood && treaty.Request() == messages.LeavePercentFood && treaty.ConditionValue()*int(float64(100-treaty.RequestValue())/100) < 3:
		// Reject treaty in which you would not be able to eat enough food to avoid critical level
		a.RejectTreaty(msg)
		return
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
		a.ActiveTreaties()[msg.TreatyID()] = treaty
		a.addToSocialFavour(msg.SenderID(), a.socialMemory[msg.SenderID()].favour+2)
	}
}
