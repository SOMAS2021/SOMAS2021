package team6

import (
	"math"
	"math/rand"

	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
)

// We only accept LeaveAmountFood and LeavePercentFood treaties that say leave ≥,> of food available.
// Hence there will never be a invalid treaty.

// The utility function with
// x - food input
// z - desired food (maximum of the function)
func calculateUtility(x, z float64, socialMotive string) float64 {

	params := newUtilityParams(socialMotive)
	// Calculates the function scaling parameter a

	var result float64

	if socialMotive == "Altruist" || socialMotive == "Narcissist" {
		// Don't scale depending on desired food
		result = params.g*math.Pow(x, 1/params.r) - params.c*x
	} else {
		a := (1 / z) * math.Pow((params.c*params.r)/params.g, params.r/(1-params.r))
		result = params.g*math.Pow(a*x, 1/params.r) - params.c*a*x
	}

	if result < 0.0 {
		return 1.2 * result
	}
	return result
}

// Evaluates our agents current utility based on the current desired food
func (a *CustomAgent6) evaluateUtility(mem memory) float64 {
	sum := food.FoodType(0)
	for _, foodAvailable := range mem {
		min := math.Min(float64(foodAvailable), float64(a.desiredFoodIntake()))
		utilityValue := calculateUtility(min, float64(a.desiredFoodIntake()), a.currBehaviour.string())
		sum += food.FoodType(utilityValue)
	}

	return float64(sum) / math.Max(float64(len(mem)), 1.0)
}

// Converts different messages "Request types" to an equivalent "food intake" value
func (a *CustomAgent6) convertToTakeFoodAmount(foodAvailable float64, requestType messages.RequestType, requestValue int) food.FoodType {

	takeFood := 0.0
	switch requestType {
	case messages.LeaveAmountFood:
		takeFood = foodAvailable - float64(requestValue)
	case messages.LeavePercentFood:
		takeFood = foodAvailable * (1.0 - float64(requestValue))
	case messages.Inform:
		takeFood = -1
	default:
		takeFood = -1
	}

	return food.FoodType(takeFood)
}

//  Decides if to accept or reject a treaty
// The functions proceeds in 2 steps:
// 1) It checks which condition applies to the treaty. If the condition is an "obvious" case, we accept/reject it directly
// 2) If the case is not obvious, we use utility functions to rate the treaty (calling the "considerTreatyUsingUtility()")
func (a *CustomAgent6) considerTreaty(t *messages.Treaty) bool {

	// HP levels based on the maximum HP value
	levels := levelsData{
		strongLevel:  a.HealthInfo().MaxHP * 3 / 5,
		healthyLevel: a.HealthInfo().MaxHP * 3 / 10,
		weakLevel:    a.HealthInfo().WeakLevel,
		critLevel:    0,
	}

	// readable constants
	// no need to convert LeaveAmountFood to TakeAmountFood, as we only propose TakeAmountFood treaties
	//convertedFood := a.convertToTakeFoodAmount(float64(t.ConditionValue()), t.Request(), t.RequestValue())
	requestAmountFood := food.FoodType(t.RequestValue())
	conditionCheck := t.ConditionOp() == messages.LE || t.ConditionOp() == messages.LT
	consultUtility := a.considerTreatyUsingUtility(t)

	// We only consider meaningful LeaveAmountFood and LeavePercentFood treaties, i.e., the treaties with RequestOp LE or LT for TakeAmountFood (GE or GT for LeaveAmountFood)
	if t.RequestOp() == messages.LE || t.RequestOp() == messages.LT {
		switch t.Condition() {
		// HP
		case messages.HP:
			return a.considerHPTreaty(conditionCheck, consultUtility, levels, t.ConditionValue())
		// Floor
		case messages.Floor:
			return consultUtility
		// AvailableFood
		case messages.AvailableFood:
			return a.considerFoodTreaty(conditionCheck, consultUtility, requestAmountFood)
		default:
			return consultUtility

		}
	}
	return false
}

