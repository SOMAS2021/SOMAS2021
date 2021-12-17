package team3

import (
	"math/rand"
	"time"

	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/agent"
	//"github.com/SOMAS2021/SOMAS2021/pkg/messages"
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

func New(baseAgent *infra.Base) (agent.Agent, error) {
	// TODO: Remove this line. See Issue #60.
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	return &CustomAgent3{
		Base: baseAgent,
		vars: team3Variables{
			stubbornness: r1.Intn(75),
			morality:     r1.Intn(100),
			mood:         r1.Intn(100),
		},
		knowledge: team3Knowledge{
			floors: []int{},
			lastHP: 100,
		},
	}, nil
}

func (a *CustomAgent3) Run() {
	//Update agent variables at the beggining of day (when HP has been reduced)
	if a.HP() < a.knowledge.lastHP {
		changeNewDay(a)
	}

	//Update agent variables at the beggining of reshuffle. (when floor has changed)
	if len(a.knowledge.floors) == 0 || a.knowledge.floors[len(a.knowledge.floors)-1] != a.Floor() {
		changeNewFloor(a)
	}
	a.Log("Custom agent 3 each run:", infra.Fields{"floor": a.Floor(), "hp": a.HP(), "Mood": a.vars.mood, "Morality": a.vars.morality})

	//receive Message

	//eat
	a.TakeFood(takeFoodCalculation(a))

	//send Message

}
