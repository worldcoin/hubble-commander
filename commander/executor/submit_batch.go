package executor

import (
	"github.com/Worldcoin/hubble-commander/models"
)

var (
	ErrNotEnoughCommitments = NewRollupError("not enough commitments")
	ErrNoLongerProposer     = NewLoggableRollupError("commander is no longer an active proposer")
)

func (c *RollupContext) SubmitBatch(batch *models.Batch, commitments []models.TxCommitmentWithTxs) error {
	select {
	case <-c.ctx.Done():
		return ErrNoLongerProposer
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

func (c *RollupContext) addCommitments(commitments []models.TxCommitmentWithTxs) error {
	for i := range commitments {
		err := c.storage.AddTxCommitment(&commitments[i].TxCommitment)
		if err != nil {
			return err
		}
	}
	return nil
}
