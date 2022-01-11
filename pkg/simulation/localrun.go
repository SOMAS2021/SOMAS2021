package simulation

import (
	"math/rand"
	"time"

	"github.com/SOMAS2021/SOMAS2021/pkg/config"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/logging"
	log "github.com/sirupsen/logrus"
)

func LocalRun(configPath string, customLogs string) {
	log.Info("Loading parameters...")

	parameters, err := config.LoadParamFromJson(configPath)
	if err != nil {
		log.Fatal("Unable to load parameters: " + err.Error())
		return
	}

	rand.Seed(time.Now().UnixNano())
	logFolderName, err := SetupLogFile(parameters)
	if err != nil {
		log.Fatal("Unable to setup log file: " + err.Error())
		return
	}

	err = logging.UpdateSimStatusJson(logFolderName, "running", -1)
	if err != nil {
		log.Fatal("Unable to setup status file: " + err.Error())
		return
	}

	parameters.CustomLog = customLogs
	RunNewSimulation(parameters, logFolderName)
}
