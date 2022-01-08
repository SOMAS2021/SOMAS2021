package team3

import (
	"math/rand"

	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
	"github.com/google/uuid"
)

//Upon receipt of message define affected emotions
// ACK MESSAGE
//If x time passed no message received/acked morale decrease
//Include if ack message same user ID occurs x+1 times, morale increase
//If stubborness = y+1, discard, a.k.a. leave unread

func (a *CustomAgent3) treatyFull() bool {
	return len(a.ActiveTreaties()) > 0
}

func (a *CustomAgent3) treatyPendingResponse() bool {
	return !(a.knowledge.treatyProposed.ID() == uuid.Nil)
}

func (a *CustomAgent3) requestHelpInCrit() {
	if a.treatyFull() || a.treatyPendingResponse() {
		msg := messages.NewRequestLeaveFoodMessage(a.ID(), a.Floor(), a.Floor()+1, a.foodReqCalc(a.HP(), a.HealthInfo().HPReqCToW)) //to higher floor?
		a.SendMessage(msg)
		a.Log("I sent a help message", infra.Fields{"message": "RequestLeaveFood"})
	} else {
		tr := messages.NewTreaty(messages.HP, 20, messages.LeavePercentFood, 95, messages.GT, messages.GT, 3, a.ID()) //generalise later
		a.knowledge.treatyProposed = *tr                                                                              //remember the treaty we proposed
		msg := messages.NewProposalMessage(a.BaseAgent().ID(), a.Floor(), a.Floor()+1, *tr)
		a.SendMessage(msg)
		a.Log("I sent a treaty")
	}
}

func (a *CustomAgent3) askHP(direction int) {
	msg := messages.NewAskHPMessage(a.BaseAgent().ID(), a.Floor(), a.Floor()+direction)
	a.SendMessage(msg)
	a.Log("I sent a message", infra.Fields{"message": "AskHP"})

}

func (a *CustomAgent3) askFoodTaken(direction int) {
	msg := messages.NewAskFoodTakenMessage(a.BaseAgent().ID(), a.Floor(), a.Floor()+direction)
	a.SendMessage(msg)
	a.Log("I sent a message", infra.Fields{"message": "AskFoodTaken"})

}

func (a *CustomAgent3) askTakeFood(HPNeighbour int, direction int) { //direction 1 or -1
	survivalFood := a.foodReqCalc(a.HP(), a.HealthInfo().HPReqCToW)
	stayFood := a.foodReqCalc(a.HP(), a.HP())

	if HPNeighbour < 40 {
		msg := messages.NewRequestTakeFoodMessage(a.BaseAgent().ID(), a.Floor(), a.Floor()+direction, stayFood)
		a.SendMessage(msg)
	} else {
		msg := messages.NewRequestTakeFoodMessage(a.BaseAgent().ID(), a.Floor(), a.Floor()+direction, survivalFood)
		a.SendMessage(msg)
	}
	a.Log("I sent a message", infra.Fields{"message": "RequestTakeFood"})

}

func (a *CustomAgent3) askLeaveFood(direction int) { //direction 1 or -1
	if direction == -1 {
		survivalFood := a.foodReqCalc(a.HP(), a.HealthInfo().HPReqCToW)
		msg := messages.NewRequestLeaveFoodMessage(a.BaseAgent().ID(), a.Floor(), a.Floor()+direction, int(a.knowledge.foodLastSeen-food.FoodType(survivalFood)))
		a.SendMessage(msg)
	} else {
		survivalFood := a.foodReqCalc(a.HP(), a.HealthInfo().HPReqCToW)
		msg := messages.NewRequestLeaveFoodMessage(a.BaseAgent().ID(), a.Floor(), a.Floor()+direction, int(a.knowledge.foodLastSeen+a.knowledge.foodLastEaten+(food.FoodType(survivalFood*4))))
		a.SendMessage(msg)
	}
	a.Log("I sent a message", infra.Fields{"message": "RequestLeaveFood"})
}

