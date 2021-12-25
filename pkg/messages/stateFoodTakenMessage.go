package messages

type StateFoodTakenMessage struct {
	*BaseMessage
	food int // Planning to change this to int, see #21
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
