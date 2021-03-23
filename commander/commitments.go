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
	ticker := time.NewTicker(500 * time.Millisecond) // TODO take from config

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
		return err
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
	commitmentID, err := createAndStoreCommitment(storage, includedTxs, cfg.FeeReceiverIndex)
	if err != nil {
		return err
	}

	err = markTransactionsAsCommitted(storage, includedTxs, *commitmentID)
	if err != nil {
		return err
	}
	return nil
}

func createAndStoreCommitment(storage *st.Storage, txs []models.Transaction, feeReceiverIndex uint32) (*int32, error) {
	combinedSignature := models.Signature{models.MakeUint256(1), models.MakeUint256(2)} // TODO: Actually combine signatures

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
	return storage.AddCommitment(&commitment)
}

func markTransactionsAsCommitted(storage *st.Storage, txs []models.Transaction, commitmentID int32) error {
	for i := range txs {
		err := storage.MarkTransactionAsIncluded(txs[i].Hash, commitmentID)
		if err != nil {
			return err
		}
	}
	return nil
}
