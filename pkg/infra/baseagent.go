package infra

import (
	"container/list"
	"errors"
	"math"
	"sync"

	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/abm"
	log "github.com/sirupsen/logrus"
)

type Fields = log.Fields

type Base struct {
	id        string
	hp        int
	floor     int
	agentType int
	inbox     *list.List
	tower     *Tower
	mx        sync.RWMutex
	logger    log.Entry
	hasEaten  bool
}

func NewBaseAgent(a *abm.ABM, agentType int, agentHP int, agentFloor int, id string) (*Base, error) {
	world := a.World()
	if world == nil {
		return nil, errors.New("agent needs a world defined to operate")
	}
	tower, ok := world.(*Tower)
	if !ok {
		return nil, errors.New("agent needs a tower world to operate")
	}
	logger := log.WithFields(log.Fields{"agent_id": id, "agent_type": agentType, "reporter": "agent"})
	return &Base{
		id:        id,
		hp:        agentHP,
		floor:     agentFloor,
		agentType: agentType,
		tower:     tower,
		inbox:     list.New(),
		logger:    *logger,
		hasEaten:  false,
	}, nil
}

func (a *Base) Log(message string, fields ...Fields) {
	if len(fields) == 0 {
		fields = append(fields, Fields{})
	}
	a.logger.WithFields(fields[0]).Info(message)
}

func (a *Base) Run() {
	a.Log("An agent cycle executed from base agent", Fields{"floor": a.floor, "hp": a.hp})
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

func (a *Base) IsAlive() bool {
	_, found := a.tower.agents[a.id]
	return found
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

func (a *Base) HasEaten() bool {
	return a.hasEaten
}

func (a *Base) UpdateHasEaten(newStatus bool) {
	a.hasEaten = newStatus
}

func (a *Base) TakeFood(amountOfFood float64) float64 {
	if a.floor == a.tower.currPlatFloor && !a.HasEaten() {
		foodTaken := math.Min(a.tower.currPlatFood, amountOfFood)
		a.updateHP(foodTaken)
		a.tower.currPlatFood -= foodTaken
		a.UpdateHasEaten(true)
		a.Log("An agent has taken food", Fields{"floor": a.floor, "amount": foodTaken})
		return foodTaken
	}
	return 0.0
}

func (a *Base) ReceiveMessage() messages.Message {
	a.mx.Lock()
	defer a.mx.Unlock()
	if a.inbox.Len() > 0 {
		msg := a.inbox.Front().Value.(messages.Message)
		(a.inbox).Remove(a.inbox.Front())
		return msg
	}
	return nil
}

func (a *Base) SendMessage(direction int, msg messages.Message) {
	if (direction == -1) || (direction == 1) {
		a.tower.SendMessage(direction, a.floor, msg)
	}
}
