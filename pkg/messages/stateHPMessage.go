package messages

import (
	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
)

// TODO (woonmoon): HP SHOULD BE AN INT!!!
// do this by supporting union return types on Statement() in
// AskMessage interface
type StateHPMessage struct {
	baseMessage *infra.BaseMessage
	hp          int
}

func NewStateHPMessage(senderFloor int, hp int) *StateHPMessage {
	msg := &StateHPMessage{
		baseMessage: infra.NewBaseMessage(senderFloor, infra.StateHP),
		hp:          hp,
	}
	return msg
}

func (msg *StateHPMessage) Statement() int {
	// TODO (woonmoon): GET RID OF THIS CONVERSION by supporting pseudo-union return types of Statement() function
	return msg.hp
}

func (msg *StateHPMessage) MessageType() infra.MessageType {
	return msg.baseMessage.MessageType()
}

func (msg *StateHPMessage) SenderFloor() int {
	return msg.baseMessage.SenderFloor()
}

func (msg *StateHPMessage) Visit(a infra.Agent) {
	a.HandleStateHP(msg)
}
