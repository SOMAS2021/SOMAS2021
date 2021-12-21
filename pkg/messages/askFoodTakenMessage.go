package messages

type askFoodTakenMessage struct{
	*baseMessage
}

func NewaskFoodTakenMessage(SenderFloor int) *askFoodTakenMessage {
	msg := &askFoodTakenMessage{
		baseMessage: NewBaseMessage(SenderFloor),
	}
	return msg
}

func (msg askFoodTakenMessage) MessageType() string {
	return "askFoodTakenMessage"
	
}
