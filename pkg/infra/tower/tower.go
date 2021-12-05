package tower

/*
abm lib - controller -> world.tick & agent.run (for all agents)

world (tower) -> Tick() [death, reshuffle etc] & each agent info (hp and floor) & (abm.agent info - custom agents)

*/

import (
	"log"
	"math/rand"
	"strconv"
	"sync"

	"github.com/divan/goabm/abm"
)

type BaseAgentCore struct {
	hp    int
	floor int
	cust  abm.Agent
}

func (tower *Tower) GetHP(id string) int {
	return tower.agents[id].hp
}

func (tower *Tower) GetFloor(id string) int {
	return tower.agents[id].floor
}

func (tower *Tower) Exists(id string) bool {
	if _, found := tower.agents[id]; found {
		return true
	} else {
		return false
	}
}

func (tower *Tower) setFloor(id string, newFloor int) {
	temp := BaseAgentCore{
		hp:    tower.agents[id].hp,
		floor: newFloor,
		cust:  tower.agents[id].cust,
	}
	tower.agents[id] = temp
}

func (tower *Tower) setHP(id string, newHP int) {
	temp := BaseAgentCore{
		hp:    newHP,
		floor: tower.agents[id].floor,
		cust:  tower.agents[id].cust,
	}
	tower.agents[id] = temp
}

type Tower struct {
	FoodOnPlatform  float64
	FloorOfPlatform uint64
	mx              sync.RWMutex
	AgentCount      int
	agents          map[string]BaseAgentCore
	AgentsPerFloor  int
	ticksPerDay     int
	missingAgents   map[int]int // key: floor, value: missing agents from that floor
}

func New(foodOnPlatform float64, floorOfPlatform uint64, agentCount int, agentsPerFloor int, ticksPerDay int) *Tower {
	t := &Tower{
		FoodOnPlatform:  foodOnPlatform,
		FloorOfPlatform: floorOfPlatform,
		AgentCount:      agentCount,
		AgentsPerFloor:  agentsPerFloor,
		ticksPerDay:     ticksPerDay,
		missingAgents:   make(map[int]int),
	}

	t.initAgents()

	return t
}

func (t *Tower) initAgents() {
	t.agents = make(map[string]BaseAgentCore, t.AgentCount)
}

func (t *Tower) killAgent(id string) {
	// this removes the agent from the list of agents in the tower
	log.Printf("Killing agent %s", id)
	deadAgentFloor := t.agents[id].floor
	t.missingAgents[deadAgentFloor]++ // can just do this instead of checking if this is found (if not found it'll automatically initialise it to 0)
	delete(t.agents, id)              // delete the agent from the map of all agents

}

func (t *Tower) death() {
	log.Printf("Killing Agents...")
	for id := range t.agents {
		if t.GetHP(id) <= 0 {
			t.killAgent(id)
		}
	}
}

func (t *Tower) replace(agentsPerFloor int) {
	log.Printf("Replacing...")
	for floor := range t.missingAgents {
		// TODO: add agents to the floor
		delete(t.missingAgents, floor)
	}

	// doesn't kill people, that is a separate function
	// go through every floor, check how many agents and add them
	// if this isn't the same amount as agents per floor, add agents
}

func (t *Tower) reshuffle(agentsPerFloor int) {

	numOfFloors := t.AgentCount / int(agentsPerFloor)
	remainingVacanies := make([]int, numOfFloors)
	log.Printf("Reshuffling...")
	// adding a max to each floor
	for i := 0; i < numOfFloors; i++ {

		remainingVacanies[i] = agentsPerFloor
	}

	// allocating agents to floors randomly
	// iterate through the uuid strings of each agent
	for id := range t.agents {

		newFloor := rand.Intn(numOfFloors) // random number in the range 1 - numOfFloors
		for remainingVacanies[newFloor] == 0 {
			newFloor = rand.Intn(numOfFloors)
		}
		t.setFloor(id, newFloor+1) // only do this to agentsLocal cause agentsABM don;t know what floor they're on
		remainingVacanies[newFloor]--
	}
}

func (t *Tower) hpDecay() {
	for id := range t.agents {
		t.setHP(id, t.agents[id].hp-3)
	}
}

var tickCounter = 1

func (t *Tower) Tick() {
	t.mx.RLock()
	defer t.mx.RUnlock()
	log.Printf("A log from the tower! Tick no: " + strconv.Itoa(tickCounter))

	reshufflePeriod := 5
	replacePeriod := 1 // replace every tick

	if tickCounter%reshufflePeriod == 0 {
		t.reshuffle(t.AgentsPerFloor)
	}

	if tickCounter%replacePeriod == 0 {
		t.replace(t.AgentsPerFloor)
	}
	//Need to agree on ticks per day so that hpDecay is updated once per day
	t.hpDecay()
	t.death()

	tickCounter += 1
}

func (t *Tower) SetAgent(agentHP, agentFloor int, id string, customAgent abm.Agent) {
	t.mx.Lock()
	t.agents[id] = BaseAgentCore{ // creating a new instance of agent in hash map
		hp:    agentHP,
		floor: agentFloor,
		cust:  customAgent,
	}
	t.mx.Unlock()
}
