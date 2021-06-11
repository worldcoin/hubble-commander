package commander

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/core/types"
)

var (
	ErrNotEnoughCommitments = NewRollupError("not enough commitments")
	ErrNoLongerProposer     = NewRollupError("commander is no longer an active proposer")
)

func (t *transactionExecutor) submitBatch(batch *models.Batch, commitments []models.Commitment) error {
	if len(commitments) < int(t.cfg.MinCommitmentsPerBatch) {
		return ErrNotEnoughCommitments
	}

	var tx *types.Transaction
	var err error

	select {
	case <-t.opts.ctx.Done():
		return ErrNoLongerProposer
	default:
	}

	if batch.Type == txtype.Transfer {
		tx, err = t.client.SubmitTransfersBatch(commitments)
	} else {
		tx, err = t.client.SubmitCreate2TransfersBatch(commitments)
	}
	if err != nil {
		return err
	}

	batch.TransactionHash = tx.Hash()
	err = t.storage.AddBatch(batch)
	if err != nil {
		return err
	}

	return t.markCommitmentsAsIncluded(commitments, batch.ID)
}

func (t *transactionExecutor) markCommitmentsAsIncluded(commitments []models.Commitment, batchID models.Uint256) error {
	for i := range commitments {
		err := t.storage.MarkCommitmentAsIncluded(commitments[i].ID, batchID)
		if err != nil {
			return err
		}
	}
	return nil
}
