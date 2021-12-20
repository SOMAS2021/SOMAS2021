package messages

type foodTakenMessage struct{
	*baseMessage
	FoodTaken float64 // Planning to change this to int, see #21
}

func NewfoodTakenMessage(SenderFloor int, foodTaken float64) *foodTakenMessage {
	msg := &foodTakenMessage{
		baseMessage: NewBaseMessage(SenderFloor),
		FoodTaken: foodTaken,
	}
	return msg
}

func (msg foodTakenMessage) MessageType() string {
	return "foodTakenMessage"

}


//return foodeaten.vari., ackmessage.
//upon receipt of foodtaken, if foodtaken.int>your.foodtaken.int, morale decrease, stubborness increase.
//if agent1 sends askFoodTakenMessage, agent2 sends foodTakenMessage.
//then agent1 has morale inc/dec, mood inc/dec, stubborness inc/dec, if agent1 foodeaten<foodtaken agent2
//agent2  has morale inc/dec, if agent1[userID] un/known. "fuck this guy".
//
