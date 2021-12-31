package logging

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

type StateLog struct {
	logmanager *LogManager
	// Loggers
	foodLogger  *log.Logger
	deathLogger *log.Logger
	// Death state
	deathCount int
}

func handleNewLoggerErr(err error) {
	if err != nil {
		fmt.Printf(err.Error())
	}
}

func NewLogState(folderpath string) *StateLog {
	// Add loggers we want for the state
	l := NewLogger(folderpath)

	// new loggers
	foodLogger, err := l.AddLogger("food", "food.json")
	handleNewLoggerErr(err)
	deathLogger, err := l.AddLogger("death", "death.json")
	handleNewLoggerErr(err)

	return &StateLog{
		logmanager:  &l,
		foodLogger:  foodLogger,
		deathLogger: deathLogger,
		deathCount:  0,
	}
}

func (ls *StateLog) LogAgentDeath(day int, tick int, agentType int) {
	ls.deathCount++
	ls.deathLogger.WithFields(log.Fields{"day": day, "tick": tick, "agent_type": agentType, "cumulativeDeaths": ls.deathCount}).Info("")
}

func (ls *StateLog) LogPlatFoodState(day int, tick int, food int) {
	ls.foodLogger.WithFields(log.Fields{"day": day, "tick": tick, "food": food}).Info("")
}
