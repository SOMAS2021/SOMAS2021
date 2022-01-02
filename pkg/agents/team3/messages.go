package team3

import (
	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
	"math/rand"
)

//Upon receipt of message define affected emotions
// ACK MESSAGE
//If x time passed no message received/acked morale decrease
//Include if ack message same user ID occurs x+1 times, morale increase
//If stubborness = y+1, discard, a.k.a. leave unread


func message(a *CustomAgent3) {
	receivedMsg := a.ReceiveMessage()
	if receivedMsg != nil {
		a.Log("Custom agent 3 each run:", infra.Fields{"floor": a.Floor(), "hp": a.HP(), "Mood": a.vars.mood, "Morality": a.vars.morality})
		receivedMsg.Visit(a)
		a.Log("Custom agent 3 each run:", infra.Fields{"floor": a.Floor(), "hp": a.HP(), "Mood": a.vars.mood, "Morality": a.vars.morality})
	} else {
		a.Log("I got no thing")
	}

	r := rand.Intn(9)
	switch r {
	case 0:
		msg := messages.NewAskFoodTakenMessage(a.Floor())
		a.SendMessage(1, msg)
		a.Log("I sent a message", infra.Fields{"message": "AskFoodTaken"})
	case 1:
		msg := messages.NewAskHPMessage(a.Floor())
		a.SendMessage(1, msg)
		a.Log("I sent a message", infra.Fields{"message": "AskHP"})
	case 2:
		msg := messages.NewAskIntendedFoodIntakeMessage(a.Floor())
		a.SendMessage(1, msg)
		a.Log("I sent a message", infra.Fields{"message": "AskIntendedFoodIntake"})
	case 3:
		msg := messages.NewRequestLeaveFoodMessage(a.Floor(), 10)
		a.SendMessage(1, msg)
		a.Log("I sent a message", infra.Fields{"message": "RequestLeaveFood"})
	case 4:
		msg := messages.NewRequestTakeFoodMessage(a.Floor(), 20)
		a.SendMessage(1, msg)
		a.Log("I sent a message", infra.Fields{"message": "RequestTakeFood"})
	case 5:
		msg := messages.NewResponseMessage(a.Floor(), true)
		a.SendMessage(1, msg)
		a.Log("I sent a message", infra.Fields{"message": "Response"})
	case 6:
		msg := messages.NewStateFoodTakenMessage(a.Floor(), 30)
		a.SendMessage(1, msg)
		a.Log("I sent a message", infra.Fields{"message": "StateFoodTaken"})
	case 7:
		msg := messages.NewStateHPMessage(a.Floor(), 40)
		a.SendMessage(1, msg)
		a.Log("I sent a message", infra.Fields{"message": "StateHP"})
	case 8:
		msg := messages.NewStateIntendedFoodIntakeMessage(a.Floor(), 50)
		a.SendMessage(1, msg)
		a.Log("I sent a message", infra.Fields{"message": "StateIntendedFoodIntake"})
	}
	a.Log("My agent is doing something", infra.Fields{"thing": "potatoe", "another_thing": "another potatoe"})
	a.TakeFood(16)
}

func (a *CustomAgent3) HandleAskHP(msg messages.AskHPMessage) {
	a.Log("I recieved an askHP message from ", infra.Fields{"floor": msg.SenderFloor()})
	if read(a){
		a.vars.stubbornness = a.vars.stubbornness - 5 //value could be different
		reply := msg.Reply(a.Floor(), a.HP())
		a.SendMessage(msg.SenderFloor()-a.Floor(), reply)
		a.Log("I recieved an askHP message from ", infra.Fields{"floor": msg.SenderFloor()})
	}	
}

