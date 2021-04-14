package commander

import "fmt"

type BatchError struct {
	Reason string
}

func NewBatchError(reason string) *BatchError {
	return &BatchError{Reason: reason}
}

func (e BatchError) Error() string {
	return fmt.Sprintf("failed to submit batch: %s", e.Reason)
}
