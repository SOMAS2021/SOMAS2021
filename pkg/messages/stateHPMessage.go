package messages

type StateHPMessage struct {
	*BaseMessage
	hp int
}

func NewStateHPMessage(senderFloor int, hp int) *StateHPMessage {
	msg := &StateHPMessage{
		NewBaseMessage(senderFloor, StateHP),
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
