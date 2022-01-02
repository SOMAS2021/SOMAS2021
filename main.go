package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"io/ioutil"
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
	devmodePtr := flag.Bool("devmode", false, "If true, disable Access Origin Control for cross-domain requests")
	flag.Parse()

	// check backend mode
	if *modePtr == "serve" { // HTTP serve mode
		port := ":" + strconv.Itoa(*portPtr)
		fs := http.FileServer(http.Dir("./build"))
		http.Handle("/", fs)
		log.Info("Listening on " + port + "...")
		if *devmodePtr {
			log.Info("DEV MODE ACTIVE")
		}

		// Simulation endpoint
		http.HandleFunc("/simulate", func(w http.ResponseWriter, r *http.Request) {
			if *devmodePtr {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			}
			parameters, err := config.LoadParamFromHTTPRequest(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// logFileName returned to be used in dashboard
			logFileName := ""

			//channel and goroutine used for timeouts
			c1 := make(chan string, 1)

			go func() {
				filenametemp := runNewSimulation(parameters)
				c1 <- filenametemp
			}()

			// Listen on our channel AND a timeout channel - which ever happens first.
			select {
			case res := <-c1:
				logFileName = res
			case <-time.After(time.Duration(parameters.SimTimeoutSeconds) * time.Second):

				http.Error(w, "Simulation Timeout", http.StatusInternalServerError)
				return
			}

			//generate the http response
			w.Header().Set("Content-Type", "application/json")

			response := config.SimulateResponse{
				LogFileName: logFileName,
			}
			err = json.NewEncoder(w).Encode(response)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		})

		// Directory fetch endpoint
		http.HandleFunc("/directory", func(w http.ResponseWriter, r *http.Request) {
			if *devmodePtr {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			}
			//read directory
			files, err := ioutil.ReadDir("./logs")
			if err != nil {
				http.Error(w, "Unable to open logs folder. There might not be any logs created yet", http.StatusInternalServerError)
				return
			}

			//put them all in a struct
			var response config.DirectoryResponse

			for _, f := range files {
				response.FolderNames = append(response.FolderNames, f.Name())
			}

			//convert struct to json and return the response

			err = json.NewEncoder(w).Encode(response)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		})

		// File fetch endpoint
		http.HandleFunc("/read", func(w http.ResponseWriter, r *http.Request) {
			if *devmodePtr {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			}

			logParams, err := config.LoadReadLogParamFromHTTPRequest(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			logFileName := filepath.Join("logs", logParams.LogFileName, logParams.LogType)
			file, err := os.Open(logFileName + ".json")
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			defer file.Close()
			var response config.ReadLogResponse
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				response.Log = append(response.Log, scanner.Text())
			}
			err = json.NewEncoder(w).Encode(response)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		})

		err := http.ListenAndServe(port, nil)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Info("Loading parameters...")

		parameters, err := config.LoadParamFromJson(*configPathPtr)
		if err != nil {
			log.Error(err)
			return
		}

		//channel and goroutine used for timeout
		c1 := make(chan string, 1)

		go func() {
			filenametemp := runNewSimulation(parameters)
			c1 <- filenametemp
		}()

		// Listen on our channel AND a timeout channel - which ever happens first.
		select {
		case <-c1:
			log.Info("Simulation Finished Successfully")
		case <-time.After(time.Duration(parameters.SimTimeoutSeconds) * time.Second):
			log.Fatal("Simulation Timeout")
		}
	}
}

// Returns the logfile name as it is needed in the HTTP response
func runNewSimulation(parameters config.ConfigParameters) (logFolderName string) {
	rand.Seed(time.Now().UnixNano())
	logFolderName, err := setupLogFile(parameters.LogFileName, parameters.LogMain)
	if err != nil {
		return
	}
	healthInfo := health.NewHealthInfo(&parameters)
	parameters.LogFileName = logFolderName
	// TODO: agentParameters - struct
	simEnv := simulation.NewSimEnv(&parameters, healthInfo)
	simEnv.Simulate()
	return logFolderName
}

func setupLogFile(parameterLogFileName string, saveMainLog bool) (ffolderName string, err error) {
	// setup logs folder if never created
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		err := os.Mkdir("logs", 0755)
		if err != nil {
			log.Error("failed to create logs directory: ", err)
			return "", err
		}
	}

	// setup simulation run folder for logs
	logFolderName := ""
	// Check if the log folder name was set in config
	if len(parameterLogFileName) != 0 {
		logFolderName = filepath.Join("logs", parameterLogFileName)
	} else {
		logFolderName = filepath.Join("logs", time.Now().Format("2006-01-02-15-04-05"))
	}
	if _, err := os.Stat(logFolderName); os.IsNotExist(err) {
		err := os.Mkdir(logFolderName, 0755)
		if err != nil {
			log.Error("failed to create custom folder directory: ", err)
			return "", err
		}
	}
	return logFolderName, nil
}
