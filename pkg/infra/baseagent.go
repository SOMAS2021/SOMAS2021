package infra

import (
	"errors"
	"fmt"
	"math"

	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/agent"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/health"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/world"

	"github.com/SOMAS2021/SOMAS2021/pkg/utils/utilFunctions"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type Agent interface {
	Run()
	IsAlive() bool
	Floor() int
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
	HandleProposeTreaty(msg messages.ProposeTreatyMessage)
	HandleTreatyResponse(msg messages.TreatyResponseMessage)
	HandlePropogate(msg messages.Message)
}

type Fields = log.Fields

type Base struct {
	id             uuid.UUID
	hp             int
	floor          int
	agentType      agent.AgentType
	inbox          chan messages.Message
	tower          *Tower
	logger         log.Entry
	hasEaten       bool
	daysAtCritical int
	age            int
	activeTreaties map[uuid.UUID]messages.Treaty
}

func NewBaseAgent(world world.World, agentType agent.AgentType, agentHP int, agentFloor int, id uuid.UUID) (*Base, error) {
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
		//TODO: Check how large to make the inbox channel. Currently set to 15.
		inbox:          make(chan messages.Message, 15),
		logger:         *tower.stateLog.Logmanager.GetLogger("main").WithFields(log.Fields{"agent_id": id, "agent_type": agentType.String(), "reporter": "agent"}),
		hasEaten:       false,
		daysAtCritical: 0,
		age:            0,
		activeTreaties: make(map[uuid.UUID]messages.Treaty),
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
	return utilFunctions.MinInt(a.hp, 100)
}

func (a *Base) Age() int {
	return a.age
}

func (a *Base) increaseAge() {
	a.age++
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

func (a *Base) ID() uuid.UUID {
	return a.id
}

func (a *Base) IsAlive() bool {
	return a.hp > 0
}

func (a *Base) AgentType() agent.AgentType {
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
	hpChange := int(a.tower.healthInfo.Width * (1 - math.Pow(math.E, -float64(foodTaken)/a.tower.healthInfo.Tau)))
	if a.hp >= a.tower.healthInfo.WeakLevel {
		a.hp = a.hp + hpChange
	} else {
		a.hp = utilFunctions.MinInt(a.tower.healthInfo.HPCritical+a.tower.healthInfo.HPReqCToW, a.hp+hpChange)
	}
}

func (a *Base) hpDecay(healthInfo *health.HealthInfo) {
	newHP := 0
	if a.hp >= healthInfo.WeakLevel {
		newHP = utilFunctions.MinInt(healthInfo.MaxHP, a.hp-(healthInfo.HPLossBase+int(float64(a.hp-healthInfo.WeakLevel)*healthInfo.HPLossSlope)))
	} else {
		if a.hp >= healthInfo.HPCritical+healthInfo.HPReqCToW {
			newHP = healthInfo.WeakLevel
			a.daysAtCritical = 0
		} else {
			newHP = healthInfo.HPCritical
			a.daysAtCritical++
		}
	}
	if newHP < healthInfo.WeakLevel {
		newHP = healthInfo.HPCritical
	}
	a.setHasEaten(false)
	if a.daysAtCritical >= healthInfo.MaxDayCritical {
		a.Log("Killing agent", Fields{"daysLived": a.Age(), "agentType": a.agentType})
		a.tower.stateLog.LogAgentDeath(a.tower.dayInfo.CurrDay, a.tower.dayInfo.CurrTick, a.agentType)
		newHP = 0
	}
	a.Log("Setting hp to " + fmt.Sprint(newHP))
	a.setHP(newHP)
}

func (a *Base) HasEaten() bool {
	return a.hasEaten
}

func (a *Base) setHasEaten(newStatus bool) {
	a.hasEaten = newStatus
}

func (a *Base) PlatformOnFloor() bool {
	return a.floor == a.tower.currPlatFloor
}

func (a *Base) TakeFood(amountOfFood food.FoodType) (food.FoodType, error) {
	if a.floor != a.tower.currPlatFloor {
		return 0, &FloorError{}
	}
	if a.hasEaten {
		return 0, &AlreadyEatenError{}
	}
	if amountOfFood < 0 {
		return 0, &NegFoodError{}
	}
	foodTaken := food.FoodType(math.Min(float64(a.tower.currPlatFood), float64(amountOfFood)))
	a.updateHP(foodTaken)
	a.tower.currPlatFood -= foodTaken
	a.setHasEaten(foodTaken > 0)
	a.Log("An agent has taken food", Fields{"floor": a.floor, "amount": foodTaken})
	return foodTaken, nil
}

func (a *Base) ReceiveMessage() messages.Message {
	select {
	case msg := <-a.inbox:
		return msg
	default:
		return nil
	}
}

func (a *Base) SendMessage(msg messages.Message) {
	a.tower.SendMessage(a.floor, msg)
}

func (a *Base) HealthInfo() *health.HealthInfo {
	return a.tower.healthInfo
}

func (a *Base) ActiveTreaties() map[uuid.UUID]messages.Treaty {
	return a.activeTreaties
}

func (a *Base) AddToActiveTreaties(treaty messages.Treaty) {
	a.activeTreaties[treaty.ID()] = treaty
}

func (a *Base) DeleteFromActiveTreaties(treaty messages.Treaty) {
	delete(a.activeTreaties, treaty.ID())
}

func (a *Base) RespondToTreaty(msg messages.ProposeTreatyMessage, senderFloor int, response bool) messages.TreatyResponseMessage {
	treatyToReview := msg.Treaty()
	if response {
		treatyToReview.SignTreaty()
		a.AddToActiveTreaties(treatyToReview)
	}
	reply := msg.Reply(a.ID(), senderFloor, response)

	return reply
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
func (a *Base) HandleProposeTreaty(msg messages.ProposeTreatyMessage)                    {}
func (a *Base) HandleTreatyResponse(msg messages.TreatyResponseMessage)                  {}

func (a *Base) HandlePropogate(msg messages.Message) {
	a.SendMessage(msg)
}
