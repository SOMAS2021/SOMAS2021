package agent2

import (
	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
)

type CustomAgent2 struct {
	*infra.Base
	myNumber int
	// new params
}

func New(baseAgent *infra.Base) (infra.Agent, error) {
	return &CustomAgent2{
		Base:     baseAgent,
		myNumber: 0,
	}, nil
}

func (a *CustomAgent2) Run() {
	a.Log("Custom agent reporting status without using fields")
	_, err := a.TakeFood(15)
	if err != nil {
		switch err.(type) {
		case *infra.FloorError:
		case *infra.NegFoodError:
		case *infra.AlreadyEatenError:
		default:
		}
	}
	receivedMsg := a.ReceiveMessage()
	if receivedMsg != nil {
		receivedMsg.Visit(a)
	} else {
		a.Log("I got no thing")
	}

}

func (a *CustomAgent2) HandleAskHP(msg messages.AskHPMessage) {
	reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), a.HP())
	direction := 1
	if a.Floor() > msg.TargetFloor() {
		direction = -1
	}
	a.SendMessage(direction, reply)
	a.Log("I recieved an askHP message from ", infra.Fields{"senderFloor": msg.SenderFloor(), "myFloor": a.Floor()})
}

func (a *CustomAgent2) HandleAskFoodTaken(msg messages.AskFoodTakenMessage) {
	reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), 10)
	direction := 1
	if a.Floor() > msg.TargetFloor() {
		direction = -1
	}
	a.SendMessage(direction, reply)
	a.Log("I recieved an askFoodTaken message from ", infra.Fields{"senderFloor": msg.SenderFloor(), "myFloor": a.Floor()})
}

func (a *CustomAgent2) HandleAskIntendedFoodTaken(msg messages.AskIntendedFoodIntakeMessage) {
	reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), 11)
	direction := 1
	if a.Floor() > msg.TargetFloor() {
		direction = -1
	}
	a.SendMessage(direction, reply)
	a.Log("I recieved an askIntendedFoodTaken message from ", infra.Fields{"senderFloor": msg.SenderFloor(), "myFloor": a.Floor()})
}

func (a *CustomAgent2) HandleRequestLeaveFood(msg messages.RequestLeaveFoodMessage) {
	reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), true)
	direction := 1
	if a.Floor() > msg.TargetFloor() {
		direction = -1
	}
	a.SendMessage(direction, reply)
	a.Log("I recieved a requestLeaveFood message from ", infra.Fields{"senderFloor": msg.SenderFloor(), "myFloor": a.Floor()})
}

func (a *CustomAgent2) HandleRequestTakeFood(msg messages.RequestTakeFoodMessage) {
	reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), true)
	direction := 1
	if a.Floor() > msg.TargetFloor() {
		direction = -1
	}
	a.SendMessage(direction, reply)
	a.Log("I recieved a requestTakeFood message from ", infra.Fields{"senderFloor": msg.SenderFloor(), "myFloor": a.Floor()})
}

func (a *CustomAgent2) HandleResponse(msg messages.BoolResponseMessage) {
	response := msg.Response()
	a.Log("I recieved a Response message from ", infra.Fields{"floor": msg.SenderFloor(), "response": response})
}

func (a *CustomAgent2) HandleStateFoodTaken(msg messages.StateFoodTakenMessage) {
	statement := msg.Statement()
	a.Log("I recieved a StateFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor(), "statement": statement})
}

func (a *CustomAgent2) HandleStateHP(msg messages.StateHPMessage) {
	statement := msg.Statement()
	a.Log("I recieved a StateHP message from ", infra.Fields{"floor": msg.SenderFloor(), "statement": statement})
}

func (a *CustomAgent2) HandleStateIntendedFoodTaken(msg messages.StateIntendedFoodIntakeMessage) {
	statement := msg.Statement()
	a.Log("I recieved a StateIntendedFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor(), "statement": statement})
}
