package metrics

import "time"

func MeasureDuration(action func() error) (*time.Duration, error) {
	startTime := time.Now()

	err := action()
	if err != nil {
		return nil, err
	}

	duration := time.Since(startTime).Round(time.Millisecond)
	return &duration, nil
}
