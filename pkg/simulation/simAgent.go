package simulation

import (
<<<<<<< HEAD
=======
	"github.com/SOMAS2021/SOMAS2021/pkg/agents/randomAgent"
	"github.com/SOMAS2021/SOMAS2021/pkg/agents/team1/agent1"
	"github.com/SOMAS2021/SOMAS2021/pkg/agents/team1/agent2"
	"github.com/SOMAS2021/SOMAS2021/pkg/agents/team3"
	// team4EvoAgent "github.com/SOMAS2021/SOMAS2021/pkg/agents/team4/agent1"
	agentTrust "github.com/SOMAS2021/SOMAS2021/pkg/agents/team4/agent2"
	team5 "github.com/SOMAS2021/SOMAS2021/pkg/agents/team5/agent1"
	"github.com/SOMAS2021/SOMAS2021/pkg/agents/team6"
	team7agent1 "github.com/SOMAS2021/SOMAS2021/pkg/agents/team7/agent1"
>>>>>>> 15ed57e... FEAT: Added Handling and Setting up message passing
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

<<<<<<< HEAD
func (sE *SimEnv) createNewAgent(tower *infra.Tower, agentType agent.AgentType, floor int) {
	sE.Log("Creating new agent", Fields{"type": agentType.String()})

=======
func (sE *SimEnv) createNewAgent(tower *infra.Tower, i, floor int) {
	sE.Log("Creating new agent")
	abs := []AgentNewFunc{agent1.New, agent2.New, team3.New, agentTrust.New, team5.New, team6.New, team7agent1.New, randomAgent.New}
	// NOTE(woonmoon): Leaving the line below commented just in case any teams want to run the 2-agent
	// 				   configuration to see how the message-passing works.
	// abs := []AgentNewFunc{agent1.New, agent2.New}
>>>>>>> 15ed57e... FEAT: Added Handling and Setting up message passing
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
