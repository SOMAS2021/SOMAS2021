package execution

import (
	"context"
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
	logFolderName, err := SetupLogFile(parameters, parameters.LogMain)
	if err != nil {
		log.Fatal("Unable to setup log file: " + err.Error())
		return
	}

	err = logging.UpdateSimStatusJson(logFolderName, "running", -1)
	if err != nil {
		log.Fatal("Unable to setup status file: " + err.Error())
		return
	}

	//Timeout related stuff
	//don't touch, don't ask me how it works
	//https://stackoverflow.com/a/50579561

	//context and channel are passed all the way down to simulation loop to check every tick if there was a timeout
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan string, 1)

	go func() {
		log.Info("Simulation started")
		parameters.CustomLog = customLogs
		RunNewSimulation(parameters, logFolderName, ctx, ch)
	}()

	// Listen on our channel AND a timeout channel - which ever happens first.
	select {
	case <-ch:
		log.Info("Simulation Finished Successfully")
		err = logging.UpdateSimStatusJson(logFolderName, "finished", GetMaxTick(logFolderName+"/story.json"))
		if err != nil {
			log.Fatal("Unable to update status file: " + err.Error())
			return
		}
	case <-time.After(time.Duration(parameters.SimTimeoutSeconds) * time.Second):
		err = logging.UpdateSimStatusJson(logFolderName, "timedout", GetMaxTick(logFolderName+"/story.json"))
		if err != nil {
			log.Fatal("Unable to update status file: " + err.Error())
			return
		}
		log.Fatal("Simulation Timeout")
	}
}
