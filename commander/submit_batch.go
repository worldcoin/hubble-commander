package commander

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/core/types"
)

var (
	ErrNotEnoughCommitments = NewRollupError("not enough commitments")
)

func (t *transactionExecutor) submitBatch(batchType txtype.TransactionType, commitments []models.Commitment) error {
	if len(commitments) < int(t.cfg.MinCommitmentsPerBatch) {
		return ErrNotEnoughCommitments
	}

	var tx *types.Transaction
	var err error

	select {
	case <-t.ctx.Done():
		return NewRollupError("commander is no longer an active proposer")
	default:
	}

	if batchType == txtype.Transfer {
		tx, err = t.client.SubmitTransfersBatch(commitments)
	} else {
		tx, err = t.client.SubmitCreate2TransfersBatch(commitments)
	}
	if err != nil {
		return err
	}

	batchNumber, err := t.storage.GetNextBatchNumber()
	if err != nil {
		return err
	}
	newPendingBatch := models.Batch{
		Type:            batchType,
		TransactionHash: tx.Hash(),
		Number:          batchNumber,
	}
	batchID, err := t.storage.AddBatch(&newPendingBatch)
	if err != nil {
		return err
	}

	err = t.markCommitmentsAsIncluded(commitments, *batchID)
	if err != nil {
		return err
	}

	return nil
}

func (t *transactionExecutor) markCommitmentsAsIncluded(commitments []models.Commitment, batchID int32) error {
	for i := range commitments {
		err := t.storage.MarkCommitmentAsIncluded(commitments[i].ID, batchID)
		if err != nil {
			return err
		}
	}
	return nil
}
