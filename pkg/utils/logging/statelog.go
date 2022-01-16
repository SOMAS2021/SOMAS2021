package logging

import (
	"fmt"
	"sync"

	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/agent"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/day"
	log "github.com/sirupsen/logrus"
)

type msgMap [12][8]int
type deathMap [8]int
type treatyResponsesCount [2][8]int // 0 for reject, 1 for accept

type StateLog struct {
	Logmanager *LogManager
	// Loggers
	foodLogger    *log.Logger
	deathLogger   *log.Logger
	storyLogger   *log.Logger
	mainLogger    *log.Logger
	utilityLogger *log.Logger
	msgLogger     *log.Logger
	// Death state
	deathCount int
	deaths     deathMap
	// Food state
	prevFood int
	// Messages state
	messages        *msgMap
	treatyResponses *treatyResponsesCount
	msgMx           *sync.Mutex
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
	msgLogger, err := l.AddLogger("messages", "messages.json")
	handleNewLoggerErr(err)
	storyLogger, err := l.AddLogger("story", storyLogName)
	handleNewLoggerErr(err)
	utilityLogger, err := l.AddLogger("utility", "utility.json")
	handleNewLoggerErr(err)
	mainLogger, err := l.AddLogger("main", mainLogName)
	handleNewLoggerErr(err)

	// Init message counters
	var msgs msgMap
	var treatyResponses treatyResponsesCount

	// Init death counter
	var deaths deathMap

	return &StateLog{
		Logmanager:      &l,
		foodLogger:      foodLogger,
		deathLogger:     deathLogger,
		mainLogger:      mainLogger,
		storyLogger:     storyLogger,
		utilityLogger:   utilityLogger,
		msgLogger:       msgLogger,
		messages:        &msgs,
		treatyResponses: &treatyResponses,
		msgMx:           &sync.Mutex{},
		deathCount:      0,
		deaths:          deaths,
		prevFood:        0,
		CustomLog:       customLog,
	}
}

// Death loggging
func (ls *StateLog) LogAgentDeath(simState *day.DayInfo, agentType agent.AgentType, age int) {
	ls.deathCount++
	temp := ls.deaths[agentType-1] + 1
	ls.deaths[agentType-1] = temp
	ls.deathLogger.
		WithFields(
			log.Fields{
				"day":              simState.CurrDay,
				"tick":             simState.CurrTick,
				"agent_type":       agentType.String(),
				"cumulativeDeaths": ls.deaths[agentType-1],
				"totalDeaths":      ls.deathCount,
				"ageUponDeath":     age,
			}).Info()
}

// Utility logging
func (ls *StateLog) LogUtility(simState *day.DayInfo, agentType agent.AgentType, utility float64, isAlive bool, floor int, floorCount int) {
	ls.utilityLogger.
		WithFields(
			log.Fields{
				"day":        simState.CurrDay,
				"tick":       simState.CurrTick,
				"agent_type": agentType.String(),
				"utility":    utility,
				"isAlive":    isAlive,
				"floor":      floor,
				"floorCount": floorCount,
			}).Info()
}

// Food logging
func (ls *StateLog) LogPlatFoodState(simState *day.DayInfo, food int, floor int, floorCount int) {
	if ls.prevFood != food {
		ls.foodLogger.
			WithFields(
				log.Fields{
					"day":        simState.CurrDay,
					"tick":       simState.CurrTick,
					"food":       food,
					"floor":      floor,
					"floorCount": floorCount,
				}).Info()
		ls.prevFood = food
	}
}

// Messages logging
func (ls *StateLog) LogMessage(simState *day.DayInfo, state AgentState, message messages.Message) {
	ls.msgMx.Lock()
	agentType := state.AgentType
	msgType := message.MessageType()
	temp := ls.messages[msgType-1][agentType-1] + 1
	ls.messages[msgType-1][agentType-1] = temp
	// Log treatyResponses accept/jerect. Tracks how many treatyResponses an agent rejects (0) and accepts (1)
	if msgType == messages.TreatyResponse {
		state := message.StoryLog()
		if state == "yes" {
			ls.treatyResponses[1][agentType-1] += 1
		} else {
			ls.treatyResponses[0][agentType-1] += 1
		}
	}
	ls.msgMx.Unlock()
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

// Simulation ended
func (ls *StateLog) SimEnd(simState *day.DayInfo) {
	// Dump messages map state
	var atypes [8]string
	for i := 0; i < 8; i++ {
		atypes[i] = agent.AgentType(i + 1).String()
	}

	var mtypes [12]string
	for i := 0; i < 12; i++ {
		mtypes[i] = messages.MessageType(i + 1).String()
	}

	ls.msgLogger.
		WithFields(
			log.Fields{
				"day":             simState.CurrDay,
				"tick":            simState.CurrTick,
				"msgcount":        ls.messages,
				"atypes":          atypes,
				"mtypes":          mtypes,
				"treatyResponses": ls.treatyResponses,
			},
		).Info()

}
