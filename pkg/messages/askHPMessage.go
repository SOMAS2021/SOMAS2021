package messages

type AskHPMessage struct {
	baseMessage *BaseMessage
}

func NewAskHPMessage(senderFloor int) *AskHPMessage {
	msg := &AskHPMessage{
		baseMessage: NewBaseMessage(senderFloor, AskHP),
	}
	return msg
}

func (msg *AskHPMessage) Reply(senderFloor int, hp int) StateMessage {
	reply := NewStateHPMessage(senderFloor, hp)
	return reply
}

func (msg *AskHPMessage) MessageType() MessageType {
	return msg.baseMessage.MessageType()
}

func (msg *AskHPMessage) SenderFloor() int {
	return msg.baseMessage.SenderFloor()
}

func (msg *AskHPMessage) Visit(a Agent) {
	a.HandleAskHP(*msg)
}
