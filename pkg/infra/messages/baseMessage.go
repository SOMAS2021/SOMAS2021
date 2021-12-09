package messages

//Define message types to enable basic protocols, voting systems ...etc
type Message interface {
	MessageType() string
}

type baseMessage struct {
	SenderFloor  int
	responseType Message
}

func NewBaseMessage(SenderFloor int) *baseMessage {
	msg := &baseMessage{
		SenderFloor: SenderFloor,
	}
	return msg
}
func (msg baseMessage) MessageType() string {
	return "baseMessage"
}

// func (msg *baseMessage) SenderFloor() uint {
// 	return msg.senderFloor
// }

func (msg *baseMessage) ResponseType() string {
	return msg.responseType.MessageType()
}