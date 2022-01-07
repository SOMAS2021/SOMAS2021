package team3

import (
	"math"
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
func (a *CustomAgent3) proposeTreatiesInmoral() {
	randomFloor := rand.Intn(a.Floor()) + 1
	tr := messages.NewTreaty(messages.Floor, a.Floor()+1, messages.LeavePercentFood, 99, messages.GT, messages.GT, 9, a.ID())
	a.knowledge.treatyProposed = *tr //remember the treaty we proposed
	msg := messages.NewProposalMessage(a.BaseAgent().ID(), a.Floor(), randomFloor, *tr)
	a.SendMessage(msg)
	a.Log("I sent a treaty")
}
func (a *CustomAgent3) proposeTreatiesMoral(direction int) { //troll treaties (then add some more b)
	tr := messages.NewTreaty(messages.HP, 20, messages.LeavePercentFood, 95, messages.GT, messages.GT, 3, a.ID())
	r := rand.Intn(3)
	switch r {
	case 0:
		tr = messages.NewTreaty(messages.HP, 20, messages.LeaveAmountFood, 20, messages.GT, messages.GT, 15, a.ID())
	case 1:
		tr = messages.NewTreaty(messages.HP, 60, messages.LeavePercentFood, 60, messages.GT, messages.GT, 5, a.ID())
	case 2:
		tr = messages.NewTreaty(messages.HP, 10, messages.LeavePercentFood, 95, messages.GT, messages.GT, 10, a.ID())
	case 3:
		tr = messages.NewTreaty(messages.AvailableFood, 50, messages.LeaveAmountFood, a.foodReqCalc(a.HP(), a.HealthInfo().HPReqCToW), messages.LT, messages.GT, 5, a.ID())
	}
	//generalise later
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
		if a.vars.morality > 50 {
			direction = -1
		}
		r := rand.Intn(4)
		switch r {
		case 0:
			a.askHP(direction)
		case 1:
			a.askFoodTaken(direction)
		case 2:
			a.askTakeFood(50, direction) //save HP knowledge
		case 3:
			a.askLeaveFood(direction)
		case 4:
			if a.treatyFull() || a.treatyPendingResponse() {
				if a.vars.morality < 10 {
					a.proposeTreatiesInmoral()
				} else {
					a.proposeTreatiesMoral(1)
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

func max(x, y, z int) int {
	temp := math.Max(float64(x), float64(y))
	return int(math.Max(float64(temp), float64(z)))
}

func min(x, y, z int) int {
	temp := math.Min(float64(x), float64(y))
	return int(math.Min(float64(temp), float64(z)))

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

func (a *CustomAgent3) HandleAskHP(msg messages.AskHPMessage) { //how are you type question
	a.Log("I recieved an askHP message from ", infra.Fields{"floor": msg.SenderFloor()})
	friendship := a.knowledge.friends[msg.SenderID()]
	if a.read() {
		if a.HP() < a.knowledge.lastHP {
			if friendship < 0.5 {
				//changeInStubbornness(a, 5, 1)
				changeInMood(a, 1, 3, -1)
				changeInMorality(a, 1, 3, -1)
			} else {
				//changeInStubbornness(a, 5, 1)
				//changeInMorality(a, 1, 3, -1)
				changeInMood(a, 1, 6, 1)
				a.updateFriendship(msg.SenderID(), 1)
			}

		} else {
			if friendship < 0.5 {
				changeInMood(a, 1, 3, 1)
				changeInMorality(a, 1, 6, 1)
				changeInStubbornness(a, 5, -1)
				a.updateFriendship(msg.SenderID(), 1)
			} else {
				changeInMood(a, 1, 6, 1)
				//changeInMorality(a, 1, 6, 1)
				changeInStubbornness(a, 5, -1)
				a.updateFriendship(msg.SenderID(), 1)
			}
		}

		reply := msg.Reply(a.BaseAgent().ID(), a.Floor(), msg.SenderFloor(), a.HP())
		a.SendMessage(reply)
		a.Log("I recieved an askHP message from ", infra.Fields{"floor": msg.SenderFloor()})
	}
}

func (a *CustomAgent3) HandleAskFoodTaken(msg messages.AskFoodTakenMessage) {
	a.Log("I recieved an askFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor()})
	friendship := a.knowledge.friends[msg.SenderID()]
	if a.read() {
		if a.HP() < a.knowledge.lastHP {
			if friendship < 0.5 {
				changeInMood(a, 1, 3, -1)
				changeInMorality(a, 1, 3, -1)
				//changeInStubbornness(a, 5, -1)
				//a.updateFriendship(msg.SenderID(), 1)
			} else {
				changeInMood(a, 1, 6, 1)
				//changeInMorality(a, 1, 6, 1)
				//changeInStubbornness(a, 5, -1)
				//a.updateFriendship(msg.SenderID(), 1)
			}
		} else {
			if friendship < 0.5 {
				//changeInMood(a, 1, 3, 1)
				//changeInMorality(a, 1, 6, 1)
				//changeInStubbornness(a, 5, -1)
				a.updateFriendship(msg.SenderID(), -1)
			} else {
				//changeInMood(a, 1, 3, 1)
				//changeInMorality(a, 1, 6, 1)
				//changeInStubbornness(a, 5, -1)
				//a.updateFriendship(msg.SenderID(), 1)
			}
		}
		reply := msg.Reply(a.BaseAgent().ID(), a.Floor(), msg.SenderFloor(), int(a.knowledge.foodLastEaten))
		a.SendMessage(reply)
		a.Log("I sent a replyFoodTaken message to ", infra.Fields{"floor": msg.SenderFloor()})
	}
}

func (a *CustomAgent3) HandleAskIntendedFoodTaken(msg messages.AskIntendedFoodIntakeMessage) {
	//friendship := a.knowledge.friends[msg.SenderID()]
	if a.read() {
		//if a.HP() < a.knowledge.lastHP {
		//if friendship < 0.5 {
		//changeInMood(a, 1, 3, 1)
		//changeInMorality(a, 1, 6, 1)
		//changeInStubbornness(a, 5, -1)
		//a.updateFriendship(msg.SenderID(), 1)
		//}
		//} else {
		//if friendship < 0.5 {
		//changeInMood(a, 1, 3, 1)
		//changeInMorality(a, 1, 6, 1)
		//changeInStubbornness(a, 5, -1)
		//a.updateFriendship(msg.SenderID(), 1)
		//} else {
		//changeInMood(a, 1, 3, 1)
		//changeInMorality(a, 1, 6, 1)
		//changeInStubbornness(a, 5, -1)
		//a.updateFriendship(msg.SenderID(), 1)
		//}
		//}
		//add critical state effect
		reply := msg.Reply(a.BaseAgent().ID(), a.Floor(), msg.SenderFloor(), a.decisions.foodToEat)
		a.SendMessage(reply)
		a.Log("I recieved an askIntendedFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor()})
	}
}

func (a *CustomAgent3) HandleRequestLeaveFood(msg messages.RequestLeaveFoodMessage) {
	friendship := a.knowledge.friends[msg.SenderID()]
	request := msg.Request()
	percentageDec := 0.8
	if a.read() {
		if a.HP() < a.knowledge.lastHP {
			if friendship < 0.5 {
				changeInMood(a, 1, 6, 1)
				changeInMorality(a, 1, 3, -1)
				changeInStubbornness(a, 5, 1)
				a.updateFriendship(msg.SenderID(), -1)
			} else {
				changeInMood(a, 1, 6, -1)
				changeInMorality(a, 1, 6, -1)
				changeInStubbornness(a, 5, 1)
				a.updateFriendship(msg.SenderID(), -1)
			}
		} else {
			if friendship < 0.5 {
				changeInMood(a, 1, 9, 1)
				changeInMorality(a, 1, 6, 1)
				changeInStubbornness(a, 5, -1)
				a.updateFriendship(msg.SenderID(), 1)
			} else {
				//changeInMood(a, 1, 3, 1)
				changeInMorality(a, 1, 6, 1)
				changeInStubbornness(a, 5, -1)
				//a.updateFriendship(msg.SenderID(), 1)
			}
		}
		if request > int(a.knowledge.foodLastSeen-a.knowledge.foodLastEaten) {
			reply := msg.Reply(a.BaseAgent().ID(), a.Floor(), msg.SenderFloor(), false)
			a.SendMessage(reply)
			if a.HP() > a.knowledge.lastHP {
				if friendship > 0.5 {
					if a.vars.morality > 50 {
						if a.vars.mood > 30 {
							a.decisions.foodToEat = a.decisions.foodToEat * percentageDec
						}
					}
				} else {
					if a.vars.morality > 70 {
						if a.vars.mood > 50 {
							a.decisions.foodToEat = a.decisions.foodToEat * percentageDec
						}
					}
				}
			}
		} else {
			if a.HP() > a.knowledge.lastHP {
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
			} else {
				reply := msg.Reply(a.BaseAgent().ID(), a.Floor(), msg.SenderFloor(), false)
				a.SendMessage(reply)
			}
		}

		reply := msg.Reply(a.BaseAgent().ID(), a.Floor(), msg.SenderFloor(), true)
		a.SendMessage(reply)
		a.Log("I recieved a requestLeaveFood message from ", infra.Fields{"floor": msg.SenderFloor()})
	}
}

func (a *CustomAgent3) HandleRequestTakeFood(msg messages.RequestTakeFoodMessage) {
	friendship := a.knowledge.friends[msg.SenderID()]
	request := msg.Request()
	if a.read() {
		if a.HP() < a.knowledge.lastHP {
			if friendship < 0.5 {
				changeInMood(a, 1, 6, -1)
				changeInMorality(a, 1, 3, -1)
				changeInStubbornness(a, 5, 1)
				a.updateFriendship(msg.SenderID(), -1)
			} else {
				changeInMood(a, 1, 6, -1)
				changeInMorality(a, 1, 6, -1)
				changeInStubbornness(a, 5, 1)
				a.updateFriendship(msg.SenderID(), -1)
			}
		} else {
			if friendship < 0.5 {
				//changeInMood(a, 1, 3, 1)
				//changeInMorality(a, 1, 6, 1)
				//changeInStubbornness(a, 5, -1)
				//a.updateFriendship(msg.SenderID(), 1)
			} else {
				changeInMood(a, 1, 6, -1)
				changeInMorality(a, 1, 6, 1)
				//changeInStubbornness(a, 5, -1)
				a.updateFriendship(msg.SenderID(), -1)
			}
		}
		if float64(request) > a.knowledge.foodMovingAvg {
			a.decisions.foodToEat = request
			reply := msg.Reply(a.BaseAgent().ID(), a.Floor(), msg.SenderFloor(), true)
			a.SendMessage(reply)
		} else {
			if a.HP() > a.knowledge.lastHP {
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
			} else {
				reply := msg.Reply(a.BaseAgent().ID(), a.Floor(), msg.SenderFloor(), false)
				a.SendMessage(reply)
			}
		}

		a.Log("I recieved a requestTakeFood message from ", infra.Fields{"floor": msg.SenderFloor()})
	}
}

func (a *CustomAgent3) HandleResponse(msg messages.BoolResponseMessage) {
	response := msg.Response()
	a.Log("I recieved a Response message from ", infra.Fields{"floor": msg.SenderFloor(), "response": response})
}

func (a *CustomAgent3) HandleStateFoodTaken(msg messages.StateFoodTakenMessage) {
	statement := msg.Statement()
	friendship := a.knowledge.friends[msg.SenderID()]
	if friendship < 0.5 {
		if statement > a.decisions.foodToEat {
			changeInMood(a, 1, 6, -1)
			changeInMorality(a, 1, 6, -1)
			changeInStubbornness(a, 5, 1)
			a.updateFriendship(msg.SenderID(), -1)
		} else {
			changeInMood(a, 1, 6, 1)
			changeInMorality(a, 1, 6, 1)
			changeInStubbornness(a, 5, -1)
			a.updateFriendship(msg.SenderID(), 1)
		}
	} else {
		if statement > a.decisions.foodToEat {
			//changeInMood(a, 1, 3, 1)
			changeInMorality(a, 1, 3, -1)
			//changeInStubbornness(a, 5, -1)
			//a.updateFriendship(msg.SenderID(), 1)
		} else {
			//changeInMood(a, 1, 3, 1)
			changeInMorality(a, 1, 6, 1)
			//changeInStubbornness(a, 5, -1)
			a.updateFriendship(msg.SenderID(), 1)
		}
	}

	a.Log("I recieved a StateFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor(), "statement": statement})
}

func (a *CustomAgent3) HandleStateHP(msg messages.StateHPMessage) {
	statement := msg.Statement()
	friendship := a.knowledge.friends[msg.SenderID()]
	if friendship < 0.5 {
		if statement > a.decisions.foodToEat {
			changeInMood(a, 1, 6, -1)
			changeInMorality(a, 1, 6, -1)
			//changeInStubbornness(a, 5, -1)
			//a.updateFriendship(msg.SenderID(), 1)
		} else {
			changeInMood(a, 1, 6, 1)
			//changeInMorality(a, 1, 6, 1)
			//changeInStubbornness(a, 5, -1)
			//a.updateFriendship(msg.SenderID(), 1)
		}
	} else {
		if statement > a.decisions.foodToEat {
			changeInStubbornness(a, 5, 1)
			changeInMood(a, 1, 3, 1)
		} else {
			changeInStubbornness(a, 5, -1)
			changeInMood(a, 1, 3, -1)
		}
	}

	a.Log("I recieved a StateHP message from ", infra.Fields{"floor": msg.SenderFloor(), "statement": statement})
}

func (a *CustomAgent3) HandleStateIntendedFoodTaken(msg messages.StateIntendedFoodIntakeMessage) {
	statement := msg.Statement()
	friendship := a.knowledge.friends[msg.SenderID()]
	if friendship < 0.5 {
		if statement > a.decisions.foodToEat {
			//changeInMood(a, 1, 3, 1)
			//changeInMorality(a, 1, 6, 1)
			changeInStubbornness(a, 5, 1)
			a.updateFriendship(msg.SenderID(), -1)
		} else {
			//changeInMood(a, 1, 3, 1)
			changeInMorality(a, 1, 6, 1)
			changeInStubbornness(a, 5, -1)
			a.updateFriendship(msg.SenderID(), 1)
		}
	} else {
		if statement > a.decisions.foodToEat {
			//changeInMood(a, 1, 3, 1)
			changeInMorality(a, 1, 3, -1)
			changeInStubbornness(a, 5, 1)
			a.updateFriendship(msg.SenderID(), -1)
		} else {
			//changeInMood(a, 1, 3, 1)
			//changeInMorality(a, 1, 6, 1)
			changeInStubbornness(a, 5, -1)
			//a.updateFriendship(msg.SenderID(), 1)
		}
	}
	a.Log("I recieved a StateIntendedFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor(), "statement": statement})
}

func (a *CustomAgent3) HandleTreatyResponse(msg messages.TreatyResponseMessage) {

	if msg.Response() {
		changeInMood(a, 5, 10, 1)
		changeInMorality(a, 5, 10, 1)
		changeInStubbornness(a, 5, -1)
		if a.knowledge.treatyProposed.ID() != uuid.Nil { //check in case something went wrong.
			treaty := a.knowledge.treatyProposed //signed our proposed treatry
			treaty.SignTreaty()
			a.AddTreaty(treaty)
		}

		//Add friendship level with agent who responded
		//msg.RequestID()
	} else {
		changeInMood(a, 5, 10, -1)
		changeInMorality(a, 5, 10, -1)
		changeInStubbornness(a, 5, 1)
		//Reduce friendship level with agent who responded
		//msg.RequestID()
	}
	a.knowledge.treatyProposed = *messages.NewTreaty(messages.HP, 0, messages.LeaveAmountFood, 0, messages.GT, messages.GT, 0, uuid.Nil) //restart the sent treaty.
}

type AgentPosition int
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

// 1. requiredAgentPosition evaluates the condition, 2. foodTakenEstimate evaluates the request,
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
