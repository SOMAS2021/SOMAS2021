package team3

import (
	"math/rand"
	"time"
)

// Function decides if we read a message or we don't, depends on our level of stubbornness, return is bool
// Stubborness of 20 means probability of 0.2 that we don't read message
func read() bool {
	rand.Seed(time.Now().UnixNano()) //creates seed to rand
	var random = rand.Intn(100)

	if random <= stubbornness {
		return true
	}
	return false
}

// Function gets as input the mini and max change we want in
func changeInMood(pointsMin, pointsMax, direction int) {
	rand.Seed(time.Now().UnixNano()) //creates seed to rand
	if direction < 0 {
		a.vars.mood = (rand.Intn(pointsMax-pointsMin) + pointsMin) * -1
	}
	if direction > 0 {
		mood = (rand.Intn(pointsMax-pointsMin) + pointsMin) * 1
	}
}
