package team3

import (
	"math/rand"

	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
)

type team3Variables struct {
	//Stubborness defines the likelihood of reading a message
	stubbornness int
	//Morality defines the willingness to help others/how much you care
	morality int
	//Mood affects the decision making and how you take things
	mood int
}

type team3Knowledge struct {
	//We know the floors we have been in
	floors []int
	//We know the last HP
	lastHP int
	//We know who we have meet
	friends []string
	//We know if we like or not the people we have met
	friendship []float64
	//Stores who is in the floor bellow
	floorBelow string
	//Stores who is in the floor above
	floorAbove string
}

type team3Decisions struct {
	//Amount of food we decided to eat.
	foodToEat int
	//Amount of food we decided to leave in the platform.
	foodToLeave int
}

type CustomAgent3 struct {
	*infra.Base
	vars      team3Variables
	knowledge team3Knowledge
	decisions team3Decisions
	//and an array of tuples for friendships
}

func New(baseAgent *infra.Base) (infra.Agent, error) {
	return &CustomAgent3{
		Base: baseAgent,
		vars: team3Variables{
			stubbornness: rand.Intn(75),
			morality:     rand.Intn(100),
			mood:         rand.Intn(100),
		},
		knowledge: team3Knowledge{
			floors:     []int{},
			lastHP:     100,
			friends:    []string{},
			friendship: []float64{},
			floorBelow: "",
			floorAbove: "",
		},
		decisions: team3Decisions{
			foodToEat:   -1,
			foodToLeave: -1,
		},
		decisions: team3Decisions{
			foodToEat:   -1,
			foodToLeave: -1,
		},
	}, nil
}

func (a *CustomAgent3) Run() {
	//Update agent variables at the beginning of day (when HP has been reduced)
	if a.HP() < a.knowledge.lastHP { //thanks to the new HP, this will also not work (as somedays we might wake up with the same HP)
		changeNewDay(a)
	}

	//Update agent variables at the beginning of reshuffle (when floor has changed)
	if len(a.knowledge.floors) == 0 || a.knowledge.floors[len(a.knowledge.floors)-1] != a.Floor() {
		changeNewFloor(a)
	}
	a.Log("Custom agent 3 each run:", infra.Fields{"floor": a.Floor(), "hp": a.HP(), "mood": a.vars.mood, "morality": a.vars.morality})

	//Update variables right after eating
	if a.HP() > a.knowledge.lastHP { //this will not work with the new HP function
		a.knowledge.lastHP = a.HP()
		a.decisions.foodToEat = -1
		a.decisions.foodToLeave = -1
	}

	//receive Message
	receivedMsg := a.Base.ReceiveMessage()
	//msg := *messages.NewBaseMessage(int(a.Floor()))
	//sender := receivedMsg.SenderFloor
	if receivedMsg != nil {
		a.Log("I got sent a message", infra.Fields{"message": receivedMsg.MessageType()})
		if receivedMsg.MessageType() == "askFoodTakenMessage" {
			a.Log("I've eaten", infra.Fields{"food": takeFoodCalculation(a)})
			msg := *messages.NewfoodTakenMessage(int(a.Floor()), takeFoodCalculation(a))
			a.SendMessage(-1, msg)
			a.Log("I sent a msg", infra.Fields{"floor": a.Floor(), "message": msg.MessageType()})
		}
	} else {
		a.Log("I got nothing", infra.Fields{"floor": a.Floor()})
	}

	//message(a)

	//eat
	_, err := a.TakeFood(food.FoodType(takeFoodCalculation(a)))
	if err != nil {
		switch err.(type) {
		case *infra.FloorError:
		case *infra.NegFoodError:
		case *infra.AlreadyEatenError:
		default:
		}
	}
	//send Message
	//if msg != *messages.NewBaseMessage(int(a.Floor())) {
	//	a.SendMessage(-1, msg)
	//	a.Log("I sent a msg", infra.Fields{"floor": a.Floor(), "message": msg.MessageType()})
	//}

}
