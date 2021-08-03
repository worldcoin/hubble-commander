package executor

import (
	"fmt"

	"github.com/Worldcoin/hubble-commander/models"
)

type RollupError struct {
	Reason string
}

func NewRollupError(reason string) *RollupError {
	return &RollupError{Reason: reason}
}

func (e RollupError) Error() string {
	return fmt.Sprintf("failed to submit batch: %s", e.Reason)
}

type DisputeType uint8

const (
	Transition DisputeType = iota
	Signature
)

type DisputableError struct {
	Type            DisputeType
	Reason          string
	CommitmentIndex int
	Proofs          []models.StateMerkleProof
}

func NewDisputableError(disputeType DisputeType, reason string) *DisputableError {
	return &DisputableError{Type: disputeType, Reason: reason, Proofs: []models.StateMerkleProof{}}
}

func NewDisputableErrorWithProofs(disputeType DisputeType, reason string, proofs []models.StateMerkleProof) *DisputableError {
	return &DisputableError{Type: disputeType, Reason: reason, Proofs: proofs}
}

func (e *DisputableError) WithCommitmentIndex(index int) *DisputableError {
	e.CommitmentIndex = index
	return e
}

func (e DisputableError) Error() string {
	return fmt.Sprintf("syncing commitment failed: %s", e.Reason)
}

type InconsistentBatchError struct {
	LocalBatch *models.Batch
}

func NewInconsistentBatchError(localBatch *models.Batch) *InconsistentBatchError {
	return &InconsistentBatchError{LocalBatch: localBatch}
}

func (e InconsistentBatchError) Error() string {
	return fmt.Sprintf("local batch #%s inconsistent with remote batch", e.LocalBatch.ID.String())
}
