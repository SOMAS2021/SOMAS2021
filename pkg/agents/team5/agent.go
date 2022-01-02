package team5

import (
	"math"

	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/health"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type Memory struct {
	foodTaken         food.FoodType // Store last known value of foodTaken by agent
	agentHP           int           // Store last known value of HP of agent
	intentionFood     food.FoodType // Store the last known value of the amount of food intended to take
	favour            int           // e.g. generosity; scale of 0 to 10, with 0 being least favoured and 10 being most favoured
	daysSinceLastSeen int           // days since last interaction
}

type CustomAgent5 struct {
	*infra.Base
	selfishness       int
	lastMeal          food.FoodType
	daysSinceLastMeal int
	hpAfterEating     int
	currentAimHP      int
	attemptFood       food.FoodType
	satisfaction      int
	rememberAge       int
	rememberFloor     int
	messagingCounter  int
	// Social network of other agents
	socialMemory      map[uuid.UUID]Memory
	surroundingAgents map[int]uuid.UUID
}

func New(baseAgent *infra.Base) (infra.Agent, error) {
	return &CustomAgent5{
		Base:              baseAgent,
		selfishness:       10, // of 0 to 10, with 10 being completely selfish, 0 being completely selfless
		lastMeal:          0,  // Stores value of the last amount of food taken
		daysSinceLastMeal: 0,  // Count of how many days since last eating
		hpAfterEating:     baseAgent.HealthInfo().MaxHP,
		currentAimHP:      baseAgent.HealthInfo().MaxHP,
		attemptFood:       0,
		satisfaction:      0,                          // Scale of -3 to 3, with 3 being satisfied and unsatisfied
		rememberAge:       -1,                         // To check if a day has passed by our age increasing
		socialMemory:      make(map[uuid.UUID]Memory), // Memory of other agents, key is agent id
		surroundingAgents: make(map[int]uuid.UUID),    // Map agent IDs of surrounding floors relative to current floor
	}, nil
}

func (a *CustomAgent5) updateAimHP() {
	a.currentAimHP = a.HealthInfo().MaxHP - ((10 - a.selfishness) * ((a.HealthInfo().MaxHP - a.HealthInfo().WeakLevel) / 10))
}

func (a *CustomAgent5) updateSelfishness() {
	a.selfishness = 10 - a.calculateAverageFavour()
}

func (a *CustomAgent5) updateSatisfaction() {
	if PercentageHP(a) >= 100 {
		a.satisfaction = 3
	}
	if a.lastMeal == 0 && a.satisfaction > -3 {
		a.satisfaction--
	}
	if PercentageHP(a) < 25 && a.satisfaction > -3 {
		a.satisfaction--
	}
	if PercentageHP(a) > 75 && a.satisfaction < 3 {
		a.satisfaction++
	}
}

func (a *CustomAgent5) dayPassed() {
	a.updateFavour()
	a.updateSelfishness()
	a.updateAimHP()
	a.attemptFood = a.calculateAttemptFood()

	// for id := range a.socialMemory {
	// 	a.Log("Memory at end of day", infra.Fields{"favour": a.socialMemory[id].favour, "agentHP": a.socialMemory[id].agentHP, "foodTaken": a.socialMemory[id].foodTaken, "intent": a.socialMemory[id].intentionFood})
	// }
	// a.Log("Selfishness at end of day", infra.Fields{"selfishness": a.selfishness})
	// a.Log("Aim HP at end of day", infra.Fields{"aim HP": a.currentAimHP})
	// a.Log("Surrounding Agent Knowledge at end of day", infra.Fields{"agent map": a.surroundingAgents})

	a.daysSinceLastMeal++
	a.incrementDaysSinceLastSeen()
	if a.rememberFloor != a.Floor() {
		a.ResetSurroundingAgents()
		a.rememberFloor = a.Floor()
	}
	a.messagingCounter = 0
	a.rememberAge = a.Age()
}

func (a *CustomAgent5) calculateAttemptFood() food.FoodType {
	if a.HP() < a.HealthInfo().WeakLevel {
		// TODO: UPDATE THIS VALUE TO A PARAMETER
		return food.FoodType(3)
	}
	foodAttempt := health.FoodRequired(a.HP(), a.currentAimHP, a.HealthInfo())
	return food.FoodType(math.Min(a.HealthInfo().Tau*3, float64(foodAttempt)))
}

func (a *CustomAgent5) Run() {
	a.Log("Reporting agent state of team 5 agent", infra.Fields{"health": a.HP(), "floor": a.Floor()})

	// Check if a day has passed
	if a.Age() > a.rememberAge {
		a.dayPassed()
	}

	a.getMessages()
	a.dailyMessages()

	// When platform reaches our floor and we haven't tried to eat, then try to eat
	if a.CurrPlatFood() != -1 && !a.HasEaten() {
		lastMeal, err := a.TakeFood(a.attemptFood)
		if err != nil {
			switch err.(type) {
			case *infra.FloorError:
			case *infra.NegFoodError:
				log.Error("Simulation - team5/agent.go: \t NegFoodError: did CalculateAttemptFood() return a negative?")
			case *infra.AlreadyEatenError:
				log.Error("Simulation - team5/agent.go: \t AlreadyEatenError occurred after checking for a.HasEaten()")
			default:
				log.Error("Simulation - team5/agent.go: \t Impossible error reached")
			}
		}
		a.lastMeal = lastMeal
		if a.lastMeal > 0 {
			a.daysSinceLastMeal = 0
		}
		a.updateSatisfaction()
		a.hpAfterEating = a.HP()
		//a.messagingCounter = 0
	}
}
