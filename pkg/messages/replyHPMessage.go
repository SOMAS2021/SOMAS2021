package messages

type replyHPMessage struct{
	*baseMessage
	HP int
}

func NewreplyHPMessage(SenderFloor int, hp int) *replyHPMessage {
	msg := &replyHPMessage{
		baseMessage: NewBaseMessage(SenderFloor),
		HP: hp,
	}
	return msg
}

func (msg replyHPMessage) MessageType() string {
	return "replyHPMessage"
	
}