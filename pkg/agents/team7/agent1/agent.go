package team7agent1

import (
	"math/rand"

	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
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

type team7Personalities struct {
	openness          int
	conscientiousness int
	extraversion      int
	agreeableness     int
	neuroticism       int
}
type floorInfo struct {
	numOfDays, avgFood int
}
type OperationalMemory struct {
	orderPrevFloors   []int
	prevFloors        map[int]floorInfo
	currentDayonFloor int
	numOfDays         int
	avgFood           int
	daysHungry        int
	seenPlatform      bool
	foodEaten         food.FoodType
	prevHP            int
}
type CurrentBehaviour struct { //determines directly the food taken
	greediness int
	kindness   int
	// likehood to read and send messages
	// affacted by openess and extraversion
	activeness int
	// affected by days on floor, days hungry
	// changes less frequently than sanity
	morality int
	// affected by days on floor, previous floor, days hungry and frequency of communictaion(extraversion),
	// neuroticism
	// more sensitive to other factors
	// lower sanity -> less logical behaviour -> (ex: no food taken when starving ?)
	sanity int
}
type CustomAgent7 struct {
	*infra.Base
	personality team7Personalities
	opMem       OperationalMemory
	behaviour   CurrentBehaviour
	// greediness   int
	// kindness     int
	// daysHungry   int
	// seenPlatform bool
	Eaten food.FoodType
	// prevHP       int
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
			greediness: 100 - agreeableness,
			kindness:   agreeableness,
		},

		opMem: OperationalMemory{
			orderPrevFloors:   []int{},
			prevFloors:        prevFloors,
			currentDayonFloor: 0,
			numOfDays:         0,
			avgFood:           0,
			daysHungry:        0,
			seenPlatform:      false,
			foodEaten:         0,
			prevHP:            100,
		},
		Eaten: 0,
	}, nil
}

