package executor

import (
	"github.com/Worldcoin/hubble-commander/models"
)

var (
	ErrNotEnoughCommitments = NewRollupError("not enough commitments")
	ErrNoLongerProposer     = NewRollupError("commander is no longer an active proposer")
)

func (c *RollupContext) SubmitBatch(batch *models.Batch, commitments []models.Commitment) error {
	if len(commitments) < int(c.cfg.MinCommitmentsPerBatch) {
		return ErrNotEnoughCommitments
	}

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

func (c *ExecutionContext) addCommitments(commitments []models.Commitment) error {
	for i := range commitments {
		err := c.storage.AddCommitment(&commitments[i])
		if err != nil {
			return err
		}
	}
	return nil
}
