package messages

import (
	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
)

type StateFoodTakenMessage struct {
	baseMessage *infra.BaseMessage
	food        int // Planning to change this to int, see #21
}

func NewStateFoodTakenMessage(senderFloor int, foodTaken int) *StateFoodTakenMessage {
	msg := &StateFoodTakenMessage{
		baseMessage: infra.NewBaseMessage(senderFloor, infra.StateFoodTaken),
		food:        foodTaken,
	}
	return msg
}

func (msg *StateFoodTakenMessage) Statement() int {
	return msg.food
}

func (msg *StateFoodTakenMessage) MessageType() infra.MessageType {
	return msg.baseMessage.MessageType()
}

func (msg *StateFoodTakenMessage) SenderFloor() int {
	return msg.baseMessage.SenderFloor()
}

func (msg *StateFoodTakenMessage) Visit(a infra.Agent) {
	a.HandleStateFoodTaken(msg)
}
