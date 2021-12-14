package abm

import (
	"sort"
	"sync"

	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/agent"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/world"
)

type ABM struct {
	mx         sync.RWMutex
	agents     []agent.Agent
	world      world.World
	reportFunc func(*ABM)
}

// New creates new ABM simulation engine with default
// parameters.
func New() *ABM {
	return &ABM{}
}

func (a *ABM) SetWorld(w world.World) {
	if a.world == nil {
		a.world = w
	}
}

func (a *ABM) World() world.World {
	return a.world
}

func (a *ABM) SetReportFunc(fn func(*ABM)) {
	a.reportFunc = fn
}

func (a *ABM) AddAgent(agent agent.Agent) {
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
			a.agents = []agent.Agent{}
		} else {
			a.agents = a.agents[1:]
		}
	} else if index == len(a.agents)-1 {
		a.agents = a.agents[:len(a.agents)-1]
	} else {
		a.agents = append(a.agents[:index], a.agents[index+1:]...)
	}
}

func (a *ABM) SimulationIterate(i int) {
	agentsToRemove := []int{}
	if a.World() != nil {
		a.World().Tick()
	}
	var wg sync.WaitGroup
	for j := 0; j < a.AgentsCount(); j++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup, i, j int) {
			if !a.agents[j].IsAlive() {
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
