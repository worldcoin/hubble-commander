package executor

import (
	"fmt"
)

type RollupError struct {
	Reason     string
	IsLoggable bool
}

func NewRollupError(reason string) *RollupError {
	return &RollupError{
		Reason:     reason,
		IsLoggable: false,
	}
}

func NewLoggableRollupError(reason string) *RollupError {
	return &RollupError{
		Reason:     reason,
		IsLoggable: true,
	}
}

func (e RollupError) Error() string {
	return fmt.Sprintf("failed to submit batch: %s", e.Reason)
}
