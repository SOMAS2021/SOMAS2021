package main

import (
	"bufio"
	"context"
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
				log.Error("Unable to load parameters from the Simulation Request: " + err.Error())
				return
			}

			rand.Seed(time.Now().UnixNano())
			logFolderName, err := setupLogFile(parameters.LogFolderName, parameters.LogMain)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Error("Unable to setup log file: " + err.Error())
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
				runNewSimulation(parameters, logFolderName, ctx, ch)
			}()

			// Listen on our channel AND a timeout channel - which ever happens first.
			select {
			case <-ch:
				log.Info("Simulation " + logFolderName + " finished successfully")
			case <-time.After(time.Duration(parameters.SimTimeoutSeconds) * time.Second):
				http.Error(w, "Simulation Timeout", http.StatusInternalServerError)
				log.Error("Simulation " + logFolderName + " timed out")
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
				log.Error("Unable to open logs folder. There might not be any logs created yet")
				return
			}

			//put them all in a struct
			var response config.DirectoryResponse

			//this initialises the array to empty, as otherwise if there are no folders it is null
			response.FolderNames = []string{}

			for _, f := range files {
				response.FolderNames = append(response.FolderNames, f.Name())
			}

			//convert struct to json and return the response

			err = json.NewEncoder(w).Encode(response)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Error("Error while encoding the response of directory HTTP request: " + err.Error())
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
				response.Log = append(response.Log, scanner.Text())
			}
			err = json.NewEncoder(w).Encode(response)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Error("Error while encoding the response of read HTTP request: " + err.Error())
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
			log.Fatal("Unable to load parameters: " + err.Error())
			return
		}

		rand.Seed(time.Now().UnixNano())
		logFolderName, err := setupLogFile(parameters.LogFolderName, parameters.LogMain)
		if err != nil {
			log.Fatal("Unable to setup log file: " + err.Error())
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
			runNewSimulation(parameters, logFolderName, ctx, ch)
		}()

		// Listen on our channel AND a timeout channel - which ever happens first.
		select {
		case <-ch:
			log.Info("Simulation Finished Successfully")
		case <-time.After(time.Duration(parameters.SimTimeoutSeconds) * time.Second):
			log.Fatal("Simulation Timeout")
		}
	}
}

// Returns the logfile name as it is needed in the HTTP response
func runNewSimulation(parameters config.ConfigParameters, logFolderName string, ctx context.Context, ch chan<- string) {
	healthInfo := health.NewHealthInfo(&parameters)
	parameters.LogFolderName = logFolderName
	// TODO: agentParameters - struct
	simEnv := simulation.NewSimEnv(&parameters, healthInfo)
	simEnv.Simulate(ctx, ch)
}

func setupLogFile(parameterLogFileName string, saveMainLog bool) (ffolderName string, err error) {
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
	if len(parameterLogFileName) != 0 {
		logFolderName = filepath.Join("logs", parameterLogFileName)
	} else {
		logFolderName = filepath.Join("logs", time.Now().Format("2006-01-02-15-04-05"))
	}
	if _, err := os.Stat(logFolderName); os.IsNotExist(err) {
		err := os.Mkdir(logFolderName, 0755)
		if err != nil {
			return "", err
		}
	}
	return logFolderName, nil
}
