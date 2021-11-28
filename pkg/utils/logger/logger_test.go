package logger

import (
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestLogger(t *testing.T) {
	var myLogger Logger
	myLogger.SetOutputFile("log.txt")
	myLogger.AddLogger("AGENT", "SELFISH", "TEAM 6")
	AgentLog(myLogger.GetLogger("AGENT", "SELFISH", "TEAM 6"))
}

func AgentLog(logger *log.Logger){
	logger.Info("An agent log")
	logger.Info("Another agent log")
}