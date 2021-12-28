package team3

import (
	"math/rand"

	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
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
}

type CustomAgent3 struct {
	*infra.Base
	vars      team3Variables
	knowledge team3Knowledge
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
			floors: []int{},
			lastHP: 100,
		},
	}, nil
}

func (a *CustomAgent3) Run() {
	//Update agent variables at the beginning of day (when HP has been reduced)
	if a.HP() < a.knowledge.lastHP {
		changeNewDay(a)
	}

	//Update agent variables at the beginning of reshuffle (when floor has changed)
	if len(a.knowledge.floors) == 0 || a.knowledge.floors[len(a.knowledge.floors)-1] != a.Floor() {
		changeNewFloor(a)
	}
	a.Log("Custom agent 3 each run:", infra.Fields{"floor": a.Floor(), "hp": a.HP(), "Mood": a.vars.mood, "Morality": a.vars.morality})

	//receive Message

	//eat
	a.TakeFood(takeFoodCalculation(a))

	//send Message

}
