package simulation

import (
	"sync"

	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
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

func (sE *SimEnv) simulationLoop(t *infra.Tower) {
	for sE.dayInfo.CurrTick <= sE.dayInfo.TotalTicks {
		sE.replaceAgents(t)
		sE.TowerTick()
		sE.AgentsRun()
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

func (sE *SimEnv) AgentsRun() {
	var wg sync.WaitGroup
	for uuid, custAgent := range sE.custAgents {
		wg.Add(1)
		go func(wg *sync.WaitGroup, custAgent agent.Agent, uuid string) {
			if custAgent.IsAlive() {
				custAgent.Run()
			} else {
				delete(sE.custAgents, uuid)
			}
			wg.Done()
		}(&wg, custAgent, uuid)
	}
	wg.Wait()
}
