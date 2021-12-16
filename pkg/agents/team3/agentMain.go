package team3

import (
	"math/rand"

	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	//"github.com/SOMAS2021/SOMAS2021/pkg/messages"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/abm"
)

type team3Variables struct {
	//Stubbornnes defines the likelyhood of reading a message
	stubbornness int
	//The willigness to help others/ how much you care
	morale int
	//A more volatile parameter that affects the decision making
	mood int
}

type CustomAgent3 struct {
	*infra.Base
	vars team3Variables
	//for later: add a array of ints to remember floors
	//and an array of tuples for friendships
}

func New(baseAgent *infra.Base) (abm.Agent, error) {
	return &CustomAgent3{
		Base: baseAgent,
		vars: team3Variables{
			stubbornness: rand.Intn(75),
			morale:       rand.Intn(100),
			mood:         rand.Intn(100),
		},
	}, nil
}

func (a *CustomAgent3) Run() {
	a.Log("Reporting agent state", infra.Fields{"health": a.HP(), "floor": a.Floor()})

	receivedMsg := a.Base.ReceiveMessage()
	if receivedMsg != nil {
		a.Log("I sent a message", infra.Fields{"message": receivedMsg.MessageType()})
	} else {
		a.Log("I got nothing")
	}

	a.Log("My agent is doing something", infra.Fields{"thing": "potatoe", "another_thing": "another potatoe"})
	a.TakeFood(10)
}
