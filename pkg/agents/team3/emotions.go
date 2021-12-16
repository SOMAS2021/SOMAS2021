package team3

import (
	"math/rand"
	"time"
)

// Function decides if we read a message or we don't, depends on our level of stubbornness, return is bool
// Stubborness of 20 means probability of 0.2 that we don't read message
func read(a *CustomAgent3) bool {
	rand.Seed(time.Now().UnixNano()) //creates seed to rand
	var random = rand.Intn(100)

	if random <= a.vars.stubbornness {
		return true
	}
	return false
}

// Function gets as input the mini and max change we want in, direction marks if we want it to go up or down
func changeInMood(pointsMin, pointsMax, direction int, a *CustomAgent3) {
	rand.Seed(time.Now().UnixNano()) //creates seed to rand
	var temp = (rand.Intn(pointsMax-pointsMin) + pointsMin)
	if direction < 0 {
		temp = a.vars.mood - temp
		if temp < 0 {
			a.vars.mood = 0
		} else {
			a.vars.mood = temp
		}
	}
	if direction > 0 {
		temp = a.vars.mood + temp
		if temp > 100 {
			a.vars.mood = 100
		} else {
			a.vars.mood = temp
		}
	}
}

// Function gets as input the mini and max change we want in, direction marks if we want it to go up or down
func changeInMorality(pointsMin, pointsMax, direction int, a *CustomAgent3) {
	rand.Seed(time.Now().UnixNano()) //creates seed to rand
	var temp = (rand.Intn(pointsMax-pointsMin) + pointsMin)
	if direction < 0 {
		temp = a.vars.morale - temp
		if temp < 0 {
			a.vars.morale = 0
		} else {
			a.vars.morale = temp
		}
	}
	if direction > 0 {
		temp = a.vars.morale + temp
		if temp > 100 {
			a.vars.morale = 100
		} else {
			a.vars.morale = temp
		}
	}
}

func changeNewDay(a *CustomAgent3) {
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
	var currentFloor = a.Floor()
	var oldFloor = a.knowledge.floors[len(a.knowledge.floors)-1]

	if currentFloor < oldFloor {
		changeInMood(5, 15, -1, a)
	} else {
		changeInMood(10, 20, 1, a)
	}

	a.knowledge.floors = append(a.knowledge.floors, currentFloor)
}
