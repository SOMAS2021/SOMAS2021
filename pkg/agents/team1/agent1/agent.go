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

	r := rand.Intn(9)
	switch r {
	case 0:
		msg := messages.NewAskFoodTakenMessage(a.ID(), a.Floor(), a.Floor()-2)
		a.SendMessage(msg)
		a.Log("I sent a message", infra.Fields{"message": "AskFoodTaken", "floor": a.Floor()})
	case 1:
		msg := messages.NewAskHPMessage(a.ID(), a.Floor(), a.Floor()+3)
		a.SendMessage(msg)
		a.Log("I sent a message", infra.Fields{"message": "AskHP", "floor": a.Floor()})
	case 2:
		msg := messages.NewAskIntendedFoodIntakeMessage(a.ID(), a.Floor(), a.Floor()+2)
		a.SendMessage(msg)
		a.Log("I sent a message", infra.Fields{"message": "AskIntendedFoodIntake", "floor": a.Floor()})
	case 3:
		msg := messages.NewRequestLeaveFoodMessage(a.ID(), a.Floor(), a.Floor()+1, 10)
		a.SendMessage(msg)
		a.Log("I sent a message", infra.Fields{"message": "RequestLeaveFood", "floor": a.Floor()})
	case 4:
		msg := messages.NewRequestTakeFoodMessage(a.ID(), a.Floor(), a.Floor()-3, 20)
		a.SendMessage(msg)
		a.Log("I sent a message", infra.Fields{"message": "RequestTakeFood", "floor": a.Floor()})
	case 5:
		msg := messages.NewStateFoodTakenMessage(a.ID(), a.Floor(), a.Floor()+2, 30)
		a.SendMessage(msg)
		a.Log("I sent a message", infra.Fields{"message": "StateFoodTaken", "floor": a.Floor()})
	case 6:
		msg := messages.NewStateHPMessage(a.ID(), a.Floor(), a.Floor()-1, 40)
		a.SendMessage(msg)
		a.Log("I sent a message", infra.Fields{"message": "StateHP", "floor": a.Floor()})
	case 7:
		msg := messages.NewStateIntendedFoodIntakeMessage(a.ID(), a.Floor(), a.Floor()+3, 50)
		a.SendMessage(msg)
		a.Log("I sent a message", infra.Fields{"message": "StateIntendedFoodIntake", "floor": a.Floor()})
	case 8:
		treaty := messages.NewTreaty(messages.HP, 100, messages.LeaveAmountFood, 30, messages.GT, messages.EQ, 1, a.ID())
		msg := messages.NewProposalMessage(a.ID(), a.Floor(), a.Floor()+1, *treaty)
		a.SendMessage(msg)
		a.Log("I sent a message", infra.Fields{"message": "ProposeTreatyMessage", "floor": a.Floor()})
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
	reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), a.HP())
	a.SendMessage(reply)
	a.Log("I recieved an askHP message from ", infra.Fields{"senderFloor": msg.SenderFloor(), "myFloor": a.Floor()})
}

func (a *CustomAgent1) HandleAskFoodTaken(msg messages.AskFoodTakenMessage) {
	reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), 10)
	a.SendMessage(reply)
	a.Log("I recieved an askFoodTaken message from ", infra.Fields{"senderFloor": msg.SenderFloor(), "myFloor": a.Floor()})
}

func (a *CustomAgent1) HandleAskIntendedFoodTaken(msg messages.AskIntendedFoodIntakeMessage) {
	reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), 11)
	a.SendMessage(reply)
	a.Log("I recieved an askIntendedFoodTaken message from ", infra.Fields{"senderFloor": msg.SenderFloor(), "myFloor": a.Floor()})
}

func (a *CustomAgent1) HandleRequestLeaveFood(msg messages.RequestLeaveFoodMessage) {
	reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), true)
	a.SendMessage(reply)
	a.Log("I recieved a requestLeaveFood message from ", infra.Fields{"senderFloor": msg.SenderFloor(), "myFloor": a.Floor()})
}

func (a *CustomAgent1) HandleRequestTakeFood(msg messages.RequestTakeFoodMessage) {
	reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), true)
	a.SendMessage(reply)
	a.Log("I recieved a requestTakeFood message from ", infra.Fields{"senderFloor": msg.SenderFloor(), "myFloor": a.Floor()})
}

func (a *CustomAgent1) HandleResponse(msg messages.BoolResponseMessage) {
	response := msg.Response()
	a.Log("I recieved a Response message from ", infra.Fields{"senderID": msg.SenderID(), "senderFloor": msg.SenderFloor(), "response": response, "myFloor": a.Floor()})
}

func (a *CustomAgent1) HandleStateFoodTaken(msg messages.StateFoodTakenMessage) {
	statement := msg.Statement()
	a.Log("I recieved a StateFoodTaken message from ", infra.Fields{"senderFloor": msg.SenderFloor(), "statement": statement, "myFloor": a.Floor()})
}

func (a *CustomAgent1) HandleStateHP(msg messages.StateHPMessage) {
	statement := msg.Statement()
	a.Log("I recieved a StateHP message from ", infra.Fields{"senderFloor": msg.SenderFloor(), "statement": statement, "myFloor": a.Floor()})
}

func (a *CustomAgent1) HandleStateIntendedFoodTaken(msg messages.StateIntendedFoodIntakeMessage) {
	statement := msg.Statement()
	a.Log("I recieved a StateIntendedFoodTaken message from ", infra.Fields{"senderFloor": msg.SenderFloor(), "statement": statement, "myFloor": a.Floor()})
}
