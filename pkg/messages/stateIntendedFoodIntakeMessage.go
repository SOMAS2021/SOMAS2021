package messages

type StateIntendedFoodIntakeMessage struct {
	*BaseMessage
	intendedFood int
}

func NewStateIntendedFoodIntakeMessage(SenderFloor int, intendedFood int, id string) *StateIntendedFoodIntakeMessage {
	msg := &StateIntendedFoodIntakeMessage{
		NewBaseMessage(SenderFloor, StateIntendedFoodIntake, id),
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
