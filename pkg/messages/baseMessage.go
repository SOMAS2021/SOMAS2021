package messages

import "github.com/google/uuid"

//Define message types to enable basic protocols, voting systems ...etc

type MessageType int

//go:generate go run golang.org/x/tools/cmd/stringer -type=MessageType

const (
	AskFoodTaken MessageType = iota + 1
	AskHP
	AskFoodOnPlatform
	AskIntendedFoodIntake
	StateFoodTaken
	StateHP
	StateFoodOnPlatform
	StateIntendedFoodIntake
	StateResponse
	ProposeTreaty
	RequestLeaveFood
	RequestTakeFood
	Response
	TreatyResponse
)

type Agent interface {
	Floor() int
	HandleAskHP(msg AskHPMessage)
	HandleAskFoodTaken(msg AskFoodTakenMessage)
	HandleAskIntendedFoodTaken(msg AskIntendedFoodIntakeMessage)
	HandleRequestLeaveFood(msg RequestLeaveFoodMessage)
	HandleRequestTakeFood(msg RequestTakeFoodMessage)
	HandleResponse(msg BoolResponseMessage)
	HandleStateFoodTaken(msg StateFoodTakenMessage)
	HandleStateHP(msg StateHPMessage)
	HandleStateIntendedFoodTaken(msg StateIntendedFoodIntakeMessage)
	HandleProposeTreaty(msg ProposeTreatyMessage)
	HandleTreatyResponse(msg TreatyResponseMessage)
	HandlePropogate(msg Message)
}

type Message interface {
	MessageType() MessageType
	SenderFloor() int
	TargetFloor() int
	ID() uuid.UUID
	Visit(a Agent)
	StoryLog() string
}

type AskMessage interface {
	Message
	Reply(senderID uuid.UUID, senderFloor int, targetFloor int, food int) StateMessage
}

type StateMessage interface {
	Message
	Statement() int
}

type RequestMessage interface {
	Message
	Request() int
	Reply(senderID uuid.UUID, senderFloor int, targetFloor int, response bool) ResponseMessage
}

type ProposalMessage interface {
	Message
	Treaty() Treaty
	Reply(senderID uuid.UUID, senderFloor int, targetFloor int, response bool) TreatyResponseMessage
}

type ResponseMessage interface {
	Message
	Response() bool
	RequestID() uuid.UUID
}

type BaseMessage struct {
	senderID    uuid.UUID
	senderFloor int
	targetFloor int
	messageType MessageType
	id          uuid.UUID
}

func NewBaseMessage(senderID uuid.UUID, senderFloor int, targetFloor int, messageType MessageType) *BaseMessage {
	msg := &BaseMessage{
		senderID:    senderID,
		senderFloor: senderFloor,
		targetFloor: targetFloor,
		messageType: messageType,
		id:          uuid.New(),
	}
	return msg
}

func (msg *BaseMessage) MessageType() MessageType {
	return msg.messageType
}

func (msg *BaseMessage) SenderFloor() int {
	return msg.senderFloor
}

func (msg *BaseMessage) TargetFloor() int {
	return msg.targetFloor
}

func (msg *BaseMessage) ID() uuid.UUID {
	return msg.id
}

func (msg *BaseMessage) SenderID() uuid.UUID {
	return msg.senderID
}

// Default a message does not have extra state info
// All the info is in the message type
func (msg *BaseMessage) StoryLog() string {
	return ""
}
