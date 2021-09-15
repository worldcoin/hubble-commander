package executor

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
)

func (c *ExecutionContext) createCommitment(
	commitmentID *models.CommitmentID,
	txType txtype.TransactionType,
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
		Type:              txType,
		FeeReceiver:       feeReceiverStateID,
		CombinedSignature: *combinedSignature,
		PostStateRoot:     *stateRoot,
		Transactions:      serializedTxs,
	}, nil
}

func (c *ExecutionContext) createCommitmentID() (*models.CommitmentID, error) {
	nextBatchID, err := c.storage.GetNextBatchID()
	if err != nil {
		return nil, err
	}
	return &models.CommitmentID{BatchID: *nextBatchID}, nil
}
