package infra

import (
	"log"
	"math/rand"
)

type Tower struct {
	currPlatFood    float64
	maxPlatFood     float64
	currPlatFloor   uint64
	agentCount      int
	agents          []*Base
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
		agents:          make([]*Base, 0),
		agentsPerFloor:  agentsPerFloor,
		ticksPerDay:     ticksPerDay,
		missingAgents:   make(map[int][]int),
		reshufflePeriod: reshufflePeriod,
		tickCounter:     1,
	}
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

func (t *Tower) AddAgent(bagent *Base) {
	t.agents = append(t.agents, bagent)
}

func (t *Tower) reshuffle(numOfFloors int) {
	remainingVacanies := make([]int, numOfFloors)
	log.Printf("Reshuffling alive agents...")
	log.Printf("Number of agents: %d", len(t.agents))
	for i := 0; i < numOfFloors; i++ { // adding a max to each floor
		remainingVacanies[i] = t.agentsPerFloor
	}
	// allocating agents to floors randomly
	// iterate through the uuid strings of each agent
	for _, agent := range t.agents {
		newFloor := rand.Intn(numOfFloors)
		for remainingVacanies[newFloor] == 0 {
			newFloor = rand.Intn(numOfFloors)
		}
		agent.setFloor(newFloor + 1)
		remainingVacanies[newFloor]--
	}
}

func (t *Tower) hpDecay() {
	// TODO: can add a parameter
	var killed []int
	for i, agent := range t.agents {
		newHP := agent.HP() - 20
		if newHP < 0 {
			agent.Die()
			killed = append(killed, i)
			log.Printf("Killing agent %s", agent.ID())
			t.missingAgents[agent.Floor()] = append(t.missingAgents[agent.Floor()], agent.agentType)
		} else {
			agent.setHP(newHP)
		}
	}
	if len(killed) > 0 {
		tmp := make([]*Base, 0)
		for i, agent := range t.agents {
			if !contains(killed, i) {
				tmp = append(tmp, agent)
			}
		}
		t.agents = tmp
	}
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func (t *Tower) ResetTower() {
	t.currPlatFood = t.maxPlatFood
	t.currPlatFloor = 1
}
