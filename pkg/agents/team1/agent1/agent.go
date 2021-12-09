package agent1

import (
	"log"

	infra "github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
	abm "github.com/SOMAS2021/SOMAS2021/pkg/utils/abm"
)

type CustomAgent1 struct {
	*infra.Base
	myNumber int
}

func New(baseAgent *infra.Base) (abm.Agent, error) {
	return &CustomAgent1{
		Base:     baseAgent,
		myNumber: 0,
	}, nil
}

func (a *CustomAgent1) Run() {
	log.Printf("Custom agent in team 1 is on floor %d has hp: %d", a.Floor(), a.HP())

	receivedMsg := a.Base.ReceiveMessage()
	if receivedMsg != nil {
		log.Printf("%d, %d I got a msg from: %s", a.Floor(), a.myNumber, receivedMsg.MessageType())
	} else {
		log.Printf("%d, %d, I got nothing", a.Floor(), a.myNumber)
	}

	if (a.myNumber)%2 == 0 {
		msg := *messages.NewAckMessage(int(a.Floor()), true)
		a.SendMessage(1, msg)
		log.Printf("%d,%d, I sent a msg: %s", a.Floor(), a.myNumber, msg.MessageType())
	} else {
		msg := *messages.NewBaseMessage(int(a.Floor()))
		a.SendMessage(1, msg)
	}

	a.takeFood(10)
	a.HP()
}

func (a *CustomAgent1) HP() int {
	return a.Base.HP()
}

func (a *CustomAgent1) Die() {
	a.Base.Die()
}

func (a *CustomAgent1) IsAlive() bool {
	return a.Base.IsAlive()
}

func (a *CustomAgent1) takeFood(foodToTake float64) float64 {
	return a.Base.TakeFood(foodToTake)
}
