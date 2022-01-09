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
	messagesSent      bool
	recievedReq       bool
	takeFood          int
	leaveFood         int
}

type CustomAgent7 struct {
	*infra.Base
	personality team7Personalities
	opMem       OperationalMemory
	behaviour   CurrentBehaviour
}

func New(baseAgent *infra.Base) (infra.Agent, error) {
	// openness := rand.Intn(100)
	// conscientiousness := rand.Intn(100)
	// extraversion := rand.Intn(100)
	// agreeableness := rand.Intn(100)
	// neuroticism := rand.Intn(100)
	openness := rand.Intn(60) + 20
	conscientiousness := rand.Intn(60) + 20
	extraversion := rand.Intn(60) + 20
	agreeableness := rand.Intn(60) + 20
	neuroticism := rand.Intn(60) + 20

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
			messagesSent:      false,
			recievedReq:       false,
			takeFood:          0,
			leaveFood:         0,
		},
	}, nil
}

func (a *CustomAgent7) Run() {

	// a.Log("Team7Agent1 reporting status:", infra.Fields{"floor": a.Floor(), "hp": a.HP(), "greed": a.behaviour.greediness, "kind": a.behaviour.kindness, "aggr": a.personality.agreeableness})

	// Operational Variables
	UserID := a.ID()
	currentHP := a.HP()
	daysCritical := a.DaysAtCritical()
	healthInfo := a.HealthInfo()
	currentFloor := a.Floor()

	// Only occurs on a floor change
	if len(a.opMem.orderPrevFloors) == 0 || a.opMem.orderPrevFloors[len(a.opMem.orderPrevFloors)-1] != currentFloor { //If day 1 or floor change

		a.Log("checking floor change........................................................................................")
		if a.Age() == 0 {
			a.Log("Personality Makeup:", infra.Fields{"A": a.personality.agreeableness, "C": a.personality.conscientiousness, "E": a.personality.extraversion, "N": a.personality.neuroticism, "O": a.personality.openness})
		}

		//Reset Variables on Floor Change
		a.opMem.currentDayonFloor = 1 // reset currentDay counter
		a.opMem.recievedReq = false   // reset Requests
		a.opMem.takeFood = 0
		a.opMem.leaveFood = 0

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

			// Update Behaviour: Higher in the tower (lower floor number) increases greed/decreases kindess and vice versa
			if currentFloor < a.opMem.orderPrevFloors[len(a.opMem.orderPrevFloors)-1] {
				// A higher openess will more significantly reduce greedines and increase kindness
				a.behaviour.greediness -= int(float64(a.opMem.orderPrevFloors[len(a.opMem.orderPrevFloors)-1]/currentFloor) * (float64(a.personality.openness) / 50))
				a.behaviour.kindness += int(float64(a.opMem.orderPrevFloors[len(a.opMem.orderPrevFloors)-1]/currentFloor) * (float64(a.personality.openness) / 50))
			} else {
				// A lower openess will more significantly increase greedines and reduce kindness
				a.behaviour.greediness += int(float64(currentFloor/a.opMem.orderPrevFloors[len(a.opMem.orderPrevFloors)-1]) * (float64(100-a.personality.openness) / 50))
				a.behaviour.kindness -= int(float64(currentFloor/a.opMem.orderPrevFloors[len(a.opMem.orderPrevFloors)-1]) * (float64(100-a.personality.openness) / 50))
			}

			// Greediness is affected by the amount of food expected on this floor as estimated from past experience
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

	// ------------------ Run() Block B: Receive messages ------------------

	receivedMsg := a.ReceiveMessage()
	if receivedMsg != nil {
		receivedMsg.Visit(a)
	} else {
		a.Log("No Messages")
	}

	// Messaging

	if currentHP < a.opMem.prevHP {
		a.opMem.messagesSent = false
	}

	if !a.opMem.messagesSent {
		msg := messages.NewAskHPMessage(UserID, currentFloor, currentFloor+1)
		a.SendMessage(msg)
		if currentFloor != 1 {
			msg = messages.NewAskHPMessage(UserID, currentFloor, currentFloor-1)
			a.SendMessage(msg)
		}

		if a.PlatformOnFloor() {
			msg2 := messages.NewAskFoodTakenMessage(UserID, currentFloor, currentFloor+1)
			a.SendMessage(msg2)
		}
		if !a.PlatformOnFloor() && a.CurrPlatFood() != -1 && a.opMem.seenPlatform && currentFloor != 1 {
			msg2 := messages.NewAskFoodTakenMessage(UserID, currentFloor, currentFloor-1)
			a.SendMessage(msg2)
		}

		a.opMem.messagesSent = true
	}

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

		// treaty := false

		// ------------------ Run() Block D.2: Take food w.r.t. current health, mood, messages and treaties ------------------

		var foodtotake food.FoodType

		satisficedHP := 50
		targetWeakFood := health.FoodRequired(currentHP, healthInfo.WeakLevel, healthInfo)
		targetSatisficedFood := health.FoodRequired(currentHP, satisficedHP, healthInfo)
		targetFullFood := health.FoodRequired(currentHP, healthInfo.MaxHP, healthInfo)
		scalingRatio := float64(targetFullFood-targetSatisficedFood) / float64(targetSatisficedFood-targetWeakFood)
		//scalingRatio := float64(healthInfo.MaxHP-satisficedHP) / float64(satisficedHP-healthInfo.WeakLevel)

		a.Log("ABCDEFG:", infra.Fields{"hp": a.HP(), "greed": a.behaviour.greediness, "kind": a.behaviour.kindness, "W": targetWeakFood, "H": targetSatisficedFood, "F": targetFullFood})

		switch {
		// Highest prioirty case - agent dies if he stays critical for 3 more days
		case currentHP <= healthInfo.HPCritical && daysCritical >= (healthInfo.MaxDayCritical-3):

			kindnessAdjuster := (float64(a.behaviour.kindness) / 1000) * float64(targetSatisficedFood-targetWeakFood)
			greedinessAdjuster := (float64(a.behaviour.greediness) / 100) * float64(targetFullFood-targetSatisficedFood) / scalingRatio
			foodtotake = targetSatisficedFood - food.FoodType(kindnessAdjuster) + food.FoodType(greedinessAdjuster)

		// Fulfilling message requests are given high priority
		case a.opMem.recievedReq:
			foodtotake = food.FoodType(a.opMem.takeFood)
			// Fulfilling treaties is given high priority

		// case treaty:
		case len(a.ActiveTreaties()) != 0:
			for _, t_active := range a.ActiveTreaties() {

				if a.PlatformOnFloor() && t_active.Request() != messages.Inform {
					available := a.CurrPlatFood()
					amount := food.FoodType(float64(t_active.RequestValue()/100)) * available
					if t_active.Request() == messages.LeaveAmountFood {
						amount = food.FoodType(t_active.RequestValue())
					}
					// case t_active.condition
					// check HP condition
					switch t_active.ConditionOp() { // if HP > 10 , Leaveamount > 15, availble = 37, foottotake = 30
					case messages.GT:
						if a.HP() > t_active.ConditionValue() || int(a.CurrPlatFood()) > t_active.ConditionValue() {
							switch t_active.RequestOp() {
							case messages.GT:
								if foodtotake > available-food.FoodType(t_active.RequestValue()) {
									foodtotake = available - food.FoodType(amount) - 1
								}
							case messages.GE:
								if foodtotake >= available-food.FoodType(t_active.RequestValue()) {
									foodtotake = available - food.FoodType(amount)
								}
							case messages.EQ:
								if foodtotake != available-food.FoodType(t_active.RequestValue()) {
									foodtotake = available - food.FoodType(amount)
								}
							}
						}
					case messages.GE:
						if a.HP() >= t_active.ConditionValue() || int(a.CurrPlatFood()) >= t_active.ConditionValue() {
							switch t_active.RequestOp() {
							case messages.GT:
								if foodtotake > available-food.FoodType(t_active.RequestValue()) {
									foodtotake = available - food.FoodType(amount) - 1
								}
							case messages.GE:
								if foodtotake >= available-food.FoodType(t_active.RequestValue()) {
									foodtotake = available - food.FoodType(amount)
								}
							case messages.EQ:
								if foodtotake != available-food.FoodType(t_active.RequestValue()) {
									foodtotake = available - food.FoodType(amount)
								}
							}
						}
					case messages.EQ:
						if a.HP() == t_active.ConditionValue() || int(a.CurrPlatFood()) == t_active.ConditionValue() || a.Floor() == t_active.ConditionValue() {
							switch t_active.RequestOp() {
							case messages.GT:
								if foodtotake > available-food.FoodType(t_active.RequestValue()) {
									foodtotake = available - food.FoodType(amount) - 1
								}
							case messages.GE:
								if foodtotake >= available-food.FoodType(t_active.RequestValue()) {
									foodtotake = available - food.FoodType(amount)
								}
							case messages.EQ:
								if foodtotake != available-food.FoodType(t_active.RequestValue()) {
									foodtotake = available - food.FoodType(amount)
								}
							}
						}
					case messages.LT:
						if a.Floor() < t_active.ConditionValue() {
							switch t_active.RequestOp() {
							case messages.GT:
								if foodtotake > available-food.FoodType(t_active.RequestValue()) {
									foodtotake = available - food.FoodType(amount) - 1
								}
							case messages.GE:
								if foodtotake >= available-food.FoodType(t_active.RequestValue()) {
									foodtotake = available - food.FoodType(amount)
								}
							case messages.EQ:
								if foodtotake != available-food.FoodType(t_active.RequestValue()) {
									foodtotake = available - food.FoodType(amount)
								}
							}
						}
					case messages.LE:
						if a.Floor() <= t_active.ConditionValue() {
							switch t_active.RequestOp() {
							case messages.GT:
								if foodtotake > available-food.FoodType(t_active.RequestValue()) {
									foodtotake = available - food.FoodType(amount) - 1
								}
							case messages.GE:
								if foodtotake >= available-food.FoodType(t_active.RequestValue()) {
									foodtotake = available - food.FoodType(amount)
								}
							case messages.EQ:
								if foodtotake != available-food.FoodType(t_active.RequestValue()) {
									foodtotake = available - food.FoodType(amount)
								}
							}
						}
					}
				}
			}
		// In this case the agent can stay critical for another 4 or more days. Hence the fulfillment of treaties and requests is prioritized over this
		case currentHP <= healthInfo.HPCritical:
			kindnessAdjuster := (float64(a.behaviour.kindness) / 1000) * float64(targetSatisficedFood-targetWeakFood)
			greedinessAdjuster := (float64(a.behaviour.greediness) / 100) * float64(targetFullFood-targetSatisficedFood) / scalingRatio
			foodtotake = targetSatisficedFood - food.FoodType(kindnessAdjuster) + food.FoodType(greedinessAdjuster)

		// Lowest priority case - it is the standard case, in the absence of messages and treaties, and the agent is not critical
		case healthInfo.HPCritical <= currentHP && currentHP <= healthInfo.MaxHP:
			kindnessAdjuster := (float64(a.behaviour.kindness) / 100) * float64(targetSatisficedFood-targetWeakFood)
			greedinessAdjuster := (float64(a.behaviour.greediness) / 100) * float64(targetFullFood-targetSatisficedFood) / scalingRatio
			foodtotake = targetSatisficedFood - food.FoodType(kindnessAdjuster) + food.FoodType(greedinessAdjuster)
		default:
			foodtotake = food.FoodType(0)
		}

		// Bound foodtotake
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

		a.Log("Team 7 has seen the platform:", infra.Fields{"foodEaten": a.opMem.foodEaten, "health": a.HP(), "daysHungry": a.opMem.daysHungry})
		a.opMem.seenPlatform = true
	}

	if !a.PlatformOnFloor() && a.opMem.seenPlatform {
		// reset variables Daily / prep for new day
		a.opMem.seenPlatform = false
		a.opMem.prevHP = a.HP()
		a.opMem.foodEaten = 0
		a.opMem.recievedReq = false
	}

	//End of Run()
}

