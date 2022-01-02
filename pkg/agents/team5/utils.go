package team5

import "github.com/google/uuid"

func PercentageHP(a *CustomAgent5) int {
	return int(float64(a.HP()) / float64(a.HealthInfo().MaxHP) * 100.0)
}

func (a *CustomAgent5) restrictToRange(lowerBound, upperBound, num int) int {
	if num < lowerBound {
		return lowerBound
	}
	if num > upperBound {
		return upperBound
	}
	return num
}

func (a *CustomAgent5) memoryIdExists(id uuid.UUID) bool {
	_, exists := a.socialMemory[id]
	return exists
}

func (a *CustomAgent5) ResetSurroundingAgents() {
	a.surroundingAgents = make(map[int]uuid.UUID)
}
