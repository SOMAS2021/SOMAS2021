package messages

import (
	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
)

type RequestTakeFoodMessage struct {
	baseMessage *infra.BaseMessage
	food        int
}

func NewRequestTakeFoodMessage(SenderFloor int, food int) *RequestTakeFoodMessage {
	msg := &RequestTakeFoodMessage{
		baseMessage: infra.NewBaseMessage(SenderFloor, infra.RequestTakeFood),
		food:        food,
	}
	return msg
}

func (msg *RequestTakeFoodMessage) Request() int {
	return msg.food
}

func (msg *RequestTakeFoodMessage) Reply(senderFloor int, response bool) infra.ResponseMessage {
	reply := NewResponseMessage(senderFloor, response)
	return reply
}

func (msg *RequestTakeFoodMessage) MessageType() infra.MessageType {
	return msg.baseMessage.MessageType()
}

func (msg *RequestTakeFoodMessage) SenderFloor() int {
	return msg.baseMessage.SenderFloor()
}

func (msg *RequestTakeFoodMessage) Visit(a infra.Agent) {
	a.HandleRequestTakeFood(msg)
}
