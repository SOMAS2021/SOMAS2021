package infra

import "github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/agent"

//Define message types to enable basic protocols, voting systems ...etc

type MessageType int

// woonmoon(TODO): See if there's a more Go-friendly way to do this...
// ideally Ask/State/Request would be three different Messagetypes
const (
	AskFoodTaken MessageType = iota + 1
	AskHP
	AskFoodOnPlatform
	AskIntendedFoodIntake
	AskIdentity
	StateFoodTaken
	StateHP
	StateFoodOnPlatform
	StateIntendedFoodIntake
	StateIdentity
	StateResponse
	RequestLeaveFood
	RequestTakeFood
	Response
)

type Message interface {
	MessageType() MessageType
	SenderFloor() int
	Visit(agent.Agent)
}

type AskMessage interface {
	Message
	Reply(senderFloor int, food float64) StateMessage
}

type StateMessage interface {
	Message
	Statement() float64
}

type RequestMessage interface {
	Message
	Request() float64
	Reply(senderFloor int, response bool) ResponseMessage
}

type ResponseMessage interface {
	Message
	Response() bool
}

type BaseMessage struct {
	senderFloor int
	messageType MessageType
}

func NewBaseMessage(senderFloor int, messageType MessageType) *BaseMessage {
	msg := &BaseMessage{
		senderFloor: senderFloor,
		messageType: messageType,
	}
	return msg
}

func (msg BaseMessage) MessageType() MessageType {
	return msg.messageType
}

func (msg BaseMessage) SenderFloor() int {
	return msg.senderFloor
}
