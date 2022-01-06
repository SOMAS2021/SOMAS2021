package messages

import (
	"strconv"

	"github.com/google/uuid"
)

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
	if msg.TargetFloor() != a.Floor() {
		a.HandlePropogate(msg)
	} else {
		a.HandleStateHP(*msg)
	}
}

func (msg *StateHPMessage) StoryLog() string {
	return strconv.Itoa(msg.hp)
}
