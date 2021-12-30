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
	ProposeTreaty
	RequestLeaveFood
	RequestTakeFood
	Response
	TreatyResponse
)

func (m MessageType) String() string {
	switch m {
	case AskFoodTaken:
		return "AskFoodTaken"
	case AskHP:
		return "AskHP"
	case AskFoodOnPlatform:
		return "AskFoodOnPlatform"
	case AskIntendedFoodIntake:
		return "AskIntendedFoodIntake"
	case AskIdentity:
		return "AskIdentity"
	case StateFoodTaken:
		return "StateFoodTaken"
	case StateHP:
		return "StateHP"
	case StateFoodOnPlatform:
		return "StateFoodOnPlatform"
	case StateIntendedFoodIntake:
		return "StateIntendedFoodIntake"
	case StateIdentity:
		return "StateIdentity"
	case StateResponse:
		return "StateResponse"
	case RequestLeaveFood:
		return "RequestLeaveFood"
	case RequestTakeFood:
		return "RequestTakeFood"
	case Response:
		return "Response"
	default:
		return "UNKNOWN"
	}
}

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
	HandleProposeTreaty(msg ProposeTreatyMessage)
	HandleTreatyResponse(msg TreatyResponseMessage)
}

type Message interface {
	MessageType() MessageType
	SenderFloor() int
	ID() uuid.UUID
	Visit(a Agent)
}

type AskMessage interface {
	Message
	Reply(senderID uuid.UUID, senderFloor int, food int) StateMessage
}

type StateMessage interface {
	Message
	Statement() int
}

type RequestMessage interface {
	Message
	Request() int
	Reply(senderID uuid.UUID, senderFloor int, response bool) ResponseMessage
}

type ProposalMessage interface {
	Message
	Treaty() Treaty
	Reply(senderID uuid.UUID, senderFloor int, response bool) TreatyResponseMessage
}

type ResponseMessage interface {
	Message
	Response() bool
	RequestID() uuid.UUID
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
