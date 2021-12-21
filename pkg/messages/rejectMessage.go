package messages

type rejectMessage struct {
	*baseMessage
} 

func NewrejectMessage(SenderFloor int) *rejectMessage {
	msg := &rejectMessage{
		baseMessage: NewBaseMessage(SenderFloor),
	}
	return msg
}

func (msg rejectMessage) MessageType() string {
	return "rejectMessage"
	
}
