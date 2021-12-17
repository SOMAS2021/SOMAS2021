package agent1

import (
	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/agent"
)

type CustomAgent1 struct {
	*infra.Base
	myNumber int
}

func New(baseAgent *infra.Base) (agent.Agent, error) {
	return &CustomAgent1{
		Base:     baseAgent,
		myNumber: 0,
	}, nil
}

func (a *CustomAgent1) Run() {
	a.Log("Reporting agent state", infra.Fields{"health": a.HP(), "floor": a.Floor()})

	// receivedMsg := a.Base.ReceiveMessage()
	// if receivedMsg != nil {
	// 	a.Log("I sent a message", infra.Fields{"message": receivedMsg.MessageType()})
	// } else {
	// 	a.Log("I got nothing")
	// }

	// if (a.myNumber)%2 == 0 {
	// 	msg := *messages.NewAckMessage(int(a.Floor()), true)
	// 	a.SendMessage(1, msg)
	// 	a.Log("I sent a message", infra.Fields{"message": msg.MessageType()})
	// } else {
	// 	msg := *messages.NewBaseMessage(int(a.Floor()))
	// 	a.SendMessage(1, msg)
	// }
	// a.Log("My agent is doing something", infra.Fields{"thing": "potatoe", "another_thing": "another potatoe"})
	a.TakeFood(40)
}
