package messages

import "github.com/google/uuid"

type ProposeTreatyMessage struct {
	*BaseMessage
	treaty Treaty
}

func NewProposalMessage(senderID uuid.UUID, senderFloor int, treaty Treaty) *ProposeTreatyMessage {
	msg := &ProposeTreatyMessage{
		NewBaseMessage(senderID, senderFloor, ProposeTreaty),
		treaty,
	}
	return msg
}

func (msg *ProposeTreatyMessage) Treaty() Treaty {
	return msg.treaty
}

func (msg *ProposeTreatyMessage) Visit(a Agent) {
	a.HandleProposeTreaty(*msg)
}

func (msg *ProposeTreatyMessage) Reply(senderID uuid.UUID, senderFloor int, response bool) TreatyResponseMessage {
	reply := *NewTreatyResponseMessage(senderID, senderFloor, response, msg.treaty.Id(), msg.ID())
	return reply
}
