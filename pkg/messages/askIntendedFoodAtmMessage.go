package messages

type askIntendedFoodAtmMessage struct{
	*baseMessage
	//foodToTake float64 // Planning to change this to int, see #21 
}

func NewaskIntendedFoodAtmMessage(SenderFloor int) *askIntendedFoodAtmMessage {
	msg := &askIntendedFoodAtmMessage{
		baseMessage: NewBaseMessage(SenderFloor),
	}
	return msg
}

func (msg askIntendedFoodAtmMessage) MessageType() string {
	return "askIntendedFoodAtmMessage"
	
}