package baseagent

import (
	"errors"
	"log"

	tower "github.com/SOMAS2021/SOMAS2021/pkg/infra/tower"
	"github.com/divan/goabm/abm"
	"github.com/google/uuid"
)

type Agent interface {
	Run()
}

type BaseAgent struct {
	hp    int
	floor int
	id    string
	tower *tower.Tower
}

func NewBaseAgent(abm *abm.ABM, floor, hp int) (*BaseAgent, error) {
	world := abm.World()
	if world == nil {
		return nil, errors.New("Agent needs a World defined to operate")
	}
	tower, ok := world.(*tower.Tower)
	if !ok {
		return nil, errors.New("Agent needs a Tower world to operate")
	}
	return &BaseAgent{
		floor: floor,
		hp:    hp,
		tower: tower,
		id:    uuid.New().String(),
	}, nil
}

func (a *BaseAgent) Run() {
	log.Printf("An agent cycle executed from base agent %d", a.floor)
}

func (a *BaseAgent) GetHP() int {
	return a.hp
}

func (a *BaseAgent) GetFloor() int {
	return a.floor
}

func (a *BaseAgent) GetID() string {
	return a.id
}
