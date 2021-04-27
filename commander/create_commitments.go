package commander

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
)

func createAndStoreCommitment(
	storage *st.Storage,
	txType txtype.TransactionType,
	feeReceiverIndex uint32,
	serializedTxs []byte,
	combinedSignature *models.Signature,
) (*models.Commitment, error) {
	stateRoot, err := st.NewStateTree(storage).Root()
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

	id, err := storage.AddCommitment(&commitment)
	if err != nil {
		return nil, err
	}

	commitment.ID = *id

	return &commitment, nil
}

func markTransactionsAsIncluded(storage *st.Storage, transactions []models.TransactionBase, commitmentID int32) error {
	for i := range transactions {
		err := storage.MarkTransactionAsIncluded(transactions[i].Hash, commitmentID)
		if err != nil {
			return err
		}
	}
	return nil
}
