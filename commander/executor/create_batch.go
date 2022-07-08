package executor

import (
	"context"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
)

func (c *TxsContext) CreateAndSubmitBatch(ctx context.Context) (*models.Batch, *int, error) {
	spanCtx, span := otel.Tracer("txsContext").Start(ctx, "CreateAndSubmitBatch")
	defer span.End()

	batch, err := c.NewPendingBatch(c.BatchType)
	if err != nil {
		return nil, nil, err
	}

	// this is where we register the pending accounts
	commitments, err := c.CreateCommitments(spanCtx)
	if err != nil {
		return nil, nil, err
	}

	// this is where we actually submit!
	err = c.SubmitBatch(spanCtx, batch, commitments)
	if err != nil {
		return nil, nil, err
	}

	return batch, ref.Int(len(commitments)), nil
}

func (c *ExecutionContext) NewPendingBatch(batchType batchtype.BatchType) (*models.Batch, error) {
	prevStateRoot, err := c.storage.StateTree.Root()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	batchID, err := c.storage.GetNextBatchID()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &models.Batch{
		ID:            *batchID,
		Type:          batchType,
		PrevStateRoot: prevStateRoot,
	}, nil
}
