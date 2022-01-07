package team7agent1

import (
	// "math"
	// "math/rand"

	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
	// "github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
	// "github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/health"
	// "github.com/google/uuid"
)

func (a *CustomAgent7) RejectTreaty(msg messages.ProposeTreatyMessage) {
	reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), false)
	a.SendMessage(reply)
	a.Log("Rejected treaty", infra.Fields{"proposerID": msg.SenderID(), "proposerFloor": msg.SenderFloor(), "treatyID": msg.TreatyID()})
}

func (a *CustomAgent7) HandleProposeTreaty(msg messages.ProposeTreatyMessage) {
	treaty := msg.Treaty()
	if treaty.Request() == messages.Inform {
		// Reject any inform treaties as I'm not sure what they are or how to handle them
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
	// decision := 5 - a.selfishness + a.socialMemory[msg.SenderID()].favour - treaty.Duration()/4
	// if decision > 0 {
	// 	treaty.SignTreaty()
	// 	a.AddTreaty(treaty)
	// 	reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), true)
	// 	a.SendMessage(reply)
	// 	a.Log("Accepted treaty", infra.Fields{"proposerID": msg.SenderID(), "proposerFloor": msg.SenderFloor(), "treatyID": msg.TreatyID()})
	// 	passingOnTo := a.Floor() + 1
	// 	if msg.SenderFloor() > a.Floor() {
	// 		passingOnTo = a.Floor() - 1
	// 	}
	// 	passItOn := messages.NewProposalMessage(a.ID(), a.Floor(), passingOnTo, treaty)
	// 	a.SendMessage(passItOn)
	// } else {
	// 	a.RejectTreaty(msg)
	// }

	// The code below can be used to accept all treaties by default.
	// treaty := msg.Treaty()
	// treaty.SignTreaty()
	// a.activeTreaties[msg.TreatyID()] = treaty
	// reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), true)
	// a.SendMessage(reply)
	// a.Log("Accepted treaty", Fields{"proposerID": msg.SenderID(), "proposerFloor": msg.SenderFloor(),
	// 	"treatyID": msg.TreatyID()})
}

func (a *CustomAgent7) HandleTreatyResponse(msg messages.TreatyResponseMessage) {
	if msg.Response() {
		treaty := a.ActiveTreaties()[msg.TreatyID()]
		treaty.SignTreaty()
		a.ActiveTreaties()[msg.TreatyID()] = treaty
		//a.addToSocialFavour(msg.SenderID(), a.socialMemory[msg.SenderID()].favour+2)
	}
}