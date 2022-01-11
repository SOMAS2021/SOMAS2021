package simulation

import (
	"context"
	"time"

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

type SimEnv struct {
	FoodOnPlatform food.FoodType
	AgentCount     map[agent.AgentType]int
	AgentHP        int
	AgentsPerFloor int
	logger         log.Entry
	dayInfo        *day.DayInfo
	healthInfo     *health.HealthInfo
	world          world.World
	stateLog       *logging.StateLog
	agentNewFuncs  map[agent.AgentType]AgentNewFunc
	parameters     *config.ConfigParameters
}

func NewSimEnv(parameters *config.ConfigParameters, healthInfo *health.HealthInfo) *SimEnv {
	stateLog := logging.NewLogState(parameters.LogFolderName, parameters.LogMain, parameters.LogStory, parameters.CustomLog)
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
		parameters: parameters,
	}
}

func (sE *SimEnv) Simulate() {
	sE.Log("Simulation Initializing")

	totalAgents := utilFunctions.Sum(sE.AgentCount)
	t := infra.NewTower(sE.FoodOnPlatform, totalAgents, sE.AgentsPerFloor, sE.dayInfo, sE.healthInfo, sE.stateLog)
	sE.SetWorld(t)

	sE.generateInitialAgents(t)

	sE.Log("Simulation Started")

	//context and channel are passed all the way down to simulation loop to check every tick if there was a timeout
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch := make(chan string, 1)

	go func() {
		log.Info("Simulation started")
		sE.simulationLoop(t, ctx)
	}()

	// Listen on our channel AND a timeout channel - which ever happens first.
	select {
	case <-ch:
		log.Info("Simulation Finished Successfully")
		err := logging.UpdateSimStatusJson(sE.parameters.LogFolderName, "finished", GetMaxTick(sE.parameters.LogFolderName+"/story.json"))
		if err != nil {
			log.Fatal("Unable to update status file: " + err.Error())
		}
	case <-time.After(time.Duration(sE.parameters.SimTimeoutSeconds) * time.Second):
		err := logging.UpdateSimStatusJson(sE.parameters.LogFolderName, "timedout", GetMaxTick(sE.parameters.LogFolderName+"/story.json"))
		if err != nil {
			log.Fatal("Unable to update status file: " + err.Error())
		}
		log.Error("Simulation Timeout")
	}

	// Assuming everything here will never timeout
	sE.Log("Simulation Ended")
	sE.Log("Summary of dead agents", infra.Fields{"Agent Type and number that died": t.DeadAgents()})
	for agentType, count := range t.DeadAgents() {
		sE.Log("dead agents", infra.Fields{"agentType": agentType.String(), "count": count})
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

	// dispatch loggers
	sE.stateLog.SimEnd(sE.dayInfo)
	ch <- "Simulation Finished"
}

func (s *SimEnv) Log(message string, fields ...Fields) {
	if len(fields) == 0 {
		fields = append(fields, Fields{})
	}
	s.logger.WithFields(fields[0]).Info(message)
}
