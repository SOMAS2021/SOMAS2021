package simulation

import (
	"sync"

	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
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
		sE.Log("", Fields{"Current Simulation Tick": sE.dayInfo.CurrTick})
		t.TowerStateLog(" start of tick")
		sE.replaceAgents(t)
		sE.AgentsRun(t)
		sE.TowerTick()
		sE.dayInfo.CurrTick++
	}
}

func (sE *SimEnv) TowerTick() {
	if sE.World() != nil {
		sE.World().Tick()
	}
}

func (sE *SimEnv) AgentsRun(t *infra.Tower) {
	var wg sync.WaitGroup
	for uuid, custAgent := range t.Agents {
		wg.Add(1)
		go func(wg *sync.WaitGroup, custAgent infra.Agent, uuid string) {
			if custAgent.IsAlive() {
				custAgent.Run()
			}
			wg.Done()
		}(&wg, custAgent, uuid)
	}
	wg.Wait()
}
