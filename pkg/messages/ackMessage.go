package messages

import "github.com/SOMAS2021/SOMAS2021/pkg/infra"

type AckMessage struct {
	baseMessage *infra.BaseMessage
	response    bool
}

func NewAckMessage(SenderFloor int, response bool) *AckMessage {
	msg := &AckMessage{
		baseMessage: infra.NewBaseMessage(SenderFloor, infra.StateResponse),
		response:    response,
	}
	return msg
}
