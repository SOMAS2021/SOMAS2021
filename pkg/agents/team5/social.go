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

func (a *CustomAgent5) updateFavourAbove() {
	if agentAbove, agentAboveKnown := a.surroundingAgents[1]; agentAboveKnown {

		H_above := float64(a.socialMemory[agentAbove].agentHP)
		F_above := float64(a.socialMemory[agentAbove].foodTaken)
		F_expected := float64(a.avgPlatFood)
		S := float64(a.sensitivity)
		F_take := float64(a.attemptFood)
		H_self := float64(PercentageHP(a))
		F_plat := float64(a.CurrPlatFood())

		F_expectedComponent := math.Exp(F_plat/((3.0*F_expected)+15.0)) * (F_plat - F_expected) / (50 * S)
		H_aboveScoreComponent := -2.0 * S * H_above * math.Pow(F_above, 2) / math.Pow(float64(a.HealthInfo().MaxHP), 3)
		H_ourScoreComponent := 2.0 * S * math.Pow(F_take, 1.3) * math.Pow(H_self, 1.7) / math.Pow(float64(a.HealthInfo().MaxHP), 3)

		C := int((math.Round(F_expectedComponent + H_aboveScoreComponent + H_ourScoreComponent)))
		a.addToSocialFavour(agentAbove, C)
		a.resetSocialKnowledge(agentAbove)
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
