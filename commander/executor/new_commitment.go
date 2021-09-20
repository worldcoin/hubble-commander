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
) (*models.Commitment, error) {
	stateRoot, err := c.storage.StateTree.Root()
	if err != nil {
		return nil, err
	}

	return &models.Commitment{
		ID:                *commitmentID,
		Type:              batchType,
		FeeReceiver:       feeReceiverStateID,
		CombinedSignature: *combinedSignature,
		PostStateRoot:     *stateRoot,
		Transactions:      serializedTxs,
	}, nil
}

func (c *RollupContext) nextCommitmentID() (*models.CommitmentID, error) {
	nextBatchID, err := c.storage.GetNextBatchID()
	if err != nil {
		return nil, err
	}
	return &models.CommitmentID{BatchID: *nextBatchID}, nil
}
