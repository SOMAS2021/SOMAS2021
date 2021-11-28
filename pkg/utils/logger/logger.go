package logger

import (
	"fmt"
	"log"
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

type Logger struct {
	outputFile string
}

func (L *Logger) SetOutputFile(a string) {
	L.outputFile = a
	fmt.Printf(L.outputFile)
}
