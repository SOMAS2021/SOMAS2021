package messages

type Team6Message struct {
	*baseMessage
	Ack bool
}

func NewTeam6Message(SenderFloor int, ack bool) *Team6Message {
	msg := &Team6Message{
		baseMessage: NewBaseMessage(SenderFloor),
		Ack:         ack,
	}
	return msg
}

func (msg Team6Message) MessageType() string {
	return "Wagwan babes"
}
