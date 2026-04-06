package dungeon

import "time"

var currentTime = time.Now

// SetNowForTesting allows external package tests to control time-sensitive dungeon behavior.
func SetNowForTesting(now func() time.Time) func() {
	previous := currentTime
	currentTime = now
	return func() {
		currentTime = previous
	}
}
