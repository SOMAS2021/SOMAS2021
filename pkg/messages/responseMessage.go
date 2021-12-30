package messages

import "github.com/google/uuid"

type BoolResponseMessage struct {
	*BaseMessage
	response bool
}

func NewResponseMessage(senderID uuid.UUID, senderFloor int, response bool) *BoolResponseMessage {
	msg := &BoolResponseMessage{
		NewBaseMessage(senderID, senderFloor, Response),
		response,
	}
	return msg
}

func (msg *BoolResponseMessage) Response() bool {
	return msg.response
}

func (msg *BoolResponseMessage) Visit(a Agent) {
	a.HandleResponse(*msg)
}
