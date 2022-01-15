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
	Behaviour() string
	BehaviourChange() string
	FoodReceived() string
	Utility() float64
	TotalFoodSeen() food.FoodType
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
	utility        float64
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
		utility:        0,
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
	return utilFunctions.MinInt(a.hp, a.tower.healthInfo.MaxHP)
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

func (a *Base) ActualCurrPlatFood() food.FoodType {
	foodOnPlatform := a.tower.currPlatFood
	platformFloor := a.tower.currPlatFloor
	if platformFloor == a.floor {
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
		a.tower.DeathCount++
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
	a.updateUtility(amountOfFood, foodTaken)
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

func (a *Base) TotalFoodSeen() food.FoodType {
	return a.totalFoodSeen
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

func (a *Base) updateUtility(foodRequested, foodTaken food.FoodType) {
	// Hyperparameters that are chosen to make utility commensurate for all agents
	// Requires that alpha > gamma > beta for diminishing returns + greed

	// The weighting assigned to the resources exactly needed to survive
	alpha := 0.2
	// The weighting assigned to bonus resources that aren't needed
	beta := 0.1
	// The weighting assigned to not getting the resources that are needed
	gamma := 0.18

	// 1.) Determines the resources it has available, gi ∈[0,1]
	foodAvailable := float64(a.CurrPlatFood()) / float64(a.tower.maxPlatFood)

	// 2.) Determines its need for resources, qi ∈[0,1]
	// If an agent is not about to die, it effectively needs 0 food
	foodToSurvive := float64(a.daysAtCritical) / float64(a.HealthInfo().MaxDayCritical)

	// 3.) Makes a provision of resources, pi ∈[0,1]
	// Agents never recontribute to the common pool (provision = 0)
	foodProvision := 0.0

	// 4.) Makes an appropriation of resources, r′i ∈[0,1]
	// Amount of food taken, normalised by the max possible intake
	foodAppropriated := float64(foodTaken) / float64(a.tower.healthInfo.MaxFoodIntake)

	// The total resources accrued at the end of a round is hence given by:
	// Ri = r′i + (gi −pi)
	// Resource generation by agent per day is what they actually take, plus the resources available, minus the provisions made
	totalFoodAccrued := foodAppropriated + (foodAvailable - foodProvision)

	// Which provides the utility of an agent, ui, as:
	// ui = αi(qi) + βi(Ri −qi), if Ri ≥qi
	// 		αi(Ri) −γi(qi −Ri), otherwise
	if totalFoodAccrued >= foodToSurvive {
		a.utility = alpha*foodToSurvive + beta*(totalFoodAccrued-foodToSurvive)
	} else {
		a.utility = alpha*totalFoodAccrued - gamma*(foodToSurvive-totalFoodAccrued)
	}
}

func (a *Base) HealthInfo() *health.HealthInfo {
	return a.tower.healthInfo
}

func (a *Base) storyState() logging.AgentState {
	return logging.AgentState{
		HP:        a.HP(),
		AgentType: a.agentType,
		Floor:     a.floor,
		Age:       a.age,
		Custom:    "",
		Utility:   a.utility,
	}
}

func (a *Base) CustomLogs() {}

func (a *Base) Behaviour() string {
	return "Not Team 6"
}

func (a *Base) BehaviourChange() string {
	return "Not Team 6"
}

func (a *Base) FoodReceived() string {
	return "Not Team 6"
}

func (a *Base) Utility() float64 {
	return a.utility
}

func (a *Base) SMCtr() []int {
	ctr := []int{
		a.tower.dayInfo.BehaviourCtr["Altruist"],
		a.tower.dayInfo.BehaviourCtr["Collectivist"],
		a.tower.dayInfo.BehaviourCtr["Selfish"],
		a.tower.dayInfo.BehaviourCtr["Narcissist"],
	}
	return ctr
}
