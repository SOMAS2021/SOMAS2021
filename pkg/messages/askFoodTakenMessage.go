package messages

import (
	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
)

type AskFoodTakenMessage struct {
	baseMessage *infra.BaseMessage
}

func NewAskFoodTakenMessage(SenderFloor int) *AskFoodTakenMessage {
	msg := &AskFoodTakenMessage{
		baseMessage: infra.NewBaseMessage(SenderFloor, infra.AskFoodTaken),
	}
	return msg
}

func (msg *AskFoodTakenMessage) Reply(senderFloor int, food float64) infra.StateMessage {
	reply := NewStateFoodTakenMessage(senderFloor, food)
	return reply
}

func (msg *AskFoodTakenMessage) MessageType() infra.MessageType {
	return msg.baseMessage.MessageType()
}

func (msg *AskFoodTakenMessage) SenderFloor() int {
	return msg.baseMessage.SenderFloor()
}

func (msg *AskFoodTakenMessage) Visit(a infra.Agent) {
	a.HandleAskFoodTaken()
}
