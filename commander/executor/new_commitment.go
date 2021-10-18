package executor

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
)

func (c *RollupContext) newCommitment(
	commitmentID *models.CommitmentID,
	batchType batchtype.BatchType,
	feeReceiverStateID uint32,
	serializedTxs []byte,
	combinedSignature *models.Signature,
) (*models.TxCommitment, error) {
	stateRoot, err := c.storage.StateTree.Root()
	if err != nil {
		return nil, err
	}

	return &models.TxCommitment{
		CommitmentBase: models.CommitmentBase{
			ID:            *commitmentID,
			Type:          batchType,
			PostStateRoot: *stateRoot,
		},
		FeeReceiver:       feeReceiverStateID,
		CombinedSignature: *combinedSignature,
		Transactions:      serializedTxs,
	}, nil
}

func (c *RollupContext) NextCommitmentID() (*models.CommitmentID, error) {
	nextBatchID, err := c.storage.GetNextBatchID()
	if err != nil {
		return nil, err
	}
	return &models.CommitmentID{BatchID: *nextBatchID}, nil
}
