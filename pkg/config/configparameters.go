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
	Team7Agent1    int           `json:"Team7Agent1"`
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
	Error       string `json:"Error"`
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

	err = CalculateDependentParameters(&tempParameters)
	if err != nil {
		return tempParameters, err
	}

	return tempParameters, nil
}

//parses the body of HTTP request (assumes json format, same as the config files), and returns the config struct
func LoadParamFromHTTPRequest(r *http.Request) (ConfigParameters, error) {

	var tempParameters ConfigParameters

	err := json.NewDecoder(r.Body).Decode(&tempParameters)
	if err != nil {
		return tempParameters, err
	}

	err = CalculateDependentParameters(&tempParameters)
	if err != nil {
		return tempParameters, err
	}

	return tempParameters, nil
}

//Some parameters depend directly on other parameters. This function calculates them and updates the original struct
func CalculateDependentParameters(parameters *ConfigParameters) error {

	//appending the sizes of the agents to the array
	parameters.NumOfAgents = append(parameters.NumOfAgents, parameters.Team1Agents, parameters.Team2Agents, parameters.Team3Agents, parameters.Team4Agents, parameters.Team5Agents, parameters.Team6Agents, parameters.Team7Agent1, parameters.RandomAgents)

	err := CheckParametersAreValid(parameters)
	if err != nil {
		return err
	}

	//do the calculations for parameters that depend on other parameters
	parameters.NumberOfFloors = utilFunctions.Sum(parameters.NumOfAgents) / parameters.AgentsPerFloor
	parameters.TicksPerDay = parameters.NumberOfFloors * parameters.TicksPerFloor
	parameters.DayInfo = day.NewDayInfo(parameters.TicksPerFloor, parameters.TicksPerDay, parameters.SimDays, parameters.ReshuffleDays)

	return nil
}

//Checking that there are no issues with the config provided, for example 0 ticks per floor, 0 agents per floor, etc.
//This can happen if they are missing from the json config / http request, as they'd all be initialised to 0 / nil
func CheckParametersAreValid(parameters *ConfigParameters) error {

	if parameters.FoodOnPlatform == 0 {
		return errors.New("foodOnPlatform not initialised or set to 0")
	}

	if utilFunctions.Sum(parameters.NumOfAgents) == 0 {
		return errors.New("no agents initialised")
	}

	if parameters.AgentHP == 0 {
		return errors.New("agentHP not initialised or set to 0")
	}

	if parameters.AgentsPerFloor == 0 {
		return errors.New("agentsPerFloor not initialised or set to 0")
	}

	if parameters.TicksPerFloor == 0 {
		return errors.New("ticksPerFloor not initialised or set to 0")
	}

	if parameters.SimDays == 0 {
		return errors.New("simDays not initialised or set to 0")
	}

	if parameters.ReshuffleDays == 0 {
		return errors.New("reshuffleDays not initialised or set to 0")
	}

	if parameters.MaxHP == 0 {
		return errors.New("maxHP not initialised or set to 0")
	}

	if parameters.WeakLevel == 0 {
		return errors.New("weakLevel not initialised or set to 0")
	}

	if parameters.Width == 0 {
		return errors.New("width not initialised or set to 0")
	}

	if parameters.Tau == 0 {
		return errors.New("wau not initialised or set to 0")
	}

	if parameters.HpReqCToW == 0 {
		return errors.New("hpReqCToW not initialised or set to 0")
	}

	if parameters.HpCritical == 0 {
		return errors.New("hpCritical not initialised or set to 0")
	}

	if parameters.MaxDayCritical == 0 {
		return errors.New("maxDayCritical not initialised or set to 0")
	}

	if parameters.HPLossBase == 0 {
		return errors.New("hpLossBase not initialised or set to 0")
	}

	if parameters.HPLossSlope == 0 {
		return errors.New("hpLossSlope not initialised or set to 0")
	}
	return nil
}
