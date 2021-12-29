package messages

type StateHPMessage struct {
	*BaseMessage
	hp int
}

func NewStateHPMessage(senderFloor int, hp int, id string) *StateHPMessage {
	msg := &StateHPMessage{
		NewBaseMessage(senderFloor, StateHP, id),
		hp,
	}
	return msg
}

func (msg *StateHPMessage) Statement() int {
	return msg.hp
}

func (msg *StateHPMessage) Visit(a Agent) {
	a.HandleStateHP(*msg)
}
