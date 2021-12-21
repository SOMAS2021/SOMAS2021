package messages

type askHPMessage struct{
	*baseMessage
}

func NewaskHPMessage(SenderFloor int) *askHPMessage {
	msg := &askHPMessage{
		baseMessage: NewBaseMessage(SenderFloor),
	}
	return msg
}

func (msg askHPMessage) MessageType() string {
	return "askHPMessage"
	
}