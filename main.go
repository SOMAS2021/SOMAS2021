package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/SOMAS2021/SOMAS2021/pkg/config"
	"github.com/SOMAS2021/SOMAS2021/pkg/simulation"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/health"
	log "github.com/sirupsen/logrus"
)

func main() {
	// processing the command-line flags
	configPathPtr := flag.String("configpath", "config.json", "path for parameter configuration json file")
	modePtr := flag.String("mode", "sim", "Execution mode. Either 'sim' for running a simulation and exiting, or 'serve' to serve a simulation endpoint")
	portPtr := flag.Int("port", 9000, "Port to run the server on if mode='serve'")
	flag.Parse()

	// check backend mode
	if *modePtr == "serve" { // HTTP serve mode
		port := ":" + strconv.Itoa(*portPtr)
		fs := http.FileServer(http.Dir("./build"))
		http.Handle("/", fs)
		log.Println("Listening on " + port + "...")
		http.HandleFunc("/simulate", func(w http.ResponseWriter, r *http.Request) {
			parameters, err := config.LoadParamFromHTTPRequest(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// logFileName returned to be used in dashboard
			logfileName := runNewSimulation(parameters)

			//generate the http response
			w.Header().Set("Content-Type", "application/json")

			response := config.Response{
				Success:     true, // this will depend on timeouts in the future, for now it is hardcoded until i figure out how timeouts work
				LogFileName: logfileName,
			}
			err = json.NewEncoder(w).Encode(response)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		})

		err := http.ListenAndServe(port, nil)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println("Loading parameters...")

		parameters, err := config.LoadParamFromJson(*configPathPtr)
		if err != nil {
			fmt.Println(err)
			return
		}

		runNewSimulation(parameters)
	}
}

// Returns the logfile name as it is needed in the HTTP response
func runNewSimulation(parameters config.ConfigParameters) (logfileName string) {
	rand.Seed(time.Now().UnixNano())
	f, err := setupLogFile(parameters.LogFileName)
	if err != nil {
		return
	}
	defer f.Close()
	healthInfo := health.NewHealthInfo(&parameters)

	// TODO: agentParameters - struct
	simEnv := simulation.NewSimEnv(&parameters, healthInfo)
	simEnv.Simulate()
	return f.Name()
}

func setupLogFile(parameterLogFileName string) (fp *os.File, err error) {
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		err := os.Mkdir("logs", 0755)
		if err != nil {
			fmt.Println("failed to create logs directory: ", err)
			return nil, err
		}
	}

	logfileName := ""
	// Check if the log file name was set in config
	if len(parameterLogFileName) != 0 {
		logfileName = filepath.Join("logs", parameterLogFileName)
	} else {
		logfileName = filepath.Join("logs", time.Now().Format("2006-01-02-15-04-05")+".json")
	}

	f, err := os.OpenFile(logfileName, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("error opening file: ", err)
		return nil, err
	}

	log.SetOutput(f)
	log.SetFormatter(&log.JSONFormatter{})
	return f, nil
}
