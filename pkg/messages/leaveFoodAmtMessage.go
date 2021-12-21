package messages

type leaveFoodAmtMessage struct{
	*baseMessage
	FoodAmt int
}

func NewleaveFoodMessage(SenderFloor int, foodAmt int) *leaveFoodAmtMessage {
	msg := &leaveFoodAmtMessage{
		baseMessage: NewBaseMessage(SenderFloor),
		FoodAmt: foodAmt ,
	}
	return msg
}

func (msg leaveFoodAmtMessage) MessageType() string {
	return "leaveFoodAmtMessage"
	
}
