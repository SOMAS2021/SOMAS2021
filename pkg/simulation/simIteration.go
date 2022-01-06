package simulation

import (
	"context"
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

// For things that need to be done once per tick, before agents take action
func (sE *SimEnv) Preface(t *infra.Tower) {
	// Update food seen metric only once for alive agents
	for _, custAgent := range t.Agents {
		agent := custAgent.BaseAgent()
		platFood := agent.CurrPlatFood()
		if agent.IsAlive() && platFood != -1 {
			agent.UpdateFoodSeen(platFood)
		}
	}
}

func (sE *SimEnv) simulationLoop(t *infra.Tower, ctx context.Context, ch chan<- string) {
	t.Reshuffle()
	for sE.dayInfo.CurrTick <= sE.dayInfo.TotalTicks {
		sE.Log("", Fields{"Current Simulation Tick": sE.dayInfo.CurrTick})
		t.TowerStateLog(" start of tick")
		sE.Preface(t)
		sE.AgentsRun(t)
		sE.TowerTick() //added reincarnation feature
		//sE.replaceAgents(t) //commented out by team2 for testing purpose
		// t.TowerStateLog(" end of tick")
		sE.dayInfo.CurrTick++

		//returns if there is a timeout
		//continously checks every tick since there is no way to kill a goroutine
		select {
		case <-ctx.Done():
			return
		default:
		}
	}
	ch <- "Simulation Finished"
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
