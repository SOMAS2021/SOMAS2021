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

//Function adds the new people we have not met to our friends
func addFriends(friend1, friend2 int, a *CustomAgent3) {
	var new1 = true
	var new2 = true
	for i := 0; i < len(a.knowledge.friends); i++ {
		if a.knowledge.friends[i] == friend1 {
			new1 = false
		}
		if a.knowledge.friends[i] == friend1 {
			new2 = false
		}
	}

	if new1 {
		a.knowledge.friends = append(a.knowledge.friends, friend1)
		a.knowledge.friendship = append(a.knowledge.friendship, 0.5)
	}
	if new2 {
		a.knowledge.friends = append(a.knowledge.friends, friend2)
		a.knowledge.friendship = append(a.knowledge.friendship, 0.5)
	}
}

// Function will return the friendship level for a specific agent, if we don't know them friendship is 0, and the position it is stored at
func friendshipLevel(friend int, a *CustomAgent3) (float64, int) {
	for i := 0; i < len(a.knowledge.friends); i++ {
		if a.knowledge.friends[i] == friend {
			return a.knowledge.friendship[i], i
		}
	}
	return 0, -1
}

// Function changes value of friendship depending on factor change -1 to 1, negative reduces frienship, positive increases
func friendshipChange(friend int, change float64, a *CustomAgent3) {

	var level, index = friendshipLevel(friend, a)

	if index >= 0 {
		if change < 0 {
			level = level + change*(level-0)
		} else {
			level = level + change*(1-level)
		}
		a.knowledge.friendship[index] = level
	}
}

// Function gets as input the mini and max change we want in, direction marks if we want it to go up or down
func changeInMood(a *CustomAgent3, pointsMin, pointsMax, direction int) {
	// TODO: Remove this line. See Issue #60.
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	points := r1.Intn(pointsMax-pointsMin) + pointsMin
	if direction <= 0 {
		a.vars.mood = a.vars.mood - points
		if a.vars.mood < 0 {
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
func changeInMorality(a *CustomAgent3, pointsMin, pointsMax, direction int) {
	points := rand.Intn(pointsMax-pointsMin) + pointsMin
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
	if a.HP() < 50 {
		changeInMorality(a, 10, 15, -1)
		changeInMood(a, 10, 15, -1)
	} else {
		changeInMorality(a, 10, 15, 1)
		changeInMood(a, 10, 15, 1)
	}
}

// Function is called when the floor changes, changes the mood when we change floors
func changeNewFloor(a *CustomAgent3) {
	currentFloor := a.Floor()
	if len(a.knowledge.floors) != 0 {
		if currentFloor < a.knowledge.floors[len(a.knowledge.floors)-1] {
			changeInMood(a, 5, 15, -1)
		} else {
			changeInMood(a, 10, 20, 1)
		}

	}
	a.knowledge.floors = append(a.knowledge.floors, currentFloor)
}
