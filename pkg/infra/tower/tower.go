package tower

import (
	"log"
	"sync"
)

type Tower struct {
	currPlatFood    float64
	maxPlatFood     float64
	currPlatFloor   int
	mx              sync.RWMutex
	agentCount      int
	agents          map[string]BaseAgentCore
	agentsPerFloor  int
	ticksPerDay     int
	missingAgents   map[int][]int // key: floor, value: types of missing agents
	reshufflePeriod int
	tickCounter     int
}

func New(currPlatFood float64, currPlatFloor int, agentCount,
	agentsPerFloor, ticksPerDay, reshufflePeriod int) *Tower {
	t := &Tower{
		currPlatFood:    currPlatFood,
		maxPlatFood:     currPlatFood,
		currPlatFloor:   currPlatFloor,
		agentCount:      agentCount,
		agentsPerFloor:  agentsPerFloor,
		ticksPerDay:     ticksPerDay,
		missingAgents:   make(map[int][]int),
		reshufflePeriod: reshufflePeriod,
		tickCounter:     1,
	}
	t.initAgents()
	return t
}

func (t *Tower) Tick() {
	t.mx.RLock()
	defer t.mx.RUnlock()

	//logs
	log.Printf("A log from the tower! Tick no: %d", t.tickCounter)
	log.Printf("The food left on the platform = %f", t.currPlatFood)

	//useful parameters
	day := 24 * 60
	numOfFloors := t.agentCount / t.agentsPerFloor
	platformMovePeriod := day / numOfFloors // can add min/max

	// Shuffle the agents
	if t.tickCounter%t.reshufflePeriod == 0 {
		t.mx.RUnlock()
		defer t.mx.RLock()
		t.reshuffle(numOfFloors)
	}
	// Move the platform
	if (t.tickCounter)%(platformMovePeriod) == 0 {
		t.currPlatFloor++
	}
	// Decrease agent HP and reset tower at end of day
	if (t.tickCounter)%(day) == 0 {
		t.mx.RUnlock()
		t.hpDecay() // decreases HP and kills if < 0
		t.mx.RLock()
		t.ResetTower()
	}

	t.tickCounter++
}

func (t *Tower) GetMissingAgents() map[int][]int {
	deadAgents := t.missingAgents
	t.missingAgents = make(map[int][]int)
	return deadAgents
}
