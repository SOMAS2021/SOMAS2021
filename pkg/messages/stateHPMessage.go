package messages

// TODO (woonmoon): HP SHOULD BE AN INT!!!
// do this by supporting union return types on Statement() in
// AskMessage interface
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
	// TODO (woonmoon): GET RID OF THIS CONVERSION by supporting pseudo-union return types of Statement() function
	return msg.hp
}

func (msg *StateHPMessage) Visit(a Agent) {
	a.HandleStateHP(*msg)
}
