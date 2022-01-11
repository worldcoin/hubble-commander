package executor

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
)

var (
	ErrNotEnoughCommitments  = NewRollupError("not enough commitments")
	ErrRollupContextCanceled = NewLoggableRollupError("rollup context canceled")
)

func (c *TxsContext) SubmitBatch(batch *models.Batch, commitments []models.CommitmentWithTxs) error {
	select {
	case <-c.ctx.Done():
		return ErrRollupContextCanceled
	default:
	}

	tx, err := c.Executor.SubmitBatch(&batch.ID, commitments)
	if err != nil {
		return err
	}

	batch.TransactionHash = tx.Hash()
	err = c.storage.AddBatch(batch)
	if err != nil {
		return err
	}

	return c.addCommitments(commitments, batch.Type)
}

func (c *TxsContext) addCommitments(commitments []models.CommitmentWithTxs, batchType batchtype.BatchType) error {
	for i := range commitments {
		var err error

		if batchType == batchtype.Transfer || batchType == batchtype.Create2Transfer {
			err = c.storage.AddTxCommitment(&commitments[i].ToTxCommitmentWithTxs().TxCommitment)
		} else if batchType == batchtype.MassMigration {
			err = c.storage.AddMMCommitment(&commitments[i].ToMMCommitmentWithTxs().MMCommitment)
		}

		if err != nil {
			return err
		}
	}
	return nil
}

func (c *DepositsContext) SubmitBatch(batch *models.Batch, vacancyProof *models.SubtreeVacancyProof) error {
	select {
	case <-c.ctx.Done():
		return ErrRollupContextCanceled
	default:
	}

	commitmentInclusionProof, err := c.proverCtx.PreviousBatchCommitmentInclusionProof(batch.ID)
	if err != nil {
		return err
	}

	tx, err := c.client.SubmitDeposits(&batch.ID, commitmentInclusionProof, vacancyProof)
	if err != nil {
		return err
	}

	batch.TransactionHash = tx.Hash()
	return c.storage.AddBatch(batch)
}
