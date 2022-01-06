package messages

import (
	"strconv"

	"github.com/google/uuid"
)

type StateIntendedFoodIntakeMessage struct {
	*BaseMessage
	intendedFood int
}

func NewStateIntendedFoodIntakeMessage(senderID uuid.UUID, senderFloor int, targetFloor int, intendedFood int) *StateIntendedFoodIntakeMessage {
	msg := &StateIntendedFoodIntakeMessage{
		NewBaseMessage(senderID, senderFloor, targetFloor, StateIntendedFoodIntake),
		intendedFood,
	}
	return msg
}

func (msg *StateIntendedFoodIntakeMessage) Statement() int {
	return msg.intendedFood
}

func (msg *StateIntendedFoodIntakeMessage) Visit(a Agent) {
	if msg.TargetFloor() != a.Floor() {
		a.HandlePropogate(msg)
	} else {
		a.HandleStateIntendedFoodTaken(*msg)
	}
}

func (msg *StateIntendedFoodIntakeMessage) StoryLog() string {
	return strconv.Itoa(msg.intendedFood)
}
