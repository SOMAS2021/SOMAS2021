package messages

type RequestLeaveFoodMessage struct {
	baseMessage *BaseMessage
	food        int
}

func NewRequestLeaveFoodMessage(SenderFloor int, food int) *RequestLeaveFoodMessage {
	msg := &RequestLeaveFoodMessage{
		baseMessage: NewBaseMessage(SenderFloor, RequestLeaveFood),
		food:        food,
	}
	return msg
}

func (msg *RequestLeaveFoodMessage) Request() int {
	return msg.food
}

func (msg *RequestLeaveFoodMessage) Reply(senderFloor int, response bool) ResponseMessage {
	reply := NewResponseMessage(senderFloor, response)
	return reply
}

func (msg *RequestLeaveFoodMessage) MessageType() MessageType {
	return msg.baseMessage.MessageType()
}

func (msg *RequestLeaveFoodMessage) SenderFloor() int {
	return msg.baseMessage.SenderFloor()
}

func (msg *RequestLeaveFoodMessage) Visit(a Agent) {
	a.HandleRequestLeaveFood(*msg)
}
