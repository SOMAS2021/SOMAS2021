package config

import (
	"encoding/json"
	"errors"
	"net/http"
)

type ConfigReadLog struct {
	LogFileName string `json:"LogFileName"`
	LogType     string `json:"LogType"`
}

type ReadLogResponse struct {
	Log []string `json:"Log"`
}

func LoadReadLogParamFromHTTPRequest(r *http.Request) (ConfigReadLog, error) {

	var configReadLog ConfigReadLog

	err := json.NewDecoder(r.Body).Decode(&configReadLog)
	if err != nil {
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
