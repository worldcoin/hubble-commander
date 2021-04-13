package commander

import (
	"log"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
)

func createCommitments(
	pendingTransactions []models.Transaction,
	storage *st.Storage,
	cfg *config.RollupConfig,
) ([]models.Commitment, error) {
	stateTree := st.NewStateTree(storage)
	commitments := make([]models.Commitment, 0, 32)

	for {
		if len(commitments) >= int(cfg.MaxCommitmentsPerBatch) {
			break
		}

		initialStateRoot, err := stateTree.Root()
		if err != nil {
			return nil, err
		}

		includedTransactions, err := ApplyTransactions(storage, pendingTransactions, cfg)
		if err != nil {
			return nil, err
		}

		pendingTransactions = removeIncludedTransactionsFromPending(includedTransactions, pendingTransactions)

		if len(includedTransactions) < int(cfg.TxsPerCommitment) {
			err = stateTree.RevertTo(*initialStateRoot)
			if err != nil {
				return nil, err
			}
			break
		}

		log.Printf("Creating a commitment from %d transactions", len(includedTransactions))
		commitment, err := createAndStoreCommitment(storage, includedTransactions, cfg.FeeReceiverIndex)
		if err != nil {
			return nil, err
		}

		commitments = append(commitments, *commitment)

		err = markTransactionsAsIncluded(storage, includedTransactions, commitment.ID)
		if err != nil {
			return nil, err
		}
	}

	return commitments, nil
}

func removeIncludedTransactionsFromPending(includedTransactions, pendingTransactions []models.Transaction) []models.Transaction {
	output_index := 0
	for i := range pendingTransactions {
		tx := pendingTransactions[i]

		for j := range includedTransactions {
			includedTransaction := includedTransactions[j]
			if includedTransaction.Hash == tx.Hash {
				pendingTransactions[output_index] = tx
				output_index++
			}
		}
	}

	pendingTransactions = pendingTransactions[:]

	return pendingTransactions
}

func createAndStoreCommitment(storage *st.Storage, txs []models.Transaction, feeReceiverIndex uint32) (*models.Commitment, error) {
	combinedSignature := models.MakeSignature(1, 2) // TODO: Actually combine signatures

	serializedTxs, err := encoder.SerializeTransactions(txs)
	if err != nil {
		return nil, err
	}

	stateRoot, err := st.NewStateTree(storage).Root()
	if err != nil {
		return nil, err
	}

	commitment := models.Commitment{
		Transactions:      serializedTxs,
		FeeReceiver:       feeReceiverIndex,
		CombinedSignature: combinedSignature,
		PostStateRoot:     *stateRoot,
	}

	id, err := storage.AddCommitment(&commitment)
	if err != nil {
		return nil, err
	}

	commitment.ID = *id

	return &commitment, nil
}

func markTransactionsAsIncluded(storage *st.Storage, txs []models.Transaction, commitmentID int32) error {
	for i := range txs {
		err := storage.MarkTransactionAsIncluded(txs[i].Hash, commitmentID)
		if err != nil {
			return err
		}
	}
	return nil
}
