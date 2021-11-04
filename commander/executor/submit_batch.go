package executor

import (
	"github.com/Worldcoin/hubble-commander/models"
)

var (
	ErrNotEnoughCommitments  = NewRollupError("not enough commitments")
	ErrRollupContextCanceled = NewLoggableRollupError("rollup context cancelled")
)

func (c *RollupContext) SubmitBatch(batch *models.Batch, commitments []models.CommitmentWithTxs) error {
	select {
	case <-c.ctx.Done():
		return ErrRollupContextCanceled
	default:
	}

	tx, err := c.Executor.SubmitBatch(c.client, commitments)
	if err != nil {
		return err
	}

	batch.TransactionHash = tx.Hash()
	err = c.storage.AddBatch(batch)
	if err != nil {
		return err
	}

	return c.addCommitments(commitments)
}

func (c *RollupContext) addCommitments(commitments []models.CommitmentWithTxs) error {
	for i := range commitments {
		err := c.storage.AddTxCommitment(&commitments[i].TxCommitment)
		if err != nil {
			return err
		}
	}
	return nil
}
