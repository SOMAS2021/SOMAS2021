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
