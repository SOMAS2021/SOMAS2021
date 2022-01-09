package team6

import (
	"math"
	"math/rand"

	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
)

// func (a *CustomAgent6) foodRange() (food.FoodType, food.FoodType) {
// 	mini := food.FoodType(0)           //initial minimum value
// 	maxi := food.FoodType(math.MaxInt) //maximum value to indicate no maximum
// 	for _, treaty := range a.ActiveTreaties() {
// 		// TODO: check if we actually know the amount of food on the platform every time foodRange is called
// 		takefoodAmount := a.convertToTakeFoodAmount(float64(a.CurrPlatFood()), treaty.Request(), treaty.RequestValue())
// 		foodAmount := a.CurrPlatFood() - takefoodAmount

// 		// deal with different request types
// 		switch treaty.RequestOp() {
// 		case 1: // GE
// 			if foodAmount > mini || foodAmount == mini {
// 				mini = foodAmount + 1 //If Greater than, then the max value is one larger
// 			}
// 		case 2:
// 			if foodAmount > mini {
// 				mini = foodAmount //If Greater or Equal, then the max value is itself
// 			}
// 		case 3:
// 			eqVal := foodAmount //if Equal then find the value as there can only be one equal operator unless the values in both treaties are the same
// 			mini, maxi = eqVal, eqVal
// 		case 4:
// 			if foodAmount < maxi || maxi == 0 {
// 				maxi = foodAmount //If Less than or Equal, then the max value is itself
// 			}
// 		case 5:
// 			if foodAmount < maxi || maxi == 0 {
// 				maxi = foodAmount - 1 //If Less than, then the max value is one smaller
// 			}
// 		default:
// 			mini, maxi = -1, -1 //unknown op code
// 		}

// 	}

// !!! We only accept treaties that say leave ≥,> of food available. Hence there will never be a invalid treaty.
// func (a *CustomAgent6) treatyValid(treaty *messages.Treaty) bool {

// 	if len(a.ActiveTreaties()) == 0 {
// 		return true
// 	}

// 	chkTrtyVal := foodAmount
// 	mini, maxi := a.foodRange()

// 	if chkTrtyVal >= mini && chkTrtyVal <= maxi {
// 		return true
// 	}
// 	return false

// }

// The utility function with
// x - food input
// z - desired food (maximum of the function)
func Utility(x, z float64, socialMotive string) float64 {

	params := NewUtilityParams(socialMotive)
	// calculate the function scaling parameter a

	if socialMotive == "Altruist" /*|| socialMotive == "Nacissist"*/ {
		// Don't scale depending on desired food
		return params.g*math.Pow(x, 1/params.r) - params.c*x
	} else {
		a := (1 / z) * math.Pow((params.c*params.r)/params.g, params.r/(1-params.r))
		return params.g*math.Pow(a*x, 1/params.r) - params.c*a*x
	}
}

func min(x, y food.FoodType) food.FoodType {
	if x < y {
		return x
	}
	return y
}

// Evaluate our agents current utility based on the current desired food
func (a *CustomAgent6) evaluateUtility(mem memory) float64 {
	sum := food.FoodType(0)
	for _, foodAvailable := range mem {
		sum += food.FoodType(Utility(float64(min(foodAvailable, a.desiredFoodIntake())), float64(a.desiredFoodIntake()), a.currBehaviour.String()))
	}

	return float64(sum) / math.Max(float64(len(mem)), 1.0)
}

// convert different message "Request types" to an equivalent "food intake" value
func (a *CustomAgent6) convertToTakeFoodAmount(foodAvailable float64, requestType messages.RequestType, requestValue int) food.FoodType {

	takeFood := 0.0
	switch requestType {
	case messages.LeaveAmountFood:
		takeFood = foodAvailable - float64(requestValue)
	case messages.LeavePercentFood:
		takeFood = foodAvailable * (1.0 - float64(requestValue))
	// case messages.TakeAmountFood:
	// 	takeFood = float64(requestValue)
	case messages.Inform:
		takeFood = -1
	default:
		takeFood = -1
	}

	return food.FoodType(takeFood)
}

