package agents

import (
	"errors"
	"log"

	"github.com/SOMAS2021/SOMAS2021/pkg/infra/tower"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/abm"
)

type Agent interface {
	Run()
	IsDead() bool
}

type Base struct {
	id    string
	tower *tower.Tower
}

func NewBaseAgent(abm *abm.ABM, uuid string) (*Base, error) {
	world := abm.World()
	if world == nil {
		return nil, errors.New("Agent needs a World defined to operate")
	}
	tower, ok := world.(*tower.Tower)
	if !ok {
		return nil, errors.New("Agent needs a Tower world to operate")
	}
	return &Base{
		tower: tower,
		id:    uuid,
	}, nil
}

func (a *Base) Run() {
	floor := a.tower.Floor(a.id)
	log.Printf("An agent cycle executed from base agent %d", floor)
}

func (a *Base) HP() int {
	return a.tower.GetHP(a.id)
}

func (a *Base) Floor() int {
	return a.tower.GetFloor(a.id)
}

func (a *Base) ID() string {
	return a.id
}

func (a *Base) IsDead() bool {
	return a.tower.Exists(a.ID())
}

func (a *Base) TakeFood(amountOfFood float64) float64 {
	return a.tower.FoodRequest(a.id, amountOfFood)
}
