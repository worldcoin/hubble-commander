package executor

import (
	"github.com/Worldcoin/hubble-commander/models"
)

var (
	ErrNotEnoughCommitments = NewRollupError("not enough commitments")
	ErrNoLongerProposer     = NewRollupError("commander is no longer an active proposer")
)

func (t *ExecutionContext) SubmitBatch(batch *models.Batch, commitments []models.Commitment) error {
	if len(commitments) < int(t.cfg.MinCommitmentsPerBatch) {
		return ErrNotEnoughCommitments
	}

	select {
	case <-t.ctx.Done():
		return ErrNoLongerProposer
	default:
	}

	tx, err := t.Executor.SubmitBatch(t.client, commitments)
	if err != nil {
		return err
	}

	batch.TransactionHash = tx.Hash()
	err = t.storage.AddBatch(batch)
	if err != nil {
		return err
	}

	return t.addCommitments(commitments)
}

func (t *ExecutionContext) addCommitments(commitments []models.Commitment) error {
	for i := range commitments {
		err := t.storage.AddCommitment(&commitments[i])
		if err != nil {
			return err
		}
	}
	return nil
}
