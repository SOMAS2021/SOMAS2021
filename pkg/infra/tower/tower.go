package tower

import (
	"log"
	"sync"
)

type Tower struct {
	FoodOnPlatform  float64
	FloorOfPlatform uint64
	mx              sync.RWMutex
	AgentCount      int
	agents          map[string]BaseAgentCore
	AgentsPerFloor  int
	ticksPerDay     int
	missingAgents   map[int][]int // key: floor, value: types of missing agents
	reshufflePeriod int
}

func New(foodOnPlatform float64, floorOfPlatform uint64, agentCount, agentsPerFloor, ticksPerDay, reshufflePeriod int) *Tower {
	t := &Tower{
		FoodOnPlatform:  foodOnPlatform,
		FloorOfPlatform: floorOfPlatform,
		AgentCount:      agentCount,
		AgentsPerFloor:  agentsPerFloor,
		ticksPerDay:     ticksPerDay,
		missingAgents:   make(map[int][]int),
		reshufflePeriod: reshufflePeriod,
	}
	t.initAgents()
	return t
}

var tickCounter = 1

func (t *Tower) Tick() {
	t.mx.RLock()
	defer t.mx.RUnlock()
	log.Printf("A log from the tower! Tick no: %d", tickCounter)

	if (tickCounter)%(t.reshufflePeriod) == 0 {
		t.reshuffle(t.AgentsPerFloor)
	}

	//TODO: Need to agree on ticks per day so that hpDecay is updated once per day
	t.hpDecay() // deacreases HP and kills if < 0

	tickCounter += 1
}

func (t *Tower) GetMissingAgents() map[int][]int {
	deadAgents := t.missingAgents
	t.missingAgents = make(map[int][]int)
	log.Printf("%v", deadAgents)
	return deadAgents
}
