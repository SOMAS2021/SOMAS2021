package messages

import "github.com/google/uuid"

type AskFoodTakenMessage struct {
	*BaseMessage
}

func NewAskFoodTakenMessage(senderID uuid.UUID, senderFloor int, targetFloor int) *AskFoodTakenMessage {
	msg := &AskFoodTakenMessage{
		NewBaseMessage(senderID, senderFloor, targetFloor, AskFoodTaken),
	}
	return msg
}

func (msg *AskFoodTakenMessage) Reply(senderID uuid.UUID, senderFloor int, targetFloor int, food int) StateMessage {
	reply := NewStateFoodTakenMessage(senderID, senderFloor, targetFloor, food)
	return reply
}

func (msg *AskFoodTakenMessage) Visit(a Agent) {
	a.HandleAskFoodTaken(*msg)
}
