package logging

import (
	"fmt"

	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/agent"
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
		fmt.Println("error creating new logger: ", err)
	}
}

func NewLogState(folderpath string) *StateLog {
	// init manager
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

func (ls *StateLog) LogAgentDeath(day int, tick int, agentType agent.AgentType) {
	ls.deathCount++
	ls.deathLogger.WithFields(log.Fields{"day": day, "tick": tick, "agent_type": agentType.String(), "cumulativeDeaths": ls.deathCount}).Info("")
}

func (ls *StateLog) LogPlatFoodState(day int, tick int, food int) {
	ls.foodLogger.WithFields(log.Fields{"day": day, "tick": tick, "food": food}).Info("")
}
