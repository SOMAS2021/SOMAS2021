package team6

import (
	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
)

// Requests the agent above to leave food on the platform
func (a *CustomAgent6) requestLeaveFood() {
	healthInfo := a.HealthInfo()
	currentHP := a.HP()

	// HP levels based on MaxHP
	levels := levelsData{
		strongLevel:  healthInfo.MaxHP * 3 / 5,
		healthyLevel: healthInfo.MaxHP * 3 / 10,
		weakLevel:    healthInfo.MaxHP * 1 / 10,
	}

	// Sets the requested amount of food (reqAmount) to -1 if the agent does not request anything
	// Sets the requested amount of food (reqAmount) to the value requested by the agent
	var reqAmount int

	switch a.currBehaviour.string() {
	case "Altruist":
		reqAmount = -1

	case "Collectivist":
		if currentHP >= levels.weakLevel {
			reqAmount = -1
		} else {
			reqAmount = 2 // 2 is what is needed to go from the critical state to the weak level
		}

	case "Selfish":
		if currentHP >= levels.strongLevel {
			reqAmount = -1
		} else {
			reqAmount = int(a.foodRequired(levels.healthyLevel, a.HealthInfo()))
		}

	case "Narcissist":
		reqAmount = 4 * int(a.HealthInfo().Tau)

	default:
		reqAmount = -1
	}

	// Sends a request to the floor above
	if reqAmount != -1 {
		msg := messages.NewRequestLeaveFoodMessage(a.ID(), a.Floor(), a.Floor()-1, reqAmount)
		a.SendMessage(msg)
		a.Log("I sent a message", infra.Fields{"message": "RequestLeaveFood", "floor": a.Floor()})
	}
}

// Requests the agent above to take a precise amount of food
// The altruist and the collectivist do not request anything like that
// The selfish and the narcissist request the other agent to take nothing
func (a *CustomAgent6) RequestTakeFood() {

	var reqAmount int

	switch a.currBehaviour.string() {
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

// Handles RequestLeaveFood messages the agent receives
// Returns true if the agent accepts the request, false otherwise
func (a *CustomAgent6) HandleRequestLeaveFood(msg messages.RequestLeaveFoodMessage) {
	healthInfo := a.HealthInfo()
	currentHP := a.HP()

	// HP levels based on maximim HP
	levels := levelsData{
		strongLevel:  healthInfo.MaxHP * 3 / 5,
		healthyLevel: healthInfo.MaxHP * 3 / 10,
		weakLevel:    healthInfo.MaxHP * 1 / 10,
	}

	var reply bool

	switch a.currBehaviour.string() {
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

	if reply {
		a.reqLeaveFoodAmount = msg.Request()
		a.Log("I received a requestLeaveFood message and my response was true")
	} else {
		a.reqLeaveFoodAmount = -1
		a.Log("I received a requestLeaveFood message and my response was false")
	}
}

// Handles RequestTakeFood messages the agent receives
// Returns false, as our agents never accept to take a precise, fixed amount of food
func (a *CustomAgent6) HandleRequestTakeFood(msg messages.RequestTakeFoodMessage) {
	reply := false
	replyMessage := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), reply)
	a.SendMessage(replyMessage)

	if reply {
		a.reqLeaveFoodAmount = msg.Request()
		a.Log("I received a requestTakeFood message and my response was true")
	} else {
		a.reqLeaveFoodAmount = -1
		a.Log("I received a requestTakeFood message and my response was false")
	}
}

// Handles AskHP messages the agent receives
// Returns the agent's HP, unless the agent is a narcissist. In this case, he does not answer.
func (a *CustomAgent6) HandleAskHP(msg messages.AskHPMessage) {
	if a.currBehaviour.string() != "Narcissist" {
		reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), a.HP())
		a.SendMessage(reply)
		a.Log("I recieved an askHP message from ", infra.Fields{"senderFloor": msg.SenderFloor(), "myFloor": a.Floor()})
	}
}

// Handles AskFoodTaken messages the agent receives
// Returns the agent's last food intake, unless the agent is a narcissist. In this case, he does not answer.
func (a *CustomAgent6) HandleAskFoodTaken(msg messages.AskFoodTakenMessage) {
	if a.currBehaviour.string() != "Narcissist" {
		reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), int(a.lastFoodTaken))
		a.SendMessage(reply)
		a.Log("I recieved an askFoodTaken message from ", infra.Fields{"senderFloor": msg.SenderFloor(), "myFloor": a.Floor()})
	}
}

// Handles AskIntendedFoodTaken messages the agent receives
// Returns the agent's last food intake, which is approximately equal the intended food intake, unless the agent is a narcissist. In this case, he does not answer.
func (a *CustomAgent6) HandleAskIntendedFoodTaken(msg messages.AskIntendedFoodIntakeMessage) {
	if a.currBehaviour.string() != "Narcissist" {
		reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), int(a.lastFoodTaken))
		a.SendMessage(reply)
		a.Log("I recieved an askIntendedFoodTaken message from ", infra.Fields{"senderFloor": msg.SenderFloor(), "myFloor": a.Floor()})
	}
}

// Handles the messages that needs to be propagated
// Sends the messages to the target floor, unless the agent is a narcissist
func (a *CustomAgent6) HandlePropagate(msg messages.ProposeTreatyMessage) {
	// The Narcissist does not propagate treaties
	if a.currBehaviour.string() != "Narcissist" {
		treatyToPropagate := messages.NewProposalMessage(msg.SenderID(), msg.SenderFloor(), msg.TargetFloor(), msg.Treaty())
		a.SendMessage(treatyToPropagate)
		a.Log("I propogated a treaty")
	}
}

// Handles the responses from other agents to our treaty proposals
func (a *CustomAgent6) HandleTreatyResponse(msg messages.TreatyResponseMessage) {
	if msg.Response() {
		// Adds accepted treaty in the active treaties of the agent
		treaty := a.proposedTreaties[msg.TreatyID()]
		a.AddTreaty(treaty)
	}
	// Deletes the treaties for which we get an answer (yes or no) from our proposed treaty list
	delete(a.proposedTreaties, msg.TreatyID())
}

// Handles the treaty proposals we get from other agents
func (a *CustomAgent6) HandleProposeTreaty(msg messages.ProposeTreatyMessage) {
	treaty := msg.Treaty()

	// Checks if we benefit from a treaty using the function "a.considerTreaty".
	// This function returns true if we should accept the treaty
	if a.considerTreaty(&treaty) {
		// Propagates the accepted treaty only if treaty doesn't already exist (avoids infinite loops)
		if _, exists := a.ActiveTreaties()[msg.TreatyID()]; !exists {
			a.proposeTreaty(treaty)
		}
		// Signs and adds the treaty to our active treaties
		treaty.SignTreaty()
		a.AddTreaty(treaty)

		// Replies with acceptance message
		reply := messages.NewTreatyResponseMessage(a.ID(), a.Floor(), msg.SenderFloor(), true, treaty.ID(), treaty.ProposerID())
		a.SendMessage(reply)
		a.Log("I accepted a treaty proposed from ", infra.Fields{"senderFloor": msg.SenderFloor(), "myFloor": a.Floor(), "my social motive": a.currBehaviour.string()})

	} else {
		a.Log("I rejected a treaty proposed from ", infra.Fields{"senderFloor": msg.SenderFloor(), "myFloor": a.Floor(), "my social motive": a.currBehaviour.string()})
	}

}
