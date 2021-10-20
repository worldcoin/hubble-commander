package eth

import "fmt"

type LogNotFoundError struct {
	logName string
}

func NewLogNotFoundError(logName string) *LogNotFoundError {
	return &LogNotFoundError{logName}
}

func (e LogNotFoundError) Error() string {
	return fmt.Sprintf("log not found in the receipt: %s", e.logName)
}
