package messages

type AskHPMessage struct {
	*BaseMessage
}

func NewAskHPMessage(senderFloor int) *AskHPMessage {
	msg := &AskHPMessage{
		NewBaseMessage(senderFloor, AskHP, ""),
	}
	return msg
}

func (msg *AskHPMessage) Reply(senderFloor int, hp int, id string) StateMessage {
	reply := NewStateHPMessage(senderFloor, hp, id)
	return reply
}

func (msg *AskHPMessage) Visit(a Agent) {
	a.HandleAskHP(*msg)
}
