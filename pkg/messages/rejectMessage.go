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

// Can be used as a reply to any message -- probably lowers morale of receiver
// Has more of an effect than just ignoring (ackMessage and nothing else)

