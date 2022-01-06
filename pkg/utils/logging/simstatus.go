package logging

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

type SimStatus struct {
	Status string `json:"status"`
}

func UpdateSimStatusJson(folderPath string, status string) error {

	statusStruct := SimStatus{
		Status: status,
	}

	statusJson, err := json.Marshal(statusStruct)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filepath.Join(folderPath, "status.json"), statusJson, 0644)
	if err != nil {
		return err
	}

	return nil
}
