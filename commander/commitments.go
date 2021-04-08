package commander

import (
	"errors"
	"log"
	"time"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
)

var (
	ErrNotEnoughTransactions = NewCommitmentError("not enough transactions")
)

func CommitmentsEndlessLoop(storage *st.Storage, cfg *config.RollupConfig) error {
	done := make(chan bool)
	return CommitmentsLoop(storage, cfg, done)
}

func CommitmentsLoop(storage *st.Storage, cfg *config.RollupConfig, done <-chan bool) error {
	ticker := time.NewTicker(cfg.CommitmentLoopInterval)

	for {
		select {
		case <-done:
			ticker.Stop()
			return nil
		case <-ticker.C:
			err := CommitTransactions(storage, cfg)
			if err != nil {
				var e *CommitmentError
				if errors.As(err, &e) {
					log.Println(e.Error())
					continue
				}
				return err
			}
		}
	}
}

func CommitTransactions(storage *st.Storage, cfg *config.RollupConfig) (err error) {
	tx, txStorage, err := storage.BeginTransaction()
	if err != nil {
		return
	}
	defer tx.Rollback(&err)

	err = unsafeCommitTransactions(txStorage, cfg)
	if err != nil {
		return
	}

	return tx.Commit()
}

func unsafeCommitTransactions(storage *st.Storage, cfg *config.RollupConfig) error {
	txs, err := storage.GetPendingTransactions()
	if err != nil {
		return err
	}

	txsCount := uint32(len(txs))
	log.Printf("%d transactions in the pool", txsCount)

	if txsCount < cfg.TxsPerCommitment {
		return ErrNotEnoughTransactions
	}

	log.Printf("Applying %d transactions", txsCount)
	includedTxs, err := ApplyTransactions(storage, txs, cfg)
	if err != nil {
		return err
	}

	txsCount = uint32(len(includedTxs))
	if txsCount != cfg.TxsPerCommitment {
		return ErrNotEnoughTransactions
	}

	log.Printf("Creating a commitment from %d transactions", len(includedTxs))
	commitment, err := createAndStoreCommitment(storage, includedTxs, cfg.FeeReceiverIndex)
	if err != nil {
		return err
	}

	err = markTransactionsAsIncluded(storage, includedTxs, commitment.ID)
	if err != nil {
		return err
	}
	return nil
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