func (a *CustomAgent3) HandleAskFoodTaken(msg messages.AskFoodTakenMessage) {
	a.Log("I recieved an askFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor()})
	if read(a){
		if a.HP() < a.knowledge.lastHP{
			a.vars.stubbornness = a.vars.stubbornness + 5
			//addfriend(a, ) need id
			if a.vars.morality < 30{
				initial_resp := messages.NewAskFoodTakenMessage(a.Floor())
				a.SendMessage(msg.SenderFloor()-a.Floor(), initial_resp)
				a.Log("I sent an initial_resp message", infra.Fields{"message": "AskFoodTaken"})
			}
			changeInMood(a, 5, 10, 1)
		}	
		reply := msg.Reply(a.Floor(), 10)
		a.SendMessage(msg.SenderFloor()-a.Floor(), reply)
		a.Log("I sent a replyFoodTaken message to ", infra.Fields{"floor": msg.SenderFloor()})
	}
}

func (a *CustomAgent3) HandleAskIntendedFoodTaken(msg messages.AskIntendedFoodIntakeMessage) {
	if read(a){
		a.vars.stubbornness = a.vars.stubbornness + 2
		if a.vars.morality < 30{
			initial_resp := messages.NewAskIntendedFoodIntakeMessage(a.Floor())
			a.SendMessage(msg.SenderFloor()-a.Floor(), initial_resp)
			a.Log("I sent an initial_resp message", infra.Fields{"message": "AskIntendedFoodIntake"})
		}
		reply := msg.Reply(a.Floor(), 11)
		a.SendMessage(msg.SenderFloor()-a.Floor(), reply)
		a.Log("I recieved an askIntendedFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor()})
	}
}	

func (a *CustomAgent3) HandleRequestLeaveFood(msg messages.RequestLeaveFoodMessage) {
	if read(a){
		if a.HP() < a.knowledge.lastHP{
			if a.vars.stubbornness > 80{
				a.vars.stubbornness = a.vars.stubbornness + 5
			} else{
				a.vars.stubbornness = a.vars.stubbornness - 2
			}
		}
		if a.vars.morality > 50{ //want to implement effects of friendship
			changeInMorality(a, 5, 10, 1)
		} else{
			changeInMorality(a, 5, 10, -1)
		}
		if a.vars.mood < 30{ //can we see when we are in critical state
			changeInMood(a, 5, 10, -1)
		}
		reply := msg.Reply(a.Floor(), true)
		a.SendMessage(msg.SenderFloor()-a.Floor(), reply)
		a.Log("I recieved a requestLeaveFood message from ", infra.Fields{"floor": msg.SenderFloor()})
	}	
}

func (a *CustomAgent3) HandleRequestTakeFood(msg messages.RequestTakeFoodMessage) {
	if read(a){
		if a.HP() < a.knowledge.lastHP{
			if a.vars.stubbornness > 80{
				a.vars.stubbornness = a.vars.stubbornness + 5
			} else {
				a.vars.stubbornness = a.vars.stubbornness - 2
			}
		}
		if a.vars.morality > 50{ //want to implement effects of friendship
			changeInMorality(a, 5, 10, 1)
		} else{
			changeInMorality(a, 5, 10, -1)
		}
		if a.vars.mood < 30{
			changeInMood(a, 5, 10, -1)
		}
		reply := msg.Reply(a.Floor(), true)
		a.SendMessage(msg.SenderFloor()-a.Floor(), reply)
		a.Log("I recieved a requestTakeFood message from ", infra.Fields{"floor": msg.SenderFloor()})
	}
}

func (a *CustomAgent3) HandleResponse(msg messages.BoolResponseMessage) {
	response := msg.Response()
	a.Log("I recieved a Response message from ", infra.Fields{"floor": msg.SenderFloor(), "response": response})
}

func (a *CustomAgent3) HandleStateFoodTaken(msg messages.StateFoodTakenMessage) {
	statement := msg.Statement()
	if statement > a.decisions.foodToEat{
		a.vars.stubbornness = a.vars.stubbornness - 5
		changeInMood(a, 5, 10, -1)
		changeInMorality(a, 5, 10, -1)
	} else{
		a.vars.stubbornness = a.vars.stubbornness + 5
		changeInMood(a, 5, 10, 1)
		changeInMorality(a, 5, 10, 1)
	}
	
	a.Log("I recieved a StateFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor(), "statement": statement})
}

func (a *CustomAgent3) HandleStateHP(msg messages.StateHPMessage) {
	statement := msg.Statement()
	if read(a){
		if statement > a.HP(){
			a.vars.stubbornness = a.vars.stubbornness - 5
			changeInMorality(a, 5, 10, -1)
		} else{
			a.vars.stubbornness = a.vars.stubbornness - 10
			changeInMorality(a, 5, 10, 1)
			changeInMood(a, 5, 10, 1)	
		}
	}
	
	a.Log("I recieved a StateHP message from ", infra.Fields{"floor": msg.SenderFloor(), "statement": statement})
}

func (a *CustomAgent3) HandleStateIntendedFoodTaken(msg messages.StateIntendedFoodIntakeMessage) {
	statement := msg.Statement()
	if read(a){
		a.vars.stubbornness = a.vars.stubbornness - 5
		if statement > a.decisions.foodToEat{
			changeInMood(a, 5, 10, -1)
		} else{
			changeInMorality(a, 5, 10, 1)
		}

	}
	a.Log("I recieved a StateIntendedFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor(), "statement": statement})
}

//	receivedMsg := a.Base.ReceiveMessage()
//	if receivedMsg != nil {
//		a.Log("I got sent a message", infra.Fields{"messageType": receivedMsg.MessageType()})
//		if receivedMsg.MessageType() == "askFoodTakenMessage" { //agents response to askFoodTakenMessage
//			a.Log("I've eaten", infra.Fields{"food": takeFoodCalculation(a)})
//			msg := *messages.NewfoodTakenMessage(int(a.Floor()), float64(takeFoodCalculation(a)))
//			//need to find out how to find where message came from
//			//if receivedMsg.SenderFloor  == a.knowledge.floorBelow{
//			a.SendMessage(-1, msg)
//			a.Log("I sent a msg", infra.Fields{"floor": a.Floor(), "message": msg.MessageType()})
//			//}
//		}
//		if receivedMsg.MessageType() == "askHPMessage" { //agents response to askHPMessage
//			a.Log("My HP is", infra.Fields{"HP": a.HP()})
//			msg := *messages.NewreplyHPMessage(int(a.Floor()), a.HP())
//			a.SendMessage(-1, msg)
//			a.Log("I sent a msg", infra.Fields{"floor": a.Floor(), "message": msg.MessageType()})
//		}
//		if receivedMsg.MessageType() == "askIntendedFoodAtmMessage" { //agents response to askIntendedFoodAtmMessage
//			a.Log("My intended food is", infra.Fields{"food": takeFoodCalculation(a)})
//			msg := *messages.NewintendedFoodAtmMessage(int(a.Floor()), float64(takeFoodCalculation(a)))
//			a.SendMessage(-1, msg)
//			a.Log("I sent a msg", infra.Fields{"floor": a.Floor(), "message": msg.MessageType()})
//		}
//		if receivedMsg.MessageType() == "leaveFoodMessage" { //agents response to leaveFoodMessage
//			a.Log("I've been asked to leave this much food", infra.Fields{"message": receivedMsg, "Food": receivedMsg}) //, infra.Fields{"foodAmt": receivedMsg.intFood})
//			//msg := *messages.NewackMessage(int(a.Floor()), true)?
//			//a.SendMessage(-1, msg)
//			a.Log("I didn't send a msg")
//		}
//		if receivedMsg.MessageType() == "takeFoodMessage" { //agents response to takeFoodMessage
//			a.Log("I've been asked to take this much food") //, infra.Fields{"foodAmt": receivedMsg.intFood})
//			//msg := *messages.NewackMessage(int(a.Floor()), true)?
//			//a.SendMessage(-1, msg)
//			a.Log("I didn't send a msg")
//		}
//		if receivedMsg.MessageType() == "whoAreYouMessage" { //agents response to whoAreYouMessage
//			a.Log("I am", infra.Fields{"ID": a.ID()})
//			msg := *messages.NewwhoIAmMessage(int(a.Floor()), a.ID())
//			a.SendMessage(-1, msg)
//			a.Log("I sent a msg", infra.Fields{"floor": a.Floor(), "message": msg.MessageType()})
//		}
//
//	} else {
//		a.Log("I got nothing", infra.Fields{"floor": a.Floor()})
//	}
//}

//if msg != *messages.NewBaseMessage(int(a.Floor())) {
//	a.SendMessage(-1, msg)
//	a.Log("I sent a msg", infra.Fields{"floor": a.Floor(), "message": msg.MessageType()})
//}
