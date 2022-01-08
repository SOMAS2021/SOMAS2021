package logging

import (
	"fmt"

	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/agent"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/day"
	log "github.com/sirupsen/logrus"
)

type StateLog struct {
	Logmanager *LogManager
	// Loggers
	foodLogger    *log.Logger
	deathLogger   *log.Logger
	storyLogger   *log.Logger
	mainLogger    *log.Logger
	utilityLogger *log.Logger
	// Death state
	deathCount int
	// Food state
	prevFood int
	// Custom log
	CustomLog string
}

type AgentState struct {
	HP        int
	AgentType agent.AgentType
	Floor     int
	Age       int
	Custom    string
	Utility   float64
}

func handleNewLoggerErr(err error) {
	if err != nil {
		fmt.Println("error creating new logger: ", err)
	}
}

func NewLogState(folderpath string, saveMainLog bool, saveStoryLog bool, customLog string) *StateLog {
	// init manager
	l := NewLogger(folderpath)

	// save main log
	mainLogName := "main.json"
	if !saveMainLog {
		mainLogName = ""
	}

	// save story log
	storyLogName := "story.json"
	if !saveStoryLog {
		storyLogName = ""
	}

	// new loggers
	foodLogger, err := l.AddLogger("food", "food.json")
	handleNewLoggerErr(err)
	deathLogger, err := l.AddLogger("death", "death.json")
	handleNewLoggerErr(err)
	storyLogger, err := l.AddLogger("story", storyLogName)
	handleNewLoggerErr(err)
	utilityLogger, err := l.AddLogger("utility", "utility.json")
	handleNewLoggerErr(err)
	utilityLogger, err := l.AddLogger("utility", "utility.json")
	handleNewLoggerErr(err)
	mainLogger, err := l.AddLogger("main", mainLogName)
	handleNewLoggerErr(err)

	return &StateLog{
		Logmanager:    &l,
		foodLogger:    foodLogger,
		deathLogger:   deathLogger,
		mainLogger:    mainLogger,
		storyLogger:   storyLogger,
		utilityLogger: utilityLogger,
		deathCount:    0,
		prevFood:      0,
		CustomLog:     customLog,
	}
}

func (ls *StateLog) LogAgentDeath(simState *day.DayInfo, agentType agent.AgentType, age int) {
	ls.deathCount++
	ls.deathLogger.
		WithFields(
			log.Fields{
				"day":              simState.CurrDay,
				"tick":             simState.CurrTick,
				"agent_type":       agentType.String(),
				"cumulativeDeaths": ls.deathCount,
				"ageUponDeath":     age,
			}).Info()
}

func (ls *StateLog) LogUtility(simState *day.DayInfo, agentType agent.AgentType, utility float64, isAlive bool) {
	ls.utilityLogger.
		WithFields(
			log.Fields{
				"day":        simState.CurrDay,
				"tick":       simState.CurrTick,
				"agent_type": agentType.String(),
				"utility":    utility,
				"isAlive":    isAlive,
			}).Info()
}

func (ls *StateLog) LogPlatFoodState(simState *day.DayInfo, food int) {
	if ls.prevFood != food {
		ls.foodLogger.
			WithFields(
				log.Fields{
					"day":  simState.CurrDay,
					"tick": simState.CurrTick,
					"food": food,
				}).Info()
		ls.prevFood = food
	}
}

// Story logging
func (state *AgentState) AgentFields() map[string]interface{} {
	return log.Fields{
		"hp":    state.HP,
		"atype": state.AgentType.String(),
		"age":   state.Age,
		"floor": state.Floor,
		"state": state.Custom,
	}
}

func (ls *StateLog) LogStoryAgentTookFood(simState *day.DayInfo, state AgentState, foodTaken int, foodLeft int) {
	ls.storyLogger.
		WithFields(
			log.Fields{
				"day":       simState.CurrDay,
				"tick":      simState.CurrTick,
				"foodTaken": foodTaken,
				"foodLeft":  foodLeft,
			}).
		WithFields(state.AgentFields()).Info("food")
}

func (ls *StateLog) LogStoryAgentSentMessage(simState *day.DayInfo, state AgentState, message messages.Message) {
	ls.storyLogger.
		WithFields(
			log.Fields{
				"day":      simState.CurrDay,
				"tick":     simState.CurrTick,
				"target":   message.TargetFloor(),
				"mtype":    message.MessageType().String(),
				"mcontent": message.StoryLog(),
			},
		).
		WithFields(state.AgentFields()).Info("message")
}

func (ls *StateLog) LogStoryAgentDied(simState *day.DayInfo, state AgentState) {
	ls.storyLogger.
		WithFields(
			log.Fields{
				"day":  simState.CurrDay,
				"tick": simState.CurrTick,
			}).
		WithFields(state.AgentFields()).Info("death")
}

func (ls *StateLog) LogStoryPlatformMoved(simState *day.DayInfo, floor int) {
	ls.storyLogger.
		WithFields(
			log.Fields{
				"day":   simState.CurrDay,
				"tick":  simState.CurrTick,
				"floor": floor,
			}).Info("platform")
}
