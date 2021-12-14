package simulation

import (
	"sync"

	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/abm"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/agent"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/world"
)

func (sE *SimEnv) SetWorld(w world.World) {
	if sE.world == nil {
		sE.world = w
	}
}

func (sE *SimEnv) World() world.World {
	return sE.world
}

func (sE *SimEnv) simulationLoop(a *abm.ABM, t *infra.Tower) {
	for sE.dayInfo.CurrTick <= sE.dayInfo.TotalTicks {
		sE.replaceAgents(t)
		// a.SimulationIterate(sE.dayInfo.CurrTick)
		sE.TowerTick()
		sE.AgentsRun(sE.dayInfo.CurrTick, t)
		if sE.reportFunc != nil {
			sE.reportFunc(sE)
		}
		sE.dayInfo.CurrTick++
	}
}

func (sE *SimEnv) SetReportFunc(fn func(*SimEnv)) {
	sE.reportFunc = fn
}

func (sE *SimEnv) TowerTick() {
	if sE.World() != nil {
		sE.World().Tick()
	}
}

func (sE *SimEnv) AgentsRun(i int, t *infra.Tower) {
	var wg sync.WaitGroup
	// first need to request tower agents.
	for _, agentPointer := range sE.custAgents {
		agentDeref := *agentPointer
		wg.Add(1)
		go func(wg *sync.WaitGroup, i int, agentDeref agent.Agent) {
			agentDeref.Run()
			wg.Done()
		}(&wg, i, agentDeref)
	}
	wg.Wait()
}
