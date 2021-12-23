package agent2

import (
	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
)

type CustomAgent2 struct {
	*infra.Base
	myNumber int
	// new params
}

func New(baseAgent *infra.Base) (infra.Agent, error) {
	return &CustomAgent2{
		Base:     baseAgent,
		myNumber: 0,
	}, nil
}

func (a *CustomAgent2) Run() {
	a.Log("Custom agent reporting status without using fields")
	a.TakeFood(15)
}

func (a *CustomAgent2) HandleAskHP()                              {}
func (a *CustomAgent2) HandleAskFoodTaken()                       {}
func (a *CustomAgent2) HandleAskIntendedFoodTaken()               {}
func (a *CustomAgent2) HandleRequestLeaveFood(request float64)    {}
func (a *CustomAgent2) HandleRequestTakeFood(request float64)     {}
func (a *CustomAgent2) HandleResponse(response bool)              {}
func (a *CustomAgent2) HandleStateFoodTaken(food float64)         {}
func (a *CustomAgent2) HandleStateHP(hp int)                      {}
func (a *CustomAgent2) HandleStateIntendedFoodTaken(food float64) {}
