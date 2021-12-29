package messages

type RequestTakeFoodMessage struct {
	*BaseMessage
	food int
}

func NewRequestTakeFoodMessage(SenderFloor int, food int) *RequestTakeFoodMessage {
	msg := &RequestTakeFoodMessage{
		NewBaseMessage(SenderFloor, RequestTakeFood, ""),
		food,
	}
	return msg
}

func (msg *RequestTakeFoodMessage) Request() int {
	return msg.food
}

func (msg *RequestTakeFoodMessage) Reply(senderFloor int, response bool, id string) ResponseMessage {
	reply := NewResponseMessage(senderFloor, response, id)
	return reply
}

func (msg *RequestTakeFoodMessage) Visit(a Agent) {
	a.HandleRequestTakeFood(*msg)
}
