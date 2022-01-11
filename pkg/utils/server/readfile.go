package server

import (
	"bufio"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/SOMAS2021/SOMAS2021/pkg/config"
	log "github.com/sirupsen/logrus"
)

func ReadFile(w http.ResponseWriter, r *http.Request, devMode bool) {
	if devMode {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	logParams, err := config.LoadReadLogParamFromHTTPRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Error("Unable to load parameters from the File Read Request: " + err.Error())
		return
	}

	logFileName := filepath.Join("logs", logParams.LogFileName, logParams.LogType)
	file, err := os.Open(logFileName + ".json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Error("Unable to open the file " + logFileName + ".json: " + err.Error())
		return
	}
	defer file.Close()
	var response config.ReadLogResponse
	response.Log = []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		a := scanner.Text()
		if !logParams.TickFilter || strings.Contains(a, "\"tick\":"+strconv.Itoa(logParams.Tick)+",") {
			response.Log = append(response.Log, scanner.Text())
		}
	}

	// special files we know are a single json
	if logParams.LogType == "config" {
		response.Log = []string{strings.Join(response.Log[:], "")}
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error("Error while encoding the response of read HTTP request: " + err.Error())
		return
	}
}
