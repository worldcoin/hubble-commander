package executor

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/core/types"
)

var (
	ErrNotEnoughCommitments = NewRollupError("not enough commitments")
	ErrNoLongerProposer     = NewRollupError("commander is no longer an active proposer")
)

func (c *ExecutionContext) SubmitBatch(batch *models.Batch, commitments []models.Commitment) error {
	if len(commitments) < int(c.cfg.MinCommitmentsPerBatch) {
		return ErrNotEnoughCommitments
	}

	var tx *types.Transaction
	var err error

	select {
	case <-c.ctx.Done():
		return ErrNoLongerProposer
	default:
	}

	if batch.Type == txtype.Transfer {
		tx, err = c.client.SubmitTransfersBatch(commitments)
	} else {
		tx, err = c.client.SubmitCreate2TransfersBatch(commitments)
	}
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
