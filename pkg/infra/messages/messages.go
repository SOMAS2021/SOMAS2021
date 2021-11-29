package messages

//Define message types to enable basic protocols, voting systems ...etc
type Message interface {
	// senderID
	// content
}

type AckMessage struct {
	msg Message
	ack bool
}
