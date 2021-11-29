// logmanager package supports a logmanager struct which encapsulates IDed subloggers
// these can be passed to e.g. agents who can log independently to a shared file

package logmanager

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

// struct definition
type LogManager struct {
	outputFile *os.File
	loggers    map[string]*log.Logger
}

// constructor
func NewLogger() LogManager {
	var myLogger LogManager
	myLogger.loggers = map[string]*log.Logger{}
	return myLogger
}

// set output file. also updates the destination of exsiting loggers
func (L *LogManager) SetOutputFile(outputFile string) (file *os.File, err error) {
	file, err = os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return file, err
	}
	L.outputFile = file

	// update all loggers to start writing again to same file
	if file != nil {
		for _, logger := range L.loggers {
			if file != nil {
				logger.Out = file
			}
		}
	}
	return file, err
}

// helper function to build logger key
func getLoggerKey(logtype string, subtype string, reporter string) string {
	return logtype + "-" + subtype + "-" + reporter
}

// managing default fields
type defaultFieldHooks struct {
	GetValues func() map[string]string
}

func (h *defaultFieldHooks) Levels() []log.Level {
	return log.AllLevels
}

func (h *defaultFieldHooks) Fire(e *log.Entry) error {
	for key, value := range h.GetValues() {
		e.Data[key] = value
	}
	return nil
}

// adding a new logger
func (L *LogManager) AddLogger(logtype string, subtype string, reporter string) (created bool, logger *log.Logger) {
	Logger := log.New()
	key := getLoggerKey(logtype, subtype, reporter)
	Fields := map[string]string{
		"type":     strings.ToUpper(logtype),
		"subtype":  strings.ToUpper(subtype),
		"reporter": strings.ToUpper(reporter),
	}

	Logger.AddHook(&defaultFieldHooks{GetValues: func() map[string]string {
		return Fields
	}})
	if L.outputFile != nil {
		Logger.Out = L.outputFile
	} else {
		Logger.Out = os.Stdout
	}
	_, exists := L.loggers[getLoggerKey(logtype, subtype, reporter)]
	if !exists {
		L.loggers[key] = Logger
	}
	return !exists, Logger
}

// getting an existing logger
func (L *LogManager) GetLogger(logtype string, subtype string, reporter string) *log.Logger {
	if L.loggers == nil {
		return nil
	}
	logger, exists := L.loggers[getLoggerKey(logtype, subtype, reporter)]
	if exists {
		return logger
	}
	return nil
}
