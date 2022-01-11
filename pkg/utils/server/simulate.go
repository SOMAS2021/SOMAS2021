package server

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"github.com/SOMAS2021/SOMAS2021/pkg/config"
	"github.com/SOMAS2021/SOMAS2021/pkg/simulation"
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
	logFolderName, err := simulation.SetupLogFile(parameters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error("Unable to setup log file: " + err.Error())
	}

	err = logging.UpdateSimStatusJson(logFolderName, "running", -1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error("Unable to setup status file: " + err.Error())
	}

	log.Info("Simulation " + logFolderName + " started")
	simulation.RunNewSimulation(parameters, logFolderName)

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
