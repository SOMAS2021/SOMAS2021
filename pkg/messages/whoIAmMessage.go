package messages

type whoIAmMessage struct {
	*baseMessage
	ID string
} 

func NewwhoIAmMessage(SenderFloor int, id string) *whoIAmMessage {
	msg := &whoIAmMessage{
		baseMessage: NewBaseMessage(SenderFloor),
		ID: id,
	}
	return msg
}

func (msg whoIAmMessage) MessageType() string {
	return "whoIAmMessage"
	
}