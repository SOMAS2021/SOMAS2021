package messages

type askIntendedFoodAmtMessage struct{
	*baseMessage
}

func NewaskIntendedFoodAmtMessage(SenderFloor int) *askIntendedFoodAmtMessage {
	msg := &askIntendedFoodAmtMessage{
		baseMessage: NewBaseMessage(SenderFloor),
	}
	return msg
}

func (msg askIntendedFoodAmtMessage) MessageType() string {
	return "askIntendedFoodAmtMessage"
	
}