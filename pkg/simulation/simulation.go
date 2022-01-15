package simulation

import (
	"context"

	"github.com/SOMAS2021/SOMAS2021/pkg/agents/randomAgent"
	"github.com/SOMAS2021/SOMAS2021/pkg/agents/team1/agent1"
	"github.com/SOMAS2021/SOMAS2021/pkg/agents/team1/agent2"
	"github.com/SOMAS2021/SOMAS2021/pkg/agents/team2"
	"github.com/SOMAS2021/SOMAS2021/pkg/agents/team3"

	team4EvoAgent "github.com/SOMAS2021/SOMAS2021/pkg/agents/team4/finalAgent"
	team5 "github.com/SOMAS2021/SOMAS2021/pkg/agents/team5"
	"github.com/SOMAS2021/SOMAS2021/pkg/agents/team6"
	team7agent1 "github.com/SOMAS2021/SOMAS2021/pkg/agents/team7/agent1"
	"github.com/SOMAS2021/SOMAS2021/pkg/config"
	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/agent"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/day"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/health"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/world"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/logging"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/utilFunctions"
	log "github.com/sirupsen/logrus"
)

type Fields = log.Fields
type AgentNewFunc func(base *infra.Base) (infra.Agent, error)

// var parameters *config.ConfigParameters
// var numAgents int = parameters.

// type State struct {
// 	ID            []string
// 	SM            []string
// 	FoodAvailable []string
// 	Utility       []string
// }

type SimEnv struct {
	FoodOnPlatform   food.FoodType
	AgentCount       map[agent.AgentType]int
	AgentHP          int
	AgentsPerFloor   int
	logger           log.Entry
	dayInfo          *day.DayInfo
	healthInfo       *health.HealthInfo
	world            world.World
	stateLog         *logging.StateLog
	agentNewFuncs    map[agent.AgentType]AgentNewFunc
	utilityCSVHeader [][]string
	deaths           int
	cumulativeDeaths int
}

func NewSimEnv(parameters *config.ConfigParameters, healthInfo *health.HealthInfo) *SimEnv {
	stateLog := logging.NewLogState(parameters.LogFolderName, parameters.LogMain, parameters.LogStory, parameters.CustomLog)
	// numAgents := parameters.NumOfAgents[agent.RandomAgent]
	return &SimEnv{
		FoodOnPlatform: parameters.FoodOnPlatform,
		AgentCount:     parameters.NumOfAgents,
		AgentHP:        parameters.AgentHP,
		dayInfo:        parameters.DayInfo,
		healthInfo:     healthInfo,
		AgentsPerFloor: parameters.AgentsPerFloor,
		logger:         *stateLog.Logmanager.GetLogger("main").WithFields(log.Fields{"reporter": "simulation"}),
		stateLog:       stateLog,
		agentNewFuncs: map[agent.AgentType]AgentNewFunc{
			agent.Team1Agent1: agent1.New,
			agent.Team1Agent2: agent2.New,
			agent.Team2:       team2.New,
			agent.Team3:       team3.New,
			agent.Team4:       team4EvoAgent.New,
			agent.Team5:       team5.New,
			agent.Team6:       team6.New,
			agent.Team7:       team7agent1.New,
			agent.RandomAgent: randomAgent.New,
		},
		utilityCSVHeader: [][]string{},
		deaths:           0,
		cumulativeDeaths: 0,
	}
}

