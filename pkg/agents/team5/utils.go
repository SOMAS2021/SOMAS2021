package team5

import (
	"math"
	"math/rand"

	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

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

func (a *CustomAgent5) checkForLeader() {
	// Random number between 3 and 12 generated, then the agent floor and selfishness are deducted from this
	diceRoll := rand.Intn(10) + 3 - a.selfishness - a.Floor()
	if diceRoll >= a.leadership {
		a.Log("An agent has become a leader", infra.Fields{"dice roll": diceRoll, "leadership": a.leadership, "selfishness": a.selfishness, "floor": a.Floor()})
		//TODO: Send treaties here about eating less food
		hpLevel := a.currentAimHP - ((a.currentAimHP-a.HealthInfo().WeakLevel)/10)*(diceRoll-a.leadership)
		a.currentProposal = messages.NewTreaty(messages.HP, hpLevel, messages.LeavePercentFood, 100, messages.GE, messages.EQ, 5, a.ID())
		a.treatySendCounter = 1
		a.Log("Agent is sending a treaty proposal", infra.Fields{"Proposed Max Hp": hpLevel})
		a.currentProposal.SignTreaty()
		a.AddTreaty(*a.currentProposal)
	}
}
