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
	prevAge           int
	prevHP            int
	foodEaten         food.FoodType
	messagesSent      bool
	msg1Sent          bool
	msg2Sent          bool
	msg3Sent          bool
	receivedReq       bool
	takeFood          int
	leaveFood         int
	// messageQueue      []int
}

type CustomAgent7 struct {
	*infra.Base
	personality team7Personalities
	opMem       OperationalMemory
	behaviour   CurrentBehaviour
}

func New(baseAgent *infra.Base) (infra.Agent, error) {
	openness := rand.Intn(60) + 20
	conscientiousness := rand.Intn(60) + 20
	extraversion := rand.Intn(60) + 20
	agreeableness := rand.Intn(60) + 20
	neuroticism := rand.Intn(60) + 20

	prevFloors := make(map[int]floorInfo)
	// messageQueue := make([]*messages.BaseMessage, 0)

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
			greediness:     ((100 - agreeableness) / 2) + 25,
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
			prevAge:           -1,
			prevHP:            100,
			foodEaten:         0,
			messagesSent:      false,
			msg1Sent:          false,
			msg2Sent:          false,
			msg3Sent:          false,
			receivedReq:       false,
			takeFood:          -1,
			leaveFood:         -1,
			// messageQueue:      messageQueue,
		},
	}, nil
}

func (a *CustomAgent7) Run() {

	a.manageFloor()

	a.manageMessages()

	a.manageFood()

	if !a.PlatformOnFloor() && a.opMem.seenPlatform {
		a.manageNewDay()
	}

	//End of Run()
}

