package team7agent1

import (
	"math"
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
	daysHungry        int
	seenPlatform      bool
	prevHP            int
	leaveFood         int
	foodEaten         food.FoodType
}

type CustomAgent7 struct {
	*infra.Base
	personality team7Personalities
	opMem       OperationalMemory
	behaviour   CurrentBehaviour
	Eaten       food.FoodType
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
	// currentHP := a.HP()

	//****experimental
	//TargetHP initialise by 70 +/- random scaling * +/- Neuroticism
	//TargetHP := 70 //+/- Neurpticism Scaled
	healthInfo := a.HealthInfo()
	//****

	currentAvailFood := a.CurrPlatFood()
	currentFloor := a.Floor()
	// prevFloor := a.opMem.orderPrevFloors[len(a.opMem.orderPrevFloors)-1]

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

			//How the above impacts stuff

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
			if expectedFood < healthInfo.HPReqCToW { //greediness only begins to be effected once the expected food is below a crticial threshold
				a.behaviour.greediness += (healthInfo.HPReqCToW - expectedFood)
			}

		}
		a.opMem.orderPrevFloors = append(a.opMem.orderPrevFloors, currentFloor) //append current floor to floor tracker
	} else { //increment currentDayonFloor counter
		a.opMem.currentDayonFloor++
	}

	//Average Food available per floor map
	// If the number of days on floor is zero set num of days to 1 and add current avail food to average
	if (a.opMem.prevFloors[currentFloor].numOfDays == 0) && (currentAvailFood != -1) && a.PlatformOnFloor() {
		a.opMem.prevFloors[currentFloor] = floorInfo{1, int(currentAvailFood)}
	} else if currentAvailFood != -1 && a.PlatformOnFloor() {
		//otherwise update number of days on floor and calculate new average
		tmp := a.opMem.prevFloors[currentFloor].avgFood * a.opMem.prevFloors[currentFloor].numOfDays
		tmp = tmp + int(currentAvailFood)
		newNumOfDays := a.opMem.prevFloors[currentFloor].numOfDays + 1
		tmp = tmp / newNumOfDays
		a.opMem.prevFloors[currentFloor] = floorInfo{tmp, newNumOfDays}
	}

	// Receive Messages
	receivedMsg := a.ReceiveMessage()
	if receivedMsg != nil {
		receivedMsg.Visit(a)
	} else {
		a.Log("No Messages")
	}

	// Eat

	//
	if a.CurrPlatFood() != -1 && a.PlatformOnFloor() {

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

		// ///******///
		// currentHealth := a.HP()

		// // Possible food to take: Food available on platform/Food required to go from ciritcal to weak, critical to healthy

		// FoodToWeakHP := healthInfo.HPReqCToW
		// currentAvailFood := a.CurrPlatFood()

		// // Influences on Food taking: greediness, kindness, memory map HOW???, Sanity [depends on messaging, extraversion???, neuroticism], activeness, morality,
		// // so if your extraverted: less interaction leads to lowered sanity/More interaction better sanity
		// // higher neuroticism -> fall faster into insanity, initialise sanity based on neuroticism, sanity = 100-neuroticism/ sanity bsaed on
		// //

		// if a.opMem.treaty {
		// 	object = treay
		// }

		// switch object {
		// case treaty
		// }

		// if treaty {
		// 	foodtotake := 100 - foodtoleave + behavioural_adjustments
		// } else if request {
		// 	foodtotake := 100 - requested_leave_food
		// } else {
		// 	foodtotake := (foodreqCriticaltoWeak + weaktoHealthy) + behavioural_adjustments
		// }

		foodtotake := food.FoodType((100) - a.behaviour.kindness + a.behaviour.greediness)

		// TREATIES!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		// ///********///

		// has the agent seen the platform
		if a.CurrPlatFood() != -1 && a.PlatformOnFloor() {
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

}

// Message Handlers

// Handle Asks and auto respond
func (a *CustomAgent7) HandleAskHP(msg messages.AskHPMessage) {
	a.Log("Recieved askHP message from ", infra.Fields{"floor": msg.SenderFloor()})
	reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor()-a.Floor(), a.HP())
	// reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), a.HP())
	// a.SendMessage(msg.SenderFloor()-a.Floor(), reply)
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


//Treaties

func (a *CustomAgent7) RejectTreaty(msg messages.ProposeTreatyMessage) {
	reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), false)
	a.SendMessage(reply)
	a.Log("Rejected treaty", infra.Fields{"proposerID": msg.SenderID(), "proposerFloor": msg.SenderFloor(), "treatyID": msg.TreatyID()})
}


func (a *CustomAgent7) HandleProposeTreaty(msg messages.ProposeTreatyMessage) {
	// The code below can be used to accept all treaties by default.
	treaty := msg.Treaty()
	treaty.SignTreaty()
	a.activeTreaties[msg.TreatyID()] = treaty
	reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), true)
	a.SendMessage(reply)
	a.Log("Accepted treaty", infra.Fields{"proposerID": msg.SenderID(), "proposerFloor": msg.SenderFloor(),
		"treatyID": msg.TreatyID()})
}

// This implementation automatically increments the signature count of the treaty if it was accepted.
func (a *CustomAgent7) HandleTreatyResponse(msg messages.TreatyResponseMessage) {
	if msg.Response() {
		treaty := a.activeTreaties[msg.TreatyID()]
		treaty.SignTreaty()
		a.activeTreaties[msg.TreatyID()] = treaty
	}
}

func (a *CustomAgent7) HandlePropogate(msg messages.Message) {
	a.SendMessage(msg)
}


