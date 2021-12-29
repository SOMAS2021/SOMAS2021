package messages

import (
	"github.com/google/uuid"
)

//Define message types to enable basic protocols, voting systems ...etc

type MessageType int

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

type Agent interface {
	HandleAskHP(msg AskHPMessage)
	HandleAskFoodTaken(msg AskFoodTakenMessage)
	HandleAskIntendedFoodTaken(msg AskIntendedFoodIntakeMessage)
	HandleRequestLeaveFood(msg RequestLeaveFoodMessage)
	HandleRequestTakeFood(msg RequestTakeFoodMessage)
	HandleResponse(msg BoolResponseMessage)
	HandleStateFoodTaken(msg StateFoodTakenMessage)
	HandleStateHP(msg StateHPMessage)
	HandleStateIntendedFoodTaken(msg StateIntendedFoodIntakeMessage)
}

type Message interface {
	MessageType() MessageType
	SenderFloor() int
	ID() string
	Visit(a Agent)
}

type AskMessage interface {
	Message
	Reply(senderFloor int, food int) StateMessage
}

type StateMessage interface {
	Message
	Statement() int
}

type RequestMessage interface {
	Message
	Request() int
	Reply(senderFloor int, response bool, uuid string) ResponseMessage
}

type ResponseMessage interface {
	Message
	Response() bool
}

type BaseMessage struct {
	senderFloor int
	messageType MessageType
	id          string
}

func NewBaseMessage(senderFloor int, messageType MessageType, id string) *BaseMessage {
	if id == "" { //empty string to generate new ids, else use messageID to reply
		id = uuid.New().String()
	}
	msg := &BaseMessage{
		senderFloor: senderFloor,
		messageType: messageType,
		id:          id,
	}
	return msg
}

func (msg BaseMessage) MessageType() MessageType {
	return msg.messageType
}

func (msg BaseMessage) SenderFloor() int {
	return msg.senderFloor
}

func (msg BaseMessage) ID() string {
	return msg.id
}
