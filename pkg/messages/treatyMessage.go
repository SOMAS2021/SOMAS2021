package messages

import (
	"github.com/google/uuid"
)

type ConditionType int // type of condition
type RequestType int
type Op int

const (
	HP ConditionType = iota + 1
	Floor
	AvailableFood
)

const (
	LeaveAmountFood RequestType = iota + 1
	LeavePercentFood
	Inform
)

const (
	GT Op = iota + 1
	GE
	EQ
	LE
	LT
)

type Treaty struct {
	condition      ConditionType
	request        RequestType
	conditionOp    Op
	requestOp      Op
	signatureCount int
	duration       int
	id             uuid.UUID
}

type Treatyer interface {
	Condition() ConditionType
	Request() RequestType
	ConditionOp() Op
	RequestOp() Op
	SignatureCount() int
	Duration() int
	Id() uuid.UUID
}

func NewTreaty(condition ConditionType, request RequestType, cop Op, rop Op, duration int) *Treaty {
	treaty := &Treaty{
		condition:      condition,
		request:        request,
		conditionOp:    cop,
		requestOp:      rop,
		signatureCount: 1,
		duration:       duration,
		id:             uuid.New(),
	}
	return treaty
}

func (t *Treaty) Condition() ConditionType {
	return t.condition
}

func (t *Treaty) Request() RequestType {
	return t.request
}

func (t *Treaty) ConditionOp() Op {
	return t.conditionOp
}

func (t *Treaty) RequestOp() Op {
	return t.requestOp
}

func (t *Treaty) SignatureCount() int {
	return t.signatureCount
}

func (t *Treaty) Duration() int {
	return t.duration
}

func (t *Treaty) Id() uuid.UUID {
	return t.id
}

type ProposeTreatyMessage struct {
	*BaseMessage
	treaty *Treaty
}

func NewProposalMessage(senderFloor int, condition ConditionType, request RequestType, cop Op, rop Op, duration int) *ProposeTreatyMessage {
	msg := &ProposeTreatyMessage{
		NewBaseMessage(senderFloor, ProposeTreaty),
		NewTreaty(condition, request, cop, rop, duration),
	}
	return msg
}

func (msg *ProposeTreatyMessage) Treaty() *Treaty {
	return msg.treaty
}

func (msg *ProposeTreatyMessage) Visit(a Agent) {
	a.HandleProposeTreaty(*msg)
}

func (msg *ProposeTreatyMessage) Reply(senderFloor int, response bool) ResponseMessage {
	reply := NewResponseMessage(senderFloor, response)
	return reply
}
