package messages

import (
	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
)

type AskHPMessage struct {
	baseMessage *infra.BaseMessage
}

func NewAskHPMessage(senderFloor int) *AskHPMessage {
	msg := &AskHPMessage{
		baseMessage: infra.NewBaseMessage(senderFloor, infra.AskHP),
	}
	return msg
}

func (msg *AskHPMessage) Reply(senderFloor int, hp float64) infra.StateMessage {
	reply := NewStateHPMessage(senderFloor, int(hp))
	return reply
}

func (msg *AskHPMessage) MessageType() infra.MessageType {
	return msg.baseMessage.MessageType()
}

func (msg *AskHPMessage) SenderFloor() int {
	return msg.baseMessage.SenderFloor()
}

func (msg *AskHPMessage) Visit(a infra.Agent) {
	a.HandleAskHP()
}
