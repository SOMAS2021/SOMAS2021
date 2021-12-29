package messages

import "github.com/google/uuid"

type RequestTakeFoodMessage struct {
	*BaseMessage
	food int
}

func NewRequestTakeFoodMessage(senderID uuid.UUID, senderFloor int, food int) *RequestTakeFoodMessage {
	msg := &RequestTakeFoodMessage{
		NewBaseMessage(senderID, senderFloor, RequestTakeFood),
		food,
	}
	return msg
}

func (msg *RequestTakeFoodMessage) Request() int {
	return msg.food
}

func (msg *RequestTakeFoodMessage) Reply(senderID uuid.UUID, senderFloor int, response bool) ResponseMessage {
	reply := NewResponseMessage(senderID, senderFloor, response)
	return reply
}

func (msg *RequestTakeFoodMessage) Visit(a Agent) {
	a.HandleRequestTakeFood(*msg)
}
