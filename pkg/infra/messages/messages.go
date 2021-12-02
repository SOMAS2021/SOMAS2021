package messages

//Define message types to enable basic protocols, voting systems ...etc
type Message interface {
	GetSenderID()
}

type baseMessage struct {
	senderID uint64
}

type AckMessage struct {
	// msg baseMessage
	// ack bool
}

func (baseMsg *baseMessage) GetSenderID() uint64 {
    return baseMsg.senderID
}
