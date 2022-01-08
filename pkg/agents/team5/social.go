package team5

import (
	"math"

	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
	"github.com/google/uuid"
)

func (a *CustomAgent5) newMemory(id uuid.UUID) {
	a.socialMemory[id] = Memory{
		foodTaken:         50,
		agentHP:           a.HealthInfo().MaxHP / 2,
		intentionFood:     50,
		favour:            0,
		daysSinceLastSeen: 0,
	}
}

func (a *CustomAgent5) updateFoodTakenMemory(id uuid.UUID, foodTaken food.FoodType) {
	if !a.memoryIdExists(id) {
		a.newMemory(id)
	}
	mem := a.socialMemory[id]
	mem.foodTaken = foodTaken
	a.socialMemory[id] = mem
}

func (a *CustomAgent5) updateAgentHPMemory(id uuid.UUID, agentHP int) {
	if !a.memoryIdExists(id) {
		a.newMemory(id)
	}
	mem := a.socialMemory[id]
	mem.agentHP = agentHP
	a.socialMemory[id] = mem
}

func (a *CustomAgent5) updateIntentionFoodMemory(id uuid.UUID, intentionFood food.FoodType) {
	if !a.memoryIdExists(id) {
		a.newMemory(id)
	}
	mem := a.socialMemory[id]
	mem.intentionFood = intentionFood
	a.socialMemory[id] = mem
}

func (a *CustomAgent5) addToSocialFavour(id uuid.UUID, change int) {
	if !a.memoryIdExists(id) {
		a.newMemory(id)
	}
	mem := a.socialMemory[id]
	mem.favour = restrictToRange(0, 10, mem.favour+change)
	a.socialMemory[id] = mem
}

func (a *CustomAgent5) updateFavour() {
	for id, mem := range a.socialMemory {
		//a.Log("Days since last seen", infra.Fields{"search for": mem.daysSinceLastSeen})
		if mem.daysSinceLastSeen < 1 {
			judgement := (a.hpAfterEating - mem.agentHP) + int(a.lastMeal-mem.foodTaken) //+ int(a.calculateAttemptFood()-mem.intentionFood)
			//a.Log("I have judged an agent", infra.Fields{"judgement": judgement})
			if judgement > 0 {
				a.addToSocialFavour(id, 1)
			}
			if judgement < 0 {
				a.addToSocialFavour(id, int(math.Max(float64(judgement)/20, -3)))
			}
		}
		if mem.daysSinceLastSeen == 6 {
			a.resetSocialKnowledge(id)
		}
	}
}

func (a *CustomAgent5) calculateFavourAbove(aboveID uuid.UUID) int {
	sens := a.sensitivity
	foodIntent := float64(a.attemptFood)
	hpSelf := float64(a.HP())
	// I would still add a check here to make sure it isn't -1.
	foodPlat := float64(a.CurrPlatFood())

	maxHP := float64(a.HealthInfo().MaxHP)

	hpAbove := float64(a.socialMemory[aboveID].agentHP)
	foodAbove := float64(a.socialMemory[aboveID].foodTaken)
	foodExp := float64(a.avgPlatFood)

	F_expectedComponent := math.Exp(foodPlat/(3.0*foodExp+15.0)) * (foodPlat - foodExp) / (50 * sens)
	H_aboveScoreComponent := -2 * sens * hpAbove * math.Pow(foodAbove, 2) / math.Pow(maxHP, 3)
	H_ourScoreComponent := 2 * sens * math.Pow(foodIntent, 1.3) * math.Pow(hpSelf, 1.7) / math.Pow(maxHP, 3)

	change := F_expectedComponent + H_aboveScoreComponent + H_ourScoreComponent
	return int(math.Round(change))
}

func (a *CustomAgent5) calculateExpectedFood() float64 {
	// average of food taken by agents around our agent, including our agent
	foodTotal := a.lastMeal
	for _, uuid := range a.surroundingAgents {
		foodTotal += a.socialMemory[uuid].foodTaken
	}
	foodExp := float64(foodTotal) / float64(len(a.surroundingAgents)+1)
	if foodExp == 0 {
		foodExp = float64(a.attemptFood) // to avoid favour shooting extremely high
	}
	return foodExp
}

func (a *CustomAgent5) calculateFavourBelow(belowID uuid.UUID) int {
	sens := float64(a.sensitivity)
	foodIntent := float64(a.attemptFood)
	hpSelf := float64(a.HP())

	maxHP := float64(a.HealthInfo().MaxHP)

	hpBelow := float64(a.socialMemory[belowID].agentHP)
	foodBelow := float64(a.socialMemory[belowID].foodTaken)
	foodAsked := float64(a.socialMemory[belowID].intentionFood)

	foodExp := a.calculateExpectedFood()

	a.updateIntentionFoodMemory(belowID, food.FoodType(foodExp))

	F_expectedComponent := math.Exp(foodAsked/(3*foodExp+1)) * (foodAsked - foodExp) / 10
	H_belowScoreComponent := -2.0 * sens * math.Pow(hpBelow, 1.7) * math.Pow(foodBelow, 1.3) / math.Pow(maxHP, 3)
	H_ourScoreComponent := 2.0 * sens * math.Pow(foodIntent, 1.3) * math.Pow(hpSelf, 1.7) / math.Pow(maxHP, 3)

	change := F_expectedComponent + H_belowScoreComponent + H_ourScoreComponent
	return int(math.Round(change))
}

func (a *CustomAgent5) updateSocialFavour(agentID uuid.UUID, change int) {
	a.addToSocialFavour(agentID, change)
	a.resetSocialKnowledge(agentID)
}

func (a *CustomAgent5) updateFavourForNeighbours() {
	if aboveID, agentAboveKnown := a.surroundingAgents[-1]; agentAboveKnown {
		change := a.calculateFavourAbove(aboveID)
		a.updateSocialFavour(aboveID, change)
	}
	if belowID, agentBelowKnown := a.surroundingAgents[1]; agentBelowKnown {
		change := a.calculateFavourBelow(belowID)
		a.updateSocialFavour(belowID, change)
	}
}

func (a *CustomAgent5) calculateAverageFavour() int {
	sum := 0
	count := 0
	for floor := range a.surroundingAgents {
		sum += a.socialMemory[a.surroundingAgents[floor]].favour
		count++
	}
	if count == 0 {
		return 10 - a.selfishness
	}
	return sum / count
}

func (a *CustomAgent5) incrementDaysSinceLastSeen() {
	for id := range a.socialMemory {
		mem := a.socialMemory[id]
		mem.daysSinceLastSeen++
		a.socialMemory[id] = mem
	}
}

func (a *CustomAgent5) resetDaysSinceLastSeen(id uuid.UUID) {
	if !a.memoryIdExists(id) {
		a.newMemory(id)
	}
	mem := a.socialMemory[id]
	mem.daysSinceLastSeen = 0
	a.socialMemory[id] = mem
}

func (a *CustomAgent5) resetSocialKnowledge(id uuid.UUID) {
	if !a.memoryIdExists(id) {
		a.newMemory(id)
	}
	mem := a.socialMemory[id]
	mem.foodTaken = 50
	mem.agentHP = a.HealthInfo().MaxHP / 2
	mem.intentionFood = 50
	a.socialMemory[id] = mem
}

func (a *CustomAgent5) updateSocialMemory(senderID uuid.UUID, senderFloor int) {
	if !a.memoryIdExists(senderID) {
		a.newMemory(senderID)
	}
	a.resetDaysSinceLastSeen(senderID)
	a.surroundingAgents[senderFloor-a.Floor()] = senderID
}
