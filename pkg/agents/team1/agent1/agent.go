package agent1

import (
	"math/rand"

	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
)

type CustomAgent1 struct {
	*infra.Base
	myNumber int
}

func New(baseAgent *infra.Base) (infra.Agent, error) {
	return &CustomAgent1{
		Base:     baseAgent,
		myNumber: 0,
	}, nil
}

func (a *CustomAgent1) Run() {
	a.Log("Reporting agent state", infra.Fields{"health": a.HP(), "floor": a.Floor()})

	receivedMsg := a.ReceiveMessage()
	if receivedMsg != nil {
		receivedMsg.Visit(a)
	} else {
		a.Log("I got no thing")
	}

	r := rand.Intn(8)
	switch r {
	case 0:
		msg := messages.NewAskFoodTakenMessage(a.ID(), a.Floor(), a.Floor()+1)
		a.SendMessage(1, msg)
		a.Log("I sent a message", infra.Fields{"message": "AskFoodTaken"})
	case 1:
		msg := messages.NewAskHPMessage(a.ID(), a.Floor(), a.Floor()+1)
		a.SendMessage(1, msg)
		a.Log("I sent a message", infra.Fields{"message": "AskHP"})
	case 2:
		msg := messages.NewAskIntendedFoodIntakeMessage(a.ID(), a.Floor(), a.Floor()+1)
		a.SendMessage(1, msg)
		a.Log("I sent a message", infra.Fields{"message": "AskIntendedFoodIntake"})
	case 3:
		msg := messages.NewRequestLeaveFoodMessage(a.ID(), a.Floor(), a.Floor()+1, 10)
		a.SendMessage(1, msg)
		a.Log("I sent a message", infra.Fields{"message": "RequestLeaveFood"})
	case 4:
		msg := messages.NewRequestTakeFoodMessage(a.ID(), a.Floor(), a.Floor()+1, 20)
		a.SendMessage(1, msg)
		a.Log("I sent a message", infra.Fields{"message": "RequestTakeFood"})
	case 5:
		msg := messages.NewStateFoodTakenMessage(a.ID(), a.Floor(), a.Floor()+1, 30)
		a.SendMessage(1, msg)
		a.Log("I sent a message", infra.Fields{"message": "StateFoodTaken"})
	case 6:
		msg := messages.NewStateHPMessage(a.ID(), a.Floor(), a.Floor()+1, 40)
		a.SendMessage(1, msg)
		a.Log("I sent a message", infra.Fields{"message": "StateHP"})
	case 7:
		msg := messages.NewStateIntendedFoodIntakeMessage(a.ID(), a.Floor(), a.Floor()+1, 50)
		a.SendMessage(1, msg)
		a.Log("I sent a message", infra.Fields{"message": "StateIntendedFoodIntake"})
	}
	a.Log("My agent is doing something", infra.Fields{"thing": "potatoe", "another_thing": "another potatoe"})
	_, err := a.TakeFood(16)
	if err != nil {
		switch err.(type) {
		case *infra.FloorError:
		case *infra.NegFoodError:
		case *infra.AlreadyEatenError:
		default:
		}
	}
}

func (a *CustomAgent1) HandleAskHP(msg messages.AskHPMessage) {
	reply := msg.Reply(a.ID(), a.Floor(), a.Floor()-1, a.HP())
	a.SendMessage(msg.SenderFloor()-a.Floor(), reply)
	a.Log("I recieved an askHP message from ", infra.Fields{"floor": msg.SenderFloor()})
}

func (a *CustomAgent1) HandleAskFoodTaken(msg messages.AskFoodTakenMessage) {
	reply := msg.Reply(a.ID(), a.Floor(), a.Floor()-1, 10)
	a.SendMessage(msg.SenderFloor()-a.Floor(), reply)
	a.Log("I recieved an askFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor()})
}

func (a *CustomAgent1) HandleAskIntendedFoodTaken(msg messages.AskIntendedFoodIntakeMessage) {
	reply := msg.Reply(a.ID(), a.Floor(), a.Floor()-1, 11)
	a.SendMessage(msg.SenderFloor()-a.Floor(), reply)
	a.Log("I recieved an askIntendedFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor()})
}

func (a *CustomAgent1) HandleRequestLeaveFood(msg messages.RequestLeaveFoodMessage) {
	reply := msg.Reply(a.ID(), a.Floor(), a.Floor()-1, true)
	a.SendMessage(msg.SenderFloor()-a.Floor(), reply)
	a.Log("I recieved a requestLeaveFood message from ", infra.Fields{"floor": msg.SenderFloor()})
}

func (a *CustomAgent1) HandleRequestTakeFood(msg messages.RequestTakeFoodMessage) {
	reply := msg.Reply(a.ID(), a.Floor(), a.Floor()-1, true)
	a.SendMessage(msg.SenderFloor()-a.Floor(), reply)
	a.Log("I recieved a requestTakeFood message from ", infra.Fields{"floor": msg.SenderFloor()})
}

func (a *CustomAgent1) HandleResponse(msg messages.BoolResponseMessage) {
	response := msg.Response()
	a.Log("I recieved a Response message from ", infra.Fields{"senderID": msg.SenderID(), "floor": msg.SenderFloor(), "response": response})
}

func (a *CustomAgent1) HandleStateFoodTaken(msg messages.StateFoodTakenMessage) {
	statement := msg.Statement()
	a.Log("I recieved a StateFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor(), "statement": statement})
}

func (a *CustomAgent1) HandleStateHP(msg messages.StateHPMessage) {
	statement := msg.Statement()
	a.Log("I recieved a StateHP message from ", infra.Fields{"floor": msg.SenderFloor(), "statement": statement})
}

func (a *CustomAgent1) HandleStateIntendedFoodTaken(msg messages.StateIntendedFoodIntakeMessage) {
	statement := msg.Statement()
	a.Log("I recieved a StateIntendedFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor(), "statement": statement})
}
