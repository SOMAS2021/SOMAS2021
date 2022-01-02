package messages

import "github.com/google/uuid"

type ProposeTreatyMessage struct {
	*BaseMessage
	treaty Treaty
}

func NewProposalMessage(senderID uuid.UUID, senderFloor int, targetFloor int, treaty Treaty) *ProposeTreatyMessage {
	msg := &ProposeTreatyMessage{
		NewBaseMessage(senderID, senderFloor, targetFloor, ProposeTreaty),
		treaty,
	}
	return msg
}

func (msg *ProposeTreatyMessage) Treaty() Treaty {
	return msg.treaty
}

func (msg *ProposeTreatyMessage) TreatyID() uuid.UUID {
	return msg.treaty.ID()
}

func (msg *ProposeTreatyMessage) Reply(senderID uuid.UUID, senderFloor int, targetFloor int, response bool) ResponseMessage {
	reply := NewTreatyResponseMessage(senderID, senderFloor, targetFloor, response, msg.treaty.ID(), msg.ID())
	return reply
}

func (msg *ProposeTreatyMessage) Visit(a Agent) {
	if msg.TargetFloor() != a.Floor() {
		a.HandlePropogate(msg)
	} else {
		a.HandleProposeTreaty(*msg)
	}
}
