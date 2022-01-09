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

func (a *CustomAgent7) HandleBadTreaty(msg messages.ProposeTreatyMessage) {
	reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), false)
	a.SendMessage(reply)
	a.Log("Rejected treaty", infra.Fields{"proposerID": msg.SenderID(), "proposerFloor": msg.SenderFloor(), "treatyID": msg.TreatyID()})
}

func (a *CustomAgent7) HandleProposeTreaty(msg messages.ProposeTreatyMessage) {
	treaty := msg.Treaty()

	if ((treaty.ConditionOp() == messages.LE || treaty.ConditionOp() == messages.LT) && treaty.Condition() == messages.HP || treaty.Condition() == messages.AvailableFood) ||
		((treaty.ConditionOp() == messages.GE || treaty.ConditionOp() == messages.GT) && treaty.Condition() == messages.Floor) ||
		(treaty.Request() == messages.Inform) ||
		((treaty.Request() == messages.LeaveAmountFood || treaty.Request() == messages.LeavePercentFood && treaty.RequestOp() == messages.LT) || treaty.RequestOp() == messages.LE) ||
		(treaty.Condition() == messages.HP && treaty.ConditionValue() < a.HealthInfo().WeakLevel*3) ||
		(a.HP() <= a.HealthInfo().HPCritical && a.DaysAtCritical() >= a.HealthInfo().MaxDayCritical-3) ||
		(treaty.Request() == messages.LeavePercentFood && treaty.RequestValue() > 50) ||
		(a.Clashoftreaties(treaty)) {
		a.HandleBadTreaty(msg)
		return
	}

	if a.personality.conscientiousness/2+a.personality.extraversion > 50 {
		treaty.SignTreaty()
		a.AddTreaty(treaty)
		reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), true)
		a.SendMessage(reply)

		a.Log("Accepted treaty", infra.Fields{"proposerID": msg.SenderID(), "proposerFloor": msg.SenderFloor(), "treatyID": msg.TreatyID()})
		propagate := a.Floor() - 1
		if msg.SenderFloor() < a.Floor() {
			propagate = a.Floor() + 1
		}
		propagationoftreaty := messages.NewProposalMessage(a.ID(), a.Floor(), propagate, treaty)
		a.SendMessage(propagationoftreaty)
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

func (a *CustomAgent7) SendTreaty() {
	tr := messages.NewTreaty(messages.HP, 20, messages.LeavePercentFood, 95, messages.GT, messages.GT, 3, a.ID()) //remember the treaty we proposed
	msg := messages.NewProposalMessage(a.BaseAgent().ID(), a.Floor(), a.Floor()+1, *tr)
	a.SendMessage(msg)
	a.Log("Team 7 Agent has sent a Treaty")
}

func (a *CustomAgent7) Clashoftreaties(t_new messages.Treaty) bool {

	//need to check this treaty with all current signed active treaties
	for _, t_active := range a.ActiveTreaties() {

		if t_new.Condition() == t_active.Condition() {

			if a.mutuallyExclusive(t_new.ConditionOp(), t_active.ConditionOp(), t_new.ConditionValue(), t_active.ConditionValue()) {
				return false // if mutually exclusive conditions of same type, we can can handle any request
			}
			if t_new.Request() == t_active.Request() {
				return true // cannot have conditions demanding different amounts of the same type of thing
			}

		} else {

			if t_new.Request() == t_active.Request() {
				// cannot have conditions demanding different amounts of the same type of thing, unless they intersect
				// ex: HP < 5 -> Food < 8, Floor < 10 -> Food > 12 ==> these cannot coexist!
				return a.mutuallyExclusive(t_new.RequestOp(), t_active.RequestOp(), t_new.RequestValue(), t_active.RequestValue())
			}
		}
	}
	return false
}

func (a *CustomAgent7) mutuallyExclusive(op1 messages.Op, op2 messages.Op, value1 int, value2 int) bool {
	switch op1 {
	case messages.LT:
		return ((op2 == messages.EQ && value1 <= value2) ||
			(op2 == messages.GT && value1 <= value2+1) ||
			(op2 == messages.GE && value1 <= value2))

	case messages.LE:
		return ((op2 == messages.EQ && value1 < value2) ||
			(op2 == messages.GT && value1 <= value2) ||
			(op2 == messages.GE && value1 < value2))

	case messages.GT:
		return ((op2 == messages.EQ && value1 >= value2) ||
			(op2 == messages.LT && value1 >= (value2+1)) ||
			(op2 == messages.LE && value1 >= value2))

	case messages.GE:
		return ((op2 == messages.EQ && value1 > value2) ||
			(op2 == messages.LT && value1 >= value2) ||
			(op2 == messages.LE && value1 > value2))

	case messages.EQ:
		return ((op2 == messages.EQ && value1 != value2) ||
			(op2 == messages.LT && value1 >= value2) ||
			(op2 == messages.LE && value1 > value2) ||
			(op2 == messages.GT && value1 <= value2) ||
			(op2 == messages.GE && value1 < value2))

	default:
		return false
	}
}
