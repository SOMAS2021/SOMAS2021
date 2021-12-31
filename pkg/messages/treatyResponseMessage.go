package messages

import "github.com/google/uuid"

type TreatyResponseMessage struct {
	*BaseMessage
	response  bool
	treatyID  uuid.UUID
	requestID uuid.UUID
}

func NewTreatyResponseMessage(senderID uuid.UUID, senderFloor int, targetFloor int, response bool, treatyID uuid.UUID, requestID uuid.UUID) *TreatyResponseMessage {
	msg := &TreatyResponseMessage{
		NewBaseMessage(senderID, senderFloor, targetFloor, TreatyResponse),
		response,
		treatyID,
		requestID,
	}
	return msg
}

func (msg *TreatyResponseMessage) Response() bool {
	return msg.response
}

func (msg *TreatyResponseMessage) RequestID() uuid.UUID {
	return msg.requestID
}

func (msg *TreatyResponseMessage) TreatyID() uuid.UUID {
	return msg.treatyID
}

func (msg *TreatyResponseMessage) Visit(a Agent) {
	a.HandleTreatyResponse(*msg)
}
