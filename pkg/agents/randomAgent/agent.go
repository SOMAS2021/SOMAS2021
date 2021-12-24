package randomAgent

import (
	"math/rand"

	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
)

type CustomAgentRandom struct {
	*infra.Base
	// new params
}

func New(baseAgent *infra.Base) (infra.Agent, error) {
	return &CustomAgentRandom{
		Base: baseAgent,
	}, nil
}

func (a *CustomAgentRandom) Run() {
	a.Log("Random agent reporting status:", infra.Fields{"floor": a.Floor(), "hp": a.HP()})
	a.TakeFood(food.FoodType(rand.Intn(100)))
}

func (a *CustomAgentRandom) HandleAskHP(msg infra.AskMessage)                    {}
func (a *CustomAgentRandom) HandleAskFoodTaken(msg infra.AskMessage)             {}
func (a *CustomAgentRandom) HandleAskIntendedFoodTaken(msg infra.AskMessage)     {}
func (a *CustomAgentRandom) HandleRequestLeaveFood(msg infra.RequestMessage)     {}
func (a *CustomAgentRandom) HandleRequestTakeFood(msg infra.RequestMessage)      {}
func (a *CustomAgentRandom) HandleResponse(msg infra.ResponseMessage)            {}
func (a *CustomAgentRandom) HandleStateFoodTaken(msg infra.StateMessage)         {}
func (a *CustomAgentRandom) HandleStateHP(msg infra.StateMessage)                {}
func (a *CustomAgentRandom) HandleStateIntendedFoodTaken(msg infra.StateMessage) {}
