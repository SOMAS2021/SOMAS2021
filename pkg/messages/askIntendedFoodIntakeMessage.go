package messages

import "github.com/google/uuid"

type AskIntendedFoodIntakeMessage struct {
	*BaseMessage
}

func NewAskIntendedFoodIntakeMessage(senderID uuid.UUID, senderFloor int, targetFloor int) *AskIntendedFoodIntakeMessage {
	msg := &AskIntendedFoodIntakeMessage{
		NewBaseMessage(senderID, senderFloor, targetFloor, AskIntendedFoodIntake),
	}
	return msg
}

func (msg *AskIntendedFoodIntakeMessage) Reply(senderID uuid.UUID, senderFloor int, targetFloor int, food int) StateMessage {
	reply := NewStateIntendedFoodIntakeMessage(senderID, senderFloor, targetFloor, food)
	return reply
}

func (msg *AskIntendedFoodIntakeMessage) Visit(a Agent) {
	if msg.TargetFloor() != a.Floor() {
		a.HandlePropogate(msg)
	} else {
		a.HandleAskIntendedFoodTaken(*msg)
	}
}
