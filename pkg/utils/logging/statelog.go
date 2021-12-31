package logging

import (
	"fmt"

	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/agent"
	log "github.com/sirupsen/logrus"
)

type StateLog struct {
	Logmanager *LogManager
	// Loggers
	foodLogger  *log.Logger
	deathLogger *log.Logger
	mainLogger  *log.Logger
	// Death state
	deathCount int
}

func handleNewLoggerErr(err error) {
	if err != nil {
		fmt.Println("error creating new logger: ", err)
	}
}

func NewLogState(folderpath string, saveMainLog bool) *StateLog {
	// init manager
	l := NewLogger(folderpath)

	// save main log
	mainLogName := "main.json"
	if !saveMainLog {
		mainLogName = ""
	}

	// new loggers
	foodLogger, err := l.AddLogger("food", "food.json")
	handleNewLoggerErr(err)
	deathLogger, err := l.AddLogger("death", "death.json")
	handleNewLoggerErr(err)
	mainLogger, err := l.AddLogger("main", mainLogName)
	handleNewLoggerErr(err)

	return &StateLog{
		Logmanager:  &l,
		foodLogger:  foodLogger,
		deathLogger: deathLogger,
		mainLogger:  mainLogger,
		deathCount:  0,
	}
}

func (ls *StateLog) LogAgentDeath(day int, tick int, agentType agent.AgentType) {
	ls.deathCount++
	ls.deathLogger.WithFields(log.Fields{"day": day, "tick": tick, "agent_type": agentType.String(), "cumulativeDeaths": ls.deathCount}).Info("")
}

func (ls *StateLog) LogPlatFoodState(day int, tick int, food int) {
	ls.foodLogger.WithFields(log.Fields{"day": day, "tick": tick, "food": food}).Info("")
}
