package messages

import "github.com/google/uuid"

type BoolResponseMessage struct {
	*BaseMessage
	response  bool
	requestId uuid.UUID
}

func NewResponseMessage(senderID uuid.UUID, senderFloor int, response bool, requestID uuid.UUID) *BoolResponseMessage {
	msg := &BoolResponseMessage{
		NewBaseMessage(senderID, senderFloor, Response),
		response,
		requestID,
	}
	return msg
}

func (msg *BoolResponseMessage) Response() bool {
	return msg.response
}

func (msg *BoolResponseMessage) RequestID() uuid.UUID {
	return msg.requestId
}

func (msg *BoolResponseMessage) Visit(a Agent) {
	a.HandleResponse(*msg)
}
