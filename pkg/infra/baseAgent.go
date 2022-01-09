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
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/logging"

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
	CustomLogs()
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
	totalFoodTaken food.FoodType
	totalFoodSeen  food.FoodType
	totalHPGained  int
	totalHPLost    int
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
		totalFoodTaken: 0,
		totalFoodSeen:  0,
		totalHPGained:  0,
		totalHPLost:    0,
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

func (a *Base) updateTreaties() {
	for _, treaty := range a.activeTreaties {
		treaty.DecrementDuration()
		if treaty.Duration() == 0 {
			delete(a.activeTreaties, treaty.ID())
		}
	}
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
	hpOld := a.hp
	hpChange := int(math.Round(a.tower.healthInfo.Width * (1 - math.Pow(math.E, -float64(foodTaken)/a.tower.healthInfo.Tau))))
	if a.hp >= a.tower.healthInfo.WeakLevel {
		a.hp = a.hp + hpChange
	} else {
		a.hp = utilFunctions.MinInt(a.tower.healthInfo.HPCritical+a.tower.healthInfo.HPReqCToW, a.hp+hpChange)
	}

	// For utility calculation
	hpDiff := a.hp - hpOld
	a.UpdateHPChange(hpDiff)
}

func (a *Base) hpDecay(healthInfo *health.HealthInfo) {
	newHP := 0
	if a.hp >= healthInfo.WeakLevel {
		newHP = utilFunctions.MinInt(healthInfo.MaxHP, a.hp-(healthInfo.HPLossBase+int(math.Round(float64(a.hp-healthInfo.WeakLevel)*healthInfo.HPLossSlope))))
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

	// For utility calculation
	hpLoss := a.hp - newHP
	a.totalHPLost += hpLoss

	a.setHasEaten(false)
	if a.daysAtCritical >= healthInfo.MaxDayCritical {
		a.Log("Killing agent", Fields{"daysLived": a.Age(), "agentType": a.agentType})
		a.tower.stateLog.LogAgentDeath(a.tower.dayInfo, a.agentType, a.Age())
		a.tower.stateLog.LogStoryAgentDied(a.tower.dayInfo, a.storyState())
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
	foodTaken := food.FoodType(utilFunctions.MinInt(int(a.tower.currPlatFood), int(amountOfFood), int(a.tower.healthInfo.MaxFoodIntake)))
	a.updateHP(foodTaken)
	a.tower.currPlatFood -= foodTaken
	a.setHasEaten(foodTaken > 0)
	a.UpdateFoodTaken(foodTaken)
	a.Log("An agent has taken food", Fields{"floor": a.floor, "amount": foodTaken})
	if foodTaken > 0 {
		a.tower.stateLog.LogStoryAgentTookFood(
			a.tower.dayInfo,
			a.storyState(),
			int(foodTaken),
			int(a.tower.currPlatFood),
		)
	}
	return foodTaken, nil
}

func (a *Base) UpdateFoodSeen(amountSeen food.FoodType) {
	a.totalFoodSeen += amountSeen
}

func (a *Base) UpdateFoodTaken(intake food.FoodType) {
	a.totalFoodTaken += intake
}

func (a *Base) UpdateHPChange(change int) {
	if change < 0 {
		a.totalHPLost -= change
	} else {
		a.totalHPGained += change
	}
}

func (a *Base) utility() float64 {
	// TODO: Improve utility calculation
	// Currently just use ratio between food taken and food seen
	// Somehow also factor
	// Ideas: piecewise/bounded function with log, concave quadratic

	hpGainFactor := netHPGainFactor(a.totalHPGained, a.totalHPLost, a.tower.healthInfo.MaxHP)
	foodIntakeFactor := foodIntakeFactor(a.totalFoodTaken, a.totalFoodSeen)
	// this means the multiplcation of these values can be larger than 1
	return simpleQuadraticUtility(hpGainFactor * foodIntakeFactor)
}

/*
HP Gain Factor.
inputs:
- gain >= 0
- loss >= 0

output range:
- y >= 0.

key values:
- = 1: agent is sustaining overall
- < 1: agent is losing hp on average
- > 1: agent gaining hp on average
- when loss is 0, default to percentage difference between gain and loss. Acknowledge that HP gain should increase utility
*/
// func hpGainFactor(gain, loss, maxhp int) float64 {
// 	if loss == 0 {
// 		return 1 + float64(gain)/float64(maxhp)
// 	}
// 	return float64(gain) / float64(loss)
// }

/*
Net HP Gain Factor.
This one uses net hp gain as opposed to hp gain/loss ratio.
*/
func netHPGainFactor(gain, loss, maxhp int) float64 {
	netHPGain := gain - loss
	multiplier := 1 + float64(netHPGain)/float64(maxhp) // reward positive net hp gain, punish otherwise
	return utilFunctions.RestrictToRange(0, math.Inf(1), multiplier)
}

/*
Food Intake Factor.
inputs:
- intake >= 0
- total >= 0
range: [0, 1]. == 1, agent has taken all the food it has seen

*/
func foodIntakeFactor(intake, seen food.FoodType) float64 {
	// If an agent hasn't seen any food, has 0 utility
	if seen == 0 {
		return 0
	}
	return float64(intake) / float64(seen)
}

// Quadratic function standard form: a(x-h)^2+k
func quadratic(x, a, h, k float64) float64 {
	return a*math.Pow((x-h), 2) + k
}

// Utility function is a strictly increasing function that starts at 0,0 ends at 1,1
func simpleQuadraticUtility(x float64) float64 {
	// -(x-1)^2 + 1
	return utilFunctions.RestrictToRange(0, 1, quadratic(x, -1, 1, 1))
}

func (a *Base) HealthInfo() *health.HealthInfo {
	return a.tower.healthInfo
}

func (a *Base) storyState() logging.AgentState {
	return logging.AgentState{
		HP:        a.hp,
		AgentType: a.agentType,
		Floor:     a.floor,
		Age:       a.age,
		Custom:    "",
		Utility:   a.utility(),
	}
}

func (a *Base) CustomLogs() {}
