package messages

import "github.com/google/uuid"

type StateHPMessage struct {
	*BaseMessage
	hp int
}

func NewStateHPMessage(senderID uuid.UUID, senderFloor int, targetFloor int, hp int) *StateHPMessage {
	msg := &StateHPMessage{
		NewBaseMessage(senderID, senderFloor, targetFloor, StateHP),
		hp,
	}
	return msg
}

func (msg *StateHPMessage) Statement() int {
	return msg.hp
}

func (msg *StateHPMessage) Visit(a Agent) {
	a.HandleStateHP(*msg)
}
