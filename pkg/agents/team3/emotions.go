package team3

import (
	"math/rand"
	"time"
)

// Function decides if we read a message or we don't, depends on our level of stubbornness, return is bool
// Stubborness of 20 means probability of 0.2 that we don't read message

//func read(a *CustomAgent3) bool {
//	rand.Seed(time.Now().UnixNano()) //creates seed to rand
//	var random = rand.Intn(100)
//
//	return random <= a.vars.stubbornness
//}

// Function gets as input the mini and max change we want in, direction marks if we want it to go up or down
func changeInMood(pointsMin, pointsMax, direction int, a *CustomAgent3) {
        // TODO: Remove this line. See Issue #60.
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	points := r1.Intn(pointsMax-pointsMin) + pointsMin
	if direction <= 0 {
		a.vars.mood = a.vars.mood - points
		if temp < 0 {
			a.vars.mood = 0
		}
	}
	if direction > 0 {
		a.vars.mood = a.vars.mood + points
		if a.vars.mood > 100 {
			a.vars.mood = 100
		}
	}
}

// Function gets as input the mini and max change we want in, direction marks if we want it to go up or down
func changeInMorality(pointsMin, pointsMax, direction int, a *CustomAgent3) {
        // TODO: Remove this line. See Issue #60.
	s1 := rand.NewSource(time.Now().UnixNano()) //creates seed to rand
	r1 := rand.New(s1)
	points := r1.Intn(pointsMax-pointsMin) + pointsMin
	if direction < 0 {
		a.vars.morality = a.vars.morality - points
		if a.vars.morality < 0 {
			a.vars.morality = 0
		}
	}
	if direction > 0 {
		a.vars.morality = a.vars.morality + points
		if a.vars.morality > 100 {
			a.vars.morality = 100
		}
	}
}

// Updates mood and morality when the floor is changed. Called at the start of each day.
func changeNewDay(a *CustomAgent3) {
	a.knowledge.lastHP = a.HP()
	if int64(a.HP()) < 50 {
		changeInMorality(10, 15, -1, a)
		changeInMood(10, 15, -1, a)
	} else {
		changeInMorality(10, 15, 1, a)
		changeInMood(10, 15, 1, a)
	}
}

// Function is called when the floor changes, changes the mood when we change floors
func changeNewFloor(a *CustomAgent3) {
	currentFloor := a.Floor()
	if len(a.knowledge.floors) != 0 {
		oldFloor := a.knowledge.floors[len(a.knowledge.floors)-1]

		if currentFloor < oldFloor {
			changeInMood(5, 15, -1, a)
		} else {
			changeInMood(10, 20, 1, a)
		}

	}
	a.knowledge.floors = append(a.knowledge.floors, currentFloor)
}