func (sE *SimEnv) Simulate(ctx context.Context, ch chan<- string) {
	sE.Log("Simulation Initializing")
	totalAgents := utilFunctions.Sum(sE.AgentCount)
	t := infra.NewTower(sE.FoodOnPlatform, totalAgents, sE.AgentsPerFloor, sE.dayInfo, sE.healthInfo, sE.stateLog)
	sE.SetWorld(t)

	sE.generateInitialAgents(t)

	// Store all agent id and teams in an array to use as headers for utility CSV file
	for id, agent := range t.Agents {
		idTeamPair := []string{id.String(), agent.BaseAgent().AgentType().String()}
		sE.utilityCSVHeader = append(sE.utilityCSVHeader, idTeamPair)
	}
	sE.utilityCSVHeader = transpose(sE.utilityCSVHeader)

	// Create CSV files
	sE.CreateBehaviourCtrData()
	sE.CreateBehaviourChangeCtrData()
	sE.CreateUtilityData()
	sE.CreateDeathData()
	sE.CreateStateData(totalAgents)

	sE.Log("Simulation Started")
	sE.simulationLoop(t, ctx)

	// Assuming everything here will never timeout
	sE.Log("Simulation Ended")
	sE.Log("Summary of dead agents", infra.Fields{"Agent Type and number that died": t.DeadAgents()})
	for agentType, count := range t.DeadAgents() {
		sE.Log("dead agents", infra.Fields{"agentType": agentType.String(), "count": count})
		sE.AddToDeathData(count)
	}

	sE.Log("Living agents at end of simulation")
	for agentID, agent := range t.Agents {
		agent := agent.BaseAgent()
		sE.Log("Agent survives till the end of the simulation", infra.Fields{"agentID": agentID, "agentType": agent.AgentType().String(), "agentAge": agent.Age()})
	}

	// custom loggers
	for _, agent := range t.Agents {
		if agent.BaseAgent().AgentType().String() == sE.stateLog.CustomLog {
			agent.CustomLogs()
		}
	}

	// Write data into respective CSV files
	sE.ExportCSV(sE.dayInfo.BehaviourCtrData, "csvFiles/socialMotives.csv")
	sE.ExportCSV(sE.dayInfo.BehaviourChangeCtrData, "csvFiles/socialMotivesChange.csv")
	sE.ExportCSV(sE.dayInfo.UtilityData, "csvFiles/utility.csv")
	sE.ExportCSV(sE.dayInfo.DeathData, "csvFiles/deaths.csv")

	// Export state data into respective CSV files
	sE.ExportCSV(sE.dayInfo.StateData.ID, "csvFiles/states/id.csv")
	sE.ExportCSV(sE.dayInfo.StateData.HP, "csvFiles/states/hp.csv")
	sE.ExportCSV(sE.dayInfo.StateData.SM, "csvFiles/states/sm.csv")
	sE.ExportCSV(sE.dayInfo.StateData.FoodAvailable, "csvFiles/states/foodAvailable.csv")
	sE.ExportCSV(sE.dayInfo.StateData.Utility, "csvFiles/states/utility.csv")

	// dispatch loggers
	sE.stateLog.SimEnd(sE.dayInfo)
	ch <- "Simulation Finished"
}

func (sE *SimEnv) setSMCtr(agent infra.Agent) {
	if agent.BaseAgent().AgentType().String() == "Team6" { // Team 6
		agentBehaviour := agent.Behaviour()
		if agentBehaviour != "Not Team 6" {
			sE.dayInfo.BehaviourCtr[agentBehaviour]++
		}
	}
}

func (sE *SimEnv) resetSMCtr() {
	sE.dayInfo.BehaviourCtr = map[string]int{
		"Altruist":     0,
		"Collectivist": 0,
		"Selfish":      0,
		"Narcissist":   0,
	}
}

func (sE *SimEnv) setSMChangeCtr(agent infra.Agent) {
	if agent.BaseAgent().AgentType().String() == "Team6" { // Team 6
		agentBehaviourChange := agent.BehaviourChange()
		if agentBehaviourChange != "Not Team 6" {
			sE.dayInfo.BehaviourChangeCtr[agentBehaviourChange]++
		}
	}
}

func (sE *SimEnv) resetSMChangeCtr() {
	sE.dayInfo.BehaviourChangeCtr = map[string]int{
		"A2A": 0,
		"A2C": 0,
		"A2S": 0,
		"A2N": 0,
		"C2A": 0,
		"C2C": 0,
		"C2S": 0,
		"C2N": 0,
		"S2A": 0,
		"S2C": 0,
		"S2S": 0,
		"S2N": 0,
		"N2A": 0,
		"N2C": 0,
		"N2S": 0,
		"N2N": 0,
	}
}

func (sE *SimEnv) addToUtilityCtr(agent infra.Agent) {
	sE.dayInfo.Utility += agent.BaseAgent().Utility()
}

func (sE *SimEnv) clearUtilityCtr() {
	sE.dayInfo.Utility = 0.0
}

func (s *SimEnv) Log(message string, fields ...Fields) {
	if len(fields) == 0 {
		fields = append(fields, Fields{})
	}
	s.logger.WithFields(fields[0]).Info(message)
}

func transpose(slice [][]string) [][]string {
	xl := len(slice[0])
	yl := len(slice)
	result := make([][]string, xl)
	for i := range result {
		result[i] = make([]string, yl)
	}
	for i := 0; i < xl; i++ {
		for j := 0; j < yl; j++ {
			result[i][j] = slice[j][i]
		}
	}
	return result
}
