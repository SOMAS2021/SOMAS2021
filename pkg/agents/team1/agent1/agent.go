package agent1

import (
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

	receivedMsg := a.Base.ReceiveMessage()
	if receivedMsg != nil {
		receivedMsg.Visit(a)
		a.Log("I sent a message", infra.Fields{"message": receivedMsg.MessageType()})
	} else {
		a.Log("I got no thing")
	}

	if (a.myNumber)%2 == 0 {
		msg := messages.NewStateFoodTakenMessage(int(a.Floor()), 50.2)
		a.SendMessage(1, msg)
		a.Log("I sent a message", infra.Fields{"message": msg.MessageType()})
	} else {
		msg := messages.NewAskHPMessage(int(a.Floor()))
		a.SendMessage(1, msg)
	}
	a.Log("My agent is doing something", infra.Fields{"thing": "potatoe", "another_thing": "another potatoe"})
	a.TakeFood(16)
}

func (a *CustomAgent1) HandleAskHP()                              {}
func (a *CustomAgent1) HandleAskFoodTaken()                       {}
func (a *CustomAgent1) HandleAskIntendedFoodTaken()               {}
func (a *CustomAgent1) HandleRequestLeaveFood(request float64)    {}
func (a *CustomAgent1) HandleRequestTakeFood(request float64)     {}
func (a *CustomAgent1) HandleResponse(response bool)              {}
func (a *CustomAgent1) HandleStateFoodTaken(food float64)         {}
func (a *CustomAgent1) HandleStateHP(hp int)                      {}
func (a *CustomAgent1) HandleStateIntendedFoodTaken(food float64) {}
