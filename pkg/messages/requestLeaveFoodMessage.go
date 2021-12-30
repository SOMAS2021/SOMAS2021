package messages

import "github.com/google/uuid"

type RequestLeaveFoodMessage struct {
	*BaseMessage
	food int
}

func NewRequestLeaveFoodMessage(senderID uuid.UUID, senderFloor int, food int) *RequestLeaveFoodMessage {
	msg := &RequestLeaveFoodMessage{
		NewBaseMessage(senderID, senderFloor, RequestLeaveFood),
		food,
	}
	return msg
}

func (msg *RequestLeaveFoodMessage) Request() int {
	return msg.food
}

func (msg *RequestLeaveFoodMessage) Reply(senderID uuid.UUID, senderFloor int, response bool) ResponseMessage {
	reply := NewResponseMessage(senderID, senderFloor, response, msg.ID())
	return reply
}

func (msg *RequestLeaveFoodMessage) Visit(a Agent) {
	a.HandleRequestLeaveFood(*msg)
}
