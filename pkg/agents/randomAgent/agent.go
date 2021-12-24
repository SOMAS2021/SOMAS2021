package randomAgent

import (
	"math/rand"

	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/agent"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/world"
	log "github.com/sirupsen/logrus"
)

type CustomAgentRandom struct {
	*infra.Base
	// new params
}

func New(world world.World, agentType int, agentHP int, agentFloor int, id string) (agent.Agent, error) {
	baseAgent, err := infra.NewBaseAgent(world, agentType, agentHP, agentFloor, id)
	if err != nil {
		log.Fatal(err)
	}
	return &CustomAgentRandom{
		Base: baseAgent,
	}, nil
}

func (a *CustomAgentRandom) Run() {
	a.Log("Random agent reporting status:", infra.Fields{"floor": a.Floor(), "hp": a.HP()})
	a.TakeFood(food.FoodType(rand.Intn(100)))
}
