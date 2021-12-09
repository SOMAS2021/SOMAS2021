package infra

import (
	"log"
)

type Tower struct {
	currPlatFood    float64
	maxPlatFood     float64
	currPlatFloor   uint64
	agentCount      int
	agents          []Base
	agentsPerFloor  int
	ticksPerDay     int
	missingAgents   map[int][]int // key: floor, value: types of missing agents
	reshufflePeriod int
	tickCounter     int
}

func NewTower(currPlatFood float64, currPlatFloor uint64, agentCount,
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
	//logs
	log.Printf("A log from the tower! Tick no: %d", t.tickCounter)
	log.Printf("The food left on the platform = %f", t.currPlatFood)

	//useful parameters
	day := 24 * 60
	numOfFloors := t.agentCount / int(t.agentsPerFloor)
	platformMovePeriod := day / numOfFloors // can add min/max

	if (t.tickCounter)%(t.reshufflePeriod) == 0 {
		t.reshuffle(numOfFloors)
	}
	if (t.tickCounter)%(platformMovePeriod) == 0 {
		t.currPlatFloor++
	}
	if (t.tickCounter)%(day) == 0 {
		t.hpDecay() // deacreases HP and kills if < 0
		t.ResetTower()
	}

	t.tickCounter++
}

func (t *Tower) GetMissingAgents() map[int][]int {
	deadAgents := t.missingAgents
	t.missingAgents = make(map[int][]int)
	return deadAgents
}

func (t *Tower) SetAgent(aType, agentHP, agentFloor int, id string) {
	newAgent := Base{
		id:        id,
		hp:        agentHP,
		floor:     agentFloor,
		agentType: aType,
		tower:     t,
	}
	t.agents = append(t.agents, newAgent)
}
