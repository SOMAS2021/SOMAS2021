package messages

import (
	"github.com/google/uuid"
)

type ConditionType int
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
	conditionValue int
	request        RequestType
	requestValue   int
	conditionOp    Op
	requestOp      Op
	signatureCount int
	duration       int
	id             uuid.UUID
	proposerID     uuid.UUID
}

type Treatyer interface {
	Condition() ConditionType
	ConditionValue() int
	Request() RequestType
	RequestValue() int
	ConditionOp() Op
	RequestOp() Op
	SignatureCount() int
	ProposerID() uuid.UUID
	Duration() int
	ID() uuid.UUID
	SignTreaty()
	SetCount(count int)
	DecrementDuration()
}

func NewTreaty(condition ConditionType, conditionValue int, request RequestType, requestValue int, cop Op, rop Op, duration int, proposerID uuid.UUID) *Treaty {
	treaty := &Treaty{
		condition:      condition,
		conditionValue: conditionValue,
		request:        request,
		requestValue:   requestValue,
		conditionOp:    cop,
		requestOp:      rop,
		signatureCount: 1,
		duration:       duration,
		id:             uuid.New(),
		proposerID:     proposerID,
	}
	return treaty
}

func (t *Treaty) Condition() ConditionType {
	return t.condition
}

func (t *Treaty) ConditionValue() int {
	return t.conditionValue
}

func (t *Treaty) Request() RequestType {
	return t.request
}

func (t *Treaty) RequestValue() int {
	return t.requestValue
}

func (t *Treaty) ConditionOp() Op {
	return t.conditionOp
}

func (t *Treaty) ProposerID() uuid.UUID {
	return t.proposerID
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

func (t *Treaty) ID() uuid.UUID {
	return t.id
}

func (t *Treaty) SignTreaty() {
	t.signatureCount++
}

func (t *Treaty) SetCount(count int) {
	t.signatureCount = count
}

func (t *Treaty) DecrementDuration() {
	t.duration--
}
