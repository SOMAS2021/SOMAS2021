package team5

import (
	"math/rand"

	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
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
	treatySendCounter int
	currentProposal   *messages.Treaty
	attemptToEat      bool
	leadership        int
	lastSeenFood      food.FoodType
	sensitivity       int
	avgPlatFood       food.FoodType
	daysSinceShuffle  int
	// Social network of other agents
	socialMemory      map[uuid.UUID]Memory
	surroundingAgents map[int]uuid.UUID
}

func New(baseAgent *infra.Base) (infra.Agent, error) {
	return &CustomAgent5{
		Base:              baseAgent,
		selfishness:       10,                           // Range from 0 to 10, with 10 being completely selfish, 0 being completely selfless
		lastMeal:          0,                            // Stores value of the last amount of food taken
		daysSinceLastMeal: 0,                            // Count of how many days since last eating
		hpAfterEating:     baseAgent.HealthInfo().MaxHP, // Stores HP value after eating in a day
		currentAimHP:      baseAgent.HealthInfo().MaxHP, // Stores aim HP for a given day
		attemptFood:       0,                            // Stores food agent will attempt to eat in a
		satisfaction:      -10,                          // Scale of -50 to 50
		rememberAge:       -1,                           // To check if a day has passed by our age increasing
		rememberFloor:     0,                            // Store the floor we are on so we can see if we have been reshuffled
		messagingCounter:  0,                            // Counter so that various messages are sent throughout the day
		treatySendCounter: 0,                            // Counter so that treaty messages can be sent
		attemptToEat:      true,                         // To check if we have already attempted to eat in a day. Needed because HasEaten() does not update if there is no food on the platform
		leadership:        rand.Intn(10),                // Initialise a random leadership value for each agent, used to determine whether they try to cause change in the tower. 0 is more likely to become a leader
		lastSeenFood:      0,                            // How much food arrived at the platform on the previous day assuming that the agent is still on the same floor
		sensitivity:       2 + rand.Intn(8),             // Set personality trait, basically determines amplification of change in favour for nearby agents
		avgPlatFood:       0,                            // Average food arriving on platform
		daysSinceShuffle:  0,                            // days since last shuffle
		socialMemory:      make(map[uuid.UUID]Memory),   // Memory of other agents, key is agent id
		surroundingAgents: make(map[int]uuid.UUID),      // Map agent IDs of surrounding floors relative to current floor
	}, nil
}

func (a *CustomAgent5) updateAimHP() {
	a.currentAimHP = a.HealthInfo().MaxHP - ((10 - a.selfishness) * ((a.HealthInfo().MaxHP - a.HealthInfo().WeakLevel) / 10))
}

func (a *CustomAgent5) updateSelfishness() {
	// Tit for tat strategy, agent will conform to the mean behaviour of their social network
	a.selfishness = 10 - a.calculateAverageFavour()
	// Make agent less selfish if going through tough times, lowers their expectations and makes them more sympathetic of others struggles
	a.selfishness = restrictToRange(0, 10, a.selfishness-a.daysSinceLastMeal)
}

func (a *CustomAgent5) updateSatisfaction() {
	seenFood := a.CurrPlatFood()
	newSatisfaction := a.satisfaction

	if seenFood >= a.lastSeenFood*12/10 {
		newSatisfaction += 2
	} else if seenFood >= a.lastSeenFood*11/10 {
		newSatisfaction += 1
	} else if seenFood <= a.lastSeenFood*8/10 {
		newSatisfaction -= 1
	} else if seenFood <= a.lastSeenFood*7/10 {
		newSatisfaction -= 2
	}

	a.satisfaction = restrictToRange(-50, 50, newSatisfaction)
}

func (a *CustomAgent5) updateAveragePlatFood() {
	a.avgPlatFood = ((a.avgPlatFood * food.FoodType(a.daysSinceShuffle)) + a.CurrPlatFood()) / food.FoodType(a.daysSinceShuffle+1)
	a.updateFavourAbove()
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
		a.resetSurroundingAgents()
		a.rememberFloor = a.Floor()
		a.daysSinceShuffle = 0
		a.avgPlatFood = 0
	}
	a.messagingCounter = 0
	a.rememberAge = a.Age()
	a.attemptToEat = true
}

func (a *CustomAgent5) Run() {
	a.Log("Reporting agent state of team 5 agent", infra.Fields{"health": a.HP(), "floor": a.Floor()})

	// Check if a day has passed
	if a.Age() > a.rememberAge {
		a.dayPassed()
	}

	a.getMessages()
	if a.treatySendCounter == 0 {
		a.dailyMessages()
	} else {
		a.treatyProposal()
	}

	// When platform reaches our floor and we haven't tried to eat, then try to eat
	if a.CurrPlatFood() != -1 && a.attemptToEat {
		a.treatyOverride()
		if a.Floor() == a.rememberFloor { // if the agent is still on the same floor it can update its satisfaction
			a.updateSatisfaction()
			a.updateAveragePlatFood()
			a.daysSinceShuffle++
		}
		a.lastSeenFood = a.CurrPlatFood()
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
