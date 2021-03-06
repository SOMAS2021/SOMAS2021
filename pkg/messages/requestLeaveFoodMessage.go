package messages

import (
	"strconv"

	"github.com/google/uuid"
)

type RequestLeaveFoodMessage struct {
	*BaseMessage
	food int
}

func NewRequestLeaveFoodMessage(senderID uuid.UUID, senderFloor int, targetFloor int, food int) *RequestLeaveFoodMessage {
	msg := &RequestLeaveFoodMessage{
		NewBaseMessage(senderID, senderFloor, targetFloor, RequestLeaveFood),
		food,
	}
	return msg
}

func (msg *RequestLeaveFoodMessage) Request() int {
	return msg.food
}

func (msg *RequestLeaveFoodMessage) Reply(senderID uuid.UUID, senderFloor int, targetFloor int, response bool) ResponseMessage {
	reply := NewResponseMessage(senderID, senderFloor, targetFloor, response, msg.ID())
	return reply
}

func (msg *RequestLeaveFoodMessage) Visit(a Agent) {
	if msg.TargetFloor() != a.Floor() {
		a.HandlePropogate(msg)
	} else {
		a.HandleRequestLeaveFood(*msg)
	}
}

func (msg *RequestLeaveFoodMessage) StoryLog() string {
	return strconv.Itoa(msg.food)
}
