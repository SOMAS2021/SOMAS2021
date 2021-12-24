package randomAgent

import (
	"math/rand"

	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
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

func (a *CustomAgentRandom) HandleAskHP(msg messages.AskHPMessage)                                {}
func (a *CustomAgentRandom) HandleAskFoodTaken(msg messages.AskFoodTakenMessage)                  {}
func (a *CustomAgentRandom) HandleAskIntendedFoodTaken(msg messages.AskIntendedFoodIntakeMessage) {}
func (a *CustomAgentRandom) HandleRequestLeaveFood(msg messages.RequestLeaveFoodMessage)          {}
func (a *CustomAgentRandom) HandleRequestTakeFood(msg messages.RequestTakeFoodMessage)            {}
func (a *CustomAgentRandom) HandleResponse(msg messages.BoolResponseMessage)                      {}
func (a *CustomAgentRandom) HandleStateFoodTaken(msg messages.StateFoodTakenMessage)              {}
func (a *CustomAgentRandom) HandleStateHP(msg messages.StateHPMessage)                            {}
func (a *CustomAgentRandom) HandleStateIntendedFoodTaken(msg messages.StateIntendedFoodIntakeMessage) {
}
