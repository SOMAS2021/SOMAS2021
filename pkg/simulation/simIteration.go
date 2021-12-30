package simulation

import (
	"sync"

	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/world"
	"github.com/google/uuid"
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
	t.Reshuffle()
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
	for id, custAgent := range t.Agents {
		wg.Add(1)
		go func(wg *sync.WaitGroup, custAgent infra.Agent, id uuid.UUID) {
			if custAgent.IsAlive() {
				custAgent.Run()
			}
			wg.Done()
		}(&wg, custAgent, id)
	}
	wg.Wait()
}
