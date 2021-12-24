package messages

type AskFoodTakenMessage struct {
	baseMessage *BaseMessage
}

func NewAskFoodTakenMessage(SenderFloor int) *AskFoodTakenMessage {
	msg := &AskFoodTakenMessage{
		baseMessage: NewBaseMessage(SenderFloor, AskFoodTaken),
	}
	return msg
}

func (msg *AskFoodTakenMessage) Reply(senderFloor int, food int) StateMessage {
	reply := NewStateFoodTakenMessage(senderFloor, food)
	return reply
}

func (msg *AskFoodTakenMessage) MessageType() MessageType {
	return msg.baseMessage.MessageType()
}

func (msg *AskFoodTakenMessage) SenderFloor() int {
	return msg.baseMessage.SenderFloor()
}

func (msg *AskFoodTakenMessage) Visit(a Agent) {
	a.HandleAskFoodTaken(*msg)
}
