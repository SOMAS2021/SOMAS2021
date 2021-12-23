package messages

import (
	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/agent"
)

type RequestLeaveFoodMessage struct {
	baseMessage *infra.BaseMessage
	food        float64
}

func NewRequestLeaveFoodMessage(SenderFloor int, food float64) *RequestLeaveFoodMessage {
	msg := &RequestLeaveFoodMessage{
		baseMessage: infra.NewBaseMessage(SenderFloor, infra.RequestLeaveFood),
		food:        food,
	}
	return msg
}

func (msg *RequestLeaveFoodMessage) Request() float64 {
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

func (msg *RequestLeaveFoodMessage) Visit(a agent.Agent) {
	a.HandleRequestLeaveFood(msg.food)
}
