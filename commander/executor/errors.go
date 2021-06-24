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

type DisputableType int

const (
	SignatureError DisputableType = iota
	TransitionError
)

type DisputableTransferError struct {
	Reason string
	Type   DisputableType
}

func NewDisputableTransferError(reason string, errorType DisputableType) *DisputableTransferError {
	return &DisputableTransferError{Reason: reason, Type: errorType}
}

func (e DisputableTransferError) Error() string {
	if e.Type == SignatureError {
		return fmt.Sprintf("failed to validate transfer signature: %s", e.Reason)
	} else {
		return fmt.Sprintf("failed to validate transfer: %s", e.Reason)
	}
}