//  Decides if to accept or reject a treaty
func (a *CustomAgent6) considerTreaty(t *messages.Treaty) bool {
	levels := levelsData{
		strongLevel:  a.HealthInfo().MaxHP * 3 / 5,
		healthyLevel: a.HealthInfo().MaxHP * 3 / 10,
		weakLevel:    a.HealthInfo().WeakLevel,
		critLevel:    0,
	}

	if t.RequestOp() == messages.GE || t.RequestOp() == messages.GT {
		switch t.Condition() {
		// HP
		case 1:
			switch a.currBehaviour.String() {
			case "Altruist":
				return true
			case "Collectivist":
				if t.ConditionOp() == messages.LE || t.ConditionOp() == messages.LT {
					a.Log("I used considerTreatyUsingUtility", infra.Fields{"my floor:": a.Floor(), "my social motive": a.currBehaviour.String()})
					return a.considerTreatyUsingUtility(t)
				} else {
					if t.ConditionValue() >= a.HealthInfo().WeakLevel {
						return true
					} else {
						a.Log("I used considerTreatyUsingUtility", infra.Fields{"my floor:": a.Floor(), "my social motive": a.currBehaviour.String()})
						return a.considerTreatyUsingUtility(t)
					}
				}

			case "Selfish":
				if t.ConditionOp() == messages.LE || t.ConditionOp() == messages.LT {
					a.Log("I used considerTreatyUsingUtility", infra.Fields{"my floor:": a.Floor(), "my social motive": a.currBehaviour.String()})
					return a.considerTreatyUsingUtility(t)
				} else {
					if t.ConditionValue() >= levels.strongLevel {
						return true
					} else {
						a.Log("I used considerTreatyUsingUtility", infra.Fields{"my floor:": a.Floor(), "my social motive": a.currBehaviour.String()})
						return a.considerTreatyUsingUtility(t)

					}
				}
			case "Narcissist":
				a.Log("I used considerTreatyUsingUtility", infra.Fields{"my floor:": a.Floor(), "my social motive": a.currBehaviour.String()})
				return a.considerTreatyUsingUtility(t)
			default:
				a.Log("I used considerTreatyUsingUtility", infra.Fields{"my floor:": a.Floor(), "my social motive": a.currBehaviour.String()})
				return a.considerTreatyUsingUtility(t)
			}
		// Floor
		case 2:
			return a.considerTreatyUsingUtility(t)
		// AvailableFood
		case 3:
			switch a.currBehaviour.String() {
			case "Altruist":
				return true
			case "Collectivist":
				if t.ConditionOp() == messages.LE || t.ConditionOp() == messages.LT {
					return a.considerTreatyUsingUtility(t)
				} else {
					if a.convertToTakeFoodAmount(float64(t.ConditionValue()), t.Request(), t.RequestValue()) <= 2 { //double-check if 2 is sufficient to go from critical to WeakLevel
						return a.considerTreatyUsingUtility(t)
					} else {
						return true
					}
				}
			case "Selfish":
				if t.ConditionOp() == messages.LE || t.ConditionOp() == messages.LT {
					return a.considerTreatyUsingUtility(t)
				} else {
					if a.convertToTakeFoodAmount(float64(t.ConditionValue()), t.Request(), t.RequestValue()) <= 60 { // change to a.HealthInfo().MaxFoodIntake
						return a.considerTreatyUsingUtility(t)
					} else {
						return true
					}
				}

			case "Narcissist":
				return a.considerTreatyUsingUtility(t)
			default:
				return a.considerTreatyUsingUtility(t)
			}

		default:
			return a.considerTreatyUsingUtility(t)
		}
	} else {
		return false
	}
}

