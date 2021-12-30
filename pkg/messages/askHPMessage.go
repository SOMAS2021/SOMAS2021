package messages

import "github.com/google/uuid"

type AskHPMessage struct {
	*BaseMessage
}

func NewAskHPMessage(senderID uuid.UUID, senderFloor int) *AskHPMessage {
	msg := &AskHPMessage{
		NewBaseMessage(senderID, senderFloor, AskHP),
	}
	return msg
}

func (msg *AskHPMessage) Reply(senderID uuid.UUID, senderFloor int, hp int) StateMessage {
	reply := NewStateHPMessage(senderID, senderFloor, hp)
	return reply
}

func (msg *AskHPMessage) Visit(a Agent) {
	a.HandleAskHP(*msg)
}
