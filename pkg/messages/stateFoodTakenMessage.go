package messages

import "github.com/google/uuid"

type StateFoodTakenMessage struct {
	*BaseMessage
	food int
}

func NewStateFoodTakenMessage(senderID uuid.UUID, senderFloor int, targetFloor int, foodTaken int) *StateFoodTakenMessage {
	msg := &StateFoodTakenMessage{
		NewBaseMessage(senderID, senderFloor, targetFloor, StateFoodTaken),
		foodTaken,
	}
	return msg
}

func (msg *StateFoodTakenMessage) Statement() int {
	return msg.food
}

func (msg *StateFoodTakenMessage) Visit(a Agent) {
	if msg.TargetFloor() != a.Floor() {
		a.HandlePropogate(msg)
	} else {
		a.HandleStateFoodTaken(*msg)
	}
}
