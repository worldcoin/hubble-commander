package commander

import "fmt"

type CommitmentError struct {
	Reason string
}

func NewCommitmentError(reason string) *CommitmentError {
	return &CommitmentError{Reason: reason}
}

func (e CommitmentError) Error() string {
	return fmt.Sprintf("Failed to commit transactions: %s", e.Reason)
}
