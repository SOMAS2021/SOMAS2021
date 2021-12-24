package infra

import (
	"math"
	"math/rand"
	"sync"

	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/day"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/health"
	log "github.com/sirupsen/logrus"
)

type Tower struct {
	currPlatFood   food.FoodType
	maxPlatFood    food.FoodType
	currPlatFloor  int
	agentCount     int
	agents         map[string]*Base
	agentsPerFloor int
	missingAgents  map[int][]int // key: floor, value: types of missing agents
	logger         log.Entry
	dayInfo        *day.DayInfo
	healthInfo     *health.HealthInfo
	mx             sync.RWMutex
}

func (t *Tower) Log(message string, fields ...Fields) {
	if len(fields) == 0 {
		fields = append(fields, Fields{})
	}
	t.logger.WithFields(fields[0]).Info(message)
}

func (t *Tower) TowerStateLog(timeOfTick string) {
	t.Log("Reporting platform status"+timeOfTick, Fields{"food_left": t.currPlatFood, "floor": t.currPlatFloor})
}

func NewTower(maxPlatFood food.FoodType, agentCount,
	agentsPerFloor int, dayInfo *day.DayInfo, healthInfo *health.HealthInfo) *Tower {
	return &Tower{
		currPlatFood:   maxPlatFood,
		maxPlatFood:    maxPlatFood,
		currPlatFloor:  1,
		agentCount:     agentCount,
		agents:         make(map[string]*Base),
		agentsPerFloor: agentsPerFloor,
		missingAgents:  make(map[int][]int),
		logger:         *log.WithFields(log.Fields{"reporter": "tower"}),
		dayInfo:        dayInfo,
		healthInfo:     healthInfo,
	}
}

func (t *Tower) Tick() {
	t.TowerStateLog(" end of tick")
	//useful parameters
	numOfFloors := t.agentCount / t.agentsPerFloor

	// Shuffle the agents
	if t.dayInfo.CurrTick%t.dayInfo.TicksPerReshuffle == 0 {
		t.reshuffle(numOfFloors)
	}
	// Move the platform
	if t.dayInfo.CurrTick%t.dayInfo.TicksPerFloor == 0 {
		t.currPlatFloor++
	}
	// Decrease agent HP and reset tower at end of day
	if t.dayInfo.CurrTick%t.dayInfo.TicksPerDay == 0 {
		t.hpDecay() // decreases HP and kills if < 0
		t.ResetTower()
		t.Log("-----------------END----OF----DAY-----------------", Fields{})
	}
}

func (t *Tower) UpdateMissingAgents() map[int][]int {
	deadAgents := t.missingAgents
	t.missingAgents = make(map[int][]int)
	return deadAgents
}

func (t *Tower) AddAgent(bagent *Base) {
	t.agents[bagent.id] = bagent
}

func (t *Tower) reshuffle(numOfFloors int) {
	remainingVacancies := make([]int, numOfFloors)
	t.Log("Reshuffling alive agents...", Fields{"agents_count": len(t.agents)})
	for i := 0; i < numOfFloors; i++ { // adding a max to each floor
		remainingVacancies[i] = t.agentsPerFloor
	}
	// allocating agents to floors randomly
	// iterate through the uuid strings of each agent
	for _, agent := range t.agents {
		newFloor := rand.Intn(numOfFloors)
		for remainingVacancies[newFloor] == 0 {
			newFloor = rand.Intn(numOfFloors)
		}
		agent.setFloor(newFloor + 1)
		remainingVacancies[newFloor]--
	}
}

func (t *Tower) hpDecay() {
	for _, agent := range t.agents {
		newHP := 0

		if agent.hp >= t.healthInfo.WeakLevel {
			newHP = int(math.Min(float64(t.healthInfo.MaxHP), float64(agent.hp)-(10.0+float64(agent.hp-t.healthInfo.WeakLevel)*0.25)))
		} else {
			if agent.hp >= t.healthInfo.HPCritical+t.healthInfo.HPReqCToW {
				newHP = t.healthInfo.WeakLevel
				agent.daysAtCritical = 0
			} else {
				newHP = t.healthInfo.HPCritical
				agent.daysAtCritical++
			}
		}

		if newHP < t.healthInfo.WeakLevel {
			newHP = t.healthInfo.HPCritical
		}

		agent.setHasEaten(false)
		if agent.daysAtCritical >= t.healthInfo.MaxDayCritical {
			t.Log("Killing agent", Fields{"agent": agent.id})
			t.missingAgents[agent.floor] = append(t.missingAgents[agent.floor], agent.agentType)
			delete(t.agents, agent.id) // maybe lock mutex?
		} else {
			agent.setHP(newHP)
		}
	}
}

func (t *Tower) SendMessage(direction int, senderFloor int, msg messages.Message) {
	for _, agent := range t.agents {
		if agent.floor == senderFloor+direction {
			agent.mx.Lock()
			agent.inbox.PushBack(msg)
			agent.mx.Unlock()
		}
	}
}

func (t *Tower) ResetTower() {
	t.currPlatFood = t.maxPlatFood
	t.currPlatFloor = 1
	t.Log("Tower Reset", Fields{})
}

func (t *Tower) TotalAgents() int {
	t.mx.RLock()
	defer t.mx.RUnlock()
	return len(t.agents)
}
