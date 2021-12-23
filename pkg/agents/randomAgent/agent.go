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

func (a *CustomAgentRandom) HandleAskHP()                              {}
func (a *CustomAgentRandom) HandleAskFoodTaken()                       {}
func (a *CustomAgentRandom) HandleAskIntendedFoodTaken()               {}
func (a *CustomAgentRandom) HandleRequestLeaveFood(request float64)    {}
func (a *CustomAgentRandom) HandleRequestTakeFood(request float64)     {}
func (a *CustomAgentRandom) HandleResponse(response bool)              {}
func (a *CustomAgentRandom) HandleStateFoodTaken(food float64)         {}
func (a *CustomAgentRandom) HandleStateHP(hp int)                      {}
func (a *CustomAgentRandom) HandleStateIntendedFoodTaken(food float64) {}