func (a *CustomAgent3) proposeTreatiesImmoral() {
	randomFloor := rand.Intn(a.Floor()) + 1
	tr := messages.NewTreaty(messages.Floor, a.Floor()+1, messages.LeavePercentFood, 99, messages.GT, messages.GT, a.knowledge.reshuffleEst/2, a.ID())
	a.knowledge.treatyProposed = *tr //remember the treaty we proposed
	msg := messages.NewProposalMessage(a.BaseAgent().ID(), a.Floor(), randomFloor, *tr)
	a.SendMessage(msg)
	a.Log("I sent a treaty")
}

func (a *CustomAgent3) proposeTreatiesMoral(direction int) { //troll treaties (then add some more b)
	tr := messages.NewTreaty(messages.AvailableFood, 50, messages.LeaveAmountFood, a.foodReqCalc(a.HP(), a.HealthInfo().HPReqCToW), messages.LT, messages.GT, 5, a.ID())
	a.knowledge.treatyProposed = *tr //remember the treaty we proposed
	msg := messages.NewProposalMessage(a.BaseAgent().ID(), a.Floor(), a.Floor()+direction, *tr)
	a.SendMessage(msg)
	a.Log("I sent a treaty")
}

func (a *CustomAgent3) ticklyMessage() {
	if a.HP() == a.HealthInfo().HPCritical {
		a.requestHelpInCrit()
	} else {
		direction := 1
		hpRecorded := a.knowledge.hpAbove
		if a.vars.morality > 50 {
			direction = -1
			hpRecorded = a.knowledge.hpBelow
		}
		if hpRecorded == -1 { //check in case
			hpRecorded = 50
		}
		r := rand.Intn(4)
		switch r {
		case 0:
			a.askHP(direction)
		case 1:
			a.askFoodTaken(direction)
		case 2:
			a.askTakeFood(hpRecorded, direction) //save HP knowledge
		case 3:
			a.askLeaveFood(direction)
		case 4:
			if a.treatyFull() || a.treatyPendingResponse() {
				if a.vars.morality < 10 {
					a.proposeTreatiesImmoral()
				} else {
					a.proposeTreatiesMoral(1)
				}
			} else {
				randomFloor := rand.Intn(a.Floor()-1) + 1
				if a.treatyFull() && !a.treatyPendingResponse() {
					for _, tr := range a.ActiveTreaties() {
						a.knowledge.treatyProposed = tr //remember the treaty we proposed
						msg := messages.NewProposalMessage(a.BaseAgent().ID(), a.Floor(), randomFloor, tr)
						a.SendMessage(msg)
						a.Log("I sent a treaty")
						break
					}
				}
			}
		}
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

//HandleAskHP handles HP messages
func (a *CustomAgent3) HandleAskHP(msg messages.AskHPMessage) { //how are you type question
	a.Log("I recieved an askHP message from ", infra.Fields{"floor": msg.SenderFloor()})
	friendship := a.knowledge.friends[msg.SenderID()]
	if a.read() {
		if a.HP() < a.knowledge.lastHP {
			if friendship < 0.5 {
				a.changeInMood(1, 3, -1)
				a.changeInMorality(1, 3, -1)
			} else {
				a.changeInMood(1, 6, 1)
				a.updateFriendship(msg.SenderID(), 1)
			}
		} else {
			if friendship < 0.5 {
				a.changeInMood(1, 3, 1)
				a.changeInMorality(1, 6, 1)
				a.changeInStubbornness(5, -1)
				a.updateFriendship(msg.SenderID(), 1)
			} else {
				a.changeInMood(1, 6, 1)
				a.changeInStubbornness(5, -1)
				a.updateFriendship(msg.SenderID(), 1)
			}
		}
		reply := msg.Reply(a.BaseAgent().ID(), a.Floor(), msg.SenderFloor(), a.HP())
		a.SendMessage(reply)
		a.Log("I recieved an askHP message from ", infra.Fields{"floor": msg.SenderFloor()})
	}
}

//HandleAskFoodTaken handles asking for Food Taken
func (a *CustomAgent3) HandleAskFoodTaken(msg messages.AskFoodTakenMessage) {
	a.Log("I recieved an askFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor()})
	friendship := a.knowledge.friends[msg.SenderID()]
	if a.read() {
		if a.HP() < a.knowledge.lastHP {
			if friendship < 0.5 {
				a.changeInMood(1, 3, -1)
				a.changeInMorality(1, 3, -1)
			} else {
				a.changeInMood(1, 6, 1)
			}
		} else {
			if friendship < 0.5 {
				a.updateFriendship(msg.SenderID(), -1)
			}
		}
		reply := msg.Reply(a.BaseAgent().ID(), a.Floor(), msg.SenderFloor(), int(a.knowledge.foodLastEaten))
		a.SendMessage(reply)
		a.Log("I sent a replyFoodTaken message to ", infra.Fields{"floor": msg.SenderFloor()})
	}
}

//HandleAskIntendedFoodTaken handles asking for intended food taken
func (a *CustomAgent3) HandleAskIntendedFoodTaken(msg messages.AskIntendedFoodIntakeMessage) {
	if a.read() {
		reply := msg.Reply(a.BaseAgent().ID(), a.Floor(), msg.SenderFloor(), a.decisions.foodToEat)
		a.SendMessage(reply)
		a.Log("I recieved an askIntendedFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor()})
	}
}

//HandleRequestLeaveFood handles asking for intended food left
func (a *CustomAgent3) HandleRequestLeaveFood(msg messages.RequestLeaveFoodMessage) {
	friendship := a.knowledge.friends[msg.SenderID()]
	request := msg.Request()
	percentageDec := 0.8
	if a.read() {
		if a.HP() < a.knowledge.lastHP {
			if friendship < 0.5 {
				a.changeInMood(1, 6, 1)
				a.changeInMorality(1, 3, -1)
				a.changeInStubbornness(5, 1)
				a.updateFriendship(msg.SenderID(), -1)
			} else {
				a.changeInMood(1, 6, -1)
				a.changeInMorality(1, 6, -1)
				a.changeInStubbornness(5, 1)
				a.updateFriendship(msg.SenderID(), -1)
			}
		} else {
			if friendship < 0.5 {
				a.changeInMood(1, 9, 1)
				a.changeInMorality(1, 6, 1)
				a.changeInStubbornness(5, -1)
				a.updateFriendship(msg.SenderID(), 1)
			} else {
				a.changeInMorality(1, 6, 1)
				a.changeInStubbornness(5, -1)
			}
		}
		if request > int(a.knowledge.foodLastSeen-a.knowledge.foodLastEaten) {
			reply := msg.Reply(a.BaseAgent().ID(), a.Floor(), msg.SenderFloor(), false)
			a.SendMessage(reply)
			if friendship > 0.5 {
				if a.vars.morality > 50 {
					if a.vars.mood > 30 {
						a.decisions.foodToEat = int(float64(a.decisions.foodToEat) * percentageDec)
					}
				}
			} else {
				if a.vars.morality > 70 {
					if a.vars.mood > 50 {
						a.decisions.foodToEat = int(float64(a.decisions.foodToEat) * percentageDec)
					}
				}
			}

		} else {
			if friendship > 0.5 {
				if a.vars.morality > 50 {
					if a.vars.mood > 30 {
						reply := msg.Reply(a.BaseAgent().ID(), a.Floor(), msg.SenderFloor(), true)
						a.SendMessage(reply)
					} else {
						reply := msg.Reply(a.BaseAgent().ID(), a.Floor(), msg.SenderFloor(), false)
						a.SendMessage(reply)
					}
				} else {
					reply := msg.Reply(a.BaseAgent().ID(), a.Floor(), msg.SenderFloor(), false)
					a.SendMessage(reply)
				}
			} else {
				if a.vars.morality > 70 {
					if a.vars.mood > 50 {
						reply := msg.Reply(a.BaseAgent().ID(), a.Floor(), msg.SenderFloor(), true)
						a.SendMessage(reply)
					} else {
						reply := msg.Reply(a.BaseAgent().ID(), a.Floor(), msg.SenderFloor(), false)
						a.SendMessage(reply)
					}
				} else {
					reply := msg.Reply(a.BaseAgent().ID(), a.Floor(), msg.SenderFloor(), false)
					a.SendMessage(reply)
				}
			}

		}

		reply := msg.Reply(a.BaseAgent().ID(), a.Floor(), msg.SenderFloor(), true)
		a.SendMessage(reply)
		a.Log("I recieved a requestLeaveFood message from ", infra.Fields{"floor": msg.SenderFloor()})
	}
}

//HandleRequestTakeFood handles asking for request food
func (a *CustomAgent3) HandleRequestTakeFood(msg messages.RequestTakeFoodMessage) {
	friendship := a.knowledge.friends[msg.SenderID()]
	request := msg.Request()
	if a.read() {
		if a.HP() < a.knowledge.lastHP {
			if friendship < 0.5 {
				a.changeInMood(1, 6, -1)
				a.changeInMorality(1, 3, -1)
				a.changeInStubbornness(5, 1)
				a.updateFriendship(msg.SenderID(), -1)
			} else {
				a.changeInMood(1, 6, -1)
				a.changeInMorality(1, 6, -1)
				a.changeInStubbornness(5, 1)
				a.updateFriendship(msg.SenderID(), -1)
			}
		} else {
			if friendship > 0.5 {
				a.changeInMood(1, 6, -1)
				a.changeInMorality(1, 6, 1)
				a.updateFriendship(msg.SenderID(), -1)
			}
		}
		if float64(request) > a.knowledge.foodMovingAvg {
			a.decisions.foodToEat = request
			reply := msg.Reply(a.BaseAgent().ID(), a.Floor(), msg.SenderFloor(), true)
			a.SendMessage(reply)
		} else {
			if friendship > 0.5 {
				if a.vars.morality > 50 {
					if a.vars.mood > 30 {
						a.decisions.foodToEat = request
						reply := msg.Reply(a.BaseAgent().ID(), a.Floor(), msg.SenderFloor(), true)
						a.SendMessage(reply)
					} else {
						reply := msg.Reply(a.BaseAgent().ID(), a.Floor(), msg.SenderFloor(), false)
						a.SendMessage(reply)
					}
				} else {
					reply := msg.Reply(a.BaseAgent().ID(), a.Floor(), msg.SenderFloor(), false)
					a.SendMessage(reply)
				}
			} else {
				if a.vars.morality > 70 {
					if a.vars.mood > 50 {
						a.decisions.foodToEat = request
						reply := msg.Reply(a.BaseAgent().ID(), a.Floor(), msg.SenderFloor(), true)
						a.SendMessage(reply)
					} else {
						reply := msg.Reply(a.BaseAgent().ID(), a.Floor(), msg.SenderFloor(), false)
						a.SendMessage(reply)
					}
				} else {
					reply := msg.Reply(a.BaseAgent().ID(), a.Floor(), msg.SenderFloor(), false)
					a.SendMessage(reply)
				}
			}
		}

		a.Log("I recieved a requestTakeFood message from ", infra.Fields{"floor": msg.SenderFloor()})
	}
}
func (a *CustomAgent3) feelingModifResponse(response bool, friendship float64, sender uuid.UUID) {
	direction := -1
	if response {
		direction = 1
	}
	if friendship > 0.5 {
		a.changeInMood(6, 12, direction)
		a.changeInMorality(6, 12, direction)
		a.changeInStubbornness(5, -1*direction)
		a.updateFriendship(sender, direction)
		a.updateFriendship(sender, direction)
	} else {
		a.changeInMood(1, 6, direction)
		a.changeInMorality(1, 6, direction)
		a.changeInStubbornness(5, -1*direction)
		a.updateFriendship(sender, direction)
	}
}

//HandleResponse handles responses
func (a *CustomAgent3) HandleResponse(msg messages.BoolResponseMessage) {
	a.feelingModifResponse(msg.Response(), a.knowledge.friends[msg.SenderID()], msg.SenderID())
	a.Log("I recieved a Response message from ", infra.Fields{"floor": msg.SenderFloor(), "response": msg.Response()})
}

func (a *CustomAgent3) feelingModifFoodTaken(statement int, friendship float64, sender uuid.UUID) {
	percentageDec := 0.9
	percentageInc := 1.3
	if friendship < 0.5 {
		if statement > a.decisions.foodToEat {
			a.changeInMood(1, 6, -1)
			a.changeInMorality(1, 6, -1)
			a.changeInStubbornness(5, 1)
			a.updateFriendship(sender, -1)
			a.decisions.foodToEat = statement
		} else {
			a.changeInMood(1, 6, 1)
			a.changeInMorality(1, 6, 1)
			a.changeInStubbornness(5, -1)
			a.updateFriendship(sender, 1)
			a.decisions.foodToEat = int(float64(a.decisions.foodToEat) * percentageDec)
		}
	} else {
		if statement > a.decisions.foodToEat {
			a.changeInMorality(1, 3, -1)
			a.decisions.foodToEat = int(float64(a.decisions.foodToEat) * percentageInc)
		} else {
			a.changeInMorality(1, 6, 1)
			a.updateFriendship(sender, 1)
			a.decisions.foodToEat = int(float64(a.decisions.foodToEat) * percentageDec)
		}
	}
}

//HandleStateFoodTaken handles a food taken statement
func (a *CustomAgent3) HandleStateFoodTaken(msg messages.StateFoodTakenMessage) {
	a.feelingModifFoodTaken(msg.Statement(), a.knowledge.friends[msg.SenderID()], msg.SenderID())
	a.Log("I recieved a StateFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor(), "statement": msg.Statement()})
}

func (a *CustomAgent3) feelingModifStateHP(statement int, friendship float64) {
	direction := 1
	if friendship < 0.5 {
		if statement > a.HP() {
			direction = -1
			a.changeInMorality(1, 6, -1)
		}
		a.changeInMood(1, 6, direction)

	} else {
		if statement < a.HP() {
			direction = -1
		}
		a.changeInStubbornness(5, direction)
		a.changeInMood(1, 3, direction)

	}
}

//HandleStateHP handles an HP statement
func (a *CustomAgent3) HandleStateHP(msg messages.StateHPMessage) {
	if msg.SenderFloor() > a.Floor() {
		a.knowledge.hpBelow = msg.Statement()
	} else {
		a.knowledge.hpAbove = msg.Statement()
	}
	a.feelingModifStateHP(msg.Statement(), a.knowledge.friends[msg.SenderID()])
	a.Log("I recieved a StateHP message from ", infra.Fields{"floor": msg.SenderFloor(), "statement": msg.Statement()})
}

func (a *CustomAgent3) feelingModifIntendedFoodTaken(statement int, friendship float64, sender uuid.UUID) {
	percentageDec := 0.9
	percentageInc := 1.3
	if friendship < 0.5 {
		if statement > a.decisions.foodToEat {
			a.changeInStubbornness(5, 1)
			a.updateFriendship(sender, -1)
			a.decisions.foodToEat = statement
		} else {
			a.changeInMorality(1, 6, 1)
			a.changeInStubbornness(5, -1)
			a.updateFriendship(sender, 1)
			a.decisions.foodToEat = int(float64(a.decisions.foodToEat) * percentageDec)

		}
	} else {
		if statement > a.decisions.foodToEat { //could incorporate a max function
			a.changeInMorality(1, 3, -1)
			a.changeInStubbornness(5, 1)
			a.updateFriendship(sender, -1)
			a.decisions.foodToEat = int(float64(a.decisions.foodToEat) * percentageInc)
		} else {
			a.changeInStubbornness(5, -1)
			a.decisions.foodToEat = int(float64(a.decisions.foodToEat) * percentageDec)
		}
	}
}

//HandleStateIntendedFoodTaken handles an intender food taken statement
func (a *CustomAgent3) HandleStateIntendedFoodTaken(msg messages.StateIntendedFoodIntakeMessage) {
	a.feelingModifIntendedFoodTaken(msg.Statement(), a.knowledge.friends[msg.SenderID()], msg.SenderID())
	a.Log("I recieved a StateIntendedFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor(), "statement": msg.Statement()})
}

//HandleTreatyResponse handles the response of a treaty
func (a *CustomAgent3) HandleTreatyResponse(msg messages.TreatyResponseMessage) {
	if msg.Response() {
		if a.knowledge.treatyProposed.ID() != uuid.Nil { //check in case something went wrong.
			treaty := a.knowledge.treatyProposed //signed our proposed treatry
			treaty.SignTreaty()
			a.AddTreaty(treaty)
		}
	}
	a.knowledge.treatyProposed = *messages.NewTreaty(messages.HP, 0, messages.LeaveAmountFood, 0, messages.GT, messages.GT, 0, uuid.Nil) //restart the sent treaty.
}

//AgentPosition signifies how strong it needs to be to approve a treaty.
type AgentPosition int

//FoodTaken signifies how much food we need to take to approve a treaty.
type FoodTaken int

const (
	Strong AgentPosition = iota + 1
	Healthy
	Average
	Weak
	SurvivalLevel
	Reject
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
	case hp == a.HealthInfo().HPCritical:
		return SurvivalLevel
	default:
		return Reject
	}
}

// Determining if a given floor means an agent is in a good / bad position relies on knowledge that the agent has no access to.
// Hence, initial approach is to assume the treaty is risky, and could possibly activate when the agent is at SurvivalLevel.
func (a *CustomAgent3) requiredFloorLevel(treaty messages.Treaty) AgentPosition {
	return Reject
}

// Same as Floor, initial approach is to assume the treaty is risky, and could possibly activate when the agent is at SurvivalLevel.
func (a *CustomAgent3) requiredAvailFoodLevel(treaty messages.Treaty) AgentPosition {
	return Reject
}

// Calculates food available to eat if request applied to current platform food, and uses this as an estimate for the general case.
func (a *CustomAgent3) reqFoodTakenEstimate(treaty messages.Treaty, percentage bool) FoodTaken {
	var foodToEatCalc int

	if treaty.RequestOp() == messages.LT || treaty.RequestOp() == messages.LE || treaty.RequestOp() == messages.EQ {
		return VeryLarge
	}

	if percentage {
		foodToEatCalc = int(float64(a.CurrPlatFood()) * float64((100.0-float64(treaty.RequestValue()))/100.0))
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

//HandleProposeTreaty 1. requiredAgentPosition evaluates the condition, 2. foodTakenEstimate evaluates the request,
// 3. agentVarsPassed uses agent params with evaluations, 4. Reply sent which accepts/rejects the treaty
func (a *CustomAgent3) HandleProposeTreaty(msg messages.ProposeTreatyMessage) {
	if a.treatyFull() {
		reply := msg.Reply(a.BaseAgent().ID(), a.Floor(), msg.SenderFloor(), false)
		a.SendMessage(reply)
	} else {

		treaty := msg.Treaty()
		var minActivationLevel AgentPosition
		var response bool

		switch treaty.Condition() {
		case messages.HP:
			minActivationLevel = a.requiredHPLevel(treaty)
		case messages.Floor:
			minActivationLevel = a.requiredFloorLevel(treaty)
		case messages.AvailableFood:
			minActivationLevel = a.requiredAvailFoodLevel(treaty)
		}

		foodTakenEstimate := a.reqFoodTakenEstimate(treaty, treaty.Request() == messages.LeavePercentFood)

		// If agent is in a bad mood, it will only accept treaties that take effect when it is in a strong position.
		// If agent has low morality, it will only accept treaties that involve it taking large amounts of food.
		agentVarsPassed := a.vars.mood > (20*int(minActivationLevel)-20) && a.vars.morality > (20*int(foodTakenEstimate)-20) && a.vars.morality < (20*int(foodTakenEstimate)+20)
		// Check duration is not too long, and use stubbornness to decide if an agent gives in and accepts treaties with at least 5 signatures
		allChecksPassed := agentVarsPassed && treaty.Duration() < 2*a.knowledge.reshuffleEst
		if treaty.SignatureCount() >= 5 && !allChecksPassed {
			allChecksPassed = a.read()
		}

		//use agent variables, foodTakenEstimate, and requiredAgentPosition to accept/reject
		if allChecksPassed && treaty.Request() != messages.Inform { // dont accept HP inform requests
			response = true
			treaty.SignTreaty()
			a.AddTreaty(treaty)
		} else { // reject other treaties
			response = false
		}
		reply := msg.Reply(a.BaseAgent().ID(), a.Floor(), msg.SenderFloor(), response)
		a.SendMessage(reply)
	}
}
