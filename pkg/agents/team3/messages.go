package team3

import (
	"math/rand"

	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
)

//Upon receipt of message define affected emotions
// ACK MESSAGE
//If x time passed no message received/acked morale decrease
//Include if ack message same user ID occurs x+1 times, morale increase
//If stubborness = y+1, discard, a.k.a. leave unread

//case 0, we are going to check our hunger and propose a treaty or, if we already have one, ask about the people bellow's health (cause why not)
func (a *CustomAgent3) sendMsgHPRelated() {
	if a.knowledge.foodLastSeen < 10 && a.HP() < a.HealthInfo().HPReqCToW { //do this properly with Eds help
		if a.knowledge.aboveFoodTreaty {
			msg := messages.NewAskHPMessage(a.BaseAgent().ID(), a.Floor(), a.Floor()-1)
			a.SendMessage(msg)
			a.Log("I sent a message", infra.Fields{"message": "AskHP"})
		} else {
			tr := messages.NewTreaty(messages.HP, 20, messages.LeaveAmountFood, 20, messages.GT, messages.GT, 20, a.ID()) //generalise later
			msg := messages.NewProposalMessage(a.BaseAgent().ID(), a.Floor(), a.Floor()+1, *tr)
			a.SendMessage(msg)
			a.Log("I sent a treaty")
		}

	} else {
		msg := messages.NewAskHPMessage(a.BaseAgent().ID(), a.Floor(), a.Floor()+1)
		a.SendMessage(msg)
		a.Log("I sent a message", infra.Fields{"message": "AskHP"})
	}
}

func (a *CustomAgent3) sendMsgFoodTaken() {
	msg := messages.NewAskHPMessage(a.BaseAgent().ID(), a.Floor(), a.Floor()-1)
	a.SendMessage(msg)
	a.Log("I sent a message", infra.Fields{"message": "AskHP"})
}

func (a *CustomAgent3) ticklyMessage() {

	r := rand.Intn(5)
	switch r {
	case 0:
		a.sendMsgHPRelated()
	case 1:
		msg := messages.NewAskHPMessage(a.BaseAgent().ID(), a.Floor(), a.Floor()-1)
		a.SendMessage(msg)
		a.Log("I sent a message", infra.Fields{"message": "AskHP"})
	case 2:
		msg := messages.NewAskIntendedFoodIntakeMessage(a.BaseAgent().ID(), a.Floor(), a.Floor()-1)
		a.SendMessage(msg)
		a.Log("I sent a message", infra.Fields{"message": "AskIntendedFoodIntake"})
	case 3:
		msg := messages.NewRequestLeaveFoodMessage(a.BaseAgent().ID(), a.Floor(), a.Floor()+1, a.requestLeaveFoodAmt())
		a.SendMessage(msg)
		a.Log("I sent a message", infra.Fields{"message": "RequestLeaveFood"})
	case 4:
		msg := messages.NewRequestTakeFoodMessage(a.BaseAgent().ID(), a.Floor(), a.Floor()+1, 10) //make func to determine value
		a.SendMessage(msg)
		a.Log("I sent a message", infra.Fields{"message": "RequestTakeFood"})
	}
}
func (a *CustomAgent3) message() {
	receivedMsg := a.ReceiveMessage()
	if receivedMsg != nil {
		a.Log("Custom agent 3 each run:", infra.Fields{"floor": a.Floor(), "hp": a.HP(), "Mood": a.vars.mood, "Morality": a.vars.morality})
		receivedMsg.Visit(a)
		a.Log("Custom agent 3 each run:", infra.Fields{"floor": a.Floor(), "hp": a.HP(), "Mood": a.vars.mood, "Morality": a.vars.morality})
	} else {
		a.ticklyMessage()
		a.Log("I got nothing")
	}

}

func max(x, y, z int) int {
	if x < y && z < y {
		return y
	}
	if x < z {
		return z
	}
	return x
}

func min(x, y, z int) int {
	if x > y && z > y {
		return y
	}
	if x > z {
		return z
	}
	return x
}

func (a *CustomAgent3) requestTakeFoodAmt() int {
	foodReqAmt := a.foodReqCalc(a.HP(), a.HP()) //food required to keep same HP
	if a.vars.morality >= 70 {
		return max(a.decisions.foodToEat, foodReqAmt, int(a.knowledge.foodLastEaten)) //we would want people to eat as much as sustainable amount due to high morality
	} else if a.vars.morality < 70 && a.vars.morality > 30 {
		return min(a.decisions.foodToEat, foodReqAmt, int(a.knowledge.foodLastEaten)) //we want people above to eat as little as sustainable
	} else {
		return min(5, int(a.knowledge.foodLastEaten), 6) //we want people to take least food possible
	}
}

