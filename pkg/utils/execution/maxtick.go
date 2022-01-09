package execution

import (
	"os"
	"regexp"
	"strconv"
)

func GetMaxTick(filename string) int {
	file, err := os.Open(filename)
	if err != nil {
		return -1
	}
	defer file.Close()

	buf := make([]byte, 100)
	stat, err := os.Stat(filename)
	start := stat.Size() - 100
	_, err = file.ReadAt(buf, start)
	if err == nil {
		str := string(buf[:])
		m := regexp.MustCompile(`\"tick\":.*,`)
		res := m.FindString(str)
		if res != "" {
			i, _ := strconv.Atoi(res[7 : len(res)-1])
			return i
		}
	}
	return -1
}
