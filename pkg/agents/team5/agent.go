package team5

import (
	"math"
	"math/rand"

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
	attemptToEat      bool
	leadership        int
	lastFood          int
	// Social network of other agents
	socialMemory      map[uuid.UUID]Memory
	surroundingAgents map[int]uuid.UUID
}

func New(baseAgent *infra.Base) (infra.Agent, error) {
	return &CustomAgent5{
		Base:              baseAgent,
		selfishness:       10,                           // of 0 to 10, with 10 being completely selfish, 0 being completely selfless
		lastMeal:          0,                            // Stores value of the last amount of food taken
		daysSinceLastMeal: 0,                            // Count of how many days since last eating
		hpAfterEating:     baseAgent.HealthInfo().MaxHP, //Stores HP value after eating in a day
		currentAimHP:      baseAgent.HealthInfo().MaxHP, //Stores aim HP for a given day
		attemptFood:       0,                            //Stores food agent will attempt to eat in a
		satisfaction:      -10,                          // Scale of -50 to 50
		rememberAge:       -1,                           // To check if a day has passed by our age increasing
		rememberFloor:     0,                            //Store the floor we are on so we can see if we have been reshuffled
		messagingCounter:  0,                            //Counter so that various messages are sent throughout the day
		lastFood:          0,                            // How much food arrived at the platform on the previous day assuming that the agent is still on the same floor
		attemptToEat:      true,                         // To check if we have already attempted to eat in a day. Needed because HasEaten() does not update if there is no food on the platform
		leadership:        rand.Intn(10),                //Initilise a random leadership value for each agent, used to determine whether they try to cause change in the tower. 0 is more likely to become a leader
		socialMemory:      make(map[uuid.UUID]Memory),   // Memory of other agents, key is agent id
		surroundingAgents: make(map[int]uuid.UUID),      // Map agent IDs of surrounding floors relative to current floor
	}, nil
}

func (a *CustomAgent5) updateAimHP() {
	a.currentAimHP = a.HealthInfo().MaxHP - ((10 - a.selfishness) * ((a.HealthInfo().MaxHP - a.HealthInfo().WeakLevel) / 10))
}

func (a *CustomAgent5) updateSelfishness() {
	a.selfishness = 10 - a.calculateAverageFavour()
}

func (a *CustomAgent5) updateSatisfaction() {
	tmp := int(a.CurrPlatFood())
	if a.lastFood*120.0/100.0 <= tmp {
		a.satisfaction += 2
	} else if a.lastFood*110.0/100.0 <= tmp {
		a.satisfaction += 1
	} else if a.lastFood*70.0/100.0 >= tmp {
		a.satisfaction -= 2
	} else if a.lastFood*80.0/100.0 >= tmp {
		a.satisfaction -= 1
	}

	if a.satisfaction > 50 {
		a.satisfaction = 50
	}

	if a.satisfaction < -50 {
		a.satisfaction = -50
	}
}

func (a *CustomAgent5) checkForLeader() {
	diceRoll := rand.Intn(10) + 3 - a.selfishness - a.Floor() //Random number between 3 and 12 generated, then the agent floor and selfishness are deducted from this
	if diceRoll >= a.leadership {
		a.Log("An agent has become a leader", infra.Fields{"dice roll": diceRoll, "leadership": a.leadership, "selfishness": a.selfishness, "floor": a.Floor()})
		//TODO: Send treaties here about eating less food
	}
}

func (a *CustomAgent5) dayPassed() {
	a.updateFavour()
	a.updateSelfishness()
	a.updateAimHP()
	a.attemptFood = a.calculateAttemptFood()
	a.checkForLeader()

	// for id := range a.socialMemory {
	// 	a.Log("Memory at end of day", infra.Fields{"favour": a.socialMemory[id].favour, "agentHP": a.socialMemory[id].agentHP, "foodTaken": a.socialMemory[id].foodTaken, "intent": a.socialMemory[id].intentionFood})
	// }
	// a.Log("Selfishness at end of day", infra.Fields{"selfishness": a.selfishness})
	// a.Log("Aim HP at end of day", infra.Fields{"aim HP": a.currentAimHP})
	// a.Log("Surrounding Agent Knowledge at end of day", infra.Fields{"agent map": a.surroundingAgents})

	a.daysSinceLastMeal++
	a.incrementDaysSinceLastSeen()
	// Needs to be fixed: you can be reshuffled but end up on the same floor.
	if a.rememberFloor != a.Floor() {
		a.ResetSurroundingAgents()
		a.rememberFloor = a.Floor()
	}
	a.messagingCounter = 0
	a.rememberAge = a.Age()
	a.attemptToEat = true
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
	if a.CurrPlatFood() != -1 && a.attemptToEat {
		if a.Floor() == a.rememberFloor { // if the agent is still on the same floor it can update its satisfaction
			a.updateSatisfaction()
		}
		a.lastFood = int(a.CurrPlatFood())
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
		a.hpAfterEating = a.HP()

		a.attemptToEat = false
	}
}