// Decides if to accept or reject a treaty using utility
func (a *CustomAgent6) considerTreatyUsingUtility(t *messages.Treaty) bool {

	// 1. Estimate the food intake of the proposed treaty

	// Calculate how much food has been available on average on this floor
	sum := food.FoodType(0)
	for _, food := range a.shortTermMemory {
		sum += food
	}
	averageFoodAvailable := float64(sum) / math.Max(float64(len(a.shortTermMemory)), 1.0)

	// convert different treaty "Request types" to a "food intake"
	estimatedTakeFood := a.convertToTakeFoodAmount(averageFoodAvailable, t.Request(), t.RequestValue())
	if estimatedTakeFood == -1 {
		return false
	}

	// check the exact request condition
	if t.RequestOp() == messages.GE || t.RequestOp() == messages.GT /*t.ConditionOp() == messages.EQ || */ {
		// The treaty is of the form "Take X (or less) food"

		// 2. Calculate the agent's utility given different outcomes (accept or reject treaty)

		// a. estimated expected utility when accepting the treaty
		treatyTrustFactor := 1.0
		treatyUtility := treatyTrustFactor * Utility(float64(estimatedTakeFood), float64(a.desiredFoodIntake()), a.currBehaviour.String())

		// b. estimated utility when rejecting the treaty
		currentShortTermUtility := a.evaluateUtility(a.shortTermMemory) // estimated utility on the current floor
		currentLongTermUtility := a.evaluateUtility(a.longTermMemory)   // estimated utility over the entire time in the tower

		// benefit of signing the treaty (> 0 == beneficial)
		shortTermBenefit := treatyUtility - currentShortTermUtility
		longTermBenefit := treatyUtility - currentLongTermUtility

		estimatedPeriod := int(a.reassignPeriodGuess)
		estimatedTimeLeft := estimatedPeriod - len(a.shortTermMemory)
		benefit := 0.0

		// only think in the short term if
		// - the duration is shorter than the time left in estimated reassignment period OR
		// - we only have short term experience OR
		// - HP is critical (survival instincts take over)
		if t.Duration() < estimatedTimeLeft ||
			len(a.shortTermMemory) == len(a.longTermMemory) ||
			a.HP() <= a.HealthInfo().HPCritical {

			benefit = shortTermBenefit

		} else {
			// → The longer the duration of a treaty, the more important is it's long term benefit
			// E.g.
			// (1) time left on this level = 4, duration = 5 --> shortTermFocus 80%
			// 		--> The treaty mostly matters in the short term
			// (2) time left on this level = 4, duration = 8 --> shortTermFocus 50%
			// 		--> The treaty counts as much on this level as for the time after
			// (3) time left on this level = 1, duration = 100 --> shortTermFocus 1%
			// 		--> The treaty mostly matters in the long term

			shortTermFocus := float64(estimatedTimeLeft) / float64(t.Duration())
			benefit = shortTermFocus*shortTermBenefit + (1.0-shortTermFocus)*longTermBenefit
		}

		// 3. If we benefit from the treaty, accept it
		return benefit > 0.0

	} else {
		// The treaty is of the form "Take X (or more) food"

		// !!! Accepting treaties of this form would potentially lead to situations where it is impossible to satisfy all treaties.
		// E.g.
		// If you agree to leave more than an absolute amount (eg. 50) (1)

		// You also agree to leave less than a relative amount (50%) (2)

		// ! At proposal, we estimate 50% of the food to be 60 -> we can take less than 60 and more than 50 to satisfy this.

		// However, we actually get only 80 (instead of 120), wherefore 50% = 40 and we can't satisfy both (1) and (2).

		// - collectivist --> accept treaty if estimatedTakeFood less than (hoping to get others to eat at least the critical level)
		// if a.currBehaviour.String() == "Collectivist" {
		// 	return estimatedTakeFood <= 2.0 // the amount of food we need to get others to eat
		// }
		// All other social motives will allways reject
		// - altruist wants to avoid eating anyhing to save more for others
		// - selfish / narcissist doesn't want others to eat more than they want
		return false
	}
}

