package team6

import (
	"log"
  "math/rand"
	"github.com/SOMAS2021/SOMAS2021/pkg/agents"
)

type behaviour int

const (
	altruist behaviour = iota
	collectivist
	selfish
	narcissist
)

type team6Config struct {
  currBehaviour behaviour
  stubborn float64
}

type CustomAgent6 struct {
	*agents.Base
  config team6Config
}

func chooseInitialBehaviour() (behaviour){
  return behaviour(rand.Intn(4))
}

func New(baseAgent *agents.Base) (agents.Agent, error) {
	return &CustomAgent6{
		Base:     baseAgent,
    config: team6Config{
      currBehaviour: chooseInitialBehaviour(),
  		stubborn:			0.5,
    },

	}, nil
}

func (a *CustomAgent6) Run() {
	log.Printf("Custom agent team 6 has floor: %d", a.Floor())
	log.Printf("Custom agent team 6 has behaviour: %d", a.config.currBehaviour)
}
