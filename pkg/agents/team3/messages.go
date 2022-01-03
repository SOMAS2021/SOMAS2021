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

func (a *CustomAgent3) updateFriendship(friend uuid.UUID, value float64) {
	friendship, _ := friendshipLevel(a, friend)
	if friendship == 0 {
		addFriend(a, friend)
	} else {
		friendshipChange(a, friend, value)
	}
}

func (a *CustomAgent3) ticklyMessage() {
	r := rand.Intn(5)
	switch r {
	case 0:
		msg := messages.NewAskFoodTakenMessage(a.BaseAgent().ID(), a.Floor(), a.Floor()-1)
		a.SendMessage(msg)
		a.Log("I sent a message", infra.Fields{"message": "AskFoodTaken"})
	case 1:
		msg := messages.NewAskHPMessage(a.BaseAgent().ID(), a.Floor(), a.Floor()-1)
		a.SendMessage(msg)
		a.Log("I sent a message", infra.Fields{"message": "AskHP"})
	case 2:
		msg := messages.NewAskIntendedFoodIntakeMessage(a.BaseAgent().ID(), a.Floor(), a.Floor()-1)
		a.SendMessage(msg)
		a.Log("I sent a message", infra.Fields{"message": "AskIntendedFoodIntake"})
	case 3:
		msg := messages.NewRequestLeaveFoodMessage(a.BaseAgent().ID(), a.Floor(), a.Floor()+1, 10)
		a.SendMessage(msg)
		a.Log("I sent a message", infra.Fields{"message": "RequestLeaveFood"})
	default:
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
	friendship, _ := friendshipLevel(a, msg.SenderID())
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
	friendship, _ := friendshipLevel(a, msg.SenderID())
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
func (a *CustomAgent3) HandleProposeTreaty(msg messages.ProposeTreatyMessage) {
	var reply messages.ResponseMessage
	// calculate the benefit of the treaty to the agent - more complex func needed
	// possible parameters: mood, food consumed history, number of friends,

	if a.knowledge.foodMovingAvg >= food.FoodType(a.foodReqCalc(50, 50)) { // satiated
		reply = msg.Reply(a.BaseAgent().ID(), a.Floor(), msg.SenderFloor()-a.Floor(), a.read()) // a.read() false "stubbornness" % of the time
	} else { // unsatiated
		if a.vars.mood > 50 {
			reply = msg.Reply(a.BaseAgent().ID(), a.Floor(), msg.SenderFloor()-a.Floor(), a.read()) // a.read() false "stubbornness" % of the time
		} else {
			reply = msg.Reply(a.BaseAgent().ID(), a.Floor(), msg.SenderFloor()-a.Floor(), false) // a.read() false "stubbornness" % of the time
		}
	}
	a.SendMessage(reply)
}