func (a *CustomAgent7) manageFood() {

	if a.PlatformOnFloor() && !a.opMem.seenPlatform {

		// ------------------ Run() Block C: Calculates average food available on this floor ------------------
		a.foodAverage()

		// ------------------ Run() Block D.1: Adjusting mood w.r.t. personality and randomness ------------------
		a.adjustMood()

		// ------------------ Run() Block D.2: Take food w.r.t. current health, mood, messages and treaties ------------------
		var foodtotake food.FoodType

		satisficedHP := a.HealthInfo().MaxHP * 3 / 10
		healthyHP := a.HealthInfo().MaxHP * 65 / 100
		targetWeakFood := health.FoodRequired(a.HP(), a.HealthInfo().WeakLevel, a.HealthInfo())
		targetSatisficedFood := health.FoodRequired(a.HP(), satisficedHP, a.HealthInfo())
		targetHealthyFood := health.FoodRequired(a.HP(), healthyHP, a.HealthInfo())
		targetFullFood := health.FoodRequired(a.HP(), a.HealthInfo().MaxHP, a.HealthInfo())

		greedinessAdjuster := (float64(a.behaviour.greediness) / 300) * float64(targetFullFood-targetSatisficedFood)
		kindnessAdjuster := (float64(a.behaviour.kindness) / 300) * float64(targetFullFood-targetSatisficedFood)

		a.Log("ABCDEFG:", infra.Fields{"hp": a.HP(), "greed": a.behaviour.greediness, "kind": a.behaviour.kindness, "W": targetWeakFood, "H": targetHealthyFood, "S": targetSatisficedFood, "F": targetFullFood, "msg1Sent": a.opMem.msg1Sent, "msg2Sent": a.opMem.msg2Sent, "msg3Sent": a.opMem.msg3Sent})
		a.Log("TREATYLEN: ", infra.Fields{"TREATYLEN": len(a.ActiveTreaties())})

		switch {
		// Highest prioirty case - agent dies if he stays critical for 3 more days (or less)
		case a.HP() <= a.HealthInfo().HPCritical && a.DaysAtCritical() >= (a.HealthInfo().MaxDayCritical-1):

			a.Log("CASE1------------------------------------------------------------------------------------------------------------------------")

			foodtotake = targetSatisficedFood - food.FoodType(kindnessAdjuster/5) + food.FoodType(greedinessAdjuster)
			a.criticalStateTreaty()
			a.desperationTreaty()

		// Fulfilling message requests are given high priority
		case a.opMem.receivedReq:
			a.Log("CASE2------------------------------------------------------------------------------------------------------------------------")
			if a.opMem.takeFood != -1 {
				foodtotake = food.FoodType(a.opMem.takeFood)
			}

		// Fulfilling treaties is given high priority
		case len(a.ActiveTreaties()) != 0:
			a.Log("CASE3------------------------------------------------------------------------------------------------------------------------")
			foodtotake = targetSatisficedFood - food.FoodType(kindnessAdjuster/5) + food.FoodType(greedinessAdjuster)

			if a.HP() <= a.HealthInfo().HPCritical {
				a.propagateTreatyUpwards()
			}

			for _, tActive := range a.ActiveTreaties() {
				if a.PlatformOnFloor() && tActive.Request() != messages.Inform {
					available := a.CurrPlatFood()
					amount := food.FoodType(float64(tActive.RequestValue()/100)) * available
					if tActive.Request() == messages.LeaveAmountFood {
						amount = food.FoodType(tActive.RequestValue())
					}
					// case tActive.condition
					// check HP condition
					switch tActive.ConditionOp() { // if HP > 10 , Leaveamount > 15, availble = 37, foottotake = 30
					case messages.GT:
						foodtotake = a.treatyGT(tActive, foodtotake, available, amount)

					case messages.GE:
						foodtotake = a.treatyGE(tActive, foodtotake, available, amount)

					case messages.EQ:
						foodtotake = a.treatyEQ(tActive, foodtotake, available, amount)

					case messages.LT:
						foodtotake = a.treatyLT(tActive, foodtotake, available, amount)

					case messages.LE:
						foodtotake = a.treatyLE(tActive, foodtotake, available, amount)

					}
				}
			}
			// In this case the agent can stay critical for another 4 or more days. Hence the fulfillment of treaties and requests is prioritized over this
		case targetWeakFood > 0:
			a.Log("CASE4------------------------------------------------------------------------------------------------------------------------")
			foodtotake = targetSatisficedFood - food.FoodType(kindnessAdjuster/5) + food.FoodType(greedinessAdjuster)

		// Lowest priority case - it is the standard case, in the absence of messages and treaties, and the agent is not critical
		default:
			a.Log("Default CASE------------------------------------------------------------------------------------------------------------------------")
			foodtotake = targetSatisficedFood - food.FoodType(kindnessAdjuster) + food.FoodType(greedinessAdjuster)
		}

		a.eating(foodtotake, targetFullFood)

		a.Log("Team 7 has seen the platform:", infra.Fields{"foodEaten": a.opMem.foodEaten, "health": a.HP(), "daysHungry": a.opMem.daysHungry})
		a.opMem.seenPlatform = true
	}

}
func (a *CustomAgent7) eating(foodtotake food.FoodType, targetFullFood food.FoodType) {

	// Bound foodtotake
	if foodtotake > targetFullFood {
		foodtotake = targetFullFood
	}
	a.Log("VWXYZ:", infra.Fields{"food": foodtotake})
	// Eat Food and Update Statuses
	eaten, err := a.TakeFood(foodtotake)

	if err != nil {
		switch err.(type) {
		case *infra.FloorError:
		case *infra.NegFoodError:
		case *infra.AlreadyEatenError:
		default:
		}
	}

	a.opMem.foodEaten = eaten
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
func (a *CustomAgent7) foodAverage() {

	// If the number of days on floor is zero set num of days to 1 and add current avail food to average
	if (a.opMem.prevFloors[a.Floor()].numOfDays == 0) && a.PlatformOnFloor() {
		a.opMem.prevFloors[a.Floor()] = floorInfo{1, int(a.CurrPlatFood())}
	} else if a.PlatformOnFloor() {
		// Otherwise update number of days on floor and calculate new average
		tmp := a.opMem.prevFloors[a.Floor()].avgFood*a.opMem.prevFloors[a.Floor()].numOfDays + int(a.CurrPlatFood())
		newNumOfDays := a.opMem.prevFloors[a.Floor()].numOfDays + 1
		tmp = tmp / newNumOfDays
		a.opMem.prevFloors[a.Floor()] = floorInfo{tmp, newNumOfDays}
	}
}

func (a *CustomAgent7) adjustMood() {
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

}

func (a *CustomAgent7) manageFloor() {
	if len(a.opMem.orderPrevFloors) == 0 || a.opMem.orderPrevFloors[len(a.opMem.orderPrevFloors)-1] != a.Floor() { //If day 1 or floor change

		a.Log("checking floor change........................................................................................")
		if a.Age() == 0 {
			a.Log("Personality Makeup:", infra.Fields{"A": a.personality.agreeableness, "C": a.personality.conscientiousness, "E": a.personality.extraversion, "N": a.personality.neuroticism, "O": a.personality.openness})
		}

		// Reset Variables on Floor Change
		a.opMem.currentDayonFloor = 1 // reset currentDay counter
		a.opMem.receivedReq = false   // reset Requests
		a.opMem.takeFood = 0
		a.opMem.leaveFood = 0

		// If the floor tracker is not empty
		if len(a.opMem.orderPrevFloors) != 0 {
			// ------------------ Run() Block A.1: Estimates the food available on the new floor based on past experience ------------------
			a.foodNewFloor()

		}
		a.opMem.orderPrevFloors = append(a.opMem.orderPrevFloors, a.Floor()) //append current floor to floor tracker

	} else { //increment currentDayonFloor counter
		a.opMem.currentDayonFloor++
	}
}

func (a *CustomAgent7) foodNewFloor() {
	closestFloorAboveCurrent := 0
	closestFloorBelowCurrent := math.Inf(1)
	beenOnCurrentFloorBefore := false

	for i := 0; i < len(a.opMem.orderPrevFloors); i++ {
		if a.opMem.orderPrevFloors[i] == a.Floor() {
			beenOnCurrentFloorBefore = true
			break
		}
	}

	if !beenOnCurrentFloorBefore {
		for i := 0; i < len(a.opMem.orderPrevFloors); i++ {
			if a.opMem.orderPrevFloors[i] < a.Floor() && a.opMem.orderPrevFloors[i] > closestFloorAboveCurrent {
				closestFloorAboveCurrent = a.opMem.orderPrevFloors[i]
			}
		}
		for i := 0; i < len(a.opMem.orderPrevFloors); i++ {
			if a.opMem.orderPrevFloors[i] > a.Floor() && a.opMem.orderPrevFloors[i] < int(closestFloorBelowCurrent) {
				closestFloorBelowCurrent = float64(a.opMem.orderPrevFloors[i])
			}
		}
	}

	var expectedFood int

	if beenOnCurrentFloorBefore {
		expectedFood = a.opMem.prevFloors[a.Floor()].avgFood
	} else {

		if closestFloorAboveCurrent != 0 && closestFloorBelowCurrent != math.Inf(1) {
			expectedFood = a.opMem.prevFloors[int(closestFloorBelowCurrent)].avgFood + (a.opMem.prevFloors[closestFloorAboveCurrent].avgFood-a.opMem.prevFloors[int(closestFloorBelowCurrent)].avgFood)*(int(closestFloorBelowCurrent)-a.Floor())/(int(closestFloorBelowCurrent)-int(closestFloorAboveCurrent))
		}
		if closestFloorAboveCurrent == 0 && closestFloorBelowCurrent != math.Inf(1) {
			if a.opMem.prevFloors[int(closestFloorBelowCurrent)].avgFood == 0 {
				expectedFood = 0
			} else {
				estimatedMealSize := a.HealthInfo().HPReqCToW + (a.personality.openness / 5)
				expectedFood = a.opMem.prevFloors[int(closestFloorBelowCurrent)].avgFood + (int(closestFloorBelowCurrent)-a.Floor())*estimatedMealSize
			}
		}
		if closestFloorAboveCurrent != 0 && closestFloorBelowCurrent == math.Inf(1) {
			estimatedMealSize := a.HealthInfo().HPReqCToW + ((100 - a.personality.openness) / 5)
			expectedFood = a.opMem.prevFloors[closestFloorAboveCurrent].avgFood - (a.Floor()-closestFloorAboveCurrent)*estimatedMealSize
		}
		if expectedFood < 0 {
			expectedFood = 0
		}

	}

	// ------------------ Run() Block A.2: Adjusts mood w.r.t. change in floor and the expected food ------------------
	a.floorChangeMood(expectedFood)

}
func (a *CustomAgent7) floorChangeMood(expectedFood int) {
	// Update Behaviour: Higher in the tower (lower floor number) increases greed/decreases kindess and vice versa
	if a.Floor() < a.opMem.orderPrevFloors[len(a.opMem.orderPrevFloors)-1] {
		// A higher openess will more significantly reduce greedines and increase kindness
		a.behaviour.greediness -= int(float64(a.opMem.orderPrevFloors[len(a.opMem.orderPrevFloors)-1]/a.Floor()) * (float64(a.personality.openness) / 50))
		a.behaviour.kindness += int(float64(a.opMem.orderPrevFloors[len(a.opMem.orderPrevFloors)-1]/a.Floor()) * (float64(a.personality.openness) / 50))
	} else {
		// A lower openess will more significantly increase greedines and reduce kindness
		a.behaviour.greediness += int(float64(a.Floor()/a.opMem.orderPrevFloors[len(a.opMem.orderPrevFloors)-1]) * (float64(100-a.personality.openness) / 50))
		a.behaviour.kindness -= int(float64(a.Floor()/a.opMem.orderPrevFloors[len(a.opMem.orderPrevFloors)-1]) * (float64(100-a.personality.openness) / 50))
	}

	// Greediness is affected by the amount of food expected on this floor as estimated from past experience
	a.behaviour.greediness -= int(float64(a.opMem.currentFloorRisk) * (float64(a.personality.conscientiousness) / 50))
	foodReqCtoW := int(health.FoodRequired(a.HealthInfo().HPCritical, a.HealthInfo().WeakLevel, a.HealthInfo()))
	if expectedFood < foodReqCtoW { //greediness only begins to be effected once the expected food is below a crticial threshold
		a.opMem.currentFloorRisk = (foodReqCtoW - expectedFood)
		a.behaviour.greediness += int(float64(a.opMem.currentFloorRisk) * (float64(a.personality.conscientiousness) / 50))
	} else {
		a.opMem.currentFloorRisk = 0
	}
	a.treatyOnFloorChange()
}

func (a *CustomAgent7) manageMessages() {
	// ------------------ Run() Block B: Receive messages ------------------
	receivedMsg := a.ReceiveMessage()
	if receivedMsg != nil {
		receivedMsg.Visit(a)
	} else {
		a.Log("No Messages")
	}

	// Messaging
	if a.opMem.prevAge < a.Age() {
		a.opMem.prevAge = a.Age()
		a.opMem.msg1Sent = false
		a.opMem.msg2Sent = false
		a.opMem.msg3Sent = false
	}

	if !a.opMem.msg1Sent {
		msg := messages.NewAskHPMessage(a.ID(), a.Floor(), a.Floor()+1)
		// a.opMem.messageQueue = append(a.opMem.messageQueue, msg)
		a.SendMessage(msg)
		if a.Floor() != 1 {
			msg = messages.NewAskHPMessage(a.ID(), a.Floor(), a.Floor()-1)
			a.SendMessage(msg)
		}
		a.opMem.msg1Sent = true
	}

	if a.PlatformOnFloor() && !a.opMem.msg2Sent {
		if a.Floor() != 1 {
			msg2 := messages.NewAskFoodTakenMessage(a.ID(), a.Floor(), a.Floor()-1)
			a.SendMessage(msg2)
		}
		a.opMem.msg2Sent = true
	}

	if !a.PlatformOnFloor() && a.CurrPlatFood() != -1 && a.opMem.seenPlatform && a.Floor() != 1 && !a.opMem.msg3Sent {
		msg2 := messages.NewAskFoodTakenMessage(a.ID(), a.Floor(), a.Floor()+1)
		a.SendMessage(msg2)
		a.opMem.msg3Sent = true
	}

}

func (a *CustomAgent7) manageNewDay() {

	a.opMem.seenPlatform = false
	a.opMem.prevHP = a.HP()
	a.opMem.foodEaten = 0
	a.opMem.receivedReq = false

}

func (a *CustomAgent7) treatyGT(tActive messages.Treaty, foodtotake food.FoodType, available food.FoodType, amount food.FoodType) food.FoodType {

	if a.HP() > tActive.ConditionValue() || int(a.CurrPlatFood()) > tActive.ConditionValue() {
		switch tActive.RequestOp() {
		case messages.GT:
			if foodtotake > available-food.FoodType(tActive.RequestValue()) {
				foodtotake = available - food.FoodType(amount) - 1
			}
		case messages.GE:
			if foodtotake >= available-food.FoodType(tActive.RequestValue()) {
				foodtotake = available - food.FoodType(amount)
			}
		case messages.EQ:
			if foodtotake != available-food.FoodType(tActive.RequestValue()) {
				foodtotake = available - food.FoodType(amount)
			}
		}
	}
	return foodtotake

}
func (a *CustomAgent7) treatyGE(tActive messages.Treaty, foodtotake food.FoodType, available food.FoodType, amount food.FoodType) food.FoodType {

	if a.HP() >= tActive.ConditionValue() || int(a.CurrPlatFood()) >= tActive.ConditionValue() {
		switch tActive.RequestOp() {
		case messages.GT:
			if foodtotake > available-food.FoodType(tActive.RequestValue()) {
				foodtotake = available - food.FoodType(amount) - 1
			}
		case messages.GE:
			if foodtotake >= available-food.FoodType(tActive.RequestValue()) {
				foodtotake = available - food.FoodType(amount)
			}
		case messages.EQ:
			if foodtotake != available-food.FoodType(tActive.RequestValue()) {
				foodtotake = available - food.FoodType(amount)
			}
		}
	}
	return foodtotake

}
func (a *CustomAgent7) treatyEQ(tActive messages.Treaty, foodtotake food.FoodType, available food.FoodType, amount food.FoodType) food.FoodType {

	if a.HP() == tActive.ConditionValue() || int(a.CurrPlatFood()) == tActive.ConditionValue() || a.Floor() == tActive.ConditionValue() {
		switch tActive.RequestOp() {
		case messages.GT:
			if foodtotake > available-food.FoodType(tActive.RequestValue()) {
				foodtotake = available - food.FoodType(amount) - 1
			}
		case messages.GE:
			if foodtotake >= available-food.FoodType(tActive.RequestValue()) {
				foodtotake = available - food.FoodType(amount)
			}
		case messages.EQ:
			if foodtotake != available-food.FoodType(tActive.RequestValue()) {
				foodtotake = available - food.FoodType(amount)
			}
		}
	}
	return foodtotake

}
func (a *CustomAgent7) treatyLT(tActive messages.Treaty, foodtotake food.FoodType, available food.FoodType, amount food.FoodType) food.FoodType {

	if a.Floor() < tActive.ConditionValue() {
		switch tActive.RequestOp() {
		case messages.GT:
			if foodtotake > available-food.FoodType(tActive.RequestValue()) {
				foodtotake = available - food.FoodType(amount) - 1
			}
		case messages.GE:
			if foodtotake >= available-food.FoodType(tActive.RequestValue()) {
				foodtotake = available - food.FoodType(amount)
			}
		case messages.EQ:
			if foodtotake != available-food.FoodType(tActive.RequestValue()) {
				foodtotake = available - food.FoodType(amount)
			}
		}
	}
	return foodtotake

}
func (a *CustomAgent7) treatyLE(tActive messages.Treaty, foodtotake food.FoodType, available food.FoodType, amount food.FoodType) food.FoodType {

	if a.Floor() <= tActive.ConditionValue() {
		switch tActive.RequestOp() {
		case messages.GT:
			if foodtotake > available-food.FoodType(tActive.RequestValue()) {
				foodtotake = available - food.FoodType(amount) - 1
			}
		case messages.GE:
			if foodtotake >= available-food.FoodType(tActive.RequestValue()) {
				foodtotake = available - food.FoodType(amount)
			}
		case messages.EQ:
			if foodtotake != available-food.FoodType(tActive.RequestValue()) {
				foodtotake = available - food.FoodType(amount)
			}
		}
	}

	return foodtotake
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

// Requests
func (a *CustomAgent7) HandleRequestTakeFood(msg messages.RequestTakeFoodMessage) {
	a.Log("Recieved requestTakeFood message from ", infra.Fields{"floor": msg.SenderFloor()})
	reqFood := msg.Request()
	if a.personality.extraversion > 5 {
		if reqFood > int(health.FoodRequired(a.HP(), a.HealthInfo().MaxHP, a.HealthInfo())) || a.DaysAtCritical() >= (a.HealthInfo().MaxDayCritical-3) || a.behaviour.greediness > a.behaviour.kindness {
			reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), false)
			a.SendMessage(reply)
		} else {
			a.opMem.receivedReq = true
			a.opMem.takeFood = reqFood
			reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), true)
			a.SendMessage(reply)
		}
	}
}

func (a *CustomAgent7) HandleRequestLeaveFood(msg messages.RequestLeaveFoodMessage) {
	a.Log("Recieved requestLeaveFood message from ", infra.Fields{"floor": msg.SenderFloor()})
	reqFoodLeft := msg.Request()
	if a.personality.extraversion > 5 {
		if a.DaysAtCritical() >= (a.HealthInfo().MaxDayCritical-3) || a.behaviour.greediness > a.behaviour.kindness {
			// a.opMem.receivedReq = false
			reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), false)
			a.SendMessage(reply)
		} else {
			a.opMem.receivedReq = true
			a.opMem.leaveFood = reqFoodLeft
			reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), true)
			a.SendMessage(reply)
		}
	}
}

// Responses
func (a *CustomAgent7) HandleResponse(msg messages.BoolResponseMessage) {
	response := msg.Response()
	a.Log("I received a Response message from ", infra.Fields{"floor": msg.SenderFloor(), "response": response})
}
