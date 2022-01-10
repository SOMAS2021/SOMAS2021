package team7agent1

import (
	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/health"
)

func (a *CustomAgent7) HandleBadTreaty(msg messages.ProposeTreatyMessage) {
	reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), false)
	a.SendMessage(reply)
	a.Log("Rejected treaty", infra.Fields{"proposerID": msg.SenderID(), "proposerFloor": msg.SenderFloor(), "treatyID": msg.TreatyID()})
}

func (a *CustomAgent7) HandleProposeTreaty(msg messages.ProposeTreatyMessage) {
	treaty := msg.Treaty()
	// switch{
	// case (treaty.ConditionOp() == messages.LE || treaty.ConditionOp() == messages.LT) && treaty.Condition() == messages.HP || treaty.Condition() == messages.AvailableFood :
	// 	a.HandleBadTreaty(msg)
	// 	return
	// case (treaty.ConditionOp() == messages.GE || treaty.ConditionOp() == messages.GT) && treaty.Condition() == messages.Floor:
	// 	a.HandleBadTreaty(msg)
	// 	return
	// case treaty.Request() == messages.Inform:
	// 	a.HandleBadTreaty(msg)
	// 	return
	// case treaty.Request() == messages.LeaveAmountFood || treaty.Request() == messages.LeavePercentFood && treaty.RequestOp() == messages.LT) || treaty.RequestOp() == messages.LE:
	// 	a.HandleBadTreaty(msg)
	// 	return
	// case treaty.Condition() == messages.HP && treaty.ConditionValue() < a.HealthInfo().WeakLevel:
	// 	a.HandleBadTreaty(msg)
	// 	return
	// case a.clashOfTreaties(treaty):
	// 	a.HandleBadTreaty(msg)
	// 	return
	// }

	if ((treaty.ConditionOp() == messages.LE || treaty.ConditionOp() == messages.LT) && treaty.Condition() == messages.HP || treaty.Condition() == messages.AvailableFood) ||
		((treaty.ConditionOp() == messages.GE || treaty.ConditionOp() == messages.GT) && treaty.Condition() == messages.Floor) ||
		(treaty.Request() == messages.Inform) ||
		((treaty.Request() == messages.LeaveAmountFood || treaty.Request() == messages.LeavePercentFood && treaty.RequestOp() == messages.LT) || treaty.RequestOp() == messages.LE) ||
		(treaty.Condition() == messages.HP && treaty.ConditionValue() < a.HealthInfo().WeakLevel) ||
		(a.clashOfTreaties(treaty)) {
		a.HandleBadTreaty(msg)
		return
	}

	//more than 5 days don't sign, don't sign if responsiveness too low
	if a.personality.conscientiousness > 33 && a.behaviour.responsiveness > 50 && treaty.Duration() < 6 {
		treaty.SignTreaty()
		a.AddTreaty(treaty)
		reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), true)
		a.SendMessage(reply)
		a.Log("Accepted treaty", infra.Fields{"proposerID": msg.SenderID(), "proposerFloor": msg.SenderFloor(), "treatyID": msg.TreatyID()})
	} else {
		a.HandleBadTreaty(msg)
	}

}

func (a *CustomAgent7) HandleTreatyResponse(msg messages.TreatyResponseMessage) {
	if msg.Response() {
		treaty := a.ActiveTreaties()[msg.TreatyID()]
		treaty.SignTreaty()
		a.ActiveTreaties()[msg.TreatyID()] = treaty
	}
}

func (a *CustomAgent7) criticalStateTreaty() {
	foodCtoW := int(health.FoodRequired(a.HealthInfo().HPCritical, a.HealthInfo().WeakLevel, a.HealthInfo()))
	if a.HP() <= a.HealthInfo().HPCritical && a.Floor() != 1 {
		tr := messages.NewTreaty(messages.HP, 25, messages.LeaveAmountFood, foodCtoW, messages.GT, messages.GT, 4, a.ID())
		msg := messages.NewProposalMessage(a.ID(), a.Floor(), a.Floor()-1, *tr)
		a.SendMessage(msg)
		a.Log("Team 7 Agent has sent a Treaty")
	}
}

