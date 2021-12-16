package team3

import (
	"math/rand"
	"time"

	"github.com/SOMAS2021/SOMAS2021/pkg/infra"

	//"github.com/SOMAS2021/SOMAS2021/pkg/messages"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/abm"
)

type team3Variables struct {
	//Stubbornnes defines the likelyhood of reading a message
	stubbornness int
	//Morality defines the willigness to help others/ how much you care
	morality int
	//Mood affects the decision making and how you take things
	mood int
}

type team3Knowledge struct {
	//We know the floors we have been in
	floors []int
}

type CustomAgent3 struct {
	*infra.Base
	vars      team3Variables
	knowledge team3Knowledge
	//and an array of tuples for friendships
}

func New(baseAgent *infra.Base) (abm.Agent, error) {
	rand.Seed(time.Now().UnixNano())
	return &CustomAgent3{
		Base: baseAgent,
		vars: team3Variables{
			//random seed
			stubbornness: rand.Intn(75),
			morality:     rand.Intn(100),
			mood:         rand.Intn(100),
		},
	}, nil
}

func (a *CustomAgent3) Run() {

	if a.knowledge.floors[0] != a.Floor() || len(a.knowledge.floors) == 0 {
		changeNewFloor(a)
	}

	//IF HP changes (but we have to remember previous HP)
	// changeNewDay(a)

	//receive Message

	//eat
	a.TakeFood(float64(takeFoodCalculation(a)))

	//send Message

}
