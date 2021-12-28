package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
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
	Team7Agent1    int           `json:"Team7Agent1"`
	Team7Agent2    int           `json:"Team7Agent2"`
	Team7Agent3    int           `json:"Team7Agent3"`
	Team7Agent4    int           `json:"Team7Agent4"`
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
	NumOfAgents    []int
	NumberOfFloors int
	TicksPerDay    int
	DayInfo        *day.DayInfo
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
	//appending the sizes of the agents to the array
	tempParameters.NumOfAgents = append(tempParameters.NumOfAgents, tempParameters.Team1Agents, tempParameters.Team2Agents, tempParameters.Team3Agents, tempParameters.Team4Agents, tempParameters.Team5Agents, tempParameters.Team6Agents, tempParameters.RandomAgents, tempParameters.Team7Agent1, tempParameters.Team7Agent2, tempParameters.Team7Agent3, tempParameters.Team7Agent4)

	//do the calculations for parameters that depend on other parameters
	tempParameters.NumberOfFloors = utilFunctions.Sum(tempParameters.NumOfAgents) / tempParameters.AgentsPerFloor
	tempParameters.TicksPerDay = tempParameters.NumberOfFloors * tempParameters.TicksPerFloor
	tempParameters.DayInfo = day.NewDayInfo(tempParameters.TicksPerFloor, tempParameters.TicksPerDay, tempParameters.SimDays, tempParameters.ReshuffleDays)

	return tempParameters, nil
}
