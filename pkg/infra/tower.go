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
	Agents         map[string]Agent
	agentsPerFloor int
	logger         log.Entry
	dayInfo        *day.DayInfo
	healthInfo     *health.HealthInfo
	mx             sync.RWMutex
	deadAgents     map[int]int
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
		Agents:         make(map[string]Agent),
		agentsPerFloor: agentsPerFloor,
		logger:         *log.WithFields(log.Fields{"reporter": "tower"}),
		dayInfo:        dayInfo,
		healthInfo:     healthInfo,
		deadAgents:     make(map[int]int),
	}
}

func (t *Tower) Tick() {
	t.TowerStateLog(" end of tick")

	// Shuffle the agents
	if t.dayInfo.CurrTick%t.dayInfo.TicksPerReshuffle == 0 {
		t.Reshuffle()
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

func (t *Tower) AddAgent(agent Agent) {
	t.Agents[agent.BaseAgent().id] = agent
}

// This function shuffles the agents by generating a random permutation of agentCount intgers,
// and maps the integers into floors by dividing each integer by the number of agents per floor.
// This function does not guarantee that an agent will be moved to a different floor.
func (t *Tower) Reshuffle() {
	t.Log("Shuffling agents")
	newFloors := rand.Perm(t.agentCount)
	i := 0
	for _, agent := range t.Agents {
		agent := agent.BaseAgent()
		newFloor := newFloors[i]/t.agentsPerFloor + 1

		t.Log("Floor change", Fields{"agent_id": agent.ID(), "old_floor": agent.Floor(), "new_floor": newFloor})

		agent.setFloor(newFloor)
		i++
	}
}

func (t *Tower) hpDecay() {
	for _, agent := range t.Agents {
		agent := agent.BaseAgent()
		newHP := 0

		if agent.hp >= t.healthInfo.WeakLevel {
			newHP = int(math.Min(float64(t.healthInfo.MaxHP), float64(agent.hp)-(float64(t.healthInfo.HPLossBase)+float64(agent.hp-t.healthInfo.WeakLevel)*t.healthInfo.HPLossSlope)))
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
			newHP = 0
		}
		agent.setHP(newHP)
	}
}

func (t *Tower) SendMessage(direction int, senderFloor int, msg messages.Message) {
	for _, agent := range t.Agents {
		agent := agent.BaseAgent()
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
	return len(t.Agents)
}

func (t *Tower) UpdateDeadAgents(agentType int) {
	t.deadAgents[agentType]++
}

func (t *Tower) DeadAgents() map[int]int {
	return t.deadAgents
}
