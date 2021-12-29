package messages

import (
	"github.com/google/uuid"
)

type BoolResponseMessage struct {
	*BaseMessage
	response bool
	returnId uuid.UUID
}

func NewResponseMessage(senderFloor int, response bool, returnId uuid.UUID) *BoolResponseMessage {
	msg := &BoolResponseMessage{
		NewBaseMessage(senderFloor, Response),
		response,
		returnId,
	}
	return msg
}

func (msg *BoolResponseMessage) Response() bool {
	return msg.response
}

func (msg *BoolResponseMessage) ReturnId() uuid.UUID {
	return msg.returnId
}

func (msg *BoolResponseMessage) Visit(a Agent) {
	a.HandleResponse(*msg)
}
