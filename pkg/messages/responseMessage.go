package messages

import (
	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
)

type ResponseMessage struct {
	baseMessage *infra.BaseMessage
	response    bool
}

func NewResponseMessage(senderFloor int, response bool) *ResponseMessage {
	msg := &ResponseMessage{
		baseMessage: infra.NewBaseMessage(senderFloor, infra.Response),
		response:    response,
	}
	return msg
}

func (msg *ResponseMessage) MessageType() infra.MessageType {
	return msg.baseMessage.MessageType()
}

func (msg *ResponseMessage) SenderFloor() int {
	return msg.baseMessage.SenderFloor()
}

func (msg *ResponseMessage) Response() bool {
	return msg.response
}

func (msg *ResponseMessage) Visit(a infra.Agent) {
	a.HandleResponse(msg.response)
}
