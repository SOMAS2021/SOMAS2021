package messages

import "github.com/google/uuid"

type TreatyResponseMessage struct {
	*BaseMessage
	response bool
	treatyID uuid.UUID
}

func NewTreatyResponseMessage(senderFloor int, response bool, treatyID uuid.UUID) *TreatyResponseMessage {
	msg := &TreatyResponseMessage{
		NewBaseMessage(senderFloor, Response),
		response,
		treatyID,
	}
	return msg
}

func (msg *TreatyResponseMessage) Response() bool {
	return msg.response
}

func (msg *TreatyResponseMessage) TreatyID() uuid.UUID {
	return msg.treatyID
}

func (msg *TreatyResponseMessage) Visit(a Agent) {
	a.HandleTreatyResponse(*msg)
}
