package messages

import (
	"github.com/google/uuid"
)

//Define message types to enable basic protocols, voting systems ...etc
type Message interface {
	MessageType() string
}

type baseMessage struct {
	SenderFloor  int
	responseType Message
	messageID    string
}

func NewBaseMessage(SenderFloor int) *baseMessage {
	msg := &baseMessage{
		SenderFloor: SenderFloor,
		messageID:   uuid.New().String(),
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
func (msg *baseMessage) MessageID() string {
	return msg.messageID
}
