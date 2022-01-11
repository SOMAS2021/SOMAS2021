package main

import (
	"flag"
	"net/http"
	"strconv"

	"github.com/SOMAS2021/SOMAS2021/pkg/simulation"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/server"
	log "github.com/sirupsen/logrus"
)

func main() {
	// processing the command-line flags
	configPathPtr := flag.String("configpath", "config.json", "path for parameter configuration json file")
	modePtr := flag.String("mode", "sim", "Execution mode. Either 'sim' for running a simulation and exiting, or 'serve' to serve a simulation endpoint")
	portPtr := flag.Int("port", 9000, "Port to run the server on if mode='serve'")
	devmodePtr := flag.Bool("devmode", false, "If true, disable Access Origin Control for cross-domain requests")
	customLogsPtr := flag.String("log", "", "Assign value of agent type to call custom logging function")

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
			server.Simulate(w, r, *devmodePtr)
		})

		// Directory fetch endpoint
		http.HandleFunc("/directory", func(w http.ResponseWriter, r *http.Request) {
			server.Directory(w, r, *devmodePtr)
		})

		// File fetch endpoint
		http.HandleFunc("/read", func(w http.ResponseWriter, r *http.Request) {
			server.ReadFile(w, r, *devmodePtr)
		})

		err := http.ListenAndServe(port, nil)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		simulation.LocalRun(*configPathPtr, *customLogsPtr)
	}
}