func (a *CustomAgent3) requestLeaveFoodAmt() int {
	if a.HP() >= 70 {
		foodReqAmt := a.foodReqCalc(a.HP(), a.HP()-5)
		if a.vars.morality >= 70 {
			foodReqAmt -= 5
		}
		if a.vars.morality <= 30 {
			foodReqAmt += 5
		}
		return foodReqAmt

	} else if a.HP() < 70 && a.HP() > 30 {
		foodReqAmt := a.foodReqCalc(a.HP(), a.HP())
		if a.vars.morality >= 70 {
			foodReqAmt -= 5
		}
		if a.vars.morality <= 30 {
			foodReqAmt += 5
		}
		return foodReqAmt
	} else {
		foodReqAmt := a.foodReqCalc(a.HP(), a.HP()+5)
		if a.vars.morality >= 70 {
			foodReqAmt -= 5
		}
		if a.vars.morality <= 30 {
			foodReqAmt += 5
		}
		return foodReqAmt
	}

}

func (a *CustomAgent3) HandleAskHP(msg messages.AskHPMessage) {
	a.Log("I recieved an askHP message from ", infra.Fields{"floor": msg.SenderFloor()})
	if a.read() {
		changeInStubbornness(a, 5, -1) //value could be different
		a.updateFriendship(msg.SenderID(), 1)
		changeInMood(a, 5, 10, 1)
		reply := msg.Reply(a.BaseAgent().ID(), a.Floor(), msg.SenderFloor()-a.Floor(), a.HP())
		a.SendMessage(reply)
		a.Log("I recieved an askHP message from ", infra.Fields{"floor": msg.SenderFloor()})
	}
}

