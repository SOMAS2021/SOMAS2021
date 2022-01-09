package execution

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/SOMAS2021/SOMAS2021/pkg/config"
	"github.com/SOMAS2021/SOMAS2021/pkg/simulation"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/health"
)

// Returns the logfile name as it is needed in the HTTP response
func RunNewSimulation(parameters config.ConfigParameters, logFolderName string, ctx context.Context, ch chan<- string) {
	healthInfo := health.NewHealthInfo(&parameters)
	parameters.LogFolderName = logFolderName
	// TODO: agentParameters - struct
	simEnv := simulation.NewSimEnv(&parameters, healthInfo)
	simEnv.Simulate(ctx, ch)
}

func SetupLogFile(parameters config.ConfigParameters) (ffolderName string, err error) {
	// setup logs folder if never created
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		err := os.Mkdir("logs", 0755)
		if err != nil {
			return "", err
		}
	}

	// setup simulation run folder for logs
	logFolderName := ""
	// Check if the log folder name was set in config
	if len(parameters.LogFolderName) != 0 {
		logFolderName = filepath.Join("logs", parameters.LogFolderName)
	} else {
		logFolderName = filepath.Join("logs", time.Now().Format("2006-01-02-15-04-05"))
	}
	if _, err := os.Stat(logFolderName); os.IsNotExist(err) {
		err := os.Mkdir(logFolderName, 0755)
		if err != nil {
			return "", err
		}
	}

	//save the config inside the folder
	jsonConfig, err := json.MarshalIndent(parameters, "", "\t")
	if err != nil {
		return "", err
	}
	err = ioutil.WriteFile(filepath.Join(logFolderName, "config.json"), jsonConfig, 0644)
	if err != nil {
		return "", err
	}

	return logFolderName, nil
}
