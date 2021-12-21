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