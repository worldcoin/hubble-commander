package executor

import (
	"errors"
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
	return fmt.Sprintf("failed to validate transfer: %s", e.Reason)
}

func IsDisputableTransferError(err error) bool {
	if err == nil {
		return false
	}
	target := &DisputableTransferError{}
	return errors.As(err, &target)
}

type BatchRaceConditionError struct {
	LocalBatch *models.Batch
}

func NewBatchRaceConditionError(localBatch *models.Batch) *BatchRaceConditionError {
	return &BatchRaceConditionError{LocalBatch: localBatch}
}

func (e BatchRaceConditionError) Error() string {
	return fmt.Sprintf("local batch #%s inconsistent with remote batch", e.LocalBatch.ID.String())
}

func IsBatchRaceConditionError(err error) bool {
	if err == nil {
		return false
	}
	target := &BatchRaceConditionError{}
	return errors.As(err, &target)
}
