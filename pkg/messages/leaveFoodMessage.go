package messages

type leaveFoodMessage struct{
	*baseMessage
	Food int
}

func NewleaveFoodMessage(SenderFloor int, food int) *leaveFoodMessage {
	msg := &leaveFoodMessage{
		baseMessage: NewBaseMessage(SenderFloor),
		Food: food ,
	}
	return msg
}

func (msg leaveFoodMessage) MessageType() string {
	return "leaveFoodMessage"
	
}