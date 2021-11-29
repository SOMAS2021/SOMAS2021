package Messages

//Define message types to enable basic protocols, voting systems ...etc
type Message struct {
	// senderID
	// content
}
type AckMessage struct {
	msg Message
	ack bool
}
