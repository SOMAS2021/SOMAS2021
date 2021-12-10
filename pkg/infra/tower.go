package infra

import (
	"math/rand"

	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
	log "github.com/sirupsen/logrus"
)

type Tower struct {
	currPlatFood    float64
	maxPlatFood     float64
	currPlatFloor   int
	agentCount      int
	agents          map[string]*Base
	agentsPerFloor  int
	ticksPerDay     int
	missingAgents   map[int][]int // key: floor, value: types of missing agents
	reshufflePeriod int
	tickCounter     int
	logger          log.Entry
}

func (t *Tower) Log(message string, fields ...Fields) {
	if len(fields) == 0 {
		fields = append(fields, Fields{})
	}
	t.logger.WithFields(fields[0]).Info(message)
}

func NewTower(currPlatFood float64, currPlatFloor, agentCount,
	agentsPerFloor, ticksPerDay, reshufflePeriod int) *Tower {
	t := &Tower{
		currPlatFood:    currPlatFood,
		maxPlatFood:     currPlatFood,
		currPlatFloor:   currPlatFloor,
		agentCount:      agentCount,
		agents:          make(map[string]*Base),
		agentsPerFloor:  agentsPerFloor,
		ticksPerDay:     ticksPerDay,
		missingAgents:   make(map[int][]int),
		reshufflePeriod: reshufflePeriod,
		tickCounter:     1,
		logger:          *log.WithFields(log.Fields{"reporter": "tower"}),
	}
	return t
}

func (t *Tower) Tick() {
	//logs
	t.Log("Reporting food left on platform", Fields{"food_left": t.currPlatFood})

	//useful parameters
	day := 24 * 60
	numOfFloors := t.agentCount / t.agentsPerFloor
	platformMovePeriod := day / numOfFloors // can add min/max

	// Shuffle the agents
	if t.tickCounter%t.reshufflePeriod == 0 {
		t.reshuffle(numOfFloors)
	}
	// Move the platform
	if t.tickCounter%platformMovePeriod == 0 {
		t.currPlatFloor++
	}
	// Decrease agent HP and reset tower at end of day
	if t.tickCounter%day == 0 {
		t.hpDecay() // decreases HP and kills if < 0
		t.ResetTower()
	}

	t.tickCounter++
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
	// TODO: can add a parameter
	for _, agent := range t.agents {
		newHP := agent.HP() - 20
		if newHP < 0 {
			log.Info("Killing agent %s", agent.ID())
			t.missingAgents[agent.Floor()] = append(t.missingAgents[agent.Floor()], agent.agentType)
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
}
