package messages

import "github.com/google/uuid"

type AskIntendedFoodIntakeMessage struct {
	*BaseMessage
}

func NewAskIntendedFoodIntakeMessage(senderID uuid.UUID, senderFloor int) *AskIntendedFoodIntakeMessage {
	msg := &AskIntendedFoodIntakeMessage{
		NewBaseMessage(senderID, senderFloor, AskIntendedFoodIntake),
	}
	return msg
}

func (msg *AskIntendedFoodIntakeMessage) Reply(senderID uuid.UUID, senderFloor int, food int) StateMessage {
	reply := NewStateIntendedFoodIntakeMessage(senderID, senderFloor, food)
	return reply
}

func (msg *AskIntendedFoodIntakeMessage) Visit(a Agent) {
	a.HandleAskIntendedFoodTaken(*msg)
}
