package messages

type AskFoodTakenMessage struct {
	*BaseMessage
}

func NewAskFoodTakenMessage(SenderFloor int) *AskFoodTakenMessage {
	msg := &AskFoodTakenMessage{
		NewBaseMessage(SenderFloor, AskFoodTaken),
	}
	return msg
}

func (msg *AskFoodTakenMessage) Reply(senderFloor int, food int) StateMessage {
	reply := NewStateFoodTakenMessage(senderFloor, food)
	return reply
}

func (msg *AskFoodTakenMessage) Visit(a Agent) {
	a.HandleAskFoodTaken(*msg)
}
