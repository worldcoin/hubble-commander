package metrics

import "time"

func MeasureDuration(action func()) time.Duration {
	startTime := time.Now()

	action()

	return time.Since(startTime).Round(time.Millisecond)
}
