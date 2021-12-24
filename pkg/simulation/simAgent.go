package simulation

import (
	"github.com/SOMAS2021/SOMAS2021/pkg/agents/randomAgent"
	"github.com/SOMAS2021/SOMAS2021/pkg/agents/team1/agent1"
	"github.com/SOMAS2021/SOMAS2021/pkg/agents/team1/agent2"
	"github.com/SOMAS2021/SOMAS2021/pkg/agents/team3"
	"github.com/SOMAS2021/SOMAS2021/pkg/agents/team6"
	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/world"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type AgentNewFunc func(world world.World, agentType int, agentHP int, agentFloor int, id string) (infra.Agent, error)

func (sE *SimEnv) generateInitialAgents(t *infra.Tower) {
	agentIndex := 1
	for i := 0; i < len(sE.AgentCount); i++ {
		for j := 0; j < sE.AgentCount[i]; j++ {
			sE.createNewAgent(t, i, agentIndex)
			agentIndex++
		}
	}
	sE.Log("", Fields{"Number of new agents created: ": agentIndex - 1})
}

func (sE *SimEnv) createNewAgent(tower *infra.Tower, i, floor int) {
	sE.Log("Creating new agent")
	abs := []AgentNewFunc{agent1.New, agent2.New, team3.New, team6.New, randomAgent.New}
	uuid := uuid.New().String()

	custAgent, err := abs[i](sE.world, i, sE.AgentHP, floor, uuid)
	if err != nil {
		log.Fatal(err)
	}
	tower.AddAgent(custAgent)
}

func (sE *SimEnv) replaceAgents(t *infra.Tower) {
	missingAgentsMap := t.UpdateMissingAgents()
	for floor := range missingAgentsMap {
		for _, agentType := range missingAgentsMap[floor] {
			sE.createNewAgent(t, agentType, floor)
		}
	}
}

func (sE *SimEnv) AgentsCount(t *infra.Tower) int {
	return t.TotalAgents()
}
