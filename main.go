package main

import (
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
	if *modePtr == "serve" {
		port := ":" + strconv.Itoa(*portPtr)
		fs := http.FileServer(http.Dir("./build"))
		http.Handle("/", fs)
		log.Println("Listening on " + port + "...")
		http.HandleFunc("/simulate", func(w http.ResponseWriter, r *http.Request) {
			// TODO add simulation end-point handling and parsing parameters from body json
			fmt.Println("Will add simulation endpoint soon!")
		})
		err := http.ListenAndServe(port, nil)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println("Loading parameters...")
		// load parameters
		parameters, err := config.LoadParamFromJson(*configPathPtr)
		if err != nil {
			fmt.Println(err)
			return
		}
		runNewSimulation(parameters)
	}
}

func runNewSimulation(parameters config.ConfigParameters) {
	rand.Seed(time.Now().UnixNano())
	f := setupLogFile()
	healthInfo := health.NewHealthInfo(&parameters)

	// TODO: agentParameters - struct
	simEnv := simulation.NewSimEnv(&parameters, healthInfo)
	simEnv.Simulate()
	f.Close()
}

func setupLogFile() (fp *os.File) {
	// Setup a log file
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		err := os.Mkdir("logs", 0755)
		if err != nil {
			fmt.Println("failed to create logs directory: ", err)
			return
		}
	}

	logfileName := filepath.Join("logs", time.Now().Format("2006-01-02-15-04-05")+".json")
	f, err := os.OpenFile(logfileName, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("error opening file: ", err)
		return
	}
	log.SetOutput(f)
	log.SetFormatter(&log.JSONFormatter{})
	return f
}
