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

type DisputableTransferError struct {
	Reason string
	Proofs []models.StateMerkleProof
}

func NewDisputableTransferError(reason error, proofs []models.StateMerkleProof) *DisputableTransferError {
	return &DisputableTransferError{Reason: reason.Error(), Proofs: proofs}
}

func NewDisputableTransferErrorWithoutProofs(reason string) *DisputableTransferError {
	return &DisputableTransferError{Reason: reason, Proofs: []models.StateMerkleProof{}}
}

func (e DisputableTransferError) Error() string {
	return fmt.Sprintf("syncing commitment failed: %s", e.Reason)
}

type DisputableTransitionError struct {
	DisputableTransferError
	CommitmentIndex int
}

func NewDisputableTransitionError(err DisputableTransferError, commitmentIndex int) *DisputableTransitionError {
	return &DisputableTransitionError{DisputableTransferError: err, CommitmentIndex: commitmentIndex}
}

type DisputableSignatureError struct {
	Reason          string
	CommitmentIndex int
	Proofs          []models.StateMerkleProof
}

func NewDisputableSignatureError(reason string) *DisputableSignatureError {
	return &DisputableSignatureError{Reason: reason}
}

func NewDisputableSignatureErrorWithProofs(reason string, proofs []models.StateMerkleProof) *DisputableSignatureError {
	return &DisputableSignatureError{Reason: reason, Proofs: proofs}
}

func (e *DisputableSignatureError) WithCommitmentIndex(index int) *DisputableSignatureError {
	e.CommitmentIndex = index
	return e
}

func (e DisputableSignatureError) Error() string {
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
