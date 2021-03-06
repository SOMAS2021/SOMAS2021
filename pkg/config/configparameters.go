package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/agent"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/day"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/utilFunctions"
)

type ConfigParameters struct {
	FoodOnPlatform       food.FoodType           `json:"FoodOnPlatform"`
	MaxFoodIntake        food.FoodType           `json:"MaxFoodIntake"`
	FoodPerAgentRatio    int                     `json:"FoodPerAgentRatio"`
	UseFoodPerAgentRatio bool                    `json:"UseFoodPerAgentRatio"`
	Team1Agents          int                     `json:"Team1Agents"`
	Team2Agents          int                     `json:"Team2Agents"`
	Team3Agents          int                     `json:"Team3Agents"`
	Team4Agents          int                     `json:"Team4Agents"`
	Team5Agents          int                     `json:"Team5Agents"`
	Team6Agents          int                     `json:"Team6Agents"`
	Team7Agents          int                     `json:"Team7Agents"`
	RandomAgents         int                     `json:"RandomAgents"`
	AgentHP              int                     `json:"AgentHP"`
	AgentsPerFloor       int                     `json:"AgentsPerFloor"`
	TicksPerFloor        int                     `json:"TicksPerFloor"`
	SimDays              int                     `json:"SimDays"`
	ReshuffleDays        int                     `json:"ReshuffleDays"`
	MaxHP                int                     `json:"maxHP"`
	WeakLevel            int                     `json:"weakLevel"`
	Width                float64                 `json:"width"`
	Tau                  float64                 `json:"tau"`
	HpReqCToW            int                     `json:"hpReqCToW"`
	HpCritical           int                     `json:"hpCritical"`
	MaxDayCritical       int                     `json:"maxDayCritical"`
	HPLossBase           int                     `json:"HPLossBase"`
	HPLossSlope          float64                 `json:"HPLossSlope"`
	LogFolderName        string                  `json:"LogFileName"`
	LogMain              bool                    `json:"LogMain"`
	LogStory             bool                    `json:"LogStory"`
	SimTimeoutSeconds    int                     `json:"SimTimeoutSeconds"`
	NumOfAgents          map[agent.AgentType]int `json:"-"`
	NumberOfFloors       int                     `json:"-"`
	TicksPerDay          int                     `json:"-"`
	DayInfo              *day.DayInfo            `json:"-"`
	CustomLog            string                  `json:"-"`
}

type SimulateResponse struct { // used for HTTP response on /simulate
	LogFolderName string `json:"LogFileName"`
}

type DirectoryResponse struct { // used for HTTP response on /directory
	FolderNames []string `json:"FolderNames"`
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
	parameters.NumOfAgents = map[agent.AgentType]int{
		agent.Team1:       parameters.Team1Agents,
		agent.Team2:       parameters.Team2Agents,
		agent.Team3:       parameters.Team3Agents,
		agent.Team4:       parameters.Team4Agents,
		agent.Team5:       parameters.Team5Agents,
		agent.Team6:       parameters.Team6Agents,
		agent.Team7:       parameters.Team7Agents,
		agent.RandomAgent: parameters.RandomAgents,
	}

	err := CheckParametersAreValid(parameters)
	if err != nil {
		return err
	}

	if parameters.UseFoodPerAgentRatio {
		parameters.FoodOnPlatform = food.FoodType(parameters.FoodPerAgentRatio * utilFunctions.Sum(parameters.NumOfAgents))
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

	if parameters.UseFoodPerAgentRatio {
		if parameters.FoodPerAgentRatio == 0 {
			return errors.New("foodPerAgentRatio not initialised or set to 0")
		}
	} else {
		if parameters.FoodOnPlatform == 0 {
			return errors.New("foodOnPlatform not initialised or set to 0")
		}
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
		return errors.New("tau not initialised or set to 0")
	}

	if parameters.MaxFoodIntake == 0 {
		return errors.New("MaxFoodIntake not initialised or set to 0")
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

	if parameters.SimTimeoutSeconds == 0 {
		return errors.New("SimTimeoutSeconds not initialised or set to 0")
	}

	return nil
}
