package abm

import (
	"log"
	"sort"
	"sync"
)

type ABM struct {
	mx     sync.RWMutex
	agents []Agent

	currIteration int
	limit         int

	world      World
	reportFunc func(*ABM)
}

// New creates new ABM simulation engine with default
// parameters.
func New() *ABM {
	return &ABM{
		limit: 1000,
	}
}

func (a *ABM) SetWorld(w World) {
	if a.world == nil {
		a.world = w
	}
}

func (a *ABM) World() World {
	return a.world
}

func (a *ABM) SetReportFunc(fn func(*ABM)) {
	a.reportFunc = fn
}

func (a *ABM) AddAgent(agent Agent) {
	a.mx.Lock()
	defer a.mx.Unlock()
	a.agents = append(a.agents, agent)
}

// TODO: corner cases?
func (a *ABM) RemoveAgent(index int) {
	a.mx.Lock()
	defer a.mx.Unlock()
	if index == 0 {
		if len(a.agents) == 1 {
			a.agents = []Agent{}
		} else {
			a.agents = a.agents[1:]
		}
	} else if index == len(a.agents)-1 {
		a.agents = a.agents[:len(a.agents)-1]
	} else {
		a.agents = append(a.agents[:index], a.agents[index+1:]...)
	}
}

func (a *ABM) AddAgents(spawnFunc func(*ABM) Agent, n int) {
	for i := 0; i < n; i++ {
		agent := spawnFunc(a)
		a.AddAgent(agent)
	}
}

func (a *ABM) LimitIterations(n int) {
	a.mx.Lock()
	defer a.mx.Unlock()
	a.limit = n
}

func (a *ABM) Limit() int {
	a.mx.RLock()
	defer a.mx.RUnlock()
	return a.limit
}

// Iteration returns current iteration (age, generation).
func (a *ABM) Iteration() int {
	a.mx.RLock()
	defer a.mx.RUnlock()
	return a.currIteration
}

func (a *ABM) SimulationIterate(i int) {
	agentsToRemove := []int{}
	a.currIteration = i
	if a.World() != nil {
		a.World().Tick()
	}
	var wg sync.WaitGroup
	for j := 0; j < a.AgentsCount(); j++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup, i, j int) {
			if !a.agents[j].IsAlive() {
				log.Printf("Removing agent from ABM")
				agentsToRemove = append(agentsToRemove, j)
			} else {
				a.agents[j].Run()
			}
			wg.Done()
		}(&wg, i, j)
	}

	wg.Wait()
	sort.Ints(agentsToRemove)
	for index := len(agentsToRemove) - 1; index > -1; index-- {
		a.RemoveAgent(index)
		if index != 0 {
			agentsToRemove = agentsToRemove[:index]
		} else {
			agentsToRemove = []int{}
		}
	}
	if a.reportFunc != nil {
		a.reportFunc(a)
	}
}

func (a *ABM) AgentsCount() int {
	a.mx.RLock()
	defer a.mx.RUnlock()
	return len(a.agents)
}

func (a *ABM) Agents() []Agent {
	a.mx.RLock()
	defer a.mx.RUnlock()
	return a.agents
}

func (a *ABM) Count(condition func(agent Agent) bool) int {
	var count int
	for _, agent := range a.Agents() {
		if condition(agent) {
			count++
		}
	}
	return count
}
