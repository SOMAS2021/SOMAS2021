package simulation

import (
	"github.com/SOMAS2021/SOMAS2021/pkg/agents/randomAgent"
	"github.com/SOMAS2021/SOMAS2021/pkg/agents/team1/agent1"
	"github.com/SOMAS2021/SOMAS2021/pkg/agents/team1/agent2"
	"github.com/SOMAS2021/SOMAS2021/pkg/agents/team3"
	team4EvoAgent "github.com/SOMAS2021/SOMAS2021/pkg/agents/team4/agent1"
	team5 "github.com/SOMAS2021/SOMAS2021/pkg/agents/team5/agent1"
	"github.com/SOMAS2021/SOMAS2021/pkg/agents/team6"
	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type AgentNewFunc func(base *infra.Base) (infra.Agent, error)

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
	abs := []AgentNewFunc{agent1.New, agent2.New, team3.New, team4EvoAgent.New, team5.New, team6.New, randomAgent.New}
	// NOTE(woonmoon): Leaving the line below commented just in case any teams want to run the 2-agent
	// 				   configuration to see how the message-passing works.
	// abs := []AgentNewFunc{agent1.New, agent2.New}
	uuid := uuid.New().String()

	bAgent, err := infra.NewBaseAgent(sE.world, i, sE.AgentHP, floor, uuid)
	if err != nil {
		log.Fatal(err)
	}

	custAgent, err := abs[i](bAgent)
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
			sE.createNewAgent(t, agent.AgentType(), agent.Floor())
		}
	}
}

func (sE *SimEnv) AgentsCount(t *infra.Tower) int {
	return t.TotalAgents()
}
