package messages

type ProposeTreatyMessage struct {
	*BaseMessage
	treaty *Treaty
}

func NewProposalMessage(senderFloor int, treaty *Treaty) *ProposeTreatyMessage {
	msg := &ProposeTreatyMessage{
		NewBaseMessage(senderFloor, ProposeTreaty),
		treaty,
	}
	return msg
}

func (msg *ProposeTreatyMessage) Treaty() *Treaty {
	return msg.treaty
}

func (msg *ProposeTreatyMessage) Visit(a Agent) {
	a.HandleProposeTreaty(*msg)
}

func (msg *ProposeTreatyMessage) Reply(senderFloor int, response bool) ResponseMessage {
	reply := NewResponseMessage(senderFloor, response)
	return reply
}
