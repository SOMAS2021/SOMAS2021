package messages

import (
	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/agent"
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

func (msg *AskIntendedFoodIntakeMessage) Reply(senderFloor int, food float64) infra.StateMessage {
	reply := NewStateIntendedFoodIntake(senderFloor, food)
	return reply
}

func (msg *AskIntendedFoodIntakeMessage) Visit(a agent.Agent) {
	a.HandleAskIntendedFoodTaken()
}
