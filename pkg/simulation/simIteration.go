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
		sE.Log("", Fields{"Current Simulation Tick": sE.dayInfo.CurrTick})
		t.TowerStateLog(" start of tick")
		sE.replaceAgents(t)
		sE.AgentsRun()
		sE.TowerTick()
		sE.dayInfo.CurrTick++
	}
}

func (sE *SimEnv) TowerTick() {
	if sE.World() != nil {
		sE.World().Tick()
	}
}

func (sE *SimEnv) AgentsRun() {
	agentsToRemove := make([]string, 0)
	var wg sync.WaitGroup
	for uuid, custAgent := range sE.custAgents {
		wg.Add(1)
		go func(wg *sync.WaitGroup, custAgent agent.Agent, uuid string) {
			if custAgent.IsAlive() {
				custAgent.Run()
			} else {
				agentsToRemove = append(agentsToRemove, uuid)
			}
			wg.Done()
		}(&wg, custAgent, uuid)
	}
	wg.Wait()
	sE.mx.Lock()
	for agentUUID := range agentsToRemove {
		delete(sE.custAgents, agentsToRemove[agentUUID])
	}
	sE.mx.Unlock()
}
