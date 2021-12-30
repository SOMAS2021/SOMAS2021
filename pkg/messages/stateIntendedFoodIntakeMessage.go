package messages

import "github.com/google/uuid"

type StateIntendedFoodIntakeMessage struct {
	*BaseMessage
	intendedFood int
}

func NewStateIntendedFoodIntakeMessage(senderID uuid.UUID, senderFloor int, intendedFood int) *StateIntendedFoodIntakeMessage {
	msg := &StateIntendedFoodIntakeMessage{
		NewBaseMessage(senderID, senderFloor, StateIntendedFoodIntake),
		intendedFood,
	}
	return msg
}

func (msg *StateIntendedFoodIntakeMessage) Statement() int {
	return msg.intendedFood
}

func (msg *StateIntendedFoodIntakeMessage) Visit(a Agent) {
	a.HandleStateIntendedFoodTaken(*msg)
}
