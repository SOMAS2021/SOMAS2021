package selfishAgent

import (
	"math/rand"

	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/health"

)

type CustomAgentSelfish struct {
	*infra.Base
	// new params
}

func New(baseAgent *infra.Base) (infra.Agent, error) {
	return &CustomAgentSelfish{
		Base: baseAgent,
	}, nil
}

func (a *CustomAgentSelfish) Run() {
	a.Log("Selfish agent reporting status:", infra.Fields{"floor": a.Floor(), "hp": a.HP()})
	healthInfo := a.HealthInfo()
	maxHP := healthInfo.MaxHP
	foodAmt := health.FoodRequired(a.HP(), maxHP, healthInfo)

	_, err := a.TakeFood(foodAmt)
	if err != nil {
		switch err.(type) {
		case *infra.FloorError:
		case *infra.NegFoodError:
		case *infra.AlreadyEatenError:
		default:
		}
	}
}
