package time_testing_helper

import (
	"errors"
	"time"
)

// Sanitize adjusts the actual time to match the expected time if the difference is less than one second.
// Otherwise, it returns an error indicating the time values are not equal.
func Sanitize(actualTime *time.Time, expectedTime time.Time) error {

	// Compare times with a threshold of 1 second
	if absDuration(actualTime.Sub(expectedTime)) < time.Second {
		*actualTime = expectedTime
	} else {
		return errors.New("Time values are not equal: actual=" + actualTime.String() + " expected=" + expectedTime.String())
	}
	return nil
}

// absDuration returns the absolute value of a time.Duration
func absDuration(d time.Duration) time.Duration {
	if d < 0 {
		return -d
	}
	return d
}
