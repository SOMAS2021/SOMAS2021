package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/SOMAS2021/SOMAS2021/pkg/config"
	"github.com/SOMAS2021/SOMAS2021/pkg/simulation"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/health"
	log "github.com/sirupsen/logrus"
)

func main() {
	rand.Seed(time.Now().UnixNano())

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
	defer f.Close()

	log.SetOutput(f)
	log.SetFormatter(&log.JSONFormatter{})

	//processing the command-line flags
	//currently only one used is -configpath, but we are likely to need more in the near future.
	// if not set, it uses default of "config.json"
	configPathPtr := flag.String("configpath", "config.json", "path for parameter configuration json file")

	parameters, err := config.LoadParamFromJson(*configPathPtr)
	if err != nil {
		fmt.Println(err)
		return
	}

	healthInfo := health.NewHealthInfo(&parameters)

	// TODO: agentParameters - struct

	simEnv := simulation.NewSimEnv(&parameters, healthInfo)

	simEnv.Simulate()
}
