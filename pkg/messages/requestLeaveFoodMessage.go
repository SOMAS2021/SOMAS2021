package messages

type RequestLeaveFoodMessage struct {
	*BaseMessage
	food int
}

func NewRequestLeaveFoodMessage(SenderFloor int, food int) *RequestLeaveFoodMessage {
	msg := &RequestLeaveFoodMessage{
		NewBaseMessage(SenderFloor, RequestLeaveFood),
		food,
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

func (msg *RequestLeaveFoodMessage) Visit(a Agent) {
	a.HandleRequestLeaveFood(*msg)
}
