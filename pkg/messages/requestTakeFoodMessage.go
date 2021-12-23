package messages

import (
	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
)

type RequestTakeFoodMessage struct {
	baseMessage *infra.BaseMessage
	food        float64
}

func NewtakeFoodMessage(SenderFloor int, food float64) *RequestTakeFoodMessage {
	msg := &RequestTakeFoodMessage{
		baseMessage: infra.NewBaseMessage(SenderFloor, infra.RequestTakeFood),
		food:        food,
	}
	return msg
}

func (msg *RequestTakeFoodMessage) Request() float64 {
	return msg.food
}

func (msg *RequestTakeFoodMessage) MessageType() infra.MessageType {
	return msg.baseMessage.MessageType()
}

func (msg *RequestTakeFoodMessage) SenderFloor() int {
	return msg.baseMessage.SenderFloor()
}

func (msg *RequestTakeFoodMessage) Visit(a infra.Agent) {
	a.HandleRequestTakeFood(msg.food)
}
