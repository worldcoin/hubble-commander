package executor

import "fmt"

type RollupError struct {
	Reason string
}

func NewRollupError(reason string) *RollupError {
	return &RollupError{Reason: reason}
}

func (e RollupError) Error() string {
	return fmt.Sprintf("failed to submit batch: %s", e.Reason)
}

type DisputableTransferError struct {
	Reason string
}

func NewDisputableTransferError(reason string) *DisputableTransferError {
	return &DisputableTransferError{Reason: reason}
}

func (e DisputableTransferError) Error() string {
	return fmt.Sprintf("failed to validate transfer: %s", e.Reason)
}
