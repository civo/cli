package utility

import (
	"fmt"
	"math"
	"time"
)

// StartTime function to start the time
func StartTime() time.Time {
	return time.Now()
}

// TrackTime function to get the start time and see
// the time between start time and the end time
func TrackTime(startTime time.Time) string {
	endTime := time.Now()
	finalEnd := endTime.Sub(startTime)
	sec := math.Round(finalEnd.Seconds())
	return secondsToMinutes(sec)
}

// utility to convert all sec to min and sec
func secondsToMinutes(inSeconds float64) string {
	minutes := int(math.Round(inSeconds)) / 60
	seconds := int(math.Round(inSeconds)) % 60
	str := fmt.Sprintf("%v min %v sec", minutes, seconds)
	return str
}
