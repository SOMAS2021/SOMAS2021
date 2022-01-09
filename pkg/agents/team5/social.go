package team5

import (
	"math"

	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
	"github.com/google/uuid"
)

func (a *CustomAgent5) newMemory(id uuid.UUID) {
	a.socialMemory[id] = Memory{
		foodTaken:         100,
		agentHP:           a.HealthInfo().MaxHP,
		intentionFood:     100,
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
	mem.favour = a.restrictToRange(0, 10, mem.favour+change)
	// a.Log("New Social Favour", infra.Fields{"new value": mem.favour})
	a.socialMemory[id] = mem
}

func (a *CustomAgent5) updateFavour() {
	for id, mem := range a.socialMemory {
		//a.Log("Days since last seen", infra.Fields{"search for": mem.daysSinceLastSeen})
		if mem.daysSinceLastSeen < 1 {
			powHP := 1.7
			powFood := 1.3
			maxHP := float64(a.HealthInfo().MaxHP)
			H_otherScoreComponent := -1 * math.Pow(float64(mem.agentHP), powHP) * math.Pow(float64(mem.foodTaken), powFood) / math.Pow(maxHP, 3)
			//a.Log("H_other_Component", infra.Fields{"h_other": H_otherScoreComponent})
			H_ourScoreComponent := math.Pow(float64(a.attemptFood), powFood) * math.Pow(float64(a.HP()), powHP) / math.Pow(maxHP, 3)
			//a.Log("H_Self_Component", infra.Fields{"h_self": H_ourScoreComponent})
			judgement := 100 * (H_otherScoreComponent + H_ourScoreComponent)
			// a.Log("I have judged an agent", infra.Fields{"judgement": judgement})
			if judgement > 0.075 {
				a.addToSocialFavour(id, 1)
			}
			if judgement < -2 {
				a.addToSocialFavour(id, int(math.Max(float64(judgement)/2, -3)))
			}
		}
		if mem.daysSinceLastSeen == 6 {
			a.resetSocialKnowledge(id)
		}
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
	mem.foodTaken = 100
	mem.agentHP = a.HealthInfo().MaxHP
	mem.intentionFood = 100
	a.socialMemory[id] = mem
}

func (a *CustomAgent5) updateSocialMemory(senderID uuid.UUID, senderFloor int) {
	if !a.memoryIdExists(senderID) {
		a.newMemory(senderID)
	}
	a.resetDaysSinceLastSeen(senderID)
	a.surroundingAgents[senderFloor-a.Floor()] = senderID
}
