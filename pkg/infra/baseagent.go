package infra

import (
	"container/list"
	"errors"
	"log"
	"math"
	"sync"

	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/abm"
)

type Base struct {
	id        string
	hp        int
	floor     int
	agentType int
	inbox     *list.List
	tower     *Tower
	mx        sync.RWMutex
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
	return &Base{
		id:        id,
		hp:        agentHP,
		floor:     agentFloor,
		agentType: agentType,
		tower:     tower,
		inbox:     list.New(),
	}, nil
}

func (a *Base) Run() {
	floor := a.floor
	log.Printf("An agent cycle executed from base agent %d", floor)
}

func (a *Base) HP() int {
	return a.hp
}

func (a *Base) CurrPlatFood() float64 {
	foodOnPlatform, platformFloor := a.tower.GetCurrPlatFood()
	if platformFloor == a.floor { // only show the food on the platform if the platform is on the same floor as the agent
		return foodOnPlatform // if agent and platform on same floor, return the amount of food on the platform
	}
	return -1.0 // otherwise return -1 to indicate the platform is not on the floor. Agents cannot see the food on the platform if it is not on their floor
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

func (a *Base) TakeFood(amountOfFood float64) float64 {
	if a.floor == a.tower.currPlatFloor {
		foodTaken := math.Min(a.tower.currPlatFood, amountOfFood)
		a.updateHP(foodTaken)
		a.tower.currPlatFood -= foodTaken
		return foodTaken
	}
	return 0.0
}

func (a *Base) ReceiveMessage() messages.Message {
	log.Printf("Tower receive message")
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
	log.Printf("agent sending message")
	if (direction == -1) || (direction == 1) {
		a.tower.SendMessage(direction, a.floor, msg)
	}
}
