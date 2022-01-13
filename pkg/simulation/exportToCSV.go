package simulation

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"

	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
)

func (sE *SimEnv) ExportCSV(data [][]string, filePath string) {

	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalln("failed to open file", err)
	}
	defer file.Close()
	w := csv.NewWriter(file)
	defer w.Flush()
	// Using WriteAll
	if err := w.WriteAll(data); err != nil {
		log.Fatalln("error writing record to file", err)
	}
}

func (sE *SimEnv) CreateBehaviourCtrData() {

	row := []string{
		"Altruist",
		"Collectivist",
		"Selfish",
		"Narcissist",
	}
	sE.dayInfo.BehaviourCtrData = append(sE.dayInfo.BehaviourCtrData, row)
}

func (sE *SimEnv) AddToBehaviourCtrData() {
	row := []string{
		strconv.Itoa(sE.dayInfo.BehaviourCtr["Altruist"]),
		strconv.Itoa(sE.dayInfo.BehaviourCtr["Collectivist"]),
		strconv.Itoa(sE.dayInfo.BehaviourCtr["Selfish"]),
		strconv.Itoa(sE.dayInfo.BehaviourCtr["Narcissist"]),
	}

	sE.dayInfo.BehaviourCtrData = append(sE.dayInfo.BehaviourCtrData, row)
}

func (sE *SimEnv) CreateBehaviourChangeCtrData() {

	row := []string{
		"A2A",
		"A2C",
		"A2S",
		"A2N",
		"C2A",
		"C2C",
		"C2S",
		"C2N",
		"S2A",
		"S2C",
		"S2S",
		"S2N",
		"N2A",
		"N2C",
		"N2S",
		"N2N",
	}
	sE.dayInfo.BehaviourChangeCtrData = append(sE.dayInfo.BehaviourChangeCtrData, row)
}

func (sE *SimEnv) AddToBehaviourChangeCtrData() {
	row := []string{
		strconv.Itoa(sE.dayInfo.BehaviourChangeCtr["A2A"]),
		strconv.Itoa(sE.dayInfo.BehaviourChangeCtr["A2C"]),
		strconv.Itoa(sE.dayInfo.BehaviourChangeCtr["A2S"]),
		strconv.Itoa(sE.dayInfo.BehaviourChangeCtr["A2N"]),
		strconv.Itoa(sE.dayInfo.BehaviourChangeCtr["C2A"]),
		strconv.Itoa(sE.dayInfo.BehaviourChangeCtr["C2C"]),
		strconv.Itoa(sE.dayInfo.BehaviourChangeCtr["C2S"]),
		strconv.Itoa(sE.dayInfo.BehaviourChangeCtr["C2N"]),
		strconv.Itoa(sE.dayInfo.BehaviourChangeCtr["S2A"]),
		strconv.Itoa(sE.dayInfo.BehaviourChangeCtr["S2C"]),
		strconv.Itoa(sE.dayInfo.BehaviourChangeCtr["S2S"]),
		strconv.Itoa(sE.dayInfo.BehaviourChangeCtr["S2N"]),
		strconv.Itoa(sE.dayInfo.BehaviourChangeCtr["N2A"]),
		strconv.Itoa(sE.dayInfo.BehaviourChangeCtr["N2C"]),
		strconv.Itoa(sE.dayInfo.BehaviourChangeCtr["N2S"]),
		strconv.Itoa(sE.dayInfo.BehaviourChangeCtr["N2N"]),
	}
	sE.dayInfo.BehaviourChangeCtrData = append(sE.dayInfo.BehaviourChangeCtrData, row)
}

func (sE *SimEnv) CreateUtilityData() {
	sE.dayInfo.UtilityData = append(sE.dayInfo.UtilityData, sE.utilityCSVHeader...)
}

func (sE *SimEnv) AddToUtilityData() {
	sE.dayInfo.UtilityData = append(sE.dayInfo.UtilityData, sE.dayInfo.Utility)
}

func (sE *SimEnv) CreateDeathData() {

	row := []string{
		"New",
		"Cumulative",
	}
	sE.dayInfo.DeathData = append(sE.dayInfo.DeathData, row)
}

func (sE *SimEnv) AddToDeathData(cumDeath int) {
	// cumDeath := sE.stateLog.DeathCount()
	// newDeath := cumDeath - sE.cumulativeDeaths
	// cumDeath := t.DeathCount
	newDeath := cumDeath - sE.cumulativeDeaths

	row := []string{
		strconv.Itoa(newDeath),
		strconv.Itoa(cumDeath),
	}
	sE.dayInfo.DeathData = append(sE.dayInfo.DeathData, row)
	sE.cumulativeDeaths = cumDeath
	sE.Log("Summary of deaths", infra.Fields{"New deaths": newDeath, "Cumulative deaths": cumDeath})
}
