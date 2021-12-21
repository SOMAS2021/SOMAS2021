package messages

type intendedFoodAtmMessage struct{
	*baseMessage
	IntendedFood float64 // Planning to change this to int, see #21 
}

func NewintendedFoodAtmMessage(SenderFloor int, intendedFood float64) *intendedFoodAtmMessage {
	msg := &intendedFoodAtmMessage{
		baseMessage: NewBaseMessage(SenderFloor),
		IntendedFood: intendedFood,
	}
	return msg
}

func (msg intendedFoodAtmMessage) MessageType() string {
	return "intendedFoodAtmMessage"
	
}