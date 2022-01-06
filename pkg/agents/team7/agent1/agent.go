package team7agent1

import (
	"math"
	"math/rand"

	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/health"
	"github.com/google/uuid"
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

type floorInfo struct {
	numOfDays, avgFood int
}

type team7Personalities struct {
	openness          int
	conscientiousness int
	extraversion      int
	agreeableness     int
	neuroticism       int
}

type CurrentBehaviour struct { //determines directly the food taken
	greediness     int
	kindness       int // likehood to read and send messages // affacted by openess and extraversion
	responsiveness int // affected by days on floor, days hungry // changes less frequently than sanity
	sanity         int

	// morality   int
	// affected by days on floor, previous floor, days hungry and frequency of communictaion(extraversion),
	// neuroticism
	// more sensitive to other factors
	// lower sanity -> less logical behaviour -> (ex: no food taken when starving ?)
}

type OperationalMemory struct {
	orderPrevFloors   []int
	prevFloors        map[int]floorInfo
	currentDayonFloor int
	numOfDays         int
	avgFood           int
	currentFloorRisk  int
	daysHungry        int
	seenPlatform      bool
	prevHP            int
	leaveFood         int
	foodEaten         food.FoodType
}

type CustomAgent7 struct {
	*infra.Base
	personality    team7Personalities
	opMem          OperationalMemory
	behaviour      CurrentBehaviour
	Eaten          food.FoodType
	activeTreaties map[uuid.UUID]messages.Treaty
}

func New(baseAgent *infra.Base) (infra.Agent, error) {
	openness := rand.Intn(100)
	conscientiousness := rand.Intn(100)
	extraversion := rand.Intn(100)
	agreeableness := rand.Intn(100)
	neuroticism := rand.Intn(100)

	prevFloors := map[int]floorInfo{}

	return &CustomAgent7{
		Base: baseAgent,
		personality: team7Personalities{
			openness:          openness,
			conscientiousness: conscientiousness,
			extraversion:      extraversion,
			agreeableness:     agreeableness,
			neuroticism:       neuroticism,
		},
		behaviour: CurrentBehaviour{
			greediness:     100 - agreeableness,
			kindness:       agreeableness,
			sanity:         100 - neuroticism,
			responsiveness: (agreeableness + openness + extraversion) / 3,
		},
		opMem: OperationalMemory{
			orderPrevFloors:   []int{},
			prevFloors:        prevFloors,
			currentDayonFloor: 0,
			numOfDays:         0,
			avgFood:           0,
			currentFloorRisk:  0,
			daysHungry:        0,
			seenPlatform:      false,
			foodEaten:         0,
			prevHP:            100,
			leaveFood:         0,
		},
		Eaten: 0,
	}, nil
}

func (a *CustomAgent7) Run() {

	a.Log("Team7Agent1 reporting status:", infra.Fields{"floor": a.Floor(), "hp": a.HP(), "greed": a.behaviour.greediness, "kind": a.behaviour.kindness, "aggr": a.personality.agreeableness})

	// Operational Variables
	// UserID := a.ID()
	currentHP := a.HP()
	healthInfo := a.HealthInfo()
	currentFloor := a.Floor()
	// prevFloor := a.opMem.orderPrevFloors[len(a.opMem.orderPrevFloors)-1]

	// ------------------ Run Block A.1: Estimates the food available on the new floor based on past experience ------------------

	//Check if floor has changed
	if len(a.opMem.orderPrevFloors) == 0 || a.opMem.orderPrevFloors[len(a.opMem.orderPrevFloors)-1] != currentFloor { //If day 1 or floor change
		a.opMem.currentDayonFloor = 1          //reset currentDay counter
		if len(a.opMem.orderPrevFloors) != 0 { //if the floor tracker is not empty

			closestFloorAboveCurrent := 0
			closestFloorBelowCurrent := math.Inf(1)
			beenOnCurrentFloorBefore := false

			for i := 0; i < len(a.opMem.orderPrevFloors) && !beenOnCurrentFloorBefore; i++ {
				if a.opMem.orderPrevFloors[i] == currentFloor {
					beenOnCurrentFloorBefore = true
				}
			}

			if !beenOnCurrentFloorBefore {
				for i := 0; i < len(a.opMem.orderPrevFloors); i++ {
					if a.opMem.orderPrevFloors[i] < currentFloor && a.opMem.orderPrevFloors[i] > closestFloorAboveCurrent {
						closestFloorAboveCurrent = a.opMem.orderPrevFloors[i]
					}
				}
				for i := 0; i < len(a.opMem.orderPrevFloors); i++ {
					if a.opMem.orderPrevFloors[i] > currentFloor && a.opMem.orderPrevFloors[i] < int(closestFloorBelowCurrent) {
						closestFloorBelowCurrent = float64(a.opMem.orderPrevFloors[i])
					}
				}
			}

			var expectedFood int

			if beenOnCurrentFloorBefore {
				expectedFood = a.opMem.prevFloors[currentFloor].avgFood
			} else {

				if closestFloorAboveCurrent != 0 && closestFloorBelowCurrent != math.Inf(1) {
					expectedFood = a.opMem.prevFloors[int(closestFloorBelowCurrent)].avgFood + (a.opMem.prevFloors[closestFloorAboveCurrent].avgFood-a.opMem.prevFloors[int(closestFloorBelowCurrent)].avgFood)*(int(closestFloorBelowCurrent)-currentFloor)/(int(closestFloorBelowCurrent)-int(closestFloorAboveCurrent))
				}
				if closestFloorAboveCurrent == 0 && closestFloorBelowCurrent != math.Inf(1) {
					if a.opMem.prevFloors[int(closestFloorBelowCurrent)].avgFood == 0 {
						expectedFood = 0
					} else {
						estimatedMealSize := healthInfo.HPReqCToW + (a.personality.openness / 5)
						expectedFood = a.opMem.prevFloors[int(closestFloorBelowCurrent)].avgFood + (int(closestFloorBelowCurrent)-currentFloor)*estimatedMealSize
					}
				}
				if closestFloorAboveCurrent != 0 && closestFloorBelowCurrent == math.Inf(1) {
					estimatedMealSize := healthInfo.HPReqCToW + ((100 - a.personality.openness) / 5)
					expectedFood = a.opMem.prevFloors[closestFloorAboveCurrent].avgFood - (currentFloor-closestFloorAboveCurrent)*estimatedMealSize
				}
				if expectedFood < 0 {
					expectedFood = 0
				}

			}

			// ------------------ Run Block A.2: Adjusts mood w.r.t. change in floor and the expected food ------------------

			//Update Behaviour
			if currentFloor < a.opMem.orderPrevFloors[len(a.opMem.orderPrevFloors)-1] {
				//only negatively impacted if openness is low
				if a.personality.openness <= 50 {
					a.behaviour.greediness += (currentFloor - a.opMem.orderPrevFloors[len(a.opMem.orderPrevFloors)-1])

				}
			} else {
				//if we have moved up our kindness increases
				a.behaviour.kindness += (a.opMem.orderPrevFloors[len(a.opMem.orderPrevFloors)-1] - currentFloor)
			}
			//greediness is effected by the amount of food expected on this floor as estimated from past experience
			a.behaviour.greediness -= a.opMem.currentFloorRisk
			if expectedFood < healthInfo.HPReqCToW { //greediness only begins to be effected once the expected food is below a crticial threshold
				a.behaviour.greediness += (healthInfo.HPReqCToW - expectedFood)
				a.opMem.currentFloorRisk = (healthInfo.HPReqCToW - expectedFood)
			} else {
				a.opMem.currentFloorRisk = 0
			}

		}
		a.opMem.orderPrevFloors = append(a.opMem.orderPrevFloors, currentFloor) //append current floor to floor tracker
	} else { //increment currentDayonFloor counter
		a.opMem.currentDayonFloor++
	}

	// ------------------ Run Block X: Receive messages ------------------

	receivedMsg := a.ReceiveMessage()
	if receivedMsg != nil {
		receivedMsg.Visit(a)
	} else {
		a.Log("No Messages")
	}

	// // Possible food to take: Food available on platform/Food required to go from ciritcal to weak, critical to healthy

	// // Influences on Food taking: greediness, kindness, memory map HOW???, Sanity [depends on messaging, extraversion???, neuroticism], activeness, morality,
	// // so if your extraverted: less interaction leads to lowered sanity/More interaction better sanity
	// // higher neuroticism -> fall faster into insanity, initialise sanity based on neuroticism, sanity = 100-neuroticism/ sanity bsaed on

	if a.PlatformOnFloor() && !a.opMem.seenPlatform {

		// ------------------ Run Block X: Calculates average food available on this floor ------------------

		// If the number of days on floor is zero set num of days to 1 and add current avail food to average
		if (a.opMem.prevFloors[a.Floor()].numOfDays == 0) && a.PlatformOnFloor() {
			a.opMem.prevFloors[a.Floor()] = floorInfo{1, int(a.CurrPlatFood())}
		} else if a.PlatformOnFloor() {
			//otherwise update number of days on floor and calculate new average
			tmp := a.opMem.prevFloors[a.Floor()].avgFood * a.opMem.prevFloors[a.Floor()].numOfDays
			tmp = tmp + int(a.CurrPlatFood())
			newNumOfDays := a.opMem.prevFloors[a.Floor()].numOfDays + 1
			tmp = tmp / newNumOfDays
			a.opMem.prevFloors[a.Floor()] = floorInfo{tmp, newNumOfDays}
		}

		// ------------------ Run Block D.1: Adjusting mood w.r.t. personality and randomness ------------------

		r1 := rand.Intn(11) - 5
		r2 := rand.Intn(11) - 5
		r3 := rand.Intn(11) - 5

		a.behaviour.kindness += (a.personality.neuroticism * r1) / 50
		a.behaviour.greediness += (a.personality.neuroticism * r2) / 50
		a.behaviour.sanity += (a.personality.neuroticism * r3) / 50 //also factor in total_messages*a.personality.extraversion !!!

		if a.behaviour.kindness < 0 {
			a.behaviour.kindness = 0
		} else if a.behaviour.kindness > 100 {
			a.behaviour.kindness = 100
		}

		if a.behaviour.greediness < 0 {
			a.behaviour.greediness = 0
		} else if a.behaviour.greediness > 100 {
			a.behaviour.greediness = 100
		}

		// ------------------ Run Block D.2: Take food w.r.t. current health, mood, messages and treaties ------------------

		var foodtotake food.FoodType

		// Cases
		switch {
		case currentHP <= healthInfo.HPCritical:
			foodtotake = health.FoodRequired(currentHP, 60, healthInfo)
		case healthInfo.WeakLevel <= currentHP && currentHP <= 60:
			foodtotake = food.FoodType(100 - a.behaviour.greediness - a.behaviour.kindness)
		default:
			foodtotake = health.FoodRequired(currentHP, healthInfo.WeakLevel, healthInfo)
		}

		// Eat Food and Update Statuses
		Eaten, err := a.TakeFood(foodtotake)

		if err != nil {
			switch err.(type) {
			case *infra.FloorError:
			case *infra.NegFoodError:
			case *infra.AlreadyEatenError:
			default:
			}
		}
		a.opMem.foodEaten = Eaten
		if a.opMem.foodEaten > 0 {
			a.opMem.daysHungry = 0
		} else {
			a.opMem.daysHungry++
		}
		a.Log("Team 7 has seen the platform:", infra.Fields{"foodEaten": a.opMem.foodEaten, "health": a.HP(), "daysHungry": a.opMem.daysHungry})
		a.opMem.seenPlatform = true
	}

	if a.PlatformOnFloor() && a.opMem.seenPlatform {
		//Get ready for new day
		a.opMem.seenPlatform = false
		a.opMem.prevHP = a.HP()
	}

}

// Message Handlers

// Handle Asks and auto respond
func (a *CustomAgent7) HandleAskHP(msg messages.AskHPMessage) {
	a.Log("Recieved askHP message from ", infra.Fields{"floor": msg.SenderFloor()})
	reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor()-a.Floor(), a.HP())
	a.SendMessage(reply)
}

