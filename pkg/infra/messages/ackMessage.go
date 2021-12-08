package messages

type AckMessage struct {
	*baseMessage
	Ack bool
}

func NewAckMessage(SenderFloor uint, ack bool) *AckMessage {
	msg := &AckMessage{
		baseMessage: NewBaseMessage(SenderFloor),
		Ack:         ack,
	}
	return msg
}

func (msg *AckMessage) MessageType() string {
	return "AckMessage"
}