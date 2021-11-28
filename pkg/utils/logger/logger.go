package logger

import (
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

type Logger struct {
	outputFile *os.File
	AgentLoggers []*log.Logger
}

func (L *Logger) SetOutputFile(outputFile string) {
	file, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	L.outputFile = file
	}

type DefaultFieldHooks struct {
	GetValues func() map[string]string
}

func (h *DefaultFieldHooks) Levels() []log.Level {
	return log.AllLevels
}

func (h *DefaultFieldHooks) Fire(e *log.Entry) error {
	for key, value := range h.GetValues() {
		e.Data[key] = value
	}
	return nil
}

func (L *Logger)  AddAgent(agentName string, reporter string){
	AgentLogger := log.New()
	Fields := map[string]string{
		"type":         "AGENT",
		"subtype":			strings.ToUpper(agentName),
		"reporter": 		strings.ToUpper(reporter),
	}

	AgentLogger.AddHook(&DefaultFieldHooks{GetValues: func() map[string]string {
		return Fields
	}})
	AgentLogger.Out = L.outputFile
	L.AgentLoggers = append(L.AgentLoggers, AgentLogger)
}