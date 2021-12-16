package team3

import (
	"math/rand"
	"time"
)

// Function decides if we read a message or we don't, depends on our level of stubbornness, return is bool
// Stubborness of 20 means probability of 0.2 that we don't read message
func read() bool {
	rand.Seed(time.Now().UnixNano())
	var random = rand.Intn(100)

	if random <= stubbornness {
		return true
	}
	return false
}
