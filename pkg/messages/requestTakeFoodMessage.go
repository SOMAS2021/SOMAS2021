package messages

type RequestTakeFoodMessage struct {
	baseMessage *BaseMessage
	food        int
}

func NewRequestTakeFoodMessage(SenderFloor int, food int) *RequestTakeFoodMessage {
	msg := &RequestTakeFoodMessage{
		baseMessage: NewBaseMessage(SenderFloor, RequestTakeFood),
		food:        food,
	}
	return msg
}

func (msg *RequestTakeFoodMessage) Request() int {
	return msg.food
}

func (msg *RequestTakeFoodMessage) Reply(senderFloor int, response bool) ResponseMessage {
	reply := NewResponseMessage(senderFloor, response)
	return reply
}

func (msg *RequestTakeFoodMessage) MessageType() MessageType {
	return msg.baseMessage.MessageType()
}

func (msg *RequestTakeFoodMessage) SenderFloor() int {
	return msg.baseMessage.SenderFloor()
}

func (msg *RequestTakeFoodMessage) Visit(a Agent) {
	a.HandleRequestTakeFood(*msg)
}
