package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/day"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/utilFunctions"
)

type ConfigParameters struct {
	FoodOnPlatform food.FoodType `json:"FoodOnPlatform"`
	Team1Agents    int           `json:"Team1Agents"`
	Team2Agents    int           `json:"Team2Agents"`
	Team3Agents    int           `json:"Team3Agents"`
	Team4Agents    int           `json:"Team4Agents"`
	Team5Agents    int           `json:"Team5Agents"`
	Team6Agents    int           `json:"Team6Agents"`
	RandomAgents   int           `json:"RandomAgents"`
	AgentHP        int           `json:"AgentHP"`
	AgentsPerFloor int           `json:"AgentsPerFloor"`
	TicksPerFloor  int           `json:"TicksPerFloor"`
	SimDays        int           `json:"SimDays"`
	ReshuffleDays  int           `json:"ReshuffleDays"`
	MaxHP          int           `json:"maxHP"`
	WeakLevel      int           `json:"weakLevel"`
	Width          float64       `json:"width"`
	Tau            float64       `json:"tau"`
	HpReqCToW      int           `json:"hpReqCToW"`
	HpCritical     int           `json:"hpCritical"`
	MaxDayCritical int           `json:"maxDayCritical"`
	HPLossBase     int           `json:"HPLossBase"`
	HPLossSlope    float64       `json:"HPLossSlope"`
	LogFileName    string        `json:"LogFileName"`
	NumOfAgents    []int
	NumberOfFloors int
	TicksPerDay    int
	DayInfo        *day.DayInfo
}

type Response struct { // used for HTTP response
	Success     bool   `json:"Success"`
	LogFileName string `json:"LogFileName"`
}

func LoadParamFromJson(path string) (ConfigParameters, error) {

	var tempParameters ConfigParameters

	// Open our jsonFile
	jsonFile, err := os.Open(path)
	if err != nil {
		return tempParameters, errors.New("Unable to open " + path)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	//parse it and put values in tempParameters
	byteValue, _ := ioutil.ReadAll(jsonFile)

	err = json.Unmarshal(byteValue, &tempParameters)
	if err != nil {
		return tempParameters, err
	}

	CalculateDependantParameters(&tempParameters)

	return tempParameters, nil
}

//parses the body of HTTP request (assumes json format, same as the config files), and returns the config struct
func LoadParamFromHTTPRequest(r *http.Request) (ConfigParameters, error) {

	var tempParameters ConfigParameters

	err := json.NewDecoder(r.Body).Decode(&tempParameters)
	if err != nil {
		return tempParameters, err
	}

	CalculateDependantParameters(&tempParameters)

	return tempParameters, nil
}

//Some parameters depend directly on other parameters. This function calculates them and updates the original struct
func CalculateDependantParameters(parameters *ConfigParameters) {

	//appending the sizes of the agents to the array
	parameters.NumOfAgents = append(parameters.NumOfAgents, parameters.Team1Agents, parameters.Team2Agents, parameters.Team3Agents, parameters.Team4Agents, parameters.Team5Agents, parameters.Team6Agents, parameters.RandomAgents)

	//do the calculations for parameters that depend on other parameters
	parameters.NumberOfFloors = utilFunctions.Sum(parameters.NumOfAgents) / parameters.AgentsPerFloor
	parameters.TicksPerDay = parameters.NumberOfFloors * parameters.TicksPerFloor
	parameters.DayInfo = day.NewDayInfo(parameters.TicksPerFloor, parameters.TicksPerDay, parameters.SimDays, parameters.ReshuffleDays)
}