func (a *CustomAgent7) propagateTreatyUpwards() {
	foodCtoW := int(health.FoodRequired(a.HealthInfo().HPCritical, a.HealthInfo().WeakLevel, a.HealthInfo()))
	for _, tr := range a.ActiveTreaties() {
		if tr.Request() == messages.LeaveAmountFood && a.Floor() != 1 {
			trnew := messages.NewTreaty(tr.Condition(), tr.ConditionValue(), messages.LeaveAmountFood, tr.RequestValue()+foodCtoW, tr.ConditionOp(), tr.RequestOp(), tr.Duration(), tr.ProposerID())
			msg := messages.NewProposalMessage(a.ID(), a.Floor(), a.Floor()-1, *trnew)
			a.SendMessage(msg)
			a.Log("Team 7 Agent has sent a Treaty")
		}
	}
}

func (a *CustomAgent7) treatyOnFloorChange() {
	foodCtoW := int(health.FoodRequired(a.HealthInfo().HPCritical, a.HealthInfo().WeakLevel, a.HealthInfo()))
	if len(a.opMem.orderPrevFloors) != 0 && a.opMem.orderPrevFloors[len(a.opMem.orderPrevFloors)-1] != a.Floor() && a.Floor() != 1 {
		if a.opMem.orderPrevFloors[len(a.opMem.orderPrevFloors)-1]/2 >= a.Floor() {
			if a.personality.conscientiousness > 50 {
				tr := messages.NewTreaty(messages.HP, 25, messages.LeaveAmountFood, foodCtoW, messages.GT, messages.GT, 4, a.ID())
				msg := messages.NewProposalMessage(a.ID(), a.Floor(), a.Floor()-1, *tr)
				a.SendMessage(msg)
				a.Log("Team 7 Agent has sent a Treaty")
			}
		}
	}
}

func (a *CustomAgent7) desperationTreaty() {
	if a.DaysAtCritical() == a.HealthInfo().MaxDayCritical-2 && a.Floor() != 1 {
		tr := messages.NewTreaty(messages.AvailableFood, 50, messages.LeavePercentFood, 50, messages.GT, messages.GE, 4, a.ID())
		msg := messages.NewProposalMessage(a.ID(), a.Floor(), a.Floor()-1, *tr)
		a.SendMessage(msg)
		a.Log("Team 7 Agent has sent a Treaty")
	}
}

func (a *CustomAgent7) clashOfTreaties(tnew messages.Treaty) bool {

	//need to check this treaty with all current signed active treaties
	for _, tActive := range a.ActiveTreaties() {

		if tnew.Condition() == tActive.Condition() {

			if a.mutuallyExclusive(tnew.ConditionOp(), tActive.ConditionOp(), tnew.ConditionValue(), tActive.ConditionValue()) {
				return false // if mutually exclusive conditions of same type, we can can handle any request
			}
			if tnew.Request() == tActive.Request() {
				return true // cannot have conditions demanding different amounts of the same type of thing
			}

		} else {

			if tnew.Request() == tActive.Request() {
				// cannot have conditions demanding different amounts of the same type of thing, unless they intersect
				// ex: HP < 5 -> Food < 8, Floor < 10 -> Food > 12 ==> these cannot coexist!
				return a.mutuallyExclusive(tnew.RequestOp(), tActive.RequestOp(), tnew.RequestValue(), tActive.RequestValue())
			}
		}
	}
	return false
}

//This function was edited using a template of code from Ben Stobbs
func (a *CustomAgent7) mutuallyExclusive(op1 messages.Op, op2 messages.Op, val1 int, val2 int) bool {
	switch op1 {
	case messages.LT:
		return ((op2 == messages.EQ && val1 <= val2) ||
			(op2 == messages.GT && val1 <= val2+1) ||
			(op2 == messages.GE && val1 <= val2))

	case messages.LE:
		return ((op2 == messages.EQ && val1 < val2) ||
			(op2 == messages.GT && val1 <= val2) ||
			(op2 == messages.GE && val1 < val2))

	case messages.GT:
		return ((op2 == messages.EQ && val1 >= val2) ||
			(op2 == messages.LT && val1 >= (val2+1)) ||
			(op2 == messages.LE && val1 >= val2))

	case messages.GE:
		return ((op2 == messages.EQ && val1 > val2) ||
			(op2 == messages.LT && val1 >= val2) ||
			(op2 == messages.LE && val1 > val2))

	case messages.EQ:
		return ((op2 == messages.EQ && val1 != val2) ||
			(op2 == messages.LT && val1 >= val2) ||
			(op2 == messages.LE && val1 > val2) ||
			(op2 == messages.GT && val1 <= val2) ||
			(op2 == messages.GE && val1 < val2))

	default:
		return false
	}
}
