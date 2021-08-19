package executor

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
)

func (t *TransactionExecutor) createCommitment(
	commitmentKey *models.CommitmentKey,
	txType txtype.TransactionType,
	feeReceiverStateID uint32,
	serializedTxs []byte,
	combinedSignature *models.Signature,
) (*models.Commitment, error) {
	stateRoot, err := t.storage.StateTree.Root()
	if err != nil {
		return nil, err
	}

	return &models.Commitment{
		ID:                *commitmentKey,
		Type:              txType,
		FeeReceiver:       feeReceiverStateID,
		CombinedSignature: *combinedSignature,
		PostStateRoot:     *stateRoot,
		IncludedInBatch:   nil,
		Transactions:      serializedTxs,
	}, nil
}

func (t *TransactionExecutor) createCommitmentKey() (*models.CommitmentKey, error) {
	nextBatchID, err := t.storage.GetNextBatchID()
	if err != nil {
		return nil, err
	}
	return &models.CommitmentKey{BatchID: *nextBatchID}, nil
}
