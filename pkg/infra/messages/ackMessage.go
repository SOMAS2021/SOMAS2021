package messages

type AckMessage struct {
	*baseMessage
	Ack bool
}

func NewAckMessage(SenderFloor int, ack bool) *AckMessage {
	msg := &AckMessage{
		baseMessage: NewBaseMessage(SenderFloor),
		Ack:         ack,
	}
	return msg
}

func (msg AckMessage) MessageType() string {
	return "AckMessage"
}