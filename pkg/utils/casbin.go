package utils

import (
	"fmt"
	"time"
)

// TimeMatchFunc is the wrapper for TimeMatch.
func TimeMatchFunc(args ...string) (bool, error) {
	if len(args) != 2 {
		return false, fmt.Errorf("args error")
	}
	return TimeMatch(args[0], args[1])
}

// TimeMatch determines whether the current time is between startTime and endTime.
// You can use "_" to indicate that the parameter is ignored
func TimeMatch(startTime, endTime string) (bool, error) {
	now := time.Now()
	if startTime != "_" {
		if start, err := time.Parse("2006-01-02 15:04:05", startTime); err != nil {
			return false, err
		} else if !now.After(start) {
			return false, nil
		}
	}
	if endTime != "_" {
		if end, err := time.Parse("2006-01-02 15:04:05", endTime); err != nil {
			return false, err
		} else if !now.Before(end) {
			return false, nil
		}
	}
	return true, nil
}