// Message Handlers

// Handle Asks
func (a *CustomAgent7) HandleAskHP(msg messages.AskHPMessage) {
	a.Log("Recieved askHP message from ", infra.Fields{"floor": msg.SenderFloor()})
	reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), a.HP())
	a.SendMessage(reply)
}

func (a *CustomAgent7) HandleAskFoodTaken(msg messages.AskFoodTakenMessage) {
	a.Log("Recieved askFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor()})
	reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), int(a.opMem.foodEaten))
	a.SendMessage(reply)
}

func (a *CustomAgent7) HandleAskIntendedFoodTaken(msg messages.AskIntendedFoodIntakeMessage) {
	a.Log("Recieved askIntendedFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor()})
	reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), int(health.FoodRequired(a.HP(), 60, a.HealthInfo())))
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
func (a *CustomAgent7) HandleRequestTakeFood(msg messages.RequestTakeFoodMessage) {
	a.Log("Recieved requestTakeFood message from ", infra.Fields{"floor": msg.SenderFloor()})
	reqFood := msg.Request()
	if reqFood > int(health.FoodRequired(a.HP(), a.HealthInfo().MaxHP, a.HealthInfo())) || a.DaysAtCritical() >= (a.HealthInfo().MaxDayCritical-3) {
		// a.opMem.recievedReq = false
		reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), false)
		a.SendMessage(reply)
	} else {
		a.opMem.recievedReq = true
		a.opMem.takeFood = reqFood
		reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), true)
		a.SendMessage(reply)
	}
}

func (a *CustomAgent7) HandleRequestLeaveFood(msg messages.RequestLeaveFoodMessage) {
	a.Log("Recieved requestLeaveFood message from ", infra.Fields{"floor": msg.SenderFloor()})
	// a.opMem.recievedReq = false
	reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), false)
	a.SendMessage(reply)
}

// Responses
func (a *CustomAgent7) HandleResponse(msg messages.BoolResponseMessage) {
	response := msg.Response()
	a.Log("I recieved a Response message from ", infra.Fields{"floor": msg.SenderFloor(), "response": response})
}
