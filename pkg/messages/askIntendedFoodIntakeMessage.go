package messages

type AskIntendedFoodIntakeMessage struct {
	*BaseMessage
}

func NewAskIntendedFoodIntakeMessage(senderFloor int) *AskIntendedFoodIntakeMessage {
	msg := &AskIntendedFoodIntakeMessage{
		NewBaseMessage(senderFloor, AskIntendedFoodIntake),
	}
	return msg
}

func (msg *AskIntendedFoodIntakeMessage) Reply(senderFloor int, food int) StateMessage {
	reply := NewStateIntendedFoodIntakeMessage(senderFloor, food)
	return reply
}

func (msg *AskIntendedFoodIntakeMessage) Visit(a Agent) {
	a.HandleAskIntendedFoodTaken(*msg)
}
