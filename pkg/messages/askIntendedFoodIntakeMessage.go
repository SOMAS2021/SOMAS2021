package messages

type AskIntendedFoodIntakeMessage struct {
	*BaseMessage
}

func NewAskIntendedFoodIntakeMessage(senderFloor int) *AskIntendedFoodIntakeMessage {
	msg := &AskIntendedFoodIntakeMessage{
		NewBaseMessage(senderFloor, AskIntendedFoodIntake, ""),
	}
	return msg
}

func (msg *AskIntendedFoodIntakeMessage) Reply(senderFloor int, food int, id string) StateMessage {
	reply := NewStateIntendedFoodIntakeMessage(senderFloor, food, id)
	return reply
}

func (msg *AskIntendedFoodIntakeMessage) Visit(a Agent) {
	a.HandleAskIntendedFoodTaken(*msg)
}
