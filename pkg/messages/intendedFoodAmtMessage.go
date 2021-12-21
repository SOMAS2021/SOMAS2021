package messages

type intendedFoodAmtMessage struct{
	*baseMessage
	IntendedFoodAmt float64 // Planning to change this to int, see #21 
}

func NewintendedFoodAmtMessage(SenderFloor int, intendedFoodAmt float64) *intendedFoodAmtMessage {
	msg := &intendedFoodAmtMessage{
		baseMessage: NewBaseMessage(SenderFloor),
		IntendedFoodAmt: intendedFoodAmt,
	}
	return msg
}

func (msg intendedFoodAmtMessage) MessageType() string {
	return "intendedFoodAmtMessage"
	
}
