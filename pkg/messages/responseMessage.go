package messages

type BoolResponseMessage struct {
	*BaseMessage
	response bool
}

func NewResponseMessage(senderFloor int, response bool) *BoolResponseMessage {
	msg := &BoolResponseMessage{
		NewBaseMessage(senderFloor, Response),
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
