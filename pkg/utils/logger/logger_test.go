package logger

import (
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestLogger(t *testing.T) {
	var myLogger Logger
	myLogger.SetOutputFile("log.txt")
	myLogger.AddAgent("SELFISH", "TEAM 5")
	myLogger.AddAgent("COLLABORATE", "TEAM 6")

	AgentLog(myLogger.AgentLoggers[1])
	AgentLog(myLogger.AgentLoggers[0])
}

func AgentLog(logger *log.Logger){
	logger.Info("An agent log")
	logger.Error("Another agent log")
}