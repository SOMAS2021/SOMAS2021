package messages

type takeFoodMessage struct{
	*baseMessage
	Food int
}

func NewtakeFoodMessage(SenderFloor int, food int) *takeFoodMessage {
	msg := &takeFoodMessage{
		baseMessage: NewBaseMessage(SenderFloor),
		Food: food ,
	}
	return msg
}

func (msg takeFoodMessage) MessageType() string {
	return "takeFoodMessage"
	
}