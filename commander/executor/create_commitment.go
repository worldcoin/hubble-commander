package executor

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
)

func (t *TransactionExecutor) createAndStoreCommitment(
	txType txtype.TransactionType,
	feeReceiverIndex uint32,
	serializedTxs []byte,
	combinedSignature *models.Signature,
) (*models.Commitment, error) {
	stateRoot, err := st.NewStateTree(t.storage).Root()
	if err != nil {
		return nil, err
	}

	commitment := models.Commitment{
		Type:              txType,
		Transactions:      serializedTxs,
		FeeReceiver:       feeReceiverIndex,
		CombinedSignature: *combinedSignature,
		PostStateRoot:     *stateRoot,
	}

	id, err := t.storage.AddCommitment(&commitment)
	if err != nil {
		return nil, err
	}

	commitment.ID = *id

	return &commitment, nil
}
