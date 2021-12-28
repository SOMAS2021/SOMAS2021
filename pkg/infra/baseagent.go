package infra

import (
	"errors"
	"math"

	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/health"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/world"
	log "github.com/sirupsen/logrus"
)

type Agent interface {
	Run()
	IsAlive() bool
	BaseAgent() *Base
	HandleAskHP(msg messages.AskHPMessage)
	HandleAskFoodTaken(msg messages.AskFoodTakenMessage)
	HandleAskIntendedFoodTaken(msg messages.AskIntendedFoodIntakeMessage)
	HandleRequestLeaveFood(msg messages.RequestLeaveFoodMessage)
	HandleRequestTakeFood(msg messages.RequestTakeFoodMessage)
	HandleResponse(msg messages.BoolResponseMessage)
	HandleStateFoodTaken(msg messages.StateFoodTakenMessage)
	HandleStateHP(msg messages.StateHPMessage)
	HandleStateIntendedFoodTaken(msg messages.StateIntendedFoodIntakeMessage)
}

type Fields = log.Fields

type Base struct {
	id             string
	hp             int
	floor          int
	agentType      int
	inbox          chan messages.Message
	tower          *Tower
	logger         log.Entry
	hasEaten       bool
	daysAtCritical int
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
		//TODO: Check how large to make the inbox channel. Currently set to 15.
		inbox:          make(chan messages.Message, 15),
		logger:         *logger,
		hasEaten:       false,
		daysAtCritical: 0,
	}, nil
}

func (a *Base) BaseAgent() *Base {
	return a
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
func (a *Base) CurrPlatFood() food.FoodType {
	foodOnPlatform := a.tower.currPlatFood
	platformFloor := a.tower.currPlatFloor
	if platformFloor == a.floor || platformFloor == a.floor+1 {
		return foodOnPlatform
	}
	return -1
}

func (a *Base) Floor() int {
	return a.floor
}

func (a *Base) ID() string {
	return a.id
}

func (a *Base) IsAlive() bool {
	return a.hp > 0
}

func (a *Base) AgentType() int {
	return a.agentType
}

func (a *Base) DaysAtCritical() int {
	return a.daysAtCritical
}

func (a *Base) setFloor(newFloor int) {
	a.floor = newFloor
}

func (a *Base) setHP(newHP int) {
	a.hp = newHP
}

// Modeled as a first order system step answer (see documentation for more information)
func (a *Base) updateHP(foodTaken food.FoodType) {
	hpChange := a.tower.healthInfo.Width * (1 - math.Pow(math.E, -float64(foodTaken)/a.tower.healthInfo.Tau))
	if a.hp >= a.tower.healthInfo.WeakLevel {
		a.hp = a.hp + int(hpChange)
	} else {
		a.hp = int(math.Min(float64(a.tower.healthInfo.HPCritical+a.tower.healthInfo.HPReqCToW), float64(a.hp)+hpChange))
	}

}

func (a *Base) HasEaten() bool {
	return a.hasEaten
}

func (a *Base) setHasEaten(newStatus bool) {
	a.hasEaten = newStatus
}

func (a *Base) TakeFood(amountOfFood food.FoodType) food.FoodType {
	if amountOfFood == 0 {
		return 0
	}
	if a.floor == a.tower.currPlatFloor && !a.hasEaten && amountOfFood > 0 {
		foodTaken := food.FoodType(math.Min(float64(a.tower.currPlatFood), float64(amountOfFood)))
		a.updateHP(foodTaken)
		a.tower.currPlatFood -= foodTaken
		a.setHasEaten(foodTaken > 0)
		a.Log("An agent has taken food", Fields{"floor": a.floor, "amount": foodTaken})
		return foodTaken
	}
	return -1
}

func (a *Base) ReceiveMessage() messages.Message {
	select {
	case msg := <-a.inbox:
		return msg
	default:
		return nil
	}
}

func (a *Base) SendMessage(direction int, msg messages.Message) {
	if (direction == -1) || (direction == 1) {
		a.tower.SendMessage(direction, a.floor, msg)
	}
}

func (a *Base) HealthInfo() *health.HealthInfo {
	return a.tower.healthInfo
}

func (a *Base) HandleAskHP(msg messages.AskHPMessage)                                    {}
func (a *Base) HandleAskFoodTaken(msg messages.AskFoodTakenMessage)                      {}
func (a *Base) HandleAskIntendedFoodTaken(msg messages.AskIntendedFoodIntakeMessage)     {}
func (a *Base) HandleRequestLeaveFood(msg messages.RequestLeaveFoodMessage)              {}
func (a *Base) HandleRequestTakeFood(msg messages.RequestTakeFoodMessage)                {}
func (a *Base) HandleResponse(msg messages.BoolResponseMessage)                          {}
func (a *Base) HandleStateFoodTaken(msg messages.StateFoodTakenMessage)                  {}
func (a *Base) HandleStateHP(msg messages.StateHPMessage)                                {}
func (a *Base) HandleStateIntendedFoodTaken(msg messages.StateIntendedFoodIntakeMessage) {}
