package messages

type StateFoodTakenMessage struct {
	*BaseMessage
	food int
}

func NewStateFoodTakenMessage(senderFloor int, foodTaken int) *StateFoodTakenMessage {
	msg := &StateFoodTakenMessage{
		NewBaseMessage(senderFloor, StateFoodTaken),
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
