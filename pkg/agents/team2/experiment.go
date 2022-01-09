package team2

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

func (a *CustomAgent2) writeToCSV(target []float64, fileName string, elementNumber int) {
	//the input should be an array of float64
	records := make([]string, 0)
	for i := 0; i < elementNumber; i++ {
		s := fmt.Sprintf("%f", target[i])
		records = append(records, s)
	}

	id := a.ID().String()
	f, err := os.OpenFile(fileName+"_"+id+".csv", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write(records); err != nil {
		log.Fatalln("error writing record to file", err)
	}
}
