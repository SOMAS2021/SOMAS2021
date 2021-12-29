package messages

type TreatyResponseMessage struct {
	*BaseMessage
	response bool
	treaty   *Treaty
}

func NewTreatyResponseMessage(senderFloor int, response bool, treaty *Treaty) *TreatyResponseMessage {
	msg := &TreatyResponseMessage{
		NewBaseMessage(senderFloor, Response),
		response,
		treaty,
	}
	return msg
}

func (msg *TreatyResponseMessage) Response() bool {
	return msg.response
}

func (msg *TreatyResponseMessage) Treaty() *Treaty {
	return msg.treaty
}

func (msg *TreatyResponseMessage) Visit(a Agent) {
	a.HandleTreatyResponse(*msg)
}