// food treaties
func (a *CustomAgent6) considerFoodTreaty(conditionCheck bool, consultUtility bool, requestAmountFood food.FoodType) bool {
	switch a.currBehaviour.string() {
	case "Altruist":
		return true
	case "Collectivist":
		// if conditionCheck || requestAmountFood <= 2 {
		// 	return consultUtility
		// }
		return true
	case "Selfish":
		if conditionCheck || requestAmountFood <= 60 {
			return consultUtility
		}
		return true

	default:
		return consultUtility
	}
}

// HP treaties
func (a *CustomAgent6) considerHPTreaty(conditionCheck bool, consultUtility bool, levels levelsData, conditionValue int) bool {
	switch a.currBehaviour.string() {
	case "Altruist":
		return true
	case "Collectivist":
		if conditionCheck || conditionValue < a.HealthInfo().WeakLevel {
			return consultUtility
		}
		return true

	case "Selfish":
		if conditionCheck || conditionValue < levels.strongLevel {
			return consultUtility
		}
		return true
	default:
		return consultUtility
	}
}

// Decides if to accept or reject a treaty using utility. Used in "considerUtility" function above
func (a *CustomAgent6) considerTreatyUsingUtility(t *messages.Treaty) bool {

	// Log

	// 1. Estimate the food intake of the proposed treaty

	// Calculates how much food has been available on average on this floor
	sum := food.FoodType(0)
	for _, food := range a.shortTermMemory {
		sum += food
	}
	averageFoodAvailable := float64(sum) / math.Max(float64(len(a.shortTermMemory)), 1.0)

	// No need to converts different treaty "Request types" to a "food intake" for TakeAmountFood treaties
	//estimatedTakeFood := a.convertToTakeFoodAmount(averageFoodAvailable, t.Request(), t.RequestValue())
	estimatedTakeFood := averageFoodAvailable
	if estimatedTakeFood == -1 {
		return false
	}

	// Checks the exact request condition. Only consider meaningful treaties: LE, LT for TakeAmountFood (GE, GT for LeaveAmountFood)
	if t.RequestOp() == messages.LE || t.RequestOp() == messages.LT /*t.ConditionOp() == messages.EQ || */ {
		// The treaty is of the form "Take X (or less) food"

		// 2. Calculate the agent's utility given different outcomes (accept or reject treaty)

		// a. estimated expected utility when accepting the treaty
		treatyTrustFactor := 1.0
		treatyUtility := treatyTrustFactor * calculateUtility(float64(estimatedTakeFood), float64(a.desiredFoodIntake()), a.currBehaviour.string())

		// b. estimated utility when rejecting the treaty
		currentShortTermUtility := a.evaluateUtility(a.shortTermMemory) // estimated utility on the current floor

		// benefit of signing the treaty (> 0 == beneficial)
		shortTermBenefit := treatyUtility - currentShortTermUtility

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

			currentLongTermUtility := a.evaluateUtility(a.longTermMemory) // estimated utility over the entire time in the tower
			longTermBenefit := treatyUtility - currentLongTermUtility

			shortTermFocus := float64(estimatedTimeLeft) / float64(t.Duration())
			benefit = shortTermFocus*shortTermBenefit + (1.0-shortTermFocus)*longTermBenefit
		}

		// 3. If we benefit from the treaty, accept it
		return benefit > 0.0

	}
	// We do not consider other type of treaties than TakeAmountFood treaties
	return false
}

func checkCondition(value int, conditionValue int, conditionOp messages.Op) bool {
	switch conditionOp {
	case messages.GT:
		return value > conditionValue
	case messages.GE:
		return value >= conditionValue
	case messages.EQ:
		return value == conditionValue
	case messages.LE:
		return value <= conditionValue
	case messages.LT:
		return value < conditionValue
	default:
		return false
	}
}

// Returns true if a given treaty applies to an agent by comparing the treaty condition to the agent's state
func (a *CustomAgent6) conditionApplies(t *messages.Treaty) bool {
	switch t.Condition() {
	case messages.HP:
		return checkCondition(a.HP(), t.ConditionValue(), t.ConditionOp())
	case messages.Floor:
		return checkCondition(a.Floor(), t.ConditionValue(), t.ConditionOp())
	case messages.AvailableFood:
		return checkCondition(int(a.CurrPlatFood()), t.ConditionValue(), t.ConditionOp())
	default:
		return false
	}
}