func (a *CustomAgent7) HandleAskFoodTaken(msg messages.AskFoodTakenMessage) {
	a.Log("Recieved askFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor()})
	reply := msg.Reply(a.ID(), msg.SenderFloor()-a.Floor(), a.Floor(), int(a.opMem.foodEaten))
	a.SendMessage(reply)
}

func (a *CustomAgent7) HandleAskIntendedFoodTaken(msg messages.AskIntendedFoodIntakeMessage) {
	a.Log("Recieved askIntendedFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor()})

	reply := msg.Reply(a.ID(), msg.SenderFloor()-a.Floor(), a.Floor(), 1)
	a.SendMessage(reply)
}

// Statements
func (a *CustomAgent7) HandleStateFoodTaken(msg messages.StateFoodTakenMessage) {
	statement := msg.Statement()
	a.Log("Recieved a StateFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor(), "statement": statement})
}

func (a *CustomAgent7) HandleStateHP(msg messages.StateHPMessage) {
	statement := msg.Statement()
	a.Log("I recieved a StateHP message from ", infra.Fields{"floor": msg.SenderFloor(), "statement": statement})
}

func (a *CustomAgent7) HandleStateIntendedFoodTaken(msg messages.StateIntendedFoodIntakeMessage) {
	statement := msg.Statement()
	a.Log("I recieved a StateIntendedFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor(), "statement": statement})
}

// Requests
func (a *CustomAgent7) HandleRequestLeaveFood(msg messages.RequestLeaveFoodMessage) {
	reply := msg.Reply(a.ID(), msg.SenderFloor()-a.Floor(), a.Floor(), true)
	a.SendMessage(reply)
	a.opMem.leaveFood = int(msg.Request())
	a.Log("Recieved requestLeaveFood message from ", infra.Fields{"floor": msg.SenderFloor()})
}

func (a *CustomAgent7) HandleRequestTakeFood(msg messages.RequestTakeFoodMessage) {
	a.Log("Recieved requestTakeFood message from ", infra.Fields{"floor": msg.SenderFloor()})
	reply := msg.Reply(a.ID(), msg.SenderFloor()-a.Floor(), a.Floor(), true)
	a.SendMessage(reply)
}

// Responses
func (a *CustomAgent7) HandleResponse(msg messages.BoolResponseMessage) {
	response := msg.Response()
	a.Log("I recieved a Response message from ", infra.Fields{"floor": msg.SenderFloor(), "response": response})
}
