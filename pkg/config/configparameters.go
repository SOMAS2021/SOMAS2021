package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"github.com/SOMAS2021/SOMAS2021/pkg/simulation" //for the sum function. That function should probably be moved somewhere else, like utils
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/day"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
)

type ConfigParameters struct {
	FoodOnPlatform food.FoodType `json:"FoodOnPlatform"`
	NumOfAgents    []int         `json:"NumOfAgents"`
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

	json.Unmarshal(byteValue, &tempParameters)

	//do the calculations for parameters that depend on other parameters
	tempParameters.NumberOfFloors = simulation.Sum(tempParameters.NumOfAgents) / tempParameters.AgentsPerFloor
	tempParameters.TicksPerDay = tempParameters.NumberOfFloors * tempParameters.TicksPerFloor
	tempParameters.DayInfo = day.NewDayInfo(tempParameters.TicksPerFloor, tempParameters.TicksPerDay, tempParameters.SimDays, tempParameters.ReshuffleDays)

	return tempParameters, nil
}
