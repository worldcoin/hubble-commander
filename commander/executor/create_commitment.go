package executor

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
)

func (t *TransactionExecutor) createAndStoreCommitment(
	txType txtype.TransactionType,
	feeReceiverStateID uint32,
	serializedTxs []byte,
	combinedSignature *models.Signature,
) (*models.Commitment, error) {
	stateRoot, err := t.storage.StateTree.Root()
	if err != nil {
		return nil, err
	}

	commitment := models.Commitment{
		Type:              txType,
		Transactions:      serializedTxs,
		FeeReceiver:       feeReceiverStateID,
		CombinedSignature: *combinedSignature,
		PostStateRoot:     *stateRoot,
	}

	err = t.storage.AddCommitment(&commitment)
	if err != nil {
		return nil, err
	}

	//commitment.IndexInBatch = *id

	return &commitment, nil
}
