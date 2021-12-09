package infra

import (
	"errors"
	"log"
	"math"

	"github.com/SOMAS2021/SOMAS2021/pkg/utils/abm"
)

type Agent interface {
	Run()
	Die()
	IsAlive() bool
}

type Base struct {
	isAlive   bool
	id        string
	hp        int
	floor     int
	agentType int
	tower     *Tower
}

func NewBaseAgent(abm *abm.ABM, uuid string) (*Base, error) {
	world := abm.World()
	if world == nil {
		return nil, errors.New("Agent needs a World defined to operate")
	}
	tower, ok := world.(*Tower)
	if !ok {
		return nil, errors.New("Agent needs a Tower world to operate")
	}
	return &Base{
		tower: tower,
		id:    uuid,
	}, nil
}

func (a *Base) Run() {
	floor := a.floor
	log.Printf("An agent cycle executed from base agent %d", floor)
}

func (a *Base) HP() int {
	return a.hp
}

func (a *Base) Floor() int {
	return a.floor
}

func (a *Base) ID() string {
	return a.id
}

func (a *Base) Die() {
	a.isAlive = false
}

func (a *Base) IsAlive() bool {
	return a.isAlive
}

func (a *Base) setFloor(newFloor int) {
	a.floor = newFloor
}

func (a *Base) setHP(newHP int) {
	a.hp = newHP
}

func (a *Base) updateHP(foodTaken float64) {
	a.hp = int(math.Min(100, float64(a.hp)+foodTaken))
}

func (a *Base) TakeFood(amountOfFood float64) float64 {
	if uint64(a.floor) == a.tower.currPlatFloor {
		foodTaken := math.Min(a.tower.currPlatFood, amountOfFood)
		a.updateHP(foodTaken)
		a.tower.currPlatFood -= foodTaken
		return foodTaken
	}
	return 0.0
	// return a.tower.FoodRequest(a.id, amountOfFood)
}
