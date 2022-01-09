package config

import (
	"encoding/json"
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type ConfigReadLog struct {
	LogFileName string `json:"LogFileName"`
	LogType     string `json:"LogType"`
	TickFilter  bool   `json:"TickFiler"`
	Tick        int    `json:"Tick"`
}

type ReadLogResponse struct {
	Log []string `json:"Log"`
}

func LoadReadLogParamFromHTTPRequest(r *http.Request) (ConfigReadLog, error) {

	var configReadLog ConfigReadLog

	err := json.NewDecoder(r.Body).Decode(&configReadLog)
	if err != nil {
		log.Error(err)
		return configReadLog, err
	}

	if configReadLog.LogFileName == "" {
		return configReadLog, errors.New("LogFileName not initialised or set to empty")
	}

	if configReadLog.LogType == "" {
		return configReadLog, errors.New("LogType not initialised or set to empty")
	}

	return configReadLog, nil
}
