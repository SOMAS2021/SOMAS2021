package agent1

import (
	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/health"
)

type CustomAgent1 struct {
	*infra.Base
}

func New(baseAgent *infra.Base) (infra.Agent, error) {
	return &CustomAgent1{
		Base: baseAgent,
	}, nil
}

func (a *CustomAgent1) Run() {
	a.Log("Reporting agent state", infra.Fields{"health": a.HP(), "floor": a.Floor()})

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
