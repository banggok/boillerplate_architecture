package time_testing_helper

import (
	"testing"
	"time"
)

// Sanitize adjusts the actual time to match the expected time if the difference is less than one second.
// Otherwise, it returns an error indicating the time values are not equal.
func Sanitize(t *testing.T, actualTime *time.Time, expectedTime time.Time) {

	// Compare times with a threshold of 1 second
	if absDuration(actualTime.Sub(expectedTime)) < time.Second {
		*actualTime = expectedTime
	} else {
		t.Error("Time values are not equal: actual=" + actualTime.String() + " expected=" + expectedTime.String())
	}
}

// absDuration returns the absolute value of a time.Duration
func absDuration(d time.Duration) time.Duration {
	if d < 0 {
		return -d
	}
	return d
}
