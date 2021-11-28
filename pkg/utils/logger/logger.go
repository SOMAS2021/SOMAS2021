package logger

import (
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

type Logger struct {
	outputFile *os.File
	Loggers map[string]*log.Logger
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

func (L *Logger)  AddLogger(logtype string, subtype string, reporter string){
	Logger := log.New()
	Fields := map[string]string{
		"type":         strings.ToUpper(logtype),
		"subtype":			strings.ToUpper(subtype),
		"reporter": 		strings.ToUpper(reporter),
	}

	Logger.AddHook(&DefaultFieldHooks{GetValues: func() map[string]string {
		return Fields
	}})
	Logger.Out = L.outputFile
	// Move this check somewhere else, maybe to instance init
	if L.Loggers == nil {
		L.Loggers = map[string]*log.Logger{}
	}
	L.Loggers[GetLoggerKey(logtype, subtype, reporter)] = Logger
}

func GetLoggerKey(logtype string, subtype string, reporter string) string{
	return logtype + "-" + subtype + "-" + reporter
}

func (L *Logger) GetLogger(logtype string, subtype string, reporter string) *log.Logger{
	// TOOD: check if map exists and if key exists
	return L.Loggers[GetLoggerKey(logtype, subtype, reporter)]
}