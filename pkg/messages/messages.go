//Define message types to enable basic protocols, voting systems ...etc
type Message struct{
    senderID
    content
}
type ackMessage struct{
    msg Message
    ack bool
}