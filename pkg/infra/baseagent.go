package infra

import (
	"container/list"
	"errors"
	"math"
	"sync"

	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/world"
	log "github.com/sirupsen/logrus"
)

type Fields = log.Fields

type Base struct {
	id          string
	hp          int
	floor       int
	agentType   int
	inbox       *list.List
	tower       *Tower
	mx          sync.RWMutex
	logger      log.Entry
	hasEaten    bool
	daysInState int
}

func NewBaseAgent(world world.World, agentType int, agentHP int, agentFloor int, id string) (*Base, error) {
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

// only show the food on the platform if the platform is on the
// same floor as the agent or directly below
func (a *Base) CurrPlatFood() float64 {
	foodOnPlatform := a.tower.currPlatFood
	platformFloor := a.tower.currPlatFloor
	if platformFloor == a.floor || platformFloor == a.floor+1 {
		return foodOnPlatform
	}
	return -1.0
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

func (a *Base) updateHP(foodTaken float64) { // first order system step answer. 63% of levelWidth obtained for tau food
	// tau is the time constant of the first order system
	tau := 0.0
	base := 0
	width := 0.0
	shift := 0.0

	switch {
	case a.hp >= a.tower.healthInfo.StrongLevel:
		tau = a.tower.healthInfo.FoodReqStrong
		base = a.tower.healthInfo.StrongLevel
		width = a.tower.healthInfo.WidthStrong
		shift = -tau * math.Log(1-(float64(a.hp)-float64(base))/width)

	case a.hp >= a.tower.healthInfo.HealthyLevel:
		tau = a.tower.healthInfo.FoodReqHealthy
		base = a.tower.healthInfo.HealthyLevel
		width = a.tower.healthInfo.WidthHealthy
		shift = -tau * math.Log(1-(float64(a.hp)-float64(base))/width)

	case a.hp >= a.tower.healthInfo.WeakLevel:
		tau = a.tower.healthInfo.FoodReqWeak
		base = a.tower.healthInfo.WeakLevel
		width = a.tower.healthInfo.WidthWeak
		shift = -tau * math.Log(1-(float64(a.hp)-float64(base))/width)

	// Critical Level - TODO: discuss its implementation
	case a.hp >= 0:
		if foodTaken >= a.tower.healthInfo.FoodReqCToW {
			a.hp = int(a.tower.healthInfo.WeakLevel) - 1
		} else {
			a.hp = 1
		}
		return
	}

	a.hp = int(float64(base) + width*(1-math.Pow(math.E, -(foodTaken+shift)/tau)))
}

func (a *Base) HasEaten() bool {
	return a.hasEaten
}

func (a *Base) setHasEaten(newStatus bool) {
	a.hasEaten = newStatus
}

func (a *Base) TakeFood(amountOfFood float64) float64 {
	if a.floor == a.tower.currPlatFloor && !a.hasEaten {
		foodTaken := math.Min(a.tower.currPlatFood, amountOfFood)
		a.updateHP(foodTaken)
		a.tower.currPlatFood -= foodTaken
		a.setHasEaten(true)
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
