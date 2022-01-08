package team5

import (
	"math"

	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/health"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func restrictToRange(lowerBound, upperBound, num int) int {
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

func (a *CustomAgent5) resetSurroundingAgents() {
	a.surroundingAgents = make(map[int]uuid.UUID)
}

func (a *CustomAgent5) calculateAttemptFood() food.FoodType {
	if a.HP() < a.HealthInfo().WeakLevel {
		// TODO: UPDATE THIS VALUE TO A PARAMETER
		return food.FoodType(3)
	}
	foodAttempt := health.FoodRequired(a.HP(), a.currentAimHP, a.HealthInfo())
	return food.FoodType(math.Min(a.HealthInfo().Tau*3, float64(foodAttempt)))
}

func statementsIntersect(op1 messages.Op, value1 int, op2 messages.Op, value2 int) bool {
	switch op1 {
	case messages.EQ:
		switch op2 {
		case messages.EQ:
			return value1 == value2
		case messages.LT:
			return value1 < value2
		case messages.LE:
			return value1 <= value2
		case messages.GT:
			return value1 > value2
		case messages.GE:
			return value1 >= value2
		default:
			log.Error("Simulation - team5/agent.go: \t Reached unreachable code in statementsIntersect", op1, value1, op2, value2)
			return true
		}

	case messages.LT:
		switch op2 {
		case messages.EQ, messages.GE:
			return value1 > value2
		case messages.LT, messages.LE:
			return true
		case messages.GT:
			return value1-1 > value2
		default:
			log.Error("Simulation - team5/agent.go: \t Reached unreachable code in statementsIntersect", op1, value1, op2, value2)
			return true
		}

	case messages.LE:
		switch op2 {
		case messages.EQ, messages.GE:
			return value1 >= value2
		case messages.LT, messages.LE:
			return true
		case messages.GT:
			return value1 > value2
		default:
			log.Error("Simulation - team5/agent.go: \t Reached unreachable code in statementsIntersect", op1, value1, op2, value2)
			return true
		}

	case messages.GT:
		switch op2 {
		case messages.EQ, messages.LE:
			return value1 < value2
		case messages.LT:
			return value1 < value2-1
		case messages.GT, messages.GE:
			return true
		default:
			log.Error("Simulation - team5/agent.go: \t Reached unreachable code in statementsIntersect", op1, value1, op2, value2)
			return true
		}

	case messages.GE:
		switch op2 {
		case messages.EQ, messages.LE:
			return value1 <= value2
		case messages.LT:
			return value1 < value2
		case messages.GT, messages.GE:
			return true
		default:
			log.Error("Simulation - team5/agent.go: \t Reached unreachable code in statementsIntersect", op1, value1, op2, value2)
			return true
		}

	default:
		// It should be impossible to get here
		log.Error("Simulation - team5/agent.go: \t Reached unreachable code in statementsIntersect", op1, value1, op2, value2)
		return true
	}
}

func (a *CustomAgent5) overrideCalculation(treaty messages.Treaty) {
	if treaty.Request() == messages.LeaveAmountFood {
		if treaty.RequestOp() == messages.GT && a.CurrPlatFood()-a.attemptFood < food.FoodType(treaty.RequestValue()) {
			a.attemptFood = food.FoodType(math.Max(0, float64(a.CurrPlatFood()-food.FoodType(treaty.RequestValue())-1)))
		}
		if treaty.RequestOp() == messages.GE && a.CurrPlatFood()-a.attemptFood <= food.FoodType(treaty.RequestValue()) {
			a.attemptFood = food.FoodType(math.Max(0, float64(a.CurrPlatFood()-food.FoodType(treaty.RequestValue()))))
		}
	}
	if treaty.Request() == messages.LeavePercentFood {
		if treaty.RequestOp() == messages.EQ {
			a.attemptFood = a.CurrPlatFood() - food.FoodType(float64(a.CurrPlatFood())*float64(treaty.RequestValue())/100)
		}
		if treaty.RequestOp() == messages.GT && a.CurrPlatFood()-a.attemptFood < food.FoodType(float64(a.CurrPlatFood())*float64(treaty.RequestValue()/100)) {
			a.attemptFood = a.CurrPlatFood() - food.FoodType(float64(a.CurrPlatFood())*float64(treaty.RequestValue())/100) - 1
		}
		if treaty.RequestOp() == messages.GE && a.CurrPlatFood()-a.attemptFood <= food.FoodType(float64(a.CurrPlatFood())*float64(treaty.RequestValue()/100)) {
			a.attemptFood = a.CurrPlatFood() - food.FoodType(float64(a.CurrPlatFood())*float64(treaty.RequestValue())/100)
		}
	}
}
