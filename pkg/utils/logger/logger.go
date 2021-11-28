package logger

import (
	"log"
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

func (L *Logger)  AddAgent(prefix string){
		AgentLogger := log.New(L.outputFile, strings.ToUpper(prefix), log.LstdFlags)
		L.AgentLoggers = append(L.AgentLoggers, AgentLogger)
}