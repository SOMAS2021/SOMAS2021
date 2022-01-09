package server

import (
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"github.com/SOMAS2021/SOMAS2021/pkg/config"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/execution"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/logging"
	log "github.com/sirupsen/logrus"
)

func Simulate(w http.ResponseWriter, r *http.Request, devMode bool) {
	if devMode {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}
	parameters, err := config.LoadParamFromHTTPRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Error("Unable to load parameters from the Simulation Request: " + err.Error())
		return
	}

	rand.Seed(time.Now().UnixNano())
	logFolderName, err := execution.SetupLogFile(parameters, parameters.LogMain)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error("Unable to setup log file: " + err.Error())
	}

	err = logging.UpdateSimStatusJson(logFolderName, "running", -1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error("Unable to setup status file: " + err.Error())
	}

	//Timeout related stuff
	//don't touch, don't ask me how it works
	//https://stackoverflow.com/a/50579561

	//context and channel are passed all the way down to simulation loop to check every tick if there was a timeout
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan string, 1)

	go func() {
		log.Info("Simulation " + logFolderName + " started")
		execution.RunNewSimulation(parameters, logFolderName, ctx, ch)
	}()

	// Listen on our channel AND a timeout channel - which ever happens first.
	select {
	case <-ch:
		log.Info("Simulation " + logFolderName + " finished successfully")
		err = logging.UpdateSimStatusJson(logFolderName, "finished", execution.GetMaxTick(logFolderName+"/story.json"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Error("Unable to update status file: " + err.Error())
		}
	case <-time.After(time.Duration(parameters.SimTimeoutSeconds) * time.Second):
		http.Error(w, "Simulation Timeout", http.StatusInternalServerError)
		log.Error("Simulation " + logFolderName + " timed out")
		err = logging.UpdateSimStatusJson(logFolderName, "timedout", execution.GetMaxTick(logFolderName+"/story.json"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Error("Unable to update status file: " + err.Error())
		}
		return
	}

	//generate the http response
	w.Header().Set("Content-Type", "application/json")

	response := config.SimulateResponse{
		LogFolderName: logFolderName,
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error("Error while encoding the response of simulation "+logFolderName+": ", err.Error())
		return
	}
}
