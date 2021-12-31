package logging

import (
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

// struct definition
type LogManager struct {
	loggers    map[string]*log.Logger
	folderPath string
}

// constructor
func NewLogger(folderPath string) LogManager {
	var myLogger LogManager
	myLogger.loggers = map[string]*log.Logger{}
	myLogger.folderPath = folderPath
	return myLogger
}

// adding a new logger
func (L *LogManager) AddLogger(key string, logpath string) (logger *log.Logger, err error) {
	Logger := log.New()
	_, exists := L.loggers[key]
	if exists {
		return nil, os.ErrExist
	}
	if logpath != "" {
		file, err := os.OpenFile(filepath.Join(L.folderPath, logpath), os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			return nil, err
		}
		Logger.SetOutput(file)
		Logger.SetFormatter(&log.JSONFormatter{})
	}
	L.loggers[key] = Logger
	return Logger, err
}

// getting an existing logger
func (L *LogManager) GetLogger(key string) *log.Logger {
	if L.loggers == nil {
		return nil
	}
	logger, exists := L.loggers[key]
	if exists {
		return logger
	}
	return nil
}
