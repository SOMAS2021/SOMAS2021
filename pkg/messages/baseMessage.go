package messages

import "github.com/google/uuid"

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
	ID() uuid.UUID
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
	Reply(senderFloor int, response bool) ResponseMessage
}

type ResponseMessage interface {
	Message
	Response() bool
	RequestId() uuid.UUID
}

type BaseMessage struct {
	senderID    uuid.UUID
	senderFloor int
	messageType MessageType
	id          uuid.UUID
}

func NewBaseMessage(senderID uuid.UUID, senderFloor int, messageType MessageType) *BaseMessage {
	msg := &BaseMessage{
		senderID:    senderID,
		senderFloor: senderFloor,
		messageType: messageType,
		id:          uuid.New(),
	}
	return msg
}

func (msg BaseMessage) MessageType() MessageType {
	return msg.messageType
}

func (msg BaseMessage) SenderFloor() int {
	return msg.senderFloor
}

func (msg BaseMessage) ID() uuid.UUID {
	return msg.id
}
func (msg BaseMessage) SenderID() uuid.UUID {
	return msg.senderID
}