// func (a *CustomAgent7) RejectTreaty(msg messages.ProposeTreatyMessage) {
// 	reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), false)
// 	a.SendMessage(reply)
// 	a.Log("Rejected treaty", infra.Fields{"proposerID": msg.SenderID(), "proposerFloor": msg.SenderFloor(), "treatyID": msg.TreatyID()})
// }

// func (a *CustomAgent7) HandleProposeTreaty(msg messages.ProposeTreatyMessage) {
// 	treaty := msg.Treaty()
// 	if treaty.Request() == messages.Inform {
// 		// Reject any inform treaties as I'm not sure what they are or how to handle them
// 		a.RejectTreaty(msg)
// 		return
// 	}

// 	if a.treatyConflicts(treaty) {
// 		// Reject treaty if it conflicts with any other active treaty
// 		a.RejectTreaty(msg)
// 		return
// 	}

// 	if (treaty.Request() == messages.LeaveAmountFood && treaty.RequestOp() == messages.EQ) || treaty.RequestOp() == messages.LE || treaty.RequestOp() == messages.LT {
// 		//Reject all treaties that ask you to leave less food, don't see why you would do this
// 		//Reject any treaty that asks you to leave a specific amount of food as this would lead to multiple people just not eating as only 1 agent can eat and fufill that criteria
// 		a.RejectTreaty(msg)
// 		return
// 	}

// 	if (treaty.ConditionOp() == messages.LE || treaty.ConditionOp() == messages.LT) && (treaty.Condition() == messages.HP || treaty.Condition() == messages.AvailableFood) {
// 		//Reject treaties that have a condition of being less than something
// 		// Why agree to treaty which is bounded by you having lower HP
// 		// Why agree to treaty which is bounded by the amount of food on platform being low
// 		// All of these are conditions where you will be desperate for food so shouldn't be limitting yourself if you do get the opurrtunity to eat
// 		a.RejectTreaty(msg)
// 		return
// 	}

// 	if (treaty.ConditionOp() == messages.GE || treaty.ConditionOp() == messages.GT) && treaty.Condition() == messages.Floor {
// 		// Why agree to treaty which is bounded by you being lower in tower
// 		a.RejectTreaty(msg)
// 		return
// 	}

// 	if treaty.Condition() == messages.HP && treaty.ConditionValue() < a.HealthInfo().WeakLevel*2 {
// 		//Reject any HP condition based treaty if the condition is too strict, and you would be put into critical by not eating
// 		a.RejectTreaty(msg)
// 		return
// 	}

// 	if treaty.Condition() == messages.Floor && treaty.ConditionValue() != 1 && treaty.Duration() >= a.HealthInfo().MaxDayCritical {
// 		//Reject any floor condition that involves us not being on the top floor and lasts for more than days you can survive in critical
// 		//This is because there is a risk of signing your own death if you agree as you may be forced to eat no food with no get out condition.
// 		a.RejectTreaty(msg)
// 		return
// 	}

// 	if treaty.Condition() == messages.AvailableFood {
// 		if treaty.Request() == messages.LeaveAmountFood && treaty.ConditionValue()-treaty.RequestValue() < 3 {
// 			// Reject treaty in which you would not be able to eat enough food to avoid critical level
// 			a.RejectTreaty(msg)
// 			return
// 		}
// 		if treaty.Request() == messages.LeavePercentFood && treaty.ConditionValue()*int(float64(100-treaty.RequestValue())/100) < 3 {
// 			// Reject treaty in which you would not be able to eat enough food to avoid critical level
// 			a.RejectTreaty(msg)
// 			return
// 		}
// 	}

// 	//Now we have ruled treaties that are always unacceptable, now to decide if agent will agree to an acceptable treaty
// 	//TODO: Develop the decision calculation.
// 	//For now: Middle range of selfishness - selfishness + opinion of agent proposing treaty
// 	decision := 5 - a.selfishness + a.socialMemory[msg.SenderID()].favour
// 	if decision > 0 {
// 		treaty.SignTreaty()
// 		a.AddTreaty(treaty)
// 		reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), true)
// 		a.SendMessage(reply)
// 		a.Log("Accepted treaty", infra.Fields{"proposerID": msg.SenderID(), "proposerFloor": msg.SenderFloor(), "treatyID": msg.TreatyID()})
// 		passingOnTo := a.Floor() + 1
// 		if msg.SenderFloor() > a.Floor() {
// 			passingOnTo = a.Floor() - 1
// 		}
// 		passItOn := messages.NewProposalMessage(a.ID(), a.Floor(), passingOnTo, treaty)
// 		a.SendMessage(passItOn)
// 	} else {
// 		a.RejectTreaty(msg)
// 	}

// 	// The code below can be used to accept all treaties by default.
// 	// treaty := msg.Treaty()
// 	// treaty.SignTreaty()
// 	// a.activeTreaties[msg.TreatyID()] = treaty
// 	// reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), true)
// 	// a.SendMessage(reply)
// 	// a.Log("Accepted treaty", Fields{"proposerID": msg.SenderID(), "proposerFloor": msg.SenderFloor(),
// 	// 	"treatyID": msg.TreatyID()})
// }

// func (a *CustomAgent7) HandleTreatyResponse(msg messages.TreatyResponseMessage) {
// 	if msg.Response() {
// 		treaty := a.ActiveTreaties()[msg.TreatyID()]
// 		treaty.SignTreaty()
// 		a.ActiveTreaties()[msg.TreatyID()] = treaty
// 		a.addToSocialFavour(msg.SenderID(), a.socialMemory[msg.SenderID()].favour+2)
// 	}
// }