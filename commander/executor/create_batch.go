package executor

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/pkg/errors"
)

func (c *TxsContext) CreateAndSubmitBatch() (*models.Batch, *int, error) {
	batch, err := c.NewPendingBatch(c.BatchType)
	if err != nil {
		return nil, nil, err
	}

	commitments, err := c.CreateCommitments()
	if err != nil {
		return nil, nil, err
	}

	err = c.SubmitBatch(batch, commitments)
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
