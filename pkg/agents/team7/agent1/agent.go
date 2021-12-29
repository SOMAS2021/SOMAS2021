package team7agent1

import (
	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
	"math/rand"
	"time"
)

/*

Openness 	high-> doesn't get effected by new floors, looks forward to new day
			low-> does get effected by new floors

Currently implemented by changing attributes on reshuffles

Conscientiousness high-> plans ahead, attentive and takes into account messgages
			low-> no planning, fails to complete assigned tasks

Not implemented yet


Extraversion high-> very likely to communicate, will share a lot of information
			low-> does not like to communicate

Not implemented yet

Agreeableness high-> caring so more likely to go hungry for the betterment of others, trustworthy
			low-> greedy, manipulative


Implemented by intialising attributes, not considering health state yet


Neuroticism high-> dramatic mood swings, struggles to recover after weak state
			low-> stable, relaxed and resilient




*/

type team7Personalities struct {
	openness int

	conscientiousness int

	extraversion int

	agreeableness int

	neuroticism int
}

type CustomAgent7 struct {
	*infra.Base
	personality team7Personalities
	greediness  int
	kindness    int
	daysAlive   int
	daysHungry  int
	prevHP      int
	prevFloors  []int
}

func New(baseAgent *infra.Base) (infra.Agent, error) {
	return &CustomAgent7{
		Base: baseAgent,
		personality: team7Personalities{
			openness:          rand.Intn(100),
			conscientiousness: rand.Intn(100),
			extraversion:      rand.Intn(100),
			agreeableness:     rand.Intn(100),
			neuroticism:       rand.Intn(100),
		},
		greediness: 0,
		kindness:   0,
		daysAlive:  0,
		daysHungry: 0,
		prevHP:     100,
		prevFloors: []int{},
	}, nil
}

func (a *CustomAgent7) Run() {

	//initialise greediness and kindness
	if a.daysAlive == 0 {
		rand.Seed(time.Now().UTC().UnixNano())
		a.greediness = 100 - a.personality.agreeableness
		a.kindness = a.personality.agreeableness
	}

	a.Log("Team7Agent1 reporting status:", infra.Fields{"floor": a.Floor(), "hp": a.HP(), "greed": a.greediness, "kind": a.kindness, "aggr": a.personality.agreeableness})

	//Check if new day has started
	if a.HP() < a.prevHP {
		a.daysAlive++
		a.prevHP = a.HP()

	}

	//Check if floor has changed
	if len(a.prevFloors) == 0 || a.prevFloors[len(a.prevFloors)-1] != a.Floor() {

		currentFloor := a.Floor()

		if len(a.prevFloors) != 0 {
			prevFloor := a.prevFloors[len(a.prevFloors)-1]

			if currentFloor < prevFloor {
				//only negatively impacted if openness is low
				if a.personality.openness <= 50 {
					a.greediness += (currentFloor - prevFloor)
				}

			} else {
				//if we have moved up our kindness increases
				a.kindness += (prevFloor - currentFloor)
			}

		}
		a.prevFloors = append(a.prevFloors, currentFloor)
	}

	//receive Message

	//eat

	//choses random number between -5 and 5
	if a.CurrPlatFood() != -1 {

		r1 := rand.Intn(11) - 10
		r2 := rand.Intn(11) - 10

		a.kindness += (a.personality.neuroticism * r1) / 50
		a.greediness += (a.personality.neuroticism * r2) / 50

		if a.kindness < 0 {
			a.kindness = 0
		} else if a.kindness > 100 {
			a.kindness = 100
		}

		if a.greediness < 0 {
			a.greediness = 0
		} else if a.greediness > 100 {
			a.greediness = 100
		}

		var foodtotake food.FoodType = food.FoodType(100 - a.kindness + a.greediness)

		if foodtotake > 0 {
			a.TakeFood(foodtotake)
		}

	}

	//send Message

}

// type CustomAgent1 struct {
// 	*infra.Base
// 	// new params
// }

// func New(baseAgent *infra.Base) (infra.Agent, error) {
// 	//create other parameters
// 	return &CustomAgent1{
// 		Base: baseAgent,
// 	}, nil
// }

// func (a *CustomAgent1) Run() {
// 	a.Log("Team7Agent1 reporting status:", infra.Fields{"floor": a.Floor(), "hp": a.HP()})

// 	//UserID := a.ID()
// 	currentHP := a.HP()
// 	//currentFloor := a.Floor()
// 	//currentAvailFood := a.CurrPlatFood()

// 	var foodtotake food.FoodType = food.FoodType(100 - currentHP)
// 	if foodtotake != 0 {
// 		a.TakeFood(foodtotake)
// 	}

// }
