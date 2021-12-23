package messages

import (
	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
)

type StateIntendedFoodIntakeMessage struct {
	baseMessage  *infra.BaseMessage
	intendedFood float64
}

func NewStateIntendedFoodIntake(SenderFloor int, intendedFood float64) *StateIntendedFoodIntakeMessage {
	msg := &StateIntendedFoodIntakeMessage{
		baseMessage:  infra.NewBaseMessage(SenderFloor, infra.StateIntendedFoodIntake),
		intendedFood: intendedFood,
	}
	return msg
}

func (msg *StateIntendedFoodIntakeMessage) Statement() float64 {
	return msg.intendedFood
}

func (msg *StateIntendedFoodIntakeMessage) MessageType() infra.MessageType {
	return msg.baseMessage.MessageType()
}

func (msg *StateIntendedFoodIntakeMessage) SenderFloor() int {
	return msg.baseMessage.SenderFloor()
}

func (msg *StateIntendedFoodIntakeMessage) Visit(a infra.Agent) {
	a.HandleStateIntendedFoodTaken(msg.intendedFood)
}
