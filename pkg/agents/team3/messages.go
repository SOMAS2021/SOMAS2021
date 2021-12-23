package team3

import (
	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
)
//Upon receipt of message define affected emotions
// ACK MESSAGE
//If x time passed no message received/acked morale decrease
//Include if ack message same user ID occurs x+1 times, morale increase
//If stubborness = y+1, discard, a.k.a. leave unread

//
func message(a *CustomAgent3) {
	receivedMsg := a.Base.ReceiveMessage()
	if receivedMsg != nil {
		a.Log("I got sent a message", infra.Fields{"messageType": receivedMsg.MessageType(), "Food": receivedMsg})
		if receivedMsg.MessageType() == "askFoodTakenMessage" { //agents response to askFoodTakenMessage
			a.Log("I've eaten", infra.Fields{"food": takeFoodCalculation(a)})
			msg := *messages.NewfoodTakenMessage(int(a.Floor()), takeFoodCalculation(a))
			//need to find out how to find where message came from
			//if receivedMsg.SenderFloor  == a.knowledge.floorBelow{
			a.SendMessage(-1, msg)
			a.Log("I sent a msg", infra.Fields{"floor": a.Floor(), "message": msg.MessageType()})
			//}
		}
		if receivedMsg.MessageType() == "askHPMessage" { //agents response to askHPMessage
			a.Log("My HP is", infra.Fields{"HP": a.HP()})
			msg := *messages.NewreplyHPMessage(int(a.Floor()), a.HP())
			a.SendMessage(-1, msg)
			a.Log("I sent a msg", infra.Fields{"floor": a.Floor(), "message": msg.MessageType()})
		}
		if receivedMsg.MessageType() == "askIntendedFoodAtmMessage" { //agents response to askIntendedFoodAtmMessage
			a.Log("My intended food is", infra.Fields{"food": takeFoodCalculation(a)})
			msg := *messages.NewintendedFoodAtmMessage(int(a.Floor()), takeFoodCalculation(a))
			a.SendMessage(-1, msg)
			a.Log("I sent a msg", infra.Fields{"floor": a.Floor(), "message": msg.MessageType()})
		}
		if receivedMsg.MessageType() == "leaveFoodMessage" { //agents response to leaveFoodMessage
			a.Log("I've been asked to leave this much food", infra.Fields{"message": receivedMsg, "Food": receivedMsg})//, infra.Fields{"foodAmt": receivedMsg.intFood})
			//msg := *messages.NewackMessage(int(a.Floor()), true)?
			//a.SendMessage(-1, msg)
			a.Log("I didn't send a msg")
		}
		if receivedMsg.MessageType() == "takeFoodMessage" { //agents response to takeFoodMessage
			a.Log("I've been asked to take this much food")//, infra.Fields{"foodAmt": receivedMsg.intFood})
			//msg := *messages.NewackMessage(int(a.Floor()), true)?
			//a.SendMessage(-1, msg)
			a.Log("I didn't send a msg")
		}
		if receivedMsg.MessageType() == "whoAreYouMessage" { //agents response to whoAreYouMessage
			a.Log("I am", infra.Fields{"ID": a.ID()})
			msg := *messages.NewwhoIAmMessage(int(a.Floor()), a.ID())
			a.SendMessage(-1, msg)
			a.Log("I sent a msg", infra.Fields{"floor": a.Floor(), "message": msg.MessageType()})
		}

	} else {
		a.Log("I got nothing", infra.Fields{"floor": a.Floor()})
	}
}


//if msg != *messages.NewBaseMessage(int(a.Floor())) {
	//	a.SendMessage(-1, msg)
	//	a.Log("I sent a msg", infra.Fields{"floor": a.Floor(), "message": msg.MessageType()})
	//}	
