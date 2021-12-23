package messages

import (
	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/agent"
)

type StateFoodTakenMessage struct {
	baseMessage *infra.BaseMessage
	food        float64 // Planning to change this to int, see #21
}

func NewStateFoodTakenMessage(senderFloor int, foodTaken float64) *StateFoodTakenMessage {
	msg := &StateFoodTakenMessage{
		baseMessage: infra.NewBaseMessage(senderFloor, infra.StateFoodTaken),
		food:        foodTaken,
	}
	return msg
}

func (msg *StateFoodTakenMessage) Statement() float64 {
	return msg.food
}

func (msg *StateFoodTakenMessage) MessageType() infra.MessageType {
	return msg.baseMessage.MessageType()
}

func (msg *StateFoodTakenMessage) SenderFloor() int {
	return msg.baseMessage.SenderFloor()
}

func (msg *StateFoodTakenMessage) Visit(a agent.Agent) {
	a.HandleStateFoodTaken(msg.food)
}
