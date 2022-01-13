package infra

import (
	"math/rand"

	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/agent"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/day"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/health"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/logging"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type Tower struct {
	currPlatFood   food.FoodType
	maxPlatFood    food.FoodType
	currPlatFloor  int
	agentCount     int
	Agents         map[uuid.UUID]Agent
	agentsPerFloor int
	logger         log.Entry
	stateLog       logging.StateLog
	dayInfo        *day.DayInfo
	healthInfo     *health.HealthInfo
	deadAgents     map[agent.AgentType]int
}

func (t *Tower) Log(message string, fields ...Fields) {
	if len(fields) == 0 {
		fields = append(fields, Fields{})
	}
	t.logger.WithFields(fields[0]).Info(message)
}

func (t *Tower) TowerStateLog(timeOfTick string) {
	t.Log("Reporting platform status"+timeOfTick, Fields{"food_left": t.currPlatFood, "floor": t.currPlatFloor})
	t.stateLog.LogPlatFoodState(t.dayInfo, int(t.currPlatFood))
}

func NewTower(maxPlatFood food.FoodType, agentCount,
	agentsPerFloor int, dayInfo *day.DayInfo, healthInfo *health.HealthInfo, stateLog *logging.StateLog) *Tower {
	return &Tower{
		currPlatFood:   maxPlatFood,
		maxPlatFood:    maxPlatFood,
		currPlatFloor:  1,
		agentCount:     agentCount,
		Agents:         make(map[uuid.UUID]Agent),
		agentsPerFloor: agentsPerFloor,
		logger:         *stateLog.Logmanager.GetLogger("main").WithFields(log.Fields{"reporter": "tower"}),
		stateLog:       *stateLog,
		dayInfo:        dayInfo,
		healthInfo:     healthInfo,
		deadAgents:     make(map[agent.AgentType]int),
	}
}

func (t *Tower) Tick() {
	// Shuffle the agents
	if t.dayInfo.CurrTick%t.dayInfo.TicksPerReshuffle == 0 {
		t.Reshuffle()
	}
	// Move the platform
	if t.dayInfo.CurrTick%t.dayInfo.TicksPerFloor == 0 {
		t.currPlatFloor++
		t.stateLog.LogStoryPlatformMoved(t.dayInfo, t.currPlatFloor)
	}
	// Decrease agent HP and reset tower at end of day
	if t.dayInfo.CurrTick%t.dayInfo.TicksPerDay == 0 {
		t.endOfDay()
		t.ResetTower()
		t.Log("-----------------END----OF----DAY-----------------", Fields{})
	}
	// Reset dead agent's HP to enable reincarnation
	//t.resetHP()
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

func (t *Tower) endOfDay() {
	t.dayInfo.CurrDay += 1
	for _, agent := range t.Agents {
		agent := agent.BaseAgent()
		agent.hpDecay(t.healthInfo)
		agent.increaseAge()
		agent.updateTreaties()
		t.stateLog.LogUtility(t.dayInfo, agent.agentType, agent.utility(), agent.IsAlive())
	}
}

func (t *Tower) resetHP() {
	for _, agent := range t.Agents {
		agent := agent.BaseAgent()
		if !agent.IsAlive() && agent.agentType.String() == "Team2" {
			agent.setHP(t.healthInfo.MaxHP)
			agent.daysAtCritical = 0
			agent.age = 0
		}
	}
}

func (t *Tower) SendMessage(senderFloor int, msg messages.Message) {
	direction := 1
	if msg.SenderFloor() > msg.TargetFloor() {
		direction = -1
	}
	for _, agent := range t.Agents {
		agent := agent.BaseAgent()
		if agent.floor == senderFloor+direction {
			select {
			case agent.inbox <- msg:
			default:
				t.Log("message send failed")
			}
		}
	}
}

func (t *Tower) ResetTower() {
	t.currPlatFood = t.maxPlatFood
	t.currPlatFloor = 1
	t.Log("Tower Reset", Fields{})
}

func (t *Tower) TotalAgents() int {
	return len(t.Agents)
}

func (t *Tower) UpdateDeadAgents(agentType agent.AgentType) {
	t.deadAgents[agentType]++
}

func (t *Tower) DeadAgents() map[agent.AgentType]int {
	return t.deadAgents
}
