package simulation

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
)

func (sE *SimEnv) ExportCSV(filePath string) {

	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalln("failed to open file", err)
	}
	defer file.Close()
	w := csv.NewWriter(file)
	defer w.Flush()
	// Using WriteAll
	if err := w.WriteAll(sE.dayInfo.BehaviourCtrData); err != nil {
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
