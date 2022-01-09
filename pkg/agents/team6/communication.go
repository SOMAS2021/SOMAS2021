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

	var reqAmount int

	switch a.currBehaviour.String() {
	case "Altruist":
		reqAmount = -1

	case "Collectivist":
		if currentHP >= levels.weakLevel {
			reqAmount = -1
		} else {
			reqAmount = 2 // to is what is needed to go from the critical state to the weak level
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

	var reqAmount int

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

	var reply bool

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
		a.updateTrust(-1, msg.SenderID()) // how dare you even ask

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
	reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), int(a.lastFoodTaken))
	a.SendMessage(reply)
	a.Log("I recieved an askIntendedFoodTaken message from ", infra.Fields{"senderFloor": msg.SenderFloor(), "myFloor": a.Floor()})
}

func (a *CustomAgent6) HandlePropagate(msg messages.ProposeTreatyMessage) {
	// The Narcissist does not propagate treaties
	if a.currBehaviour.String() != "Narcissist" {
		treatyToPropagate := messages.NewProposalMessage(msg.SenderID(), msg.SenderFloor(), msg.TargetFloor(), msg.Treaty())
		a.SendMessage(treatyToPropagate)
		a.Log("I propogated a treaty")
	}
}

func (a *CustomAgent6) HandleTreatyResponse(msg messages.TreatyResponseMessage) {
	if msg.Response() {
		treaty := a.proposedTreaties[msg.TreatyID()]
		a.AddTreaty(treaty)
		a.updateTrust(3, msg.SenderID()) // great - they must be cool
	} else {
		a.updateTrust(-2, msg.SenderID()) // we trust them less if they refuse our treaty - must be up to something
	}
	delete(a.proposedTreaties, msg.TreatyID())
}

func (a *CustomAgent6) HandleProposeTreaty(msg messages.ProposeTreatyMessage) {
	treaty := msg.Treaty()

	// add trust

	// }
	// check if we benefit from a treaty
	if a.considerTreaty(&treaty) {
		// Propagate only if treaty doesn't already exist (avoids infinite loops)
		if _, exists := a.ActiveTreaties()[msg.TreatyID()]; !exists {
			a.ProposeTreaty(treaty)
		}
		treaty.SignTreaty()
		a.AddTreaty(treaty)

		// reply with acceptance message
		reply := messages.NewTreatyResponseMessage(a.ID(), a.Floor(), msg.SenderFloor(), true, treaty.ID(), treaty.ProposerID())
		a.SendMessage(reply)
		a.updateTrust(2, msg.SenderID()) // good treaty - these guys are probably nice :)
		a.Log("I accepted a treaty proposed from ", infra.Fields{"senderFloor": msg.SenderFloor(), "myFloor": a.Floor(), "my social motive": a.currBehaviour.String()})

	} else {
		a.updateTrust(-1, msg.SenderID()) // bad treaty - these guys are trying to sabotage us >:)
		a.Log("I rejected a treaty proposed from ", infra.Fields{"senderFloor": msg.SenderFloor(), "myFloor": a.Floor(), "my social motive": a.currBehaviour.String()})
	}

}
