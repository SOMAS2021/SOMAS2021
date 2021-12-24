package messages

type AskIntendedFoodIntakeMessage struct {
	baseMessage *BaseMessage
}

func NewAskIntendedFoodIntakeMessage(senderFloor int) *AskIntendedFoodIntakeMessage {
	msg := &AskIntendedFoodIntakeMessage{
		baseMessage: NewBaseMessage(senderFloor, AskIntendedFoodIntake),
	}
	return msg
}

func (msg *AskIntendedFoodIntakeMessage) MessageType() MessageType {
	return msg.baseMessage.MessageType()
}

func (msg *AskIntendedFoodIntakeMessage) SenderFloor() int {
	return msg.baseMessage.SenderFloor()
}

func (msg *AskIntendedFoodIntakeMessage) Reply(senderFloor int, food int) StateMessage {
	reply := NewStateIntendedFoodIntakeMessage(senderFloor, food)
	return reply
}

func (msg *AskIntendedFoodIntakeMessage) Visit(a Agent) {
	a.HandleAskIntendedFoodTaken(*msg)
}
