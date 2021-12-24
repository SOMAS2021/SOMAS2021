package messages

type StateFoodTakenMessage struct {
	baseMessage *BaseMessage
	food        int // Planning to change this to int, see #21
}

func NewStateFoodTakenMessage(senderFloor int, foodTaken int) *StateFoodTakenMessage {
	msg := &StateFoodTakenMessage{
		baseMessage: NewBaseMessage(senderFloor, StateFoodTaken),
		food:        foodTaken,
	}
	return msg
}

func (msg *StateFoodTakenMessage) Statement() int {
	return msg.food
}

func (msg *StateFoodTakenMessage) MessageType() MessageType {
	return msg.baseMessage.MessageType()
}

func (msg *StateFoodTakenMessage) SenderFloor() int {
	return msg.baseMessage.SenderFloor()
}

func (msg *StateFoodTakenMessage) Visit(a Agent) {
	a.HandleStateFoodTaken(*msg)
}
