package team3

import (
	"math/rand"

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
func (a *CustomAgent3) updateFriendship(friend uuid.UUID, change int) {
	level, found := a.knowledge.friends[friend]
	b := 0.3
	c := 0.4
	if !found {
		a.knowledge.friends[friend] = 0.4 + (float64(a.vars.morality)/100)*0.2
		return
	}
	if change < 0 {
		level -= b * (level)
	} else {
		level += c * (1 - level)
	}
	a.knowledge.friends[friend] = level

}

// Function gets as input the mini and max change we want in, direction marks if we want it to go up or down
func (a *CustomAgent3) changeInMood(pointsMin, pointsMax, direction int) {
	points := rand.Intn(pointsMax-pointsMin) + pointsMin
	if direction < 0 {
		a.vars.mood -= points
		if a.vars.mood < 0 {
			a.vars.mood = 0
		}
	} else {
		a.vars.mood += points
		if a.vars.mood > 100 {
			a.vars.mood = 100
		}
	}
}

// Function gets as input the mini and max change we want in, direction marks if we want it to go up or down
func (a *CustomAgent3) changeInMorality(pointsMin, pointsMax, direction int) {
	points := rand.Intn(pointsMax-pointsMin) + pointsMin
	if direction < 0 {
		a.vars.morality -= points
		if a.vars.morality < 0 {
			a.vars.morality = 0
		}
	} else {
		a.vars.morality += points
		if a.vars.morality > 100 {
			a.vars.morality = 100
		}
	}
}

func (a *CustomAgent3) changeInStubbornness(change, direction int) {
	if direction < 0 {
		a.vars.stubbornness -= change
		if a.vars.stubbornness < 0 {
			a.vars.stubbornness = 0
		}
	} else {
		a.vars.stubbornness += change
		if a.vars.stubbornness > 75 {
			a.vars.stubbornness = 75
		}
	}
}

// Updates mood and morality when the floor is changed. Called at the start of each day.
func (a *CustomAgent3) changeNewDay() {
	a.knowledge.lastHP = a.HP()
	a.knowledge.rememberedAge = a.BaseAgent().Age()
	if a.HP() < 25 {
		a.changeInMorality(5, 10, -1)
		a.changeInMood(8, 16, -1)
	} else if a.HP() < 40 {
		a.changeInMorality(0, 4, -1)
		a.changeInMood(1, 6, -1)
	} else if a.HP() < 75 {
		a.changeInMorality(0, 4, 1)
		a.changeInMood(1, 6, 1)
	} else {
		a.changeInMorality(3, 7, 1)
		a.changeInMood(5, 10, 1)
	}
	a.knowledge.foodMovingAvg = float64(a.knowledge.foodMovingAvg)*0.9 + float64(a.knowledge.foodLastEaten)*0.1
}

func (a *CustomAgent3) reshuffleEstimator() {
	if a.knowledge.reshuffleEst == -1 {
		a.knowledge.reshuffleEst = a.Age()
	} else {
		if a.Age()-a.knowledge.reshuffleEst < a.knowledge.reshuffleEst { //attempt to correct in case we resshufle to same floor twice
			a.knowledge.reshuffleEst = a.Age() - a.knowledge.reshuffleEst
		}
	}
}

// Function is called when the floor changes, changes the mood when we change floors
func (a *CustomAgent3) changeNewFloor() {

	//Calculate mood based on what floors we have been in
	if len(a.knowledge.floors) != 0 {
		lastFloor := a.knowledge.floors[len(a.knowledge.floors)-1]
		currentFloor := a.Floor()
		beenInHigher := false
		beenInLower := false
		a.reshuffleEstimator()
		a.knowledge.hpAbove = -1
		a.knowledge.hpBelow = -1

		for _, floor := range a.knowledge.floors {
			if floor > currentFloor {
				beenInLower = true
			} else if floor < currentFloor {
				beenInHigher = true
			}
			if beenInLower && beenInHigher {
				break
			}
		}
		if !beenInHigher {
			a.changeInMood(5, 10, 1)
		} else if !beenInLower {
			a.changeInMood(5, 15, -1) //Situation is specially bad, so we get more depressed
		} else if currentFloor < lastFloor {
			a.changeInMood(0, 5, 1)
		} else {
			a.changeInMood(0, 5, -1)
		}
	}
	a.knowledge.floors = append(a.knowledge.floors, a.Floor())
}
