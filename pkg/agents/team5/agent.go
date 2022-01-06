package team5

import (
	"math"
	"math/rand"

	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
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
	treatySendCounter int
	currentProposal   *messages.Treaty
	attemptToEat      bool
	leadership        int
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
		satisfaction:      0,                            // Scale of -3 to 3, with 3 being satisfied and unsatisfied
		rememberAge:       -1,                           // To check if a day has passed by our age increasing
		rememberFloor:     0,                            // Store the floor we are on so we can see if we have been reshuffled
		messagingCounter:  0,                            // Counter so that various messages are sent throughout the day
		treatySendCounter: 0,                            // Counter so that treaty messages can be sent
		attemptToEat:      true,                         // To check if we have already attempted to eat in a day. Needed because HasEaten() does not update if there is no food on the platform
		leadership:        rand.Intn(10),                // Initialise a random leadership value for each agent, used to determine whether they try to cause change in the tower. 0 is more likely to become a leader
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
	a.selfishness = a.restrictToRange(0, 10, a.selfishness-a.daysSinceLastMeal)
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

func (a *CustomAgent5) treatiesCanCoexist(t1 messages.Treaty, t2 messages.Treaty) bool {
	conditionsIntersect := statementsIntersect(t1.ConditionOp(), t1.ConditionValue(), t2.ConditionOp(), t2.ConditionValue())

	// If the treaties are based on the same type of condition, they can coexist
	// when the conditions are mutually exclusive (can never occurr simultaneously)
	if t1.Condition() == t2.Condition() && !conditionsIntersect {
		return true
	}

	// Otherwise, if the treaties are based on different conditions, or the conditions
	// can occurr simultaneously, treaties of different types should be rejected
	// because it is impossible to deduce if their requests can be fulfilled simultaneously
	if t1.Request() != t2.Request() {
		return false
	}

	// Otherwise, treaties can coexist only if the requests can be fulfilled simultaneously
	return statementsIntersect(t1.RequestOp(), t1.RequestValue(), t2.RequestOp(), t2.RequestValue())
}

func (a *CustomAgent5) treatyConflicts(treaty messages.Treaty) bool {
	for _, activeTreaty := range a.ActiveTreaties() {
		if !a.treatiesCanCoexist(treaty, activeTreaty) {
			return true
		}
	}
	return false
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

func (a *CustomAgent5) treatyOverride() {
	for _, treaty := range a.ActiveTreaties() {
		switch {
		case treaty.Condition() == messages.HP && treaty.ConditionOp() == messages.EQ && a.HP() == treaty.ConditionValue():
			fallthrough
		case treaty.Condition() == messages.HP && treaty.ConditionOp() == messages.GT && a.HP() > treaty.ConditionValue():
			fallthrough
		case treaty.Condition() == messages.HP && treaty.ConditionOp() == messages.GE && a.HP() >= treaty.ConditionValue():
			fallthrough
		case treaty.Condition() == messages.Floor && treaty.ConditionOp() == messages.EQ && a.Floor() == treaty.ConditionValue():
			fallthrough
		case treaty.Condition() == messages.Floor && treaty.ConditionOp() == messages.LT && a.Floor() < treaty.ConditionValue():
			fallthrough
		case treaty.Condition() == messages.Floor && treaty.ConditionOp() == messages.LE && a.Floor() <= treaty.ConditionValue():
			fallthrough
		case treaty.Condition() == messages.AvailableFood && treaty.ConditionOp() == messages.EQ && a.CurrPlatFood() == food.FoodType(treaty.ConditionValue()):
			fallthrough
		case treaty.Condition() == messages.AvailableFood && treaty.ConditionOp() == messages.GT && a.CurrPlatFood() > food.FoodType(treaty.ConditionValue()):
			fallthrough
		case treaty.Condition() == messages.AvailableFood && treaty.ConditionOp() == messages.GE && a.CurrPlatFood() >= food.FoodType(treaty.ConditionValue()):
			a.overrideCalculation(treaty)
		}
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
	if a.treatySendCounter == 0 {
		a.dailyMessages()
	} else {
		a.treatyProposal()
	}

	// When platform reaches our floor and we haven't tried to eat, then try to eat
	if a.CurrPlatFood() != -1 && a.attemptToEat {
		a.treatyOverride()
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
		a.updateSatisfaction()
		a.attemptToEat = false
	}
}
