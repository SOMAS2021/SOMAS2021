package agent1

import (
	"log"

	"github.com/SOMAS2021/SOMAS2021/pkg/agents"
	"github.com/SOMAS2021/SOMAS2021/pkg/infra/messages"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/abm"
)

type CustomAgent1 struct {
	*agents.Base
	myNumber int
}

func New(baseAgent *agents.Base) (abm.Agent, error) {
	return &CustomAgent1{
		Base:     baseAgent,
		myNumber: 0,
	}, nil
}

func (a *CustomAgent1) Run() {
	// log.Printf("Custom agent in team 1 is on floor %d has hp: %d", a.Floor(), a.HP())
	// a.takeFood(10)
	// a.HP()
	// a.IsDead()

	receivedMsg := a.Base.ReceiveMessage()
	if receivedMsg != nil {
		log.Printf("%d, %d I got a msg from: %s", a.Floor(), a.myNumber, receivedMsg.MessageType())
	} else {
		log.Printf("%d, %d, I got nothing", a.Floor(), a.myNumber)
	}

	//log.Printf("%s, %d", a.ID(), a.Floor())

	if (a.myNumber)%2 == 0 {

		msg := *messages.NewAckMessage(int(a.Floor()), true)

		a.SendMessage(1, msg)

		log.Printf("%d,%d, I sent a msg: %s", a.Floor(), a.myNumber, msg.MessageType())
	} else {

		msg := *messages.NewBaseMessage(int(a.Floor()))

		a.SendMessage(1, msg)
	}

	a.TakeFood(1)
}
