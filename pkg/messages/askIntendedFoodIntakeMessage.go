package messages

import (
	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
)

type AskIntendedFoodIntakeMessage struct {
	baseMessage *infra.BaseMessage
}

func NewAskIntendedFoodIntakeMessage(SenderFloor int) *AskIntendedFoodIntakeMessage {
	msg := &AskIntendedFoodIntakeMessage{
		baseMessage: infra.NewBaseMessage(SenderFloor, infra.AskIntendedFoodIntake),
	}
	return msg
}

func (msg *AskIntendedFoodIntakeMessage) MessageType() infra.MessageType {
	return msg.baseMessage.MessageType()
}

func (msg *AskIntendedFoodIntakeMessage) SenderFloor() int {
	return msg.baseMessage.SenderFloor()
}

func (msg *AskIntendedFoodIntakeMessage) Reply(senderFloor int, food int) infra.StateMessage {
	reply := NewStateIntendedFoodIntakeMessage(senderFloor, food)
	return reply
}

func (msg *AskIntendedFoodIntakeMessage) Visit(a infra.Agent) {
	a.HandleAskIntendedFoodTaken(msg)
}
