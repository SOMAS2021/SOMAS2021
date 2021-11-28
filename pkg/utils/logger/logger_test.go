package logger

import (
	"testing"
)

func TestLogger(t *testing.T) {
	var myLogger Logger
	myLogger.SetOutputFile("log.txt")
	myLogger.AddAgent("TEAM 6: ")
	myLogger.AddAgent("TEAM 5: ")
	myLogger.AgentLoggers[1].Println("Some logs")
	myLogger.AgentLoggers[0].Println("Some logs")
}
