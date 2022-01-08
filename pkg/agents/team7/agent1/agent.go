package team7agent1

import (
	"math"
	"math/rand"

	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/health"
)

/*
Openness 	high-> doesn't get effected by new floors, looks forward to new day
			low-> does get negatively effected by new floors

Conscientiousness high-> plans ahead, attentive and takes into account messgages
			low-> no planning, fails to complete assigned tasks

Extraversion high-> very likely to communicate, will share a lot of information
			low-> does not like to communicate

Agreeableness high-> caring so more likely to go hungry for the betterment of others, trustworthy
			low-> greedy, manipulative

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
	responsiveness int // affected by days on floor, days hungry
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
	foodEaten         food.FoodType
	receivedReq       bool
	leaveFood         int
	takeFood          int
	treaty            int
	// activeTreaties map[uuid.UUID]messages.Treaty
}

type CustomAgent7 struct {
	*infra.Base
	personality team7Personalities
	opMem       OperationalMemory
	behaviour   CurrentBehaviour
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
			prevHP:            100,
			foodEaten:         0,
			receivedReq:       false,
			leaveFood:         0,
			takeFood:          0,
			treaty:            0,
		},
	}, nil
}

func (a *CustomAgent7) Run() {

	// a.Log("Team7Agent1 reporting status:", infra.Fields{"floor": a.Floor(), "hp": a.HP(), "greed": a.behaviour.greediness, "kind": a.behaviour.kindness, "aggr": a.personality.agreeableness})

	// Operational Variables
	currentHP := a.HP()
	healthInfo := a.HealthInfo()
	currentFloor := a.Floor()
	daysCritical := a.DaysAtCritical()

	// Only occurs on a floor change
	if len(a.opMem.orderPrevFloors) == 0 || a.opMem.orderPrevFloors[len(a.opMem.orderPrevFloors)-1] != currentFloor { //If day 1 or floor change
		a.Log("checking floor change........................................................................................")
		if a.Age() == 0 {
			a.Log("Personality Makeup:", infra.Fields{"A": a.personality.agreeableness, "C": a.personality.conscientiousness, "E": a.personality.extraversion, "N": a.personality.neuroticism, "O": a.personality.openness})
		}

		a.opMem.currentDayonFloor = 1          //reset currentDay counter
		if len(a.opMem.orderPrevFloors) != 0 { //if the floor tracker is not empty

			// ------------------ Run() Block A.1: Estimates the food available on the new floor based on past experience ------------------

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

			// ------------------ Run() Block A.2: Adjusts mood w.r.t. change in floor and the expected food ------------------

			// Update Behaviour
			if currentFloor < a.opMem.orderPrevFloors[len(a.opMem.orderPrevFloors)-1] {
				//only negatively impacted if openness is low
				if a.personality.openness <= 50 {
					a.behaviour.greediness += (a.opMem.orderPrevFloors[len(a.opMem.orderPrevFloors)-1] / currentFloor)

				}
			} else {
				//if we have moved up our kindness increases
				a.behaviour.kindness += (a.opMem.orderPrevFloors[len(a.opMem.orderPrevFloors)-1] / currentFloor)
			}
			//greediness is effected by the amount of food expected on this floor as estimated from past experience
			a.behaviour.greediness -= int(float64(a.opMem.currentFloorRisk) * (float64(a.personality.conscientiousness) / 50))
			foodReqCtoW := int(health.FoodRequired(healthInfo.HPCritical, healthInfo.WeakLevel, healthInfo))
			if expectedFood < foodReqCtoW { //greediness only begins to be effected once the expected food is below a crticial threshold
				a.opMem.currentFloorRisk = (foodReqCtoW - expectedFood)
				a.behaviour.greediness += int(float64(a.opMem.currentFloorRisk) * (float64(a.personality.conscientiousness) / 50))
			} else {
				a.opMem.currentFloorRisk = 0
			}

		}
		a.opMem.orderPrevFloors = append(a.opMem.orderPrevFloors, currentFloor) //append current floor to floor tracker
	} else { //increment currentDayonFloor counter
		a.opMem.currentDayonFloor++
	}

	// ------------------ Run() Block B: Receive messages (Every Tick) ------------------

	// receivedMsg := a.ReceiveMessage()
	// if receivedMsg != nil {
	// 	receivedMsg.Visit(a)
	// } else {
	// 	a.Log("No Messages")
	// }

	// Nayans Stuff

	// msg := messages.NewAskHPMessage(a.ID(), a.Floor(), 1)
	// a.SendMessage(msg)
	// msg = messages.NewAskHPMessage(a.ID(), a.Floor(), -1)
	// a.SendMessage(msg)

	// Runs Once a day (When platform is visible)
	if a.PlatformOnFloor() && !a.opMem.seenPlatform {

		// ------------------ Run() Block C: Calculates average food available on this floor ------------------

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

		// ------------------ Run() Block D.1: Adjusting mood w.r.t. personality and randomness ------------------

		r1 := rand.Intn(11) - 5
		r2 := rand.Intn(11) - 5

		a.behaviour.kindness += int(float64(a.personality.neuroticism*r1) / 50)
		a.behaviour.greediness += int(float64(a.personality.neuroticism*r2) / 50)

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

		// ------------------ Run() Block D.2: Take food w.r.t. current health, mood, messages and treaties ------------------

		var foodtotake food.FoodType

		// Defined Satisfied as 45 HP
		satisfiedHP := 45
		targetWeakFood := health.FoodRequired(currentHP, healthInfo.WeakLevel, healthInfo)
		targetHealthyFood := health.FoodRequired(currentHP, satisfiedHP, healthInfo)
		targetFullFood := health.FoodRequired(currentHP, healthInfo.MaxHP, healthInfo)
		a.Log("ABCDEFG:", infra.Fields{"hp": a.HP(), "greed": a.behaviour.greediness, "kind": a.behaviour.kindness, "W": targetWeakFood, "H": targetHealthyFood, "F": targetFullFood})
		switch {
		// Highest prioirty case - agent dies if he stays critical for 3 more days
		case currentHP <= healthInfo.HPCritical && daysCritical >= (healthInfo.MaxDayCritical-3):
			a.Log("CASE1--------------------------------------------------------------------------------------------")
			kindnessAdjuster := (float64(a.behaviour.kindness) / 1000) * float64(targetHealthyFood-targetWeakFood)
			greedinessAdjuster := (float64(a.behaviour.greediness) / 100) * float64(targetFullFood-targetHealthyFood)
			foodtotake = targetHealthyFood - food.FoodType(kindnessAdjuster) + food.FoodType(greedinessAdjuster)

		// // Fulfilling message requests are given high priority
		// case a.opMem.leaveFood >= 0:
		// 	foodtotake = a.CurrPlatFood() - food.FoodType(a.opMem.leaveFood)
		// case a.opMem.takeFood >= 0:
		// 	foodtotake = food.FoodType(a.opMem.takeFood)

		// // Fulfilling treaties is given high priority
		// case a.opMem.treaty != 0:
		// 	// foodtotake = a.opMem.treatyTake
		// 	foodtotake = a.CurrPlatFood() - food.FoodType(a.opMem.leaveFood)

		// In this case the agent can stay critical for another 4 or more days. Hence the fulfillment of treaties and requests is prioritized over this
		case currentHP <= healthInfo.HPCritical:
			a.Log("CASE2--------------------------------------------------------------------------------------------")
			kindnessAdjuster := (float64(a.behaviour.kindness) / 1000) * float64(targetHealthyFood-targetWeakFood)
			greedinessAdjuster := (float64(a.behaviour.greediness) / 100) * float64(targetFullFood-targetHealthyFood)
			foodtotake = targetHealthyFood - food.FoodType(kindnessAdjuster) + food.FoodType(greedinessAdjuster)

		// Lowest priority case - it is the standard case, in the absence of messages and treaties, and the agent is not critical
		default:
			a.Log("CASE3--------------------------------------------------------------------------------------------")
			kindnessAdjuster := (float64(a.behaviour.kindness) / 100) * float64(targetHealthyFood-targetWeakFood)
			greedinessAdjuster := (float64(a.behaviour.greediness) / 100) * float64(targetFullFood-targetHealthyFood)
			foodtotake = targetHealthyFood - food.FoodType(kindnessAdjuster) + food.FoodType(greedinessAdjuster)
		}

		if foodtotake > health.FoodRequired(currentHP, healthInfo.MaxHP, healthInfo) {
			foodtotake = health.FoodRequired(currentHP, healthInfo.MaxHP, healthInfo)
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
			a.behaviour.greediness -= int(0.5 * float64(a.opMem.daysHungry*a.opMem.daysHungry))
			a.opMem.daysHungry = 0
		} else {
			a.behaviour.greediness -= int(0.3 * float64(a.opMem.daysHungry*a.opMem.daysHungry))
			a.opMem.daysHungry++
			a.behaviour.greediness += int(0.3 * float64(a.opMem.daysHungry*a.opMem.daysHungry))
		}

		// a.Log("Team 7 has seen the platform:", infra.Fields{"foodEaten": a.opMem.foodEaten, "health": a.HP(), "daysHungry": a.opMem.daysHungry})
		a.opMem.seenPlatform = true
	}

	if !a.PlatformOnFloor() && a.opMem.seenPlatform {
		//Get ready for new day
		a.opMem.seenPlatform = false
		a.opMem.prevHP = a.HP()
		a.opMem.receivedReq = false
	}

	//End of Run()
}

// Message Handlers

// Handle Asks
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
	reply := msg.Reply(a.ID(), msg.SenderFloor()-a.Floor(), a.Floor(), int(health.FoodRequired(a.HP(), 60, a.HealthInfo())))
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
	//check health, if above critical check greediness/kindness, if critical do not leave food
	if a.HP() <= a.HealthInfo().HPCritical {
		reply := msg.Reply(a.ID(), msg.SenderFloor()-a.Floor(), a.Floor(), false)
		a.SendMessage(reply)
	} else if a.behaviour.greediness > a.behaviour.kindness {
		reply := msg.Reply(a.ID(), msg.SenderFloor()-a.Floor(), a.Floor(), false)
		a.SendMessage(reply)
	} else {
		reply := msg.Reply(a.ID(), msg.SenderFloor()-a.Floor(), a.Floor(), true)
		a.SendMessage(reply)
		a.opMem.receivedReq = true
		a.opMem.leaveFood = int(msg.Request())
		a.Log("Recieved requestLeaveFood message from ", infra.Fields{"floor": msg.SenderFloor()})
	}

}

func (a *CustomAgent7) HandleRequestTakeFood(msg messages.RequestTakeFoodMessage) {
	//check health, if above critical check greediness/kindness, if critical do not leave food
	if a.HP() <= a.HealthInfo().HPCritical {
		reply := msg.Reply(a.ID(), msg.SenderFloor()-a.Floor(), a.Floor(), false)
		a.SendMessage(reply)
	} else if a.behaviour.greediness > a.behaviour.kindness {
		reply := msg.Reply(a.ID(), msg.SenderFloor()-a.Floor(), a.Floor(), false)
		a.SendMessage(reply)
	} else {
		reply := msg.Reply(a.ID(), msg.SenderFloor()-a.Floor(), a.Floor(), true)
		a.SendMessage(reply)
		a.opMem.receivedReq = true
		a.opMem.leaveFood = int(msg.Request())
		a.Log("Recieved requestTakeFood message from ", infra.Fields{"floor": msg.SenderFloor()})
	}

}

// Responses
func (a *CustomAgent7) HandleResponse(msg messages.BoolResponseMessage) {
	response := msg.Response()
	a.Log("I recieved a Response message from ", infra.Fields{"floor": msg.SenderFloor(), "response": response})
}
