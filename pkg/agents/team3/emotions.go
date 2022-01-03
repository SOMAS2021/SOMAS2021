package team3

import (
	"math/rand"

	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
	"github.com/google/uuid"
)

// Function decides if we read a message or we don't, depends on our level of stubbornness, return is bool
// Stubborness of 20 means probability of 0.2 that we don't read message

func (a *CustomAgent3) read() bool {
	random := rand.Intn(100)

	return random >= a.vars.stubbornness
}

//Function updates friendship depending on factor change -1 to 1, negative reduces frienship, positive increases
//If they are not in our list, we add them.
func (a *CustomAgent3) updateFriendship(friend uuid.UUID, change float64) {
	level, found := a.knowledge.friends[friend]
	if !found {
		a.knowledge.friends[friend] = 0.4 + (float64(a.vars.morality)/100)*0.2
	} else {
		if change < 0 {
			level = level - change*(level)
		} else {
			level = level + change*(1-level)
		}
		a.knowledge.friends[friend] = level
	}
}

// Function gets as input the mini and max change we want in, direction marks if we want it to go up or down
func changeInMood(a *CustomAgent3, pointsMin, pointsMax, direction int) {
	// TODO: Remove this line. See Issue #60.
	points := rand.Intn(pointsMax-pointsMin) + pointsMin
	if direction < 0 {
		a.vars.mood -= points
		if a.vars.mood < 0 {
			a.vars.mood = 0
		}
	}
	if direction > 0 {
		a.vars.mood += points
		if a.vars.mood > 100 {
			a.vars.mood = 100
		}
	}
}

// Function gets as input the mini and max change we want in, direction marks if we want it to go up or down
func changeInMorality(a *CustomAgent3, pointsMin, pointsMax, direction int) {
	points := rand.Intn(pointsMax-pointsMin) + pointsMin
	if direction < 0 {
		a.vars.morality -= points
		if a.vars.morality < 0 {
			a.vars.morality = 0
		}
	}
	if direction > 0 {
		a.vars.morality += points
		if a.vars.morality > 100 {
			a.vars.morality = 100
		}
	}
}

func changeInStubbornness(a *CustomAgent3, change, direction int) {
	if direction < 0 {
		a.vars.stubbornness -= change
		if a.vars.stubbornness < 0 {
			a.vars.stubbornness = 0
		}
	}
	if direction > 0 {
		a.vars.stubbornness += change
		if a.vars.stubbornness > 75 {
			a.vars.stubbornness = 75
		}
	}
}

// Updates mood and morality when the floor is changed. Called at the start of each day.
//TO DO: Trial this and do sth different.
func changeNewDay(a *CustomAgent3) {
	a.knowledge.lastHP = a.HP()
	a.knowledge.agentAge = a.BaseAgent().Age()
	if a.HP() < 50 {
		changeInMorality(a, 10, 15, -1)
		changeInMood(a, 10, 15, -1)
	} else {
		changeInMorality(a, 10, 15, 1)
		changeInMood(a, 10, 15, 1)
	}
	a.knowledge.foodMovingAvg = food.FoodType(float64(a.knowledge.foodMovingAvg)*0.9 + float64(a.knowledge.foodLastEaten)*0.1)
}

// Function is called when the floor changes, changes the mood when we change floors
func changeNewFloor(a *CustomAgent3) {
	//Restart neighbour functions as we don´t know who is above/below anymore
	a.knowledge.floorAbove = uuid.Nil
	a.knowledge.floorBelow = uuid.Nil

	//Calculate mood based on what floors we have been in
	currentFloor := a.Floor()

	//If we are still in the same floor as before we keep the same mood
	if len(a.knowledge.floors) != 0 {
		lastFloor := a.knowledge.floors[len(a.knowledge.floors)-1]
		beenInHigher := false
		beenInLower := false

		for _, floor := range a.knowledge.floors {
			if floor < currentFloor {
				beenInLower = true
			}
			if floor > currentFloor {
				beenInHigher = true
			}
			if beenInLower && beenInHigher {
				break
			}
		}
		if !beenInHigher {
			changeInMood(a, 5, 10, 1)
		}
		if beenInHigher && currentFloor > lastFloor {
			changeInMood(a, 0, 5, 1)
		}
		if beenInLower && currentFloor < lastFloor {
			changeInMood(a, 0, 5, -1)
		}
		if !beenInLower {
			changeInMood(a, 5, 15, -1) //Situation is specially bad, so we get more depressed
		}

	}
	a.knowledge.floors = append(a.knowledge.floors, currentFloor)
}
