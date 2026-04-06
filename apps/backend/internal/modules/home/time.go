package home

import "time"

var currentTime = time.Now

func SetNowForTesting(now func() time.Time) func() {
	previous := currentTime
	currentTime = now
	return func() {
		currentTime = previous
	}
}
