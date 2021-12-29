package messages

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
	HandleProposeTreaty(msg ProposeTreatyMessage)
}

type Message interface {
	MessageType() MessageType
	SenderFloor() int
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

type ProposalMessage interface {
	Message
	Treaty() Treaty
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
