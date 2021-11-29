package logmanager

import (
	"bufio"
	"os"
	"strings"
	"testing"

	log "github.com/sirupsen/logrus"
)

// helper function
// delete file
func removeFile(filepath string) {
	os.Remove(filepath)
}

// helper function
// delete file
func fileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	return !os.IsNotExist(err)
}

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// this would be an agent instance
// we give it a logger
// it does whatever it wants with it
func AgentLog(logger *log.Logger) {
	logger.Info("An agent log")
	logger.Info("Another agent log")
}

const filepath = "log.txt"

func Test1(t *testing.T) {
	t.Log("Logger using constructor and output file")
	var myLogger LogManager = NewLogger()
	removeFile(filepath)
	_, _ = myLogger.SetOutputFile(filepath)
	myLogger.AddLogger("AGENT", "SELFISH", "TEAM 6")
	AgentLog(myLogger.GetLogger("AGENT", "SELFISH", "TEAM 6"))
	if !fileExists(filepath) {
		t.Errorf("File %s does not exist", filepath)
	}
	removeFile(filepath)
}

func Test2(t *testing.T) {
	t.Log("Logger using constructor without output file")
	var myLogger LogManager = NewLogger()
	removeFile(filepath)
	myLogger.AddLogger("AGENT", "SELFISH", "TEAM 6")
	AgentLog(myLogger.GetLogger("AGENT", "SELFISH", "TEAM 6"))
	if fileExists(filepath) {
		t.Errorf("File %s should not exist", filepath)
	}
}

func Test3(t *testing.T) {
	t.Log("Verifying logger output contains default fields")
	var myLogger LogManager = NewLogger()
	removeFile(filepath)
	_, _ = myLogger.SetOutputFile(filepath)
	myLogger.AddLogger("AGENT", "SELFISH", "TEAM 6")
	AgentLog(myLogger.GetLogger("AGENT", "SELFISH", "TEAM 6"))
	lines, err := readLines(filepath)
	if err != nil {
		t.Errorf("readLines: %s", err)
	}
	keys := [3]string{"AGENT", "SELFISH", "TEAM 6"}
	for _, line := range lines {
		for _, key := range keys {
			if !strings.Contains(line, key) {
				t.Errorf("\"%s\" key was not found in logs", key)
			}
		}
	}
	removeFile(filepath)
}

func Test4(t *testing.T) {
	t.Log("Changing output")
	var myLogger LogManager = NewLogger()
	logger := myLogger.GetLogger("AGENT", "SELFISH", "TEAM 6")
	if logger != nil {
		t.Errorf("Expected nil logger")
	}
	myLogger.AddLogger("AGENT", "SELFISH", "TEAM 6")
	logger = myLogger.GetLogger("AGENT", "SELFISH", "TEAM 6")
	removeFile(filepath)
	_, _ = myLogger.SetOutputFile(filepath)
	AgentLog(logger)
	if !fileExists(filepath) {
		t.Errorf("File %s does not exist", filepath)
	}
	removeFile(filepath)
}

func Test5(t *testing.T) {
	t.Log("Coverage edge cases")
	var myLogger LogManager
	logger := myLogger.GetLogger("AGENT", "SELFISH", "TEAM 6")
	if logger != nil {
		t.Errorf("Expected nil logger")
	}
	_, err := myLogger.SetOutputFile("")
	if !os.IsNotExist(err) {
		t.Errorf("Expected file not found error")
	}
}