// Returns treaties to propose to other agents, based on the current behaviour and randomness
func (a *CustomAgent6) constructTreaty() messages.Treaty {
	// Health levels
	levels := levelsData{
		strongLevel:  a.HealthInfo().MaxHP * 3 / 5,
		healthyLevel: a.HealthInfo().MaxHP * 3 / 10,
		weakLevel:    a.HealthInfo().WeakLevel,
		critLevel:    0,
	}

	// var proposedTreaty := messages.NewTreaty(1, a.HealthInfo().MaxHP, 2, 1, 1, 1, int(2*a.reassignPeriodGuess), a.ID())
	// Altruist and Narcissist do not propose treaties
	switch a.currBehaviour.string() {
	case "Collectivist":
		switch rand.Intn(2) {
		case 0:
			// ConditionType, conditionValue, RequestType, requestValue, cop, rop, duration, proposerID
			// proposedTreaty = messages.NewTreaty(messages.HP, a.HealthInfo().WeakLevel, messages.LeavePercentFood, 1, messages.GT, messages.GE, int(2*a.reassignPeriodGuess), a.ID())
			return *messages.NewTreaty(messages.HP, a.HealthInfo().WeakLevel, messages.TakeAmountFood, 0, messages.GE, messages.LE, int(2*a.reassignPeriodGuess), a.ID())
		default:
			// Different treaty
			// proposedTreaty = messages.NewTreaty(messages.HP, a.HealthInfo().WeakLevel, messages.LeavePercentFood, 1, messages.GT, messages.GE, int(4*a.reassignPeriodGuess), a.ID())
			return *messages.NewTreaty(messages.HP, a.HealthInfo().WeakLevel, messages.TakeAmountFood, 0, messages.GE, messages.LE, int(2*a.reassignPeriodGuess), a.ID())
			//return *messages.NewTreaty(messages.HP, a.HealthInfo().WeakLevel, messages.TakeAmountFood, 2, messages.LT, messages.LE, int(2*a.reassignPeriodGuess), a.ID())
		}

	case "Selfish":
		switch rand.Intn(2) {
		case 0:
			// ConditionType, conditionValue, RequestType, requestValue, cop, rop, duration, proposerID
			// proposedTreaty = messages.NewTreaty(messages.HP, levels.strongLevel, messages.LeavePercentFood, 1, messages.GT, messages.GE, int(2*a.reassignPeriodGuess), a.ID())
			return *messages.NewTreaty(messages.HP, levels.strongLevel, messages.TakeAmountFood, 0, messages.GT, messages.LE, int(2*a.reassignPeriodGuess), a.ID())
		default:
			// Different treaty
			// proposedTreaty = messages.NewTreaty(messages.HP, levels.strongLevel, messages.LeavePercentFood, 1, messages.GT, messages.GE, int(4*a.reassignPeriodGuess), a.ID())
			return *messages.NewTreaty(messages.HP, levels.strongLevel, messages.TakeAmountFood, 0, messages.GT, messages.LE, int(4*a.reassignPeriodGuess), a.ID())
		}
	default:
		return *messages.NewTreaty(messages.HP, a.HealthInfo().MaxHP, messages.TakeAmountFood, 0, messages.GT, messages.LE, int(4*a.reassignPeriodGuess), a.ID())
	}
}

// Sends a message to the adjacent floors containing a treaty
func (a *CustomAgent6) proposeTreaty(treaty messages.Treaty) {
	// Proposes treaty to the 10 floor around us
	for i := -5; i < 6; i++ {
		// Propose treaty only if the trust in that "direction" is >0
		if (i > 0 && a.trustTeams[a.neighbours.below] > 0) || (i < 0 && a.trustTeams[a.neighbours.above] > 0) {
			targetFloor := a.Floor() + i
			proposedTreaty := messages.NewProposalMessage(a.ID(), a.Floor(), targetFloor, treaty)
			a.SendMessage(proposedTreaty)
			a.Log("I proposed a treaty", infra.Fields{"ConditionValue": treaty.ConditionValue(), "RequestValue": treaty.RequestValue(), "My Floor": a.Floor(), "My Social Motive": a.currBehaviour})
		}
	}
	a.proposedTreaties[treaty.ID()] = treaty
}
