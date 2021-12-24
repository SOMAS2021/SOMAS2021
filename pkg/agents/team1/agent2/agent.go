package agent2

import (
	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
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
	a.TakeFood(15)
	receivedMsg := a.Base.ReceiveMessage()
	if receivedMsg != nil {
		receivedMsg.Visit(a)
	} else {
		a.Log("I got no thing")
	}

}

func (a *CustomAgent2) HandleAskHP(msg infra.AskMessage) {
	reply := msg.Reply(a.Floor(), a.HP())
	a.SendMessage(msg.SenderFloor()-a.Floor(), reply)
	a.Log("I recieved an askHP message from ", infra.Fields{"floor": msg.SenderFloor()})
}

func (a *CustomAgent2) HandleAskFoodTaken(msg infra.AskMessage) {
	reply := msg.Reply(a.Floor(), 10)
	a.SendMessage(msg.SenderFloor()-a.Floor(), reply)
	a.Log("I recieved an askFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor()})
}

func (a *CustomAgent2) HandleAskIntendedFoodTaken(msg infra.AskMessage) {
	reply := msg.Reply(a.Floor(), 11)
	a.SendMessage(msg.SenderFloor()-a.Floor(), reply)
	a.Log("I recieved an askIntendedFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor()})
}

func (a *CustomAgent2) HandleRequestLeaveFood(msg infra.RequestMessage) {
	reply := msg.Reply(a.Floor(), true)
	a.SendMessage(msg.SenderFloor()-a.Floor(), reply)
	a.Log("I recieved a requestLeaveFood message from ", infra.Fields{"floor": msg.SenderFloor()})
}

func (a *CustomAgent2) HandleRequestTakeFood(msg infra.RequestMessage) {
	reply := msg.Reply(a.Floor(), true)
	a.SendMessage(msg.SenderFloor()-a.Floor(), reply)
	a.Log("I recieved a requestTakeFood message from ", infra.Fields{"floor": msg.SenderFloor()})
}

func (a *CustomAgent2) HandleResponse(msg infra.ResponseMessage) {
	response := msg.Response()
	a.Log("I recieved a Response message from ", infra.Fields{"floor": msg.SenderFloor(), "response": response})
}

func (a *CustomAgent2) HandleStateFoodTaken(msg infra.StateMessage) {
	statement := msg.Statement()
	a.Log("I recieved a StateFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor(), "statement": statement})
}

func (a *CustomAgent2) HandleStateHP(msg infra.StateMessage) {
	statement := msg.Statement()
	a.Log("I recieved a StateHP message from ", infra.Fields{"floor": msg.SenderFloor(), "statement": statement})
}

func (a *CustomAgent2) HandleStateIntendedFoodTaken(msg infra.StateMessage) {
	statement := msg.Statement()
	a.Log("I recieved a StateIntendedFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor(), "statement": statement})
}
