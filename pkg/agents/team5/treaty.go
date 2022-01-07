package team5

import (
	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
)

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

func (a *CustomAgent5) treatyOverride() {
	for _, treaty := range a.ActiveTreaties() {
		if treaty.SignatureCount() > 1 {
			switch {
			case treaty.Condition() == messages.HP && treaty.ConditionOp() == messages.EQ && a.HP() == treaty.ConditionValue():
				fallthrough
			case treaty.Condition() == messages.HP && treaty.ConditionOp() == messages.GT && a.HP() > treaty.ConditionValue():
				fallthrough
			case treaty.Condition() == messages.HP && treaty.ConditionOp() == messages.GE && a.HP() >= treaty.ConditionValue():
				fallthrough
			case treaty.Condition() == messages.Floor && treaty.ConditionOp() == messages.EQ && a.Floor() == treaty.ConditionValue():
				fallthrough
			case treaty.Condition() == messages.Floor && treaty.ConditionOp() == messages.LT && a.Floor() < treaty.ConditionValue():
				fallthrough
			case treaty.Condition() == messages.Floor && treaty.ConditionOp() == messages.LE && a.Floor() <= treaty.ConditionValue():
				fallthrough
			case treaty.Condition() == messages.AvailableFood && treaty.ConditionOp() == messages.EQ && a.CurrPlatFood() == food.FoodType(treaty.ConditionValue()):
				fallthrough
			case treaty.Condition() == messages.AvailableFood && treaty.ConditionOp() == messages.GT && a.CurrPlatFood() > food.FoodType(treaty.ConditionValue()):
				fallthrough
			case treaty.Condition() == messages.AvailableFood && treaty.ConditionOp() == messages.GE && a.CurrPlatFood() >= food.FoodType(treaty.ConditionValue()):
				a.overrideCalculation(treaty)
			}
		}
	}
}

func (a *CustomAgent5) treatiesCanCoexist(t1 messages.Treaty, t2 messages.Treaty) bool {
	conditionsIntersect := statementsIntersect(t1.ConditionOp(), t1.ConditionValue(), t2.ConditionOp(), t2.ConditionValue())

	// If the treaties are based on the same type of condition, they can coexist
	// when the conditions are mutually exclusive (can never occur simultaneously)
	if t1.Condition() == t2.Condition() && !conditionsIntersect {
		return true
	}

	// Otherwise, if the treaties are based on different conditions, or the conditions
	// can occurr simultaneously, treaties of different types should be rejected
	// because it is impossible to deduce if their requests can be fulfilled simultaneously
	if t1.Request() != t2.Request() {
		return false
	}

	// Otherwise, treaties can coexist only if the requests can be fulfilled simultaneously
	return statementsIntersect(t1.RequestOp(), t1.RequestValue(), t2.RequestOp(), t2.RequestValue())
}

func (a *CustomAgent5) treatyConflicts(treaty messages.Treaty) bool {
	for _, activeTreaty := range a.ActiveTreaties() {
		if !a.treatiesCanCoexist(treaty, activeTreaty) {
			return true
		}
	}
	return false
}

func (a *CustomAgent5) RejectTreaty(msg messages.ProposeTreatyMessage) {
	reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), false)
	a.SendMessage(reply)
	a.Log("Rejected treaty", infra.Fields{"proposerID": msg.SenderID(), "proposerFloor": msg.SenderFloor(), "treatyID": msg.TreatyID()})
}

func (a *CustomAgent5) HandleProposeTreaty(msg messages.ProposeTreatyMessage) {
	treaty := msg.Treaty()

	// Eliminate unacceptable treaties
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
}

func (a *CustomAgent5) HandleTreatyResponse(msg messages.TreatyResponseMessage) {
	if msg.Response() {
		treaty, ok := a.ActiveTreaties()[msg.TreatyID()]
		if ok {
			treaty.SetCount(treaty.SignatureCount() + 1)
			a.ActiveTreaties()[msg.TreatyID()] = treaty
			a.addToSocialFavour(msg.SenderID(), a.socialMemory[msg.SenderID()].favour+2)
		}
	}
}
