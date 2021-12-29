package simulation

import (
	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/agent"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func (sE *SimEnv) generateInitialAgents(t *infra.Tower) {
	floor := 1
	for agentType, agentCount := range sE.AgentCount {
		for i := 0; i < agentCount; i++ {
			sE.createNewAgent(t, agentType, floor)
			floor++
		}
	}
	sE.Log("", Fields{"Number of new agents created: ": floor - 1})
}

func (sE *SimEnv) createNewAgent(tower *infra.Tower, agentType agent.AgentType, floor int) {
	sE.Log("Creating new agent", Fields{"type": agentType.String()})

	// NOTE(woonmoon): Leaving the line below commented just in case any teams want to run the 2-agent
	// 				   configuration to see how the message-passing works.
	// abs := []AgentNewFunc{agent1.New, agent2.New}

	uuid := uuid.New()

	bAgent, err := infra.NewBaseAgent(sE.world, agentType, sE.AgentHP, floor, uuid)
	if err != nil {
		log.Fatal(err)
	}

	custAgent, err := sE.agentNewFuncs[agentType](bAgent)
	if err != nil {
		log.Fatal(err)
	}
	tower.AddAgent(custAgent)
}

func (sE *SimEnv) replaceAgents(t *infra.Tower) {
	for uuid, agent := range t.Agents {
		agent := agent.BaseAgent()
		if !agent.IsAlive() {
			delete(t.Agents, uuid)
			t.UpdateDeadAgents(agent.AgentType())
			sE.createNewAgent(t, agent.AgentType(), agent.Floor())
		}
	}
}

func (sE *SimEnv) AgentsCount(t *infra.Tower) int {
	return t.TotalAgents()
}
