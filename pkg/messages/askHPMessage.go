package messages

import "github.com/google/uuid"

type AskHPMessage struct {
	*BaseMessage
}

func NewAskHPMessage(senderID uuid.UUID, senderFloor int, targetFloor int) *AskHPMessage {
	msg := &AskHPMessage{
		NewBaseMessage(senderID, senderFloor, targetFloor, AskHP),
	}
	return msg
}

func (msg *AskHPMessage) Reply(senderID uuid.UUID, senderFloor int, targetFloor int, hp int) StateMessage {
	reply := NewStateHPMessage(senderID, senderFloor, targetFloor, hp)
	return reply
}

func (msg *AskHPMessage) Visit(a Agent) {
	a.HandleAskHP(*msg)
}