func (a *CustomAgent7) Run() {

	a.Log("Team7Agent1 reporting status:", infra.Fields{"floor": a.Floor(), "hp": a.HP(), "greed": a.behaviour.greediness, "kind": a.behaviour.kindness, "aggr": a.personality.agreeableness})

	// Operational Variables
	// UserID := a.ID()
	// currentHP := a.HP()
	currentAvailFood := a.CurrPlatFood()
	currentFloor := a.Floor()
	prevFloor := a.opMem.orderPrevFloors[len(a.opMem.orderPrevFloors)-1]

	//Check if floor has changed
	if len(a.opMem.orderPrevFloors) == 0 || prevFloor != currentFloor { //If day 1 or floor change
		a.opMem.currentDayonFloor = 1          //reset currentDay counter
		if len(a.opMem.orderPrevFloors) != 0 { //if the floor tracker is not empty

			//Update Behaviour
			if currentFloor < prevFloor {
				//only negatively impacted if openness is low
				if a.personality.openness <= 50 {
					a.behaviour.greediness += (currentFloor - prevFloor)
				}
			} else {
				//if we have moved up our kindness increases
				a.behaviour.kindness += (prevFloor - currentFloor)
			}

		}
		a.opMem.orderPrevFloors = append(a.opMem.orderPrevFloors, currentFloor) //append current floor to floor tracker
	} else { //increment currentDayonFloor counter
		a.opMem.currentDayonFloor++
	}

	//Average Food available per floor map
	// If the number of days on floor is zero set num of days to 1 and add current avail food to average
	if (a.opMem.prevFloors[currentFloor].numOfDays == 0) && (currentAvailFood != -1) {
		a.opMem.prevFloors[currentFloor] = floorInfo{1, int(currentAvailFood)}
	} else if currentAvailFood != -1 {
		//otherwise update number of days on floor and calculate new average
		tmp := a.opMem.prevFloors[currentFloor].avgFood * a.opMem.prevFloors[currentFloor].numOfDays
		tmp = tmp + int(currentAvailFood)
		newNumOfDays := a.opMem.prevFloors[currentFloor].numOfDays + 1
		tmp = tmp / newNumOfDays
		a.opMem.prevFloors[currentFloor] = floorInfo{tmp, newNumOfDays}
	}

	// Receive Message
	receivedMsg := a.ReceiveMessage()
	if receivedMsg != nil {
		receivedMsg.Visit(a)
	} else {
		a.Log("No Messages")
	}

	// Sending Messages
	msg := messages.NewRequestLeaveFoodMessage(currentFloor, 100)
	a.SendMessage(1, msg)
	a.SendMessage(-1, msg)

	//eat

	//choses random number between -5 and 5
	if a.CurrPlatFood() != -1 {

		r1 := rand.Intn(11) - 5
		r2 := rand.Intn(11) - 5

		a.behaviour.kindness += (a.personality.neuroticism * r1) / 50
		a.behaviour.greediness += (a.personality.neuroticism * r2) / 50

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

		foodtotake := food.FoodType(100 - a.behaviour.kindness + a.behaviour.greediness)

		//has the agent seen the platform
		if a.CurrPlatFood() != -1 && !a.opMem.seenPlatform {
			//foodEaten, err := a.TakeFood(foodtotake)
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

		if (a.CurrPlatFood() == -1 && a.opMem.seenPlatform) || a.CurrPlatFood() == 100 {
			//Get ready for new day
			a.opMem.seenPlatform = false
			a.opMem.prevHP = a.HP()
		}

	}

	//send Message
}

// Message Handlers
// Asks
func (a *CustomAgent7) HandleAskHP(msg messages.AskHPMessage) {
	reply := msg.Reply(a.Floor(), a.HP())
	a.SendMessage(msg.SenderFloor()-a.Floor(), reply)
	a.Log("I recieved an askHP message from ", infra.Fields{"floor": msg.SenderFloor()})
}

func (a *CustomAgent7) HandleAskFoodTaken(msg messages.AskFoodTakenMessage) {
	reply := msg.Reply(a.Floor(), 10)
	a.SendMessage(msg.SenderFloor()-a.Floor(), reply)
	a.Log("I recieved an askFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor()})
}

func (a *CustomAgent7) HandleAskIntendedFoodTaken(msg messages.AskIntendedFoodIntakeMessage) {
	reply := msg.Reply(a.Floor(), 11)
	a.SendMessage(msg.SenderFloor()-a.Floor(), reply)
	a.Log("I recieved an askIntendedFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor()})
}

// Requests
func (a *CustomAgent7) HandleRequestLeaveFood(msg messages.RequestLeaveFoodMessage) {
	reply := msg.Reply(a.Floor(), true)
	a.SendMessage(msg.SenderFloor()-a.Floor(), reply)
	a.leaveFood = food.FoodType(msg.Request())
	a.Log("I recieved a requestLeaveFood message from ", infra.Fields{"floor": msg.SenderFloor()})
}

func (a *CustomAgent7) HandleRequestTakeFood(msg messages.RequestTakeFoodMessage) {
	reply := msg.Reply(a.Floor(), true)
	a.SendMessage(msg.SenderFloor()-a.Floor(), reply)
	a.Log("I recieved a requestTakeFood message from ", infra.Fields{"floor": msg.SenderFloor()})
}

// Responses
func (a *CustomAgent7) HandleResponse(msg messages.BoolResponseMessage) {
	response := msg.Response()
	a.Log("I recieved a Response message from ", infra.Fields{"floor": msg.SenderFloor(), "response": response})
}

// Statements
func (a *CustomAgent7) HandleStateFoodTaken(msg messages.StateFoodTakenMessage) {
	statement := msg.Statement()
	a.Log("I recieved a StateFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor(), "statement": statement})
}

func (a *CustomAgent7) HandleStateHP(msg messages.StateHPMessage) {
	statement := msg.Statement()
	a.Log("I recieved a StateHP message from ", infra.Fields{"floor": msg.SenderFloor(), "statement": statement})
}

func (a *CustomAgent7) HandleStateIntendedFoodTaken(msg messages.StateIntendedFoodIntakeMessage) {
	statement := msg.Statement()
	a.Log("I recieved a StateIntendedFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor(), "statement": statement})
}
