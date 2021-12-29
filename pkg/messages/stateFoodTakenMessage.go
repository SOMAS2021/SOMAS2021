package messages

type StateFoodTakenMessage struct {
	*BaseMessage
	food int
}

func NewStateFoodTakenMessage(senderFloor int, foodTaken int, id string) *StateFoodTakenMessage {
	msg := &StateFoodTakenMessage{
		NewBaseMessage(senderFloor, StateFoodTaken, id),
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