func (a *CustomAgent3) HandleAskFoodTaken(msg messages.AskFoodTakenMessage) {
	a.Log("I recieved an askFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor()})
	if a.read() {
		if a.HP() < a.knowledge.lastHP {
			changeInStubbornness(a, 5, 1)

			changeInMood(a, 5, 10, -1)
		}
		reply := msg.Reply(a.BaseAgent().ID(), a.Floor(), msg.SenderFloor()-a.Floor(), int(a.knowledge.foodLastEaten))
		a.SendMessage(reply)
		a.Log("I sent a replyFoodTaken message to ", infra.Fields{"floor": msg.SenderFloor()})
	}
}

func (a *CustomAgent3) HandleAskIntendedFoodTaken(msg messages.AskIntendedFoodIntakeMessage) {
	friendship := a.knowledge.friends[msg.SenderID()]
	if a.read() {
		changeInStubbornness(a, 2, 1)
		if friendship != 0 {
			changeInMood(a, 5, 10, 1)
		}
		//add critical state effect
		reply := msg.Reply(a.BaseAgent().ID(), a.Floor(), msg.SenderFloor()-a.Floor(), a.decisions.foodToEat)
		a.SendMessage(reply)
		a.Log("I recieved an askIntendedFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor()})
	}
}

func (a *CustomAgent3) HandleRequestLeaveFood(msg messages.RequestLeaveFoodMessage) {
	friendship := a.knowledge.friends[msg.SenderID()]
	if a.read() {
		if a.HP() < a.knowledge.lastHP {
			if a.vars.stubbornness > 70 {
				changeInStubbornness(a, 5, 1)
			} else {
				changeInStubbornness(a, 2, -1)
			}
		}
		if a.vars.morality > 50 { //want to implement effects of sender.floor
			changeInMorality(a, 5, 10, 1)
		} else {
			if friendship != 0 {
				changeInMorality(a, 5, 10, -1)
			}
		}
		if a.vars.mood < 30 { //can we see when we are in critical state
			changeInMood(a, 5, 10, -1)
		}
		reply := msg.Reply(a.BaseAgent().ID(), a.Floor(), msg.SenderFloor()-a.Floor(), true)
		a.SendMessage(reply)
		a.Log("I recieved a requestLeaveFood message from ", infra.Fields{"floor": msg.SenderFloor()})
	}
}

func (a *CustomAgent3) HandleRequestTakeFood(msg messages.RequestTakeFoodMessage) {
	if a.read() {
		if a.HP() < a.knowledge.lastHP {
			if a.vars.stubbornness > 70 {
				changeInStubbornness(a, 5, 1)
			} else {
				changeInStubbornness(a, 2, -1)
			}
		}
		if a.vars.morality > 50 { //want to implement effects of friendship
			changeInMorality(a, 5, 10, 1)
		} else {
			changeInMorality(a, 5, 10, -1)
		}
		if a.vars.mood < 30 {
			changeInMood(a, 5, 10, -1)
		}
		reply := msg.Reply(a.BaseAgent().ID(), a.Floor(), msg.SenderFloor()-a.Floor(), true)
		a.SendMessage(reply)
		a.Log("I recieved a requestTakeFood message from ", infra.Fields{"floor": msg.SenderFloor()})
	}
}

func (a *CustomAgent3) HandleResponse(msg messages.BoolResponseMessage) {
	response := msg.Response()
	a.Log("I recieved a Response message from ", infra.Fields{"floor": msg.SenderFloor(), "response": response})
}

func (a *CustomAgent3) HandleStateFoodTaken(msg messages.StateFoodTakenMessage) {
	statement := msg.Statement()
	if statement > a.decisions.foodToEat {
		changeInStubbornness(a, 5, -1)
		changeInMood(a, 5, 10, -1)
		changeInMorality(a, 5, 10, -1)
	} else {
		changeInStubbornness(a, 5, 1)
		changeInMood(a, 5, 10, 1)
		changeInMorality(a, 5, 10, 1)
	}

	a.Log("I recieved a StateFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor(), "statement": statement})
}

func (a *CustomAgent3) HandleStateHP(msg messages.StateHPMessage) {
	statement := msg.Statement()
	if a.read() {
		if statement > a.HP() {
			changeInStubbornness(a, 5, -1)
			changeInMorality(a, 5, 10, -1)
		} else {
			changeInStubbornness(a, 10, -1)
			changeInMorality(a, 5, 10, 1)
			changeInMood(a, 5, 10, 1)
		}
	}

	a.Log("I recieved a StateHP message from ", infra.Fields{"floor": msg.SenderFloor(), "statement": statement})
}

func (a *CustomAgent3) HandleStateIntendedFoodTaken(msg messages.StateIntendedFoodIntakeMessage) {
	statement := msg.Statement()
	if a.read() {
		changeInStubbornness(a, 5, -1)
		if statement > a.decisions.foodToEat {
			changeInMood(a, 5, 10, -1)
		} else {
			changeInMorality(a, 5, 10, 1)
		}

	}
	a.Log("I recieved a StateIntendedFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor(), "statement": statement})
}

func (a *CustomAgent3) HandleTreatyResponse(msg messages.TreatyResponseMessage) {

	if msg.Response() {
		changeInMood(a, 5, 10, 1)
		changeInMorality(a, 5, 10, 1)
		changeInStubbornness(a, 5, -1)
		//Add friendship level with agent who responded
		//msg.RequestID()
	} else {
		changeInMood(a, 5, 10, -1)
		changeInMorality(a, 5, 10, -1)
		changeInStubbornness(a, 5, 1)
		//Reduce friendship level with agent who responded
		//msg.RequestID()
	}

}

// func (a *CustomAgent3) HandleProposeTreaty(msg messages.ProposeTreatyMessage) {
// 	var reply messages.ResponseMessage
// 	// calculate the benefit of the treaty to the agent - more complex func needed
// 	// possible parameters: stubbornness, mood, food consumed history, number of friends,

// 	if a.knowledge.foodMovingAvg >= float64(a.foodReqCalc(50, 50)) { // satiated
// 		reply = msg.Reply(a.BaseAgent().ID(), a.Floor(), msg.SenderFloor()-a.Floor(), a.read()) // a.read() false "stubbornness" % of the time
// 	} else { // unsatiated
// 		if a.vars.mood > 50 {
// 			reply = msg.Reply(a.BaseAgent().ID(), a.Floor(), msg.SenderFloor()-a.Floor(), a.read()) // a.read() false "stubbornness" % of the time
// 		} else {
// 			reply = msg.Reply(a.BaseAgent().ID(), a.Floor(), msg.SenderFloor()-a.Floor(), false) // a.read() false "stubbornness" % of the time
// 		}
// 	}
// 	a.SendMessage(reply)
// }

type AgentPosition int
type FoodTaken int

const (
	Strong AgentPosition = iota + 1
	Healthy
	Average
	Weak
	SurvivalLevel
)

const (
	VeryLarge FoodTaken = iota + 1
	Large
	Moderate
	Little
	SurvivalAmount
	TooLittle
)

// Returns the AgentPosition (relative strength measure) of the agent when at the minimum HP defined by the condition and conditionOp
func (a *CustomAgent3) requiredHPLevel(treaty messages.Treaty) AgentPosition {
	if treaty.ConditionOp() == messages.LT || treaty.ConditionOp() == messages.LE || treaty.ConditionValue() > 100 {
		return SurvivalLevel
	}
	switch hp := treaty.ConditionValue(); {
	case hp >= 75:
		return Strong
	case hp >= 55:
		return Healthy
	case hp >= 35:
		return Average
	case hp >= a.HealthInfo().WeakLevel:
		return Weak
	default:
		return SurvivalLevel
	}
}

// Determining if a given floor means an agent is in a good / bad position relies on knowledge that the agent has no access to.
// Hence, initial approach is to assume the treaty is risky, and could possibly activate when the agent is at SurvivalLevel.
func (a *CustomAgent3) requiredFloorLevel(treaty messages.Treaty) AgentPosition {
	return SurvivalLevel
}

// Same as Floor, initial approach is to assume the treaty is risky, and could possibly activate when the agent is at SurvivalLevel.
func (a *CustomAgent3) requiredAvailFoodLevel(treaty messages.Treaty) AgentPosition {
	return SurvivalLevel
}

// Calculates food available to eat if request applied to current platform food, and uses this as an estimate for the general case.
func (a *CustomAgent3) reqFoodTakenEstimate(treaty messages.Treaty, percentage bool) FoodTaken {
	var foodToEatCalc int
	if percentage {
		foodToEatCalc = int(int(a.CurrPlatFood()) * (100 - treaty.RequestValue()))
	} else {
		foodToEatCalc = int(int(a.CurrPlatFood()) - treaty.RequestValue())
	}
	switch foodToEat := foodToEatCalc; {
	case foodToEat > a.foodReqCalc(85, 85):
		return VeryLarge
	case foodToEat > a.foodReqCalc(60, 60):
		return Large
	case foodToEat > a.foodReqCalc(40, 40):
		return Moderate
	case foodToEat > a.foodReqCalc(a.HealthInfo().WeakLevel+15, a.HealthInfo().WeakLevel+15):
		return Little
	case foodToEat > a.foodReqCalc(a.HealthInfo().WeakLevel, a.HealthInfo().WeakLevel):
		return SurvivalAmount
	default:
		return TooLittle
	}
}

// 1. requiredAgentPosition evaluates the condition, 2. foodTakenEstimate evaluates the request,
// 3. agentVarsPassed uses agent params with evaluations, 4. Reply sent which accepts/rejects the treaty
func (a *CustomAgent3) HandleProposeTreaty(msg messages.ProposeTreatyMessage) {
	treaty := msg.Treaty()
	var requiredAgentPosition AgentPosition
	var foodTakenEstimate FoodTaken
	var reply messages.ResponseMessage

	switch treaty.Condition() {
	case messages.HP:
		requiredAgentPosition = a.requiredHPLevel(treaty)
	case messages.Floor:
		requiredAgentPosition = a.requiredFloorLevel(treaty)
	case messages.AvailableFood:
		requiredAgentPosition = a.requiredAvailFoodLevel(treaty)
	}

	foodTakenEstimate = a.reqFoodTakenEstimate(treaty, treaty.Request() == messages.LeavePercentFood)

	// Maybe take the duration and signatures into account

	// If agent is in a bad mood, it will only accept treaties that take effect when it is in a strong position.
	// If agent has low morality, it will only accept treaties that involve it taking large amounts of food.
	agentVarsPassed := a.vars.mood > (20*int(requiredAgentPosition)-20) && a.vars.morality > (20*int(foodTakenEstimate)-20) && a.vars.morality < (20*int(foodTakenEstimate)+20)

	//use agent variables, foodTakenEstimate, and requiredAgentPosition to accept/reject
	if treaty.Request() == messages.Inform && treaty.Condition() == messages.HP { // accept HP inform requests
		reply = msg.Reply(a.BaseAgent().ID(), a.Floor(), msg.SenderFloor()-a.Floor(), true)
	} else if treaty.Request() == messages.Inform { // reject Floor and AvailableFood inform requests
		reply = msg.Reply(a.BaseAgent().ID(), a.Floor(), msg.SenderFloor()-a.Floor(), false)
	} else if agentVarsPassed { // For leavePercentFood and leaveAmountFood requests, use agent state and predicted benefit of treaty
		reply = msg.Reply(a.BaseAgent().ID(), a.Floor(), msg.SenderFloor()-a.Floor(), a.read()) // accept acceptable treaties "stubbornness" % of time
	} else {
		reply = msg.Reply(a.BaseAgent().ID(), a.Floor(), msg.SenderFloor()-a.Floor(), false) // reject unaccceptable treaties
	}
	a.SendMessage(reply)

}
