package messages

type StateIntendedFoodIntakeMessage struct {
	*BaseMessage
	intendedFood int
}

func NewStateIntendedFoodIntakeMessage(SenderFloor int, intendedFood int) *StateIntendedFoodIntakeMessage {
	msg := &StateIntendedFoodIntakeMessage{
		NewBaseMessage(SenderFloor, StateIntendedFoodIntake),
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
