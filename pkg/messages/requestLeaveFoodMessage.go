package messages

import (
	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
)

type RequestLeaveFoodMessage struct {
	baseMessage *infra.BaseMessage
	food        int
}

func NewRequestLeaveFoodMessage(SenderFloor int, food int) *RequestLeaveFoodMessage {
	msg := &RequestLeaveFoodMessage{
		baseMessage: infra.NewBaseMessage(SenderFloor, infra.RequestLeaveFood),
		food:        food,
	}
	return msg
}

func (msg *RequestLeaveFoodMessage) Request() int {
	return msg.food
}

func (msg *RequestLeaveFoodMessage) Reply(senderFloor int, response bool) infra.ResponseMessage {
	reply := NewResponseMessage(senderFloor, response)
	return reply
}

func (msg *RequestLeaveFoodMessage) MessageType() infra.MessageType {
	return msg.baseMessage.MessageType()
}

func (msg *RequestLeaveFoodMessage) SenderFloor() int {
	return msg.baseMessage.SenderFloor()
}

func (msg *RequestLeaveFoodMessage) Visit(a infra.Agent) {
	a.HandleRequestLeaveFood(msg)
}
