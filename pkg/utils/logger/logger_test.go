package logger

import (
	"testing"
)

func TestGreetsGitHub(t *testing.T) {
	// log.SetFlags(log.LstdFlags)
	// file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.SetOutput(file)
	// log.Println("LOG MESSAGE")
	// log.Println("LOG EKFBHJDF")
	var myLogger Logger
	myLogger.SetOutputFile("lo.txt")
}
