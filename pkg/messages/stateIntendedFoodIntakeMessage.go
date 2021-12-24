package messages

import (
	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
)

type StateIntendedFoodIntakeMessage struct {
	baseMessage  *infra.BaseMessage
	intendedFood int
}

func NewStateIntendedFoodIntakeMessage(SenderFloor int, intendedFood int) *StateIntendedFoodIntakeMessage {
	msg := &StateIntendedFoodIntakeMessage{
		baseMessage:  infra.NewBaseMessage(SenderFloor, infra.StateIntendedFoodIntake),
		intendedFood: intendedFood,
	}
	return msg
}

func (msg *StateIntendedFoodIntakeMessage) Statement() int {
	return msg.intendedFood
}

func (msg *StateIntendedFoodIntakeMessage) MessageType() infra.MessageType {
	return msg.baseMessage.MessageType()
}

func (msg *StateIntendedFoodIntakeMessage) SenderFloor() int {
	return msg.baseMessage.SenderFloor()
}

func (msg *StateIntendedFoodIntakeMessage) Visit(a infra.Agent) {
	a.HandleStateIntendedFoodTaken(msg)
}
