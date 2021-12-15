package simulation

import (
	"github.com/SOMAS2021/SOMAS2021/pkg/agents/team1/agent1"
	"github.com/SOMAS2021/SOMAS2021/pkg/agents/team1/agent2"
	"github.com/SOMAS2021/SOMAS2021/pkg/agents/team6"
	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func (sE *SimEnv) createNewAgent(tower *infra.Tower, i, floor int) {
	// TODO: clean this looping, make a nice abs map
	sE.Log("Creating new agent")
	abs := []AgentNewFunc{agent1.New, agent2.New, team6.New}

	uuid := uuid.New().String()
	bagent, err := infra.NewBaseAgent(sE.world, i, sE.AgentHP, floor, uuid)
	if err != nil {
		log.Fatal(err)
	}

	custagent, err := abs[i](bagent)
	if err != nil {
		log.Fatal(err)
	}

	sE.custAgents[uuid] = custagent
	tower.AddAgent(bagent)
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
	sE.mx.RLock()
	defer sE.mx.RUnlock()
	return t.TotalAgents()
}
