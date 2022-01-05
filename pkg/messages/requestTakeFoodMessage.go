package messages

import (
	"strconv"

	"github.com/google/uuid"
)

type RequestTakeFoodMessage struct {
	*BaseMessage
	food int
}

func NewRequestTakeFoodMessage(senderID uuid.UUID, senderFloor int, targetFloor int, food int) *RequestTakeFoodMessage {
	msg := &RequestTakeFoodMessage{
		NewBaseMessage(senderID, senderFloor, targetFloor, RequestTakeFood),
		food,
	}
	return msg
}

func (msg *RequestTakeFoodMessage) Request() int {
	return msg.food

}

func (msg *RequestTakeFoodMessage) Reply(senderID uuid.UUID, senderFloor int, targetFloor int, response bool) ResponseMessage {
	reply := NewResponseMessage(senderID, senderFloor, targetFloor, response, msg.ID())
	return reply
}

func (msg *RequestTakeFoodMessage) Visit(a Agent) {
	if msg.TargetFloor() != a.Floor() {
		a.HandlePropogate(msg)
	} else {
		a.HandleRequestTakeFood(*msg)
	}
}

func (msg *RequestTakeFoodMessage) StoryLog() string {
	return strconv.Itoa(msg.food)
}
