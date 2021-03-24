package commander

import "fmt"

type CommitmentError struct {
	Reason string
}

func NewCommitmentError(reason string) *CommitmentError {
	return &CommitmentError{Reason: reason}
}

func (e CommitmentError) Error() string {
	return fmt.Sprintf("failed to commit transactions: %s", e.Reason)
}

type BatchError struct {
	Reason string
}

func NewBatchError(reason string) *BatchError {
	return &BatchError{Reason: reason}
}

func (e BatchError) Error() string {
	return fmt.Sprintf("failed to submit batch: %s", e.Reason)
}
