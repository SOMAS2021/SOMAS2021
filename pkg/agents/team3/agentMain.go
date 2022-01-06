package team3

import (
	"math/rand"

	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
	"github.com/google/uuid"
)

type team3Variables struct {
	//Stubborness defines the likelihood of reading a message
	stubbornness int
	//Morality defines the willingness to help others/how much you care
	morality int
	//Mood affects the decision making and how you take things
	mood int
}

type team3Knowledge struct {
	//We know the floors we have been in
	floors []int
	//We know the last HP
	lastHP int
	//Stores who we have met and how we feel about them
	friends map[uuid.UUID]float64
	//Stores food last eaten
	foodLastEaten food.FoodType
	//Stores food last eaten
	foodLastSeen food.FoodType
	//Stores moving average of food consumed
	foodMovingAvg float64
	//Stores how old (in days) the agent is
	agentAge int
	//Stores whether we already have a treaty with the person above
	treatyProposed messages.Treaty
	//Stores our estimation of resshufle
	reshuffleEst int
}

type team3Decisions struct {
	//Amount of food we decided to eat.
	foodToEat int
	//Amount of food we decided to leave in the platform.
	foodToLeave int
}

type CustomAgent3 struct {
	*infra.Base
	vars      team3Variables
	knowledge team3Knowledge
	decisions team3Decisions
	//and an array of tuples for friendships
}

func New(baseAgent *infra.Base) (infra.Agent, error) {
	return &CustomAgent3{
		Base: baseAgent,
		vars: team3Variables{
			stubbornness: rand.Intn(75),
			morality:     rand.Intn(100),
			mood:         rand.Intn(100),
		},
		knowledge: team3Knowledge{
			floors:         []int{},
			lastHP:         100,
			friends:        make(map[uuid.UUID]float64),
			foodLastEaten:  food.FoodType(0),
			foodLastSeen:   food.FoodType(-1),
			foodMovingAvg:  0, //a.foodReqCalc(50, 50)
			agentAge:       0,
			treatyProposed: *messages.NewTreaty(messages.HP, 0, messages.LeaveAmountFood, 0, messages.GT, messages.GT, 0, uuid.Nil),
			reshuffleEst:   -1,
		},
		decisions: team3Decisions{
			foodToEat:   -1,
			foodToLeave: -1,
		},
	}, nil
}

func (a *CustomAgent3) Run() {
	//Update agent variables at the beginning of day (when we age)
	if a.knowledge.agentAge < a.BaseAgent().Age() {
		changeNewDay(a)
	}

	if a.knowledge.agentAge == 1 {
		a.knowledge.foodMovingAvg = float64(a.foodReqCalc(50, 50))
	}

	//Update agent variables at the beginning of reshuffle (when floor has changed)
	if len(a.knowledge.floors) == 0 || a.knowledge.floors[len(a.knowledge.floors)-1] != a.Floor() {
		changeNewFloor(a)
	}
	a.Log("Agent 3 each run:", infra.Fields{"floor": a.Floor(), "hp": a.HP(), "mood": a.vars.mood, "morality": a.vars.morality, "stubbornness": a.vars.stubbornness})

	//check if platform is in our floor, if so, eat.
	foodTaken, err := a.TakeFood(food.FoodType((a.takeFoodCalculation())))
	if err != nil {
		switch err.(type) {
		case *infra.FloorError:
		case *infra.NegFoodError:
		case *infra.AlreadyEatenError:
		default:
		}
	} else {
		//Update variables right after eating
		a.knowledge.foodLastSeen = a.BaseAgent().CurrPlatFood()
		a.knowledge.lastHP = a.HP()
		a.knowledge.foodLastEaten = food.FoodType(foodTaken)
		a.decisions.foodToEat = -1
		a.decisions.foodToLeave = -1
	}

	a.message()

	//send Message

}
