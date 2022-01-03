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
func ticklyMessage(a *CustomAgent3) {
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
func message(a *CustomAgent3) {
	receivedMsg := a.ReceiveMessage()
	if receivedMsg != nil {
		a.Log("Custom agent 3 each run:", infra.Fields{"floor": a.Floor(), "hp": a.HP(), "Mood": a.vars.mood, "Morality": a.vars.morality})
		receivedMsg.Visit(a)
		a.Log("Custom agent 3 each run:", infra.Fields{"floor": a.Floor(), "hp": a.HP(), "Mood": a.vars.mood, "Morality": a.vars.morality})
	} else {
		ticklyMessage(a)
		a.Log("I got nothing")
	}

}

func (a *CustomAgent3) HandleAskHP(msg messages.AskHPMessage) {
	a.Log("I recieved an askHP message from ", infra.Fields{"floor": msg.SenderFloor()})
	if read(a) {
		a.vars.stubbornness = a.vars.stubbornness - 5 //value could be different
		reply := msg.Reply(a.BaseAgent().ID(), a.Floor(), msg.SenderFloor()-a.Floor(), a.HP())
		a.SendMessage(reply)
		a.Log("I recieved an askHP message from ", infra.Fields{"floor": msg.SenderFloor()})
	}
}

func (a *CustomAgent3) HandleAskFoodTaken(msg messages.AskFoodTakenMessage) {
	a.Log("I recieved an askFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor()})
	if read(a) {
		if a.HP() < a.knowledge.lastHP {
			a.vars.stubbornness = a.vars.stubbornness + 5
			//addfriend(a, ) need id
			//if a.vars.morality < 30 {
			//can we reject this message or send a response of false?
			//}
			changeInMood(a, 5, 10, 1)
		}
		reply := msg.Reply(a.BaseAgent().ID(), a.Floor(), msg.SenderFloor()-a.Floor(), int(a.knowledge.foodLastEaten))
		a.SendMessage(reply)
		a.Log("I sent a replyFoodTaken message to ", infra.Fields{"floor": msg.SenderFloor()})
	}
}

func (a *CustomAgent3) HandleAskIntendedFoodTaken(msg messages.AskIntendedFoodIntakeMessage) {
	if read(a) {
		a.vars.stubbornness = a.vars.stubbornness + 2
		//if a.vars.morality < 30 {

		//}
		reply := msg.Reply(a.BaseAgent().ID(), a.Floor(), msg.SenderFloor()-a.Floor(), a.decisions.foodToEat)
		a.SendMessage(reply)
		a.Log("I recieved an askIntendedFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor()})
	}
}

func (a *CustomAgent3) HandleRequestLeaveFood(msg messages.RequestLeaveFoodMessage) {
	if read(a) {
		if a.HP() < a.knowledge.lastHP {
			if a.vars.stubbornness > 80 {
				a.vars.stubbornness = a.vars.stubbornness + 5
			} else {
				a.vars.stubbornness = a.vars.stubbornness - 2
			}
		}
		if a.vars.morality > 50 { //want to implement effects of friendship
			changeInMorality(a, 5, 10, 1)
		} else {
			changeInMorality(a, 5, 10, -1)
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
	if read(a) {
		if a.HP() < a.knowledge.lastHP {
			if a.vars.stubbornness > 80 {
				a.vars.stubbornness = a.vars.stubbornness + 5
			} else {
				a.vars.stubbornness = a.vars.stubbornness - 2
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
		a.vars.stubbornness = a.vars.stubbornness - 5
		changeInMood(a, 5, 10, -1)
		changeInMorality(a, 5, 10, -1)
	} else {
		a.vars.stubbornness = a.vars.stubbornness + 5
		changeInMood(a, 5, 10, 1)
		changeInMorality(a, 5, 10, 1)
	}

	a.Log("I recieved a StateFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor(), "statement": statement})
}

func (a *CustomAgent3) HandleStateHP(msg messages.StateHPMessage) {
	statement := msg.Statement()
	if read(a) {
		if statement > a.HP() {
			a.vars.stubbornness = a.vars.stubbornness - 5
			changeInMorality(a, 5, 10, -1)
		} else {
			a.vars.stubbornness = a.vars.stubbornness - 10
			changeInMorality(a, 5, 10, 1)
			changeInMood(a, 5, 10, 1)
		}
	}

	a.Log("I recieved a StateHP message from ", infra.Fields{"floor": msg.SenderFloor(), "statement": statement})
}

func (a *CustomAgent3) HandleStateIntendedFoodTaken(msg messages.StateIntendedFoodIntakeMessage) {
	statement := msg.Statement()
	if read(a) {
		a.vars.stubbornness = a.vars.stubbornness - 5
		if statement > a.decisions.foodToEat {
			changeInMood(a, 5, 10, -1)
		} else {
			changeInMorality(a, 5, 10, 1)
		}

	}
	a.Log("I recieved a StateIntendedFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor(), "statement": statement})
}
