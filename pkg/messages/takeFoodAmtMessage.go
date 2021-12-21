package messages

type takeFoodAmtMessage struct{
	*baseMessage
	Food int
}

func NewtakeFoodAmtMessage(SenderFloor int, food int) *takeFoodAmtMessage {
	msg := &takeFoodAmtMessage{
		baseMessage: NewBaseMessage(SenderFloor),
		Food: food ,
	}
	return msg
}

func (msg takeFoodAmtMessage) MessageType() string {
	return "takeFoodAmtMessage"
	
}
