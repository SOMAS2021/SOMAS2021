package simulation

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/world"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/utilFunctions"
	"github.com/google/uuid"
)

func (sE *SimEnv) SetWorld(w world.World) {
	if sE.Tower == nil {
		sE.Tower = w
	}
}

func (sE *SimEnv) World() world.World {
	return sE.Tower
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

func (sE *SimEnv) simulationLoop(t *infra.Tower, ctx context.Context) {
	t.Reshuffle()
	totalAgents := utilFunctions.Sum(sE.AgentCount)

	for sE.dayInfo.CurrTick <= sE.dayInfo.TotalTicks {
		sE.initStates(totalAgents)
		sE.Log("", Fields{"Current Simulation Tick": sE.dayInfo.CurrTick})
		t.TowerStateLog(" start of tick")

		sE.Preface(t)
		sE.AgentsRun(t)
		sE.TowerTick()
		sE.replaceAgents(t)
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
}

func (sE *SimEnv) initStates(totalAgents int) {
	sE.dayInfo.State.ID = []string{}
	sE.dayInfo.State.HP = []string{}
	sE.dayInfo.State.SM = []string{}
	sE.dayInfo.State.FoodAvailable = []string{}
	sE.dayInfo.State.Utility = []string{}

	var i int
	for i < totalAgents {
		sE.dayInfo.State.ID = append(sE.dayInfo.State.ID, strconv.Itoa(0))
		sE.dayInfo.State.HP = append(sE.dayInfo.State.HP, strconv.Itoa(0))
		sE.dayInfo.State.SM = append(sE.dayInfo.State.SM, strconv.Itoa(0))
		sE.dayInfo.State.FoodAvailable = append(sE.dayInfo.State.FoodAvailable, strconv.Itoa(0))
		sE.dayInfo.State.Utility = append(sE.dayInfo.State.Utility, strconv.Itoa(0))
		i++
	}
}

func (sE *SimEnv) TowerTick() {
	if sE.World() != nil {
		sE.World().Tick()
	}
}

func (sE *SimEnv) AgentsRun(t *infra.Tower) {
	var wg sync.WaitGroup
	if sE.dayInfo.CurrTick%sE.dayInfo.TicksPerDay == 0 {
		sE.resetSMCtr()
		sE.resetSMChangeCtr()
	}
	// For each agent
	for id, custAgent := range t.Agents {
		wg.Add(1)
		// To do once a day before agent runs
		// Log the ctr for agent social motives
		if sE.dayInfo.CurrTick%sE.dayInfo.TicksPerDay == 0 {
			sE.setSMCtr(custAgent)
			sE.AddToState(custAgent)
		}

		go func(wg *sync.WaitGroup, custAgent infra.Agent, id uuid.UUID) {
			if custAgent.IsAlive() {
				custAgent.Run()
			}
			wg.Done()
		}(&wg, custAgent, id)
		// To do once a day after agent runs
		if sE.dayInfo.CurrTick%sE.dayInfo.TicksPerDay == 0 {
			sE.setSMChangeCtr(custAgent)
			sE.addToUtilityCtr(custAgent)
		}
	}
	wg.Wait()
	// To do once a day after agents finish running
	if sE.dayInfo.CurrTick%sE.dayInfo.TicksPerDay == 0 {
		sE.stateLog.LogSocialMotivesCtr(sE.dayInfo)
		sE.AddToBehaviourCtrData()
		sE.AddToBehaviourChangeCtrData()
		sE.AddToUtilityData(t)
		sE.AddToDeathData(t.DeathCount)
		sE.clearUtilityCtr()
		sE.AddToStateData()
	}
}

func (sE *SimEnv) AddToState(agent infra.Agent) {
	if agent.BaseAgent().AgentType().String() != "Team6" {
		return
	}
	floorIdx := agent.BaseAgent().Floor() - 1
	// Agent ID
	sE.dayInfo.State.ID[floorIdx] = agent.BaseAgent().ID().String()
	// Agent HP
	sE.dayInfo.State.HP[floorIdx] = strconv.Itoa(agent.BaseAgent().HP())
	// Agent SM
	sE.dayInfo.State.SM[floorIdx] = agent.Behaviour()
	// Agent Food Available
	sE.dayInfo.State.FoodAvailable[floorIdx] = agent.FoodReceived()
	// Agent Utility
	sE.dayInfo.State.Utility[floorIdx] = fmt.Sprintf("%f", agent.BaseAgent().Utility())
}
