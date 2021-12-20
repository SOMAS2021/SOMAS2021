package messages

type whoAreYouMessage struct {
	*baseMessage
	ID string
} 

func NewwhoAreYouMessage(SenderFloor int, id string) *whoAreYouMessage {
	msg := &whoAreYouMessage{
		baseMessage: NewBaseMessage(SenderFloor),
		ID: id,
	}
	return msg
}

func (msg whoAreYouMessage) MessageType() string {
	return "whoAreYouMessage"
	
}