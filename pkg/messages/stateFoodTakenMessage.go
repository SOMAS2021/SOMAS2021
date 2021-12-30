package messages

import "github.com/google/uuid"

type StateFoodTakenMessage struct {
	*BaseMessage
	food int
}

func NewStateFoodTakenMessage(senderID uuid.UUID, senderFloor int, foodTaken int) *StateFoodTakenMessage {
	msg := &StateFoodTakenMessage{
		NewBaseMessage(senderID, senderFloor, StateFoodTaken),
		foodTaken,
	}
	return msg
}

func (msg *StateFoodTakenMessage) Statement() int {
	return msg.food
}

func (msg *StateFoodTakenMessage) Visit(a Agent) {
	a.HandleStateFoodTaken(*msg)
}
