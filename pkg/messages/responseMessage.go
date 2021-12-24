package messages

type BoolResponseMessage struct {
	baseMessage BaseMessage
	response    bool
}

func NewResponseMessage(senderFloor int, response bool) *BoolResponseMessage {
	msg := &BoolResponseMessage{
		baseMessage: *NewBaseMessage(senderFloor, Response),
		response:    response,
	}
	return msg
}

func (msg *BoolResponseMessage) MessageType() MessageType {
	return msg.baseMessage.MessageType()
}

func (msg *BoolResponseMessage) SenderFloor() int {
	return msg.baseMessage.SenderFloor()
}

func (msg *BoolResponseMessage) Response() bool {
	return msg.response
}

func (msg *BoolResponseMessage) Visit(a Agent) {
	a.HandleResponse(*msg)
}
