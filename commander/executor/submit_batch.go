package executor

import (
	"context"

	"github.com/Worldcoin/hubble-commander/models"
	"go.opentelemetry.io/otel"
)

var ErrRollupContextCanceled = NewLoggableRollupError("rollup context canceled")

func (c *TxsContext) SubmitBatch(ctx context.Context, batch *models.Batch, commitments []models.CommitmentWithTxs) error {
	spanCtx, span := otel.Tracer("txsContext").Start(ctx, "SubmitBatch")
	defer span.End()

	select {
	case <-c.ctx.Done():
		return ErrRollupContextCanceled
	default:
	}

	tx, err := c.Executor.SubmitBatch(spanCtx, &batch.ID, commitments)
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

func (c *TxsContext) addCommitments(commitments []models.CommitmentWithTxs) error {
	for i := range commitments {
		err := c.storage.AddCommitment(commitments[i].ToCommitment())
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
