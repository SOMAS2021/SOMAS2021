package randomAgent

import (
	"math/rand"

	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/agent"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
)

type CustomAgentRandom struct {
	*infra.Base
	// new params
}

func New(baseAgent *infra.Base) (agent.Agent, error) {
	//create other parameters
	return &CustomAgentRandom{
		Base: baseAgent,
	}, nil
}

func (a *CustomAgentRandom) Run() {
	a.Log("Random agent reporting status:", infra.Fields{"floor": a.Floor(), "hp": a.HP()})
	a.TakeFood(food.FoodType(rand.Intn(100)))
}