func (a *CustomAgent6) conditionApplies(t *messages.Treaty) bool {

	switch t.Condition() {
	// Condition : HP
	case 1:
		switch t.ConditionOp() {
		// GT
		case 1:
			return a.HP() > t.ConditionValue()
		//GE
		case 2:
			return a.HP() >= t.ConditionValue()
		//EQ
		case 3:
			return a.HP() == t.ConditionValue()
		// LE
		case 4:
			return a.HP() <= t.ConditionValue()
		// LT
		case 5:
			return a.HP() < t.ConditionValue()
		default:
			return true
		}
	// Condition : floor
	case 2:
		switch t.ConditionOp() {
		// GT
		case 1:
			return a.Floor() > t.ConditionValue()
		//GE
		case 2:
			return a.Floor() >= t.ConditionValue()
		//EQ
		case 3:
			return a.Floor() == t.ConditionValue()
		// LE
		case 4:
			return a.Floor() <= t.ConditionValue()
		// LT
		case 5:
			return a.Floor() < t.ConditionValue()
		default:
			return true
		}
	// Condition : AvailableFood
	case 3:
		switch t.ConditionOp() {
		// GT
		case 1:
			return int(a.CurrPlatFood()) > t.ConditionValue()
		//GE
		case 2:
			return int(a.CurrPlatFood()) >= t.ConditionValue()
		//EQ
		case 3:
			return int(a.CurrPlatFood()) == t.ConditionValue()
		// LE
		case 4:
			return int(a.CurrPlatFood()) <= t.ConditionValue()
		// LT
		case 5:
			return int(a.CurrPlatFood()) < t.ConditionValue()
		default:
			return true
		}
	default:
		return true
	}
}

func (a *CustomAgent6) ConstructTreaty() messages.Treaty {
	// Health levels
	levels := levelsData{
		strongLevel:  a.HealthInfo().MaxHP * 3 / 5,
		healthyLevel: a.HealthInfo().MaxHP * 3 / 10,
		weakLevel:    a.HealthInfo().WeakLevel,
		critLevel:    0,
	}
	proposedTreaty := messages.NewTreaty(1, a.HealthInfo().MaxHP, 2, 1, 1, 1, int(2*a.reassignPeriodGuess), a.ID())
	// Altruist and Narcissist do not propose treaties
	switch a.currBehaviour.String() {
	case "Collectivist":
		switch rand.Intn(2) {
		case 0:
			// ConditionType, conditionValue, RequestType, requestValue, cop, rop, duration, proposerID
			proposedTreaty = messages.NewTreaty(1, a.HealthInfo().WeakLevel, 2, 1, 1, 1, int(2*a.reassignPeriodGuess), a.ID())
		case 1:
			// Different treaty
			proposedTreaty = messages.NewTreaty(1, a.HealthInfo().WeakLevel, 2, 1, 1, 1, int(4*a.reassignPeriodGuess), a.ID())
		}

	case "Selfish":
		switch rand.Intn(2) {
		case 0:
			// ConditionType, conditionValue, RequestType, requestValue, cop, rop, duration, proposerID
			proposedTreaty = messages.NewTreaty(1, levels.strongLevel, 2, 1, 1, 1, int(2*a.reassignPeriodGuess), a.ID())
		case 1:
			// Different treaty
			proposedTreaty = messages.NewTreaty(1, levels.strongLevel, 2, 1, 1, 1, int(4*a.reassignPeriodGuess), a.ID())
		}
	}

	return *proposedTreaty
}

func (a *CustomAgent6) ProposeTreaty(treaty messages.Treaty) {
	// propose treaty to the 10 floor around us
	for i := -5; i < 6; i++ {
		// do not propose the treaty to yourself
		if i != 0 {
			targetFloor := a.Floor() + i
			proposedTreaty := messages.NewProposalMessage(a.ID(), a.Floor(), targetFloor, treaty)
			a.SendMessage(proposedTreaty)
			a.Log("I proposed a treaty", infra.Fields{"ConditionValue": treaty.ConditionValue(), "RequestValue": treaty.RequestValue(), "My Floor": a.Floor(), "My Social Motive": a.currBehaviour})
		}
	}

	// targetFloor := a.Floor() + i
	// proposedTreaty := messages.NewProposalMessage(a.ID(), a.Floor(), targetFloor, treaty)
	// a.SendMessage(proposedTreaty)
	// a.Log("I proposed a treaty", infra.Fields{"ConditionValue": treaty.ConditionValue(), "RequestValue": treaty.RequestValue(), "My Floor": a.Floor(), "My Social Motive": a.currBehaviour})

	a.proposedTreaties[treaty.ID()] = treaty
}