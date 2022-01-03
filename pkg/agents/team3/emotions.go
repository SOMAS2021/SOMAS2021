package team3

import (
	"math/rand"
	"time"
)

// Function decides if we read a message or we don't, depends on our level of stubbornness, return is bool
// Stubborness of 20 means probability of 0.2 that we don't read message

func read(a *CustomAgent3) bool {
	random := rand.Intn(100)

	return random <= a.vars.stubbornness
}

////Function adds a new person to our freindship list if they are not there yet.
//func addFriend(a *CustomAgent3, friend string) {
//	newFriend := true
//	for i := 0; i < len(a.knowledge.friends); i++ {
//		if a.knowledge.friends[i] == friend {
//			newFriend = false
//			break
//		}
//
//	}
//
//	if newFriend {
//		a.knowledge.friends = append(a.knowledge.friends, friend)
//		a.knowledge.friendship = append(a.knowledge.friendship, 0.4+(float64(a.vars.morality)/100)*0.2)
//	}
//
//}

// Function will return the friendship level for a specific agent, if we don't know them friendship is 0, and the position it is stored at
//func friendshipLevel(a *CustomAgent3, friend string) (float64, int) {
//	for i := 0; i < len(a.knowledge.friends); i++ {
//		if a.knowledge.friends[i] == friend {
//			return a.knowledge.friendship[i], i
//		}
//	}
//	return 0, -1
//}

// Function changes value of friendship depending on factor change -1 to 1, negative reduces frienship, positive increases
//func friendshipChange(a *CustomAgent3, friend string, change float64) {
//
//	var level, index = friendshipLevel(a, friend)
//
//	if index >= 0 {
//		if change < 0 {
//			level = level - change*(level)
//		} else {
//			level = level + change*(1-level)
//		}
//		a.knowledge.friendship[index] = level
//	}
//
//}

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
//TO DO: Trial this and do sth different.
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
	//Restart neighbour functions as we don´t know who is above/below anymore
	a.knowledge.floorAbove = ""
	a.knowledge.floorBelow = ""

	//Calculate mood based on what floors we have been in
	currentFloor := a.Floor()

	//If we are still in the same floor as before we keep the same mood
	if len(a.knowledge.floors) != 0 {
		lastFloor := a.knowledge.floors[len(a.knowledge.floors)-1]
		beenInHigher := false
		beenInLower := false

		for _, floor := range a.knowledge.floors {
			if a.knowledge.floors[i] < currentFloor {
				beenInLower = true
			}
			if a.knowledge.floors[i] > currentFloor {
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
