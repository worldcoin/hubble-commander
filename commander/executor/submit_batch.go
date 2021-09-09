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

func (t *ExecutionContext) SubmitBatch(batch *models.Batch, commitments []models.Commitment) error {
	if len(commitments) < int(t.cfg.MinCommitmentsPerBatch) {
		return ErrNotEnoughCommitments
	}

	var tx *types.Transaction
	var err error

	select {
	case <-t.ctx.Done():
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
