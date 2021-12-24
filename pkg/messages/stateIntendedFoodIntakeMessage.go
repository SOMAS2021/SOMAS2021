package messages

type StateIntendedFoodIntakeMessage struct {
	baseMessage  *BaseMessage
	intendedFood int
}

func NewStateIntendedFoodIntakeMessage(SenderFloor int, intendedFood int) *StateIntendedFoodIntakeMessage {
	msg := &StateIntendedFoodIntakeMessage{
		baseMessage:  NewBaseMessage(SenderFloor, StateIntendedFoodIntake),
		intendedFood: intendedFood,
	}
	return msg
}

func (msg *StateIntendedFoodIntakeMessage) Statement() int {
	return msg.intendedFood
}

func (msg *StateIntendedFoodIntakeMessage) MessageType() MessageType {
	return msg.baseMessage.MessageType()
}

func (msg *StateIntendedFoodIntakeMessage) SenderFloor() int {
	return msg.baseMessage.SenderFloor()
}

func (msg *StateIntendedFoodIntakeMessage) Visit(a Agent) {
	a.HandleStateIntendedFoodTaken(*msg)
}
