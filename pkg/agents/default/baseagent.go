package baseagent

import (
	"errors"
	"log"

	tower "github.com/SOMAS2021/SOMAS2021/pkg/infra/tower"
	"github.com/divan/goabm/abm"
)

type BaseAgent struct {
	HP    int
	Floor int
	Tower *tower.Tower
}

func New(abm *abm.ABM, floor, hp int) (*BaseAgent, error) {
	world := abm.World()
	if world == nil {
		return nil, errors.New("Agent needs a World defined to operate")
	}
	tower, ok := world.(*tower.Tower)
	if !ok {
		return nil, errors.New("Agent needs a Tower world to operate")
	}
	return &BaseAgent{
		Floor: floor,
		HP:    hp,
		Tower: tower,
	}, nil
}

func (a *BaseAgent) Run() {
	log.Printf("An agent cycle executed from base agent %d", a.Floor)
}
