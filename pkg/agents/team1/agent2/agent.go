package agent2

import (
	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/world"

	log "github.com/sirupsen/logrus"
)

type CustomAgent2 struct {
	*infra.Base
	myNumber int
	// new params
}

func New(world world.World, agentType int, agentHP int, agentFloor int, id string) (infra.Agent, error) {
	baseAgent, err := infra.NewBaseAgent(world, agentType, agentHP, agentFloor, id)
	if err != nil {
		log.Fatal(err)
	}
	return &CustomAgent2{
		Base:     baseAgent,
		myNumber: 0,
	}, nil
}

func (a *CustomAgent2) Run() {
	a.Log("Custom agent reporting status without using fields")
	a.TakeFood(15)
}
